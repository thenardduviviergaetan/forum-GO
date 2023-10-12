package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

var GlobalSessions = map[string]Session{}

type Session struct {
	Username string
	UserID   int64
	EndLife  time.Time
}

func IsExpired(s Session) bool {
	return s.EndLife.Before(time.Now())
}

func CheckSession(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("no session token")
		return
	}
	session_token := c.Value

	session := GlobalSessions[session_token]

	var userToken string
	err = db.QueryRow("SELECT session_token FROM users where id = ?", session.UserID).Scan(&userToken)
	if err != nil {
		fmt.Println("no session token in database")
		return
	}
	if userToken != session_token {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if IsExpired(session) {
		// err = db.QueryRow("UPDATE users SET session_token where id = ?", session.UserID)
		fmt.Println("time out")
		stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE id = ?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(nil, session.UserID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
