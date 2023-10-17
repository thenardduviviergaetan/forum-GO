package forum

import (
	"fmt"

	models "forum/pkg/models"
)

func Returncomment(app *App_db, currentuser int64) {
	var tab_comment []models.Comment
	var comment models.Comment
	rows, err := app.DB.Query("SELECT id,authorid,postid,content,creation FROM comment WHERE postid = ?", app.Data.CurrentPost.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&comment.ID, &comment.AuthorID, &comment.Postid, &comment.Content, &comment.CreationDate)
		comment.Postid = app.Data.CurrentPost.ID
		err = app.DB.QueryRow("SELECT username FROM users where id = ?", comment.AuthorID).Scan(&comment.Author)
		if err != nil {
			fmt.Println(err)
			return
		}
		comment.Ifcurrentuser = comment.AuthorID == currentuser
		comment.User_like, comment.User_dislike = linkcomment(app, comment.ID)
		comment.Like = len(comment.User_like)
		comment.Dislike = len(comment.User_dislike)
		tab_comment = append(tab_comment, comment)
		comment = models.Comment{}
	}
	app.Data.CurrentPost.Tab_comment = tab_comment
}

func linkcomment(app *App_db, id int64) (tablike map[int64]bool, tabdislike map[int64]bool) {
	tablike, tabdislike = make(map[int64]bool), make(map[int64]bool)
	rows, err := app.DB.Query("SELECT userid,like FROM linkcomment WHERE commentid = ?", id)
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
