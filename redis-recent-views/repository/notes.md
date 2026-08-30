# Idiomatic Go & Error Handling

## 1. Go mein Error Handling ka Basic Pattern

Go exceptions ke instead errors ko explicitly return karta hai.

```go
result, err := doSomething()

if err != nil {
    return err
}
```

Ye Go ka normal aur idiomatic pattern hai.

### Short form

```go
if err := doSomething(); err != nil {
    return err
}
```

`err` sirf `if` block ke scope mein available hota hai.

---

# 2. `.Result()` vs `.Err()`

Go Redis commands mein agar command ka result chahiye:

```go
length, err := r.client.LPush(ctx, key, productID).Result()
```

Agar sirf ye check karna hai ki operation successful hua ya nahi:

```go
if err := r.client.LPush(ctx, key, productID).Err(); err != nil {
    return err
}
```

### Rule

```text
Result chahiye?
    ↓
.Result()

Sirf success/failure check karna hai?
    ↓
.Err()
```

Example:

```go
// LPush returns new list length
length, err := r.client.LPush(ctx, key, productID).Result()
```

But:

```go
// LTrim returns "OK", which we don't need
if err := r.client.LTrim(ctx, key, 0, 9).Err(); err != nil {
    return err
}
```

---

# 3. Error ko Destroy Mat Karo

Avoid:

```go
return errors.New(err.Error())
```

Aur avoid:

```go
return errors.New("failed: " + err.Error())
```

Better:

```go
return fmt.Errorf("failed to perform operation: %w", err)
```

`%w` original error ko wrap karta hai.

---

# 4. `%w` vs `%v`

### `%w`

```go
fmt.Errorf("database query failed: %w", err)
```

Original error preserve hota hai.

Isliye baad mein:

```go
errors.Is(err, target)
errors.As(err, &target)
```

use kar sakte ho.

### `%v`

```go
fmt.Errorf("database query failed: %v", err)
```

Error ko sirf formatted value ki tarah include karta hai.

### Rule

Jab existing error ko wrap karke return kar rahe ho:

```text
Prefer %w
```

---

# 5. Error mein Context Add Karo

Bad:

```go
return err
```

Ye technically correct hai, but debugging mein context missing ho sakta hai.

Better:

```go
return fmt.Errorf(
    "failed to add recent view: %w",
    err,
)
```

Agar multiple layers hain:

```text
Redis
  ↓
Repository
  ↓
Service
  ↓
Handler
```

Error useful context carry kar sakta hai:

```text
failed to add recent view:
failed to update Redis:
connection refused
```

---

# 6. Error Messages Lowercase Rakho

Idiomatic Go error strings generally lowercase hote hain.

Prefer:

```go
errors.New("failed to add recent view")
```

Instead of:

```go
errors.New("Failed to Add Recent View")
```

Usually ending punctuation bhi avoid karo:

```go
"failed to add recent view"
```

not:

```go
"Failed to add recent view."
```

Reason: errors commonly wrap hote hain:

```text
failed to add recent view: redis connection refused
```

---

# 7. Error ko Unnecessarily Recreate Mat Karo

Bad:

```go
if err != nil {
    return errors.New(err.Error())
}
```

Better:

```go
if err != nil {
    return err
}
```

Aur agar context add karna hai:

```go
return fmt.Errorf("failed to add recent view: %w", err)
```

---

# 8. Log Kahan Karna Hai?

Har layer mein same error ko log mat karo.

Avoid:

```text
Repository → log
Service    → log
Handler    → log
```

Isse duplicate logs milenge.

Better:

```text
Repository
    ↓
return error
    ↓
Service
    ↓
return error
    ↓
Handler / Error Middleware
    ↓
log + HTTP response
```

### General principle

> Jis layer ke paas error ko properly handle karne ki responsibility hai, ideally wahi log kare.

Lower layers usually:

```go
return fmt.Errorf("operation failed: %w", err)
```

Higher layer:

```text
log error
+
safe response to client
```

---

# 9. Internal Error Client ko Mat Dikhao

Suppose Redis returns:

```text
dial tcp 127.0.0.1:6379:
connection refused
```

Client ko ye mat bhejo:

```json
{
    "error": "dial tcp 127.0.0.1:6379: connection refused"
}
```

Instead:

```json
{
    "error": "internal server error"
}
```

Server logs mein actual error:

```text
failed to add recent view:
dial tcp 127.0.0.1:6379:
connection refused
```

### Principle

```text
Internal error
      ↓
wrap/context
      ↓
log internally
      ↓
safe public error
```

---

# 10. `errors.Is`

`%w` ka important benefit hai `errors.Is()`.

Define:

```go
var ErrNotFound = errors.New("not found")
```

Lower layer:

```go
return fmt.Errorf(
    "wallet lookup failed: %w",
    ErrNotFound,
)
```

Higher layer:

```go
if errors.Is(err, ErrNotFound) {
    // handle as not found
}
```

Wrapping ke baad bhi original error identify ho sakta hai.

---

# 11. `errors.As`

`errors.As()` ka use specific error type ko retrieve karne ke liye hota hai.

Conceptually:

```go
var appErr *AppError

if errors.As(err, &appErr) {
    // access AppError
}
```

Useful when different layers specific error types return karte hain.

---

# 12. `redis.Nil` ko Correctly Handle Karo

Redis mein:

```go
value, err := r.client.Get(ctx, key).Result()
```

Agar key exist nahi karti:

```go
err == redis.Nil
```

Ye necessarily infrastructure failure nahi hai.

It's a:

> **Cache miss**

So:

```go
if err == redis.Nil {
    // cache miss → go to database
} else if err != nil {
    // actual Redis error
}
```

Important distinction:

```text
redis.Nil
   ↓
expected condition
   ↓
cache miss

other error
   ↓
actual Redis/infrastructure failure
```

---

# 13. Don't Confuse Errors with Normal Control Flow

Har unusual situation ko error nahi banana chahiye.

Example:

```text
Cache key doesn't exist
```

For your cache:

```go
if err == redis.Nil {
    // normal cache miss
}
```

This is expected behavior.

But:

```text
Redis connection refused
```

is a real error.

---

# 14. Your Recent View Example

Idiomatic version:

```go
func (r *RedisRepositoryImpl) AddRecentView(
    ctx context.Context,
    userID string,
    productID string,
) error {

    key := "recent_views:" + userID

    // Remove existing occurrence.
    if err := r.client.LRem(
        ctx,
        key,
        0,
        productID,
    ).Err(); err != nil {
        return fmt.Errorf(
            "failed to remove product from recent views: %w",
            err,
        )
    }

    // Add product to the front.
    if err := r.client.LPush(
        ctx,
        key,
        productID,
    ).Err(); err != nil {
        return fmt.Errorf(
            "failed to add product to recent views: %w",
            err,
        )
    }

    // Keep only the 10 most recent products.
    if err := r.client.LTrim(
        ctx,
        key,
        0,
        9,
    ).Err(); err != nil {
        return fmt.Errorf(
            "failed to trim recent views: %w",
            err,
        )
    }

    return nil
}
```

The important pattern is:

```text
Operation
   ↓
if err := ...; err != nil
   ↓
fmt.Errorf("what failed: %w", err)
   ↓
return
```

---

# 15. Don't Over-Compress Go Code

Idiomatic doesn't mean:

> "Code ko jitna chhota ho sake utna chhota karo."

Good:

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

Bad style:

```go
if e:=x();e!=nil{return fmt.Errorf("x:%w",e)}
```

Even though both compile.

### Go philosophy

> **Simple, explicit, readable and boring code is good code.**

---

# 16. Idiomatic Go Mental Model

When writing Go, think:

```text
Do operation
     ↓
Check error immediately
     ↓
Add useful context if needed
     ↓
Return
```

Example:

```go
if err := repository.Save(ctx, data); err != nil {
    return fmt.Errorf("failed to save wallet: %w", err)
}
```

Not:

```text
Do 10 things
     ↓
Check all errors at the end
```

---

# 17. Quick Cheat Sheet

```text
┌─────────────────────────────────────────────┐
│          IDIOMATIC GO ERROR HANDLING        │
├─────────────────────────────────────────────┤
│                                             │
│ Check immediately                           │
│ if err != nil { return err }                │
│                                             │
│ Need context?                               │
│ fmt.Errorf("operation failed: %w", err)     │
│                                             │
│ Need result?                                │
│ .Result()                                   │
│                                             │
│ Only need error?                            │
│ .Err()                                      │
│                                             │
│ Inspect wrapped error?                      │
│ errors.Is / errors.As                       │
│                                             │
│ Expected Redis cache miss?                 │
│ redis.Nil                                   │
│                                             │
│ Don't expose internal errors to clients.   │
│                                             │
│ Don't log the same error at every layer.   │
│                                             │
│ Keep code simple and readable.             │
│                                             │
└─────────────────────────────────────────────┘
```

## One-line rule to remember

> **Check errors immediately, preserve them with `%w`, add useful context, handle them at the appropriate layer, and keep internal details away from the client.**

