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
