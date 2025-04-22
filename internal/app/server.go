package app

import (
	"forum-project/internal/middleware"
	"net/http"
	"os"
	"time"
)

func (app *App) setupServer() {
	// server configuration
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(app.mux, app.logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	app.server = server
}
