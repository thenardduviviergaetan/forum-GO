package forum

import (
	"database/sql"
	//"strconv"
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

func AskModerator(db *sql.DB, r *http.Request, asked int, id int) error {

	_, err := db.Exec("UPDATE users SET askedmod=? WHERE id=?", asked, id)
	if err != nil {
        return err
    }
	return nil
}
