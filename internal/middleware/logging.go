package middleware

import (
	"forum-project/internal/app"
	"net/http"
)

func Logging(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			app.Logger.Info(
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(rw, r)
		})
	}
}
