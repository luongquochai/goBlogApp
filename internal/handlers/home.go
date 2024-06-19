package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Function helpers
func (h *AuthHandler) HomePage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout.html", "internal/templates/home.html"))
	data := PageData{Title: "Home"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) HomePageAuth(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	log.Printf("User ID: %v\n", userID)
	// Ensure userID is not nil and type assertion is correct
	if userID == nil {
		http.Error(c.Writer, "User not logged in", http.StatusUnauthorized)
		return
	}

	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(c.Writer, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:           "Home",
		CurrentAuthorID: userIDInt,
	}

	tmpl := template.Must(template.ParseFiles("internal/templates/layout_auth.html", "internal/templates/home_auth.html"))
	if err := tmpl.ExecuteTemplate(c.Writer, "layout_auth.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
