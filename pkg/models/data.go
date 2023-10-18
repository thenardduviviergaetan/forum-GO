package forum

type Data struct {
	Connected      bool
	Posts          []Post
	CurrentPost    Post
	CurrentComment Comment
}
