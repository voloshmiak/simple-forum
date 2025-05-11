package routes

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
)

func RegisterHomeRoutes(appConfig *config.AppConfig, hh *handlers.HomeHandler) {
	appConfig.Mux.HandleFunc("GET /home", hh.GetHome)
}
