package handler

import (
	"forum-project/internal/app"
	"forum-project/internal/model"
	"net/http"
)

type HomeHandler struct {
	app *app.App
}

func NewHomeHandler(app *app.App) *HomeHandler {
	return &HomeHandler{
		app: app,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "home.page", new(model.Page))
	if err != nil {
		h.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "about.page", new(model.Page))
	if err != nil {
		h.handleError(rw, "Unable to render template", err, http.StatusInternalServerError)
		return
	}
}

func (h *HomeHandler) handleError(rw http.ResponseWriter, msg string, err error, code int) {
	http.Error(rw, msg, code)
	h.app.Logger.Error(msg, "error", err.Error())
}
