Race Condition — Rate Limiter

❌ Wrong approach: GET → CHECK → INCR

Request A                 Request B
   │                         │
   │ GET → 4                 │
   │                         │
   │                         │ GET → 4
   │                         │
   │ check: 4 < 5 ✓          │
   │                         │ check: 4 < 5 ✓
   │                         │
   │ INCR → 5                │
   │                         │
   │                         │ INCR → 6
   │                         │
   ▼                         ▼
 ALLOW                     ALLOW ❌

Problem:

Both requests read 4 before either increments it.

GET → 4
GET → 4

Both think:

4 < 5 → ALLOW

So both requests pass, even though only one should have been allowed.

✅ Better approach: INCR → CHECK
Request A              Request B
   │                      │
   │ INCR                 │
   ▼                      │
   5                      │
   │                      │ INCR
   │                      ▼
   │                      6
   │                      │
   ▼                      ▼
 ALLOW ✓                REJECT ❌

Because Redis INCR is atomic, each request gets a unique incremented value.

Request A → INCR → 5 → 5 > 5 ❌ → ALLOW
Request B → INCR → 6 → 6 > 5 ✅ → 429
Core idea
GET → CHECK → INCR
     ↑
   Race condition

INCR → CHECK
  ↑
Atomic operation

Remember: Redis INCR itself is atomic, but a sequence like GET → CHECK → INCR is not atomic as a whole.



# decoding trick
URL PATH
/users/123
     ↓
mux.Vars(req)

QUERY
/users?id=123
     ↓
req.URL.Query().Get("id")

BODY
{
    "id": 123
}
     ↓
json.NewDecoder(req.Body).Decode(...)