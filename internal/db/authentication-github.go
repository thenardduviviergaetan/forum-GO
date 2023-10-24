package forum

import (
	"encoding/json"
	"fmt"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"net/http"
	"net/url"
	"strconv"
)

// GithubAuthHandler is the handler for the "login" and "register" page using Github.
// Both page are going to the same handler as we juste need to verify that the user can connect to
// a github account. The function redirects to a the appy providing a "callback" url that will
// be waiting for the data from github.
func (app *App_db) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"http://localhost:8080/github/callback")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// GithubCallbackHandler is a Handler that is waiting for the answer from Gitlab Login,
// once login is received the function retrieves retrieves the user token, then the user data
// and hand all that info to the LoggedInHandler.
func (app *App_db) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, errToken := middle.GetGithubToken(code, githubClientID, githubClientSecret)
	if errToken != nil {
		errRedirect(w, r, fmt.Sprint(errToken))
		return
	}

	githubData, dataErr := middle.GetGithubData(token)
	if dataErr != nil {
		errRedirect(w, r, fmt.Sprint(dataErr))
		return
	}

	app.GithubSessionHandler(w, r, githubData)
}

// GithubSessionHandler handles the connection once that everything necessary has been sent by github,
// the function will then unmarshall the data needed for the DataBase.
// If a user exist with the same
func (app *App_db) GithubSessionHandler(w http.ResponseWriter, r *http.Request, githubData []byte) {
	var gitUser models.GithubUser

	if githubData == nil {
		errRedirect(w, r, "Unauthorized access, try login in")
		return
	}

	w.Header().Set("Content-type", "application/json")
	if errUnmarshal := json.Unmarshal(githubData, &gitUser); errUnmarshal != nil {
		errRedirect(w, r, "An error occured while trying to login, please try again")
		return
	}

	usr := models.User{
		UserType: 1,
		Username: gitUser.Login,
		Password: strconv.Itoa(gitUser.ID),
		Email:    fmt.Sprintf("%s@template.github.com", gitUser.Login),
	}
	errCheckAuth := middle.CheckGithubRegister(app.DB, &usr)

	//Checks for user in the database, if user does not exist then it adds it to the database.
	switch {
	case errCheckAuth == nil:
		if errCreate := app.CreateUser(&usr); errCreate != nil {
			errRedirect(w, r, "An error occured while trying to login, please try again")
		}
		fallthrough
	case errCheckAuth.Error() == "username or email already exist":
		if err := middle.AuthGithub(app.DB, w, r, &usr); err != nil {
			errRedirect(w, r, "An error occured while trying to login, please try again")
			return
		}
	default:
		errRedirect(w, r, "An error occured while trying to login, please try again")
		return
	}

	//Final redirect to Main or Profile Page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func errRedirect(w http.ResponseWriter, r *http.Request, s string) {
	http.Redirect(w, r, "/login?error="+url.QueryEscape(s), http.StatusInternalServerError)
}
