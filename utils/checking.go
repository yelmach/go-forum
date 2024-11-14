package utils

import (
	"regexp"
)

func CheckEmailFormat(email string) (bool, error) {
	isValid, err := regexp.MatchString(`\w+@[a-z]+\.[a-z]+`, email)
	if err != nil {
		return false, err
	} else if !isValid {
		return false, nil
	}
	return true, nil
}

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
