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