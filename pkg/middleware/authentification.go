package forum

import (
	"database/sql"
	"errors"
	"net/http"

	models "forum/pkg/models"
	s "forum/sessions"

	"golang.org/x/crypto/bcrypt"
)

// Check user credentials
func Auth(db *sql.DB, w http.ResponseWriter, r *http.Request, user *models.User) error {
	email, password := r.FormValue("email"), r.FormValue("password")

	err := db.QueryRow("SELECT id,username,email, pwd FROM users WHERE email=?", email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return errors.New("invalid email")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	s.SetToken(db, w, r, user)
	return nil
}
