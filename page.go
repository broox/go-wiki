package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// A struct to represent a wiki page
type Page struct {
    Title string
    Body []byte
}

// Add the save() function to our Page struct
func (p *Page) save(db *sql.DB) error {
    existingPage, err := loadPage(p.Title, db)
    if err != nil {
        return err
    }

    p.Body = []byte("bleh")

    if existingPage == nil {
        insert, err := db.Prepare("INSERT INTO pages (title, body, created_at) VALUES (?,?,NOW())")
        if err != nil {
            return err
        }
        defer insert.Close()
        _, err = insert.Exec(p.Title, p.Body)
    } else {
        update, err := db.Prepare("UPDATE pages SET body = ?, updated_at = NOW() WHERE title = ?")
        if err != nil {
            return err
        }
        defer update.Close()
        _, err = update.Exec(p.Body, p.Title)
    }
    return err
}

// loadPage is a function that loads a page by title
// It returns a Page struct, and an optional error
// Feels like a static method on a Page "object"
// Maybe a PageCollection type?
func loadPage(title string, db *sql.DB) (*Page, error) {
    query, err := db.Prepare("SELECT title, body FROM pages WHERE title = ?")
    if err != nil {
        return nil, err
    }

    var body []byte

    err = query.QueryRow(title).Scan(&title, &body)
    if err != nil {
        return nil, err
    }
    return &Page{ Title: title, Body: body }, nil
}