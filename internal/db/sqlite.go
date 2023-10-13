package forum

import (
	"database/sql"
	"errors"

	. "forum/pkg/models"

	"github.com/mattn/go-sqlite3"
)

type App_db struct {
	DB *sql.DB
}

func InitDB(db *sql.DB) *App_db {
	return &App_db{DB: db}
}

func (app *App_db) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			username TEXT NOT NULL, 
			password TEXT NOT NULL,
			email TEXT NOT NULL,
<<<<<<< Updated upstream
=======
			validation INTEGER NOT NULL,
			time DATETIME NOT NULL,
>>>>>>> Stashed changes
			session_token TEXT);

		CREATE TABLE IF NOT EXISTS post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			category TEXT NOT NULL,
			title TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			like INTEGER NOT NULL,
			dislikes INTEGER NOT NULL);
	`
	_, err := app.DB.Exec(query)
	return err
}

func (app *App_db) Create(post Post) (*Post, error) {
	res, err := app.DB.Exec("INSERT INTO post(author, category, title, content, like, dislike) VALUES (?,?,?,?,?,?)")
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, err
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	post.ID = id

	return &post, nil
}
