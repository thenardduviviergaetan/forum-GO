package forum

import (
	"database/sql"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Auth(db *sql.DB, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var valid_email, valid_password string
	err := db.QueryRow("SELECT email, password FROM users WHERE email=?", email).Scan(&valid_email, &valid_password)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return errors.New("invalid email")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(valid_password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
