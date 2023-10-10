package forum

type Post struct {
	ID       int64  `json:"id"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Like     int    `json:"like"`
	Dislike  int    `json:"dislike"`
}
