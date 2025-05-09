package main

import (
	"context"
	"errors"
	"fmt"
	"forum-project/internal/app"
	"forum-project/internal/database"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"forum-project/internal/mylogger"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	app := &app.Config{}

	// init logger
	logger := mylogger.NewLogger()
	app.Logger = logger

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
	app.Templates = templateManager

	// init database
	conn, err := database.New()
	if err != nil {
		logger.Error("Failed to init database", "error", err)
		os.Exit(1)
	}

	app.Database = conn

	// init repositories
	postRepository := repository.NewPostRepository(conn)
	topicRepository := repository.NewTopicRepository(conn)
	userRepository := repository.NewUserRepository(conn)

	// init services
	postService := service.NewPostService(postRepository)
	topicService := service.NewTopicService(topicRepository)
	userService := service.NewUserService(userRepository)

	app.TopicService = topicService
	app.PostService = postService
	app.UserService = userService

	//init handlers
	hh := handlers.NewHomeHandler(app)
	th := handlers.NewTopicHandler(logger, templateManager, topicService, postService)
	ph := handlers.NewPostHandler(logger, templateManager, postService, topicService)
	uh := handlers.NewUserHandler(logger, templateManager, userService)

	// init mux
	mux := registerRoutes(hh, th, ph, uh)

	// init server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(mux, logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(app, done)

	logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	logger.Info("Server shutdown complete")
}

func gracefulShutdown(app *app.Config, done chan bool) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	app.Logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := app.Database.Close()
	if err != nil {
		app.Logger.Error(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}
