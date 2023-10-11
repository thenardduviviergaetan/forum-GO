package forum

import (
	middle "forum/pkg/middleware"
	"net/http"
	"net/url"
	"text/template"
)

func (app *App_db) LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/login.html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := r.URL.Query().Get("error")
	if r.Method == "POST" {
		// user := &models.User{}
		if err := middle.Auth(app.DB, r); err != nil {
			if err.Error() == "invalid email" {
				errMsg = "Invalid email address"
			}
			if err.Error() == "invalid password" {
				errMsg = "Invalid password"
			}
			http.Redirect(w, r, "/login?error="+url.QueryEscape(errMsg), http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
	if err := tmpl.Execute(w, map[string]string{"ErrorMessage": errMsg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
