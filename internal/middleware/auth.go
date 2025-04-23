package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) { return []byte("secret-key"), nil })
		if err != nil {
			http.Error(rw, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(rw, "Invalid token", http.StatusUnauthorized)
		}

		next.ServeHTTP(rw, r)
	})
}
