package forum

import (
	"database/sql"
	"errors"
	models "forum/pkg/models"
)

// Prevent duplicate credentials in database during register procedure
func CheckAdminRegister(db *sql.DB, confirmation string, user *models.User) error {

	if confirmation != user.Password {
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
	return nil
}
