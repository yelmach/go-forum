package controllers

import (
	"errors"

	"forum/models"
	"forum/utils"
)

// ;
// INSERT INTO categories (categori,post_id)VALUES()
func CreateCategories(categories models.Categories) error {
	var count int
	err := utils.DataBase.QueryRow("SELECT COUNT(*) FROM categories WHERE categori = ? nameCategori", categories.Categori).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("categori already exist")
	}
	stmt, err := utils.DataBase.Prepare("INSERT INTO categories (categori,post_id) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(categories.Categori, categories.PostId)
	if err != nil {
		return err
	}
	return nil
}
