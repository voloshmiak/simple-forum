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
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		h.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}

func (h *HomeHandler) GetAbout(rw http.ResponseWriter, r *http.Request) {
	err := h.app.Templates.Render(rw, r, "about.page", new(model.Page))
	if err != nil {
		http.Error(rw, "Unable to render template", http.StatusInternalServerError)
		h.app.Logger.Error(err.Error(), "method", r.Method, "status",
			http.StatusInternalServerError, "path", r.URL.Path)
		return
	}
}
