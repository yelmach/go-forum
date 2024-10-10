package main

import (
	"fmt"
	"net/http"

	"forum/internal/server"
)

const (
	RED   = "\x1b[0;31m"
	GREEN = "\x1b[0;32m"
	RESET = "\x1b[0m"
)

func main() {
	// Initialize the database
	server.InitDB()
	// Serve login page
	http.HandleFunc("/", server.LoginPage)
	// Handle login form submission
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/sign_up", server.Sign_UpHandler)
	http.HandleFunc("/index", server.Index)
	// Start the server
	fmt.Println(GREEN + "http://localhost:8080")
	fmt.Println(RED + "Ctrl+C to stop it" + RESET)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error %s\n", err)
	}
}
