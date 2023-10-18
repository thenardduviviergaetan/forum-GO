package forum

import (
	"database/sql"
	. "forum/pkg/models"

	_ "github.com/mattn/go-sqlite3"
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
			rank INTEGER NOT NULL,
			label TEXT NOT NULL);

		CREATE TABLE IF NOT EXISTS categories(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			time DATETIME NOT NULL   
		);

		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			userstypeid INTEGER NOT NULL,
			username TEXT NOT NULL, 
			password TEXT NOT NULL,
			email TEXT NOT NULL,
			validation INTEGER NOT NULL,
			askedmod INTEGER DEFAULT 0,
			time DATETIME NOT NULL,
			session_token TEXT,
			FOREIGN KEY(userstypeid)REFERENCES userstype(id) ON DELETE CASCADE);
		
		CREATE TABLE IF NOT EXISTS post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			authorid INTEGER NOT NULL,
			author TEXT NOT NULL,
			categoryid INTEGER NOT NULL,
			title TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			like INTEGER NOT NULL,
			dislikes INTEGER NOT NULL,
			creation CURRENT_TIMESTAMP,
			flaged INTEGER DEFAULT 0,
			FOREIGN KEY(categoryid) REFERENCES categories(id) ON DELETE CASCADE,
			FOREIGN KEY(authorid) REFERENCES users(id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS comment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			authorid INTEGER NOT NULL,
			postid INTEGER NOT NULL,
			content TEXT NOT NULL,
			creation CURRENT_TIMESTAMP,
			FOREIGN KEY(authorid) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(postid) REFERENCES post(id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS linkcomment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			commentid INTEGER NOT NULL,
			like BOOLEAN NOT NULL,
			FOREIGN KEY(userid) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(commentid) REFERENCES comment(id) ON DELETE CASCADE
		);
	`
	_, err := app.DB.Exec(query)

	//creation usertype
	var count int
	errChecker := app.DB.QueryRow("SELECT COUNT(*) FROM userstype").Scan(&count)
	if errChecker == sql.ErrNoRows || count == 0 {
		_, err = app.DB.Exec("INSERT INTO userstype(rank, label) VALUES (?,?)", 1, "user")
		if err != nil {
			return err
		}
		_, err = app.DB.Exec("INSERT INTO userstype(rank, label) VALUES (?,?)", 2, "moderator")
		if err != nil {
			return err
		}
		_, err = app.DB.Exec("INSERT INTO userstype(rank, label) VALUES (?,?)", 3, "admin")
		if err != nil {
			return err
		}
	}
	return err
}
