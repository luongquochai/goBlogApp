package handlers

import (
	"database/sql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db *sql.DB
}

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	handler := &AuthHandler{db: db}

	// Set up sessions
	router.Use(sessions.Sessions("mysession", store))

	router.Static("/static", "./internal/static")

	// Public routes
	router.GET("/", handler.HomePage)
	router.GET("/login", handler.Login)
	router.POST("/login", handler.LoginPost)
	router.GET("/register", handler.Register)
	router.POST("/register", handler.RegisterPost)
	router.GET("/logout", handler.Logout) // Logout route

	// Authorization routes
	authourized := router.Group("/")
	authourized.Use(AuthMiddleware())
	{
		authourized.GET("/home", handler.HomePageAuth)
		authourized.GET("/change-password", handler.ChangePassword)
		authourized.POST("/change-password", handler.ChangePasswordPost)
	}

}
