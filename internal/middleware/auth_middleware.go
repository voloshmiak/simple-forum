package middleware

import (
	"context"
	"errors"
	"forum-project/internal/app"
	"forum-project/internal/model"
	"forum-project/internal/service"
	"net/http"
)

func AuthMiddleware(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}
			claims, err := service.ValidateToken(cookie.Value, app.Config.JWT.Secret)
			if err != nil {
				if errors.Is(err, service.ErrInvalidToken) {
					http.Error(rw, "Invalid token", http.StatusUnauthorized)
					return
				}
				http.Error(rw, "Invalid token", http.StatusInternalServerError)
				app.Logger.Error("Failed to validate token", "error", err)
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
