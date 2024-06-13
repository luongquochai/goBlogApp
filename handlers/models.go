package handlers

import (
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/luongquochai/goBlog/util"
)

var config, _ = util.LoadConfig("../")

var (
	key       = []byte(config.SUPER_SECRET_KEY)
	store     = sessions.NewCookieStore(key)
	jwtKey    = []byte(config.MY_SECRET_KEY)
	templates = template.Must(template.ParseGlob("templates/*.html"))
	// dbConn    *sql.DB
)

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}
