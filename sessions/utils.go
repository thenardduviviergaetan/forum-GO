package forum

import (
	"database/sql"
	"fmt"
	. "forum/pkg/models"
)

func Check_PostCreation(db *sql.DB, token string) (Session, bool) {
	CheckActive()
	session, ok := GlobalSessions[token]
	return session, ok
}

func Check_PostModification(db *sql.DB, post Post, token string) bool {
	CheckActive()
	session, ok := GlobalSessions[token]

	if !ok || post.Author != session.Username {
		return false
	}
	return true
}

func Check_PostDeletion(db *sql.DB, post Post, token string) bool {
	//Refresh actual sessions
	CheckActive()
	session, ok := GlobalSessions[token]

	if !ok {
		return false
	}

	var user User
	errQuery := db.QueryRow("SELECT * FROM users where id = ?", session.UserID).Scan(&user)

	if errQuery != nil {
		fmt.Println(errQuery)
	}

	if post.Author != user.Username || user.UserType < 2 {
		return false
	}
	return true
}
