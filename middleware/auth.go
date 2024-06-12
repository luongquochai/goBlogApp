package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/luongquochai/goBlog/util"
)

var config, _ = util.LoadConfig("../")
var store = sessions.NewCookieStore([]byte(config.SECRET_KEY))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "ZLUOHAI")
		if session.Values["authenticated"] != true {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
