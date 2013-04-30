package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles(viewPath+"edit.html", viewPath+"view.html"))

func renderTemplate(writer http.ResponseWriter, filename string, page *Page) {
    err := templates.ExecuteTemplate(writer, filename+".html", page)
    if err != nil  {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}
