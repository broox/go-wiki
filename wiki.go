package main

import (
    "net/http"
    "regexp"
)

const viewPath = "views/"
const dataPath = "data/"
const pathPrefix = "/view/"
const lenPath = len(pathPrefix)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Wrap the CRUD handlers to validate the title in a single place
func makeHandler(handler func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        title := request.URL.Path[lenPath:]
        if !titleValidator.MatchString(title) {
            http.NotFound(writer, request)
            return
        }
        handler(writer, request, title)
    }
}

// 301 root directory requests to FrontPage
func goHome(writer http.ResponseWriter, request *http.Request) {
    http.Redirect(writer, request, "/view/FrontPage", http.StatusFound)
}

func main() {
    http.HandleFunc("/", goHome)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.ListenAndServe(":8080", nil)
}
