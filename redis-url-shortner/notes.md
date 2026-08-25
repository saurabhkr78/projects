Without Redis

You might have:

Client
  ↓
GET /aB92x
  ↓
Application
  ↓
PostgreSQL
  ↓
SELECT original_url
  ↓
Redirect

Every click potentially hits your database.

Imagine a popular short URL gets:

100,000 requests

You could end up doing a huge number of DB reads.

With Redis
100,000 requests
       ↓
      Redis
       ↓
GET aB92x
       ↓
original URL
       ↓
redirect

Redis is extremely good for this because URL mappings are:

frequently read
simple key → value data
often cacheable
relatively small
But there's an important distinction

A production URL shortener would usually still have a persistent database:

             POST /shorten
                  │
                  ▼
              PostgreSQL
                  │
                  └── persistent mapping
                       │
                       ▼
                    Redis
                  (cache)

Creation:

Long URL
   ↓
Generate short ID
   ↓
Store in DB
   ↓
Store/cache in Redis

Redirect:

Short ID
   ↓
Redis GET
   │
   ├── HIT → redirect
   │
   └── MISS
        ↓
       DB
        ↓
      Redis SET
        ↓
      redirect

That's the classic cache-aside pattern.

So URL shortener would teach you something your counter/rate-limiter/OTP projects haven't:

Redis can act as a high-speed cache in front of persistent storage.

# so in our projec the flow is
GET /aB92x
     │
     ▼
Redis GET("aB92x")
     │
     ├── HIT ──────────────→ redirect
     │
     └── MISS
          │
          ▼
       Fake DB
          │
          ▼
       Found URL
          │
          ▼
      Redis SET
          │
          ▼
       redirect


# you can hide the actual Redis library from your service by defining the methods in client itself which is  going to be used in service layer after dependency injection e.g

Add:

func (c *Client) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return c.RedisClient.Get(ctx, key).Result()
}

func (c *Client) Set(
	ctx context.Context,
	key string,
	value string,
) error {
	return c.RedisClient.Set(ctx, key, value, 0).Err()
}