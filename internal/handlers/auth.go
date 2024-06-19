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
	postHandler := &PostHandler{db: db}

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
	// router.GET("/api/posts", postHandler.ListPosts) // List posts on homepage
	router.GET("/list-posts", postHandler.ListPosts)

	// Authorization routes
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/home", handler.HomePageAuth)
		authorized.GET("/change-password", handler.ChangePassword)
		authorized.POST("/change-password", handler.ChangePasswordPost)

		// Post routes
		authorized.POST("/create-post", postHandler.CreatePost)
		authorized.GET("/cancel", postHandler.CancelHandler)
		// authorized.GET("/list-posts", postHandler.ListPosts)
		authorized.PUT("/edit-post/:id", postHandler.UpdatePostByAuthor)
		authorized.DELETE("/delete-post/:id", postHandler.DeletePost)
	}

}
