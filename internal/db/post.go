package forum

import (
	"fmt"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	session "forum/sessions"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//TODO add defaults to all switch cases with error handling

// PostCreateHandler is the handler that is responsable for the creation of new posts. It will first ensure that a user is logged in
// and if so will first display a screen where the user can enter data and once it is sent to the webstie the data will be enterred
// in the DataBase.
func (app *App_db) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Checking for rights to access this page
	cookie, errCookie := r.Cookie("session_token")
	if errCookie != nil {
		postError(fmt.Errorf("not user detected"), false, w, r)
		return
	}

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

		postTmpl := models.PostTemplates{
			IsSigned: true,
		}

		if err := tmpl.Execute(w, postTmpl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "POST":
		//Checking for rights to create a post.
		session, ok := session.Check_PostCreation(app.DB, cookie.Value)

		if !ok {
			postError(fmt.Errorf("no rights to interact"), true, w, r)
			return
		}

		errParse := r.ParseForm()
		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
			return
		}

		//TODO check if any value from the form is null
		if emptyCheck(r.PostForm) {
			fmt.Println("Cannot use empty values in a POST")
			return
		}

		post := &models.Post{
			AuthorID: session.UserID,
			Author:   session.Username,
			Category: r.FormValue("categories"),
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
		}

		err := app.DB.QueryRow("SELECT id, username FROM users where session_token = ?", cookie.Value).Scan(&post.AuthorID, &post.Author)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(post.AuthorID)
		fmt.Println(post.Author)
		fmt.Println(post.Category)

		//Storing data in the DataBase.
		if errCreaPost := middle.CreatePost(app.DB, post); errCreaPost != nil {
			http.Error(w, errCreaPost.Error(), http.StatusInternalServerError)
			return
		}
	default:
		postError(fmt.Errorf("illegal actions"), true, w, r)
		return
	}
}

// PostUpdateHandler is the handler that update the posts on the forum. It will first check if a user is logged and if it has the rights to modify
// the post accesed via the URL and if so it will update the DataBase with the new data given by the user.
func (app *App_db) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//Checking for rights to access this page.
	cookie, errCookie := r.Cookie("session_token")
	if errCookie != nil {
		postError(fmt.Errorf("not user detected"), false, w, r)
		return
	}

	//Retrieving ID from URL.
	url := r.URL.String()
	tmpID := string(url[strings.LastIndex(url, "/")+1:])
	postId, errConv := strconv.ParseInt(tmpID, 10, 64)
	if len(url) == len(tmpID) || errConv != nil {
		postError(fmt.Errorf("incorrect post ID"), true, w, r)
		return
	}

	//Retrieves the post related to the ID to find the actual user rights on it.
	post, errQuery := middle.RetrievePost(app.DB, postId)
	if errQuery != nil {
		postError(fmt.Errorf("incorrect post ID"), true, w, r)
		return
	} else if !session.Check_PostModification(app.DB, post, cookie.Value) {
		postError(fmt.Errorf("no rights to interact"), true, w, r)
		return
	}

	switch r.Method {
	case "GET":
		var postTmpl = models.PostTemplates{
			Post:     post,
			IsSigned: true,
		}

		tmpl, errParse := template.ParseFiles(
			"web/templates/post-create.html",
			"web/templates/head.html",
			"web/templates/navbar.html",
			"web/templates/footer.html",
		)

		if err := tmpl.Execute(w, postTmpl); err != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		errParse := r.ParseForm()
		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
			return
		}

		//TODO check if any value from the form is null
		if emptyCheck(r.PostForm) {
			fmt.Println("Cannot use empty values in a POST")
			return
		}

		newPost := models.Post{
			Category: r.FormValue("category"),
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
		}

		err := middle.UpdatePost(app.DB, postId, &newPost)

		if err != nil {
			fmt.Println(err)
		}
	default:
		postError(fmt.Errorf("illegal actions"), true, w, r)
		return
	}
}

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	default:
	}
}

func emptyCheck(form url.Values) bool {
	for _, key := range form {
		for _, val := range key {
			if val == "" {
				return true
			}
		}
	}

	return false
}

// Error management, each error is linked to a specific message to briefly explain to the user what might have gone wrong.
func postError(err error, isSigned bool, w http.ResponseWriter, r *http.Request) {
	tmpl, errParse := template.ParseFiles(
		"web/templates/post-error.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)

	postTmpl := models.PostTemplates{IsSigned: isSigned}

	switch err.Error() {
	case "not user detected":
		postTmpl.Err = "You don't seem to be logged in"
	case "no rights to interact":
		postTmpl.Err = "You don't have the rights to modify this post"
	case "incorrect post ID":
		postTmpl.Err = "The post you're trying to find doesn't seem to exist"
	case "illegal actions":
		postTmpl.Err = "Action not allowed"
	}

	if errParse != nil {
		fmt.Println("Error while Parsing templates for postError")
		return
	}

	//TODO find a solution to the IsSigned problem : have to create a struct for the HEADER
	//to be able to send the data anyway and have other types of data sent as well.
	if err := tmpl.Execute(w, postTmpl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
