package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	"log"
	"net/http"
	"time"
)

func (app *App_db) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := middle.GetCookie(w, r)
	if err != nil {
		fmt.Println("error logout")
		return
	}
	sessionToken := cookie.Value

	// delete(models.Sessions, sessionToken)

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
