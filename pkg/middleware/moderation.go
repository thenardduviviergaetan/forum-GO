package forum

import (
	"database/sql"
	//"errors"
	//"strconv"
	//"net/http"
	"log"
	//"time"
	models "forum/pkg/models"
)

func FetchFlagedCom(db *sql.DB) []models.Comment {

	rows, err := db.Query("SELECT id, authorid, postid, content, creation, flaged FROM comment")
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
	//to finish
	rows, err := db.Query("SELECT id, authorid, postid, content, creation, flaged FROM post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
        err = rows.Scan(&post.ID, &post.AuthorID, &post.Postid, &post.Content, &post.CreationDate, &post.Flaged)
        if err != nil {
            log.Fatal(err)
        }
		posts = append(posts, post)
	}
	return posts
}
