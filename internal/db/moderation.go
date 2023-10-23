package forum

import (
	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	"html/template"

	//"log"
	"net/http"
	//"time"
)

func (app *App_db) ModHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/moderation.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/post-flaged.html",
		"web/templates/comment-flaged.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	//check if user has mod right
	if !s.GlobalSessions[c.Value].Admin && !s.GlobalSessions[c.Value].Moderator && !s.GlobalSessions[c.Value].Modlight {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	type Context struct {
		Comments  []models.Comment
		Posts     []models.Post
		Connected bool
		Moderator bool
		Modlight  bool
		Admin     bool
	}
	var context Context
	context.Comments = middle.FetchFlagedCom(app.DB)
	context.Posts = middle.FetchFlagedPost(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.Moderator = s.GlobalSessions[c.Value].Modlight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App_db) ComModHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}
	tmpl, err := template.ParseFiles(
		"web/templates/moderation.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/post-flaged.html",
		"web/templates/comment-flaged.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//check if user has mod right
	if !s.GlobalSessions[c.Value].Admin && !s.GlobalSessions[c.Value].Moderator && !s.GlobalSessions[c.Value].Modlight {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	type Context struct {
		Comments  []models.Comment
		Connected bool
		Moderator bool
		Modlight  bool
		Admin     bool
	}
	var context Context
	context.Comments = middle.FetchFlagedCom(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = s.GlobalSessions[c.Value].Moderator
	context.Modlight = s.GlobalSessions[c.Value].Modlight
	context.Admin = s.GlobalSessions[c.Value].Admin

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
