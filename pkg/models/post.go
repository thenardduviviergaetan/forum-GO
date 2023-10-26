package forum

type Post struct {
	ID             int64  `json:"id"`
	AuthorID       int64  `json:"author_id"`
	Author         string `json:"author"`
	Ifcurrentuser  bool
	Img            string `json:"img"`
	Ifimg          bool
	Categories     []int
	CategoriesName []string
	Title          string `json:"title"`
	Content        string `json:"content"`
	Like           int
	User_like      map[int64]bool
	Dislike        int
	User_dislike   map[int64]bool
	CreationDate   string `json:"time"`
	Flaged         int    `json:"flaged"`
	Tab_comment    []Comment
}
