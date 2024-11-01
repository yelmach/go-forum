package models

import (
	"database/sql"
)

type Post struct {
	Id         int
	By         string
	Title      string
	Content    string
	ImageURL   sql.NullString // Using sql.NullString to handle potential NULL values
	CreatedAt  string
	Comments   []Comment
	Categories []string
	Likes      []int
	Dislikes   []int
}

type Comment struct {
	By        string
	Content   string
	CreatedAt string
	Likes     []int
	Dislikes  []int
}
