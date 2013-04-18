package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)
	if err != nil {
		fmt.Fprintf(w, "<h1>Error loading page: %s</h1><div>%s</div>", title, err)
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
        http.HandleFunc(pathPrefix, viewHandler)
        http.ListenAndServe(":8080", nil)
}
