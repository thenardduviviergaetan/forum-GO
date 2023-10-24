package main

import (
	"crypto/tls"
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
	http.HandleFunc("/com_moderation", app.ComModHandler)
	http.HandleFunc("/profile", app.ProfileHandler)
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", app.LogoutHandler)

	//Post related handlers
	http.HandleFunc("/category", app.CategoryHandler)
	http.HandleFunc("/post/create", app.PostCreateHandler)
	http.HandleFunc("/post", app.PostHandler)
	http.HandleFunc("/post/id", app.PostIdHandler)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         ":8080", // Replace ":8080" with ":443" for production
		Handler:      nil,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	fmt.Println("Listening on port 8080 for development(should be 443 for prod)...")

	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
