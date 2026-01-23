# 🔍 Context Usage Guide

Complete guide on how we use `context.Context` throughout the application for better observability, cancellation, and tracing.

---

## 📚 Why Context?

Using `context.Context` properly provides:

✅ **Request Tracing** - Track requests across services  
✅ **Cancellation** - Cancel operations when client disconnects  
✅ **Timeouts** - Set deadlines for operations  
✅ **Logging** - Include request metadata in logs  
✅ **Propagation** - Pass values between layers  

---

## 🎯 Our Context Values

We store these values in context throughout the application:

```go
// Request ID - unique identifier for each request
ctx = context.WithValue(ctx, "request_id", "req_abc123...")

// User ID - authenticated user making the request
ctx = context.WithValue(ctx, "user_id", int64(42))

// Trace ID - for distributed tracing (OpenTelemetry)
ctx = context.WithValue(ctx, "trace_id", "trace_xyz...")

// IP Address - client IP for auditing
ctx = context.WithValue(ctx, "client_ip", "192.168.1.1")
```

---

## 🔧 Implementation Examples

### 1. **Stripe Service** - Context with API Calls

**Before (WRONG):**
```go
func (s *StripeService) CreateStripeCustomer(
    _ context.Context,  // ❌ Ignored!
    userID int64,
    email string,
) (string, error) {
    params := &stripe.CustomerParams{
        Email: stripe.String(email),
    }
    return customer.New(params)
}
```

**After (CORRECT):**
```go
func (s *StripeService) CreateStripeCustomer(
    ctx context.Context,  // ✅ Used!
    userID int64,
    email string,
) (string, error) {
    // Log with context
    logWithContext(ctx, "Creating Stripe customer for user %d", userID)

    // Pass context to Stripe API
    params := &stripe.CustomerParams{
        Params: stripe.Params{
            Context: ctx,  // ✅ Enables cancellation & timeout
        },
        Email: stripe.String(email),
    }

    cust, err := customer.New(params)
    if err != nil {
        logErrorWithContext(ctx, err, "Failed to create customer")
        return "", err
    }

    logWithContext(ctx, "Customer %s created", cust.ID)
    return cust.ID, nil
}
```

### 2. **Logging with Context**

```go
// Helper function for logging with context values
func logWithContext(ctx context.Context, format string, args ...interface{}) {
    // Extract values from context
    requestID := ctx.Value("request_id")
    userID := ctx.Value("user_id")

    // Build prefix
    prefix := "[Service]"
    if requestID != nil {
        prefix += fmt.Sprintf(" [req:%v]", requestID)
    }
    if userID != nil {
        prefix += fmt.Sprintf(" [user:%v]", userID)
    }

    // Log with context info
    msg := fmt.Sprintf(format, args...)
    log.Printf("%s %s", prefix, msg)
}

// Example output:
// 2026/01/23 12:34:56 [Stripe] [req:abc123] [user:42] Creating customer
// 2026/01/23 12:34:57 [Stripe] [req:abc123] [user:42] Customer cus_xxx created
```

### 3. **Timeout with Context**

```go
// Set timeout for Stripe API call
func (s *StripeService) CreateSubscriptionWithTimeout(
    ctx context.Context,
    customerID string,
    priceID string,
) (string, error) {
    // Create context with 10-second timeout
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Stripe API will respect this timeout
    return s.CreateStripeSubscription(ctx, customerID, priceID)
}
```

### 4. **Cancellation with Context**

```go
// HTTP handler example
func (h *Handler) HandleCreateSubscription(w http.ResponseWriter, r *http.Request) {
    // Get context from request (includes cancellation)
    ctx := r.Context()

    // If client disconnects, ctx.Done() will be closed
    // and Stripe API call will be canceled

    subscriptionID, err := h.stripeService.CreateStripeSubscription(
        ctx,  // ✅ Stripe SDK respects cancellation
        customerID,
        priceID,
    )

    if err != nil {
        // Check if error is due to cancellation
        if ctx.Err() == context.Canceled {
            log.Println("Request canceled by client")
            return
        }
        // Handle other errors...
    }
}
```

---

## 🌐 Context Flow in Our Application

```
HTTP Request
    ↓
[API Gateway Middleware]
    ├─ Add request_id to context
    ├─ Add user_id from JWT to context
    ├─ Add trace_id for tracing
    ↓
[HTTP Handler]
    ↓
[Service Layer] ← Context passed here
    ├─ Stripe API (with context)
    ├─ Database queries (with context)
    ├─ gRPC calls (with context)
    ├─ Logging (using context values)
    ↓
[Repository Layer] ← Context passed here
    ├─ SQL queries (with context)
    ↓
[Database]
```

---

## 🎨 Middleware for Context Enrichment

```go
// Middleware to add request ID to context
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Generate unique request ID
        requestID := generateRequestID()

        // Add to context
        ctx := context.WithValue(r.Context(), "request_id", requestID)

        // Add to response header for client
        w.Header().Set("X-Request-ID", requestID)

        // Continue with enriched context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware to add user ID from JWT to context
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from JWT
        userID, err := extractUserIDFromJWT(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Add to context
        ctx := context.WithValue(r.Context(), "user_id", userID)

        // Continue with enriched context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 📊 Structured Logging with Context

```go
// Use structured logging library (e.g., zerolog, zap)
import "github.com/rs/zerolog/log"

func logWithContext(ctx context.Context, msg string) {
    logger := log.With().
        Str("request_id", getStringFromContext(ctx, "request_id")).
        Int64("user_id", getInt64FromContext(ctx, "user_id")).
        Str("trace_id", getStringFromContext(ctx, "trace_id")).
        Logger()

    logger.Info().Msg(msg)
}

// Example output (JSON):
// {"level":"info","request_id":"abc123","user_id":42,"trace_id":"xyz","time":"2026-01-23T12:34:56Z","message":"Creating Stripe customer"}
```

---

## 🔒 Best Practices

### ✅ DO:

1. **Always accept `context.Context` as first parameter**
   ```go
   func DoSomething(ctx context.Context, arg1 string, arg2 int) error
   ```

2. **Pass context to downstream calls**
   ```go
   result, err := s.repository.GetData(ctx, id)
   result2, err := s.externalAPI.FetchData(ctx, id)
   ```

3. **Use context for cancellation signals**
   ```go
   select {
   case <-ctx.Done():
       return ctx.Err()  // Canceled or timed out
   case result := <-ch:
       return result
   }
   ```

4. **Check context errors**
   ```go
   if err != nil {
       if ctx.Err() != nil {
           return fmt.Errorf("operation canceled: %w", ctx.Err())
       }
       return err
   }
   ```

### ❌ DON'T:

1. **Don't ignore context with `_`**
   ```go
   // ❌ WRONG
   func BadFunc(_ context.Context, id int64) error
   
   // ✅ CORRECT
   func GoodFunc(ctx context.Context, id int64) error
   ```

2. **Don't store context in structs**
   ```go
   // ❌ WRONG
   type Service struct {
       ctx context.Context  // Don't do this!
   }
   
   // ✅ CORRECT - pass as parameter
   func (s *Service) DoWork(ctx context.Context) error
   ```

3. **Don't pass `nil` context**
   ```go
   // ❌ WRONG
   result, err := service.DoWork(nil)
   
   // ✅ CORRECT
   result, err := service.DoWork(context.Background())
   ```

4. **Don't use context for optional parameters**
   ```go
   // ❌ WRONG - use context for config
   ctx = context.WithValue(ctx, "page_size", 50)
   
   // ✅ CORRECT - use explicit parameters
   func GetUsers(ctx context.Context, pageSize int) ([]User, error)
   ```

---

## 🧪 Testing with Context

```go
func TestCreateStripeCustomer(t *testing.T) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Add test values to context
    ctx = context.WithValue(ctx, "request_id", "test_req_123")
    ctx = context.WithValue(ctx, "user_id", int64(99))

    // Run test
    customerID, err := service.CreateStripeCustomer(ctx, 99, "test@example.com")
    
    assert.NoError(t, err)
    assert.NotEmpty(t, customerID)
}

func TestCancellation(t *testing.T) {
    // Create context that can be canceled
    ctx, cancel := context.WithCancel(context.Background())

    // Start operation in goroutine
    errCh := make(chan error)
    go func() {
        errCh <- service.LongRunningOperation(ctx)
    }()

    // Cancel after 100ms
    time.Sleep(100 * time.Millisecond)
    cancel()

    // Should return context.Canceled error
    err := <-errCh
    assert.ErrorIs(t, err, context.Canceled)
}
```

---

## 📈 Monitoring & Tracing

### OpenTelemetry Integration

```go
import "go.opentelemetry.io/otel"

func (s *StripeService) CreateStripeCustomer(
    ctx context.Context,
    userID int64,
    email string,
) (string, error) {
    // Start span for tracing
    ctx, span := otel.Tracer("stripe").Start(ctx, "CreateStripeCustomer")
    defer span.End()

    // Add attributes to span
    span.SetAttributes(
        attribute.Int64("user_id", userID),
        attribute.String("email", email),
    )

    // Stripe API call (context propagates trace info)
    customerID, err := s.createCustomerInStripe(ctx, userID, email)
    if err != nil {
        span.RecordError(err)
        return "", err
    }

    span.SetAttributes(attribute.String("customer_id", customerID))
    return customerID, nil
}
```

---

## 📚 Additional Resources

- [Go Context Package](https://pkg.go.dev/context)
- [Context Best Practices](https://go.dev/blog/context)
- [Stripe Go SDK Context Support](https://github.com/stripe/stripe-go)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)

---

**✅ Summary:**

Context is not just a "required parameter" - it's a powerful tool for:
- 🔍 Observability (logging, tracing)
- ⏱️ Cancellation & timeouts
- 🔗 Request propagation
- 📊 Monitoring

Always use it properly! 🚀
