package forum

import (
	"database/sql"
	models "forum/pkg/models"
)

func CreatePost(db *sql.DB, post models.Post) error {
	// _, err := db.Exec("INSERT INTO post(author, category, title, content, )")

	return nil
}
