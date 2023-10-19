package forum

type Comments struct {
	ID       int64  `json:"id"`
	Post_ID  int64  `json:"post_id"`
	Author   string `json:"author"`
	Likes    int64  `json:"likes"`
	Dislikes int64  `json:"dislikes"`
}
