package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
)

func (app *App_db) PosteditHandler(w http.ResponseWriter, r *http.Request, currentuser int64) {
	tmpl, err := template.ParseFiles(
		"web/templates/edit-post.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.Data.Categories = middle.FetchCat(app.DB, int64(app.Data.CurrentPost.Categoryid))

	Returncurentpost(app, w, r, currentuser)
	renderpost_id(w, tmpl, app)
}

func (app *App_db) PostIdHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/post-id.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/comment.html",
		"web/templates/comment-create.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var currentuser int64
	c, _ := r.Cookie("session_token")
	if c != nil {
		currentuser = s.GlobalSessions[c.Value].UserID
	}
	if !Returncurentpost(app, w, r, currentuser) {
		return
	}
	Returncomment(app, currentuser)
	switch r.Method {
	case "POST":
		if app.Data.Connected {
			// create comment
			if r.FormValue("content") != "" {
				var comment models.Comment
				comment.AuthorID = currentuser
				comment.Content = r.FormValue("content")
				comment.Postid = app.Data.CurrentPost.ID
				middle.Createcomment(app.DB, &comment)
			}
			// like comment
			if r.FormValue("like") != "" {
				like := strings.Split(r.FormValue("like"), " ")[0] == "true"
				idcomment, _ := strconv.Atoi(strings.Split(r.FormValue("like"), " ")[1])
				middle.Updatelike(app.DB, int64(idcomment), currentuser, like)
			}
			// like post
			if r.FormValue("like-post") != "" {
				like := strings.Split(r.FormValue("like-post"), " ")[0] == "true"
				idpost, _ := strconv.Atoi(strings.Split(r.FormValue("like-post"), " ")[1])
				middle.Updatelikepost(app.DB, int64(idpost), currentuser, like)
				Returncurentpost(app, w, r, currentuser)
			}
			// delete comment
			if r.FormValue("delete") != "" {
				idcomment, _ := strconv.Atoi(r.FormValue("delete"))
				if s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin || s.GlobalSessions[c.Value].Modlight {
					middle.DelCom(app.DB, r)
				} else {
					middle.Removecomment(app.DB, int64(idcomment), currentuser)
				}
			}
			// edit comment
			if r.FormValue("edit-comment") != "" {
				idcomment, _ := strconv.Atoi(r.FormValue("edit-comment"))
				if s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin || s.GlobalSessions[c.Value].Modlight {
					app.CommentHandler(w, r, int64(idcomment), -1)
				} else {
					app.CommentHandler(w, r, int64(idcomment), currentuser)
				}
				return
			}
			if r.FormValue("comment-editor") != "" {
				var comment models.Comment
				comment.Content = r.FormValue("content-editor")
				id, _ := strconv.Atoi(r.FormValue("comment-editor"))
				comment.ID = int64(id)
				middle.Updatecomment(app.DB, &comment)
			}
			// flag comment
			if r.FormValue("report") != "" {
				middle.FlagComment(app.DB, r)
			}
			// flag post
			if r.FormValue("report-post") != "" {
				middle.FlagPost(app.DB, r)
				http.Redirect(w, r, "id?id="+r.FormValue("report-post"), http.StatusFound)
			}
			if app.Data.CurrentPost.AuthorID == currentuser || s.GlobalSessions[c.Value].Moderator || s.GlobalSessions[c.Value].Admin {
				// delete post
				if r.FormValue("delete-post") != "" {
					// idpost, _ := strconv.Atoi(r.FormValue("delete-post"))
					// middle.RemovePost(app.DB, int64(idpost))
					middle.RemovePost(app.DB, app.Data.CurrentPost.ID)
					http.Redirect(w, r, "/post", http.StatusFound)
				}
				// edit post
				if r.FormValue("edit-post") != "" {
					app.PosteditHandler(w, r, currentuser)
					return
				}
				if r.FormValue("post-editor") != "" {
					var post models.Post
					post.Content = r.FormValue("content-editor")
					// id, _ := strconv.Atoi(r.FormValue("post-editor"))
					// post.ID = int64(id)
					post.ID = app.Data.CurrentPost.ID
					cat, _ := strconv.Atoi(r.FormValue("categories-editor"))
					post.Categoryid = cat
					post.Title = r.FormValue("title-editor")
					middle.UpdatePost(app.DB, &post)
					Returncurentpost(app, w, r, currentuser)
				}
			}
		}
	}
	Returncomment(app, currentuser)
	renderpost_id(w, tmpl, app)
}

func Returncurentpost(app *App_db, w http.ResponseWriter, r *http.Request, currentuser int64) bool {
	var post models.Post
	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "invalid query", http.StatusBadRequest)
			return false
		}

		err = app.DB.QueryRow("SELECT * FROM post where id = ?", id).Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Categoryid,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", post.Categoryid).Scan(&post.Category)

		post.User_like, post.User_dislike = linkpost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		post.Ifcurrentuser = post.AuthorID == currentuser
		app.Data.CurrentPost = post
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No such post", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return false
		}
	}
	return true
}

func renderpost_id(w http.ResponseWriter, tmpl *template.Template, app *App_db) {
	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler that shows the post creation page and ensures that users are certified to create posts.
func (app *App_db) PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Checking for rights to access this page
	cookie, errCookie := r.Cookie("session_token")
	if errCookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	s.CheckActive()
	_, ok := s.GlobalSessions[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		app.Data.Categories = middle.FetchCat(app.DB, 0)

		if err := tmpl.Execute(w, app.Data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "POST":
		var post *models.Post

		errParse := r.ParseForm()
		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusInternalServerError)
			return
		}

		cat, _ := strconv.Atoi(r.FormValue("categories"))
		post = &models.Post{
			Categoryid: cat,
			Title:      r.FormValue("title"),
			Content:    r.FormValue("content"),
		}

		err := app.DB.QueryRow("SELECT id, username FROM users where session_token = ?", cookie.Value).Scan(&post.AuthorID, &post.Author)
		if err != nil {
			log.Fatal(err)
		}

		id, errCreaPost := middle.CreatePost(app.DB, post)
		if errCreaPost != nil {
			http.Error(w, errCreaPost.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post/id?id="+strconv.Itoa(id), http.StatusFound)
	}
}

func linkpost(app *App_db, postid int64) (tablike map[int64]bool, tabdislike map[int64]bool) {
	tablike, tabdislike = make(map[int64]bool), make(map[int64]bool)
	rows, err := app.DB.Query("SELECT userid,likes FROM linkpost WHERE postid = ?", postid)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var userid int64
	var like bool
	for rows.Next() {
		rows.Scan(&userid, &like)
		if like {
			tablike[userid] = true
		} else {
			tabdislike[userid] = true
		}
	}
	return
}
