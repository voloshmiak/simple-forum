package main

import (
	"errors"
	"fmt"
	"forum-project/internal/handlers"
	"forum-project/internal/render"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	// initialize logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// initialize environment variables
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Failed to load .env file")
		os.Exit(1)
	}

	// initialize renderer
	renderer, err := render.NewRenderer()
	if err != nil {
		logger.Error("Failed to create renderer", err)
		os.Exit(1)
	}

	// initialize mux
	mux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(logger, renderer)
	ph := handlers.NewPostHandler(logger, renderer)

	// routing
	mux.HandleFunc("GET /topics/", th.GetTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopic)
	mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPosts)
	mux.HandleFunc("GET /posts/{id}", ph.GetPost)

	// server configuration
	server := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start serving server
	logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		logger.Error("Error listening and serving", err)
	}
}
