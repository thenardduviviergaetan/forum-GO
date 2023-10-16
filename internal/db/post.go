package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"log"
	"net/http"
)

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

		if err := tmpl.Execute(w, true); err != nil {
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

		fmt.Println(post.AuthorID)
		fmt.Println(post.Author)
		fmt.Println(post.Category)

		if errCreaPost := middle.CreatePost(app.DB, post); errCreaPost != nil {
			http.Error(w, errCreaPost.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
