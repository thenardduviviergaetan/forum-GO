package forum

import (
	//"database/sql"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"html/template"
	"log"
	"net/http"
	//"fmt"
	//"time"
)

func (app *App_db) AdminHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"web/templates/admin.html",
		"web/templates/head.html",
		"web/templates/navbar.html",
		"web/templates/footer.html",
		"web/templates/comment-flaged.html",
		"web/templates/post-flaged.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//check if user is admin
	if !app.Data.Admin {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == "POST" {
		if len(r.FormValue("deletion")) > 0 {
			if err := middle.RmUser(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("delmod")) > 0 {
			if err := middle.Delmod(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("addmod")) > 0 {
			if err := middle.Addmod(app.DB, r); err != nil {
				log.Fatal(err)
			}
		// } else if len(r.FormValue("catitle")) > 0 {
		// 	if err := middle.AddCategory(app.DB, r); err != nil {
		// 		log.Fatal(err)
		// 	}
		} else if len(r.FormValue("delcat")) > 0 {
			if err := middle.DelCategory(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("delpost")) > 0 {
			if err := middle.DelPost(app.DB, r); err != nil {
				log.Fatal(err)
			}
		} else if len(r.FormValue("delcom")) > 0 {
			if err := middle.DelCom(app.DB, r); err != nil {
				log.Fatal(err)
			}
		}
	}

	type Context struct {
		Userlst		[]models.User
		Categories	[]models.Categories
		Comments	[]models.Comment
		Posts		[]models.Post
		Connected	bool
		Moderator	bool
		Admin		bool
	}
	var context Context
	context.Userlst = middle.FetchUsers(app.DB)
	context.Categories = middle.FetchCat(app.DB)
	context.Comments = middle.FetchFlagedCom(app.DB)
	context.Posts = middle.FetchFlagedPost(app.DB)
	context.Connected = app.Data.Connected
	context.Moderator = app.Data.Moderator
	context.Admin = app.Data.Admin

	if err := tmpl.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
