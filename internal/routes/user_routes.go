package routes

import (
	"forum-project/internal/handlers"
	"net/http"
)

func RegisterUserRoutes(mux *http.ServeMux, uh *handlers.UserHandler) {
	mux.HandleFunc("GET /login", uh.GetLogin)
	mux.HandleFunc("POST /login", uh.PostLogin)
	mux.HandleFunc("GET /logout", uh.GetLogout)
	mux.HandleFunc("GET /signup", uh.GetRegister)
	mux.HandleFunc("POST /signup", uh.PostRegister)
}
