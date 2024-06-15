package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/luongquochai/goBlogApp/internal/handlers"
	"github.com/luongquochai/goBlogApp/internal/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	router := gin.Default()

	handlers.RegisterRoutes(router, db)

	log.Println("Server started at: ", config.PORT)
	log.Fatal(router.Run(config.PORT))
}
