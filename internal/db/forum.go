package forum

import (
	"html/template"
	"net/http"

	models "forum/pkg/models"
	s "forum/sessions"
)

// Display the home page handler
func (app *App_db) ForumHandler(w http.ResponseWriter, r *http.Request) {
	app.Data.Posts = nil
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

	app.Data.Connected = func() bool {
		if _, err := r.Cookie("session_token"); err == nil {
			s.CheckSession(app.DB, w, r)
			return true
		}
		s.CheckActive()
		return false
	}()

	GetRecentPosts(app)

	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetRecentPosts(app *App_db) error {
	var post models.Post
	rows, err := app.DB.Query("SELECT * FROM post ORDER BY rowid LIMIT 5")
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Category,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			return err
		}

		app.Data.Posts = append(app.Data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
