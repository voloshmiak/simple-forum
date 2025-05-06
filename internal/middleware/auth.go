package middleware

import (
	"context"
	"forum-project/internal/auth"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := auth.ValidateTokenFromRequest(r)
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		user := claims["user"].(map[string]interface{})
		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		role := user.(map[string]interface{})["role"].(string)

		if role != "admin" {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
