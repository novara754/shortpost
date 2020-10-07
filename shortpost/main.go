package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/go-sessions/v3"
	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/login.html"))
var config = MustLoadConfig()
var db = MustOpenSQL(config.PostgresURL)
var sess = sessions.New(sessions.Config{
	Cookie:  "shortpost",
	Expires: 2 * time.Hour,
})

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

	session := sess.Start(w, r)
	var currentUser *User = nil
	if currentUserID, err := session.GetInt64("current-user-id"); err == nil {
		currentUser, err = GetUserByID(db, currentUserID)
	}

	data := struct {
		Posts       []Post
		CurrentUser *User
	}{
		Posts:       posts,
		CurrentUser: currentUser,
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	authorID, err := sess.Start(w, r).GetInt64("current-user-id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	err = InsertPost(authorID, content)
	if err != nil {
		log.Fatalf("Failed to insert post: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session := sess.Start(w, r)

	if r.Method == "GET" {
		if session.Get("current-user-id") != nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		errors := r.URL.Query()["error"]
		templates.ExecuteTemplate(w, "login.html", struct {
			Action string
			Error  bool
		}{
			Action: "login",
			Error:  len(errors) == 1 && errors[0] == "1",
		})
	} else if r.Method == "POST" {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		user, err := AuthenticateUser(username, password)

		if err == sql.ErrNoRows || err == bcrypt.ErrMismatchedHashAndPassword {
			http.Redirect(w, r, "/login?error=1", http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Fatalf("Failed to authenticate user: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Set("current-user-id", user.ID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	session := sess.Start(w, r)

	if r.Method == "GET" {
		if session.Get("current-user-id") != nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		errors := r.URL.Query()["error"]
		templates.ExecuteTemplate(w, "login.html", struct {
			Action string
			Error  bool
		}{
			Action: "register",
			Error:  len(errors) == 1 && errors[0] == "1",
		})
	} else if r.Method == "POST" {
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		user, err := InsertUser(username, password)
		if err != nil {
			log.Fatalf("Failed to insert user: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Set("current-user-id", user.ID)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := sess.Start(w, r)
	if session.Get("current-user-id") != nil {
		session.Destroy()
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	if err := CreateUserTable(db); err != nil {
		log.Fatalf("Failed to create users table: %s", err.Error())
	}

	if err := CreatePostTable(db); err != nil {
		log.Fatalf("Failed to create posts table: %s", err.Error())
	}

	fs := http.FileServer(http.Dir("public/"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new-post", newPostHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatalf("Server encountered an error: %s", http.ListenAndServe(":8080", nil))
}
