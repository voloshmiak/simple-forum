package httperror

import (
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponder struct {
	logger *slog.Logger
}

func NewErrorResponder(logger *slog.Logger) *ErrorResponder {
	return &ErrorResponder{logger: logger}
}

func (e *ErrorResponder) BadRequest(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusBadRequest)
}

func (e *ErrorResponder) InternalServer(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusInternalServerError)
}

func (e *ErrorResponder) NotFound(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusNotFound)
}

func (e *ErrorResponder) Unauthorized(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	msg = fmt.Sprintf("%s: %s", msg, err)
	http.Error(rw, msg, http.StatusUnauthorized)
}
