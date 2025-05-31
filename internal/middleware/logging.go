package middleware

import (
	"forum-project/internal/application"
	"net/http"
)

func Logging(app *application.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			app.Logger.Info(
				"HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(rw, r)
		})
	}
}
