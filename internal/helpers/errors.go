package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorHandler struct {
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (e *ErrorHandler) BadRequest(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusBadRequest)
}

func (e *ErrorHandler) InternalServer(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusInternalServerError)
}

func (e *ErrorHandler) NotFound(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusNotFound)
}

func (e *ErrorHandler) Unauthorized(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusUnauthorized)
}
