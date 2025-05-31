package responder

import (
	"log/slog"
	"net/http"
)

type Responder struct {
	logger *slog.Logger
}

func NewErrorResponder(logger *slog.Logger) *Responder {
	return &Responder{logger: logger}
}

func (e *Responder) BadRequest(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	http.Error(rw, msg, http.StatusBadRequest)
}

func (e *Responder) InternalServer(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	http.Error(rw, msg, http.StatusInternalServerError)
}

func (e *Responder) NotFound(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	http.Error(rw, msg, http.StatusNotFound)
}

func (e *Responder) Unauthorized(rw http.ResponseWriter, msg string, err error) {
	e.logger.Error(msg, "error", err)
	http.Error(rw, msg, http.StatusUnauthorized)
}
