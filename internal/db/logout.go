package forum

import (
	"fmt"
	s "forum/sessions"
	"log"
	"net/http"
	"time"
)

func (app *App_db) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := s.GetCookie(w, r)
	if err != nil {
		fmt.Println("error logout")
		return
	}
	sessionToken := cookie.Value

	delete(s.GlobalSessions, sessionToken)

	// Update the user data in the database
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
	http.Redirect(w, r, "/", http.StatusFound)
}
