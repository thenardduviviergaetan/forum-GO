package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
	"net/http"
)

// Prevent duplicate credentials in database during register procedure
func CheckRegister(db *sql.DB, r *http.Request, user *models.User) error {
	username, email := r.FormValue("username"), r.FormValue("email")
	password := r.FormValue("password")
	if password != r.FormValue("confirmation") {
		return errors.New("passwords do not match")
	}

	err := db.QueryRow(
		"SELECT username,email FROM users WHERE username=? OR email=?",
		user.Username,
		user.Email).Scan(&user.Username, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return errors.New("username or email already exist")
	}
	user.Username, user.Email, user.Password = username, email, password
	return nil
}
