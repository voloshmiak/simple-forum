package main

import (
	"fmt"
	"forum-project/internal/handlers"
	"forum-project/internal/render"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	// init logger
	logger := zap.NewExample().Sugar()

	// init renderer
	renderer := render.NewRenderer(logger)

	// set debug mode
	renderer.Debug = true

	// init mux
	mux := http.NewServeMux()

	// init handlers
	th := handlers.NewTopicHandler(logger, renderer)
	ph := handlers.NewPostHandler(logger, renderer)

	// routing
	mux.HandleFunc("GET /topics/", th.GetTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopic)
	mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPosts)
	mux.HandleFunc("GET /posts/{id}", ph.GetPost)

	// server configuration
	server := http.Server{
		Addr:        ":8090",
		Handler:     mux,
		ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	// running
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal("Error listening and serving ", err)
	}
}
