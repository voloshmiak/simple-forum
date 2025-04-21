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

type App struct {
	logger       *slog.Logger
	database     *sql.DB
	mux          *http.ServeMux
	topicService *service.TopicService
	postService  *service.PostService
}

func New() *App {
	app := &App{}
	app.setupLogger()

	// initialize environment variables
	if err := godotenv.Load(".env"); err != nil {
		app.logger.Error("Failed to load .env file")
		os.Exit(1)
	}

	// initialize renderer
	renderer, err := render.NewRenderer()
	if err != nil {
		app.logger.Error("Failed to create renderer", err)
		os.Exit(1)
	}

	app.setupDatabase()
	app.setupServices()
	app.setupMux(renderer)

	return app
}

func (app *App) gracefulShutdown(server *http.Server, logger *slog.Logger, done chan bool) {
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

	app.database.Close()

	done <- true
}

func (app *App) Run() {
	// server configuration
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(app.mux, app.logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool, 1)

	go app.gracefulShutdown(server, app.logger, done)

	// start serving server
	app.logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		app.logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	app.logger.Info("Server shutdown complete")
}

func (app *App) setupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	app.logger = logger
}

func (app *App) setupDatabase() {
	conn, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		os.Getenv("HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("USER"), os.Getenv("PASSWORD")))
	if err != nil {
		app.logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}

	app.database = conn
}

func (app *App) setupServices() {
	// initialize repositories
	postRepository := repository.NewPostRepository(app.database)
	topicRepository := repository.NewTopicRepository(app.database)

	postService := service.NewPostService(postRepository)
	app.postService = postService

	topicService := service.NewTopicService(topicRepository)
	app.topicService = topicService
}

func (app *App) setupMux(renderer *render.Renderer) {
	// initialize mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(app.logger, renderer, app.topicService)
	ph := handlers.NewPostHandler(app.logger, renderer, app.postService)

	// guests routing
	mux.HandleFunc("GET /topics/", th.GetAllTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopicByID)
	mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPostsByTopicID)
	mux.HandleFunc("GET /posts/{id}", ph.GetPostByID)

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

	app.mux = mux
}
