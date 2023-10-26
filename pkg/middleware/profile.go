package forum

import (
	"database/sql"
	models "forum/pkg/models"
	"log"
	"net/http"
)

func FetchUser(db *sql.DB, cookie string) models.User {
	var currentUser models.User
	err := db.QueryRow("SELECT id, user_type_id, username, email, valid, asked_mod, creation FROM users WHERE session_token=?", cookie).Scan(
		&currentUser.ID,
		&currentUser.UserType,
		&currentUser.Username,
		&currentUser.Email,
		&currentUser.Validation,
		&currentUser.AskedMod,
		&currentUser.CreationDate)
	if err != nil {
		log.Fatal(err)
	}
	return currentUser
}

func AskModerator(db *sql.DB, r *http.Request, asked int, id int) error {

	_, err := db.Exec("UPDATE users SET asked_mod=? WHERE id=?", asked, id)
	if err != nil {
		return err
	}
	return nil
}
