package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ww := &wrappedWriter{
			ResponseWriter: rw,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		logger.Info(
			"HTTP request",
			"method", r.Method,
			"status", strconv.Itoa(ww.statusCode),
			"path", r.URL.Path,
		)
	})
}
