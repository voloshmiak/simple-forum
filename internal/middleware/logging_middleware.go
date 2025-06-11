package middleware

import (
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

type logger struct {
	h http.Handler
	l *slog.Logger
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	wrappedRW := newResponseWriter(rw)

	l.h.ServeHTTP(wrappedRW, r)

	duration := time.Since(start).Seconds()

	l.l.Info(
		"HTTP request",
		"status", wrappedRW.statusCode,
		"method", r.Method,
		"path", r.URL.Path,
		"duration", duration,
	)
}

func NewLogging(l *slog.Logger) func(http.Handler) http.Handler {
	fn := func(h http.Handler) http.Handler {
		return &logger{
			h: h,
			l: l,
		}
	}
	return fn
}
