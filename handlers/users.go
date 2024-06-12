package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	db "github.com/luongquochai/goBlog/db/models"
	"github.com/luongquochai/goBlog/util"
)

var config, _ = util.LoadConfig("../")
var store = sessions.NewCookieStore([]byte(config.SECRET_KEY))

type UserHandler struct {
	Queries *db.Queries
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUser, err := h.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := h.Queries.GetUserByUsername(r.Context(), user.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := util.CheckPasswordHash(user.Password, dbUser.PasswordHash); err != nil {
		http.Error(w, "Password is not corrected", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "ZLUOHAI")
	session.Values["authenticated"] = true
	session.Values["user-id"] = dbUser.ID
	session.Save(r, w)
	fmt.Println(session)

	json.NewEncoder(w).Encode(dbUser)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ZLUOHAI")
	session.Values["authenticated"] = false
	sessions.Save(r, w)
}
