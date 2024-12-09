package models

type User struct {
	Id        int
	Username  string
	Password  string
	Email     string
	SessionId string
}

type Post struct {
	UserId     int
	Title      string
	Content    string
	Categories []string
}

type Comment struct {
	PostId  int
	UserId  int
	Content string
}

type Reaction struct {
	UserId    int
	PostId    int
	CommentId int
	IsLike    bool
	IsDislike bool
}
