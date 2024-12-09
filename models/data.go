package models

type PostApi struct {
	Id              int          `json:"id"`
	By              string       `json:"by"`
	Title           string       `json:"title"`
	Content         string       `json:"content"`
	CreatedAt       string       `json:"createdAt"`
	Comments        []CommentApi `json:"comments"`
	TotalComments   int          `json:"totalComments"`
	HasMoreComments bool         `json:"hasMoreComments"`
	Categories      []string     `json:"categories"`
	Likes           []int        `json:"likes"`
	Dislikes        []int        `json:"dislikes"`
}

type CommentApi struct {
	Id        int    `json:"id"`
	By        string `json:"by"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	Likes     []int  `json:"likes"`
	Dislikes  []int  `json:"dislikes"`
}

type Error struct {
	Error struct {
		Status int    `json:"status"`
		Code   string `json:"code"`
	} `json:"error"`
}
