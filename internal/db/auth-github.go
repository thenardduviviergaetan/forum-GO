package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"net/http"
)

// GithubAuthHandler is the handler for the "login" and "register" page using Github.
// The function redirects to a the API providing a "callback" url that will
// be waiting for the data from github.
func (app *App_db) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"https://localhost:8080/github/callback")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// GithubCallbackHandler is a Handler that is waiting for the answer from Gitlab Login,
// once login is received the function retrieves the user token, then the user data
// and hand all that info to the LoggedInHandler.
func (app *App_db) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, errToken := middle.GetGithubToken(code, githubClientID, githubClientSecret)
	if errToken != nil {
		LoginErrRedirect(w, r, fmt.Sprint(errToken))
		return
	}

	githubData, dataErr := middle.GetGithubData(token)
	if dataErr != nil {
		LoginErrRedirect(w, r, fmt.Sprint(dataErr))
		return
	}

	app.ThirdPartyLoginHandler(w, r, githubData, "github")
}
