package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/goBlog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *sql.DB
}

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	handler := &AuthHandler{db: db}

	router.Static("/static", "./static")

	router.GET("/", handler.HomePage)
	router.GET("/login", handler.Login)
	router.POST("/login", handler.LoginPost)
	router.GET("/register", handler.Register)
	router.POST("/register", handler.RegisterPost)
	router.GET("/change-password", handler.ChangePassword)
	router.POST("/change-password", handler.ChangePasswordPost)

}

type PageData struct {
	Title string
}

// Function helpers
func (h *AuthHandler) HomePage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout.html", "internal/templates/home.html"))
	data := PageData{Title: "Home"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout.html", "internal/templates/login.html"))
	data := PageData{Title: "Login"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := models.New(h.db).GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid creadentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) Register(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout.html", "internal/templates/register.html"))
	data := PageData{Title: "Login"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) RegisterPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = models.New(h.db).CreateUser(c.Request.Context(), models.CreateUserParams{
		Username:       username,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout.html", "internal/templates/changepassword.html"))
	data := PageData{Title: "Login"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) ChangePasswordPost(c *gin.Context) {
	username := c.PostForm("username")
	newPassword := c.PostForm("newpassword")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	err = models.New(h.db).UpdatePassword(c.Request.Context(), models.UpdatePasswordParams{
		Username:       username,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
