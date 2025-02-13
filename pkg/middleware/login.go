package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
	s "forum/sessions"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Check user credentials
func Auth(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	var user models.User
	email, password := r.FormValue("email"), r.FormValue("password")

	err := db.QueryRow("SELECT id,username,email,pwd,user_type_id FROM users WHERE email=?", email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.UserType)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return errors.New("Email not found")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("Wrong Password")
	}
	s.SetToken(db, w, r, &user)
	return nil
}
