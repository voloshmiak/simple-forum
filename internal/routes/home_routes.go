package routes

import (
	"forum-project/internal/handlers"
	"net/http"
)

func RegisterHomeRoutes(mux *http.ServeMux, hh *handlers.HomeHandler) {
	// public routing
	mux.HandleFunc("GET /home", hh.GetHome)
	mux.HandleFunc("GET /about", hh.GetAbout)
}
