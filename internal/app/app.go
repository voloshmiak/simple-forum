package app

import (
	"context"
	"errors"
	"fmt"
	"forum-project/internal/config"
	"forum-project/internal/database"
	"forum-project/internal/handlers"
	"forum-project/internal/helpers"
	"forum-project/internal/middleware"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type App struct {
	config *config.AppConfig
}

func New() *App {
	// create a new app config instance
	config := &config.AppConfig{}

	// init logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	config.Logger = logger

	// load environment variables
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}

	// init template manager
	templateManager, err := template.NewManager()
	if err != nil {
		logger.Error("Failed to create renderer", "error", err)
		os.Exit(1)
	}

	config.Templates = templateManager

	// init database
	conn, err := database.New()
	if err != nil {
		logger.Error("Failed to init database", "error", err)
		os.Exit(1)
	}

	config.Database = conn

	// init repositories
	postRepository := repository.NewPostRepository(conn)
	topicRepository := repository.NewTopicRepository(conn)
	userRepository := repository.NewUserRepository(conn)

	// init services
	postService := service.NewPostService(postRepository)
	topicService := service.NewTopicService(topicRepository)
	userService := service.NewUserService(userRepository)

	config.PostService = postService
	config.TopicService = topicService
	config.UserService = userService

	// init error handler
	errorHandler := helpers.NewErrorHandler(logger)

	config.Errors = errorHandler

	//init handlers
	hh := handlers.NewHomeHandler(config)
	th := handlers.NewTopicHandler(config)
	ph := handlers.NewPostHandler(config)
	uh := handlers.NewUserHandler(config)

	// init server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(registerRoutes(hh, th, ph, uh), logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	config.Server = server

	return &App{config: config}
}

func (app *App) gracefulShutdown(done chan bool) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	app.config.Logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.config.Server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.config.Logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := app.config.Database.Close()
	if err != nil {
		app.config.Logger.Error(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}

func (app *App) Run() {
	done := make(chan bool, 1)

	go app.gracefulShutdown(done)

	app.config.Logger.Info(fmt.Sprintf("Starting server on port %s", app.config.Server.Addr))

	if err := app.config.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.config.Logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	app.config.Logger.Info("Server shutdown complete")
}
