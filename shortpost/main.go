package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))
var config = MustLoadConfig()
var db = MustOpenSQL(config.PostgresURL)

func reverse(slice []interface{}) {
	len := len(slice)
	for i := 0; i < int(len/2); i++ {
		tmp := slice[i]
		slice[i] = slice[len-i]
		slice[len-i] = tmp
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	posts, err := GetAllPost(db)
	if err != nil {
		log.Fatalf("Failed to query posts: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "index.html", posts)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	author := strings.TrimSpace(r.FormValue("author"))
	content := strings.TrimSpace(r.FormValue("content"))

	err := InsertPost(author, content)
	if err != nil {
		log.Fatalf("Failed to insert post: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	if err := CreatePostTable(db); err != nil {
		log.Fatalf("Failed to create posts table: %s", err.Error())
	}

	fs := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatalf("Server encountered an error: %s", http.ListenAndServe(":8080", nil))
}
