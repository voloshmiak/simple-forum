package handlers

import (
	"fmt"
	"forum-project/internal/models"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
)

type HomeHandler struct {
	logger    *slog.Logger
	templates *template.Manager
}

func NewHomeHandler(logger *slog.Logger, templates *template.Manager) *HomeHandler {
	return &HomeHandler{
		logger:    logger,
		templates: templates,
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.templates.Render(rw, r, "about.page", &models.ViewData{})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Unable to template template: %s", err))
		http.Error(rw, fmt.Sprintf("Unable to template template: %s", err), http.StatusInternalServerError)
	}
}
