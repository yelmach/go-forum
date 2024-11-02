package models

import (
	"database/sql"
)

type Post struct {
	Id       int            `json:"id"`
	By       string         `json:"by"`
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	ImageURL sql.NullString `json:"imageUrl"`
	// Using sql.NullString to handle potential NULL values
	CreatedAt  string    `json:"createdAt"`
	Comments   []Comment `json:"comments"`
	Categories []string  `json:"categories"`
	Likes      []int     `json:"likes"`
	Dislikes   []int     `json:"dislikes"`
}

type Comment struct {
	By        string `json:"by"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	Likes     []int  `json:"likes"`
	Dislikes  []int  `json:"dislikes"`
}
