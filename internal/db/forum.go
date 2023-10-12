package forum

import (
	s "forum/sessions"
	"html/template"
	"net/http"
)

// Display the home page handler
func (app *App_db) ForumHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isLogin := false
	if _, err := r.Cookie("session_token"); err == nil {
		isLogin = true
		s.CheckSession(app.DB, w, r)
	}
	if err := tmpl.Execute(w, isLogin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
