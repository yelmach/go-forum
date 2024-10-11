package main

import (
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

	http.ListenAndServe(":8080", app)
}
