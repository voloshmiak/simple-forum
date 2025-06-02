package handler

import (
	"forum-project/internal/application"
	"forum-project/internal/model"
	"net/http"
)

type HomeHandler struct {
	app *application.App
}

func NewHomeHandler(app *application.App) *HomeHandler {
	return &HomeHandler{
		app: app,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "home.page", new(model.Page))
	if err != nil {
		InternalServer(h.app.Logger, rw, "Unable to render template", err)
		return
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "about.page", new(model.Page))
	if err != nil {
		InternalServer(h.app.Logger, rw, "Unable to render template", err)
		return
	}
}
