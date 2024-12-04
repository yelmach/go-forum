package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/routers"
)

func main() {
	if err := database.InitDb(); err != nil {
		log.Fatalln(err)
	}

	defer database.DataBase.Close()

	rootMux := http.NewServeMux()
	routers.SetupRoutes(rootMux)

	fmt.Println("Server running on port: 8080")
	fmt.Println("URL: http://localhost:8080")

	if err := http.ListenAndServe(":8080", rootMux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
