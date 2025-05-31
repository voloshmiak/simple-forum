package middleware

import (
	"context"
	"forum-project/internal/application"
	"forum-project/internal/auth"
	"forum-project/internal/model"
	"net/http"
)

func AuthMiddleware(app *application.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			claims, err := auth.GetClaimsFromRequest(r, app.Config.JWT.Secret)
			if err != nil {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user := claims["user"].(map[string]interface{})

			userIDFloat := user["id"].(float64)
			userIDInt := int(userIDFloat)
			username := user["username"].(string)
			role := user["role"].(string)

			authorizedUser := &model.AuthorizedUser{
				ID:       userIDInt,
				Username: username,
				Role:     role,
			}

			ctx := context.WithValue(r.Context(), "user", authorizedUser)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
