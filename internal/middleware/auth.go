package middleware

import (
	"forum-project/internal/auth"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func UserAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := auth.ValidateTokenFromRequest(r)
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func AdminAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := auth.ValidateTokenFromRequest(r)
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
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
