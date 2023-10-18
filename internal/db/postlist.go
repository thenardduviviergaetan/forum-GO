package forum

import (
	"html/template"
	"net/http"
	// "fmt"

	models "forum/pkg/models"
	s "forum/sessions"
)

func (app *App_db) PostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	app.Data.Posts = nil
	tmpl, err := template.ParseFiles(
		"web/templates/post.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Has("created") {
		app.ApplyFilter(w, r)
	} else {
		rows, err := app.DB.Query("SELECT * FROM post")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		for rows.Next() {
			err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Categoryid,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", post.Categoryid ).Scan(&post.Category)
		post.User_like, post.User_dislike = linkpost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		app.Data.Posts = append(app.Data.Posts, post)
	}
	
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
	
	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO check for categories and liked filter 
func (app *App_db) ApplyFilter(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	f := r.URL.Query().Get("created")
	
	c, _ := r.Cookie("session_token")
	authid := c.Value

	if f == "true" {
		rows, err := app.DB.Query("SELECT * FROM post WHERE authorid = ? ;", s.GlobalSessions[authid].UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		for rows.Next() {
			err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Categoryid,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", post.Categoryid ).Scan(&post.Category)
		post.User_like, post.User_dislike = linkpost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		app.Data.Posts = append(app.Data.Posts, post)
	}
	}
}