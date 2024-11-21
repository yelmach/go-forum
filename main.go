package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/routers"
)

func main() {
	err := database.InitDb()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer database.DataBase.Close()

	rootMux := http.NewServeMux()

	routers.SetupRoutes(rootMux)

	fmt.Println("Server running on port: 8080")
	fmt.Println("URL: http://localhost:8080")
	http.ListenAndServe(":8080", rootMux)
}
