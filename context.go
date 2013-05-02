package main

import(
	"database/sql"
	"net/http"
    "log"
    "fmt"
)

// A struct to store things in the context of a request
type Context struct {
    Database *sql.DB
}

// A method on the Context struct to close the DB connection
func (context *Context) Close() {
    context.Database.Close()
}

// Create a new instance of Context with a database
func NewContext(request *http.Request) (*Context, error) {
    connection := fmt.Sprintf("%s:%s@%s/%s", config.Database.User,
                                             config.Database.Password,
                                             config.Database.Host,
                                             config.Database.Name)

    log.Printf("Connecting to the `%s` database", config.Database.Name)
    db, err := sql.Open("mysql",connection)
    if err != nil {
        return nil, err
    }

    return &Context {
        Database: db,
    }, nil
}