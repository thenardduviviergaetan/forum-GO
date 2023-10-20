package forum

import (
	"database/sql"
	//"errors"
	"strconv"
	"net/http"
	"log"
	//"time"
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
		comments = append(comments, comment)
	}
	return comments
}

func FetchFlagedPost(db *sql.DB) []models.Post {

	rows, err := db.Query("SELECT id, authorid, categoryid, content, creation, flaged FROM post WHERE flaged=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
        err = rows.Scan(&post.ID, &post.AuthorID, &post.Categoryid, &post.Content, &post.CreationDate, &post.Flaged)
        if err != nil {
            log.Fatal(err)
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
