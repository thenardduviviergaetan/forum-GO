package forum

import (
	"database/sql"
	"errors"
	. "forum/pkg/models"

	"github.com/mattn/go-sqlite3"
)

type App_db struct {
	DB   *sql.DB
	Data Data
}

func InitDB(db *sql.DB) *App_db {
	return &App_db{DB: db}
}

func (app *App_db) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS userstype(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			rank TEXT NOT NULL);

		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			userstypeid INTEGER NOT NULL,
			username TEXT NOT NULL, 
			password TEXT NOT NULL,
			email TEXT NOT NULL,
			validation INTEGER NOT NULL,
			time DATETIME NOT NULL,
			session_token TEXT,
			FOREIGN KEY(userstypeid)REFERENCES userstype(id) ON DELETE CASCADE);
		
			CREATE TABLE IF NOT EXISTS post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			authorid INTEGER NOT NULL,
			author TEXT NOT NULL,
			category TEXT NOT NULL,
			title TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			like INTEGER NOT NULL,
			dislikes INTEGER NOT NULL,
			creation CURRENT_TIMESTAMP,
			FOREIGN KEY(authorid)REFERENCES users(id) ON DELETE CASCADE);
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
