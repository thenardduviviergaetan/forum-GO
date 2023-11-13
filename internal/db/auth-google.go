package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"net/http"
	"strings"
)

var googleClientID = "googleUser"
var googleClientSecret = "googleUserSecret"

// GoogleAuthHandler is the handler for the "login" and "register" page using Google.
// The function redirects to a the API providing a "callback" url that will
// be waiting for the data from github.
func (app *App_db) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	var origin string = strings.TrimRight(strings.ReplaceAll(r.URL.Path, "/google/auth/", ""), "?")

	if r.URL.Path == "/google/auth/login" || r.URL.Path == "/google/auth/register" {
		redirectURL := fmt.Sprintf("%s?%s&%s&%s&%s",
			"https://accounts.google.com/o/oauth2/auth",
			"client_id="+googleClientID,
			"response_type=code",
			"redirect_uri=https://localhost:8080/google/callback/"+origin,
			"scope=profile email")

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		AuthErrRedirect(w, r, "Unauthorised action", origin)
	}
}

// GoogleCallbackHandler is a Handler that is waiting for the answer from Google Login,
// once login is received the function retrieves the user token, then the user data
// and hand all that info to the ThirdPartyLoginHandler.
func (app *App_db) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var origin string = strings.TrimRight(strings.ReplaceAll(r.URL.Path, "/google/callback/", ""), "?")

	if r.URL.Query().Has("code") {
		code := r.URL.Query().Get("code")
		token, token_id, errToken := middle.GetGoogleToken(code, googleClientID, googleClientSecret, origin)
		if errToken != nil {
			AuthErrRedirect(w, r, "Unauthorised action", origin)
			return
		}

		googleData, errData := middle.GetGoogleData(token, token_id)
		switch {
		case errData != nil || googleData == nil:
			AuthErrRedirect(w, r, "Error getting data from Google", origin)
		case origin == "login":
			app.ThirdPartyLoginHandler(w, r, googleData, "google")
		case origin == "register":
			app.ThirdPartyRegisterHandler(w, r, googleData, "google")
		default:
			AuthErrRedirect(w, r, "Unauthorised action", origin)
		}
	} else {
		AuthErrRedirect(w, r, "Unauthorised action", origin)
	}
}
