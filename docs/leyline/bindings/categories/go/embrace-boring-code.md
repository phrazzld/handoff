---
id: embrace-boring-code
last_modified: '2025-06-17'
version: '0.1.0'
derived_from: simplicity
enforced_by: 'code review, go fmt, go vet, team standards'
---

# Binding: Embrace Boring Go Code

Write straightforward, explicit Go code that solves problems without cleverness. Leverage Go's intentional simplicity rather than fighting it with elaborate abstractions. Choose composition over inheritance, explicit error handling over hidden magic, and clear function signatures over flexible interfaces.

## Rationale

Go was designed with grug's philosophy in mind: simple, readable, boring code that gets the job done. The language deliberately omits features that enable clever-but-complex patterns (generics were added reluctantly, inheritance doesn't exist, reflection is discouraged). This isn't a limitation—it's a feature that prevents the complexity demon from taking hold.

Grug's wisdom aligns perfectly with Go's design: no magic, explicit over implicit, composition over inheritance, and simple error handling. When you fight these principles by trying to make Go "more sophisticated," you're working against both the language and your future maintainers. The most productive Go code is often the most boring Go code.

The complexity demon whispers lies specific to Go: "You need elaborate interfaces for testability." "This generic abstraction will make everything flexible." "Reflection makes this code more powerful." But Go thrives on directness. A simple struct with clear methods beats a complex interface hierarchy. Explicit error checking beats hidden exception magic. Straightforward composition beats inheritance trees.

## Rule Definition

**MUST** prefer explicit error handling over abstraction layers that hide errors.

**MUST** use struct composition instead of attempting inheritance-like patterns.

**SHOULD** keep interfaces small and focused (ideally 1-3 methods).

**SHOULD** choose clear, verbose code over clever, terse solutions.

**SHOULD** use Go's standard library patterns instead of reimplementing abstractions from other languages.

## Go's Natural Grug Alignment

### Go's Grug Alignment

**No Magic:** Explicit error handling, no hidden exceptions
**Composition:** Embed structs and interfaces, don't mimic inheritance
**Simplicity:** Standard library first, avoid complex abstractions

## Error Handling Patterns

**✅ Explicit Checking:** Check each error, wrap with context using `fmt.Errorf`
**✅ Clear Context:** Add meaningful context when wrapping errors
**✅ Simple Retry:** Explicit loop with backoff, no complex retry libraries
**❌ Error Hiding:** Don't create abstractions that hide Go's explicit error model

## Simple Patterns

**Small Interfaces:** 1-3 methods max, compose when needed
**Clear Signatures:** Explicit parameters, avoid `interface{}` flexibility
**Struct Config:** Use configuration structs instead of complex option patterns
**Standard Library:** Use built-in types and functions before creating custom abstractions

## Common Anti-Patterns

**❌ Premature Interfaces:** Start with concrete types, add interfaces when you have multiple implementations
**❌ Reflection Overuse:** Use explicit structs instead of reflection for configuration
**❌ Generic Everything:** Post-1.18, avoid making everything generic; use specific types for clarity
**❌ Framework Thinking:** Don't recreate inheritance or complex patterns from other languages

## Testing Patterns

**Table Tests:** Use standard table-driven test patterns
**Constructor Injection:** Pass dependencies through constructors, not frameworks
**Simple Helpers:** Create focused test helper functions, avoid complex test machinery
**Real Dependencies:** Use real dependencies in tests when possible, simple mocks when needed

## Standard Library Examples

**Built-in Types:** Use `map[string]*User`, `[]Event`, avoid custom collection wrappers
**HTTP Handlers:** Standard library HTTP is sufficient for most APIs
**Sorting/Processing:** Use `sort.Slice`, `strings`, `json` packages before external dependencies
**Concurrency:** Use channels and goroutines, not complex async libraries

## Success Indicators

### Code Quality Metrics

**Readability:** New team members understand code without extensive explanation
**Debugging:** Issues are easy to locate and reproduce
**Testing:** Tests are straightforward to write and maintain
**Performance:** Simple code often performs better than complex abstractions

### Team Productivity

**Faster Development:** Less time spent on architecture decisions
**Easier Maintenance:** Fewer surprises in existing code
**Confident Refactoring:** Simple code is safe to change
**Reduced Bugs:** Explicit code has fewer hidden edge cases

### Go-Specific Benefits

**Fast Compilation:** Simple code compiles quickly
**Clear Stack Traces:** Explicit code produces readable error traces
**Tool Integration:** `go fmt`, `go vet`, and linters work better with simple code
**Memory Efficiency:** Straightforward code has predictable memory usage

## Related Patterns

**Simplicity Above All:** Go's design philosophy naturally enforces simplicity principles.

**Explicit Over Implicit:** Go requires explicit error handling and dependency management.

**No Secret Suppression:** Go's explicit nature makes hidden complexity visible through compilation errors.
