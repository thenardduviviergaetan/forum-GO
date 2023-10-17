package forum

import (
	//"database/sql"
	//middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"net/http"
	"log"
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
        err = rows.Scan(user.ID, user.Username, user.Email, user.UserType, user.Validation, user.CreationDate)
        if err != nil {
			//error
            log.Fatal(err)
        }
		userlst = append(userlst, user)
    }


	isLogin := false
	if _, err := r.Cookie("session_token"); err == nil {
		isLogin = true
	}

	type Context struct {
		isLogin	bool
	}
	var context Context
	context.isLogin = isLogin

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}