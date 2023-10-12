package forum

import (
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"net/http"
	"net/url"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

func (app *App_db) CreateUser(user *models.User) error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(
		"INSERT INTO users(username, password, email) VALUES (?,?,?)",
		user.Username,
		string(hashPass),
		user.Email)

	return err
}

func (app *App_db) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/register.html",
		"web/templates/head.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	errMsg := r.URL.Query().Get("error")
	if r.Method == "POST" {
		user := &models.User{}
		if err := middle.CheckRegister(app.DB, r, user); err != nil {
			if err.Error() == "username or email already exist" {
				errMsg = "Username or Email already exist!"
			}
			if err.Error() == "passwords do not match" {
				errMsg = "Passwords do not match!"
			}
			http.Redirect(w, r, "/register?error="+url.QueryEscape(errMsg), http.StatusFound)
			return
		}
		if err := app.CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if err := tmpl.Execute(w, map[string]string{"ErrorMessage": errMsg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
