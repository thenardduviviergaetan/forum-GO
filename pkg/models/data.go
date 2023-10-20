package forum

type Data struct {
	Connected      bool
	Admin          bool
	Moderator      bool
	Posts          []Post
	CurrentPost    Post
	CurrentComment Comment
	Categories     []Categories
}
