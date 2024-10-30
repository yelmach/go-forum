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
type Categories struct {
	Id       int
	Categori string
	PostId   int
}
type Posts struct {
	Id      int
	Title   string
	Content string
}
