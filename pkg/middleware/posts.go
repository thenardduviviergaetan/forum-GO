package forum

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	_, err := db.Exec("INSERT INTO post(authorid, author, categoryid, img, title, content, creation) VALUES(?,?,?,?,?,?, date())",
		post.AuthorID, post.Author, post.Categoryid, "", post.Title, post.Content)
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
	err = os.RemoveAll("web/static/upload/img/post" + strconv.Itoa(int(idpost)))
	if err != nil {
		fmt.Println(err)
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

func UpdateImgPoste(db *sql.DB, idpost int64, newimg string) error {
	var lastimg string
	rows, err := db.Query("SELECT img FROM post WHERE id = ?", idpost)
	for rows.Next() {
		rows.Scan(&lastimg)
	}
	if err != nil {
		fmt.Println("Update img post err1: ", err)
		return err
	}
	if lastimg != "" {
		err := os.Remove("web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/" + lastimg)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("UPDATE post SET img = ? WHERE id = ?", newimg, idpost)
		if err != nil {
			fmt.Println("Update img post err2: ", err)
			return err
		}
	} else {
		_, err := db.Exec("UPDATE post SET img = ? WHERE id = ?", newimg, idpost)
		if err != nil {
			fmt.Println("Update img post err3: ", err)
			return err
		}
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

func FlagPost(db *sql.DB, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("report-post"))
	_, err := db.Exec("UPDATE post SET flaged = ? WHERE id = ?", 1, id)
	if err != nil {
		log.Fatal(err)
	}
}
