---
id: dependency-injection-patterns
last_modified: '2025-06-03'
version: '0.1.0'
derived_from: orthogonality
enforced_by: 'Go interfaces, constructor functions, dependency injection frameworks, code review'
---

# Binding: Design Modular Systems Through Go's Interface-Based Dependency Injection

Use Go's interface system and constructor patterns to create loosely coupled, testable components that depend on abstractions rather than concrete implementations. Proper dependency injection enables component independence while maintaining clear boundaries and making testing straightforward.

## Rationale

This binding implements our orthogonality tenet by using Go's interface system to create truly independent components that can be developed, tested, and modified without affecting each other. Dependency injection is fundamental to building maintainable systems because it allows you to compose functionality from independent parts while keeping those parts loosely coupled.

Like building with standardized electrical components, each Go struct declares what interfaces it needs without caring about specific implementations. A LED light strip declares "I need 12V DC power" through its connector interface, but doesn't care whether that power comes from a wall adapter, battery pack, or solar panel. This flexibility allows you to swap power sources without rewiring, test with a bench power supply, and combine components in ways the original designers never imagined.

Without dependency injection, your components become hardwired like appliances permanently connected to specific outlets. Moving requires rewiring the house, testing requires access to exact outlets, and if outlets break, appliances become useless. Go's interface system provides standardized "plugs and sockets" that allow components to connect without hardwiring, enabling the flexibility and testability that orthogonal design requires.

## Rule Definition

Dependency injection patterns must establish these Go-specific practices:

- **Interface-Based Dependencies**: Define interfaces for all external dependencies and inject implementations through constructors. Components should depend on interfaces, not concrete types.
- **Constructor Injection**: Use constructor functions that accept dependencies as parameters, making all dependencies explicit and ensuring components are fully initialized when created.
- **Interface Segregation**: Design focused interfaces that represent single capabilities rather than large, monolithic interfaces. This allows components to depend only on the functionality they actually use.
- **Dependency Inversion**: High-level modules should not depend on low-level modules. Both should depend on abstractions (interfaces) that the high-level modules define.

Common patterns this binding requires:

- Small, focused interfaces with single responsibilities
- Constructor functions that accept all dependencies as parameters
- Dependency interfaces defined by consumers, not providers
- Clear separation between interface definitions and implementations
- Composition over inheritance through interface embedding

What this explicitly prohibits:

- Direct instantiation of dependencies within structs
- Large, monolithic interfaces with multiple responsibilities
- Circular dependencies between packages
- Global variables or singletons for managing dependencies
- Constructor functions that hide dependency requirements

## Practical Implementation

1. **Define Consumer-Driven Interfaces**: Create small, focused interfaces that
   represent only the capabilities each component actually needs.

   ```go
   // Consumer-driven interfaces for user service
   type UserRepository interface {
       GetByID(ctx context.Context, id string) (*User, error)
       Create(ctx context.Context, user *User) error
       Update(ctx context.Context, user *User) error
   }

   type NotificationSender interface {
       SendWelcomeEmail(ctx context.Context, user *User) error
   }

   type PasswordHasher interface {
       Hash(password string) (string, error)
       Verify(password, hash string) error
   }

   // Service depends only on interfaces it needs
   type UserService struct {
       repo   UserRepository
       notify NotificationSender
       hasher PasswordHasher
   }

   // Constructor injection makes all dependencies explicit
   func NewUserService(
       repo UserRepository,
       notify NotificationSender,
       hasher PasswordHasher,
   ) *UserService {
       return &UserService{
           repo:   repo,
           notify: notify,
           hasher: hasher,
       }
   }

   func (s *UserService) RegisterUser(ctx context.Context, email, password string) error {
       // Hash password
       hashedPassword, err := s.hasher.Hash(password)
       if err != nil {
           return fmt.Errorf("password hashing failed: %w", err)
       }

       // Create user
       user := &User{
           ID:       generateID(),
           Email:    email,
           Password: hashedPassword,
       }

       if err := s.repo.Create(ctx, user); err != nil {
           return fmt.Errorf("user creation failed: %w", err)
       }

       // Send welcome email
       if err := s.notify.SendWelcomeEmail(ctx, user); err != nil {
           // Log error but don't fail registration
           log.Printf("Failed to send welcome email: %v", err)
       }

       return nil
   }
   ```

2. **Implement Provider Structs**: Create concrete implementations that satisfy
   the interfaces defined by consumers.

   ```go
   // Database implementation of UserRepository
   type PostgresUserRepository struct {
       db *sql.DB
   }

   func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
       return &PostgresUserRepository{db: db}
   }

   func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
       var user User
       err := r.db.QueryRowContext(ctx,
           "SELECT id, email, password FROM users WHERE id = $1", id,
       ).Scan(&user.ID, &user.Email, &user.Password)

       if err == sql.ErrNoRows {
           return nil, nil
       }
       if err != nil {
           return nil, fmt.Errorf("database query failed: %w", err)
       }

       return &user, nil
   }

   func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
       _, err := r.db.ExecContext(ctx,
           "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
           user.ID, user.Email, user.Password,
       )
       if err != nil {
           return fmt.Errorf("database insert failed: %w", err)
       }
       return nil
   }

   // Email implementation of NotificationSender
   type EmailNotificationSender struct {
       client EmailClient
   }

   func NewEmailNotificationSender(client EmailClient) *EmailNotificationSender {
       return &EmailNotificationSender{client: client}
   }

   func (e *EmailNotificationSender) SendWelcomeEmail(ctx context.Context, user *User) error {
       return e.client.SendEmail(ctx, EmailMessage{
           To:      user.Email,
           Subject: "Welcome!",
           Body:    "Welcome to our service!",
       })
   }

   // Bcrypt implementation of PasswordHasher
   type BcryptPasswordHasher struct {
       cost int
   }

   func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
       return &BcryptPasswordHasher{cost: cost}
   }

   func (b *BcryptPasswordHasher) Hash(password string) (string, error) {
       hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
       if err != nil {
           return "", fmt.Errorf("bcrypt hash generation failed: %w", err)
       }
       return string(hash), nil
   }

   func (b *BcryptPasswordHasher) Verify(password, hash string) error {
       return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
   }
   ```

3. **Wire Dependencies in Main Function**: Compose the application by creating
   all dependencies and injecting them through constructors.

   ```go
   func main() {
       // Initialize infrastructure dependencies
       db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
       if err != nil {
           log.Fatal("Failed to connect to database:", err)
       }
       defer db.Close()

       emailClient := NewSMTPEmailClient(os.Getenv("SMTP_HOST"))

       // Create implementations
       userRepo := NewPostgresUserRepository(db)
       notificationSender := NewEmailNotificationSender(emailClient)
       passwordHasher := NewBcryptPasswordHasher(12)

       // Inject dependencies into service
       userService := NewUserService(userRepo, notificationSender, passwordHasher)

       // Create HTTP handler with service dependency
       userHandler := NewUserHandler(userService)

       // Start server
       http.Handle("/users", userHandler)
       log.Fatal(http.ListenAndServe(":8080", nil))
   }

   // HTTP handler depends on service interface
   type UserHandler struct {
       service UserRegistrar
   }

   type UserRegistrar interface {
       RegisterUser(ctx context.Context, email, password string) error
   }

   func NewUserHandler(service UserRegistrar) *UserHandler {
       return &UserHandler{service: service}
   }

   func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
       if r.Method != http.MethodPost {
           http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
           return
       }

       var req struct {
           Email    string `json:"email"`
           Password string `json:"password"`
       }

       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
           http.Error(w, "Invalid JSON", http.StatusBadRequest)
           return
       }

       if err := h.service.RegisterUser(r.Context(), req.Email, req.Password); err != nil {
           http.Error(w, "Registration failed", http.StatusInternalServerError)
           return
       }

       w.WriteHeader(http.StatusCreated)
   }
   ```

## Examples

```go
// ❌ BAD: Direct dependencies, untestable, tightly coupled
type UserService struct {
    // Direct dependency on concrete types
    db     *sql.DB
    mailer *smtp.Client
}

func NewUserService() *UserService {
    // Hidden dependencies, makes testing impossible
    db, _ := sql.Open("postgres", "production-db-url")
    mailer, _ := smtp.Dial("smtp.production.com:587")

    return &UserService{
        db:     db,
        mailer: mailer,
    }
}

func (s *UserService) CreateUser(email string) error {
    // Directly using concrete implementations
    _, err := s.db.Exec("INSERT INTO users (email) VALUES (?)", email)
    if err != nil {
        return err
    }

    // Tightly coupled to SMTP implementation
    msg := "Welcome!"
    return s.mailer.Mail(email, msg)
}
```

```go
// ✅ GOOD: Interface-based dependencies, testable, loosely coupled
type UserRepository interface {
    Create(ctx context.Context, user *User) error
}

type EmailSender interface {
    SendWelcome(ctx context.Context, email string) error
}

type UserService struct {
    repo   UserRepository
    sender EmailSender
}

// All dependencies explicitly declared and injected
func NewUserService(repo UserRepository, sender EmailSender) *UserService {
    return &UserService{
        repo:   repo,
        sender: sender,
    }
}

func (s *UserService) CreateUser(ctx context.Context, email string) error {
    user := &User{Email: email}

    if err := s.repo.Create(ctx, user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    if err := s.sender.SendWelcome(ctx, email); err != nil {
        // Log but don't fail the operation
        log.Printf("Failed to send welcome email: %v", err)
    }

    return nil
}

// Easy to test with mock implementations
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockSender := &MockEmailSender{}

    service := NewUserService(mockRepo, mockSender)

    err := service.CreateUser(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.True(t, mockRepo.CreateCalled)
    assert.True(t, mockSender.SendWelcomeCalled)
}
```

## Related Bindings

- [interface-design](../../docs/bindings/categories/go/interface-design.md): Dependency injection requires well-designed interfaces that follow the Interface Segregation Principle and represent single capabilities rather than large, monolithic contracts.

- [package-design](../../docs/bindings/categories/go/package-design.md): Effective dependency injection depends on proper package structure that avoids circular dependencies and places interface definitions in appropriate packages.

- [no-internal-mocking](../../core/no-internal-mocking.md): Dependency injection enables testing without internal mocking by making all external dependencies explicit and replaceable through interfaces.
