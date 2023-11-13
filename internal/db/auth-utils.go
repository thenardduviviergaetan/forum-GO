package forum

import (
	"encoding/json"
	"errors"
	"fmt"
	models "forum/pkg/models"
	"net/http"
	"net/url"
	"strings"
)

// LoginErrRedirect is a simple function to redirect every error during login to the main login page
// including in the URL the error as a Query.
func AuthErrRedirect(w http.ResponseWriter, r *http.Request, s string, auth_type string) {
	switch auth_type {
	case "login":
		http.Redirect(w, r, "/login?error="+url.QueryEscape(s), http.StatusSeeOther)
	case "register":
		http.Redirect(w, r, "/register?error="+url.QueryEscape(s), http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// unmarshalData gets the data retrieved from the Login Third Party Query and parse it into
// the according type depending on the "loginType" entered.
func unmarshalData(data []byte, loginType string) (models.User, error) {
	var usr = models.User{UserType: 1}
	var errUnmarshal error

	switch loginType {
	case "google":
		var googleUser models.GoogleUser
		errUnmarshal = json.Unmarshal(data, &googleUser)
		usr.Username = strings.TrimSpace(googleUser.Name)
		usr.Email = googleUser.Email
		usr.Password = googleUser.Email
	case "github":
		var gitUsr models.GithubUser
		errUnmarshal = json.Unmarshal(data, &gitUsr)
		usr.Username = gitUsr.Login
		usr.Email = fmt.Sprintf("%s@template.github.com", gitUsr.Login)
		usr.Password = usr.Email
	default:
		return usr, errors.New("Unrecognized Data Type from external website")
	}

	return usr, errUnmarshal
}
