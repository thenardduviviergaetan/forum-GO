package forum

import (
	"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"net/http"
	"time"
)

// Display the home page handler
func (app *App_db) ForumHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				CheckSessionToken(app.DB, w, r)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	isLogin := false
	if _, err := r.Cookie("session_token"); err == nil {
		isLogin = true
	}

	if err := tmpl.Execute(w, isLogin); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckSessionToken(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	c, err := middle.GetCookie(w, r)
	if err != nil {
		return
	}
	sessionToken := c.Value

	userSession, exists := models.Sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.IsExpired() {
		delete(models.Sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		// return

		var userSessionT string
		err = db.QueryRow("SELECT session_token FROM users WHERE session_token=?", sessionToken).Scan(&userSessionT)
		if err != nil {
			if err != sql.ErrNoRows {
				http.Redirect(w, r, "/logout", http.StatusUnauthorized)
			}
		}
	}
}
