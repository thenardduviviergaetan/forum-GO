package forum

import (
	"database/sql"
	"fmt"
	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post *models.Post) (int, error) {
	_, err := db.Exec("INSERT INTO post(authorid, author, category, title, content, likes, dislikes, creation) VALUES(?,?,?,?,?,?,?, date())",
		post.AuthorID, post.Author, post.Category, post.Title, post.Content, 0, 0)
	if err != nil {
		fmt.Println("ERROR CREATE POST", err)
	}

	sql, _ := db.Exec("SELECT last_insert_rowid()")
	id, _ := sql.LastInsertId()
	return int(id), nil
}

func UpdatePost(db *sql.DB, id int, post models.Post) {

}
