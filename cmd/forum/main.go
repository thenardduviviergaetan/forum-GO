package main

import (
	"database/sql"
	"fmt"
	. "forum/internal/db"
	"log"
	"net/http"
)

func main() {

	db, err := sql.Open("sqlite3", "config/db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := InitDB(db)
	app.DB.Exec("PRAGMA foreign_keys = ON")
	if err := app.Migrate(); err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/", app.ForumHandler)
	http.HandleFunc("/admin", app.AdminHandler)
	http.HandleFunc("/moderation", app.ModHandler)
	http.HandleFunc("/profile", app.ProfileHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", app.LogoutHandler)

	//Alternative Github Authentication
	http.HandleFunc("/github/auth", app.GithubAuthHandler)
	http.HandleFunc("/github/callback", app.GithubCallbackHandler)

	//Alternative Google Authentication
	http.HandleFunc("/google/auth", app.GoogleAuthHandler)
	http.HandleFunc("/google/callback", app.GoogleCallbackHandler)

	//Post related handlers
	http.HandleFunc("/category", app.CategoryHandler)
	http.HandleFunc("/post/create", app.PostCreateHandler)
	http.HandleFunc("/post", app.PostHandler)
	http.HandleFunc("/post/id", app.PostIdHandler)
	// http.HandleFunc("/post/update", handlers.PostUpdateHandler)
	// http.HandleFunc("/post/delete", handlers.PostDeleteHandler)

	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
