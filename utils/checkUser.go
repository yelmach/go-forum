package utils

import (
	"errors"
	"regexp"

	"forum/database"
	"forum/models"
)

// CheckUserExist checks user if already registered
func CheckUserExist(user models.User) error {
	var isExist bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
	if err := database.DataBase.QueryRow(query, user.Email, user.Username).Scan(&isExist); err != nil {
		return err
	}
	if isExist {
		return errors.New("credentials already exists")
	}
	return nil
}

// CheckUsernameFormat checks username format if it's valid
func CheckUsernameFormat(username string) (bool, error) {
	valid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{3,21}$`, username)
	if err != nil || !valid {
		return false, err
	}

	return true, nil
}

// CheckEmailFormat checks email format if it is valid
func CheckEmailFormat(email string) (bool, error) {
	if len(email) > 60 {
		return false, nil
	}

	isValid, err := regexp.MatchString(`^[a-zA-Z-0-9]+\.?[a-zA-Z-0-9]+@[a-zA-Z-0-9]+\.[a-z]+$`, email)
	if err != nil {
		return false, err
	} else if !isValid {
		return false, nil
	}
	return true, nil
}

// CheckPasswordFormat checks if password written correct
func CheckPasswordFormat(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	isSpecial := regexp.MustCompile(`[^\w\s]`)
	isLower := regexp.MustCompile(`[a-z]`)
	isUpper := regexp.MustCompile(`[A-Z]`)
	isDigit := regexp.MustCompile(`[0-9]`)
	return isLower.MatchString(password) && isUpper.MatchString(password) && isDigit.MatchString(password) && isSpecial.MatchString(password)
}
