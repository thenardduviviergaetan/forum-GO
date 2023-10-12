package main

import (
	"database/sql"
	"fmt"
	. "forum/internal/db"
	. "forum/pkg/handlers"
	handlers "forum/pkg/handlers"
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

	http.HandleFunc("/", handlers.ForumHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", LogoutHandler)

	//Post related handlers
	http.HandleFunc("/post/create", app.PostCreateHandler)
	// http.HandleFunc("/post/update", handlers.PostUpdateHandler)
	// http.HandleFunc("/post/delete", handlers.PostDeleteHandler)

	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
