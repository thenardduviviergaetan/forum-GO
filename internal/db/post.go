package forum

import (
	"fmt"
	models "forum/pkg/models"
	"html/template"
	"net/http"
)

func (app *App_db) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	token := "test_token"
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
		fmt.Println("A post has been submitted")
		errParse := r.ParseForm()

		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
		}

		//TODO check Category to be sure that it exist

		post := models.Post{
			Author:   "Tristan",
			Category: r.FormValue("categories"),
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
		}

	}
}

func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
