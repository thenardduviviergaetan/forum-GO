package forum

import (
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	s "forum/sessions"
	"html/template"
	"net/http"
	"fmt"
)

// PostHandler is a method for the App_db struct that handles HTTP requests related to posts.
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
		fmt.Println("1")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	app.Data.Categories = middle.FetchCat(app.DB, []int{0})
	if r.URL.Query().Has("created") ||
		r.URL.Query().Has("liked") ||
		(r.URL.Query().Has("categories") && r.URL.Query().Get("categories") != "") {
		ApplyFilter(app, w, r)
	} else {
		rows, err := app.DB.Query("SELECT * FROM post")
		if err != nil {
			fmt.Println("2")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		for rows.Next() {
			err := rows.Scan(
				&post.ID,
				&post.AuthorID,
				&post.Author,
				&post.Title,
				&post.Content,
				&post.CreationDate,
				&post.Flaged,
			)
			if err != nil {
				fmt.Println("3")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//get catids from mid table
			catrows, erro := app.DB.Query("SELECT categoryid FROM linkcatpost WHERE postid=?", post.ID)
			if erro != nil {
				fmt.Println("4")
				http.Error(w, erro.Error(), http.StatusInternalServerError)
				return
			}
			for catrows.Next() {
				var catid int
				err = catrows.Scan(&catid)
				if err != nil {
					fmt.Println("5")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				post.Categories = append(post.Categories, catid)
				var catitle string
				err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", catid).Scan(&catitle)
				if err == rows.Err() && err != nil {
					http.Redirect(w, r, "/", http.StatusFound)
				}
				post.CategoriesName = append(post.CategoriesName, catitle)
			}
			post.User_like, post.User_dislike = linkpost(app, post.ID)
			post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
			app.Data.Posts = append(app.Data.Posts, post)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	if err := tmpl.Execute(w, app.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreatedFilter retrieves posts from the database authored by the user with the session token,
// scans the data into a Post model, checks if the post exists in a provided slice of posts (if provided),
// and appends the post to the app's Posts data.
func CreatedFilter(app *App_db, w http.ResponseWriter, r *http.Request, t []models.Post) {
	var post models.Post
	if t != nil {
		app.Data.Posts = nil
	}
	c, _ := r.Cookie("session_token")

	rows, err := app.DB.Query("SELECT * FROM post WHERE authorid = ? ;", s.GlobalSessions[c.Value].UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get catids from mid table
		catrows, erro := app.DB.Query("SELECT categoryid FROM linkcatpost WHERE postid=?", post.ID)
		if erro != nil {
			http.Error(w, erro.Error(), http.StatusInternalServerError)
		}
		for catrows.Next() {
			var catid int
			err = rows.Scan(&catid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			post.Categories = append(post.Categories, catid)
			var catitle string
			err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", catid).Scan(&catitle)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			post.CategoriesName = append(post.CategoriesName, catitle)
		}
		app.Data.Posts = append(app.Data.Posts, post)

		post.User_like, post.User_dislike = linkpost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		if t != nil {
			if middle.HasPost(t, post) {
				app.Data.Posts = append(app.Data.Posts, post)
			}
		} else {
			app.Data.Posts = append(app.Data.Posts, post)
		}
	}
}

// LikedFilter retrieves posts liked by the user from the database, populates their details,
// and appends them to the app's Posts data. If a slice of posts is provided,
// only posts existing in that slice are appended.
func LikedFilter(app *App_db, w http.ResponseWriter, r *http.Request, t []models.Post) {
	var post models.Post
	if t != nil {
		app.Data.Posts = nil
	}
	c, _ := r.Cookie("session_token")

	var tmp int
	rows, err := app.DB.Query("SELECT postid FROM linkpost WHERE userid = ? ;", s.GlobalSessions[c.Value].UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for rows.Next() {
		err := rows.Scan(&tmp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, err := app.DB.Query("SELECT * FROM post WHERE id = ? ;", tmp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for rows.Next() {
			err := rows.Scan(
				&post.ID,
				&post.AuthorID,
				&post.Author,
				&post.Title,
				&post.Content,
				&post.CreationDate,
				&post.Flaged,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//get catids from mid table
			catrows, erro := app.DB.Query("SELECT categoryid FROM linkcatpost WHERE postid=?", post.ID)
			if erro != nil {
				http.Error(w, erro.Error(), http.StatusInternalServerError)
			}
			for catrows.Next() {
				var catid int
				err = rows.Scan(&catid)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				post.Categories = append(post.Categories, catid)
				var catitle string
				err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", catid).Scan(&catitle)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				post.CategoriesName = append(post.CategoriesName, catitle)
			}
			app.Data.Posts = append(app.Data.Posts, post)

			post.User_like, post.User_dislike = linkpost(app, post.ID)
			post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
			if t != nil {
				if middle.HasPost(t, post) {
					app.Data.Posts = append(app.Data.Posts, post)
				}
			} else {
				app.Data.Posts = append(app.Data.Posts, post)
			}
		}
	}
}

// CatFilter retrieves posts from the database based on a given category ID,
// processes each post's data including likes, dislikes, and category title,
// then applies additional filters based on URL query parameters or adds the post to the app's data.
func CatFilter(app *App_db, w http.ResponseWriter, r *http.Request) {
	type PostCatLink struct {
		ID			int
		CategoryId	int
		PostId		int
	}
	var catmp []PostCatLink
	var tmp []models.Post
	cat_id := r.URL.Query().Get("categories")

	if cat_id == "" {
		return
	}
	rows, err := app.DB.Query("SELECT * FROM linkcatpost WHERE categoryid = ?", cat_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	for rows.Next() {
		var postCatLink PostCatLink
		err := rows.Scan(&postCatLink.ID, &postCatLink.CategoryId, &postCatLink.PostId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		catmp = append(catmp, postCatLink)
	}
	if len(catmp) == 0 {
		http.Redirect(w, r, "/post", http.StatusFound)
		return
	}

	var catitle string
	err = app.DB.QueryRow("SELECT title FROM categories WHERE id=?", catmp[0].CategoryId).Scan(&catitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, v := range catmp {
		var post models.Post
		post.Categories = append(post.Categories, v.CategoryId)
		post.CategoriesName = append(post.CategoriesName, catitle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		erro := app.DB.QueryRow("SELECT * FROM post WHERE id=?", v.PostId).Scan(
			&post.ID,
			&post.AuthorID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.CreationDate,
			&post.Flaged,
			)
		if erro != nil {
			http.Error(w, erro.Error(), http.StatusInternalServerError)
		}
		
		post.User_like, post.User_dislike = linkpost(app, post.ID)
		post.Like, post.Dislike = len(post.User_like), len(post.User_dislike)
		tmp = append(tmp, post)
		switch {
		case r.URL.Query().Get("created") == "true":
			CreatedFilter(app, w, r, tmp)
		case r.URL.Query().Get("liked") == "true":
			LikedFilter(app, w, r, tmp)
		default:
			app.Data.Posts = append(app.Data.Posts, post)
		}
	}
}

// ApplyFilter checks the URL query parameters and applies the appropriate
// filter (categories, created, or liked) to the app database.
func ApplyFilter(app *App_db, w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Query().Get("categories") != "":
		CatFilter(app, w, r)

	case r.URL.Query().Get("created") == "true":
		CreatedFilter(app, w, r, nil)

	case r.URL.Query().Get("liked") == "true":
		LikedFilter(app, w, r, nil)
	}
}
