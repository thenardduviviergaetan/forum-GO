package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"net/http"
	"strings"
)

var githubClientID = "githubUser"
var githubClientSecret = "githubUserSecret"

// GithubAuthHandler is the handler for the "login" and "register" page using Github.
// The function redirects to a the API providing a "callback" url that will
// be waiting for the data from github.
func (app *App_db) GithubAuthHandler(w http.ResponseWriter, r *http.Request) {
	var origin string = strings.TrimRight(strings.ReplaceAll(r.URL.Path, "/github/auth/", ""), "?")

	if r.URL.Path == "/github/auth/login" || r.URL.Path == "/github/auth/register" {
		redirectURL := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
			githubClientID,
			"https://localhost:8080/github/callback/"+origin)

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		AuthErrRedirect(w, r, "Unauthorised action", origin)
	}
}

// GithubCallbackHandler is a Handler that is waiting for the answer from Github Login,
// once login is received the function retrieves the user token, then the user data
// and hand all that info to the ThirdPartyLoginHandler.
func (app *App_db) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var origin string = strings.TrimRight(strings.ReplaceAll(r.URL.Path, "/github/callback/", ""), "?")

	if r.URL.Query().Has("code") {
		code := r.URL.Query().Get("code")
		token, errToken := middle.GetGithubToken(code, githubClientID, githubClientSecret)
		if errToken != nil {
			AuthErrRedirect(w, r, fmt.Sprint(errToken), origin)
			return
		}

		githubData, errData := middle.GetGithubData(token)

		switch {
		case errData != nil || githubData == nil:
			AuthErrRedirect(w, r, "Error getting data from Github", origin)
		case origin == "login":
			app.ThirdPartyLoginHandler(w, r, githubData, "github")
		case origin == "register":
			app.ThirdPartyRegisterHandler(w, r, githubData, "github")
		default:
			AuthErrRedirect(w, r, "Unauthorised action", origin)
		}
	} else {
		AuthErrRedirect(w, r, "Unauthorised action", origin)
	}
}
