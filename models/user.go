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
type PostContent struct {
	User_id    int
	Title      string
	Content    string
	Image_url  string
	Created_at string
}
