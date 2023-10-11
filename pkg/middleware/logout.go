package forum

import (
	"fmt"
	models "forum/pkg/models"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := GetCookie(w, r)
	if err != nil {
		fmt.Println("error logout")
		return
	}
	sessionToken := cookie.Value
	delete(models.Sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
