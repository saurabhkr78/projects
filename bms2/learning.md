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


# Traditional CRUD (No Interface)

Usually you'll see something like:

type BookRepository struct {
	db *pgxpool.Pool
}

Then in main.go:

repo := &BookRepository{
	db: db,
}

service := NewService(repo)

And the service looks like:

type Service struct {
	repo *BookRepository
}

Notice this line:

repo *BookRepository

The service is tightly coupled to PostgreSQL.

Architecture
Service
    │
    ▼
BookRepository
    │
    ▼
PostgreSQL

Everything depends directly on PostgreSQL.

Our Current Project

Now look at what we're building.

type Repository interface {
	Create(...)
	GetByID(...)
	List(...)
	Update(...)
	Delete(...)
}

Service becomes:

type Service struct {
	repo Repository
}

Not

repo *PostgresRepository
Architecture
Service
    │
    ▼
Repository Interface
    ▲
    │
PostgresRepository

Now the service doesn't know PostgreSQL exists.

What's the practical difference?
Without Interface
type Service struct {
	repo *PostgresRepository
}

Tomorrow you decide to use MongoDB.

Now you must change:

repo *MongoRepository

Everywhere.

With Interface
type Service struct {
	repo Repository
}

Tomorrow:

repo := NewMongoRepository(client)

Service code?

Zero changes.

Previous CRUD
Service
   │
   ▼
Concrete PostgreSQL Repository
Current CRUD
Service
   │
   ▼
Repository Interface
   │
   ▼
Concrete PostgreSQL Repository

That single extra layer of abstraction is what makes the architecture more flexible and testable without changing the business logic.


# is there any diff when using struct and interface dependncy injection makes in main.go?
Without Interface

Suppose we don't have an interface.

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

Then main.go becomes

repo := book.NewRepository(db)
service := book.NewService(repo)

The type of repo is

*PostgresRepository
With Interface

Now our constructor returns

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{
		db: db,
	}
}

main.go is still

repo := book.NewRepository(db)
service := book.NewService(repo)

Looks identical.

But now the type of repo is

Repository

not

*PostgresRepository
That's the only visible difference

Look carefully.

Without interface

repo

↓

*PostgresRepository

With interface

repo

↓

Repository

The code you write in main.go is almost the same.

Why does this matter?

Suppose tomorrow you write

repo := book.NewMongoRepository(client)

What is the type?

Still

Repository

because

NewMongoRepository(...)

also returns

Repository

So

service := book.NewService(repo)

doesn't change.

Visualize it
Without Interface
main.go

repo
 │
 ▼
*PostgresRepository

↓

Service
With Interface

# two kinds of string in go?
1. Double Quotes ("")

These are called interpreted string literals.

Go interprets special characters inside them.

Example:

name := "Saurabh"
Escape characters work
fmt.Println("Hello\nWorld")

Output:

Hello
World

\n became a newline.

Another example:

fmt.Println("He said \"Hello\"")

Output:

He said "Hello"

Notice we had to escape the quotes.

2. Backticks (` `)

These are called raw string literals.

Go does not interpret escape characters.

Everything inside is taken literally.

Example:

fmt.Println(`Hello\nWorld`)

Output:

Hello\nWorld

Notice:

It printed \n literally.

Why do we use backticks for SQL?

Imagine writing SQL with double quotes.

query := "INSERT INTO books (title, author, isbn, published_at) VALUES ($1, $2, $3, $4)"

This works.

But for longer queries:

query := "SELECT id, title, author FROM books WHERE id = $1 ORDER BY created_at DESC"

It becomes hard to read.

With backticks:

query := `
INSERT INTO books
(title, author, isbn, published_at)
VALUES ($1, $2, $3, $4)
`

Much cleaner.

You can preserve formatting exactly as you would write the SQL.


# 1. Are we using raw SQL instead of an ORM?

Yes.

Our repository is using raw SQL through the pgx driver.

_, err := r.db.Exec(
	ctx,
	`
	INSERT INTO books
	(title, author, isbn)
	VALUES ($1, $2, $3)
	`,
	book.Title,
	book.Author,
	book.ISBN,
)

We write:

SQL
Parameters
Mapping

ourselves.

How would an ORM do the same thing?

Suppose we use GORM.

Create
db.Create(&book)

That's it.

Behind the scenes GORM generates:

INSERT INTO books
(title, author, isbn)
VALUES (...);
Get By ID

Instead of

err := r.db.QueryRow(...).Scan(...)

ORM

var book Book

err := db.First(&book, id).Error

Behind the scenes

SELECT *
FROM books
WHERE id = 1;

Then GORM automatically scans every field.

Get All

Instead of

rows, _ := r.db.Query(...)

ORM

var books []Book

err := db.Find(&books).Error

ORM internally

Executes SQL
Iterates rows
Calls Scan
Appends to slice

You don't see any of it.

Update

Instead of

UPDATE books
SET title=$1

You write

db.Save(&book)

or

db.Model(&book).Updates(book)
Delete

Instead of

DELETE FROM books
WHERE id=$1

ORM

db.Delete(&book)
What does ORM hide?

Let's compare.

pgx
Book

↓

SQL

↓

Exec()

↓

Database
ORM
Book

↓

Reflection

↓

Generated SQL

↓

Database

The ORM generates SQL for you.

Why many Go companies prefer pgx

Because they like seeing

SELECT ...

JOIN ...

WHERE ...

GROUP BY ...

instead of wondering what SQL the ORM generated.

2. What does rows store?

Excellent question.

Look at this.

rows, err := r.db.Query(...)

What is

rows

?

It is NOT

[]Book

It is NOT

Book

Think of it as a cursor.

Suppose PostgreSQL returned

id	title
1	Atomic Habits
2	Clean Code
3	Go Programming

The database doesn't immediately create

[]Book

Instead it says

"I have three rows."

and gives you

Cursor

↓

Row 1

↓

Row 2

↓

Row 3

That cursor is

rows

Then

rows.Next()

moves

Cursor

↓

Row 1

↓

Cursor

↓

Row 2

↓

Cursor

↓

Row 3

Then

rows.Scan(...)

copies the current row into your struct.

Why defer rows.Close()?

This is extremely important.

Imagine PostgreSQL.

PostgreSQL

↓

Opens Cursor

↓

Sends Rows

Until you call

rows.Close()

the database keeps that cursor open.

Open cursors consume resources.

Think of it like opening a file.

file, _ := os.Open(...)

Eventually you must

file.Close()

Same here.

rows.Close()

That's why

defer rows.Close()

is immediately written after checking the error.

Even if an error happens later,

Go guarantees

rows.Close()

will run.

Why rows.Err()?

Suppose

Book 1

Book 2

Database Connection Lost

Your loop finishes.

How do you know something went wrong?

rows.Err()

checks whether iteration itself failed.

pgx Cheat Sheet
Method	Returns	When to Use	Example
Exec()	CommandTag, error	INSERT, UPDATE, DELETE	Create Book
QueryRow()	Row	SELECT returning exactly one row	Get book by ID
Query()	Rows	SELECT returning multiple rows	List books
Scan()	error	Copy columns into Go variables	Convert DB row → Book
Ping()	error	Verify DB connection	Startup health check
Close()	void	Close connection pool	Shutdown application
Begin()	Tx, error	Start a transaction	Money transfer
Commit()	error	Save transaction	Complete order
Rollback()	error	Undo transaction	Payment failure
Connection Pool Methods

Since r.db is a *pgxpool.Pool, you also have:

Method	Purpose	Practical Example
Exec()	Run SQL with no returned rows	INSERT, UPDATE, DELETE
QueryRow()	Fetch one row	Get user by ID
Query()	Fetch multiple rows	List all books
Ping()	Check DB connectivity	Application startup
Acquire()	Manually borrow one connection from the pool	Rare advanced cases
Close()	Shut down the entire pool	Application shutdown
Rows Methods

rows itself provides methods:

Method	Purpose
Next()	Move to the next row
Scan()	Read the current row into Go variables
Err()	Check whether iteration failed
Close()	Release the cursor/resources

# the Repository is an ideal place for an interface because the service depends on it. The Service itself can remain a concrete BookService until there's a real need for multiple implementations or mocking at that boundary. This keeps the code simpler and more idiomatic.
"Accept interfaces, return structs." for service layer in general
When do companies create Service interfaces?

Only if they need another implementation.

Example:

BookService

▲
│
├── DefaultBookService

└── CachedBookService

or

BookService

▲
│
├── RealService

└── MockService

Otherwise, it's unnecessary abstraction.

This follows an important Go philosophy

Accept interfaces, return structs.

Look at our repository constructor.

Earlier I suggested:

func NewRepository(db *pgxpool.Pool) Repository

Many senior Go developers would actually write:

func NewRepository(db *pgxpool.Pool) *PostgresRepository

because they don't return interfaces either.

Instead, they let the consumer decide whether to store it as an interface.

For example:

repo := book.NewRepository(db)

service := book.NewService(repo)

Since NewService accepts

repo Repository

Go automatically converts:

*PostgresRepository

↓

Repository
This is actually the idiomatic Go approach

The proverb is:

Accept interfaces, return concrete types.

Meaning:

func NewRepository(...) *PostgresRepository
func NewService(repo Repository) *BookService

Notice:

Constructors return structs.
Consumers accept interfaces.

So, we dont abstract service layer methods


# when multiple implementation is needed then only abstract " so repo can be implemented by many db so methods are wrapped in interface for abstraction but service can be only implemented by book service so no interface"

# since our project is small so right now our service layer methods are thin as the project grows we maintain SRP and add business logic like 
Validating ISBN format.
Preventing duplicate ISBNs.
Rejecting future publication dates.
Checking permissions.
Enforcing business constraints.

# The handler is responsible for:

Receiving the HTTP request.
Reading the JSON body.
Reading path/query parameters.
Calling the service.
Returning an HTTP response.

# in every Go web framework (standard net/http, Gin, Fiber, Echo, Chi).

The flow is:

HTTP Request
      │
      ▼
Handler Function
      │
      ├── Read JSON
      ├── Validate JSON format
      ├── Call Service
      └── Write JSON Response


# 1. net/http (Standard Library)

What is it?

The HTTP server built into Go. Every major Go framework is built on top of it in some way.

Pros
No external dependency.
Stable and maintained with Go itself.
You understand HTTP fundamentals.
Minimal abstraction.
Easy to debug.
Excellent for interviews.
Cons
More boilerplate.
No built-in routing features like path parameters or middleware helpers.
You often combine it with another router (e.g. Chi).
Best for
Internal APIs
Microservices
Learning HTTP
Long-lived backend services
Companies that value simplicity

Example companies often rely heavily on the standard library (sometimes with a router layered on top):

Google
Cloudflare
Many Kubernetes ecosystem projects
2. Chi

What is it?

A lightweight router that builds directly on net/http.

net/http
     ▲
     │
    Chi
Pros
Very idiomatic Go.
Tiny abstraction.
Middleware support.
REST routing.
Compatible with every net/http middleware.
Easy to migrate away from.
Cons
Slightly more verbose than Gin/Fiber.
Fewer convenience helpers.
Best for
Production REST APIs
SaaS backends
Enterprise services
Teams that like idiomatic Go
Why companies like Chi

Because the request handler is still

func(w http.ResponseWriter, r *http.Request)

So you're never locked into Chi.

3. Gin

Probably the most popular Go web framework.

net/http
      ▲
      │
     Gin
Pros
Huge ecosystem.
Large community.
Built-in validation.
JSON helpers.
Middleware.
Easy routing.
Easy for beginners.
Cons
More abstraction than Chi.
Uses its own gin.Context.
More magic.
Some middleware only works with Gin.
Best for
CRUD APIs
Startup MVPs
Admin panels
General REST services

Example

Instead of

json.NewDecoder(r.Body).Decode(&req)

Gin gives

c.BindJSON(&req)

Less code.

4. Fiber

Fiber is inspired by Express.js.

Unlike Gin and Chi, Fiber does not use net/http directly in its request handling model; it is built on fasthttp, which is a different HTTP implementation.

Pros
Very fast.
Express-like API.
Low memory usage.
Easy for Node.js developers.
Cons
Not based on the standard net/http interfaces.
Some net/http middleware won't work directly.
Smaller ecosystem than Gin.
Best for
High-throughput APIs.
Simple services.
Teams already familiar with Express.js.
Performance

People often compare benchmarks.

Reality:

Database Query

↓

5–20 ms

HTTP Framework

↓

0.05 ms

The database usually dominates request time.

So unless you're serving hundreds of thousands of requests per second, framework choice rarely determines performance.


# Suppose your frontend sends:

POST /books HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer xyz

{
    "title": "Atomic Habits",
    "author": "James Clear"
}

Before your handler is called, net/http has already done a lot of work.

Browser / Postman
        │
        ▼
TCP Connection
        │
        ▼
Go HTTP Server
        │
        ▼
Parse HTTP Request
        │
        ▼
Create *http.Request
        │
        ▼
Call Your Handler

Notice something.

You don't create *http.Request.

Go creates it for you.

# The handler signature

Every Go HTTP handler looks like this.

func(w http.ResponseWriter, r *http.Request)

Why exactly these two?

Because HTTP is always

Client

↓

Request

↓

Server

↓

Response

↓

Client

The handler needs exactly two things.

The incoming request
A way to send the response

That's precisely what these parameters represent.

*http.Request

Think of it as

Everything the client sent to the server.

Imagine the client sends

POST /books?id=10 HTTP/1.1

Host: localhost:8080

Authorization: Bearer xyz

Content-Type: application/json

{
    "title":"Atomic Habits"
}

Go converts this into one huge struct r *http.Request 
this huge struct contains
http.Request

├── Method
├── URL
├── Header
├── Body
├── Context
├── Cookies
├── RemoteAddr
└── ... etc

# What can ResponseWriter do?

Three things.

1. Set Headers

Example

w.Header().Set(
    "Content-Type",
    "application/json",
)

Now the client knows

I'm receiving JSON.
2. Status Code

Example

w.WriteHeader(http.StatusCreated)

Sends

201 Created

Other examples

http.StatusOK

↓

200

http.StatusBadRequest

↓

400

http.StatusNotFound

↓

404

3. Response Body

Suppose

book := BookResponse{
    Title:"Atomic Habits",
}

We write

json.NewEncoder(w).Encode(book)

Go converts

BookResponse

↓

JSON

↓

Writes it into

w

↓

Client receives

{
    "title":"Atomic Habits"
}

# These two interfaces power almost every Go web framework

When you later see Gin:

func(c *gin.Context)

or Fiber:

func(c *fiber.Ctx)

remember that they're essentially combining and wrapping the ideas of:

*http.Request (incoming request)
http.ResponseWriter (outgoing response)


# Instead of returning only:

201 Created

I'd return the created resource (or at least its ID).

For example:

{
  "id": 1,
  "title": "Atomic Habits",
  "author": "James Clear",
  "isbn": "123456789"
}

That means changing the flow slightly:

Repository returns the created Book (or generated ID).
Service converts it to BookResponse.
Handler encodes it with:
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(response)

# in handler file 
You're probably noticing duplication.

For example,

http.Error(...)

appears many times.

So does

json.NewEncoder(...)

and

strconv.ParseInt(...)

This is exactly what happens in real projects.

The first version of a handler is often repetitive.

Later, we'll introduce:

helper functions for JSON responses,
centralized error handling,
middleware,
validation.

Those improvements will reduce the duplication while keeping the handlers easy to read.

# the working of routes.go file
Let's understand every line
Why ServeMux?
mux := http.NewServeMux()

Think of it as a routing table.

Initially

ServeMux

(empty)

After registering routes

ServeMux

POST /books        → CreateBook()

GET /books         → GetAllBooks()

GET /books/{id}    → GetBookByID()

PUT /books/{id}    → UpdateBook()

DELETE /books/{id} → DeleteBook()

So ServeMux is literally a map of routes.

Why do we pass *Handler?

Because we need access to its methods.

handler.CreateBook

is a method on

Handler

So we inject the handler.

Exactly like:

Repository

↓

Service

↓

Handler

↓

Routes

Every layer receives the dependency it needs.

What is this?
mux.HandleFunc(
	"POST /books",
	handler.CreateBook,
)

Read it in plain English:

"When an HTTP POST request comes to /books, execute handler.CreateBook."

Similarly,

mux.HandleFunc(
	"GET /books/{id}",
	handler.GetBookByID,
)

means:

"When an HTTP GET request comes to /books/{id}, execute handler.GetBookByID."


flow:
Browser / Postman
        │
        ▼
HTTP Request
        │
        ▼
ServeMux
        │
        ▼
Route Match
        │
        ▼
Handler
        │
        ▼
Service
        │
        ▼
Repository
        │
        ▼
PostgreSQL

# The responsibility of main.go file is only
ts only job is to create dependencies, wire them together, and start the server


# phase 2

# we are cleaning our handler with helper function or util functions
File 1 : internal/http/response.go
package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
Let's understand every line
Why data any?

Suppose

CreateBook returns

BookResponse

GetAll returns

[]BookResponse

Health returns

HealthResponse

We don't know beforehand.

So

data any

means

Accept any Go type.

Equivalent to old

interface{}
Why return error?

Because

json.NewEncoder(...).Encode(...)

can fail.

Example

Broken connection.

So

return err

lets the handler decide.


# File 2 : internal/http/error.go
package http

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(
	w http.ResponseWriter,
	status int,
	message string,
) {

	_ = WriteJSON(
		w,
		status,
		ErrorResponse{
			Error: message,
		},
	)
}

Instead of

http.Error(...)

client receives

{
    "error":"invalid request body"
}

Much more RESTful.

Why use JSON instead of http.Error()?

http.Error() returns plain text:

invalid request body

Most APIs today return structured JSON:

{
  "error": "invalid request body"
}

This is easier for frontend applications to parse and display.


One important principle:

# The frontend should never have to guess the response format.

Whether the request succeeds or fails, the JSON structure should always be predictable.

Our API Response Standard
Success Response
{
    "success": true,
    "data": { ... }
}

or

{
    "success": true,
    "data": [
        ...
    ]
}
Error Response
{
    "success": false,
    "error": {
        "message": "book not found"
    }
}

Notice

The frontend always checks

success

Then

success == true

↓

Read data

or

success == false

↓

Read error.message

No guessing.