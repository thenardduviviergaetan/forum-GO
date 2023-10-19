package forum

import (
	//"database/sql"
	//middle "forum/pkg/middleware"
	//models "forum/pkg/models"
	"html/template"
	//"log"
	"net/http"
	//"fmt"
	//"time"
)

func (app *App_db) ModHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/moderation.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//check if user has mod right
	if !app.Data.Admin || !app.Data.Moderator {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}