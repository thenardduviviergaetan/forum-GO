package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB *sql.DB
}

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &App{DB: db}

	http.HandleFunc("/register", app.registerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (app *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := register(app.DB, username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully")
}

func register(db *sql.DB, username, password string) error {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
	if err != nil {
		return err
	}
	statement.Exec()

	statement, err = db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(username, password)
	return err
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	c, err := GetCookie(w, r)
	if err != nil {
		fmt.Fprint(w, "Unauthorized")
		return
	}

	session, err := db.GetSession(c.Value) // Get session from database
	if err != nil || session.IsExpired() {
		fmt.Fprint(w, "Unauthorized")
		return
	}

	// The user is authenticated and the session is valid.
	// You can now handle the request.
	fmt.Fprintf(w, "Hello, %s!", session.Username)
}
