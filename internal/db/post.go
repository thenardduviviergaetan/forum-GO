package forum

import (
	"database/sql"
	"fmt"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (app *App_db) PostIdHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	tmpl, err := template.ParseFiles(
		"web/templates/post-id.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "invalid query", http.StatusBadRequest)
			return
		}

		err = app.DB.QueryRow("SELECT * FROM post where id = ?", id).Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Category,
			&post.Title,
			&post.Content,
			&post.Like,
			&post.Dislike,
			&post.CreationDate,
			&post.Flaged,
		)
		app.Data.CurrentPost = post
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No such post", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

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

	rows, err := app.DB.Query("SELECT * FROM post")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Category,
			&post.Title,
			&post.Content,
			&post.Like,
			&post.Dislike,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		app.Data.Posts = append(app.Data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App_db) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	// token := "test_token"
	// if token, err := r.Cookie("session_token"); err == nil {
	// 	http.Redirect(w, r, "http://localhost:8080/", http.StatusUnauthorized)
	// 	return
	// }

	switch r.Method {
	case "GET":
		tmpl, err := template.ParseFiles(
			"web/templates/post-create.html",
			"web/templates/head.html",
			"web/templates/navbar.html",
			"web/templates/footer.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, app.Data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "POST":
		c, _ := r.Cookie("session_token")
		// current := c.Value
		fmt.Println("A post has been submitted")
		errParse := r.ParseForm()

		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
			return
		}

		//TODO check Category to be sure that it exist
		//TODO retrieve user ID to store in the post
		var post *models.Post

		post = &models.Post{
			// AuthorID: current.UserID,
			// Author:   current.Username,
			Category: r.FormValue("categories"),
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
		}

		// query := (`SELECT userid, username FROM users WHERE session_token = ?`)
		err := app.DB.QueryRow("SELECT id, username FROM users where session_token = ?", c.Value).Scan(&post.AuthorID, &post.Author)
		if err != nil {
			log.Fatal(err)
		}

		id, errCreaPost := middle.CreatePost(app.DB, post)
		if errCreaPost != nil {
			http.Error(w, errCreaPost.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post?id="+strconv.Itoa(id), http.StatusFound)
	}
}

func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
