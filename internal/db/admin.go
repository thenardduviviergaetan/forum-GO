package forum

import (
	//"database/sql"
	//middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"log"
<<<<<<< HEAD
	"net/http"
=======
	"fmt"
>>>>>>> 257ec57 (work on admin continued)
	//"time"
)

func (app *App_db) AdminHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := app.DB.Query("SELECT id, userstypeid, username, email, validation, time FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var userlst []models.User
	for rows.Next() {
		var user models.User
<<<<<<< HEAD
		err = rows.Scan(user.ID, user.Username, user.Email, user.UserType, user.Validation, user.CreationDate)
		if err != nil {
			log.Fatal(err)
		}
=======
        err = rows.Scan(&user.ID, &user.UserType, &user.Username, &user.Email, &user.Validation, &user.CreationDate)
        if err != nil {
            log.Fatal(err)
        }
>>>>>>> 257ec57 (work on admin continued)
		userlst = append(userlst, user)
	}

	isLogin := false
	if _, err := r.Cookie("session_token"); err == nil {
		isLogin = true
	}

	type Context struct {
<<<<<<< HEAD
		isLogin bool
=======
		isLogin	bool
		userlst []models.User
>>>>>>> 257ec57 (work on admin continued)
	}
	var context Context
	context.isLogin = isLogin
	context.userlst = userlst

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
