package utils

import (
	"regexp"
	"strings"
)


func CheckEmailFormat(email string) bool {
	// heytesting@gmail.com
	// heytesting@hotmail.com
	// heytesting@univ.edu

	return strings.Contains(email, "@") 
		
}

func CheckPasswordFormat(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	isUper, _ := regexp.MatchString(`[A-Z]`, password) 
	if !isUper {
		return false
	}	

	isLower, _ := regexp.MatchString(`[a-z]`, password) 
	if !isLower {
		return false
	}
	return true
}
