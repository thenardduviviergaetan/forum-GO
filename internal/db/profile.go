package forum

import (
	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"log"
	"net/http"
	//"time"
)

func (app *App_db) ProfileHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/profile.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if len(r.FormValue("askmod")) > 0 {
			if err := middle.AskModerator(app.DB, r); err != nil {
				log.Fatal(err)
			}
		}
	}

	type Context struct {
		User models.User
	}
	var context Context
	if cookie, err := r.Cookie("session_token"); err == nil {
		context.User = middle.FetchUser(app.DB, cookie.Value)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}