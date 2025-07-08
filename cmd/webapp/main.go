package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"simple-forum/internal/auth"
	"simple-forum/internal/config"
	"simple-forum/internal/database"
	"simple-forum/internal/handler"
	"simple-forum/internal/middleware"
	"simple-forum/internal/repository"
	"simple-forum/internal/service"
	"simple-forum/internal/template"
	"syscall"
	"time"

	"github.com/justinas/nosurf"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Config
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// Migrate
	err = database.Migrate(cfg.DB.Addr, cfg.Path.ToMigrations)
	if err != nil {
		return err
	}

	// Connection
	conn, err := database.Connect(cfg.DB.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Logger
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Authenticator
	a := auth.NewJWTAuthenticator(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Templates
	t := template.NewTemplates(cfg.Env, cfg.Path.ToTemplates, a)

	// Repository
	postRepository := repository.NewPostRepository(conn)
	topicRepository := repository.NewTopicRepository(conn)
	userRepository := repository.NewUserRepository(conn)

	// Service
	postService := service.NewPostService(postRepository)
	topicService := service.NewTopicService(topicRepository)
	userService := service.NewUserService(userRepository)

	// Handlers
	hh := handler.NewHomeHandler(l, t)
	ph := handler.NewPostHandler(l, a, t, postService, topicService)
	th := handler.NewTopicHandler(l, a, t, postService, topicService)
	uh := handler.NewUserHandler(l, a, t, userService)

	// Mux
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// Middleware
	adminMiddleware := middleware.PermissionMiddleware(l, postService, "admin")
	authorMiddleware := middleware.PermissionMiddleware(l, postService, "author")
	sharedMiddleware := middleware.PermissionMiddleware(l, postService, "admin", "author")
	authMiddleware := middleware.AuthMiddleware(a)
	loggingMiddleware := middleware.LoggingMiddleware(l)

	// ToStatic
	fileserver := http.FileServer(http.Dir(filepath.ToSlash(cfg.Path.ToStatic)))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	// Home
	mux.HandleFunc("GET /home", hh.GetHome)
	mux.HandleFunc("GET /about", hh.GetAbout)

	// User
	mux.HandleFunc("GET /login", uh.GetLogin)
	mux.HandleFunc("POST /login", uh.PostLogin)
	mux.HandleFunc("GET /logout", uh.GetLogout)
	mux.HandleFunc("GET /signup", uh.GetRegister)
	mux.HandleFunc("POST /signup", uh.PostRegister)

	// Post
	mux.HandleFunc("GET /topics/{topicID}/posts/{postID}", ph.GetPost)
	authMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authMux.HandleFunc("GET /posts/{postID}/edit", authorMiddleware(http.HandlerFunc(ph.GetEditPost)))
	authMux.HandleFunc("POST /posts/{postID}/edit", authorMiddleware(http.HandlerFunc(ph.PostEditPost)))
	authMux.HandleFunc("GET /posts/{postID}/delete", sharedMiddleware(http.HandlerFunc(ph.GetDeletePost)))

	mux.Handle("/user/", http.StripPrefix("/user", authMiddleware(authMux))) // grouping

	// Topic
	mux.HandleFunc("GET /topics", th.GetTopics)
	mux.HandleFunc("GET /topics/{topicID}", th.GetTopic)
	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{topicID}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/delete", th.GetDeleteTopic)

	mux.Handle("/admin/", http.StripPrefix("/admin", authMiddleware(adminMiddleware(adminMux)))) // grouping

	// Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      nosurf.New(loggingMiddleware(mux)),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Listening
	l.Info("Starting server on port: " + server.Addr)

	done := make(chan bool)

	go func() {
		signs := make(chan os.Signal)
		signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)

		<-signs
		l.Info("Shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err = server.Shutdown(shutdownCtx); err != nil {
			l.Error("Server forced to shutdown", "error", err)
		}

		done <- true
	}()

	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-done
	l.Info("Graceful shutdown complete")

	return nil
}
