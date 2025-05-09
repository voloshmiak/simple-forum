package handlers

import (
	"forum-project/internal/app"
	"forum-project/internal/models"
	"net/http"
)

type HomeHandler struct {
	app *app.Config
}

func NewHomeHandler(app *app.Config) *HomeHandler {
	return &HomeHandler{
		app: app,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "home.page", &models.ViewData{})
	if err != nil {
		h.app.Logger.ServerInternalError(rw, "Unable to render template", err)
	}
}
