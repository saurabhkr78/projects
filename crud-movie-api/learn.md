# this is an nested struct
type movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Year     string    `json:"year"`
	Director *director `json:"director"`
}
type director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
# to start server in go
http.ListenAndServe(":8080", r)

		<host>:<port>
			Listen on all interfaces
			localhost:8080
			127.0.0.1:8080
			192.168.x.x:8080
			your-machine-ip:8080
			Listen only on localhost:http.ListenAndServe("localhost:8080", nil)
			/*
		ListenAndServe asks the operating system to open a TCP port
		(e.g. 8080) and start listening for incoming connections.

		The OS creates and returns a listening socket.

		When a client (browser, curl, mobile app, etc.) connects,
		the server accepts the connection, reads the HTTP request,
		and passes it to the router (ServeMux).

		The router finds the matching handler based on the request
		path and method.

		The handler processes the request and writes a response
		through http.ResponseWriter.

		The server then sends the HTTP response back to the client.
# r:=mux.NewRouter()
mux.NewRouter() creates a router that maps incoming HTTP
requests to handler functions based on URL patterns,
HTTP methods, query parameters, and route variables.
It is an alternative to Go's default ServeMux and provides
more advanced routing capabilities.


you create your own router instance r := mux.NewRouter()
Without a router:
http.HandleFunc("/hello", helloHandler)
http.HandleFunc("/form", formHandler)


What is a Router?

A router's job is:

Incoming Request
      ↓
Check URL Path
      ↓
Find Matching Handler
      ↓
Execute Handler

Example:

r.HandleFunc("/hello", helloHandler)
r.HandleFunc("/movies", getMovies)

Request:

GET /hello

Router says:

"/hello" matches helloHandler

and calls:

helloHandler(w, r)


# Why use Gorilla Mux?
The standard router is simple.But Gorilla Mux adds features.

1.Route Parameters
r.HandleFunc("/movies/{id}", getMovie)

2.Method-based Routing
r.HandleFunc("/movies", getMovies).Methods("GET")
instead of checking for method inside the handler explicitly

3.Query Matching
r.HandleFunc("/search", searchHandler).
	Queries("type", "movie")

# how does response writer and request reader works 
http.Request is a large struct containing headers, URL,body, and other metadata so,Pass the address of the request object, not a copy.Copying all of that for every handler call would be inefficient.

 ResponseWriter is actually an interface, not a struct.Interfaces already contain a reference to the underlying object.

Q. CAN WE MODIFY THE REQUEST OBJECT IN THE HANDLER? 
YES,ParseForm() modifies fields inside the request.If Go passed a copy.changes would affect only the copy.Passing a pointer avoids copying and allows methods to update request state.

# the way to inspect the request body/struct by the client 
r.Method
r.URL.Path
r.URL.Query().Get()

r.ParseForm()
r.FormValue()

r.Header.Get()

r.Body

r.Cookie()
r.Cookies()

r.Context()

r.Host
r.RemoteAddr

r.UserAgent()
r.Referer()

r.FormFile()



# This is the kind of cheat sheet I'd keep while learning Go web development.

# *http.Request Cheat Sheet

```go
func handler(w http.ResponseWriter, r *http.Request)
```

Think:

```text
r = Everything client sent
```

# Method

```go
r.Method
```

Examples:

```go
r.Method == http.MethodGet
r.Method == http.MethodPost
r.Method == http.MethodPut
r.Method == http.MethodDelete
```

Common use:

```go
if r.Method != http.MethodPost {
	http.Error(w, "method not allowed", 405)
	return
}
```

# URL

## Path

```go
r.URL.Path
```

Request:

```text
/users/123
```

Output:

```text
/users/123
```

## Query Parameters

```go
r.URL.Query()
```

Request:

```text
/users?page=2
```

Access:

```go
page := r.URL.Query().Get("page")
```

## Raw Query

```go
r.URL.RawQuery
```

Output:

```text
page=2&limit=10
```

# Header

## Read Header

```go
r.Header.Get("Authorization")
r.Header.Get("Content-Type")
r.Header.Get("User-Agent")
```

## All Headers

```go
fmt.Println(r.Header)
```

# Body

## Read Raw Body

```go
body, err := io.ReadAll(r.Body)
```

## Decode JSON

```go
json.NewDecoder(r.Body).Decode(&data)
```

Very common for APIs.

# Form

## Parse Form

```go
r.ParseForm()
```

Must be called before using `r.Form`.

## Get Single Value

```go
name := r.FormValue("name")
```

Most common.

## All Form Values

```go
fmt.Println(r.Form)
```

## Access Specific Form Key

```go
r.Form["name"]
```

Returns:

```go
[]string
```

# PostForm

POST-only form data.

```go
r.ParseForm()

fmt.Println(r.PostForm)
```

Access:

```go
r.PostForm["email"]
```

# Host

Host requested by client.

```go
r.Host
```

Example:

```text
localhost:8080
```

Useful for:

```text
Multi-tenant apps
Reverse proxies
```

# Remote Address

Client IP.

```go
r.RemoteAddr
```

Output:

```text
192.168.1.10:54321
```

# Cookies

## Get Cookie

```go
cookie, err := r.Cookie("session")
```

Value:

```go
cookie.Value
```

## All Cookies

```go
cookies := r.Cookies()
```

Loop:

```go
for _, c := range cookies {
	fmt.Println(c.Name, c.Value)
}
```

# Context

## Get Context

```go
ctx := r.Context()
```

## Get Value

```go
userID := ctx.Value("userID")
```

Typically set by middleware.

## Deadline

```go
deadline, ok := ctx.Deadline()
```

## Cancellation

```go
select {
case <-ctx.Done():
	fmt.Println("request cancelled")
}
```

Very common in production.

# File Uploads

## Parse Multipart Form

```go
r.ParseMultipartForm(10 << 20)
```

10 MB limit.

## Get Uploaded File

```go
file, header, err := r.FormFile("avatar")
```

## Filename

```go
header.Filename
```

# Useful Helpers

## Content Type

```go
r.Header.Get("Content-Type")
```

## Authorization Token

```go
token := r.Header.Get("Authorization")
```

## User Agent

```go
ua := r.UserAgent()
```

## Referrer

```go
ref := r.Referer()
```

## TLS Information

```go
r.TLS
```

Non-nil if HTTPS.

# how json/encoding works

it has marshal function

```go
json.Marshal(data)
```

then sent backt he encoded data to response writer

```go
w.Write(data)
```

so the shortcut for this is

```go
json.NewEncoder(w).Encode(data)
```

in production while returning the data client need to know which type of data they are receiving. that what header tells

Headers must be sent before the response body.

```go
w.Header().Set(
	"Content-Type",
	"application/json",
)
```

The client may guess:

```text
type/subtype
text/plain
application/octet-stream
unknown
application/xml
application/pdf
image/png
```

Different clients may behave differently.

# var Movie []movie
Movie is a variable that holds a slice of movie structs.

1.Static seed data

var movies = []Movie{
	...
}

2.Runtime additions

movies = append(movies, newMovie)

Rule of thumb:
Known upfront?      → initialize slice literal.
Added later?        → append().

A slice literal is simply a way to create and initialize a slice in one step.
var numbers = []int{1, 2, 3, 4, 5}


# why not to start server first

If you start the server first, `ListenAndServe()` is a blocking function.

It opens the port, starts accepting requests, and then waits indefinitely for incoming requests.

Since it blocks the current goroutine, execution never goes beyond that line until the server stops or returns an error.

As a result, any code written after `ListenAndServe()` will never execute while the server is running.

For example, if route registration, database initialization, or movie data creation happens after `ListenAndServe()`, that code will never run.

Therefore, all application setup should happen before starting the server:
here,
1. Initialize data/database
2. Create router
3. Register routes/handlers
4. Start server using `ListenAndServe()`

# A common startup sequence is:

1. Load configuration
2. Connect to database
3. Initialize repositories/services
4. Create router
5. Register routes
6. Start server

# Meaning of the name “iota” and use of assigning id to movies
“a very small amount”
movie.ID = strconv.Itoa(rand.Intn(1000000))

# Better alternatives
Option 1: UUID 
import "github.com/google/uuid"

movie.ID = uuid.New().String()

Example:

"550e8400-e29b-41d4-a716-446655440000"

Option 2: Database auto-increment

SQL:

id SERIAL PRIMARY KEY

Database generates:

1, 2, 3, 4...


Option 3: Counter (simple in-memory)
var counter int

movie.ID = strconv.Itoa(counter)
counter++

# how in prodcution crete response is returned
1. return the new movie created in the response for confirmation,id,stored values
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
1.2  There are slight variations in production depends on payload size and urgency to see the created object
1. return full object { created movie }
2. minimal response 
{
  "message": "movie created"
}

# * string means
suppose json sends {} and go result " " .

Now you cannot tell:
Did client send empty string ""?
OR
Did client not send FirstName at all? both becomes "".


pointer can be nil but other data types cannot be


| Situation           | Value              |
| ------------------- | ------------------ |
| client sends name   | pointer with value |
| client doesn’t send | nil                |

*string is a pointer to a string, used in APIs to distinguish between “field not provided” (nil) and “field provided with a value”.
* string field points to either address if input provided else to nil.
when you make *string Now you can distinguish:
| Situation      | Value              |
| -------------- | ------------------ |
| not provided   | nil                |
| provided empty | pointer to ""      |
| provided value | pointer to "Nolan" |


# ... is called the variadic expansion operator. 
“Take this slice and pass its elements one by one to the function.”
without(...)
append(a,b) means append b(single value not slice) into a

with(...) append ALL elements inside slice b individually


# Game of [:]
s[:]   - all elements (entire slice)

s[a:b] - elements from index a to b-1

s[:i]  - everything BEFORE i (0 to i-1)

s[i:]  - everything FROM i to end

s[i:j] - elements from i to j-1

s[:len(s)] - entire slice (explicit form)

s[len(s):] - empty slice

s[:i] + s[i+1:]         - remove element at index i

s[i:i]                  - empty slice at position i

append(s[:i], x)        - insert x at index i (basic idea)

s[:]                    - full slice copy view (shares memory)

# {} this is called place holder