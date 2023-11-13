package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"net/http"
	"net/url"
	"text/template"
)

// Connect existing user to forum by comparing their credentials with database
func (app *App_db) LoginHandler(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"web/templates/login.html",
		"web/templates/head.html",
		"web/templates/footer.html",
	)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	errMsg := r.URL.Query().Get("error")
	if r.Method == "POST" {
		// user := &models.User{}
		// if err := middle.Auth(app.DB, w, r, user); err != nil {

		if err := middle.Auth(app.DB, w, r); err != nil {
			if err.Error() == "invalid email" {
				errMsg = "Invalid email address"
			}
			if err.Error() == "invalid password" {
				errMsg = "Invalid password"
			}
			http.Redirect(w, r, "/login?error="+url.QueryEscape(errMsg), http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if err := template.Execute(w, map[string]string{"ErrorMessage": errMsg}); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}

// GithubSessionHandler handles the connection once that everything necessary has been sent by github,
// the function will then unmarshall the data needed for the DataBase.
// If a user exist with the same
func (app *App_db) ThirdPartyLoginHandler(w http.ResponseWriter, r *http.Request, data []byte, loginType string) {
	w.Header().Set("Content-type", "application/json")

	usr, errUnmarshal := unmarshalData(data, loginType)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
		AuthErrRedirect(w, r, fmt.Sprintf("Error trying to login with %s", loginType), "login")
	}

	errCheckRegistered := middle.CheckThirdPartyRegister(app.DB, &usr)

	//Checks for user in the database, if user does not exist then it adds it to the database.
	switch {
	case errCheckRegistered == nil:
		AuthErrRedirect(w, r, "You don't seem to be registered, please first create an account", "login")
		return
	case errCheckRegistered.Error() == "username or email already exist":
		var form = make(url.Values, 0)
		form.Add("password", usr.Password)
		form.Add("email", usr.Email)
		r.URL.RawQuery = form.Encode()
		if err := middle.Auth(app.DB, w, r); err != nil {
			fmt.Println("Error during login ? ", err)
			AuthErrRedirect(w, r, "An error occured while trying to login, please try again", "login")
			return
		}
	default:
		AuthErrRedirect(w, r, "An error occured while trying to login, please try again", "login")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
