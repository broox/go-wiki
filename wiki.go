package main

import (
    "net/http"
    "html/template"
    "regexp"
    "errors"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

const viewPath = "views/"
const dataPath = "data/"

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

// Wrap the CRUD handlers to validate the title in a single place
func makeHandler(handler func (http.ResponseWriter, *http.Request, string, *sql.DB)) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        db := openDB()
        title := request.URL.Path[lenPath:]
        if !titleValidator.MatchString(title) {
            http.NotFound(writer, request)
            return
        }
        handler(writer, request, title, db)
        defer db.Close()
    }
}

// add a view to load wiki pages by title at /view/
func viewHandler(writer http.ResponseWriter, request *http.Request, title string, db *sql.DB) {
    page, err := loadPage(title, db)
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

func editHandler(writer http.ResponseWriter, request *http.Request, title string, db *sql.DB) {
    page, err := loadPage(title, db)
    if err != nil {
        page = &Page{ Title: title }
    }
    renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string, db *sql.DB) {
    body := request.FormValue("body")
    p := &Page{ Title: title, Body: []byte(body) }
    err := p.save(db)
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

// 301 root directory requests to FrontPage
func goHome(writer http.ResponseWriter, request *http.Request) {
    http.Redirect(writer, request, "/view/FrontPage", http.StatusFound)
}

// Open a connection to a database for wiki storage
func openDB() (db *sql.DB) {
    db, err := sql.Open("mysql","root:@/gowiki")
    if err != nil {
           panic(err)
    }
    return db
}

func main() {
    http.HandleFunc("/", goHome)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.ListenAndServe(":8080", nil)
}
