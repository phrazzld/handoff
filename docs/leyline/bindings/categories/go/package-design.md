---
derived_from: modularity
id: package-design
last_modified: '2025-05-14'
version: '0.1.0'
enforced_by: code review & style guides
---
# Binding: Organize Go Code Into Purpose-Driven Packages

Design Go packages as cohesive units with clear, focused responsibilities and
well-defined boundaries. Each package should have a single purpose, contain related
functionality, maintain high internal cohesion, and expose a minimal, well-documented
API that hides implementation details.

## Rationale

This binding implements our modularity tenet by establishing clear, well-defined boundaries for Go code organization. Package design is the primary mechanism for modularity in Go, determining how components interact and compose together.

Poor package boundaries create hidden dependencies and unexpected side effects. Purpose-driven packages allow developers to reason about one part of the system without holding the entire codebase in their head—turning complex systems into manageable pieces.

## Rule Definition

**Core Requirements:**

- **Single Purpose**: Each package has one well-defined purpose expressible in a short sentence
- **Clear Naming**: Concise, lowercase, single words describing package contents (avoid "factory", "manager", "util")
- **Standard Layout**: Follow Go project layout with `/cmd`, `/internal`, `/pkg` directories
- **Domain Organization**: Group by business domains/features, not technical layers (prefer `internal/user` over `internal/controllers`)
- **High Cohesion**: All code works together for unified purpose
- **Low Coupling**: Minimal dependencies on other packages
- **No Circular Dependencies**: Import graphs form directed acyclic graph (DAG)
- **Intentional APIs**: Export only types/functions related to package purpose; keep implementation details unexported
- **Interface Acceptance**: Prefer accepting interfaces and returning concrete types

## Practical Implementation

**Complete Package Organization Pattern:**

```go
// Standard Go project layout
project-root/
├── cmd/myapp/main.go        # Application entry point
├── internal/                # Private application code
│   ├── user/               # Domain-based organization
│   │   ├── user.go         # Core types and domain logic
│   │   ├── service.go      # Business operations
│   │   ├── repository.go   # Data access interface
│   │   └── handler.go      # HTTP handlers
│   ├── order/              # Order domain
│   │   ├── order.go
│   │   ├── service.go
│   │   └── repository.go
│   └── app/               # Application assembly
│       └── app.go         # Dependency injection
├── pkg/                   # Public library code (use sparingly)
├── go.mod
└── go.sum

// Domain package example: internal/order/order.go
package order

// Exported types - part of package API
type Order struct {
    ID         string
    CustomerID string
    Items      []Item
    Status     OrderStatus
    CreatedAt  time.Time
}

type Item struct {
    ProductID string
    Quantity  int
    Price     decimal.Decimal
}

type OrderStatus string

const (
    StatusPending   OrderStatus = "pending"
    StatusConfirmed OrderStatus = "confirmed"
    StatusShipped   OrderStatus = "shipped"
)

// Exported interfaces define package contracts
type Service interface {
    Create(ctx context.Context, customerID string, items []Item) (*Order, error)
    Get(ctx context.Context, id string) (*Order, error)
    Update(ctx context.Context, order *Order) error
}

type Repository interface {
    Store(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
    Update(ctx context.Context, order *Order) error
}

// Unexported implementation
type service struct {
    repo Repository
}

// Factory function for dependency injection
func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, customerID string, items []Item) (*Order, error) {
    order := &Order{
        ID:         generateID(),
        CustomerID: customerID,
        Items:      items,
        Status:     StatusPending,
        CreatedAt:  time.Now(),
    }

    if err := validateOrder(order); err != nil {
        return nil, fmt.Errorf("invalid order: %w", err)
    }

    return order, s.repo.Store(ctx, order)
}

// unexported helper functions
func generateID() string { /* implementation */ }
func validateOrder(order *Order) error { /* implementation */ }

// Application assembly: internal/app/app.go
package app

type App struct {
    UserService  user.Service
    OrderService order.Service
}

func NewApp(cfg Config) (*App, error) {
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("connecting to database: %w", err)
    }

    // Dependency injection with interfaces
    userRepo := user.NewRepository(db)
    userService := user.NewService(userRepo)

    orderRepo := order.NewRepository(db)
    orderService := order.NewService(orderRepo)

    return &App{
        UserService:  userService,
        OrderService: orderService,
    }, nil
}
```

## Examples

```go
// ❌ BAD: Mixed responsibilities and circular dependencies
package utils  // Generic name with mixed concerns

// User-related functions mixed with order functions
func ValidateUserEmail(email string) bool { /* ... */ }
func CalculateOrderTotal(items []Item) decimal.Decimal { /* ... */ }
func FormatTimestamp(t time.Time) string { /* ... */ }

// ❌ BAD: Circular dependency
package user
import "github.com/org/project/internal/payment"

func (u *User) CanMakePayment(amount decimal.Decimal) bool {
    return payment.CheckUserCredit(u.ID, amount)  // Creates cycle
}

// ✅ GOOD: Focused packages with clear boundaries
package user  // Single responsibility - user domain

// PaymentChecker interface breaks circular dependency
type PaymentChecker interface {
    CheckCredit(userID string, amount decimal.Decimal) bool
}

type Service struct {
    paymentChecker PaymentChecker  // Dependency injection
}

func NewService(paymentChecker PaymentChecker) *Service {
    return &Service{paymentChecker: paymentChecker}
}

func (s *Service) CanMakePayment(userID string, amount decimal.Decimal) bool {
    return s.paymentChecker.CheckCredit(userID, amount)
}

func ValidateEmail(email string) bool { /* user-specific validation */ }

// package order - separate domain with focused responsibility
package order

type Service struct {
    userService user.Service  // Clear dependency direction
}

func (s *Service) Create(ctx context.Context, customerID string, items []Item) (*Order, error) {
    // Implementation using userService
}

func CalculateTotal(items []Item) decimal.Decimal { /* order-specific logic */ }

// package payment - implements user interface without importing user
package payment

type Service struct {
    db Database
}

// Ensure Service implements user.PaymentChecker
var _ user.PaymentChecker = (*Service)(nil)

func (s *Service) CheckCredit(userID string, amount decimal.Decimal) bool {
    creditLimit, err := s.fetchCreditLimitFromDB(userID)
    return err == nil && creditLimit >= amount
}
```

## Related Bindings

- [modularity](../../tenets/modularity.md): This binding implements modularity principles through Go package design and clear boundaries.

- [dependency-inversion](../../core/dependency-inversion.md): Use interfaces in consumer packages to manage coupling and create flexible, testable systems.

- [package-structure](../python/package-structure.md): Python equivalent for domain-based organization with clear module boundaries.

- [interface-design](../../docs/bindings/categories/go/interface-design.md): Well-designed interfaces help define clean package boundaries and APIs.

- [hex-domain-purity](../../core/hex-domain-purity.md): Package organization supports hexagonal architecture by reflecting domain boundaries.
