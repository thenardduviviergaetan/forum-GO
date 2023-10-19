package forum

import (
	"database/sql"
	"fmt"
	models "forum/pkg/models"
	"time"
)

func CreatePost(db *sql.DB, post *models.Post) error {
	_, err := db.Exec("INSERT INTO post(authorid, author, category, title, content, like, dislike, creationDate) VALUES(?,?,?,?,?,?,?,?)",
		post.AuthorID, post.Author, post.Category, post.Title, post.Content, 0, 0, time.Now())

	//TODO handle error
	if err != nil {
		fmt.Println("Bonjour", err)
	}

	return nil
}

func UpdatePost(db *sql.DB, id int64, post *models.Post) error {
	_, err := db.Exec("UPDATE post SET title = ?,category = ?,content = ? WHERE id = ?", post.Title, post.Category, post.Content, id)

	if err != nil {
		fmt.Println("Didn't work")
		fmt.Println(err)
	}

	return nil
}

// Takes an ID and finds the correspondinf post in the DataBase.
func RetrievePost(db *sql.DB, id int64) (models.Post, error) {
	post := models.Post{}
	errQuery := db.QueryRow("SELECT creationDate,author,category,title,content,like,dislike,creationDate FROM post WHERE id = ?", id).Scan(
		&post.CreationDate, &post.Author, &post.Category, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.CreationDate)

	return post, errQuery
}
