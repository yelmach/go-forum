package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/routers"
	"forum/utils"
)

func main() {
	err := utils.InitDb()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer utils.DataBase.Close()
	mux := http.NewServeMux()
	routers.SetupRoutes(mux)

	fmt.Println("Server running on port: 8080")
	fmt.Println("URL: http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
