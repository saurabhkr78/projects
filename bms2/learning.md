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

# panic
panic()

↓

Run defer()

↓

Print Stack Trace

↓

Program Ends


log.Fatal
log.Println()

↓

os.Exit(1)

↓

Program Ends

No defer.

No stack trace.

| Situation                             | Use                                             |
| ------------------------------------- | ----------------------------------------------- |
| Database connection failed at startup | `log.Fatal()` (or return the error from `main`) |
| HTTP request validation failed        | Return `400`, don't panic                       |
| SQL query failed                      | Return an error                                 |
| File not found                        | Return an error                                 |
| Programmer bug / impossible state     | `panic()`                                       |
| Library code                          | Return errors, don't call `log.Fatal()`         |

#The 3-Level Error Mantra
Can the program continue?

        │
        ├── YES
        │      ↓
        │   Return an error
        │
        └── NO
               │
               ├── Startup failed?
               │       ↓
               │    log.Fatal()
               │
               └── Programmer bug /
                   Impossible state?
                       ↓
                    panic()

## Rule 1 (99% of your code)

If the caller can handle it, return an error.

Examples:

return err

Examples:

User entered wrong password ✅
Book not found ✅
Database timeout ✅
Validation failed ✅
File doesn't exist ✅

Never panic here.

## Rule 2

If the application cannot even start, use log.Fatal().

Examples:

cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}
db, err := database.Connect(cfg)
if err != nil {
    log.Fatal(err)
}

Without config or DB,

Server cannot run.

Terminate.

## Rule 3

If "this should NEVER happen", use panic().

Examples:

switch state {
case Running:
case Stopped:
default:
    panic("unknown state")
}

or

panic("unreachable code")

or

panic("nil pointer where it should never be nil")

This is a programmer mistake, not a user mistake.


# higher order function 
function taking a function and returning another function?"
That's called a higher-order function, and it's the foundation of middleware.

# go magic
Go says:

"If a function has this signature..."

func(
    http.ResponseWriter,
    *http.Request,
)

"...I can automatically turn it into an http.Handler."

That's why this works:

mux.HandleFunc(
    "GET /books",
    handler.GetAllBooks,
)

No extra code needed.

# middleware:This "thing" between the client and the handler is called middleware.

This raises one question.

If middleware runs before the handler...

How does it call the handler afterward?

That's the entire purpose of this function:

func Logging(
    next http.Handler,
) http.Handler

Notice something strange.

It accepts a handler...

and returns another handler.

This is the heart of middleware.

# stuct can store state

# Why do you think Go made Handler an interface with just one method instead of simply using functions everywhere?
Go didn't create http.Handler to replace functions.

Go created http.Handler so that anything
(function, struct, router, middleware, proxy)
can handle an HTTP request through one common contract:

ServeHTTP().

Why http.Handler is an interface with one method?
1. To define one common contract

The HTTP server only needs to know one thing:

"Can you handle an HTTP request?"

So Go defines:

type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

Now the server doesn't care what the handler is.

It only cares that it has

ServeHTTP()
Example
type BookHandler struct{}

func (b *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Books")
}

or

type UserHandler struct{}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Users")
}

The server treats both the same.

HTTP Server

↓

ServeHTTP()
2. Interfaces allow different implementations

Without interfaces, the server would only accept functions.

But sometimes we want a struct.

Example:

type AuthMiddleware struct {
	secret string
}

The middleware needs to remember

JWT Secret

Functions cannot store state.

Structs can.

By implementing

ServeHTTP()

that struct becomes a handler.

3. Middleware becomes possible

Suppose every request must be logged.

Without interfaces

Client

↓

Handler

With interfaces

Client

↓

Logging Middleware

↓

Handler

Both Logging Middleware and Handler implement

ServeHTTP()

so they can wrap each other.

Example:

func Logging(next http.Handler) http.Handler

Notice

Input

↓

Handler

Output

↓

Handler

Everything speaks the same language.

4. Routers become possible

Your router itself is also a handler.

Example:

mux := http.NewServeMux()

Internally

type ServeMux struct {
	...
}

implements

ServeHTTP()

So

http.ListenAndServe(":8080", mux)

works because

ServeMux

↓

implements

↓

Handler

The server doesn't know it's a router.

5. Reverse proxies become possible

Even a reverse proxy can be a handler.

Example

Client

↓

Reverse Proxy

↓

Backend Server

The proxy implements

ServeHTTP()

The HTTP server doesn't know whether it's serving files, proxying requests, or routing to APIs.

6. Easy testing

Suppose

type MockHandler struct{}

implements

ServeHTTP()

Now during tests

server := httptest.NewServer(&MockHandler{})

No real application needed.

Interfaces make mocking easy.

7. Go likes small interfaces

Imagine if Go designed

type Handler interface {
	ServeHTTP()
	Log()
	Close()
	Validate()
}

Every handler would have to implement all four methods.

Ridiculous.

Instead

type Handler interface {
	ServeHTTP()
}

One responsibility.

One method.

Simple.

This follows Go's philosophy:

Small interfaces are easier to implement and compose.

8. Functions can still be used (HandlerFunc)

Go didn't remove functions.

Instead it created

type HandlerFunc func(http.ResponseWriter, *http.Request)

and gave it

func (f HandlerFunc) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	f(w, r)
}

So now

A normal function

func GetBooks(w http.ResponseWriter, r *http.Request)

automatically behaves like a handler.

Best of both worlds.

# Go `http.Handler` Interface — Method Sets, Named Types aur Middleware Pattern

---

## 1. Interface ka Basic Rule

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**Rule**: Jo bhi type `ServeHTTP(ResponseWriter, *Request)` naam ka method provide karega, wo `Handler` maana jayega.

Matlab Go mein interface satisfy karne ke liye explicitly "implements" likhne ki zarurat nahi — bas method signature match hona chahiye.

---

## 2. Approach 1 — Struct se Handler Interface Implement Karna

```go
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Hello from myHandler!")
}
```

`myHandler` struct ke paas exactly wahi method hai — `ServeHTTP(ResponseWriter, *Request)` — isliye `myHandler` ko `Handler` maan liya jayega.

```go
var h Handler = &myHandler{}   // ✅ kaam karega
```

Agar ye method na hota, to:

```go
var h Handler = &myHandler{}   // ❌ error: myHandler does not implement Handler
```

kyunki `myHandler` ke paas `ServeHTTP(ResponseWriter, *Request)` method nahi hoga.

---

## 3. Approach 2 — Function Type se Handler Interface Implement Karna

**Kab struct use karo, kab function type?**

- Agar sirf function chahiye, **extra data store karne ki zarurat nahi** → **function type** use karo.
- Agar function ke saath extra state/data bhi chahiye (jaise counter, cache, config) → **struct type** use karo, kyunki struct ke paas flexibility hoti hai apni memory rakhne ki.

```go
type handlerFunc func(ResponseWriter, *Request)

// Go mein hum kisi bhi type pe method attach kar sakte hain,
// including function types. Isliye handlerFunc pe bhi
// ServeHTTP method attach kar sakte hain.
func (f handlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

---

## 4. Concept Samajhne Ke Liye — Named Types Pe Method Attach Karna

Upar wale concept ko samajhne ke liye kuch simple examples:

### Example 1 — Struct

```go
type Person struct {
    Name string
}

func (p Person) Greet() {
    fmt.Println("Hello, my name is", p.Name)
}

p := Person{Name: "John"}   // p, Person struct ka instance hai, Name field "John" set kiya
p.Greet()                    // p ke paas Greet method ka access hai
```

Yahan `p` ke paas `Name` field hai, isliye `p.Greet()` kaam karega. Agar `Name` field nahi hota, to bhi `Greet()` method kaam karta — kyunki method struct type pe attach hai, field pe nahi. (Field yahan sirf method ke andar use ho raha hai.)

### Example 2 — Named type based on `string`

```go
type Name string

func (n Name) Greet() {
    fmt.Println("Hello, my name is", n)
}

n := Name("John")   // n, Name type ka instance hai, value "John" set kiya
n.Greet()            // n ke paas Greet method ka access hai
```

### Example 3 — Named type based on `int`
# Every middleware follows this pattern:

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Before

		next.ServeHTTP(w, r)

		// After
	})
}

Think of it like a parent calling a child:

Parent
│
├── Do work before
│
├── Call child
│
└── Do work after child returns

Every middleware is both:

a child of the middleware above it, and
a parent of the middleware below it.

That's the essence of middleware chaining.


# http,FuncHandler vvi question to understand middleware
Achha sawaal — chalo ek kadam peeche jaake dekhte hain ki asli dikkat kya thi jiske liye ye poora setup banaya gaya.

Dikkat #1: Middleware chahiye tha, bina har handler ko badle

Socho tumhare paas 20 handlers hain — GetBook, GetUser, DeleteOrder, etc. Sabpe tum logging, auth check, panic-recovery jaisi cheezein chahte ho. Do tarike hain:

Bura tarika: Har handler ke andar khud logging likho:

go
func GetBook(w http.ResponseWriter, r *http.Request) {
    log.Println("before")
    // asli logic
    log.Println("after")
}

20 handlers, 20 jagah same code copy-paste — repetitive, maintain karna mushkil, aur "logic" (asli kaam) aur "cross-cutting concern" (logging) mix ho gaye.

Chahiye tha: Ek generic wrapper jo kisi bhi handler ko le, use "cover" kar de, aur naya handler wapas de — asli handler ko touch kiye bina:

go
loggedGetBook := Logging(GetBook)
loggedGetUser  := Logging(GetUser)

Yehi decorator/middleware pattern hai. Iske liye ek common "shape" chahiye jisme har handler fit ho — taaki Logging function generically kaam kar sake, chahe andar koi bhi handler ho.

Dikkat #2: Go mein "common shape" ka matlab interface hota hai

net/http ne decide kiya ki wo shape hoga: kisi bhi cheez ke paas ServeHTTP(w, r) method ho. Isliye:

go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

Ab Logging function aise likh sakte hain:

go
func Logging(next Handler) Handler {
    // next ko wrap karke naya Handler return karo
}

Isse Logging ko pata hi nahi chalega andar GetBook hai ya GetUser ya khud koi aur wrapped middleware — bas usse ServeHTTP chahiye.

Dikkat #3: Lekin plain function ke paas method nahi hota

Yahi asli twist hai. GetBook sirf ek function hai:

go
func GetBook(w http.ResponseWriter, r *http.Request) { ... }

Iska type hai func(ResponseWriter, *Request) — ek unnamed function type. Go mein method sirf named types pe lag sakta hai. Toh GetBook ke paas ServeHTTP method laga hi nahi sakte directly. Matlab GetBook khud Handler interface satisfy nahi karta.

Toh dikkat ye ban gayi: tumhare paas simple functions hain (jo likhna easy hai), lekin interface ko method chahiye (jo sirf named types pe milta hai). In dono ko jodne wala kuch chahiye.

Solution: HandlerFunc — ek adapter
go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

Ye ek chhota sa "adapter" type hai jiska kaam hi ye hai: kisi bhi plain function ko Handler interface satisfy karne layak bana dena, bina us function ko struct banaye ya uska code change kiye. Bas ek line conversion:

go
http.HandlerFunc(GetBook)   // ab ye Handler hai
Toh total problem ek line mein

"Main chahta hoon ki plain functions (jo likhna simple hai) ko generic Handler-wrapping machinery (jo interface maangti hai) ke saath use kar sakoon — bina har function ko struct banaye, aur bina copy-paste logging likhe."

HandlerFunc + type conversion is poori problem ka solution hai — ek pul (bridge) jo "function" aur "interface" ke beech ka gap bharta hai, taaki middleware pattern (Logging(next)) kaam kar sake kisi bhi plain function handler ke upar.


```go
type age int

func (a age) eligible() bool {
    return a >= 18
}

a := age(20)
a.eligible()   // a variable check karega ki value 18 se bada hai ya nahi — true/false return karega
```

**Inference**: Go mein hum kisi bhi named type (chahe woh struct ho, string ho, int ho, ya function ho) pe method attach kar sakte hain. Bas type ka naam hona chahiye.

### Example 4 — Function type (`handlerFunc`) same logic follow karta hai

```go
type handlerFunc func(ResponseWriter, *Request)

func (f handlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

hf := handlerFunc(GetBook)
hf.ServeHTTP(w, r)   // ServeHTTP call hoga, aur andar GetBook(w, r) execute hoga
```

- `handlerFunc` type banaya jo `ServeHTTP(ResponseWriter, *Request)` method provide karta hai — isliye `handlerFunc` ko `Handler` maan liya jayega.
- Ab `handlerFunc` type ke instance ke paas `ServeHTTP` method ka access hai.
- Is instance ko `Handler` interface wale variable mein assign kar sakte hain, kyunki `handlerFunc` type ke paas required method hai.
- Jab bhi `ServeHTTP` method call hoga, `handlerFunc` instance ke andar jo function store hai — wahi execute hoga.

---

## 5. Ab Problem Aati Hai — 20 Handlers Ka Case

Socho mere paas ek application ke liye 20 handlers hain:

```go
func GetBook(w http.ResponseWriter, r *http.Request)    { ... }
func GetUser(w http.ResponseWriter, r *http.Request)    { ... }
func DeleteOrder(w http.ResponseWriter, r *http.Request){ ... }
// ... 17 more
```

Aur har handler par mujhe **logging, auth check, rate limiting, panic recovery** etc. apply karna hai.

### Naive approach

```go
func GetBook(w http.ResponseWriter, r *http.Request) {
    log.Println("before")
    // asli logic
    log.Println("after")
}
```

Har handler ke andar ye code likhna padega.

**Dikkat:**
- 20 jagah copy-paste — DRY violation
- "Business logic" aur "cross-cutting concern" mix ho jaate hain
- Kal agar logging format badalna ho, 20 jagah edit karna padega

### Better approach — Generic Wrapper

Iski jagah hum ek **generic wrapper** bana sakte hain jo har handler ke liye logging, auth check, rate limiting, panic recovery etc. apply kare — bina original handler ka code change kiye.

Wrapper ek handler leta hai, aur ek naya handler wapas karta hai — lekin us naye handler ke aage-peeche logging/auth/etc. execute hota hai.

Main mein call karne pe aisa dikhega:

```go
loggedGetBook := Logging(GetBook)
loggedGetUser  := Logging(GetUser)
```

Isko **decorator / middleware pattern** kehte hain.

Lekin `Logging` (aur baaki cross-cutting concerns) ko ye pata nahi hona chahiye ki andar konsa specific handler call ho raha hai.

---

## 6. Solution Design — Step by Step

**Observation 1**: Handler ek HTTP function hai jiska signature hai:

```go
func(http.ResponseWriter, *http.Request)
```

**Observation 2**: Mujhe ek aisi cheez chahiye jiske andar koi bhi handler fit ho jaaye.

**Kyun?**
Kyunki mujhe har handler ke around logging, auth check, rate limiting, panic recovery jaise cross-cutting concerns lagane hain — lekin `Logging()` ko ye nahi pata hona chahiye ki uske andar `GetBook` aaya hai ya `GetUser`. Usko sirf ek "handler" milna chahiye jise wo execute kar sake.

**Isliye pehle handler ko store karna hoga.**

Go mein function bhi ek value hota hai. Aur kisi bhi value ko store karne ke liye uska ek **type** chahiye.

Kyunki saare handlers ka signature same hai, hum un sabke liye ek common function type bana sakte hain:

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

Ab is `HandlerFunc` type ke variable mein `GetBook` bhi store ho sakta hai, `GetUser` bhi, `DeleteOrder` bhi:

```go
var h HandlerFunc

h = GetBook
// ya
h = GetUser
// ya
h = DeleteOrder
```

Ab `Logging()` ke paas ek generic handler aa gaya. Use ye jaanne ki zarurat hi nahi ki andar kaunsa handler stored hai. Jab request aayegi, `Logging()` bas us stored handler ko call kar dega:

```go
func Logging(h HandlerFunc) HandlerFunc {
    // ...
}
```

Aur jab actual request aayegi, `Logging` ke paas jo bhi handler store hai, usko bas call kar dega:

```go
h(w, r)
```

- Agar `h` ke andar `GetBook` stored hai, to ye `GetBook(w, r)` ban jayega.
- Agar `h` ke andar `GetUser` stored hai, to ye `GetUser(w, r)` ban jayega.

**Isliye `Logging` generic ban gayi.** Use fark hi nahi padta andar kaunsa handler hai — uska kaam sirf handler ke pehle aur baad mein cross-cutting concern execute karna hai.

---

## 7. Core Takeaways

1. **Interface satisfaction implicit hai** — jo type required method provide kare, wo automatically interface satisfy karta hai.
2. **Method sirf named types pe lag sakta hai** — struct, string-based type, int-based type, ya function-based type, sab pe method attach ho sakta hai, bashart type ka naam ho.
3. **Struct vs function type** — extra state chahiye to struct, sirf function wrap karna hai to function type.
4. **Function bhi ek value hai Go mein**, aur value ko store karne ke liye type chahiye — isi wajah se `HandlerFunc` type banaya jaata hai, taaki alag-alag handlers ko ek common type ke through generically store/pass kiya ja sake.
5. **Middleware/decorator pattern** isi generic storage aur interface satisfaction ke upar based hai — wrapper ko andar ke specific handler se koi matlab nahi, use sirf itna pata hai ki jo bhi stored hai use `(w, r)` ke saath call karna hai.






# Most beginners think this happens
book := h.service.GetByID(...)

return book

❌ Wrong.

Handlers don't return responses.

What actually happens

Your handler gets

func GetBook(
    w http.ResponseWriter,
    r *http.Request,
)

When you do

httphelper.WriteJSON(
    w,
    http.StatusOK,
    book,
)

inside WriteJSON

json.NewEncoder(w).Encode(book)

The encoder writes bytes into

w

Think of w as a network pipe.

Book Struct

↓

json.Marshal()

↓

{"id":1,"title":"Go"}

↓

ResponseWriter (w)

↓

TCP Socket

↓

Browser / Postman

Notice

The handler never says

return book

It says

Write it into w.
Then what is ServeHTTP() returning?

This surprises many people.

func (h *BookHandler) ServeHTTP(...) {
    ...
}

returns

void

Nothing.

Why?

Because the response has already been written.

By the time

next.ServeHTTP(w,r)

finishes,

the response is already flowing toward the client.

Then a very interesting problem appears...

Suppose your handler writes

httphelper.WriteJSON(
    w,
    http.StatusCreated,
    ...
)

or

httphelper.WriteError(
    w,
    http.StatusBadRequest,
    ...
)

or

httphelper.WriteError(
    w,
    http.StatusInternalServerError,
    ...
)

Question:

How can Logging Middleware know whether the handler returned

200

or

201

or

404

or

500

?

It only has

w

But

http.ResponseWriter

has no method like

StatusCode()

😄

This is a design challenge.

This is where Go becomes really beautiful.

We solve it by creating our own ResponseWriter.

Something like

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

This is called embedding, and it's one of Go's most elegant features.

Our custom responseWriter behaves exactly like the original http.ResponseWriter, but it also remembers the status code.

Then the flow becomes:

Handler

↓

WriteHeader(201)

↓

Our ResponseWriter

↓

Stores

statusCode = 201

↓

Forwards to original ResponseWriter

↓

Client

Now Logging Middleware can print:

POST /books
Status: 201
Duration: 3ms



# mind-blowing moment

You already answered correctly:

Who sends the response?

Answer:

ResponseWriter

Now I'm going to ask a much harder question.

Suppose your handler does

httphelper.WriteJSON(
    w,
    http.StatusCreated,
    book,
)

Where is

201

stored?

Inside

w

But...

ResponseWriter doesn't expose it.

There is no

w.StatusCode()

method.

So...

How can Logging Middleware print

POST /books

201

3ms

?

That mystery leads to one of Go's most beautiful techniques:

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

You'll suddenly understand:

Embedding
Interface delegation
Method promotion
How frameworks like Gin and Chi capture status codes
Why Go doesn't need inheritance
# chainig of middleware works from inside out to call in reverse order using for loop

# what if need to pass a middleware to all the handler then 
1. we have multiple middleware so we chain them 
2.Every middleware has exactly the same signature.
e,g
func Logging(next http.Handler) http.Handler

func Recovery(next http.Handler) http.Handler

func CORS(next http.Handler) http.Handler

func RequestID(next http.Handler) http.Handler

3.So let's give this function type a name
exactly like type HandlerFunc func(http.ResponseWriter, *http.Request)
or we create type middleware func(http.RepsonseWriter ,*http.Request) a middleware is any function that accepts and  returns a handler
4.Now we can store them together in one slice
middlewares := []Middleware{
    Logging,
    Recovery,
    Authentication,
}
5.Now create Chain()
Instead of writing

Logging(
    Recovery(
        Authentication(handler),
    ),
)

we write

handler = Chain(
    handler,
    Logging,
    Recovery,
    Authentication,
)

Chain loops over the slice and wraps each middleware automatically.
6.Every middleware only knows one thing:

next.ServeHTTP(w, r)

It has no idea whether next is:

another middleware 
the final handler 
a router 

It simply forwards the request.

7. How different framework implement chaining of middleware
This is exactly how frameworks work

in Gin you write:

router.Use(
    gin.Logger(),
    gin.Recovery(),
)

In Chi:

r.Use(
    middleware.Logger,
    middleware.Recoverer,
)

In Fiber:

app.Use(logger.New())
app.Use(recover.New())

Internally, they are all building a middleware chain very similar to the Chain() function we'll eventually write.


# q: imagine thousands of req comes and some of them failed but dont know which failed beacuse no status code!!
aslo http.ResponseWriter does not have any function like where u can read the status

ans: If Go's ResponseWriter doesn't remember the status code for us... can we build our own ResponseWriter that does?
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}
When you understand why this works, you'll also understand:

Embedding
Method promotion
Interface composition
How Gin, Chi, and other frameworks capture HTTP status codes

# HandlerFunc is a new type whose underlying type is a function.
type HandlerFunc func(http.ResponseWriter , *http.Request)
Then we pass our anonymous function
func(w http.ResponseWriter, r *http.Request) {
    ...
}

Its signature matches exactly.

So Go says

"Okay, I can store this function inside a HandlerFunc."


# why do we call HandlerFunc method as adapter
every http handler looks similar in method parameter
so to attain common abstraction
go introduced
type Handler interface{
    ServeHTTP(http.ResponseWriter,*http.request)
}
the problem was a normal function cannnot implement ServeHTTP() because it 
so i does not satisfy  http.Handler interface

instead of forcing dev to write this every time

``type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    fmt.Fprintln(w, "Hello")
}
``
they created a new named function type

type HandlerFunc func (http.ResponseWriter,*http.request)

Since Go allows methods on named types, they attached:
func (f HandlerFunc) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request,
) {
    f(w, r)
}
Now any ordinary function can be converted into a HandlerFunc:

hf := http.HandlerFunc(Hello)

Since HandlerFunc has a ServeHTTP() method, it automatically satisfies:

type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}

As a result, the HTTP server can treat both:

Struct-based handlers
Function-based handlers

exactly the same.



# reciever object
type Book struct {
    Title string
}

func (b Book) Print() {
    fmt.Println(b.Title)
}

Question:

Who is

b

?

It is the object on which the method is called.



# How does a logging middleware know the HTTP status code?

Your current middleware probably looks like this:

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        next.ServeHTTP(w, r)

        log.Printf("%s %s took %v",
            r.Method,
            r.URL.Path,
            time.Since(start),
        )
    })
}

It can log:

✅ Method
✅ Path
✅ Time

But it cannot log:

❌ 200
❌ 404
❌ 500

Why?

Because http.ResponseWriter has no method like:

w.StatusCode()

So we'll build our own wrapper around it.

This is the next concept

We'll create something like:

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

At first glance, it looks like a tiny struct.

In reality, it teaches four major Go concepts:

Embedding
Method promotion
Method overriding
How decorators/wrappers work

Those four concepts appear everywhere in Go—not just in net/http.

# learn the concept of decorator/wrapping ,embedding

# decorator

what a Decorator (Wrapper) does.

It doesn't change the original object.

It only adds something before or after it.

like e.g

Customer
    │
    ▼
Record Time
    │
    ▼
Wear Gloves
    │
    ▼
Wash Hands
    │
    ▼
Cook Food

the cook never changes simply wrapping more behaviour around it(before or after)

Step 3: Same Idea in Go

Suppose we have a function.

func SayHello() {
    fmt.Println("Hello")
}

Output

Hello

Now suppose we want logging.

We could change it.

func SayHello() {
    fmt.Println("Starting...")
    fmt.Println("Hello")
    fmt.Println("Finished...")
}

Works.

But what if we have 500 functions?

You'll repeat logging everywhere.

Bad.

Instead we wrap it.

Original

func SayHello() {
    fmt.Println("Hello")
}

Wrapper

func Logging(next func()) func() {
    return func() {
        fmt.Println("Starting")
        next()
        fmt.Println("Finished")
    }
}

Now

hello := Logging(SayHello)

hello()

Output

Starting
Hello
Finished

Did SayHello change?

No.

We decorated it.

Step 4: Why is it called "next"?

Because the wrapper doesn't know what it's wrapping.

It simply says

"Whatever comes next, I'll call it."

Wrapper
   │
   ▼
next()

Maybe next is

SayHello

Maybe

Login

Maybe

DeleteUser

Wrapper doesn't care.

Step 5: HTTP Server Example

Suppose this is our handler.

func Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Home")
}

Without middleware

Client
   │
   ▼
Home Handler

Now we want logging.

Instead of changing Home,

we wrap it.

func Logging(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        log.Println("Request Started")

        next.ServeHTTP(w, r)

        log.Println("Request Finished")

    })
}

Now flow becomes

Client

   │

   ▼

Logging Middleware

   │

   ▼

Home Handler
Step 6: Visualize the Call Stack

Suppose request arrives.

Execution

Logging enters

↓

Print Start

↓

next.ServeHTTP()

↓

Home Handler

↓

returns

↓

Print Finish

↓

returns

Exactly like function calls.

Step 7: Multiple Wrappers

Now imagine three wrappers.

Recovery

↓

Logging

↓

Authentication

↓

Handler

Execution

Recovery enters

↓

Logging enters

↓

Authentication enters

↓

Handler runs

↓

Authentication exits

↓

Logging exits

↓

Recovery exits

Like nested boxes.

Recovery(
    Logging(
        Authentication(
            Handler
        )
    )
)
Step 8: Why "Wrapper"?

Because it literally wraps another object.

Wrapper
    │
    ▼
Original

It surrounds it.

Nothing more.

Step 9: Why "Decorator"?

This term comes from the classic Decorator Design Pattern.

Imagine a plain coffee.

Coffee

Add milk.

Milk(Coffee)

Add sugar.

Sugar(
      Milk(
          Coffee
      )
)

Coffee didn't change.

You decorated it.

Go middleware works exactly the same way.

Step 10: Middleware = Decorator

This is why middleware always looks like

func Middleware(next http.Handler) http.Handler

It takes

Handler

returns

New Handler

The new handler simply surrounds the old handler.

Old Handler

↓

Wrapped Handler

↓

Returned Handler
Step 11: The Mental Model

Imagine this diagram whenever you see middleware:

Incoming Request

        │

        ▼

+----------------+
|   Recovery     |
+----------------+
        │

        ▼

+----------------+
|   Logging      |
+----------------+
        │

        ▼

+----------------+
| Authentication |
+----------------+
        │

        ▼

+----------------+
|   Handler      |
+----------------+
        │

        ▼

Outgoing Response

Each layer:

does something before,
optionally calls next,
does something after.

That is the entire Decorator Pattern in Go.


# Embedding in go
Many people confuse embedding with inheritance from Java/C++. It is not inheritance. Go follows a different philosophy: "composition over inheritance".

Step 1: Imagine a Person

Suppose we have a Person.

type Person struct {
	Name string
	Age  int
}

Usage:

p := Person{
	Name: "Saurabh",
	Age: 22,
}

fmt.Println(p.Name)
fmt.Println(p.Age)

Output

Saurabh
22

Simple.

Step 2: Now We Need an Employee

Every employee is also a person.

In Java you might write

class Employee extends Person

Go says:

Don't inherit.

Embed.

type Employee struct {
	Person
	Salary int
}

Notice something?

There is no field name.

Not

Person Person

Just

Person

This is called embedding.

Step 3: Memory Layout

Think of it like this.

Employee

+----------------+
| Person         |
|  Name          |
|  Age           |
+----------------+
| Salary         |
+----------------+

Employee literally contains a Person.

It's not pointing to another object.

It's inside it.

Step 4: Creating One
e := Employee{
	Person: Person{
		Name: "Saurabh",
		Age: 22,
	},
	Salary: 50000,
}

Normally we'd expect

fmt.Println(e.Person.Name)

Right?

But Go does something magical.

Step 5: Field Promotion

Go automatically promotes embedded fields.

So both work.

fmt.Println(e.Person.Name)

AND

fmt.Println(e.Name)

Output

Saurabh

Why?

Because Go says:

If Employee doesn't have Name,
check the embedded Person.

It's as if Go searches inside.

Visual

Employee

Name?

↓

Employee has no Name

↓

Look inside Person

↓

Found Name
Step 6: Methods Are Promoted Too

Suppose Person has a method.

func (p Person) Introduce() {
	fmt.Println("Hi, I'm", p.Name)
}

Employee never defines it.

Yet

e.Introduce()

works.

Why?

Go looks here.

Employee

↓

No Introduce()

↓

Look inside Person

↓

Found Introduce()

Output

Hi, I'm Saurabh

Again,

Employee didn't inherit it.

Employee contains Person.

Step 7: Why Embed?

Without embedding

type Employee struct {
	Person Person
	Salary int
}

Access

e.Person.Name
e.Person.Age
e.Person.Introduce()

With embedding

type Employee struct {
	Person
	Salary int
}

Access

e.Name
e.Age
e.Introduce()

Less typing.

Cleaner API.

Step 8: Multiple Embedding

Go even allows this.

type Address struct {
	City string
}

type Employee struct {
	Person
	Address
}

Now

e.Name
e.City

work directly.

Go searches

Employee

↓

Name?

↓

Person

↓

Found

and

Employee

↓

City?

↓

Address

↓

Found
Step 9: Method Override?

Suppose Person has

func (p Person) Speak() {
	fmt.Println("Person speaking")
}

Employee defines

func (e Employee) Speak() {
	fmt.Println("Employee speaking")
}

Now

e.Speak()

prints

Employee speaking

The Employee method hides the promoted one.

You can still access the embedded version explicitly:

e.Person.Speak()

Output

Person speaking
Step 10: Where You'll See Embedding in Go
1. HTTP
type MyHandler struct {
	http.ServeMux
}

Your handler automatically gets all the methods of ServeMux, like Handle and ServeHTTP.

2. Database
type Repository struct {
	*pgxpool.Pool
}

Now the repository can directly call methods such as

repo.Query(...)
repo.Exec(...)

because the pool is embedded.

3. Configuration
type Server struct {
	Config
}

Instead of

server.Config.Port

you can write

server.Port
4. Logger
type App struct {
	*log.Logger
}

Then

app.Println("started")

instead of

app.Logger.Println("started")
Step 11: Embedding Interfaces

You can even embed interfaces.

type ReaderWriter struct {
	io.Reader
	io.Writer
}

Now any type implementing both interfaces satisfies ReaderWriter.

This is a common way to compose behaviors in Go.



1. Interfaces define contracts
type BookRepository interface {
    Create(ctx context.Context, book Book) error
    GetByID(ctx context.Context, id int) (Book, error)
}

Service doesn't care whether data comes from:

PostgreSQL
MongoDB
Redis
Mock for testing

It only knows:

"Give me something that satisfies BookRepository."

2. Embedding composes types

Suppose every HTTP handler needs a logger.

Without embedding:

type Handler struct {
    logger *slog.Logger
}

You write:

h.logger.Info("Creating book")

With embedding:

type Handler struct {
    *slog.Logger
}

Now you can write:

h.Info("Creating book")

The logger's methods are promoted to Handler.

3. Decorators wrap behavior

Your actual handler:

func CreateBook(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Book Created")
}

Now wrap it:

Logging(
    Authentication(
        Recovery(
            CreateBook,
        ),
    ),
)

Request flow:

Request
   │
   ▼
Logging
   │
   ▼
Authentication
   │
   ▼
Recovery
   │
   ▼
CreateBook

The handler never changes. Each wrapper adds one responsibility.

Notice something?

Each concept solves a completely different problem.

Interface

"I don't care who you are, as long as you can do this."

Service
   │
depends on
   ▼
BookRepository
Embedding

"I already have this inside me."

Handler
   │
contains
   ▼
Logger
Decorator

"Before I let you work, I'll do something."

Request

↓

Logging

↓

Authentication

↓

Handler


# logging and wrapped response writer need to learn properly 


# panic
func C() {
	panic("Boom!")
}

func B() {
	C()
}

func A() {
	B()
}

func main() {
	A()
}
### excecution 
main

↓

A

↓

B

↓

C

↓

panic

When panic happens,

Go starts unwinding the stack.

panic

↑

C exits

↑

B exits

↑

A exits

↑

main exits

↑

Program crashes

Panic means:"I cannot continue safely." not for normal business errror.
instead of returning error i gave up.

e.g if amazon one request panics. should the whole api stop? obv no,
instead 1 req panic i should enter into recover mode and continue serving req2.
but there is one rule recover works inside the defer function mean 
func main() {

	defer func() {

		if err := recover(); err != nil {
			fmt.Println("Recovered:", err)
		}

	}()

	panic("Boom!")

	fmt.Println("Never reached")
}
program continues..

why inside defer ?
defer mean execute this function just before returning ,even if panic occurs.
Function starts

↓

Register defer

↓

Work

↓

panic

↓

Run defer

↓

Return

# defer is like registering the action beforehand e.g"I'll do this before leaving"
defer doesn't execute immediately ,it executes when this function returns
func ReadFile() {

	file := Open()

	defer file.Close()

	// read file

}
Execution:
Open File

↓

Register Close()

↓

Read

↓

Read

↓

Read

↓

Function ends

↓

Close()


Without it:
file := Open()

// many returns...

file.Close()
You might forget to close the file on one of the return paths.

With defer:

file := Open()
defer file.Close()

No matter how the function exits, Close() runs.

# how defer works in case of multiple defer
it works like LIFO

# what makes recovery possible
func main() {

	defer fmt.Println("Cleaning")

	panic("Boom!")

}
Question:

Will "Cleaning" execute?
Most beginners think

Panic means everything stops.

Wrong.

Output

Cleaning

panic: Boom!

Why?
register defer ->panic->run all the defer->continue panic->program crash

# This is why recover() works

Now imagine

func B() {
	defer func() {

		if err := recover(); err != nil {
			fmt.Println("Recovered")
		}

	}()

	panic("Boom")
}

Where is the panic?

Inside B.

Where is the deferred function?

Inside B.

Perfect.

The panic reaches B's deferred functions before B returns.


# Understand the working of function registering defer ,panic and LIFO behaviour f defer to understand below flow
# When a handler panics:

Handler

↓

panic

↓

Recovery Middleware defer

↓

recover()

↓

Write 500 response

↓

Request ends safely

↓

Server continues serving other requests

The server stays alive because the middleware catches and stops the panic before it can escape the request-handling goroutine.




# Architected struct-based alternative and juxtaposed methodologies
Architected struct-based alternative and juxtaposed methodologies
Bilkul kar sakte ho — struct ek alternative route hai wahi problem solve karne ka. Farak sirf itna hai ki function ko struct ke andar field ke roop mein store karna padega, kyunki struct khud function nahi hai, wo ek container hai.

go
package main

import "fmt"

// A = interface we must satisfy
type A interface {
	M(x, y string)
}

// STRUCT approach: instead of naming the function itself,
// we wrap the function INSIDE a struct as a field
type B struct {
	f func(x, y string) // function stored as a field, not as the type itself
}

// method attached to the struct B (struct is a named type, so this works)
func (b B) M(x, y string) {
	b.f(x, y) // call the function stored inside the struct
}

// C = plain function
func C(x, y string) {
	fmt.Println("C running:", x, y)
}

func main() {
	// no direct type conversion possible now (C is not type B)
	// instead we construct a B, putting C into the field
	var e A = B{f: C}

	e.M("hello", "world")
	// call chain: e.M(...) -> B.M(...) -> b.f(x,y) -> C(x,y) actually runs
}
Struct vs Func-type — side by side
Func-type (HandlerFunc style)	Struct (B{f: C} style)
Kaise store karte ho function ko	Type khud function hai	Function ek field ke andar rehta hai
Conversion	B(C) — direct type conversion	B{f: C} — struct literal banao
Method ke andar call	f(x, y) — receiver khud function hai	b.f(x, y) — field access karke call
Extra data store kar sakte ho?	Nahi — sirf function hi hai, kuch aur nahi	Haan — struct mein aur fields bhi rakh sakte ho (jaise logger, config, counter)
Overhead	Bilkul minimal — bas ek func wrapper	Thoda zyada — struct allocate hota hai
Struct kab use karoge (real reason)
Agar tumhe sirf function wrap karna hai, HandlerFunc-style (func-type) hi kaafi hai — Go standard library isiliye ye use karti hai, simplicity ke liye.

Lekin agar tumhe function ke saath extra state bhi carry karni ho — jaise:

go
type B struct {
	f      func(x, y string)
	prefix string   // extra data
	count  int       // extra data
}

func (b *B) M(x, y string) {
	b.count++
	fmt.Println(b.prefix, "call #", b.count)
	b.f(x, y)
}
— tab struct zaroori ban jaata hai, kyunki func-type (type B func(...)) sirf function store kar sakta hai, koi extra field nahi rakh sakta. Struct ka fayda hai flexibility — jab middleware ko apni khud ki memory/state chahiye ho (jaise request counter, cache, ya config).

One-line summary: func-type = "function ko naya naam do". struct = "function ko ek box ke andar rakho, box mein aur cheezein bhi daal sakte ho."


# cross cutting concerns in api development
Middleware is not the only way to solve cross-cutting concerns.

Cross-cutting concerns are a problem:

"Many parts of the application need the same behavior."

Middleware is one solution, and it's the best solution for HTTP request/response processing.

In other contexts, cross-cutting concerns might be handled using:

Function wrappers
Decorators (common in Python/Java)
Interceptors (gRPC)
Filters (Java/Spring)
Aspect-Oriented Programming (AOP)

"Does every endpoint need this?"

If YES, it's probably a cross-cutting concern.


✅ Logging
✅ Authentication
✅ Authorization
✅ Rate limiting
✅ Request ID generation
✅ Metrics collection
✅ Panic recovery
✅ CORS
✅ Compression

Cross-cutting concerns are common tasks that "cut across" many parts of an application, so instead of writing them in every function, we put them in middleware or shared components.


# the concept of chaining of middle ware over any handler is done using chaining function 
// which is same as .Use() in frameworks like gin and others


# global vs routes specific middleware
Global middleware

Applied to every request.

app := middleware.Chain(
	router,
	middleware.Recovery,
	middleware.Logging,
	middleware.RequestID,
)
Route-specific middleware

Applied only to protected routes.

protected := middleware.Chain(
	bookHandler,
	middleware.Authentication,
)

This separation is very common in production services because not every endpoint has the same requirements.
# why request id for each request is required?
Imagine your API receives 500 requests per second.

Two users hit the same endpoint simultaneously.
user1:POST /books
user2:POST /books

Your logging middleware prints

Started POST /books
Started POST /books

Book inserted
Book inserted

Completed 201
Completed 201

Now imagine one request took 50ms and another took 4 seconds.

Which "Book inserted" belongs to which request?

You have no idea.


Now imagine your API looks like this:

Client

↓

Load Balancer

↓

API

↓

Redis

↓

Postgres

↓

Kafka

↓

Email Service

One request touches six different systems.

If something fails, how do you follow that request?

You can't.


Every incoming request gets a unique identifier. so, all logs for a single request can be grouped together.

This is why almost every production service has request IDs.
but should we store it?
globally naah ?
bcoz this will change the value once the new req arrives it sets the new value.

so as we know every req has a context (r.Context())
wherever this request goes it carries context That's exactly what context.Context is for.

we'll create a new request that carries the updated context.

This follows Go's design: http.Request should be treated as immutable for context changes.

# request_id middleware responsibility:
When a request arrives:

Request

↓

Generate ID

↓

Store ID

↓

Call next handler

That's it.

Notice it has one responsibility.

# The request scope is in each layer 
and each layer need some metadata
Request ID
User ID
Trace ID
Deadline
Cancellation signal
Locale
Tenant ID

Instead of passing all these values as separate parameters:

go bundle them in single object ctx context.Context

Q. they we donot attach the request id with http.request itself by modifiying the http request struct?
=>http.Request is part of Go's standard library. You cannot modify its definition.
Instead, the standard library provides the extension mechanism:
ctx := r.Context()
You can attach your own values without changing the type itself.

# Why shouldn't you use a string as a context key?

Because context.Context keys should be unique across packages. Using plain strings can lead to accidental key collisions
Defining an unexported custom type, such as type contextKey string, ensures that even if another package uses the same string value, the keys remain distinct because their types differ.

Let's prove it with actual Go code.

Case 1: Using a string key (Collision)

Imagine two packages.

logging package
package logging

import "context"

func AddRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, "request_id", "abc123")
}
auth package
package auth

import "context"

func AddUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, "request_id", "john")
}

Notice something?

Both use

"request_id"

as the key.

Now in main.go

ctx := context.Background()

ctx = logging.AddRequestID(ctx)
ctx = auth.AddUser(ctx)

fmt.Println(ctx.Value("request_id"))

What happens?

Let's execute it mentally.

Step 1
ctx = logging.AddRequestID(ctx)

Context contains

"request_id" -> "abc123"
Step 2
ctx = auth.AddUser(ctx)

This creates a new derived context with the same key:

"request_id" -> "john"

Now the context chain looks like

ctx2
 │
 ├── "request_id" -> "john"
 │
 ▼
ctx1
 │
 ├── "request_id" -> "abc123"
 │
 ▼
Background

Now Go searches from the newest context backwards.

When you ask

ctx.Value("request_id")

Go finds

john

first.

It never reaches "abc123".

So you've effectively hidden the earlier value.

Why?

Because

"request_id"

equals

"request_id"

Go compares:

string == string

Result:

true

So it considers them the same key.

Case 2: Custom Types

Now let's fix it.

logging package
package logging

type contextKey string

const requestIDKey contextKey = "request_id"

func AddRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey, "abc123")
}
auth package
package auth

type authKey string

const requestIDKey authKey = "request_id"

func AddUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey, "john")
}

Now look carefully.

The values are still

"request_id"

But the types are

logging.contextKey

and

auth.authKey

The context now contains

logging.contextKey("request_id")
             ↓
         abc123


auth.authKey("request_id")
          ↓
        john

These are two different keys.

Now imagine Go checking equality.

It compares

logging.contextKey("request_id")

with

auth.authKey("request_id")

Are they equal?

No.

Because Go first checks the type.

logging.contextKey

≠

auth.authKey

Different types.

Comparison fails.

No collision.

Here's the actual proof

Run this.

package main

import "fmt"

type contextKey string
type authKey string

func main() {

	var a contextKey = "request_id"
	var b authKey = "request_id"

	fmt.Printf("%T\n", a)
	fmt.Printf("%T\n", b)
}

Output

main.contextKey
main.authKey

Even though both store

request_id

their identities are different.

Another proof

This won't even compile.

package main

type contextKey string
type authKey string

func main() {

	var a contextKey = "request_id"

	var b authKey = a
}

Compiler says something like

cannot use a (type contextKey)
as authKey

Why?

Because Go treats them as completely different types.


# Packages should define "keys" as an unexported type to avoid collisions.to create unique identities for context keys across packages.
TODO://


#  a pkg/logger package, stop using the standard log package directly inside middleware.

Instead:

logger.Info(
    "request completed",
    "request_id", requestID,
    "method", r.Method,
    "path", r.URL.Path,
    "status", rw.statusCode,
    "duration", time.Since(start),
)

This is called structured logging.

Instead of one formatted string:

[RequestID=8c3f1d4e] POST /books 201 4ms

you produce structured fields:

level=INFO
request_id=8c3f1d4e
method=POST
path=/books
status=201
duration=4ms

This format is much easier for log aggregation tools (like Elasticsearch, Loki, or Datadog) to search and filter, which is why it's commonly used in production Go services.







# how request id is attached to the context as context.Context is immutable so logger can log the request id
Step 0: The request arrives

Go creates a request.

Request
│
├── Method = POST
├── URL = /books
└── Context = context.Background()

Imagine the context as an empty backpack.

Request
   │
   ▼
Context
┌──────────────┐
│              │
│   Empty      │
│              │
└──────────────┘

Now the request enters the first middleware.

Step 1: RequestID middleware starts
func RequestID(next http.Handler) http.Handler {

We generate an ID.

id := requestid.New()

Suppose

id = "ab82ef9137..."

Now we have two things.

Context (empty)

RequestID

ab82ef9137...

But they are not connected yet.

Step 2: IntoContext()

Now we call

ctx := requestid.IntoContext(
    r.Context(),
    id,
)

Internally

func IntoContext(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, key, id)
}

Now here's the MOST IMPORTANT PART.

Does context.WithValue() modify the old context?

NO

It creates a new context.

Think of it like this.

Old context

Context A

(empty)

After

context.WithValue(...)

Go creates

Context B

request_id = ab82ef9137...

Notice

Context A

(empty)

still exists.

Nothing was modified.

A new context was created.

Step 3: r.WithContext()

Now you have

Old Request

↓

Context A

and

New Context

↓

Context B

But the request still points to

Context A

So we do

r = r.WithContext(ctx)

Again

Does this modify the request?

No.

It creates a new request.

Before

Request A

↓

Context A

After

Request B

↓

Context B

Now the request carries the Request ID.

Step 4: next.ServeHTTP()

Now we call

next.ServeHTTP(w, r)

Notice carefully.

Which request are we passing?

Not

Request A

We pass

Request B

This request contains

Context

↓

request_id = ab82ef9137...
Step 5: Logging middleware

Now Logging starts.

func Logging(...)

It receives

Request B

not Request A.

Therefore

requestid.FromContext(
    r.Context(),
)

looks inside

Context B

and finds

request_id

↓

ab82ef9137...

Logging didn't generate it.

Logging didn't store it.

Logging simply reads it.

Step 6: Logging logs it
logger.Info(...)

prints

request_id=ab82ef9137...
method=POST
path=/books

Then Logging calls

next.ServeHTTP(...)
Step 7: Router

Router receives

Request B

The same request.

Router doesn't create another request.

It forwards it.

Step 8: Handler

Handler receives

Request B

Again

requestid.FromContext(
    r.Context(),
)

returns

ab82ef9137...

Exactly the same value.

When a request reaches the server, Go creates an http.Request with a default context. The first middleware is RequestID. It generates a unique ID using crypto/rand, then calls context.WithValue() to create a new derived context containing the Request ID. Since contexts are immutable, the original context isn't modified. Next, r.WithContext(newCtx) creates a new request associated with the new context. This new request is passed to the next middleware via next.ServeHTTP(). The Logging middleware receives this updated request, extracts the Request ID using requestid.FromContext(), logs it along with other request details, and forwards the same request to the router. The router dispatches it to the matching handler, and because the same request object continues through the handler, service, and repository, all of them can access the same Request ID from the context. This is why context is commonly used to carry request-scoped metadata throughout the lifetime of an HTTP request.

# Q. Why do we define a custom type for context keys instead of using a plain string?
context.Context stores values as key-value pairs, where the key is of type any (interface{}).

If we use a plain string as the key:

ctx = context.WithValue(ctx, "request_id", id)

another package might also use:

ctx = context.WithValue(ctx, "request_id", somethingElse)

Since both keys are of type string and have the same value ("request_id"), they collide.

To prevent this, we define our own unexported named type:

type contextKey string

const requestIDKey contextKey = "request_id"

Now the key is:

Type  : contextKey
Value : "request_id"

instead of

Type  : string
Value : "request_id"

Although the text is the same, the types are different, so Go treats them as different keys.

This gives each package its own namespace for context keys and prevents accidental collisions between unrelated packages.

Example

❌ Unsafe

ctx = context.WithValue(ctx, "user", user)

Another package:

ctx = context.WithValue(ctx, "user", logger)

Both use the same key (string("user")), so they collide.

✅ Safe

type contextKey string

const userKey contextKey = "user"

ctx = context.WithValue(ctx, userKey, user)

Another package's "user" key or even another contextKey type from a different package will not collide.

We use a custom unexported type for context keys because context.Context compares keys by both type and value. A named type gives the package its own namespace, preventing collisions with keys from other packages that might use the same string.
Context keys are identified by (type, value), not just by the string value.
A custom named type provides package-level namespacing and prevents key collisions.


# Why do we generate a Request ID?

In a production system, thousands of HTTP requests are processed concurrently.

Each incoming request is assigned a unique Request ID so that every log, error, trace, and downstream service call related to that request can be correlated.

Example:

Request 1 → RequestID = a1b2c3

Request 2 → RequestID = x9y8z7

Request 3 → RequestID = p4q5r6

This makes debugging and distributed tracing possible.

Why do we need requestIDKey?

After generating the Request ID, we need to store it in the request's context so every middleware and handler processing that request can access it.

ctx = context.WithValue(ctx, requestIDKey, requestID)

Here:

requestIDKey → the identifier (key)
requestID    → the actual unique ID (value)

Later:

id := ctx.Value(requestIDKey)

retrieves the same Request ID.

Why is requestIDKey a custom type?

The context is shared across all middleware and packages handling the same request.

Many packages may store their own values in that shared context.

If everyone used plain string keys like:

"user"
"request_id"
"trace_id"

different packages could accidentally use the same key, causing collisions.

Therefore we define:

type contextKey string

const requestIDKey contextKey = "request_id"

A custom named type creates a unique namespace for the package's keys, preventing accidental collisions.


# authentication middlware
Authentication middleware does one thing:

Read JWT

↓

Validate JWT

↓

Extract User ID

↓

Put User ID into Context

flow:
Client

Authorization: Bearer xxxxx

        │
        ▼

Authentication Middleware

        │

Read Header

        │

Validate JWT

        │

Extract UserID

        │

IntoContext(userID)

        │

next.ServeHTTP()


# What changes in Context?

Before authentication

Context

request_id

After authentication

Context

request_id

user_id

Later the handler simply reads

requestID := requestid.FromContext(r.Context())

userID := auth.FromContext(r.Context())

Notice

Handler doesn't know JWT exists.


# notice
requestid exposes:

requestid.IntoContext()
requestid.FromContext()

The package owns everything related to Request IDs.

auth exposes:

auth.IntoContext()
auth.FromContext()

The package owns everything related to authenticated users.

Notice that callers never touch the keys.

That's encapsulation.


# jwt 
Think of a JWT as a Digital ID Card

Suppose you work at Google.

Your employee ID card says:

----------------------------
Google Employee Card

ID      : 42
Name    : Saurabh
Role    : Backend Engineer
Email   : saurabh@google.com
Expires : 31 Dec 2026

----------------------------

Everything written inside the ID card is information about you.

In JWT terminology, this information is called Claims.

A JWT has three parts
Header
.
Payload
.
Signature

Example

xxxxx.yyyyy.zzzzz

The middle part (Payload) contains the claims.

Payload (Claims)

After decoding the payload, it may look like this:

{
    "user_id": "42",
    "email": "saurabh@gmail.com",
    "role": "admin",
    "exp": 1753330000,
    "iat": 1753320000
}

Everything inside this JSON is called Claims.

Think of it as:

JWT

↓

Payload

↓

Claims
There are two kinds of Claims
1. Standard Claims

These are defined by the JWT specification.

Examples:

{
    "exp": 1753330000,
    "iat": 1753320000,
    "iss": "my-api",
    "sub": "42"
}

Meaning:

exp → Expiration Time
iat → Issued At
iss → Issuer
sub → Subject (usually user ID)

These are understood by JWT libraries.

2. Custom Claims

These are created by your application.

Example:

{
    "user_id":"42",
    "email":"abc@gmail.com",
    "role":"admin"
}

JWT doesn't know what these mean.

They're meaningful only to your application.

That's why we created
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

Notice two things.

Your custom claims
UserID
Email
Role

These become

{
    "user_id":"42",
    "email":"abc@gmail.com",
    "role":"admin"
}
Standard claims
jwt.RegisteredClaims

This gives you fields like

ExpiresAt
IssuedAt
Issuer
Subject

which become

{
    "exp": ...,
    "iat": ...
}

Suppose we create a token for you.

Claims object

Claims{
    UserID: "42",
    Email: "saurabh@gmail.com",
    Role: "admin",

    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: ...
        IssuedAt: ...
    },
}

This becomes

{
    "user_id":"42",
    "email":"saurabh@gmail.com",
    "role":"admin",

    "exp":1753330000,
    "iat":1753320000
}

Then the JWT library

Signs
↓

Encodes

↓

Creates Token

which finally looks like

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
.
eyJ1c2VyX2lkIjoiNDIiLCJyb2xlIjoiYWRtaW4i...
.
lQmH3Lk...


# Why do we parse Claims later?

Suppose a client sends

Authorization: Bearer eyJhbGc...

Your middleware validates the signature.

Then it extracts the payload.

That payload becomes

Claims{
    UserID:"42",
    Email:"abc@gmail.com",
    Role:"admin",
}

Finally, you convert it into your domain model:

return User{
	ID:    claims.UserID,
	Email: claims.Email,
	Role:  claims.Role,
}, nil

Now the rest of your application doesn't care about JWT anymore—it just works with a User.

One important design insight

Notice how Claims and User are different types, even though they have similar fields.

type Claims struct {
	UserID string
	Email  string
	Role   string
	jwt.RegisteredClaims
}

type User struct {
	ID    string
	Email string
	Role  string
}

That's intentional.

Claims represents how user information is stored inside a JWT (transport format).
User represents how your application models an authenticated user (domain model).

Separating these types means that if the JWT format changes—for example, you add tenant_id or change claim names—you only update the mapping in jwt.go. The rest of your handlers, services, and repositories continue working with the same User type. This separation of transport models from domain models is a common production design pattern.


# jwt sentinel errors
Authentication Flow
Client
    │
    ▼
Authorization: Bearer eyJhbGc...

    │
    ▼
Authentication Middleware

    │
    ▼
Extract Token

    │
    ▼
Validate JWT

    │
    ├──────────────► Is token malformed?
    │                    │
    │                    ▼
    │              ErrMalformedToken
    │
    ├──────────────► Is signature valid?
    │                    │
    │                    ▼
    │              ErrInvalidToken
    │
    ├──────────────► Is token expired?
    │                    │
    │                    ▼
    │              ErrExpiredToken
    │
    ├──────────────► Are required claims present?
    │                    │
    │                    ▼
    │              ErrMissingClaims
    │
    ▼
User

Now let's understand each one.

1. ErrMalformedToken
ErrMalformedToken = errors.New("malformed token")
When does this happen?

The client sends something that isn't a JWT.

Example:

Authorization: Bearer hello

or

Authorization: Bearer 12345

A JWT has three parts:

xxxxx.yyyyy.zzzzz

(header.payload.signature)

Suppose the client sends

abc

There are no dots.

The JWT library can't even parse it.

It immediately fails.

JWT Parser
     │
     ▼

abc

↓

❌ Not a JWT

↓

ErrMalformedToken
2. ErrInvalidToken
ErrInvalidToken = errors.New("invalid token")

Now suppose the token looks valid.

xxxxx.yyyyy.zzzzz

But someone changed the payload.

Example

Original token

role=user

Attacker edits it to

role=admin

without knowing your secret key.

Now the signature no longer matches.

The server calculates

Expected Signature

↓

abc123

Received

xyz999

Mismatch.

JWT

↓

Signature Verification

↓

❌

↓

ErrInvalidToken

This is probably the most common authentication failure.

3. ErrExpiredToken

Every JWT usually contains

{
  "exp": 1753320000
}

Meaning

This token expires at 10:30 AM.

Suppose

Current time

11:45 AM

The library checks

Now > exp ?

↓

YES

↓

Expired

Return

ErrExpiredToken
4. ErrMissingClaims

Suppose your application requires

{
    "user_id": "...",
    "email": "...",
    "role": "..."
}

But someone sends

{
    "email":"abc@gmail.com"
}

No

user_id

Or

role

exists.

The token is correctly signed.

The signature is valid.

The token isn't expired.

But your application cannot authenticate the user because required information is missing.

So

Claims

↓

Missing user_id

↓

ErrMissingClaims
Where do these errors come from in code?

Imagine this.

func (j *JWT) ValidateToken(token string) (User, error) {

    parsedToken, err := jwt.ParseWithClaims(...)

First possibility

if err != nil {

    if errors.Is(err, jwt.ErrTokenMalformed) {
        return User{}, ErrMalformedToken
    }

    if errors.Is(err, jwt.ErrTokenExpired) {
        return User{}, ErrExpiredToken
    }

    return User{}, ErrInvalidToken
}

Now parsing succeeded.

We extract claims.

claims := parsedToken.Claims.(*Claims)

Suppose

claims.UserID == ""

Then

return User{}, ErrMissingClaims

Otherwise

return User{
    ID: claims.UserID,
    Email: claims.Email,
    Role: claims.Role,
}, nil

# the concept of sentinel errors
A sentinel error is a predefined, reusable error value that represents a specific condition. Instead of creating a new error every time, you define it once and compare against it later.
In Go, it's usually created like this:

var ErrNotFound = errors.New("not found")

Then anywhere in your code, you return that exact error:

func FindUser(id int) error {
    if id != 1 {
        return ErrNotFound
    }
    return nil
}

And the caller checks it:

err := FindUser(2)

if errors.Is(err, ErrNotFound) {
    fmt.Println("User doesn't exist")
}
Whenever the "user not found" situation happens, the function returns thesame error value. that's why it's called sentinel
Without a sentinel error

Imagine writing:

func FindUser(id int) error {
    if id != 1 {
        return errors.New("user not found")
    }
    return nil
}

Now every call creates a new error object.

Call 1
errors.New("user not found")
      ↓
Address A

Call 2
errors.New("user not found")
      ↓
Address B

If you compare them:

err == errors.New("user not found")

it will be false because they're different error values.

With a sentinel error
var ErrUserNotFound = errors.New("user not found")

func FindUser(id int) error {
    if id != 1 {
        return ErrUserNotFound
    }
    return nil
}

Now every failure returns the same error variable.

Call 1 ───────► ErrUserNotFound

Call 2 ───────► ErrUserNotFound

Call 3 ───────► ErrUserNotFound

So checking becomes reliable.



In your JWT project

You recently defined:

var (
    ErrInvalidToken   = errors.New("invalid token")
    ErrExpiredToken   = errors.New("token expired")
    ErrMalformedToken = errors.New("malformed token")
    ErrMissingClaims  = errors.New("missing claims")
)

These are sentinel errors.

Your parser can return them:

return ErrExpiredToken

or wrap them:

return fmt.Errorf("parse jwt: %w", ErrExpiredToken)

And callers can reliably check:

switch {
case errors.Is(err, ErrExpiredToken):
    // Tell client the token has expired

case errors.Is(err, ErrInvalidToken):
    // Tell client the token is invalid

case errors.Is(err, ErrMalformedToken):
    // Tell client the JWT format is incorrect

default:
    // Internal server error
}

# Why use errors.Is instead of ==?

Imagine you have a person

His name is John.

var John = Person{Name: "John"}

Now suppose I give you John.

John
 ▲
 │
person

If I ask,

Is person == John?

Answer:

YES

Because it's the same person.

Now imagine I put John inside a car.
+-----------+
|   CAR     |
|           |
|  John     |
+-----------+

Now the variable points to the car.

person
   │
   ▼
+-----------+
| CAR        |
|            |
| John       |
+-----------+

If I ask

Is the car == John?

NO

Because a car is not John.

John is inside the car.

This is exactly what %w does.

Without %w
return ErrUserNotFound
err
 │
 ▼

ErrUserNotFound

Both point to the same thing.

So

err == ErrUserNotFound

✅ True

With %w
return fmt.Errorf("find user 42: %w", ErrUserNotFound)

Go creates a new error.

err
 │
 ▼

+--------------------------------+
| find user 42:                  |
|                                |
| ErrUserNotFound                |
+--------------------------------+

Notice

err is NOT ErrUserNotFound.

It is a new error that contains ErrUserNotFound.

So why is == false?

Because you're comparing

Wrapper Error

with

ErrUserNotFound

Like comparing

Car
==
John

Are they the same object?

NO
Then what does errors.Is do?

errors.Is is like asking

"Is John inside this car?"

It opens the car.

Car

↓

Open it

↓

John

Finds John.

Returns

true
Let's use your JWT example.

You have

var ErrExpiredToken = errors.New("token expired")

Later you write

return fmt.Errorf("parse jwt failed: %w", ErrExpiredToken)

The returned error becomes

parse jwt failed: token expired

But internally it looks like

+------------------------------------+
| parse jwt failed                   |
|                                    |
| contains                           |
|                                    |
| ErrExpiredToken                    |
+------------------------------------+

Now look carefully.

Using ==
err == ErrExpiredToken

Go asks

Is this big box exactly ErrExpiredToken?

Answer:

NO
Using errors.Is
errors.Is(err, ErrExpiredToken)

Go asks

Does this box contain ErrExpiredToken?

It opens the box.

Box

↓

ErrExpiredToken

↓

Found

↓

true
The only thing I want you to remember is this:

== checks:

"Are these exactly the same error object?"

errors.Is() checks:

"Is this error, or anything wrapped inside it, the error I'm looking for?"



# how authentication works
1. claims.go
Responsibility

Define what information is stored inside the JWT.

type Claims struct {
    UserID string
    Email  string
    Role   string

    jwt.RegisteredClaims
}

Think of a passport.

A passport contains

Name
Nationality
Passport Number
Expiry

Similarly, a JWT contains

UserID
Email
Role
Expiry
IssuedAt
Issuer

This file only defines the structure.

It doesn't create tokens.

It doesn't validate them.

It doesn't know HTTP.

Just a data model.

2. errors.go
Responsibility

Defines authentication errors.

ErrExpiredToken
ErrInvalidToken
ErrMalformedToken
ErrMissingClaims

Think of airport security.

Instead of saying

Something went wrong

they give precise reasons

Passport expired

Passport damaged

Passport fake

Those are your sentinel errors.

Every other file understands these errors.

3. user.go

Responsibility

Represents the authenticated user inside your application.

type User struct {
    ID
    Email
    Role
}

Notice

It has NO

ExpiresAt

IssuedAt

Issuer

Why?

Because after authentication, your application doesn't care about JWT metadata.

It only cares about

Who is the user?
4. context.go

Responsibility

Store and retrieve the authenticated user from the request context.

IntoContext()

FromContext()

Think of this as attaching an ID badge to the request.

Incoming Request

↓

Authenticate

↓

Attach User Badge

↓

Pass request forward

Later,

Any handler can simply do

user, ok := auth.FromContext(ctx)

without knowing anything about JWT.

5. jwt.go

This is the engine.

Responsibility

Manage JWT lifecycle.

It does

GenerateToken()

ValidateToken()

This file knows

JWT library
Signing method
Secret
Expiry
Claims

It DOES NOT know

HTTP
Context
Middleware
Handlers
GenerateToken()

Input

User

↓

Creates

Claims

↓

Creates

JWT

↓

Signs it

↓

Returns

Token String

Flow

User

↓

Claims

↓

JWT

↓

Signed Token
ValidateToken()

Input

JWT String

↓

Verify signature

↓

Check expiry

↓

Decode payload

↓

Return Claims

Token

↓

Claims

Nothing else.

6. authentication.go

This is the security guard.

Its responsibility is

Authenticate every incoming request.

It does NOT know

how JWT is created
how signature works

It only knows

I have a token.

Ask JWTManager to validate it.
Complete Data Flow

Suppose a user logs in.

Step 1

User logs in.

POST /login

Backend verifies

Email

Password
Step 2

Backend creates

user := auth.User{
    ID: "123",
    Email: "abc@gmail.com",
    Role: "admin",
}
Step 3

JWT Manager

GenerateToken(user)

takes

User

↓

creates

Claims

↓

creates JWT

↓

returns

eyJhbGc...

Send to client.

Now imagine another request.

Client sends

GET /books

Authorization:
Bearer eyJhbGc...
Authentication Middleware

Reads

Authorization Header

↓

Extracts

eyJhbGc...

↓

Calls

ValidateToken()
JWT Manager

Receives

eyJhbGc...

↓

Verifies signature

↓

Checks expiry

↓

Parses payload

↓

Returns

Claims

Middleware now has

Claims
UserID

Email

Role

Expiry

Middleware converts

Claims

↓

User

because handlers don't need

Expiry

Issuer

IssuedAt

Only

ID

Email

Role

Middleware stores

User

inside

Context
Context

↓

User

Finally

next.ServeHTTP()

is called.

Handler

Handler simply writes

user, ok := auth.FromContext(r.Context())

It immediately gets

ID

Email

Role

without parsing JWT again.

End-to-End Flow Diagram
                  LOGIN
                    │
                    ▼
              Verify Credentials
                    │
                    ▼
                 auth.User
                    │
                    ▼
           JWTManager.GenerateToken()
                    │
                    ▼
                 Claims
                    │
                    ▼
               Signed JWT Token
                    │
             ───── Client ─────
                    │
                    ▼
         Authorization: Bearer <token>
                    │
                    ▼
     Authentication Middleware
                    │
                    ▼
     JWTManager.ValidateToken()
                    │
                    ▼
                 Claims
                    │
                    ▼
          Convert Claims → User
                    │
                    ▼
      auth.IntoContext(ctx, user)
                    │
                    ▼
            next.ServeHTTP()
                    │
                    ▼
                 Handler
                    │
                    ▼
auth.FromContext(r.Context()) → User
Why this design scales well

Each file answers exactly one question:

File	Responsibility	Input	Output
claims.go	What data goes inside a JWT?	—	Claims type
errors.go	Which auth errors exist?	—	Sentinel errors
user.go	What is an authenticated user?	—	User type
context.go	How do we carry the authenticated user through the request?	User, context.Context	Updated context or retrieved User
jwt.go	How do we create and validate JWTs?	User or token string	Signed token or Claims
authentication.go	How do we protect HTTP endpoints?	HTTP request	Authenticated request with User in context

Notice how the data transforms as it moves through the system:

User (application model)
        │
        ▼
Claims (JWT payload)
        │
        ▼
JWT Token (serialized string)
        │
        ▼
Claims (decoded after validation)
        │
        ▼
User (application model again)
        │
        ▼
Context
        │
        ▼
Handlers

That's the key architectural idea: JWT is only a transport format for identity. The rest of your application works with User, not with JWT internals.

# bootstraping and dependency injection (composition root)

"How do I assemble an application?"

That process is called the Composition Root (also known as the Bootstrap Layer).

Let's start with your current main.go

Right now, your main.go does this:

Load Config
      │
      ▼
Connect Database
      │
      ▼
Create Repository
      │
      ▼
Create Service
      │
      ▼
Create Handler
      │
      ▼
Create JWTManager
      │
      ▼
Create Router
      │
      ▼
Register Routes
      │
      ▼
Apply Middleware
      │
      ▼
Start HTTP Server

Notice something?

It isn't doing business logic.

It's just connecting objects together.

This is called composition.

What is Composition?

Imagine you're building a gaming PC.

You buy:

CPU
Motherboard
RAM
SSD
GPU
Power Supply

Each component works independently.

Then you assemble them.

That assembly process is composition.

Your application is exactly the same.

You have:

Config
Database
Repository
Service
Handler
JWTManager
Router
Middleware
Server

Individually they're useless.

When connected together,

they become an application.

Composition Root

The place where everything is assembled is called the Composition Root.

In small Go applications,

it's usually

main.go

because everything is wired there.

Dependency Graph

Every object depends on another object.

Example

Repository
    ▲
    │
Database

Service depends on Repository.

Service
    ▲
    │
Repository

Handler depends on Service.

Handler
    ▲
    │
Service

Authentication Middleware depends on JWTManager.

Authentication
      ▲
      │
 JWTManager

Server depends on Router.

Server
    ▲
    │
Router

Together

Server
   ▲
Router
   ▲
Handler
   ▲
Service
   ▲
Repository
   ▲
Database

This is called the dependency graph.

Dependency Injection

Notice something.

Repository isn't creating Database.

Instead

main.go creates Database.

Then

repo := NewRepository(db)

Database is injected.

Likewise

service := NewService(repo)

Repository is injected.

Likewise

handler := NewHandler(service)

Service is injected.

Nothing creates its own dependencies.

That's Dependency Injection.

Why?

Imagine Repository did this.

func NewRepository() *Repository {

    db := database.Connect()

    return &Repository{
        db: db,
    }
}

Now Repository is tightly coupled to Database.

You can't

mock database
replace postgres
test repository

Everything becomes difficult.

Instead

repo := NewRepository(db)

Repository doesn't care where db came from.

That's inversion of control.

Why does main.go become huge?

Suppose your application grows.

Today

Book

Tomorrow

Book
User
Order
Payment
Inventory
Notification
Analytics
Search

Each has

Repository

Service

Handler

Now main.go becomes

bookRepo := ...
bookService := ...
bookHandler := ...

userRepo := ...
userService := ...
userHandler := ...

orderRepo := ...
orderService := ...
orderHandler := ...

paymentRepo := ...
paymentService := ...
paymentHandler := ...

300+ lines.

Nothing wrong.

Just noisy.

Solution

Move composition into another package.

Example

cmd/
    api/
        main.go

internal/
    app/
        app.go
app.go

This becomes the new composition root.

func New() *App {

    cfg := ...

    db := ...

    repo := ...

    service := ...

    handler := ...

    router := ...

    server := ...

    return &App{
        server: server,
    }
}
main.go

Now becomes

func main() {

    app := app.New()

    log.Fatal(app.Run())
}

Beautiful.

This is called Bootstrapping

The process

Create Config

↓

Create Database

↓

Create Dependencies

↓

Create Server

↓

Run

is called

Application Bootstrapping.

Why "Composition Root"?

Because

all dependencies meet here.

Nothing else creates dependencies.

Everything starts here.

Think of a tree.

                 main.go
                     │
      ┌──────────────┼──────────────┐
      ▼              ▼              ▼
   Config         Database     JWTManager
                     │
                     ▼
               Repository
                     │
                     ▼
                 Service
                     │
                     ▼
                 Handler
                     │
                     ▼
                  Router
                     │
                     ▼
                 Middleware
                     │
                     ▼
                   Server

Everything originates from one root.

Hence

Composition Root.

Why do big companies care?

Imagine a service at companies like Stripe, Uber, or Netflix.

Instead of 8 objects,

they may have hundreds:

Database connections
Redis clients
Kafka producers
Kafka consumers
S3 clients
JWT manager
Metrics
Tracing
Logger
Feature flags
Cache
Mailer
Search client
Payment gateway
External APIs

All of these need to be created in the right order and passed to the components that depend on them. Keeping that wiring in one place—the composition root—makes it much easier to understand how the application is assembled, swap implementations for testing, and avoid hidden dependencies.

A subtle but important rule

The composition root knows about everything.

Your main.go (or app.New()) is allowed to import:

config
database
book
auth
middleware

But the reverse should not happen.

For example:

book/service should not import main.
repository should not create its own database connection.
middleware should not call config.Load() by itself.

This keeps dependencies flowing in one direction:

Composition Root
        │
        ▼
 Infrastructure (DB, Config, Logger)
        │
        ▼
 Repository
        │
        ▼
 Service
        │
        ▼
 Handler
        │
        ▼
 HTTP Server

Each layer receives what it needs from above instead of reaching out to build its own dependencies. That's one of the foundations of maintainable Go applications and a pattern you'll encounter repeatedly in production codebases.


# eaasiest way to write docker file
# Production Go Dockerfile Notes

``` dockerfile
// What am I containerizing? Is it Go, Java, Python, C#?
// What base image should I use for that language?

// This image contains a lot of tools:
//
// golang:1.25
//
// ├── Linux OS
// ├── Go Compiler
// ├── Go Toolchain
// ├── Standard Library
// ├── Git
// ├── Shell
// ├── Build Tools
// ├── Package Cache
// └── Your Source Code
//
// All of this is necessary only to build the executable.
//
// After running:
//
// go build -o server ./cmd/api
//
// we get:
//
// server
//
// This is a Linux executable (machine code) generated by the Go compiler.

FROM golang:1.25 AS builder

// What will be the working directory inside the container?
// If this directory does not exist, Docker automatically creates it.

WORKDIR /app

// Now we need the project's dependencies.
// We first copy go.mod and go.sum because they define all the dependencies.
// Why not copy the entire project first?
// We'll understand that after learning Docker cache.

COPY go.mod go.sum ./

// Download all dependencies.
// Since go.mod and go.sum are already inside /app,
// Go can download every required module.

RUN go mod download

// Now all downloaded dependencies are stored inside the builder image.

// Q. Why do we copy dependency files first?
//
// A. Docker builds images layer by layer.
// If a layer changes, every layer after it must be rebuilt.
//
// If we copied the entire source code first:
//
// COPY . .
// RUN go mod download
//
// then changing even one Go file would force Docker
// to download every dependency again.
//
// Instead:
//
// COPY go.mod go.sum ./
// RUN go mod download
// COPY . .
//
// If only source code changes and go.mod/go.sum stay the same,
// Docker reuses the cached dependency layer.
// So RUN go mod download is skipped, making builds much faster.

// Now copy the entire project into the working directory.
//
// First "."  -> current directory on the host machine.
// Second "." -> current WORKDIR (/app) inside the container.

COPY . .

// Now we have both:
// - Source code
// - Dependencies
//
// So we can build the application.
//
// Many production projects have multiple entry points:
//
// cmd/api
// cmd/worker
// cmd/scheduler
//
// Go builds packages, not individual files.
// So instead of building main.go,
// we build the package that contains main().
//
// Go automatically compiles every .go file inside that package.
//
// If tomorrow you add:
//
// cmd/api/routes.go
// cmd/api/config.go
//
// Go automatically includes them during compilation.
//
// Production command:
//
// RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api
//
// This compiles the entire cmd/api package,
// along with all imported packages,
// and creates a Linux executable named "server".

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

// Alternative examples:
//
// RUN go build -o worker ./cmd/worker
// RUN go build -o scheduler ./cmd/scheduler

// ------------------------------------------------------------------

// Now the builder stage contains:
//
// Source Code
// Dependencies
// Go Compiler
// Git
// Cache
// server (compiled executable)
//
// We no longer need the compiler or source code.
//
// We only need the executable.
//
// So we start a brand-new, lightweight Linux image.

FROM alpine:latest

// This is a completely new image.
// Nothing from the builder stage exists here.

// Set the working directory again.
// Every stage has its own filesystem.

WORKDIR /app

// Copy only the executable from the builder stage.
//
// COPY --from=<stage-name> <source> <destination>
//
// --from=builder changes the source from:
//
// Host Machine
//      ↓
// Builder Stage
//
// Docker copies:
//
// Builder:
/app/server
//
//      ↓
//
// Runtime:
/app/server
//
// Only the executable is copied.
//
// It does NOT copy:
//
// - Source code
// - Go compiler
// - Git
// - go.mod
// - go.sum
// - Downloaded dependencies
//
// This is the main advantage of multi-stage builds.

// About Docker flags:
//
// -  : Short option (usually one letter)
//
//      -t
//      -p
//      -d
//      -f
//
// -- : Long option (full descriptive name)
//
//      --tag
//      --publish
//      --detach
//      --file
//      --from

COPY --from=builder /app/server /app/server

// Now the runtime image contains only:
//
// /app
// └── server
//
// It is much smaller because all build tools
// stayed behind in the builder stage.

// Inform Docker that this application listens on port 8080.
// EXPOSE is documentation for the image.
// It does NOT publish the port automatically.

EXPOSE 8080

// Start the application.
//
// Docker executes:
//
// ./server
//
// when the container starts.

CMD ["./server"]
```


# we have to follow this 
□ What am I containerizing? (Go, Node, Python, Java...)

□ What base image fits?

□ What dependencies must be installed?

□ How do I build the app?

□ What artifact should run?

□ What port does it listen on?

□ What command starts it?

□ What other services does it need?
   • PostgreSQL?
   • Redis?
   • Kafka?
   • Nginx?

□ How will services communicate?
   • Network
   • Environment variables

□ What data must persist?
   • Volumes

□ How should it recover?
   • restart policy

| Question                               | Dockerfile | Docker Compose | Proof                                                                      |
| -------------------------------------- | ---------- | -------------- | -------------------------------------------------------------------------- |
| ✅ What am I containerizing?            | ✔          |                | `FROM golang:1.25` tells Docker we're containerizing a **Go** application. |
| ✅ What base image fits?                | ✔          |                | `golang:1.25` for building, `alpine:latest` for running.                   |
| ✅ What dependencies must be installed? | ✔          |                | `COPY go.mod go.sum ./` → `RUN go mod download`                            |
| ✅ How do I build the app?              | ✔          |                | `RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api`                |
| ✅ What artifact should run?            | ✔          |                | The compiled executable `server`                                           |
| ✅ What port does it listen on?         | ✔          |                | `EXPOSE 8080`                                                              |
| ✅ What command starts it?              | ✔          |                | `CMD ["./server"]`                                                         |
| ✅ What other services does it need?    |            | ✔              | PostgreSQL, Redis, Kafka, Nginx are defined as services in Compose.        |
| ✅ How will services communicate?       |            | ✔              | Networks + environment variables (`DB_HOST=postgres`)                      |
| ✅ What data must persist?              |            | ✔              | Volumes (`postgres-data`)                                                  |
| ✅ How should it recover?               |            | ✔              | `restart: unless-stopped`                                                  |

A Dockerfile never asks about PostgreSQL, Redis, Kafka, or volumes because a Dockerfile is only concerned with one image.
Compose answers the remaining questions.

□ What services?
      ↓
   app
   postgres
   redis

□ How do they communicate?
      ↓
   backend network

□ Environment variables?
      ↓
   DB_HOST=postgres

□ Persistent data?
      ↓
   postgres-data volume

□ Recovery?
      ↓
   restart: unless-stopped


Q. why you docker file should not have ohter services
This is considered bad practice.
Why?
Because containers should have one primary responsibility.
Go Container:Runs API

Postgres Container:Runs Database

Redis Container:Runs Cache
## and Docker Compose is responsible for connecting these independent containers together.
Docker Compose orchestrates multiple containers (your app, databases, caches, proxies, message brokers, etc.) so they run together as a complete system.


# how to write docker compose 
ask yourself what services starts my application?
and we know that one docker file produce one image
and we need docker compose that can run all the containers together by making connection to each other 

so,docker compose starts with services beacause everything inside your system is service  and compose manages services containers
services:

  app:

  postgres:

  redis:

For every service ask

Do I build it?

or

Do I download it?

1.if ur code build it using docker image and then use in service
app:
  build:
    context: .

2.else it's an image 
image: postgres:17
image: redis:8-alpine


now u have all the services
now what users need to access?
the user usually access app or api
so we expose the port so inbound connection from user can come to the the container port
so it's usually HOST : CONTAINER

ports:

  - "8080:8080"

# how do services find each other?
So how does the app find PostgreSQL?
Docker automatically gives every service a DNS name.
Docker Starts a DNS Server

When you run

docker compose up

Docker does several things automatically.

Create Network

↓

Start DNS Server

↓

Start Containers

↓

Register Every Service

Notice the important step.

Register Every Service
Suppose your compose file is
services:

  app:

  postgres:

  redis:

Docker internally creates something like

Docker DNS

postgres

↓

172.20.0.2


redis

↓

172.20.0.3


app

↓

172.20.0.4

You never see this.

Docker does it automatically.


So how does the Go app connect?

Your compose file

services:

  app:

    environment:

      DB_HOST: postgres

Inside Go

host := os.Getenv("DB_HOST")

returns

postgres

Now Go tries to connect.

postgres:5432

Question

How does it become an IP?

Here's the Hidden Step
Go Application

↓

Connect to

postgres:5432

↓

Linux asks

"Who is postgres?"

↓

Docker DNS

↓

postgres

↓

172.20.0.2

↓

TCP Connection

↓

PostgreSQL

The application never knows the IP.

It only knows

postgres

Docker handles the translation.

### In Docker Compose, the service name becomes the hostname. Docker's built-in DNS automatically resolves that hostname to the correct container IP, so containers communicate using service names instead of hardcoded IP addresses.

That's why in production Go code you'll often see connection strings like:

connStr := "host=postgres port=5432 user=postgres password=secret dbname=app sslmode=disable"


# should data survive?
suppose we do docker compose down/docker rm postgres mean database container disappears
customer data lost

so ask which services have important data?
postgres
redis
mongodb
sql
kafka
the services have important data need persitance volume for data persistance and volume to give data a space to live permanently.
so where do the data lives ? outside the container in persistent storage(database files)

<Named docker Volume> : <Directory Inside Container>

Docker-managed storage:Location inside container

e.g - postgres-data:/var/lib/postgresql/data

## the docker manages volumes seperately from containers
## PostgreSQL stores all its database files here "/var/lib/postgresql/data"

so 
volumes:
- postgres-data:/var/lib/postgresql/data


# now we have starts the services
Which services must start first?
if your 
Your API starts.
Database isn't ready.
Application crashes.

so we start dependency of api/app first then api/app.
like here postgres then redis and likewise

# how do services communicate?
compose create networks

We have multiple isolated containers.
Every container has its own localhost.
If your Go application does

localhost:5432

What happens?

It searches inside its own container. and postgres not found 
Because PostgreSQL is in another container.

so containers need addresses
Docker automatically creates a private network.
Every container gets its own IP address.
Docker includes a DNS server.

Suppose your compose file says

services:

  postgres:

  redis:

  app:

Docker automatically creates DNS records.

postgres

↓

172.18.0.3

redis

↓

172.18.0.4

Now your application simply connects to

host := "postgres"

instead of

host := "172.18.0.3"

### How does the app know the name?
We pass it using environment variables.

services:

  app:

    environment:

      DB_HOST: postgres

      REDIS_HOST: redis


networks:
  backend:
Create a private virtual LAN (Local Area Network) named "backend".

How do services join it?

Example

services:

  app:
    networks:
      - backend

  postgres:
    networks:
      - backend

  redis:
    networks:
      - backend

networks:
  backend:

Create a private network called backend, and let every service attached to it talk to the others securely.
Read it like English.

Create a network named backend.

↓

Connect app to backend.

↓

Connect postgres to backend.

↓

Connect redis to backend.

Now imagine the network.

                 backend

          ┌─────────────────────┐

          app

          postgres

          redis

          └─────────────────────┘



Why do production projects define backend?

Because they often have more than one network.

Example:

Internet
      │
      ▼
    Nginx

──────── frontend network ────────

API

──────── backend network ─────────

PostgreSQL

Redis



Notice something important.

The database is not exposed to the internet.

Only the API can reach it.

This is a security benefit.

# What happens after a crash?

Server crashes at 3 AM.

Production shouldn't wait for a human.
so it should start automatically 
restart: unless-stopped



```
What am I deploying?

↓

Go API

↓

Needs PostgreSQL

↓

Needs Redis

↓

Needs persistent storage

↓

Needs networking

↓

Needs restart policy

↓

Now translate those decisions into YAML.

```