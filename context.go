package main

import(
	"database/sql"
	"net/http"
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
    db, err := sql.Open("mysql","root:@/gowiki")
    if err != nil {
        return nil, err
    }

    return &Context {
        Database: db,
    }, nil
}