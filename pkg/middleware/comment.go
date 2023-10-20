package forum

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	models "forum/pkg/models"
)

func Createcomment(db *sql.DB, comment *models.Comment) (int, error) {
	_, err := db.Exec("INSERT INTO comment(authorid, postid, content,img, creation) VALUES(?,?,?,?, datetime())",
		comment.AuthorID, comment.Postid, comment.Content, "")
	if err != nil {
		fmt.Println("Create comment : ", err)
		return 0, err
	}

	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	comment.ID = id
	return int(id), nil
}

func Updatecomment(db *sql.DB, comment *models.Comment) error {
	_, err := db.Exec("UPDATE comment SET content = ? WHERE id = ?", comment.Content, comment.ID)
	if err != nil {
		fmt.Println("Update comment : ", err)
		return err
	}
	return nil
}

func Reportcomment(db *sql.DB, comment int) error {
	_, err := db.Exec("UPDATE comment SET flaged = ? WHERE id = ?", 1, comment)
	if err != nil {
		fmt.Println("Update comment : ", err)
		return err
	}
	return nil
}

func Removecomment(db *sql.DB, idcomment, currentuser int64) error {
	var idpost int64
	var img string
	err := db.QueryRow("Select postid,img From comment Where id = ?", idcomment).Scan(&idpost, &img)
	if err != nil {
		fmt.Println(err)
	}
	err = os.RemoveAll("web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/comment/" + img)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("DELETE FROM comment WHERE id = ? AND authorid = ? ", idcomment, currentuser)
	if err != nil {
		fmt.Println("Remove comment : ", err)
		return err
	}
	return nil
}

func Updatelike(db *sql.DB, idcomment, iduser int64, like bool) {
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT * FROM linkcomment WHERE userid = ? AND commentid = ?) AS exist", iduser, idcomment).Scan(&exist)
	if err != nil {
		fmt.Println("Update like : ", err)
		return
	}
	if exist {
		var current bool
		err := db.QueryRow("SELECT likes FROM linkcomment WHERE userid = ? AND commentid = ?", iduser, idcomment).Scan(&current)
		if err != nil {
			fmt.Println("Update like : ", err)
			return
		}
		if current == like {
			_, err := db.Exec("DELETE FROM linkcomment WHERE userid = ? AND commentid = ?", iduser, idcomment)
			if err != nil {
				fmt.Println("Update like : ", err)
				return
			}
		} else {
			_, err := db.Exec("UPDATE linkcomment SET likes = ? WHERE userid = ? AND commentid = ?", like, iduser, idcomment)
			if err != nil {
				fmt.Println("Update like : ", err)
				return
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO linkcomment(userid,commentid,likes) VALUES(?,?,?)", iduser, idcomment, like)
		if err != nil {
			fmt.Println("Update like : ", err)
			return
		}
	}
}

func UpdateImgComment(db *sql.DB, idpost, idcomment int64, newimg string) error {
	var lastimg string
	rows, err := db.Query("SELECT img FROM comment WHERE id = ?", idcomment)
	for rows.Next() {
		rows.Scan(&lastimg)
	}
	if err != nil {
		fmt.Println("Update img post err1: ", err)
		return err
	}
	if lastimg != "" {
		err := os.Remove("web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/comment/" + lastimg)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("UPDATE comment SET img = ? WHERE id = ?", newimg, idcomment)
		if err != nil {
			fmt.Println("Update img post err2: ", err)
			return err
		}
	} else {
		_, err := db.Exec("UPDATE comment SET img = ? WHERE id = ?", newimg, idcomment)
		if err != nil {
			fmt.Println("Update img post err3: ", err)
			return err
		}
	}
	return nil
}
