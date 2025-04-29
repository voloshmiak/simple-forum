package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func UserAuthorization(next http.Handler) http.Handler {
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
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func AdminAuthorization(next http.Handler) http.Handler {
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
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		if role != "admin" {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
