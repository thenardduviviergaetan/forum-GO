package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	_, err := db.Exec("INSERT INTO post(author_id, author, title, content, creation) VALUES(?,?,?,?, datetime())",
		post.AuthorID, post.Author, post.Title, post.Content)
	if err != nil {
		return -1, err
	}
	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	for _, v := range post.Categories {
		_, err = db.Exec("INSERT INTO link_cat_post(category_id, post_id) VALUES(?,?)", v, id)
		if err != nil {
			fmt.Println("ERROR CREATE POST", err)
		}
	}
	return int(id), nil
}

func RemovePost(db *sql.DB, id_post int64) error {
	_, err := db.Exec("DELETE FROM post WHERE id = ?", id_post)
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
	_, err := db.Exec("DELETE FROM link_cat_post WHERE post_id = ?", post.ID)
	if err != nil {
		return err
	}
	for _, v := range post.Categories {
		_, err := db.Exec("INSERT INTO link_cat_post(category_id, post_id) VALUES(?,?)", v, post.ID)
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

func UpdateLikePost(db *sql.DB, id_post, id_user int64, like bool) error {
	var exist bool
	err := db.QueryRow("SELECT EXISTS( SELECT * FROM link_post WHERE user_id = ? AND post_id = ?) AS exist", id_user, id_post).Scan(&exist)
	if err != nil {
		fmt.Println("Update like post: ", err)
		return err
	}
	if exist {
		var current bool
		err := db.QueryRow("SELECT likes FROM link_post WHERE user_id = ? AND post_id = ?", id_user, id_post).Scan(&current)
		if err != nil {
			fmt.Println("Update like post: ", err)
			return err
		}
		if current == like {
			_, err := db.Exec("DELETE FROM link_post WHERE user_id = ? AND post_id = ?", id_user, id_post)
			if err != nil {
				fmt.Println("Update like post: ", err)
				return err
			}
		} else {
			_, err := db.Exec("UPDATE link_post SET likes = ? WHERE user_id = ? AND post_id = ?", like, id_user, id_post)
			if err != nil {
				fmt.Println("Update like post: ", err)
				return err
			}
		}
	} else {
		_, err := db.Exec("INSERT INTO link_post(user_id,post_id,likes) VALUES(?,?,?)", id_user, id_post, like)
		if err != nil {
			fmt.Println("Update like post: ", err)
			return err
		}
	}
	return nil
}

func FlagPost(db *sql.DB, r *http.Request) error {
	id, _ := strconv.Atoi(r.FormValue("report-post"))
	_, err := db.Exec("UPDATE post SET flagged = ? WHERE id = ?", 1, id)
	if err != nil {
		return err
	}
	return nil
}
