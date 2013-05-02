package main

import (
    "net/http"
    "regexp"
    "github.com/gorilla/pat"
    "log"
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

        title := request.URL.Query().Get(":title")

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
    config := &Config{}
    if err := config.FromJson("./config.json"); err != nil {
        log.Fatal("Cannot find config file. ", err)
        return
    }

    log.Println("Running Wiki on",config.Port)

    router := pat.New()
    router.HandleFunc("/", goHome)
    router.Get("/view/{title}", makeHandler(viewHandler))
    router.Get("/edit/{title}", makeHandler(editHandler))
    router.Post("/save/{title}", makeHandler(saveHandler))
    if err := http.ListenAndServe(config.Port, router); err != nil {
        log.Fatal("Error starting webserver. ", err)
        return
    }
}
