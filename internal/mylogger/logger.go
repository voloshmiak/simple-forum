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

func (el *WrappedLogger) ServerInternalError(rw http.ResponseWriter, msg string, err error) {
	el.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusInternalServerError)
}

func (el *WrappedLogger) NotFoundError(rw http.ResponseWriter, msg string, err error) {
	el.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusNotFound)
}

func (el *WrappedLogger) BadRequestError(rw http.ResponseWriter, msg string, err error) {
	el.Logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusBadRequest)
}
