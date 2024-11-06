package models

import "time"

type User struct {
	Id        int
	Username  string
	Password  string
	Email     string
	SessionId string
}

type Session struct {
	Username  string
	ExpiresAt time.Time
}

type Post struct {
	User_id     int
	Title       string
	Content     string
	Category_id []string
	Image_url   string
}

type Comment struct {
	Post_id    int
	User_id    int
	Content    string
}

type Reaction struct {
	User_id    int
	Post_id    int
	Comment_id int
	Is_like    bool
}