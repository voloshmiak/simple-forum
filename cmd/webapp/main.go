package main

import (
	"context"
	"errors"
	"fmt"
	"forum-project/internal/config"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"forum-project/internal/repository"
	"forum-project/internal/routes"
	"forum-project/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load environment variables
	if err := config.LoadEnv(); err != nil {
		panic(err)
	}

	// create app config instance
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	// init mux
	mux := http.NewServeMux()

	// serve static
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/static"))))

	// Home
	hh := handlers.NewHomeHandler(appConfig)
	routes.RegisterHomeRoutes(mux, hh)

	// Post
	postRepository := repository.NewPostRepository(appConfig.Database)
	postService := service.NewPostService(postRepository)
	appConfig.PostService = postService
	ph := handlers.NewPostHandler(appConfig)
	routes.RegisterPostRoutes(mux, ph, appConfig)

	// Topic
	topicRepository := repository.NewTopicRepository(appConfig.Database)
	topicService := service.NewTopicService(topicRepository)
	appConfig.TopicService = topicService
	th := handlers.NewTopicHandler(appConfig)
	routes.RegisterTopicRoutes(mux, th)

	// User
	userRepository := repository.NewUserRepository(appConfig.Database)
	userService := service.NewUserService(userRepository)
	appConfig.UserService = userService
	uh := handlers.NewUserHandler(appConfig)
	routes.RegisterUserRoutes(mux, uh)

	// Server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      middleware.Logging(mux, appConfig.Logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	appConfig.Server = server

	done := make(chan bool, 1)

	go gracefulShutdown(done, appConfig)

	appConfig.Logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		appConfig.Logger.Error(fmt.Sprintf("Error listening and serving: %s", err))
		os.Exit(1)
	}

	// wait for graceful shutdown to complete
	<-done
	appConfig.Logger.Info("Server shutdown complete")
}

func gracefulShutdown(done chan bool, appConfig *config.AppConfig) {
	// wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	appConfig.Logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := appConfig.Server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		appConfig.Logger.Error(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := appConfig.Database.Close()
	if err != nil {
		appConfig.Logger.Error(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}
