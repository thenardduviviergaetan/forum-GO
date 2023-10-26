package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"net/http"
)

func (app *App_db) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("%s?%s&%s&%s&%s",
		"https://accounts.google.com/o/oauth2/auth",
		"client_id="+googleClientID,
		"response_type=code",
		"redirect_uri=https://localhost:8080/google/callback",
		"scope=profile email")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *App_db) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("code") {
		code := r.URL.Query().Get("code")
		token, token_id, errToken := middle.GetGoogleToken(code, googleClientID, googleClientSecret)
		if errToken != nil {
			fmt.Println("We have an error", token)
			return
		}

		googleData, errData := middle.GetGoogleData(token, token_id)
		switch {
		case errData != nil:
			fmt.Println("error retrieving data", errData)
			return
		case googleData == nil:
			fmt.Println("no data retrieved")
			return
		}

		app.ThirdPartyLoginHandler(w, r, googleData, "google")
	} else {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
}
