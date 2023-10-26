package forum

import (
	"database/sql"
	//"errors"
	"strconv"
	"net/http"
	"log"
	//"time"
	"fmt"
	models "forum/pkg/models"
)

func FetchFlagedCom(db *sql.DB) []models.Comment {

	rows, err := db.Query("SELECT id, authorid, postid, content, creation, flaged FROM comment WHERE flaged=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
        err = rows.Scan(&comment.ID, &comment.AuthorID, &comment.Postid, &comment.Content, &comment.CreationDate, &comment.Flaged)
        if err != nil {
            log.Fatal(err)
        }
		//get author name and category name
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", comment.AuthorID).Scan(&comment.Author)
		if err != nil {
            log.Fatal(err)
        }
		err = db.QueryRow("SELECT title FROM post WHERE id = ?", comment.Postid).Scan(&comment.Post)
		if err != nil {
            log.Fatal(err)
        }
		comments = append(comments, comment)
	}
	return comments
}

func FetchFlagedPost(db *sql.DB) []models.Post {

	rows, err := db.Query("SELECT id, authorid, content, creation, flaged FROM post WHERE flaged=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
        err = rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.CreationDate, &post.Flaged)
        if err != nil {
            log.Fatal(err)
        }
		//get author name and category name
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", post.AuthorID).Scan(&post.Author)
		if err != nil {
            log.Fatal(err)
        }

		//get the categories
		catrows, erro := db.Query("SELECT categoryid FROM linkcatpost WHERE postid=?", post.ID)
		if erro != nil {
			log.Fatal(err)
		}
		defer catrows.Close()
		for catrows.Next() {
			fmt.Println("nice")
			var catid int
			err = catrows.Scan(&catid)
			if err != nil {
				log.Fatal(err)
			}
			//get the categories names
			temp := ""
			err = db.QueryRow("SELECT title FROM categories WHERE id = ?", catid).Scan(&temp)
			if err != nil {
				log.Fatal(err)
			}
			post.CategoriesName = append(post.CategoriesName, temp)
			post.Categories = append(post.Categories, catid)
		}
		posts = append(posts, post)
	}
	return posts
}

func DelPost(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("delpost"))
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM post WHERE id=?", id)
	if err != nil {
        return err
    }
	return nil
}

func DelPostFlag(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("delpostflag"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE post SET flaged=? WHERE id=?", 0, id)
	if err != nil {
        return err
    }
	return nil
}

func DelCom(db *sql.DB, r *http.Request) error {

	id := 0
	var err error
	if len(r.FormValue("delcom")) == 0 {
		id, err = strconv.Atoi(r.FormValue("delete"))
	} else {
		id, err = strconv.Atoi(r.FormValue("delcom"))
	}
	
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM comment WHERE id=?", id)
	if err != nil {
        return err
    }
	return nil
}

func DelComFlag(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("delcomflag"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE comment SET flaged=? WHERE id=?", 0, id)
	if err != nil {
        return err
    }
	return nil
}
