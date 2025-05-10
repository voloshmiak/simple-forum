package middleware

import (
	"context"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		claims, err := auth.GetClaimsFromRequest(r)
		if err != nil {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := claims["user"].(map[string]interface{})

		userIDFloat := user["id"].(float64)
		userIDInt := int(userIDFloat)
		username := user["username"].(string)
		email := user["email"].(string)
		role := user["role"].(string)

		authorizedUser := &models.AuthorizedUser{
			ID:       userIDInt,
			Username: username,
			Email:    email,
			Role:     role,
		}

		ctx := context.WithValue(r.Context(), "user", authorizedUser)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		role := user.(*models.AuthorizedUser).Role

		if role != "admin" {
			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
