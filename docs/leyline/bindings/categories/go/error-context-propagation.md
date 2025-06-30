---
id: error-context-propagation
last_modified: '2025-06-02'
version: '0.1.0'
derived_from: fix-broken-windows
enforced_by: 'Go error handling patterns, code review, static analysis'
---
# Binding: Propagate Error Context Through Go's Error Wrapping Patterns

Use Go's error wrapping mechanisms to preserve context and create clear error chains that enable effective debugging and monitoring. Properly wrapped errors maintain the full context of failure paths while providing actionable information for both developers and operations teams.

## Rationale

This binding implements our fix-broken-windows tenet by preventing the decay of error information quality over time. When errors lose context as they bubble up through the call stack, debugging becomes exponentially more difficult, leading to longer incident resolution times.

Go's error wrapping capabilities provide tools to maintain investigative trails. Without proper error context, debugging production issues becomes like trying to solve a mystery with only the final outcome and no clues about how you got there.

## Rule Definition

**Core Requirements:**

- **Error Wrapping**: Use `fmt.Errorf` with `%w` verb to preserve original errors while adding operational context
- **Context Preservation**: Include relevant identifiers (user IDs, request IDs, resource identifiers) without exposing sensitive information
- **Error Chain Inspection**: Use `errors.Is()` and `errors.As()` for type checking through multiple wrap layers
- **Sentinel Errors**: Define package-level sentinel errors with `errors.New()` for specific error conditions
- **Structured Error Types**: Create custom error types with structured data (codes, metadata, retry information)
- **Error Boundaries**: Translate technical errors to user-facing messages while preserving technical context for logging

**Essential Patterns:**
- Operation context: `"processing user order: %w"`
- Resource identification: `"failed to update user %s: %w"`
- Correlation IDs for request tracing
- Retry indicators for operational decisions

## Practical Implementation

**Comprehensive Error Context Pattern:**

```go
// Core error wrapping with operational context
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
    correlationID := getCorrelationID(ctx)

    order, err := s.orderRepo.GetOrder(ctx, orderID)
    if err != nil {
        return fmt.Errorf("failed to retrieve order %s (correlation_id: %s): %w",
                          orderID, correlationID, err)
    }

    if err := s.validateOrder(order); err != nil {
        return fmt.Errorf("validation failed for order %s (correlation_id: %s): %w",
                          orderID, correlationID, err)
    }

    if err := s.processPayment(ctx, order); err != nil {
        return fmt.Errorf("payment processing failed for order %s amount $%.2f (correlation_id: %s): %w",
                          orderID, order.Amount, correlationID, err)
    }

    return nil
}

// Structured error types with metadata
type PaymentError struct {
    Code      string
    Message   string
    Amount    float64
    Retryable bool
    Cause     error
}

func (e *PaymentError) Error() string {
    retryStatus := "non-retryable"
    if e.Retryable {
        retryStatus = "retryable"
    }
    return fmt.Sprintf("payment of $%.2f failed (code: %s, %s): %s",
                       e.Amount, e.Code, retryStatus, e.Message)
}

func (e *PaymentError) Unwrap() error {
    return e.Cause
}

type ValidationError struct {
    Field   string
    Message string
    Cause   error
}

func (e *ValidationError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("validation failed for field '%s': %s (caused by: %v)",
                          e.Field, e.Message, e.Cause)
    }
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
    return e.Cause
}

// Error boundary for user-facing responses
type ErrorBoundary struct {
    logger Logger
}

func (eb *ErrorBoundary) HandleError(ctx context.Context, w http.ResponseWriter, err error) {
    correlationID := getCorrelationID(ctx)

    // Log full technical error with context
    eb.logger.Error("request failed",
        "correlation_id", correlationID,
        "error", err.Error(),
        "error_chain", fmt.Sprintf("%+v", err))

    // Set correlation ID header for debugging
    w.Header().Set("X-Correlation-ID", correlationID)

    // Handle specific error types with appropriate responses
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        http.Error(w, fmt.Sprintf("Invalid %s: %s", validationErr.Field, validationErr.Message),
                   http.StatusBadRequest)
        return
    }

    var paymentErr *PaymentError
    if errors.As(err, &paymentErr) {
        if paymentErr.Code == "INSUFFICIENT_FUNDS" || paymentErr.Code == "CARD_DECLINED" {
            http.Error(w, paymentErr.Message, http.StatusPaymentRequired)
            return
        }

        // Retryable payment errors
        if paymentErr.Retryable {
            http.Error(w, "Payment processing temporarily unavailable",
                       http.StatusServiceUnavailable)
            return
        }
    }

    // Generic server error for unexpected errors
    http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
}

// Context propagation with correlation IDs
type ContextKey string

const CorrelationIDKey ContextKey = "correlation_id"

func getCorrelationID(ctx context.Context) string {
    if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
        return id
    }
    return "unknown"
}

// Middleware to inject correlation ID
func CorrelationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        correlationID := r.Header.Get("X-Correlation-ID")
        if correlationID == "" {
            correlationID = generateCorrelationID()
        }

        ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
        w.Header().Set("X-Correlation-ID", correlationID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Error chain inspection
func (s *Service) handleDatabaseError(err error) error {
    // Check for specific sentinel errors
    if errors.Is(err, sql.ErrNoRows) {
        return &NotFoundError{Resource: "user", Cause: err}
    }

    // Check for custom error types
    var timeoutErr *TimeoutError
    if errors.As(err, &timeoutErr) {
        return &PaymentError{
            Code:      "TIMEOUT",
            Message:   "database operation timed out",
            Retryable: true,
            Cause:     err,
        }
    }

    return fmt.Errorf("database operation failed: %w", err)
}
```

## Examples

```go
// ❌ BAD: Losing error context and making debugging difficult
func (s *OrderService) ProcessOrder(orderID string) error {
    order, err := s.repo.GetOrder(orderID)
    if err != nil {
        return err // Lost: which order failed to load?
    }

    err = s.validateOrder(order)
    if err != nil {
        return errors.New("validation failed") // Lost: what validation? Original error?
    }

    err = s.processPayment(order)
    if err != nil {
        return err // Lost: payment processing context
    }

    return nil
}

// ✅ GOOD: Comprehensive error context propagation
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
    correlationID := getCorrelationID(ctx)

    order, err := s.repo.GetOrder(ctx, orderID)
    if err != nil {
        return fmt.Errorf("failed to retrieve order %s (correlation_id: %s): %w",
                          orderID, correlationID, err)
    }

    if err := s.validateOrder(order); err != nil {
        return fmt.Errorf("validation failed for order %s (correlation_id: %s): %w",
                          orderID, correlationID, err)
    }

    if err := s.processPayment(ctx, order); err != nil {
        return fmt.Errorf("payment processing failed for order %s amount $%.2f (correlation_id: %s): %w",
                          orderID, order.Amount, correlationID, err)
    }

    return nil
}

// ❌ BAD: Generic error handling without actionable information
func (s *PaymentService) ChargeCard(amount float64, cardToken string) error {
    resp, err := s.paymentAPI.Charge(amount, cardToken)
    if err != nil {
        return errors.New("payment failed") // No retry info, error codes, or context
    }

    if resp.Status != "success" {
        return errors.New("payment was declined") // No specific reason or handling guidance
    }

    return nil
}

// ✅ GOOD: Structured error handling with actionable information
func (s *PaymentService) ChargeCard(ctx context.Context, amount float64, cardToken string) error {
    correlationID := getCorrelationID(ctx)

    resp, err := s.paymentAPI.Charge(amount, cardToken)
    if err != nil {
        return &PaymentError{
            Code:      "API_ERROR",
            Message:   fmt.Sprintf("payment API call failed (correlation_id: %s)", correlationID),
            Amount:    amount,
            Retryable: true,
            Cause:     err,
        }
    }

    switch resp.Status {
    case "success":
        return nil
    case "insufficient_funds":
        return &PaymentError{
            Code:      "INSUFFICIENT_FUNDS",
            Message:   "card has insufficient funds",
            Amount:    amount,
            Retryable: false,
            Cause:     nil,
        }
    case "temporary_failure":
        return &PaymentError{
            Code:      "TEMPORARY_FAILURE",
            Message:   "temporary payment processing issue",
            Amount:    amount,
            Retryable: true,
            Cause:     nil,
        }
    default:
        return &PaymentError{
            Code:      "UNKNOWN_STATUS",
            Message:   fmt.Sprintf("unexpected payment status: %s", resp.Status),
            Amount:    amount,
            Retryable: false,
            Cause:     nil,
        }
    }
}

// Error handling with structured inspection
func (s *OrderService) handlePaymentError(err error) error {
    var paymentErr *PaymentError
    if errors.As(err, &paymentErr) {
        // Log structured payment error details
        s.logger.Error("payment processing failed",
            "error_code", paymentErr.Code,
            "amount", paymentErr.Amount,
            "retryable", paymentErr.Retryable)

        // Retry logic based on error metadata
        if paymentErr.Retryable {
            return s.schedulePaymentRetry(paymentErr)
        }

        return fmt.Errorf("payment failed: %s", paymentErr.Message)
    }

    return err
}
```

## Related Bindings

- [error-wrapping](../../docs/bindings/categories/go/error-wrapping.md): This binding builds on basic error wrapping patterns to create comprehensive error chains with operational context.

- [technical-debt-tracking](../../core/technical-debt-tracking.md): Poor error context is technical debt that compounds over time; proper error handling prevents debugging debt accumulation.

- [use-structured-logging](../../core/use-structured-logging.md): Error contexts provide data for structured logging, while correlation IDs link errors to request traces.

- [fix-broken-windows](../../tenets/fix-broken-windows.md): Properly wrapped errors preserve debugging capability and prevent system decay from poor error handling.
