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
	// check if data provided exists
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("email, username and password are required")
	}

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
	if err := database.DataBase.QueryRow(query, user.Email, user.Username).Scan(&isExist); err != nil {
		return err
	}
	if isExist {
		return errors.New("credentials already exists")
	}
	return nil
}

// CheckEmailFormat checks email format if it is valid
func CheckEmailFormat(email string) (bool, error) {
	isValid, err := regexp.MatchString(`\w+@[a-z]+\.[a-z]+`, email)
	if err != nil {
		return false, err
	} else if !isValid {
		return false, nil
	}
	return true, nil
}

// CheckPasswordFormat checks if password written correct
func CheckPasswordFormat(password string) bool {
	if len(password) < 8 {
		return false
	}
	isSpecial := regexp.MustCompile(`[^\w\s]`)
	isLower := regexp.MustCompile(`[a-z]`)
	isUpper := regexp.MustCompile(`[A-Z]`)
	isDigit := regexp.MustCompile(`[0-9]`)
	return isLower.MatchString(password) && isUpper.MatchString(password) && isDigit.MatchString(password) && isSpecial.MatchString(password)
}
