package handlers

import (
	"forum-project/internal/config"
	"forum-project/internal/models"
	"net/http"
)

type HomeHandler struct {
	app *config.AppConfig
}

func NewHomeHandler(app *config.AppConfig) *HomeHandler {
	return &HomeHandler{
		app: app,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "home.page", new(models.Page))
	if err != nil {
		h.app.ErrorResponder.InternalServer(rw, "Unable to render template", err)
	}
}
