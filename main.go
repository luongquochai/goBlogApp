package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/goBlog/db"
	"github.com/luongquochai/goBlog/handlers"
	"github.com/luongquochai/goBlog/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	dbConn := db.InitDB(config.DB_SOURCE)
	db.InitSchema(dbConn)
	r := gin.Default()
	// Serve static files
	r.Static("/static", "./static")
	// Load HTML templates from the "templates" directory
	r.LoadHTMLGlob("templates/*.html")

	// Public routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(handlers.AuthMiddleware)
	{
		auth.GET("/dashboard", handlers.ShowHomePageAfterLogin)
		// TODO:
		// auth.GET("/change-password", ShowChangePasswordPage)
		// auth.POST("/change-password", ChangePassword)
		auth.GET("/tasks", handlers.GetTasks)
		auth.POST("/tasks", handlers.CreateTask)
		auth.PUT("/tasks/:id", handlers.UpdateTask)
		auth.DELETE("/tasks/:id", handlers.DeleteTask)
	}

	r.Run(":8000")
}
