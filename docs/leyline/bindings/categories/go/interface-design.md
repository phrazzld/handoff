---
derived_from: testability
enforced_by: code review & linting
id: interface-design
last_modified: '2025-05-14'
version: '0.1.0'
---
# Binding: Design Small, Focused Interfaces in Consumer Packages

In Go, define interfaces where they are used (consumer packages), not where they are implemented. Keep interfaces small and focused, preferably with only 1-3 methods, and design them based on the specific behaviors required by the consuming code rather than an implementation's capabilities.

## Rationale

This binding directly implements our testability tenet by enabling clean, flexible testing without over-mocking. When interfaces are small and defined by consumers, they naturally facilitate testing by allowing easy substitution of implementations without forcing tight coupling to concrete types.

Small interfaces defined by consumers act as "adapters" that specify only the behavior needed in a particular context, rather than forcing implementations to conform to large, rigid contracts. This approach creates natural seams in your codebase where alternate implementations—including test doubles—can be easily inserted without complex mocking frameworks.

The difference between small, consumer-defined interfaces and large "header interfaces" becomes most apparent when requirements change. With the small interface approach, a change in one area of code affects only the immediate interfaces needed by that area, while the rest of the system remains stable. This principle is captured in the Go proverb: "The bigger the interface, the weaker the abstraction"—larger interfaces create more rigid, brittle connections between components.

## Rule Definition

This binding establishes specific requirements for Go interface design:

- **Consumer Ownership**: Interfaces MUST be defined in the package that uses the behavior, not the package that implements it. Define interfaces based on the specific behaviors needed by consuming code. Implementation packages should return concrete types, not interfaces.

- **Interface Size**: Interfaces MUST be small and focused. Prefer single-method interfaces where possible. Most interfaces should contain no more than 3 methods. Larger interfaces should be composed from smaller ones. Never create "header interfaces" that simply mirror the methods of a concrete type.

- **Interface Naming**: Interface names MUST clearly describe behaviors. Names should describe the behavior, not the implementation. Prefer the `-er` suffix for active interfaces (e.g., `Reader`, `Writer`, `Validator`). Avoid implementation-specific words in the name.

- **Interface Implementation**: Implementation of interfaces MUST be implicit. Use compile-time checks to verify interface implementation. Return concrete types from factories, not interfaces. Accept interfaces as parameters to provide flexibility.

- **Avoid Empty Interface**: Usage of the empty interface (`interface{}` or `any`) SHOULD be minimized. The empty interface conveys no behavioral requirements and bypasses type safety. Prefer generic types (Go 1.18+) for type-safe polymorphism where applicable.

## Practical Implementation

1. **Define Interfaces at the Point of Use**: Place interfaces in the package where they're consumed to create clear dependency relationships and enable easier testing.

2. **Verify Interface Compliance**: Use compile-time checks to ensure types implement interfaces with patterns like `var _ Interface = (*ConcreteType)(nil)`.

3. **Design Small, Composable Interfaces**: Break down complex behaviors into smaller interfaces that can be composed together when needed.

4. **Accept Interfaces, Return Concrete Types**: Design functions and methods to promote flexibility by accepting interfaces as parameters while returning concrete types.

5. **Avoid the Empty Interface When Possible**: Use generics for type-safe polymorphism instead of `interface{}` when you need flexibility with type safety.

## Examples

```go
// ✅ GOOD: Small interfaces defined based on specific use cases
// Package: internal/auth
type UserAuthenticator interface {
    FindByEmail(email string) (*User, error)
    UpdatePassword(id string, password string) error
}

// Package: internal/user/admin
type UserManager interface {
    FindByID(id string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id string) error
}

// Package: internal/rbac
type RoleBasedUserFinder interface {
    FindByRole(role string) ([]*User, error)
}
```

```go
// ✅ GOOD: Interface defined in consumer package
package notification

// EmailSender defines the contract for sending emails.
// This is defined here because the notification package consumes this behavior.
type EmailSender interface {
    SendEmail(to string, subject string, body string) error
}

// Service sends notifications through various channels.
type Service struct {
    emailSender EmailSender
}

// NewService creates a notification service with the given dependencies.
func NewService(emailSender EmailSender) *Service {
    return &Service{
        emailSender: emailSender,
    }
}

// Send delivers a notification via email.
func (s *Service) Send(to string, message string) error {
    return s.emailSender.SendEmail(to, "New Notification", message)
}
```

```go
// package: internal/email/service.go
package email

// Service implements email sending capability.
// Note: This doesn't reference the notification.EmailSender interface.
type Service struct {
    smtpServer string
    username   string
    password   string
}

// NewService creates an email service.
func NewService(smtpServer, username, password string) *Service {
    return &Service{
        smtpServer: smtpServer,
        username:   username,
        password:   password,
    }
}

// SendEmail sends an email through SMTP.
func (s *Service) SendEmail(to string, subject string, body string) error {
    // Implementation details...
    return nil
}
```

```go
// Testing with simple mock implementations
package notification

import (
    "testing"
    "myapp/internal/email"
)

// Compile-time check to ensure email.Service implements EmailSender
var _ EmailSender = (*email.Service)(nil)

// Simple test implementation
type mockEmailSender struct {
    sentEmails []struct {
        to      string
        subject string
        body    string
    }
}

func (m *mockEmailSender) SendEmail(to string, subject string, body string) error {
    m.sentEmails = append(m.sentEmails, struct {
        to      string
        subject string
        body    string
    }{to, subject, body})
    return nil
}

func TestNotificationService(t *testing.T) {
    // Arrange
    mock := &mockEmailSender{}
    service := NewService(mock)

    // Act
    err := service.Send("user@example.com", "Hello, World!")

    // Assert
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if len(mock.sentEmails) != 1 {
        t.Fatalf("expected 1 email, got %d", len(mock.sentEmails))
    }
    if mock.sentEmails[0].to != "user@example.com" {
        t.Errorf("wrong recipient, got %s", mock.sentEmails[0].to)
    }
}
```

```go
// ✅ GOOD: Small, focused interfaces that compose
// Reader abstracts reading operations.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer abstracts writing operations.
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Closer abstracts resource cleanup.
type Closer interface {
    Close() error
}

// ReadWriter combines reading and writing operations.
type ReadWriter interface {
    Reader
    Writer
}

// ReadWriteCloser adds resource cleanup to ReadWriter.
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

```go
// ✅ GOOD: Accept interfaces, return concrete types
// repository.go
package repository

// Database defines what the repository needs from a database connection.
type Database interface {
    Query(query string, args ...interface{}) (Rows, error)
    Exec(query string, args ...interface{}) (Result, error)
}

// Rows is a subset of sql.Rows functionality needed by repository.
type Rows interface {
    Scan(dest ...interface{}) error
    Next() bool
    Close() error
}

// UserRepository manages user data storage.
type UserRepository struct {
    db Database
}

// NewUserRepository creates a concrete UserRepository.
func NewUserRepository(db Database) *UserRepository {
    return &UserRepository{db: db}
}
```

```go
// ✅ GOOD: Using generics (Go 1.18+) for type safety
type Cache[T any] struct {
    items map[string]T
}

func NewCache[T any]() *Cache[T] {
    return &Cache[T]{
        items: make(map[string]T),
    }
}

func (c *Cache[T]) Set(key string, value T) {
    c.items[key] = value
}

func (c *Cache[T]) Get(key string) (T, bool) {
    value, exists := c.items[key]
    return value, exists
}

// Type-safe usage
userCache := NewCache[*User]()
userCache.Set("user", &User{Name: "Alice"})
user, exists := userCache.Get("user")
if exists {
    // No type assertion needed, user is already *User
    fmt.Println(user.Name)
}
```

```go
// ❌ BAD: Large "header interface" that mirrors a concrete implementation
type UserRepository interface {
    FindByID(id string) (*User, error)
    FindByEmail(email string) (*User, error)
    FindAll() ([]*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id string) error
    FindByUsername(username string) (*User, error)
    FindByRole(role string) ([]*User, error)
    UpdatePassword(id string, password string) error
    UpdateEmail(id string, email string) error
    // ... and many more methods
}

// ❌ BAD: Interface defined alongside implementation
package database

type Repository interface {
    Find(id string) (interface{}, error)
    Save(entity interface{}) error
}

type PostgresRepository struct {
    // ...
}
```

## Related Bindings

- [dependency-inversion](../../docs/bindings/core/dependency-inversion.md): Go's consumer-defined interfaces directly implement dependency inversion by having high-level modules define interfaces that low-level modules implement. This binding provides the concrete Go pattern for achieving dependency inversion, which creates more testable, modular code by decoupling components and focusing on behaviors rather than implementations.

- [type-hinting](../python/type-hinting.md): Python equivalent using explicit type hints and interfaces to define clear contracts between components.

- [go-package-design](../../docs/bindings/categories/go/package-design.md): Interface design and package design work together closely in Go. Well-designed interfaces define the boundaries between packages and enable clean dependency management. The consumer-defined interface approach helps prevent circular dependencies between packages.

- [testability](../../docs/tenets/testability.md): Small interfaces make testing dramatically easier, as they can be readily mocked with simple test implementations rather than complex mocking frameworks. This binding implements the "No Mocking Internal Components" guideline by providing natural seams for substituting test implementations at package boundaries.

- [hex-domain-purity](../../docs/bindings/core/hex-domain-purity.md): Go interfaces are a key enabler for hexagonal architecture by providing the adapters between the domain and external concerns. Domain logic can remain pure by depending on interfaces rather than concrete implementations of infrastructure components.
