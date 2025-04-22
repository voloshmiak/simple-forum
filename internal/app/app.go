package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/database"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	"forum-project/internal/template"
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
	server       *http.Server
	templates    *template.Manager
	topicService *service.TopicService
	postService  *service.PostService
}

func New() *App {
	app := &App{}

	// initialize logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	app.logger = logger

	// initialize environment variables
	if err := godotenv.Load(".env"); err != nil {
		app.logger.Error("Failed to load .env file")
		os.Exit(1)
	}

	// initialize template manager
	templatesManager, err := template.NewManager(true)
	if err != nil {
		app.logger.Error("Failed to create renderer", err)
		os.Exit(1)
	}
	app.templates = templatesManager

	// initialize db
	conn, err := database.Init()
	if err != nil {
		app.logger.Error("Failed to init database", err)
		os.Exit(1)
	}
	app.database = conn

	// initialize repositories
	postRepository := repository.NewPostRepository(app.database)
	topicRepository := repository.NewTopicRepository(app.database)

	// initialize services
	postService := service.NewPostService(postRepository)
	app.postService = postService

	topicService := service.NewTopicService(topicRepository)
	app.topicService = topicService

	app.registerRoutes()

	app.setupServer()

	return app
}

func (app *App) gracefulShutdown(done chan bool) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	app.logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := app.database.Close()
	if err != nil {
		app.logger.Error(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}

func (app *App) Run() {

	done := make(chan bool, 1)

	go app.gracefulShutdown(done)

	app.logger.Info(fmt.Sprintf("Starting server on port %s", app.server.Addr))

	if err := app.server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		app.logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	app.logger.Info("Server shutdown complete")
}
