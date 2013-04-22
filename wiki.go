package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

// A struct to represent a wiki page
type Page struct {
	Title string
	Body []byte
}

const viewPath = "views/"
const dataPath = "data/"

// Add the save() function to our Page struct
func (p *Page) save() error {
	filename := dataPath + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage is a function that loads a page by title
// It returns a Page struct, and an optional error
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(dataPath+filename)
	if err != nil {
		return nil, err
	}
	return &Page{ Title: title, Body: body }, nil
}

const pathPrefix = "/view/"
const lenPath = len(pathPrefix)
var templates = template.Must(template.ParseFiles(viewPath+"edit.html", viewPath+"view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Validate a title and return it if valid
func getTitle(writer http.ResponseWriter, request *http.Request) (title string, err error) {
	title = request.URL.Path[lenPath:]
	if !titleValidator.MatchString(title) {
		http.NotFound(writer, request)
		err = errors.New("Invalid Page Title")
	}
	return
}

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

// add a view to load wiki pages by title at /view/
func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		// If page can't be found, redirect to the form so we can create it
		http.Redirect(writer, request, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{ Title: title }
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
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
	err := templates.ExecuteTemplate(writer, filename+".html", page)
	if err != nil  {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
        http.HandleFunc("/view/", makeHandler(viewHandler))
        http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
        http.ListenAndServe(":8080", nil)
}
