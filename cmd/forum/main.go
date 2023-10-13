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
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", app.LogoutHandler)

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
