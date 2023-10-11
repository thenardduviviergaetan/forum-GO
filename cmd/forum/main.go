package main

import (
	"database/sql"
	"fmt"
	. "forum/internal/db"
	. "forum/pkg/handlers"
	. "forum/pkg/middleware"
	"log"
	"net/http"
)

func main() {

	db, err := sql.Open("sqlite3", "./config/db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := InitDB(db)
	if err := app.Migrate(); err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/", ForumHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", LogoutHandler)

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
