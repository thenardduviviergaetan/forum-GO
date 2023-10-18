package forum

import (
	"database/sql"
	"strconv"
	"net/http"
	"log"
	models "forum/pkg/models"
)

func FetchUser(db *sql.DB, cookie string) models.User {
	var currentUser models.User
	err := db.QueryRow("SELECT id, userstypeid, username, email, valide, askedmod, creation FROM users WHERE session_token=?", cookie).Scan(
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

func AskModerator(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("askmod"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET askedmod=? WHERE id=?", 1, id)
	if err != nil {
        return err
    }
	return nil
}
