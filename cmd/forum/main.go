package main

import (
	"database/sql"
	"fmt"
	s "forum/sessions"
	"log"
	"net/http"
	"time"

	. "forum/internal/db"
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

	limiter := s.NewBucket(5, time.Second)

	img := http.FileServer(http.Dir("web/static/upload/img"))
	http.Handle("/img/", http.StripPrefix("/img", img))

	s.HandleWithLimiter("/", app.ForumHandler, limiter)
	s.HandleWithLimiter("/admin", app.AdminHandler, limiter)
	s.HandleWithLimiter("/moderation", app.ModHandler, limiter)
	s.HandleWithLimiter("/com_moderation", app.ComModHandler, limiter)
	s.HandleWithLimiter("/profile", app.ProfileHandler, limiter)
	s.HandleWithLimiter("/login", app.LoginHandler, limiter)
	s.HandleWithLimiter("/register", app.RegisterHandler, limiter)
	s.HandleWithLimiter("/logout", app.LogoutHandler, limiter)

	//Alternative Github Authentication
	s.HandleWithLimiter("/github/auth/", app.GithubAuthHandler, limiter)
	s.HandleWithLimiter("/github/callback/", app.GithubCallbackHandler, limiter)
	// http.HandleFunc("/github/auth/", app.GithubAuthHandler)
	// http.HandleFunc("/github/callback/", app.GithubCallbackHandler)

	//Alternative Google Authentication
	s.HandleWithLimiter("/google/auth/", app.GoogleAuthHandler, limiter)
	s.HandleWithLimiter("/google/callback/", app.GoogleCallbackHandler, limiter)
	// http.HandleFunc("/google/auth/", app.GoogleAuthHandler)
	// http.HandleFunc("/google/callback/", app.GoogleCallbackHandler)

	// Post related handlers
	http.HandleFunc("/category", app.CategoryHandler)
	http.HandleFunc("/post/create", app.PostCreateHandler)
	http.HandleFunc("/post", app.PostHandler)
	http.HandleFunc("/post/id", app.PostIdHandler)

	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
