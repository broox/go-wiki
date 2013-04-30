package main

import (
    "io/ioutil"
    "fmt"
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

// Create links out of [PageTitle] text
// FIXME: The output of this escaped to prevent XSS
// We would need to link the titles at the template level rather than on Body so as
// to not unescape other potentially dangerous content
func LinkTitle(bytes []byte) []byte {
    title := bytes[1:len(bytes)-1]
    link := fmt.Sprintf("<a href=\"/view/%s\">%s</a>", title, title)
    bytes = []byte(link)
    return bytes
}