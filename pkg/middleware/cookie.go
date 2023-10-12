package forum

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// Set Token and send it to the server session and user cookie
func SetToken(db *sql.DB, w http.ResponseWriter, r *http.Request, ID int64) {
	sessionToken, _ := uuid.NewV4()
	expiresAt := time.Now().Add(30 * time.Second)

	// models.Sessions[sessionToken.String()] = models.Session{
	// 	Username: username,
	// 	EndLife:  expiresAt,
	// }

	// Update the user data in the database
	stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(sessionToken.String(), ID)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiresAt,
	})
}

// Return token from cookie or error if token does not exist
func GetCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	c, err := r.Cookie(("session_token"))
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	return c, nil
}
