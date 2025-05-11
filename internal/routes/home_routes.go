package routes

import (
	"forum-project/internal/handlers"
	"net/http"
)

func RegisterHomeRoutes(mux *http.ServeMux, hh *handlers.HomeHandler) {
	mux.HandleFunc("GET /home", hh.GetHome)
}
