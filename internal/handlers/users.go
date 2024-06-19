package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/luongquochai/goBlogApp/internal/models"
	"github.com/luongquochai/goBlogApp/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type UserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	if username == "" || password == "" {
		var user UserParams
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		username = user.Username
		password = user.Password
	}

	// TODO: Dynamic inpu login: username or email
	// TODO: Adjust sqlc: GetUserByEmail
	user, err := models.New(h.db).GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid creadentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}
	// Set session
	session := sessions.Default(c)
	session.Set("user", username)
	session.Set("user_id", user.ID)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/home")
	// c.JSON(http.StatusOK, user)
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
	email := c.PostForm("email")
	password := c.PostForm("password")

	if username == "" || email == "" || password == "" {
		var user models.CreateUserParams
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		username = user.Username
		email = user.Email
		password = user.HashedPassword
	}

	log.Printf("data: %v %v %v", username, email, password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	//TODO: check email is exist or not
	emailExists := util.EmailExists(h.db, email)
	if emailExists != nil {
		c.String(http.StatusExpectationFailed, "Email is already exists!!!")
		return
	}

	_, err = models.New(h.db).CreateUser(c.Request.Context(), models.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
	// c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("internal/templates/layout_auth.html", "internal/templates/changepassword.html"))
	data := PageData{Title: "Change Password"}
	if err := tmpl.ExecuteTemplate(c.Writer, "layout_auth.html", data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) ChangePasswordPost(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")
	if username == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	user, err := models.New(h.db).GetUserByUsername(c.Request.Context(), username.(string))
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid session")
		return
	}

	// Compare old password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(oldPassword)); err != nil {
		c.String(http.StatusUnauthorized, "Old password is incorrect")
		return
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to hash new password")
		return
	}

	// Update the password in the database
	err = models.New(h.db).UpdatePassword(c.Request.Context(), models.UpdatePasswordParams{
		Username:       user.Username,
		HashedPassword: string(hashedNewPassword),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
