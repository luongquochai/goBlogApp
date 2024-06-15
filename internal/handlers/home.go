package handlers

import (
	"html/template"
	"net/http"

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
	tmpl := template.Must(template.ParseFiles("internal/templates/layout_auth.html", "internal/templates/home_auth.html"))
	data := PageData{Title: "Home"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout_auth.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
