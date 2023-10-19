package forum

import (
	"database/sql"
	"fmt"

	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	_, err := db.Exec("INSERT INTO post(authorid, author, categoryid, title, content, creation) VALUES(?,?,?,?,?, date())",
		post.AuthorID, post.Author, post.Categoryid, post.Title, post.Content)
	if err != nil {
		fmt.Println("ERROR CREATE POST", err)
	}

	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	return int(id), nil
}

func RemovePost(db *sql.DB, idpost int64) error {
	_, err := db.Exec("DELETE FROM post WHERE id = ?", idpost)
	if err != nil {
		fmt.Println("Remove post : ", err)
		return err
	}
	return nil
}

func UpdatePost(db *sql.DB, post *models.Post) error {
	_, err := db.Exec("UPDATE post SET content = ? , title = ? , categoryid = ? WHERE id = ?", post.Content, post.Title, post.Categoryid, post.ID)
	if err != nil {
		fmt.Println("Update comment : ", err)
		return err
	}
	return nil
}

func Updatelikepost(db *sql.DB, idpost, iduser int64, like bool) {
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT * FROM linkpost WHERE userid = ? AND postid = ?) AS exist", iduser, idpost).Scan(&exist)
	if err != nil {
		fmt.Println("Update like post: ", err)
		return
	}
	if exist {
		var current bool
		err := db.QueryRow("SELECT likes FROM linkpost WHERE userid = ? AND postid = ?", iduser, idpost).Scan(&current)
		if err != nil {
			fmt.Println("Update like post: ", err)
			return
		}
		if current == like {
			_, err := db.Exec("DELETE FROM linkpost WHERE userid = ? AND postid = ?", iduser, idpost)
			if err != nil {
				fmt.Println("Update like post: ", err)
				return
			}
		} else {
			_, err := db.Exec("UPDATE linkpost SET likes = ? WHERE userid = ? AND postid = ?", like, iduser, idpost)
			if err != nil {
				fmt.Println("Update like post: ", err)
				return
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO linkpost(userid,postid,likes) VALUES(?,?,?)", iduser, idpost, like)
		if err != nil {
			fmt.Println("Update like post: ", err)
			return
		}
	}
}
