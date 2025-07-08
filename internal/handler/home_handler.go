package handler

import (
	"log/slog"
	"net/http"
	"simple-forum/internal/model"
	"simple-forum/internal/template"
)

type HomeHandler struct {
	l *slog.Logger
	t *template.Templates
}

func NewHomeHandler(l *slog.Logger, t *template.Templates) *HomeHandler {
	return &HomeHandler{
		l: l,
		t: t,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.t.Render(rw, r, "home.page", new(model.Page))
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		h.l.Error(msg, "error", err.Error())
		return
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.t.Render(rw, r, "about.page", new(model.Page))
	if err != nil {
		msg := "Unable to render template"
		http.Error(rw, msg, http.StatusInternalServerError)
		h.l.Error(msg, "error", err.Error())
		return
	}
}
