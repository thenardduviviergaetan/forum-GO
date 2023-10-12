package forum

import (
	"database/sql"
	"errors"
	"fmt"
	models "forum/pkg/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Check user credentials
func Auth(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var username, valid_email, valid_password string
	var ID int64
	err := db.QueryRow("SELECT id,username,email, password FROM users WHERE email=?", email).Scan(&ID, &username, &valid_email, &valid_password)
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

	SetToken(db, w, r, ID)
	fmt.Println(models.Sessions)
	return nil
}
