# production level folder structure for monolithic in go lang
bookstore/
│
├── cmd/
│   └── server/
│       └── main.go          # Entry point
│
├── internal/
│   ├── book/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   │
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   │
│   ├── auth/
│   ├── order/
│   └── payment/
│
├── pkg/
│   ├── database/
│   ├── logger/
│   ├── config/
│   └── middleware/
│
├── migrations/
│
├── configs/
│
├── api/
│
├── scripts/
│
├── go.mod
└── README.md


You organize by feature.

# flow
Client

↓

Router

↓

Handler

↓

Service

↓

Repository

↓

Database


cmd

# Production apps usually don't have main.go in the root.

Instead:

cmd/

server/

main.go

Why?

Because one project may have several executables.

Example:

cmd/

server/

main.go

worker/

main.go

migration/

main.go

cli/

main.go

One codebase.

Four programs.

# internal

Go has a special folder:

internal/

Packages inside internal cannot be imported from outside the module.

It's Go's built-in encapsulation.

Perfect for application code.


# pkg

pkg usually contains reusable packages.

Example:

pkg/

logger

database

jwt

validator

middleware

These aren't tied to a specific feature.



# packages required while working with mongodb database api  
1. Install the MongoDB driver

The only essential package is:

go get go.mongodb.org/mongo-driver/mongo

This also pulls in the related packages you'll use.

2. Packages you'll commonly import
import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

Let's understand each one.

mongo
"go.mongodb.org/mongo-driver/mongo"

Contains the main MongoDB types:

Client
Collection
Database
CRUD methods (InsertOne, Find, UpdateOne, etc.)
options
"go.mongodb.org/mongo-driver/mongo/options"

Used to configure things.

Example:

clientOptions := options.Client().
    ApplyURI("mongodb://localhost:27017")
bson
"go.mongodb.org/mongo-driver/bson"

MongoDB stores documents as BSON (Binary JSON).

You'll use it for queries:

bson.M{"name": "Go"}

which means

{
  "name": "Go"
}
primitive
"go.mongodb.org/mongo-driver/bson/primitive"

Used for MongoDB's special types.

Most commonly:

primitive.ObjectID

instead of

uint

Example model:

type Book struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name string           `bson:"name" json:"name"`
}
context

Every MongoDB operation requires a context.

Example:

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

collection.InsertOne(ctx, book)

This lets operations time out or be cancelled cleanly.

time

Used with context timeouts:

5 * time.Second
If you're using environment variables (recommended)

Install:

go get github.com/joho/godotenv

Then:

godotenv.Load()

instead of hardcoding:

mongodb://localhost:27017
If you want request validation

A popular package is:

go get github.com/go-playground/validator/v10

Example:

type Book struct {
    Name string `validate:"required"`
}
# Comparison with your MySQL
| MySQL (GORM)  | MongoDB           |
| ------------- | ----------------- |
| `gorm.Open()` | `mongo.Connect()` |
| `*gorm.DB`    | `*mongo.Client`   |
| Table         | Collection        |
| Row           | Document          |
| `db.Create()` | `InsertOne()`     |
| `db.Find()`   | `Find()`          |
| `db.First()`  | `FindOne()`       |
| `db.Save()`   | `UpdateOne()`     |
| `db.Delete()` | `DeleteOne()`     |


# hy do we even need a Handler struct?
In your previous project you had:

func CreateBook(w http.ResponseWriter, r *http.Request) {
    ...
}

It worked perfectly.

So why change it?

in Old approach
Router

↓

CreateBook()

↓

Repository

Every handler was just a standalone function.

What's the problem?

Imagine after a month your project grows.

Your handlers need:

BookService
Logger
Validator
Configuration
Cache

How will CreateBook() access all these?


Option 1:

var service *Service
var logger *Logger
var cache *Cache

Globals.

❌ Bad idea.

Option 2:

Pass everything as parameters.

func CreateBook(
    service *Service,
    logger *Logger,
    cache *Cache,
    w http.ResponseWriter,
    r *http.Request,
)

The router cannot call this. It expects func(http.ResponseWriter, *http.Request) only.

# Production solution

Put the dependencies inside a struct.

type Handler struct {
    service *Service
}

Now the handler owns its dependencies.

The router doesn't call a function.

It calls a method.

h.CreateBook(...)



# Easy way to remember
func Add(a, b int)

No receiver → Function

func (h *Handler) CreateBook()

Has receiver → Method

(h *Handler)
│
├── h           → Receiver variable
├── *Handler    → Receiver type
└── (h *Handler)→ Receiver     

"h is the receiver variable that refers to the Handler object. (h *Handler) is the receiver declaration. Because the function has a receiver, CreateBook is called a method."


# Why no return from the handler methods?

Because the result is written directly to the http.ResponseWriter.

Think of w as an open connection to the client.

Instead of returning a value like this:

book := CreateBook()
return book

you do:

json.NewEncoder(w).Encode(book)

or

http.Error(w, "Book not found", http.StatusNotFound)

The client receives the response through w.


# clean layer architecture
Router knows Handler.
Handler knows Service.
Service knows Repository.
Repository will know MongoDB.

# One important observation

noticed a pattern?

handler.go:

type Handler struct {
	service *Service
}

service.go:

type Service struct {
	repository *Repository
}

Soon, repository.go will look like:

type Repository struct {
	collection *mongo.Collection
}

Each layer owns exactly one dependency—the next layer down.
It's a common design pattern in production Go because it keeps dependencies explicit and makes testing and maintenance much easier.

# decode this func (s *Service) CreateBook(book *Book) (*Book, error)
s is receiver variable 
(s *Service) is receiver
inside this method, I'll refer to the Service object as s.
*Service:Receiver type mean this method belogs to service

(book *Book)

Method parameter.
book:Parameter name.
*Book: parameter type
The caller of this merhod must provide a pointer to a Book.


# to write main.go file
Build from the bottom

Always ask:

What does this object need?

Repository needs

Collection

Service needs

Repository

Handler needs

Service

Router needs

Handler

Server needs

Router

So the order becomes

Mongo Client

↓

Database

↓

Collection

↓

Repository

↓

Service

↓

Handler

↓

Router

↓

HTTP Server

That order almost writes main.go for you.

# env varibales are stored at os level so written in caps to differentiate them with other varibles
os.Getenv("port")

# this is the hierarchy in mongo
Hierarchy:

Mongo Client: is root obj

↓

Database

↓

Collection

↓

Documents


# In real production code, you'll often see service and repository methods accept a context.Context parameter:

func (s *Service) CreateBook(ctx context.Context, book *Book) (*Book, error)

and

func (r *Repository) CreateBook(ctx context.Context, book *Book) (*Book, error)

# What does ctx contain?
1.1. Deadline
2. Cancellation signal
3. values
3. Values

Example:

User ID

Trace ID

Request ID

Authentication data
It contains information about the request

note: The same context travels through every layer.


# Why is context.Background() used in Connect()?

When connecting at application startup, there is no HTTP request.

So you create a fresh root context:

context.Background()

and add a timeout:

context.WithTimeout(...)
During an HTTP request

Later, inside a handler, you won't create a new background context.

Instead you'll use:

ctx := r.Context()

because the request already has one.

Example:

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	books, err := h.service.GetBooks(ctx)
}

Now if the client closes the browser, the request context is canceled automatically, and that cancellation can reach your repository and MongoDB operations.


# Step 1: Create client options
clientOptions := options.Client()

Conceptually:

ClientOptions
│
├── URI = ""
├── MaxPoolSize = default
├── MinPoolSize = default
├── RetryWrites = default
├── TLS = default
└── ...

At this point, nothing is connected yet.

Step 2: Configure the options
clientOptions := options.Client().
    ApplyURI(uri)

Now the configuration becomes:

ClientOptions
│
├── URI = mongodb://localhost:27017
├── MaxPoolSize = default
├── MinPoolSize = default
├── RetryWrites = default
└── ...
Step 3: Connect
client, err := mongo.Connect(ctx, clientOptions)

Now the driver uses those settings to create the client.
# Why not just pass the URI?

Instead of:

mongo.Connect(ctx, uri)

the driver accepts an options object because there are many settings you may want to configure.

For example:

options.Client().
    ApplyURI(uri).
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetRetryWrites(true)

This is much more flexible than adding lots of parameters to mongo.Connect().

# Common client options
| Option             | Purpose                           |
| ------------------ | --------------------------------- |
| `ApplyURI()`       | MongoDB connection string         |
| `SetMaxPoolSize()` | Maximum database connections      |
| `SetMinPoolSize()` | Minimum connections kept open     |
| `SetRetryWrites()` | Retry failed write operations     |
| `SetAppName()`     | Application name shown by MongoDB |
| `SetTLSConfig()`   | Configure TLS/SSL                 |

# Every file starts with:

package book

That means all of these files are compiled together into one package.

Conceptually:

book package
│
├── RegisterRoutes()    ← from routes.go
├── NewHandler()        ← from handler.go
├── NewService()        ← from service.go
├── NewRepository()     ← from repository.go
└── Book               ← from model.go

Notice that the filename disappeared.

# every DB operation should eventually timeout.

Different from startup.

Startup:

10 seconds

CRUD operations:

3–5 seconds

# During an HTTP request

Never create a new background context.

Instead use:

ctx := r.Context()

and pass it through:

Handler

↓

Service

↓

Repository

↓

MongoDB

# `Handler` Type

```go
type Handler struct {
	service *Service
}
```

A `Handler` stores a pointer to a `Service`.

The handler itself does not contain business logic. Instead, it depends on a `Service` to perform business operations. This is an example of **dependency injection**, where the required dependency is provided from outside instead of being created inside the handler.

---

# Constructor

By Go convention, functions whose names start with `New` are constructors. Their job is to create and initialize an object.

```go
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
```

Let's break it down:

### Method parameter

```go
service *Service
```

The constructor expects a pointer to a `Service`. This dependency will be stored inside the `Handler`.

### Return type

```go
*Handler
```

The constructor returns a pointer to a newly created `Handler`.

### Inside the constructor

```go
return &Handler{
	service: service,
}
```

This creates a new `Handler` struct and initializes its `service` field with the `service` parameter passed to the constructor.

Think of it like:

```go
handler := Handler{
	service: service,
}

return &handler
```

The `&` operator returns the memory address (pointer) of the newly created `Handler`.

After construction, the object looks conceptually like this:

```
Handler
│
└── service ─────► Service
```

The handler does **not** contain a copy of the `Service`; it stores only a pointer to the same `Service` instance.

---

# Receiver

```go
func (h *Handler) CreateBook(...)
```

* `h` → Receiver variable
* `*Handler` → Receiver type
* `(h *Handler)` → Receiver declaration

The receiver variable `h` refers to the current `Handler` object.

Since the handler already stores a pointer to a `Service`, every handler method can access it using:

```go
h.service
```

For example:

```go
createdBook, err := h.service.CreateBook(ctx, &book)
```

---

# Complete Flow

```
Repository
      │
      ▼
NewRepository(collection)
      │
      ▼
Service
      │
      ▼
NewService(repository)
      │
      ▼
Handler
      │
      ▼
NewHandler(service)
      │
      ▼
handler.CreateBook()
      │
      ▼
h.service.CreateBook()
```

This pattern repeats throughout the application:

* A **struct** stores its dependencies.
* A **constructor (`New...`)** injects those dependencies.
* The **methods** use the injected dependencies.

# | Function       | Returns                  | Purpose                           |
| -------------- | ------------------------ | --------------------------------- |
| `InsertOne()`  | `*mongo.InsertOneResult` | Insert one document               |
| `Find()`       | `*mongo.Cursor`          | Find multiple documents           |
| `cursor.All()` | `[]Book`                 | Decode all documents into a slice |

# One thing to understand deeply

This line:

cursor, err := r.collection.Find(ctx, bson.M{})

does not return []Book.

It returns a cursor.

Think of a cursor like a database iterator.

MongoDB

↓

Cursor

↓

Document 1

↓

Document 2

↓

Document 3

Then:

cursor.All(ctx, &books)

takes every document from the cursor and decodes them into:

[]Book

That's why Find() and InsertOne() feel different:

InsertOne() performs an action and returns metadata (InsertedID).
Find() returns a stream of matching documents (a cursor), which you then decode into Go structs. This is the standard pattern you'll use for any query that can return multiple documents.
# BSON stands for:

Binary JSON

# HOW TO work with bson
First, what is BSON?

MongoDB does not store JSON internally.

It stores BSON.

BSON stands for:

Binary JSON

Think of it like this:

Go Struct
      ↓
BSON
      ↓
MongoDB

When you insert a book:

book := Book{
    Title:  "Go",
    Author: "Saurabh",
    Year:   2026,
}

The MongoDB driver converts it into BSON.

Conceptually:

{
    "title": "Go",
    "author": "Saurabh",
    "year": 2026
}

The driver handles this conversion automatically because of your bson struct tags:

type Book struct {
    Title string `bson:"title"`
}
What is bson.M?

This is the one you'll use the most.

Definition (simplified):

type M map[string]interface{}

So,

bson.M

is just a shortcut for:

map[string]interface{}
Example
bson.M{
    "title": "Go",
}

means

{
    "title": "Go"
}
Empty bson.M
bson.M{}

means

{}

Empty filter.

Translation:

Match every document.

That's why

collection.Find(ctx, bson.M{})

returns all books.

Think of bson.M as a filter

Collection:

[
  {
    "title":"Go",
    "author":"A"
  },
  {
    "title":"Java",
    "author":"B"
  }
]

Query:

bson.M{
    "title":"Go",
}

MongoDB reads it as:

{
   "title":"Go"
}

Result:

[
  {
      "title":"Go",
      "author":"A"
  }
]
Common bson.M queries
Find all
bson.M{}

↓

{}
Find by title
bson.M{
    "title":"Go",
}

↓

{
   "title":"Go"
}
Find by year
bson.M{
    "year":2026,
}

↓

{
   "year":2026
}
Multiple conditions
bson.M{
    "title":"Go",
    "year":2026,
}

↓

{
    "title":"Go",
    "year":2026
}

Equivalent SQL:

WHERE title='Go'
AND year=2026
What about operators?

MongoDB has operators beginning with $.

Example:

Greater than

bson.M{
    "year": bson.M{
        "$gt":2020,
    },
}

This becomes

{
    "year":{
        "$gt":2020
    }
}

Meaning:

WHERE year > 2020

Less than

bson.M{
    "year": bson.M{
        "$lt":2025,
    },
}

Greater than or equal

bson.M{
    "year": bson.M{
        "$gte":2020,
    },
}

In operator

bson.M{
    "author": bson.M{
        "$in":[]string{
            "A",
            "B",
        },
    },
}

Equivalent SQL

WHERE author IN ('A','B')
Why is it called M?

Because

bson.M

means

Map

MongoDB team chose

M

for

Map
Other BSON types

Besides bson.M, you'll see several related types.

1. bson.M
bson.M{
    "title":"Go",
}

Underlying type

map[string]interface{}

Order is not guaranteed.

Most commonly used.

2. bson.D

Very important.

Definition:

type D []E

where E is:

type E struct {
    Key string
    Value interface{}
}

Example

bson.D{
    {"title","Go"},
    {"year",2026},
}

Same query:

{
    "title":"Go",
    "year":2026
}

But

order is preserved.

When do we use bson.D?

Mostly when order matters.

Example:

Aggregation

Indexes

Sort

Command documents

3. bson.A

"A" means Array.

Example

bson.A{
    "Go",
    "Java",
    "Python",
}

↓

[
   "Go",
   "Java",
   "Python"
]
4. bson.E

Single key-value pair.

bson.E{
    Key:"title",
    Value:"Go",
}

Normally you won't use it directly.

It's mainly used inside bson.D.


# For your bookstore project


// Get all books
bson.M{}

// Get one book by ID
bson.M{
    "_id": objectID,
}

// Find books by author
bson.M{
    "author": "Robert C. Martin",
}

// Update a book
bson.M{
    "$set": bson.M{
        "title": "Clean Code",
    },
}

# mongodb concepts 
InsertOneResult → gives you InsertedID.
Cursor → decode many documents with cursor.All(...).
SingleResult → decode one document with Decode(...).

# in updatebook method of repository Why call GetBookByID() afterwards?

UpdateOne() does not return the updated document.

It returns:

UpdateResult

├── MatchedCount
├── ModifiedCount
└── UpsertedID

It doesn't contain the updated Book.

So a common pattern is:

Update

↓

Get Updated Document

↓

Return Book

# One production note to save an extra query in update situation

Returning r.GetBookByID(ctx, id) after a successful update is a reasonable pattern for many REST APIs because clients often expect the latest state of the resource.

However, it does perform two database operations:

UpdateOne()
FindOne()

If you want to update and return the updated document in a single round trip, MongoDB also provides FindOneAndUpdate(). It's commonly used when the updated document is needed immediately and avoids the extra query.

# DeleteOne() returns:

result, err := collection.DeleteOne(...)

The result contains:

DeleteResult

└── DeletedCount

Example:

DeletedCount = 1

means:

One document was deleted.

If:

DeletedCount = 0

it means:

No document matched the filter.

# | Operation | MongoDB Method | Returns           | Important Result Field          |
| --------- | -------------- | ----------------- | ------------------------------- |
| Create    | `InsertOne()`  | `InsertOneResult` | `InsertedID`                    |
| Get All   | `Find()`       | `Cursor`          | `cursor.All()`                  |
| Get One   | `FindOne()`    | `SingleResult`    | `Decode()`                      |
| Update    | `UpdateOne()`  | `UpdateResult`    | `MatchedCount`, `ModifiedCount` |
| Delete    | `DeleteOne()`  | `DeleteResult`    | `DeletedCount`                  |


# Imagine after one year.

You decide to change MongoDB to PostgreSQL.

If your code looked like:

func CreateBook(w http.ResponseWriter, r *http.Request) {
    collection.InsertOne(...)
}

Now every handler needs to change.

That's a nightmare.

With your architecture:

Handler

↓

Service

↓

Repository

↓

MongoDB

Only the repository changes.

The handler and service don't care whether you're using:

MongoDB
PostgreSQL
MySQL
Redis

They simply call the repository.

Another example

Suppose today you expose a REST API.

Tomorrow the company wants gRPC.

Old design:

REST Handler

↓

MongoDB

You'd rewrite almost everything.

With your design:

REST Handler

↓

Service

↓

Repository

Tomorrow:

gRPC Handler

↓

Service

↓

Repository

The service and repository stay the same.

Only the transport layer changes.

# Company Type	Common Choice of db or ORM?
Startups / MVPs	GORM
Product companies	pgx or database/sql
High-performance systems	pgx
Enterprise applications	Depends on the team's standards

# when changing database from mongo to postgres what i need to change?
Let's assume your application is already running:

Client
   │
   ▼
Router
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
MongoDB

Now your company says:

"We're moving from MongoDB to PostgreSQL."

Most of your code does not change.

Layer by layer
1. Handler

Current:

func (h *Handler) CreateBook(...) {
    h.service.CreateBook(...)
}

Does it know MongoDB exists?

No.

It only knows about the service.

✅ No changes.

2. Service

Current:

func (s *Service) CreateBook(ctx context.Context, book *Book) (*Book, error) {
    return s.repository.CreateBook(ctx, book)
}

Does it know whether the repository uses:

MongoDB?
PostgreSQL?
MySQL?

No.

✅ No changes.

3. Repository

This is where the database-specific code lives.

Current MongoDB repository:

type Repository struct {
    collection *mongo.Collection
}

This changes.

PostgreSQL version might become:

type Repository struct {
    db *sql.DB
}

or if using GORM:

type Repository struct {
    db *gorm.DB
}
Repository methods change

MongoDB:

r.collection.InsertOne(...)

becomes

PostgreSQL:

db.Exec(...)

or

db.QueryRow(...)

or with GORM:

db.Create(book)

Only this layer changes.

4. Database package

Current:

pkg/database/

mongodb.go

After migration:

pkg/database/

postgres.go

Instead of

mongo.Connect(...)

you write

sql.Open(...)

or

gorm.Open(...)
5. Config

Current:

MONGO_URI=...

becomes

POSTGRES_DSN=...

Maybe also:

POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
6. Main

Current:

client, err := database.Connect(cfg.MongoURI)

becomes

db, err := database.Connect(cfg.PostgresDSN)

Then

repo := book.NewRepository(db)

instead of

repo := book.NewRepository(collection)
What doesn't change?

These files stay almost identical:

internal/book/

handler.go      ✅

service.go      ✅

routes.go       ✅

Only:

repository.go

changes.

Visual comparison
MongoDB
Repository

↓

collection.InsertOne()

↓

MongoDB
PostgreSQL
Repository

↓

INSERT INTO books ...

↓

PostgreSQL

Same handler.

Same service.

Different repository.

Even the API stays the same

Your client still sends:

POST /books

Body:

{
  "title":"Go",
  "author":"Saurabh",
  "year":2026
}

The client doesn't even know you've switched databases.

This is called abstraction

The repository hides how the data is stored.

The service only knows:

s.repository.CreateBook(...)

It doesn't know whether that method eventually does:

InsertOne(...)

or

INSERT INTO books ...

or

db.Create(...)
A small improvement many production teams make

Many teams go one step further and introduce an interface.

For example:

type Repository interface {
    CreateBook(ctx context.Context, book *Book) (*Book, error)
    GetBooks(ctx context.Context) ([]Book, error)
    GetBookByID(ctx context.Context, id string) (*Book, error)
    UpdateBook(ctx context.Context, id string, book *Book) (*Book, error)
    DeleteBook(ctx context.Context, id string) error
}

Then the service depends on the interface instead of a concrete MongoDB repository.

Service
    │
    ▼
Repository Interface
    ▲            ▲
    │            │
Mongo Repo   Postgres Repo

With this design, switching databases often means changing only the dependency wiring in main.go:

service := book.NewService(postgresRepo)

instead of:

service := book.NewService(mongoRepo)

The service code itself doesn't change at all.

This is why companies invest in this architecture

You're not writing extra layers just to make the code longer.

You're buying flexibility:

Swap MongoDB for PostgreSQL.
Replace REST with gRPC.
Add caching later.
Mock the repository during tests.

Those changes stay localized because each layer has a single responsibility and only knows about the layer directly beneath it. That's one of the biggest reasons layered architectures remain popular in production systems.