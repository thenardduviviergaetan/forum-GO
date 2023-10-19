package forum

import (
	"time"
)

type Post struct {
	ID           int64     `json:"id"`
	AuthorID     int64     `json:"author_id"`
	Author       string    `json:"author"`
	Category     string    `json:"category"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Like         int       `json:"like"`
	Dislike      int       `json:"dislike"`
	CreationDate time.Time `json:"time"`
}

type PostTemplates struct {
	Post     Post
	Err      string
	IsSigned bool
}
