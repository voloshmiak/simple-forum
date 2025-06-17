package middleware

import (
	"fmt"
	"net/http"
	"simple-forum/internal/app"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(app *app.App) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedRW := newResponseWriter(rw)

			next.ServeHTTP(wrappedRW, r)

			duration := fmt.Sprintf("%fs", time.Since(start).Seconds())

			app.Logger.Info(
				"HTTP request",
				"status", wrappedRW.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
			)
		}
	}
}
