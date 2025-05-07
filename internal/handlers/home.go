package handlers

import (
	"forum-project/internal/models"
	"forum-project/internal/mylogger"
	"forum-project/internal/template"
	"net/http"
)

type HomeHandler struct {
	logger    *mylogger.WrappedLogger
	templates *template.Manager
}

func NewHomeHandler(logger *mylogger.WrappedLogger, templates *template.Manager) *HomeHandler {
	return &HomeHandler{
		logger:    logger,
		templates: templates,
	}
}

func (h *HomeHandler) GetHome(rw http.ResponseWriter, r *http.Request) {
	err := h.templates.Render(rw, r, "home.page", &models.ViewData{})
	if err != nil {
		h.logger.ServerInternalError(rw, "Unable to render template", err)
	}
}
