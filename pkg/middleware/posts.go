package forum

import (
	"database/sql"
	"fmt"
	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post models.Post) error {
	_, err := db.Exec("INSERT INTO post(author, category, title, content, like, dislikes) VALUES(?,?,?,?,?,?)",
		post.Author, post.Category, post.Title, post.Content, 0, 0)

	if err != nil {
		fmt.Println("Bonjour", err)
	}

	return nil
}

func UpdatePost(db *sql.DB, id int, post models.Post) {

}
