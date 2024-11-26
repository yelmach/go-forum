package main

import (
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

	routers.ListenAndServe(rootMux)
	// http.ListenAndServe(":8080", rootMux)
}
