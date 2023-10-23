package forum

import (
	//"database/sql"
	middle "forum/pkg/middleware"
	s "forum/sessions"

	//models "forum/pkg/models"
	"html/template"
	"log"
	"net/http"

	//"fmt"
	"strconv"
	//"time"
)

func (app *App_db) CategoryHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/category-create.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	//check if user is admin
	if !s.GlobalSessions[c.Value].Admin {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == "POST" {
		if r.FormValue("creatcat") == "create" {
			if err := middle.AddCategory(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := middle.ModCategory(app.DB, r); err != nil {
				log.Fatal(err)
			}
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
	}

	type Context struct {
		Connected   bool
		Moderator   bool
		Admin       bool
		Modlight    bool
		ID          int
		Title       string
		Description string
	}
	var context Context
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.Modlight = s.GlobalSessions[c.Value].Modlight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "invalid query", http.StatusBadRequest)
			return
		}
		context.ID = id
		title := ""
		description := ""
		err = app.DB.QueryRow("SELECT title, descriptions FROM categories WHERE id=?", id).Scan(&title, &description)
		if err != nil {
			http.Error(w, "invalid query", http.StatusBadRequest)
			return
		}
		context.Title = title
		context.Description = description
	}

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
