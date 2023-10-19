package forum

import (
	s "forum/sessions"
	"log"
	"net/http"
	"time"
)

// Disconnect user, remove session and clear cookies before redirecting to home page
func (app *App_db) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	sessionToken := cookie.Value
	delete(s.GlobalSessions, sessionToken)

	stmt, err := app.DB.Prepare("UPDATE users SET session_token = NULL WHERE session_token = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(sessionToken)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	//TODO change status to allow instant redirect
	http.Redirect(w, r, "/", http.StatusFound)
}
