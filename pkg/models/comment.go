package forum

type Comment struct {
	ID            int64  `json:"id"`
	AuthorID      int64  `json:"authorid"`
	Author        string `json:"author"`
	Ifcurrentuser bool
	Postid        int64  `json:"postid"`
	Content       string `json:"content"`
	Like          int
	User_like     map[int64]bool
	Dislike       int
	User_dislike  map[int64]bool
	Flaged		  int
	CreationDate  string `json:"time"`
}
