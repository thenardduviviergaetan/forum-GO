package forum

import (
	"bytes"
	"encoding/json"
	"fmt"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"io"
	"log"
	"net/http"
)

func (app *App_db) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=%s&scope=https%%3A//www.googleapis.com/auth/%s",
		googleClientID,
		"http://localhost:8080/google/callback",
		"token",
		"userinfo.profile")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *App_db) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query().Has("code"))

	token := r.URL.Query().Get("code")
	//TODO check if token is empty

	var resBody bytes.Buffer
	_, errCopy := io.Copy(&resBody, r.Body)
	if errCopy != nil {
		fmt.Println(errCopy)
	}

	fmt.Println("Byte buffer", resBody.String())

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		fmt.Println(err)
	}

	if GoogleOauthTokenRes == nil {
		fmt.Println("Nothing to see")
		return
	}

	tokenBody := &models.GoogleAuthToken{
		Access_token: GoogleOauthTokenRes["access_token"].(string),
		Id_token:     GoogleOauthTokenRes["id_token"].(string),
	}

	fmt.Println(tokenBody.Access_token)

	googleData, dataErr := middle.GetGoogleData(token, "")
	if dataErr != nil {
		fmt.Println("Erreur bro")
		errRedirect(w, r, fmt.Sprint(dataErr))
	}

	app.GoogleSessionHandler(w, r, googleData)
}

func (app *App_db) GoogleSessionHandler(w http.ResponseWriter, r *http.Request, googleData []byte) {
	// var googleUser models.GoogleUser

	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(googleData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprint(w, prettyJSON.String())
}

// func errRedirect(w http.ResponseWriter, r *http.Request, s string) {
// 	http.Redirect(w, r, "/login?error="+url.QueryEscape(s), http.StatusInternalServerError)
// }
