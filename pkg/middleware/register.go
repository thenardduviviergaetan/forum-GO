package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
	"net/http"
)

func CheckRegister(db *sql.DB, r *http.Request, user *models.User) error {

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmation := r.FormValue("confirmation")
	if password != confirmation {
		return errors.New("passwords do not match")
	}

	err := db.QueryRow("SELECT email FROM users WHERE email =  ?", email).Scan(&user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return errors.New("email already exist")
	}

	user.Username = username
	user.Email = email
	user.Password = password

	return nil
}
