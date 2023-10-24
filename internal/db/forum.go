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
		if c, err := r.Cookie("session_token"); err == nil {
			s.CheckSession(app.DB, w, r)
			app.Data.Moderator = s.GlobalSessions[c.Value].Moderator
			app.Data.Admin = s.GlobalSessions[c.Value].Admin
			app.Data.Modlight = s.GlobalSessions[c.Value].Modlight
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
	rows, err := app.DB.Query("SELECT * FROM post ORDER BY rowid DESC LIMIT 5")
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Categoryid1,
			&post.Categoryid2,
			&post.Categoryid3,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			return err
		}
		err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", post.Categoryid1).Scan(&post.Category1)

		app.Data.Posts = append(app.Data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
