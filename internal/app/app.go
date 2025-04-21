package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"forum-project/internal/render"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
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

	// initialize db
	conn, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		os.Getenv("HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("USER"), os.Getenv("PASSWORD")))
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer conn.Close()

	// initialize repositories
	postRepository := repository.NewPostRepository(conn)

	// initialize services
	postService := service.NewPostService(postRepository)

	// initialize mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(logger, renderer)
	ph := handlers.NewPostHandler(logger, renderer, postService)

	// guests routing
	mux.HandleFunc("GET /topics/", th.GetTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopic)
	// mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPosts)
	mux.HandleFunc("GET /posts/{id}", ph.GetPost)

	// authorized users routing
	authorizedMux.HandleFunc("POST /posts", ph.CreatePost)
	authorizedMux.HandleFunc("PUT /posts/{id}", ph.UpdatePost)
	authorizedMux.HandleFunc("DELETE /posts/{id}", ph.DeletePost)

	mux.Handle("/", authorizedMux)

	// admin routing
	adminMux.HandleFunc("POST /topics", th.CreateTopic)
	adminMux.HandleFunc("PUT /topics/{id}", th.UpdateTopic)
	adminMux.HandleFunc("DELETE /topics/{id}", th.DeleteTopic)

	mux.Handle("/admin", adminMux)

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
