package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/database"
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
	logger       *slog.Logger
	database     *sql.DB
	mux          *http.ServeMux
	server       *http.Server
	templates    *template.Manager
	topicService *service.TopicService
	postService  *service.PostService
	userService  *service.UserService
}

func New() *App {
	app := &App{}

	app.initLogger()
	app.loadEnviroment()
	app.initTemplates()
	app.initDatabase()
	app.initServices()
	app.initRouter()
	app.initServer()

	return app
}

func (app *App) initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	app.logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func (app *App) loadEnviroment() {
	if err := godotenv.Load(".env"); err != nil {
		app.logger.Error("Failed to load .env file")
		os.Exit(1)
	}
}

func (app *App) initTemplates() {
	templateManager, err := template.NewManager()
	if err != nil {
		app.logger.Error("Failed to create renderer", "error", err)
		os.Exit(1)
	}
	app.templates = templateManager
}

func (app *App) initDatabase() {
	conn, err := database.Init()
	if err != nil {
		app.logger.Error("Failed to init database", "error", err)
		os.Exit(1)
	}
	app.database = conn
}

func (app *App) initServices() {
	// initialize repositories
	postRepository := repository.NewPostRepository(app.database)
	topicRepository := repository.NewTopicRepository(app.database)
	userRepository := repository.NewUserRepository(app.database)

	app.postService = service.NewPostService(postRepository)
	app.topicService = service.NewTopicService(topicRepository)
	app.userService = service.NewUserService(userRepository)
}

func (app *App) initServer() {
	app.server = &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(app.mux, app.logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
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
