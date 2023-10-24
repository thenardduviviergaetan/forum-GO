package forum

import (
	"database/sql"
	//"strconv"
	"net/http"
	"log"
	"sort"
	"fmt"
	models "forum/pkg/models"
)

func FetchUser(db *sql.DB, cookie string) models.User {
	var currentUser models.User
	err := db.QueryRow("SELECT id, userstypeid, username, email, valide, askedmod, creation FROM users WHERE session_token=?", cookie).Scan(
						&currentUser.ID, 
						&currentUser.UserType, 
						&currentUser.Username, 
						&currentUser.Email, 
						&currentUser.Validation, 
						&currentUser.AskedMod, 
						&currentUser.CreationDate)
	if err != nil {
		log.Fatal(err)
	}
	return currentUser
}

func AskModerator(db *sql.DB, r *http.Request, asked int, id int) error {

	_, err := db.Exec("UPDATE users SET askedmod=? WHERE id=?", asked, id)
	if err != nil {
        return err
    }
	return nil
}

func FetchLikes(db *sql.DB, user int) models.Liked {

	//liked posts
	rows, err := db.Query("SELECT postid FROM linkpost WHERE userid = ? AND likes = ?", user, true)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var postlist []models.Post
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		var posts models.Post
		err = db.QueryRow("SELECT id, authorid, categoryid, title, content, creation, flaged FROM post WHERE id = ?", id).Scan(
				&posts.ID,
				&posts.AuthorID,
				&posts.Categoryid,
				&posts.Title,
				&posts.Content,
				&posts.CreationDate,
				&posts.Flaged)
		if err != nil {
			log.Fatal(err)
		}
		postlist = append(postlist, posts)
	}

	//liked comment
	rows, err = db.Query("SELECT commentid FROM linkcomment WHERE userid = ? AND likes = ?", user, true)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var commentlist []models.Comment
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		var comments models.Comment
		err = db.QueryRow("SELECT id, authorid, postid, content, creation, flaged FROM comment WHERE id = ?", id).Scan(
				&comments.ID,
				&comments.AuthorID,
				&comments.Postid,
				&comments.Content,
				&comments.CreationDate,
				&comments.Flaged)
		if err != nil {
			log.Fatal(err)
		}
		commentlist = append(commentlist, comments)
	}
	sort.Slice(postlist, func(i, j int) bool {
		return postlist[i].CreationDate > postlist[j].CreationDate
	})
	sort.Slice(commentlist, func(i, j int) bool {
		return commentlist[i].CreationDate > commentlist[j].CreationDate
	})
	if len(postlist) > 5 {
		postlist = postlist[:5]
	}
	if len(commentlist) > 5 {
		commentlist = commentlist[:5]
	}
	var liked models.Liked
	liked.Posted = postlist
	liked.Commented = commentlist
	
	return liked
}

func FetchDislikes(db *sql.DB, user int) models.Disliked {
	//liked posts
	rows, err := db.Query("SELECT postid FROM linkpost WHERE userid = ? AND likes = ?", user, false)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var postlist []models.Post
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		var posts models.Post
		err = db.QueryRow("SELECT id, authorid, categoryid, title, content, creation, flaged FROM post WHERE id = ?", id).Scan(
				&posts.ID,
				&posts.AuthorID,
				&posts.Categoryid,
				&posts.Title,
				&posts.Content,
				&posts.CreationDate,
				&posts.Flaged)
		if err != nil {
			log.Fatal(err)
		}
		postlist = append(postlist, posts)
	}

	//liked comment
	rows, err = db.Query("SELECT commentid FROM linkcomment WHERE userid = ? AND likes = ?", user, false)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var commentlist []models.Comment
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		var comments models.Comment
		err = db.QueryRow("SELECT id, authorid, postid, content, creation, flaged FROM comment WHERE id = ?", id).Scan(
				&comments.ID,
				&comments.AuthorID,
				&comments.Postid,
				&comments.Content,
				&comments.CreationDate,
				&comments.Flaged)
		if err != nil {
			log.Fatal(err)
		}
		commentlist = append(commentlist, comments)
	}
	sort.Slice(postlist, func(i, j int) bool {
		return postlist[i].CreationDate > postlist[j].CreationDate
	})
	sort.Slice(commentlist, func(i, j int) bool {
		return commentlist[i].CreationDate > commentlist[j].CreationDate
	})
	if len(postlist) > 5 {
		postlist = postlist[:5]
	}
	if len(commentlist) > 5 {
		commentlist = commentlist[:5]
	}
	var disliked models.Disliked
	disliked.Posted = postlist
	disliked.Commented = commentlist
	
	return disliked
}

func FetchComments(db *sql.DB, user int) {
	
	rows, err := db.Query("SELECT id FROM comment WHERE authorid = ?", user)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var commentlist []models.Compostfusion
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		fmt.Println(id)
		if err != nil {
			log.Fatal(err)
		}
		var compostfusion models.Compostfusion
		// err = db.QueryRow("SELECT id, authorid, postid, content, creation, flaged FROM comment WHERE id = ?", id).Scan(
		// 		&compostfusion.Commentorig.ID,
		// 		&compostfusion.Commentorig.AuthorID,
		// 		&compostfusion.Commentorig.Postid,
		// 		&compostfusion.Commentorig.Content,
		// 		&compostfusion.Commentorig.CreationDate,
		// 		&compostfusion.Commentorig.Flaged)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// commentlist = append(commentlist, compostfusion)
	}
}

// func FetchPosts(db *sql.DB, user int) models.User {
	
// }