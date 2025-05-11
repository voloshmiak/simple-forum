package routes

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
)

func RegisterUserRoutes(appConfig *config.AppConfig, uh *handlers.UserHandler) {
	appConfig.Mux.HandleFunc("GET /login", uh.GetLogin)
	appConfig.Mux.HandleFunc("POST /login", uh.PostLogin)
	appConfig.Mux.HandleFunc("GET /logout", uh.GetLogout)
	appConfig.Mux.HandleFunc("GET /signup", uh.GetRegister)
	appConfig.Mux.HandleFunc("POST /signup", uh.PostRegister)
}
