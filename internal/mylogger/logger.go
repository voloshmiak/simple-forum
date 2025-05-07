package mylogger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type WrappedLogger struct {
	*slog.Logger
}

func NewLogger() *WrappedLogger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return &WrappedLogger{Logger: logger}
}

func (wl *WrappedLogger) ServerInternalError(rw http.ResponseWriter, msg string, err error) {
	wl.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusInternalServerError)
}

func (wl *WrappedLogger) NotFoundError(rw http.ResponseWriter, msg string, err error) {
	wl.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusNotFound)
}

func (wl *WrappedLogger) BadRequestError(rw http.ResponseWriter, msg string, err error) {
	wl.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusBadRequest)
}

func (wl *WrappedLogger) UnauthorizedError(rw http.ResponseWriter, msg string, err error) {
	wl.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusUnauthorized)
}
