package forum

import (
	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	"html/template"
	"log"
	"net/http"
	"strconv"
	//"fmt"
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

	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	if r.Method == "POST" {
		if len(r.FormValue("askmod")) > 0 {
			id, err := strconv.Atoi(r.FormValue("askmod"))
			if err != nil {
				log.Fatal(err)
			}
			if err := middle.AskModerator(app.DB, r, 1, id); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("asklightmod")) > 0 {
			id, err := strconv.Atoi(r.FormValue("asklightmod"))
			if err != nil {
				log.Fatal(err)
			}
			if err := middle.AskModerator(app.DB, r, 2, id); err != nil {
				log.Fatal(err)
			}
		}
	}

	type Context struct {
		User      models.User
		Connected bool
		Moderator bool
		Modlight  bool
		Admin     bool
		Liked	  models.Liked
		Disliked  models.Disliked
	}
	var context Context
	if cookie, err := r.Cookie("session_token"); err == nil {
		context.User = middle.FetchUser(app.DB, cookie.Value)
		context.Connected = app.Data.Connected
		context.Moderator = s.GlobalSessions[c.Value].Moderator
		context.Modlight = s.GlobalSessions[c.Value].Modlight
		context.Admin = s.GlobalSessions[c.Value].Admin
		context.Liked = middle.FetchLikes(app.DB, int(context.User.ID))
		context.Disliked = middle.FetchDislikes(app.DB, int(context.User.ID))
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	middle.FetchComments(app.DB, int(context.User.ID))

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
