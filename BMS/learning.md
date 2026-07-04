# grouped variable declaration
var(
    name string
    id int
)
# net/http receives the HTTP request.
gorilla/mux matches the route

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


# server response have three parts 
an HTTP response as having three parts:

┌────────────────────────────┐
│ 1. Status Line             │
├────────────────────────────┤
│ 2. Headers                 │
├────────────────────────────┤
│ 3. Body                    │
└────────────────────────────┘

1. status line controlled by:

w.WriteHeader(http.StatusOK) or automatically by

w.Write(...)

if you never called WriteHeader().

2. Headers

Example:

Content-Type: application/json
Content-Length: 42

Controlled by:

w.Header().Set("Content-Type", "application/json")

Example:

w.Header().Set("Content-Type", "application/json")
w.Header().Set("Cache-Control", "no-cache")

Result:

Content-Type: application/json
Cache-Control: no-cache

3. Body

Example:

{
    "id":1,
    "name":"Go"
}

Controlled by:

w.Write(...)

or

json.NewEncoder(w).Encode(book)

Notice:

Encode() internally writes to w, so it becomes the response body.

# diff between write() and WriteHeader().

Case 1
w.Write([]byte("Hello"))

Go notices:

"No status code has been sent."

So internally it behaves like:

w.WriteHeader(http.StatusOK)
w.Write([]byte("Hello"))

Response:

HTTP/1.1 200 OK

Hello
Case 2
w.WriteHeader(http.StatusNotFound)
w.Write([]byte("Book not found"))

Response:

HTTP/1.1 404 Not Found

Book not found
What if you do this?
w.Write([]byte("Hello"))

w.WriteHeader(http.StatusNotFound)

Too late!

The first Write() already sent:

HTTP/1.1 200 OK

Go can't change it afterward.

So the response is still:

HTTP/1.1 200 OK

Hello

and the later WriteHeader() has no effect.



# The repository layer error are db related errors does not know http and controller layer errors are error before db call e.g could be a request validtion errror

# there are teo common styles
1. controller fxn doesn't return anything let the function handle everything internally.
func CreateBook(w http.ResponseWriter, r *http.Request)
2. the controller returns which goes to middleware for further action
func CreateBook(w http.ResponseWriter, r *http.Request) error



# When do we use Marshal() else json.NewEncoder(w).Encode(book) is enough ?

When you need the JSON as bytes before sending it somewhere.

Examples:

Save JSON to a file:

data, _ := json.Marshal(book)
os.WriteFile("book.json", data, 0644)

Send to another API:

data, _ := json.Marshal(book)
http.Post(url, "application/json", bytes.NewReader(data))

Store JSON in Redis:

data, _ := json.Marshal(book)
redis.Set(key, data)

In these cases, you need the []byte, so Marshal() is the right tool.

# Why do we convert the ID from string to uint?

HTTP is a text-based protocol, so all values coming from the request URL are received as strings.

Gorilla Mux extracts path variables using:

    vars := mux.Vars(r)

`mux.Vars(r)` returns a map of type:

    map[string]string

For the route:

    /book/{id}

and the request:

    GET /book/5

`mux.Vars(r)` returns:

    map[string]string{
        "id": "5",
    }

Notice that `"5"` is a string, not a uint.

Gorilla Mux does not know whether `{id}` represents:
- an integer
- a username
- a UUID
- a slug

Therefore, it always returns path variables as strings.

Since the repository function expects a numeric ID (e.g., `uint`), the controller must convert the string to the appropriate type before passing it to the repository.

# strconv.ParseUint() expects a string because it's designed to parse text into an unsigned integer.
Its signature is:

func ParseUint(s string, base int, bitSize int) (uint64, error)

Let's decode each parameter.

1. s string

The text you want to convert.

Example:

idStr := "5"

or

idStr := vars["id"]
2. base int

The number system.

Common values:

10 → Decimal (normal numbers)
2 → Binary
8 → Octal
16 → Hexadecimal

For HTTP IDs, use:

10

because:

/book/25

is decimal.

3. bitSize int

The maximum size of the integer.

Examples:

8   // uint8
16  // uint16
32  // uint32
64  // uint64

Most Go code uses:

64

because it gives the largest range.


# Why doesn't ParseUint return uint?

Because uint is platform-dependent:

On a 32-bit system: uint is 32 bits.
On a 64-bit system: uint is 64 bits.

If ParseUint returned uint, its behavior would vary depending on the machine.

Instead, the Go standard library always returns a fixed type:

uint64

This makes it predictable across all platforms.

so that's why while passing to the repository function getbookbyid () we need to convert it to unit from unit64 which returned by the parseUint.




# Every controller repeats this pattern:

vars := mux.Vars(r)

bookIDStr := vars["bookId"]

bookID64, err := strconv.ParseUint(...)

You'll probably write it again in:

GetBookByID
UpdateBook
DeleteBook

That's perfectly fine for now.

Remember what we discussed earlier:

"Write it twice before you abstract it."

You're already on the third use.

At this point, it's reasonable to think about extracting a helper like:

func parseBookID(r *http.Request) (uint, error)



# r := mux.NewRouter()
Intuition
Create an empty router object.

Currently it knows nothing.

No routes.

No handlers.

Just an empty routing engine.

Think of it like creating an empty telephone directory.

Initially:

(empty)


# flow from main.go 
main()

↓

Create router

↓

Register routes

↓

Start HTTP server

↓

Client sends request

↓

Router matches route

↓

Controller executes

↓

Repository executes

↓

Database

↓

Repository

↓

Controller

↓

HTTP Response


# In larger production applications, it's more common for main() to explicitly initialize everything, because it makes startup order obvious and easier to test:

main()

↓

Load config

↓

Connect database

↓

Run migrations

↓

Create repositories

↓

Create services

↓

Create controllers

↓

Register routes

↓

Start server



# the main file 

func main() {
	// Initialize database.
	configs.Connect()
	log.Println("Database connected")

	// Create router.
	r := mux.NewRouter()

	// Register routes.
	routes.RegisterBookStoreRoutes(r)

	// Create HTTP server.
	/*you're creating an HTTP server object.

	Think of it like creating a car before driving it.
	The struct stores all the settings the server will use.
	*/
	server := &http.Server{
		Addr:        ":8080", //Listen on TCP port 8080.
		ReadTimeout: 5 * time.Second,
		/*3. ReadTimeout
		ReadTimeout: 5 * time.Second,

		Imagine a malicious client.

		He connects.

		But instead of sending

		GET /book

		he sends

		G

		waits 30 seconds

		then

		E

		waits another minute...

		This is called a Slowloris attack.

		Without a timeout,

		your server waits forever.

		Eventually,

		thousands of attackers connect.

		Now your server has thousands of stuck connections.

		With

		ReadTimeout: 5 * time.Second,

		you're saying

		If the client doesn't finish sending the request within 5 seconds, disconnect them.
		*/
		WriteTimeout: 10 * time.Second,
		/*
					4. WriteTimeout
			WriteTimeout: 10 * time.Second,

			Suppose your server has already processed the request.

			Now it wants to send

			{
			    ...
			}

			But the client reads extremely slowly.

			Without timeout,

			your server keeps waiting.

			With

			WriteTimeout: 10 * time.Second,

			you're saying

			If I can't finish writing my response within 10 seconds,

			close the connection.
		*/
		IdleTimeout: 60 * time.Second,
		/*
					5. IdleTimeout

			Imagine

			GET /book

			Response sent.

			Connection stays open.

			Nothing happens.

			One minute later...

			Still nothing.

			Eventually you may have

			20,000 idle connections.

			They're doing absolutely nothing.

			This wastes memory.

			So

			IdleTimeout: 60 * time.Second,

			means

			If a connection stays idle for 60 seconds,

			close it.
		*/
		Handler: r, //"Whenever a request comes in, who should handle it?"
		//MaxHeaderBytes:    1 << 20, //1 << 20 is 1 megabyte. This is the maximum size of request headers your server will accept.If someone sends 500 MB headers the server rejects them Protects against attacks.
	}


# Every left shift multiplies by 2.
1 << n = 2ⁿ

Therefore
1 << 20

means

2²⁰

which equals

1,048,576
Why is this used for memory?

Computers measure memory in powers of 2.

1024 bytes = 1 KB

1024 KB = 1 MB

Since

1024 = 2¹⁰

then

1 MB

=

1024 × 1024

=

2¹⁰ × 2¹⁰

=

2²⁰

which is

1 << 20

So

MaxHeaderBytes: 1 << 20,

means

Allow maximum HTTP headers of

1,048,576 bytes

≈ 1 MB