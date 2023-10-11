package forum

import (
	models "forum/pkg/models"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// Set Token and send it to the server session and user cookie
func SetToken(w http.ResponseWriter, r *http.Request, username string) {
	sessionToken, _ := uuid.NewV4()
	expiresAt := time.Now().Add(120 * time.Second)

	models.Sessions[sessionToken.String()] = models.Session{
		Username: username,
		EndLife:  expiresAt,
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
