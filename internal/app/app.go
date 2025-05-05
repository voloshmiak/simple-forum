package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/database"
	"forum-project/internal/handlers"
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
	// init logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// load environment variables
	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Failed to load .env file")
		os.Exit(1)
	}

	// init template manager
	templateManager, err := template.NewManager()
	if err != nil {
		logger.Error("Failed to create renderer", "error", err)
		os.Exit(1)
	}

	// init database
	conn, err := database.New()
	if err != nil {
		logger.Error("Failed to init database", "error", err)
		os.Exit(1)
	}

	// init repositories
	postRepository := repository.NewPostRepository(conn)
	topicRepository := repository.NewTopicRepository(conn)
	userRepository := repository.NewUserRepository(conn)

	// init services
	postService := service.NewPostService(postRepository)
	topicService := service.NewTopicService(topicRepository)
	userService := service.NewUserService(userRepository)

	//init handlers
	th := handlers.NewTopicHandler(logger, templateManager, topicService)
	ph := handlers.NewPostHandler(logger, templateManager, postService, topicService)
	uh := handlers.NewUserHandler(logger, templateManager, userService)

	// init mux
	mux := registerRoutes(th, ph, uh)

	// init server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(mux, logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &App{
		logger:       logger,
		database:     conn,
		mux:          mux,
		server:       server,
		templates:    templateManager,
		topicService: topicService,
		postService:  postService,
		userService:  userService,
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
