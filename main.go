package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/routers"
	"forum/utils"
)

func main() {
	err := utils.InitDb() // initial data base
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer utils.DataBase.Close()

	app := utils.NewApp() // multipluxer
	routers.InitPagesRouter(app)
	routers.InitApiRouter(app)

	fmt.Println("Server running on port: 8080")
	fmt.Println("URL: http://localhost:8080")
	http.ListenAndServe(":8080", app)
}
