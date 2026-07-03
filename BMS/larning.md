# grouped variable declaration
var(
    name string
    id int
)
# _ means in import
Import this package ONLY for its side effects. Do NOT give me a name to use it.
When Go imports a package, it automatically runs:init() functions inside that package. this is side effect
Does NOT expose any functions or variables to your code.

# What happens if you REMOVE _?
 Go will complain: “imported but not used”
 MySQL driver may NOT register
GORM won’t recognize "mysql"

# var (db *gorm.DB)
Create a global database connection pointer accessible everywhere in this package.
# *gorm.DB
pointer to database connection object
# Why pointer?
Because DB connection is:

heavy object
shared resource
must not be copied

So Go uses: single shared pointer to DB pool

# var db *gorm.DB:
Here variable is of *gorm.DB type

 gorm.DB
gorm → library name
DB → database engine object
A struct provided by GORM
Represents:
database connection
query builder
transaction handler
config state

* means this variable stores the memory address of a DB object, not the object itself

# when you write gorm.Open() :You are getting something that represents:

“Everything needed to talk to a database”
That includes:
connection pool
current config
dialect (mysql/postgres)
logger
query builder state
transaction state



# GORM uses reflection

This is the magic you'll learn later in Go.

GORM looks at:

type Book struct {
	Name        string
	Author      string
	Publication string
	Price       float64
}

and asks:

How many fields?
What are their names?
What are their Go types?
Any tags?
Which field is the primary key?

It discovers something like:

Field: Name
Type: string

Field: Author
Type: string

Field: Price
Type: float64

This process is called reflection.


# GORM maps Go types to SQL types
 | Go Type   | MySQL Type   |
| --------- | ------------ |
| string    | VARCHAR(255) |
| int       | INT          |
| float64   | DOUBLE       |
| bool      | BOOLEAN      |
| time.Time | DATETIME     |
# Name string `json:"name"` tags are used by the encoding/json package to map JSON to Go fields. GORM uses the field names and its own gorm tags for database behavior. If you want to control the database schema (primary keys, column names, indexes, constraints), you'll add gorm tags type Book struct 
{
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"size:100;not null" json:"name"`
}

# gorm:"primaryKey;autoIncrement;not null"
gorm read it as 
primaryKey
autoIncrement
not null
as everything inside the backtick is a string 

# working of auto migration 
Does books table exist?
        │
   Yes ─┴─ No
    │        │
    ▼        ▼
Compare    Create table
schema
    │
    ▼
Add missing columns if needed

# When the application starts, Connect() opens the database and stores the *gorm.DB object. GetDB() returns that same shared object. AutoMigrate() then compares the Book struct with the existing database schema and creates or updates the table if necessary. It does not recreate the table on every startup.

# &Book{}: Create an empty Book value for the struct and give the pointer to that empty book

# since db connection resides inside the model files in this project so we need to have all db fxn like create ,update ,delete etc in the same file for each entity/model


# working of first() to find the item in db 
Give me the first matching row.

When you pass an ID as the second argument, GORM interprets it as:

"Find the row whose primary key equals this ID."
intially existing book is not intialized or empty 
when we pass &existingBook mean Here is the memory location where you should put the data.

# What does First() return?
It returns a *gorm.DB, not a Book.
GORM returns information about the query
contains something like:

Error: nil
RowsAffected: 1

GORM wants to return metadata.

Suppose the book doesn't exist.

How would this work?

book := db.First(...)

What would it return?

An empty Book?

A nil Book?

An error?

Instead, GORM separates them.

The data goes here:

existingBook

The query status goes here:

result

# Find()

Returns all matching rows.


# Production rule
| Method                                        | Best Use               
| --------------------------------------------- | -----------------------
| `db.First(&book, id)`                         | Primary key lookup     
| `db.Where("id = ?", id).First(&book)`         | Works, but more verbose 
| `db.Where("author = ?", author).First(&book)` | Non-primary key lookup 


# An update can be performed directly in the controller using GORM (First(), Save(), or Updates()), but in production it's standard practice to keep database operations in a separate function (repository/model layer). The controller should handle HTTP concerns, while the update function handles database logic.
this is done to minimize the toomany responsibilty on controller.


# Why do we Marshal and Unmarshal JSON?

The database (through GORM) returns data as a Go struct (e.g., Book). A Go struct is an in-memory Go object, not a format that can be sent over HTTP.

HTTP sends bytes, and APIs commonly use JSON as the data format.

So before sending a response to the client, we **marshal** the Go struct into JSON.

Client <----JSON----> Server <----Struct----> Database

Request:
JSON → Unmarshal → Go struct

Response:
Go struct → Marshal → JSON

When you write:

json.Marshal(book)

Go uses reflection to inspect the Book struct at runtime.

type Book struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

It reads the struct tags:

json:"id"
json:"name"

and automatically generates:

{
    "id": 1,
    "name": "Go"
}

You don't manually build the JSON string. The `encoding/json` package handles the conversion for you.