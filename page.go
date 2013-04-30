package main

import (
    "io/ioutil"
)

// A struct to represent a wiki page
type Page struct {
    Title string
    Body []byte
}

// Add the save() function to our Page struct
func (p *Page) save() error {
    filename := dataPath + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

// loadPage is a function that loads a page by title
// It returns a Page struct, and an optional error
// Feels like a static method on a Page "object"
// Maybe a PageCollection type?
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(dataPath+filename)
    if err != nil {
        return nil, err
    }
    return &Page{ Title: title, Body: body }, nil
}
