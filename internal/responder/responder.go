package responder

import (
	"log/slog"
	"net/http"
)

type HonestResponder struct {
	l *slog.Logger
}

func NewHonestResponder(l *slog.Logger) *HonestResponder {
	return &HonestResponder{l}
}

func (r *HonestResponder) InternalServerError(rw http.ResponseWriter, msg string, err error) {
	http.Error(rw, msg, http.StatusInternalServerError)
	r.l.Error(msg, "error", err.Error())
}

func (r *HonestResponder) BadRequest(rw http.ResponseWriter, msg string, err error) {
	http.Error(rw, msg, http.StatusBadRequest)
	r.l.Error(msg, "error", err.Error())
}

func (r *HonestResponder) Unauthorized(rw http.ResponseWriter, msg string, err error) {
	http.Error(rw, msg, http.StatusUnauthorized)
	r.l.Error(msg, "error", err.Error())
}

func (r *HonestResponder) Forbidden(rw http.ResponseWriter, msg string, err error) {
	http.Error(rw, msg, http.StatusForbidden)
	r.l.Error(msg, "error", err.Error())
}

func (r *HonestResponder) NotFound(rw http.ResponseWriter, msg string, err error) {
	http.Error(rw, msg, http.StatusNotFound)
	r.l.Error(msg, "error", err.Error())
}
