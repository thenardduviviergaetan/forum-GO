package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	. "forum/internal/db"
	s "forum/sessions"
	"log"
	"net/http"
	"time"
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

	s.HandleWithLimiter("/", app.ForumHandler, limiter)
	s.HandleWithLimiter("/admin", app.AdminHandler, limiter)
	s.HandleWithLimiter("/moderation", app.ModHandler, limiter)
	s.HandleWithLimiter("/com_moderation", app.ComModHandler, limiter)
	s.HandleWithLimiter("/profile", app.ProfileHandler, limiter)
	s.HandleWithLimiter("/login", app.LoginHandler, limiter)
	s.HandleWithLimiter("/register", app.RegisterHandler, limiter)
	s.HandleWithLimiter("/logout", app.LogoutHandler, limiter)

	//Post related handlers
	s.HandleWithLimiter("/category", app.CategoryHandler, limiter)
	s.HandleWithLimiter("/post/create", app.PostCreateHandler, limiter)
	s.HandleWithLimiter("/post", app.PostHandler, limiter)
	s.HandleWithLimiter("/post/id", app.PostIdHandler, limiter)

	cert := "cert.pem"
	key := "cert-key.pem"

	serverTLSCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Error loading certificate and key: %v", err)
	}

	srv := &http.Server{
		Addr:    ":8080", // Replace ":8080" with ":443" for production
		Handler: nil,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			InsecureSkipVerify:       true,
			Certificates:             []tls.Certificate{serverTLSCert},
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	defer srv.Close()
	fmt.Println("Listening on port 8080 for development(should be 443 for prod)...")
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
