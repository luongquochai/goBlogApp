package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	models "github.com/luongquochai/goBlog/db/models"
	"github.com/luongquochai/goBlog/handlers"
	"github.com/luongquochai/goBlog/middleware"
	"github.com/luongquochai/goBlog/util"
)

func main() {
	// Load Config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	db, err := sql.Open(config.DB_Driver, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	queries := models.New(db)
	userHandler := &handlers.UserHandler{Queries: queries}
	postHandler := &handlers.PostHandler{Queries: queries}

	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	// Post routes
	router.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(postHandler.CreatePost))).Methods("POST")
	router.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postHandler.GetPost))).Methods("GET")
	router.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postHandler.UpdatePost))).Methods("PUT")
	router.Handle("/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(postHandler.DeletePost))).Methods("DELETE")
	router.Handle("/posts", http.HandlerFunc(postHandler.ListPosts)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
