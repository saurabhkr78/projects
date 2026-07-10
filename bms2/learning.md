# the clean architecture we are following implements the Single Responsibility Principle (SRP).
# The Config struct is basically the Go representation of your .env file.
For example, if your .env contains:

APP_ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/bookdb?sslmode=disable


then your Config struct will look like:

type Config struct {
	AppEnv     string
	Port       string
	DatabaseURL string
}
read it like 
.env file
        │
        ▼
Read values
        │
        ▼
Config struct
        │
        ▼
Used throughout the application

# Why do we convert .env into a struct?

Imagine you need the database URL.


Option 1 (Bad)

Every package reads the environment itself.

os.Getenv("DATABASE_URL")

Now your codebase has:

handler
    │
    ├── os.Getenv(...)
service
    │
    ├── os.Getenv(...)
database
    │
    ├── os.Getenv(...)
logger
    │
    ├── os.Getenv(...)

If you ever rename DATABASE_URL, you have to update it everywhere.

# Option 2 (Production)

Read the environment once during startup.

.env
   │
   ▼
LoadConfig()
   │
   ▼
Config Struct
   │
   ├── DatabaseURL
   ├── Port
   ├── LogLevel
   └── AppEnv

Now every package simply uses the Config object.

This has several benefits:

One place to load configuration.
One place to validate it.
Easier to test (you can create different Config values in tests).
No scattered os.Getenv() calls throughout the code.

# Our Config struct will eventually be passed to other packages:

main.go
   │
   ├── config.Load()
   │
   ▼
Config
   │
   ├── database.New(cfg)
   ├── logger.New(cfg)
   └── server.New(cfg)

This is called Dependency Injection,

# Responsibility of pkg/config

Our package should have only one responsibility:

Load the application configuration and validate it.

That's it.



# Every package exposes an API.

For example,

The database package may expose:

database.Connect()

The logger package may expose:

logger.New()

So what should the config package expose?

There are two common designs.

Option 1
cfg, err := config.Load()
Option 2
cfg, err := config.New()

Both work.

Which one should we choose?

I want Option 1.

Why?

Because we're loading configuration, not creating a brand new object from scratch.

The function name tells the story.

cfg, err := config.Load()

reads naturally.

# In Go:

Lowercase → private (only usable inside the same package)
Uppercase → public/exported (usable from other packages)

# THE DATABASE URL have the database name so no need to write in seperate field
postgres://postgres:password@localhost:5432/bookdb?sslmode=disable

# How it actually works

There are two separate things happening.

Step 1

We have a file:

PORT=8080
DATABASE_URL=postgres://...
APP_ENV=development

At this point, Go cannot read this file automatically.

It's just a text file.

Step 2

A library like:

github.com/joho/godotenv

reads the .env file.

Think of it as:

.env file
      │
      ▼
godotenv
Step 3

godotenv loads those values into the operating system's environment variables for the current process.

Now the OS knows:

PORT=8080

DATABASE_URL=...

APP_ENV=development
Step 4

Now your Go program can ask the OS:

os.Getenv("PORT")

or

os.Getenv("DATABASE_URL")


# Load function implementation
Load .env.
Read each environment variable.
Populate the Config struct.
Validate required values.
Return *Config.

# when to use getenv and lookupenv
Getenv
PORT exists?
        │
        ├── Yes → "8080"
        │
        └── No  → ""

You only get the value.

LookupEnv
PORT exists?
        │
        ├── Yes → ("8080", true)
        │
        └── No  → ("", false)

You get both:

the value
whether the variable exists

# follows two important principles:

Single Responsibility Principle (SRP) – the config package owns configuration.
Dependency Injection – other packages receive a *Config; they don't fetch environment variables themselves.

# Go Concepts Learned

✓ Exported vs Unexported

Load()

↓

Public

load()

↓

Private

------------------------

Struct Literal

cfg := &Config{
    Port: "8080",
}

------------------------

Pointer to Struct

cfg := &Config{}

cfg is a pointer to Config.

------------------------

errors.New()

Static error.

fmt.Errorf()

Formatted error.


# Application Starts
        │
        ▼
config.Load()
        │
        ├── godotenv.Load()
        │
        ├── os.Getenv()
        │
        ├── Create Config
        │
        ├── Validate()
        │
        ▼
Return Config
        │
        ▼
main.go
        │
        ├── database.New(cfg)
        ├── logger.New(cfg)
        └── server.New(cfg)

        # Let's imagine our project has grown.

# Our .env now contains:

APP_ENV=development
PORT=8080

DATABASE_URL=postgres://postgres:password@localhost:5432/bookdb?sslmode=disable

JWT_SECRET=super-secret-key

LOG_LEVEL=debug

READ_TIMEOUT=10s
WRITE_TIMEOUT=10s

REDIS_URL=localhost:6379

Without a helper, your code becomes repetitive:

cfg := &Config{
	AppEnv:       os.Getenv("APP_ENV"),
	Port:         os.Getenv("PORT"),
	DatabaseURL:  os.Getenv("DATABASE_URL"),
	JWTSecret:    os.Getenv("JWT_SECRET"),
	LogLevel:     os.Getenv("LOG_LEVEL"),
	ReadTimeout:  os.Getenv("READ_TIMEOUT"),
	WriteTimeout: os.Getenv("WRITE_TIMEOUT"),
	RedisURL:     os.Getenv("REDIS_URL"),
}

There's nothing wrong with this, but every line repeats:

os.Getenv(...)
Solution 1 (Most Common)

Create a helper.

func getEnv(key string) string {
	return os.Getenv(key)
}

Now your code becomes:

cfg := &Config{
	AppEnv:       getEnv("APP_ENV"),
	Port:         getEnv("PORT"),
	DatabaseURL:  getEnv("DATABASE_URL"),
	JWTSecret:    getEnv("JWT_SECRET"),
	LogLevel:     getEnv("LOG_LEVEL"),
	ReadTimeout:  getEnv("READ_TIMEOUT"),
	WriteTimeout: getEnv("WRITE_TIMEOUT"),
	RedisURL:     getEnv("REDIS_URL"),
}

Looks cleaner.

Solution 2 (Production Favorite ⭐)

Suppose some variables are mandatory.

We can create:

func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}

	return value, nil
}

Now

dbURL, err := mustGetEnv("DATABASE_URL")
if err != nil {
	return nil, err
}

This removes the need for a separate Validate() because the helper validates as it reads.

Solution 3 (Optional Variables)

Some variables are optional.

Example

LOG_LEVEL=

If it's missing, use a default.

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

Usage:

cfg := &Config{
	LogLevel: getEnv("LOG_LEVEL", "info"),
}

Very common in production.

Solution 4 (Our Future Project ⭐⭐⭐⭐⭐)

Eventually our config.go might look like this:

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

Notice something:

Every line reads almost like English.

# each package should have one responsibility.
# We are designing packages with a single responsibility."

Because SRP applies to more than just functions:

Functions should have one responsibility.
Structs should have one responsibility.
Packages should have one responsibility.

# context.Context carries request-scoped information 

# There are three different things:

Database (PostgreSQL, MySQL, MongoDB)
Driver/Client (pgx, MongoDB driver, Redis client)
Connection Pool (managed by the driver/client)

# "One request = One database connection."

That's not true.

The correct mental model is:

1000 HTTP Requests

↓

One Shared Connection Pool

↓

20 Actual Database Connections

The pool decides which request gets which connection and when it's returned.

Using MySQL you might see:

db, err := sql.Open(...) already maintains a connection pool.
MongoDB driver:

client, err := mongo.Connect(...)

That client also manages a pool of connections.


# This syntax
func (cfg *Config) Connect() (*pgxpool.Pool, error) {
}

means:

Connect is a method of the Config type.

Read it like English:

"A Config object can connect."

Example:

cfg, _ := config.Load()

db, err := cfg.Connect()

Notice you're calling the method on the object.

Compare it with a package function
func Connect(cfg *Config) (*pgxpool.Pool, error) {
}

Usage:

cfg, _ := config.Load()

db, err := database.Connect(cfg)

Now you're calling a package function, not a method.


# Version 1 (Current)
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

Question:

Who created this context?

context.Background()

Nobody.

The database package created it.

Why is this not ideal?

Suppose tomorrow your application receives

Ctrl + C

or Kubernetes says

Stop the application.

Who should tell every package to stop?

The database package?

No.

The config package?

No.

The main application.

Who owns the application?

Always remember this.

Application

↓

main.go

Everything starts from

func main()

Therefore,

main owns the application's lifecycle.

# Why doesn't the entity have JSON tags?

For example:

Title string `json:"title"`

Answer:

Because this is our domain entity.

It shouldn't know anything about HTTP or JSON.

Those belong to the DTO layer.

This is one of the benefits of Clean Architecture.

Why no db tags?

For example:

Title string `db:"title"`

Again,

the entity shouldn't know about PostgreSQL.

The repository is responsible for mapping database rows to the entity.

Why no validation?

For example:

Title string `validate:"required"`

Because validation belongs to the service layer (business rules) or the request DTO (input validation), not the entity itself.

# Different DTO structs are used for different request and response payloads.
package book

import "time"

type CreateBookRequest struct { to the server
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
}

type UpdateBookRequest struct { to the server
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
}

type BookResponse struct { from server to the client
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

Each API operation has a different data contract. A CreateBookRequest contains only the fields needed to create a book, an UpdateBookRequest contains only editable fields, and a BookResponse contains the data the server wants to expose to the client. This keeps the API contract separate from the domain entity and prevents exposing or accepting fields that shouldn't cross the API boundary.

# The repository is responsible for persisting and retrieving Book entities.
storing and retreiving the book that's all

# The repository should be an interface, not a struct.

Why?

Because today we use:

PostgreSQL

Tomorrow we might use:

MongoDB

or

Mock Repository (Tests)

The service shouldn't care.

# Imagine this.

Today

Service

↓

PostgreSQL

Tomorrow

Service

↓

MongoDB

Should the service change?

No.

The service only knows:

type Repository interface

The implementation can change.

--The service depends on the interface, not PostgreSQL.

# Why context.Context?

Every repository method accepts:

ctx context.Context

because every database operation should be cancellable and should respect request deadlines.

We'll use it like this later:

row := db.QueryRow(ctx, ...)

For now, just remember:

Every operation that touches the database accepts a context.Context.

# Why is this called Loose Coupling?

Imagine we didn't use an interface.

Service would do:

type Service struct {
	repo *PostgresRepository
}

Now if tomorrow we move to MongoDB...

Everything breaks.

You have to change

*PostgresRepository

↓

*MongoRepository

Everywhere.

With an interface

type Service struct {
	repo Repository
}

Now Service doesn't care.

You can inject

PostgresRepository

or

MongoRepository

or

MockRepository

without changing the service.

That's loose coupling.


# This is the biggest reason interfaces exist.

Suppose you're testing your service.

Do you really want your tests to connect to PostgreSQL?

❌ No.

Instead:

type MockRepository struct{}

implements

Create()
GetByID()
List()
Update()
Delete()

Now your tests run without a real database.

This is one of the biggest practical benefits of programming to an interface.

# note:The database-specific repository implements the interface.