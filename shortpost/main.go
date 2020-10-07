package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))
var config = MustLoadConfig()
var db = MustOpenSQL(config.PostgresURL)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "index.html", []Post{
		Post{AuthorName: "Peter", Content: "This is a post."},
		Post{AuthorName: "Peter", Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse efficitur, purus at gravida sagittis, purus mi hendrerit enim, sit amet mattis magna quam sed orci. Nunc nisl erat, ullamcorper ut pulvinar a, fringilla sed est porttitor."},
	})
}

func main() {
	if err := CreatePostTable(db); err != nil {
		log.Fatalf("Failed to create posts table: %s", err.Error())
	}

	fs := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatalf("Server encountered an error: %s", http.ListenAndServe(":8080", nil))
}
