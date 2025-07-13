package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Authenticator interface {
	GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error)
}

func AuthMiddleware(a Authenticator) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			claims, err := a.GetClaimsFromRequest(r)
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
