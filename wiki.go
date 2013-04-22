package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
)

// A struct to represent a wiki page
type Page struct {
	Title string
	Body []byte
}

// Add the save() function to our Page struct
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage is a function that loads a page by title
// It returns a Page struct, and an optional error
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{ Title: title, Body: body }, nil
}

const pathPrefix = "/view/"
const lenPath = len(pathPrefix)

// add a view to load wiki pages by title at /view/
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[lenPath:]
	page, err := loadPage(title)
	if err != nil {
		// If page can't be found, redirect to the form so we can create it
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[lenPath:]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{ Title: title }
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[lenPath:]
	body := request.FormValue("body")
	p := &Page{ Title: title, Body: []byte(body) }
	err := p.save()
	if err != nil  {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/"+title, http.StatusFound)
}

func renderTemplate(writer http.ResponseWriter, filename string, page *Page) {
	view, err := template.ParseFiles(filename + ".html")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = view.Execute(writer, page)
	if err != nil  {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
        http.HandleFunc(pathPrefix, viewHandler)
        http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
        http.ListenAndServe(":8080", nil)
}
