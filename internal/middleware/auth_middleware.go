package middleware

import (
	"context"
	"net/http"
	"simple-forum/internal/app"
)

func AuthMiddleware(app *app.App) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			claims, err := app.Authenticator.GetClaimsFromRequest(r)
			if err != nil {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user := claims["user"].(map[string]interface{})

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(rw, r.WithContext(ctx))
		}
	}
}
