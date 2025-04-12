package main

import (
	"context"
	"errors"
	"fmt"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"forum-project/internal/render"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(server *http.Server, logger *slog.Logger, done chan bool) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	done <- true
}

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
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(mux, logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(server, logger, done)

	// start serving server
	logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	logger.Info("Server shutdown complete")
}
