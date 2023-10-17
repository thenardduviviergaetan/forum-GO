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

	//check if user is admin
	if cookie, err := r.Cookie("session_token"); err == nil {
		var userstypeid int
		err = app.DB.QueryRow("SELECT userstypeid FROM users WHERE session_token=?", cookie.Value).Scan(&userstypeid)
		if err != nil || userstypeid != 3 {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == "POST" {
		if len(r.FormValue("deletion")) > 0 {
			if err := middle.RmUser(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("delmod")) > 0 {
			if err := middle.Delmod(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("addmod")) > 0 {
			if err := middle.Addmod(app.DB, r); err != nil {
				log.Fatal(err)
			}
		}
	}

	userlst := middle.FetchUsers(app.DB)
	categorylst := middle.FetchCategory(app.DB)
	type Context struct {
		Userlst []models.User
		Categorylst []models.Category
	}
	var context Context
	context.Userlst = userlst
	context.Categorylst = categorylst

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
