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

func LoginErrRedirect(w http.ResponseWriter, r *http.Request, s string) {
	http.Redirect(w, r, "/login?error="+url.QueryEscape(s), http.StatusInternalServerError)
}

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
