---
derived_from: simplicity
enforced_by: code review & race detector
id: concurrency-patterns
last_modified: '2025-05-14'
version: '0.1.0'
---
# Binding: Use Goroutines and Channels Judiciously with Explicit Coordination

Implement Go concurrency using clear patterns that prevent leaks, race conditions, and
deadlocks. Use goroutines only when they genuinely simplify design or improve
performance, pass context for propagating cancellation, use channels with clear
ownership semantics, and apply synchronization primitives correctly when sharing memory.

## Rationale

This binding implements our simplicity tenet by preventing uncontrolled concurrency complexity. Each `go` statement creates potential non-determinism that can compound into systems with intermittent, hard-to-reproduce bugs.

Disciplined use of Go's concurrency features—context propagation, channel ownership, and proper synchronization—creates predictable, understandable concurrent systems that remain simple even while leveraging concurrent execution power.

## Rule Definition

**Core Requirements:**

- **Goroutine Lifecycle**: Every goroutine MUST have a clear termination strategy through context cancellation, done channels, or defined exit conditions. Use worker pools to bound concurrency.

- **Context Propagation**: Functions performing I/O or spawning goroutines MUST accept `context.Context` as the first parameter and respect cancellation by checking `ctx.Done()`.

- **Channel Ownership**: Exactly one goroutine owns and closes each channel. Use typed channel directions in function parameters and include cancellation cases in `select` statements.

- **Synchronization**: Use `sync.Mutex` for shared data, `sync.WaitGroup` for coordinating goroutines, and `sync.Once` for initialization. Always defer unlock statements.

- **Race Safety**: All code MUST pass `go test -race ./...` and properly synchronize shared data access.

## Practical Implementation

**Comprehensive Concurrency Pattern Demonstrating All Key Principles:**

```go
// Complete example showing goroutine lifecycle, context, channels, and synchronization
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Thread-safe cache with proper synchronization
type Cache struct {
    mu      sync.RWMutex
    entries map[string]string
}

func NewCache() *Cache {
    return &Cache{entries: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, found := c.entries[key]
    return value, found
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.entries[key] = value
}

// Worker pool with bounded concurrency and proper lifecycle
func ProcessItems(ctx context.Context, items []string, cache *Cache) error {
    const maxWorkers = 3
    sem := make(chan struct{}, maxWorkers) // Semaphore for concurrency control
    errCh := make(chan error, 1)           // First error channel

    var wg sync.WaitGroup

    // Process each item with bounded concurrency
    for i, item := range items {
        // Acquire semaphore or handle cancellation
        select {
        case sem <- struct{}{}:
            // Got semaphore slot, continue
        case <-ctx.Done():
            return ctx.Err()
        case err := <-errCh:
            return err // Return first error
        }

        wg.Add(1)
        go func(i int, item string) {
            defer wg.Done()
            defer func() { <-sem }() // Always release semaphore

            // Simulate work with context cancellation check
            select {
            case <-time.After(100 * time.Millisecond):
                // Work completed
                cache.Set(fmt.Sprintf("item_%d", i), item)
            case <-ctx.Done():
                return // Exit on cancellation
            }
        }(i, item)
    }

    // Wait for completion with timeout handling
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    case err := <-errCh:
        return err
    }
}

// Pipeline pattern with channel ownership
func CreatePipeline(ctx context.Context) <-chan string {
    // Generator owns and closes output channel
    output := make(chan string)

    go func() {
        defer close(output) // Generator closes channel

        for i := 0; i < 10; i++ {
            item := fmt.Sprintf("item_%d", i)

            select {
            case output <- item:
                // Item sent successfully
            case <-ctx.Done():
                return // Exit on cancellation
            }
        }
    }()

    return output // Return receive-only channel
}

// Main function demonstrating proper usage
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    cache := NewCache()

    // Create pipeline and collect items
    items := []string{}
    for item := range CreatePipeline(ctx) {
        items = append(items, item)
    }

    // Process items with worker pool
    if err := ProcessItems(ctx, items, cache); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Verify results
    for i := 0; i < len(items); i++ {
        if value, found := cache.Get(fmt.Sprintf("item_%d", i)); found {
            fmt.Printf("Cached: %s\n", value)
        }
    }
}

// Testing concurrent operations
func TestConcurrentAccess(t *testing.T) {
    cache := NewCache()
    var wg sync.WaitGroup

    // Start multiple readers and writers
    for i := 0; i < 10; i++ {
        wg.Add(2)

        // Writer goroutine
        go func(id int) {
            defer wg.Done()
            cache.Set(fmt.Sprintf("key_%d", id), fmt.Sprintf("value_%d", id))
        }(i)

        // Reader goroutine
        go func(id int) {
            defer wg.Done()
            _, _ = cache.Get(fmt.Sprintf("key_%d", id))
        }(i)
    }

    wg.Wait()
}
```

## Examples

```go
// ❌ BAD: Goroutine leak, race condition, unclear ownership
func BadConcurrency() {
    // Goroutine with no cancellation mechanism
    go func() {
        for {
            time.Sleep(1 * time.Second)
            doWork() // Runs forever, no way to stop
        }
    }()

    // Race condition on shared data
    counter := 0
    for i := 0; i < 100; i++ {
        go func() {
            counter++ // Unsynchronized access
        }()
    }

    // Channel with unclear ownership
    ch := make(chan int)
    go func() {
        ch <- 42 // Who closes this channel?
    }()
    result := <-ch // Potential deadlock if sender fails
}

// ✅ GOOD: Proper lifecycle, synchronization, and ownership
func GoodConcurrency(ctx context.Context) error {
    // Synchronized counter
    var mu sync.Mutex
    counter := 0

    // Worker with proper cancellation
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return // Clean shutdown
            case <-ticker.C:
                mu.Lock()
                counter++
                mu.Unlock()
            }
        }
    }()

    // Channel with clear ownership
    ch := make(chan int, 1)
    go func() {
        defer close(ch) // Generator owns and closes channel

        select {
        case ch <- 42:
        case <-ctx.Done():
        }
    }()

    // Safe consumption
    select {
    case result := <-ch:
        fmt.Printf("Result: %d\n", result)
    case <-ctx.Done():
        return ctx.Err()
    }

    return nil
}
```

## Related Bindings

- [error-wrapping](../../docs/bindings/categories/go/error-wrapping.md): Error handling in concurrent systems requires context propagation across goroutines for effective debugging.

- [interface-design](../../docs/bindings/categories/go/interface-design.md): Well-designed interfaces enable safe concurrent operations by defining clear component interaction contracts.

- [pure-functions](../../core/pure-functions.md): Pure functions are inherently thread-safe and reduce the need for complex synchronization.

- [immutable-by-default](../../core/immutable-by-default.md): Immutable data eliminates entire classes of concurrency bugs by preventing shared state modification.
