package forum

import (
	"database/sql"
	"fmt"
	"strconv"
	"net/http"
	"log"

	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	_, err := db.Exec("INSERT INTO post(authorid, author, title, content, creation) VALUES(?,?,?,?, datetime())",
		post.AuthorID, post.Author, post.Title, post.Content)
	if err != nil {
		fmt.Println("ERROR CREATE POST", err)
	}
	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	for _, v := range post.Categories {
		_, err = db.Exec("INSERT INTO linkcatpost(categoryid, postid) VALUES(?,?)", v, id)
		if err != nil {
			fmt.Println("ERROR CREATE POST", err)
		}
	}
	return int(id), nil
}

func RemovePost(db *sql.DB, idpost int64) error {
	_, err := db.Exec("DELETE FROM linkcatpost WHERE postid = ?", idpost)
	if err != nil {
		fmt.Println("Remove post : ", err)
		return err
	}
	_, err = db.Exec("DELETE FROM post WHERE id = ?", idpost)
	if err != nil {
		fmt.Println("Remove post : ", err)
		return err
	}
	return nil
}

func UpdatePost(db *sql.DB, post *models.Post) error {
	_, err := db.Exec("UPDATE post SET content = ? , title = ? WHERE id = ?", post.Content, post.Title, post.ID)
	if err != nil {
		fmt.Println("Update comment : ", err)
		return err
	}
	return nil
}

func UpdateCategory(db *sql.DB, post *models.Post) error {

	//delete everything and reinsert everything method
	_, err := db.Exec("DELETE FROM linkcatpost WHERE postid = ?", post.ID)
	if err != nil {
		return err
	}
	for _, v := range post.Categories {
		_, err := db.Exec("INSERT INTO linkcatpost(categoryid, postid) VALUES(?,?)", v, post.ID)
		if err != nil {
			return err
		}
	}
	return err

	//more complicated method and not tested, update the many to many table without deleting everything
	// rows, err := db.Query("SELECT id, categoryid FROM linkcatpost WHERE postid = ?)", post.ID)
	// if err != nil {
	// 	return err
	// }
	// defer rows.Close()
	// todel := []int{}
	// idfound := []int{}
	// for rows.Next() {
	// 	var tempid		int
	// 	var tempcat		int
	// 	found := false
	// 	rows.Scan(&tempid, &tempcat)
	// 	for _, v := range post.Categories {
	// 		if tempcat == v {
	// 			idfound = append(idfound, v)
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		todel = append(todel, tempid)
	// 	}
	// }
	// for _, v := range post.Categories {
	// 	for i, subv := range idfound {
	// 		if v != subv && i == len(idfound)-1 {
	// 			//add category
	// 			_, err := db.Exec("INSERT INTO linkcatpost(categoryid, postid) VALUES(?,?)",
	// 				v, post.ID)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }
	// for _, v := range todel {
	// 	_, err := db.Exec("DELETE FROM linkcatpost WHERE id = ?", v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// return err
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
