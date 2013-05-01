package main

import (
    "net/http"
    "regexp"
    "github.com/gorilla/mux"
)

const viewPath = "views/"

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Wrap the CRUD handlers to validate the title in a single place
func makeHandler(handler func (http.ResponseWriter, *http.Request, string, *Context)) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        context, err := NewContext(request)
        if err != nil {
            panic(err) // FIXME
        }

        vars := mux.Vars(request)
        title := vars["title"]

        if !titleValidator.MatchString(title) {
            http.NotFound(writer, request)
            return
        }
        handler(writer, request, title, context)
        defer context.Close()
    }
}

// 301 root directory requests to FrontPage
func goHome(writer http.ResponseWriter, request *http.Request) {
    http.Redirect(writer, request, "/view/FrontPage", http.StatusFound)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", goHome)
    router.HandleFunc("/view/{title}", makeHandler(viewHandler)).Methods("GET")
    router.HandleFunc("/edit/{title}", makeHandler(editHandler)).Methods("GET")
    router.HandleFunc("/save/{title}", makeHandler(saveHandler)).Methods("POST")
    if err := http.ListenAndServe(":8080", router); err != nil {
        panic(err)
    }
}
