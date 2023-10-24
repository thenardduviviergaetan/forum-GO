package forum

type Liked struct {
	Posted		[]Post
	Commented	[]Comment
}

type Disliked struct {
	Posted		[]Post
	Commented	[]Comment
}

type Compostfusion struct {
	Postlink	Post
	Commentorig	Comment
}