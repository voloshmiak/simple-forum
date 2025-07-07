package handler

import (
	"net/http"
	"simple-forum/internal/app"
	"simple-forum/internal/model"
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
		h.app.Responder.InternalServerError(rw, "Unable to render template", err)
		return
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "about.page", new(model.Page))
	if err != nil {
		h.app.Responder.InternalServerError(rw, "Unable to render template", err)
		return
	}
}
