package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	db "github.com/luongquochai/goBlog/db"
	models "github.com/luongquochai/goBlog/db/models"
	"golang.org/x/crypto/bcrypt"
)

// Helper functions to render templates
func renderTemplate(c *gin.Context, tmpl string, data interface{}) {
	// if err := templates.ExecuteTemplate(c.Writer, tmpl, data); err != nil {
	// 	c.String(http.StatusInternalServerError, err.Error())
	// }

	err := templates.ExecuteTemplate(c.Writer, tmpl, data)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error rendering template %s: %v", tmpl, err)
		// Respond with an error page or message
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}
}

func ShowHomePage(c *gin.Context) {
	renderTemplate(c, "layout.html", nil)
}

func ShowHomePageAfterLogin(c *gin.Context) {
	renderTemplate(c, "index.html", nil)
}

func ShowRegisterPage(c *gin.Context) {
	renderTemplate(c, "register.html", nil)
}

func ShowLoginPage(c *gin.Context) {
	renderTemplate(c, "login.html", nil)
}

func ShowChangePasswordPage(c *gin.Context) {
	renderTemplate(c, "change_password.html", nil)
}

func Register(c *gin.Context) {

	// Print raw request body for debugging
	// rawBody, _ := c.GetRawData()
	// fmt.Println("Raw Request Body:", string(rawBody))

	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	user, err := db.Queries.CreateUser(context.Background(), models.CreateUserParams{
		Username: creds.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	// Process form data (for example, save to database)
	// For demonstration purposes, assume registration is successful
	successMessage := "Registration successful!"

	// Render success message in HTML
	c.HTML(http.StatusOK, "register.html", gin.H{
		"SuccessMessage": successMessage,
	})

	c.JSON(http.StatusOK, user)
}
func Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := db.Queries.GetUserByUsername(context.Background(), creds.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	session, _ := store.Get(c.Request, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = creds.Username
	session.Save(c.Request, c.Writer)

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// func ChangePassword(c *gin.Context) {
// 	session, _ := store.Get(c.Request, "session")
// 	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	var creds Credentials
// 	if err := c.ShouldBindJSON(&creds); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
// 		return
// 	}

// 	_, err = db.Queries.UpdateUserPassword(context.Background(), db.UpdateUserPasswordParams{
// 		Password: string(hashedPassword),
// 		Username: session.Values["username"].(string),
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
// }

func AuthMiddleware(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	c.Next()
}
