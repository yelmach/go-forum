package controllers

import (
	"forum/models"
	"forum/utils"
)

// ;
// INSERT INTO categories (categori,post_id)VALUES()
func CreateCategories(categories models.Categories) error {
	resalt := ""
	err := utils.DataBase.QueryRow(`SELECT * FROM categories WHERE categori = ? nameCategori`, categories.Categori).Scan(&resalt)
	if err != nil {
		return err
	}
	return nil
}
