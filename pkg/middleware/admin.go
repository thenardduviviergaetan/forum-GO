package forum

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	//"fmt"
	"strconv"

	//"time"
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

func RmUser(db *sql.DB, r *http.Request) error {
	id, err := strconv.Atoi(r.FormValue("deletion"))
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func Addmod(db *sql.DB, r *http.Request) error {
	id, err := strconv.Atoi(r.FormValue("addmod"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET userstypeid=?, askedmod=? WHERE id=?", 2, 0, id)
	if err != nil {
		return err
	}
	return nil
}

func Delmod(db *sql.DB, r *http.Request) error {
	id, err := strconv.Atoi(r.FormValue("delmod"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET userstypeid=? WHERE id=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

func FetchUsers(db *sql.DB) []models.User {
	rows, err := db.Query("SELECT id, userstypeid, username, email, valide, askedmod, creation FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var userlst []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.UserType, &user.Username, &user.Email, &user.Validation, &user.AskedMod, &user.CreationDate)
		if err != nil {
			log.Fatal(err)
		}
		userlst = append(userlst, user)
	}
	return userlst
}

// func FetchCat(db *sql.DB) []models.Categories {
// 	rows, err := db.Query("SELECT id, title, description FROM categories")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	var categorylst []models.Categories
// 	for rows.Next() {
// 		var categories models.Categories
//         err = rows.Scan(&categories.ID, &categories.Title, &categories.Description)
//         if err != nil {
//             log.Fatal(err)
//         }
// 		categorylst = append(categorylst, categories)
// 	}
// 	return categorylst
// }
