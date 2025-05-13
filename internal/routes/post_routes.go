package routes

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func RegisterPostRoutes(mux *http.ServeMux, ph *handlers.PostHandler, appConfig *config.AppConfig) {
	// init mux for authorized routes
	authorizedMux := http.NewServeMux()

	// initialize middleware
	auth := middleware.AuthMiddleware
	isPostAuthor := middleware.IsPostAuthor(appConfig)
	isPostAuthorOrAdmin := middleware.IsPostAuthorOrAdmin(appConfig)

	// public routing
	mux.HandleFunc("GET /topics/{topicID}/posts/{postID}", ph.GetPost)

	// authorized routing
	authorizedMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.Handle("GET /posts/{postID}/edit", isPostAuthor(http.HandlerFunc(ph.GetEditPost)))
	authorizedMux.Handle("POST /posts/{postID}/edit", isPostAuthor(http.HandlerFunc(ph.PostEditPost)))
	authorizedMux.Handle("GET /posts/{postID}/delete", isPostAuthorOrAdmin(http.HandlerFunc(ph.GetDeletePost)))

	mux.Handle("/user/", http.StripPrefix("/user", auth(authorizedMux)))
}
