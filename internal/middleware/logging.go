package middleware

import (
	"log/slog"
	"net/http"
)

func Logging(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		logger.Info(
			"HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(rw, r)
	})
}
