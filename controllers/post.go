package controllers

import (
	"errors"
	"fmt"
	"forum/database"
	"forum/models"
	"strconv"
)

func CreatePost(postContent models.Post) error {
	C_post, err := database.DataBase.Prepare(`INSERT INTO posts(user_id, title, content, image_url) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer C_post.Close()

	res, err := C_post.Exec(postContent.User_id, postContent.Title, postContent.Content, postContent.Image_url)
	if err != nil {
		return err
	}

	// Get the ID of the newly created post
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	C_postCategories, err := database.DataBase.Prepare(`INSERT INTO post_categories(post_id, category_id) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement for post_categories: %w", err)
	}
	defer C_postCategories.Close()

	for _, categoryID := range postContent.Category_id {
		Category_id, _ := strconv.Atoi(categoryID)
		if _, err := C_postCategories.Exec(postID, Category_id); err != nil {
			return fmt.Errorf("error linking post to category %s: %w", categoryID, err)
		}
	}
	return nil
}

func CreateCategorie(name_categorie string) error {
	var count int

	err := database.DataBase.QueryRow("SELECT COUNT(*) FROM categories WHERE name = ?", name_categorie).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking category existence: %w", err)
	}
	if count > 0 {
		return errors.New("category already exists")
	}

	C_categories, err := database.DataBase.Prepare(`INSERT INTO categories (name) VALUES (?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer C_categories.Close()

	if _, err := C_categories.Exec(name_categorie); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CreateComment(comment models.Comment) error {
	C_comment, err := database.DataBase.Prepare(`INSERT INTO comments(post_id, user_id, content) VALUES(?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}

	defer C_comment.Close()

	if _, err = C_comment.Exec(comment.Post_id, comment.User_id, comment.Content); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}
	return nil
}

func CreateReaction(reactions models.Reaction) error {
	var query string
	var id int
	isPost := reactions.Post_id != 0
	if isPost {
		query = `INSERT INTO reactions (user_id, post_id, is_like) VALUES (?, ?, ?)`
		id = reactions.Post_id
	} else {
		query = `INSERT INTO reactions (user_id, comment_id, is_like) VALUES (?, ?, ?)`
		id = reactions.Comment_id
	}
	C_reaction, err := database.DataBase.Prepare(query)
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}
	if _, err = C_reaction.Exec(reactions.User_id, id, reactions.Is_like); err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}
