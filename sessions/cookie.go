package forum

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	models "forum/pkg/models"

	"github.com/gofrs/uuid"
)

// Set Token and send it to the server session and user cookie
func SetToken(db *sql.DB, w http.ResponseWriter, r *http.Request, user *models.User) {
	sessionToken, _ := uuid.NewV4()
	expiresAt := time.Now().Add(3600 * time.Second)

	GlobalSessions[sessionToken.String()] = Session{
		Username: user.Username,
		UserID:   user.ID,
		EndLife:  expiresAt,
	}

	stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(sessionToken.String(), user.ID)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Expires: expiresAt,
	})
}
