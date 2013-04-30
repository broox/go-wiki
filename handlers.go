package main

import (
    "net/http"
    "regexp"
)

// add a view to load wiki pages by title at /view/
func viewHandler(writer http.ResponseWriter, request *http.Request, title string, context *Context) {
    page, err := loadPage(title, context.Database)
    if err != nil {
        // If page can't be found, redirect to the form so we can create it
        http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
        return
    }

    r := regexp.MustCompile("\\[([a-zA-Z]+)\\]")
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    page.Body = r.ReplaceAllFunc(page.Body, LinkTitle)
    renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string, context *Context) {
    page, err := loadPage(title, context.Database)
    if err != nil {
        page = &Page{ Title: title }
    }
    renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string, context *Context) {
    body := request.FormValue("body")
    p := &Page{ Title: title, Body: []byte(body) }
    err := p.save(context.Database)
    if err != nil  {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}