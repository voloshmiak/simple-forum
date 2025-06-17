package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
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

func LoggingMiddleware(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedRW := newResponseWriter(rw)

			next.ServeHTTP(wrappedRW, r)

			duration := fmt.Sprintf("%fs", time.Since(start).Seconds())

			l.Info(
				"HTTP request",
				"status", wrappedRW.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"duration", duration,
			)
		})
	}
}
