package routes

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func RegisterPostRoutes(mux *http.ServeMux, ph *handlers.PostHandler, appConfig *config.AppConfig) {
	authorizedMux := http.NewServeMux()

	canEditPost := middleware.IsPostAuthor(appConfig)
	canDeletePost := middleware.IsPostAuthorOrAdmin(appConfig)

	mux.HandleFunc("GET /topics/{topicID}/posts/{postID}", ph.GetPost)

	// authorized users routing
	authorizedMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.Handle("GET /posts/{postID}/edit", canEditPost(http.HandlerFunc(ph.GetEditPost)))
	authorizedMux.Handle("POST /posts/{postID}/edit", canEditPost(http.HandlerFunc(ph.PostEditPost)))
	authorizedMux.Handle("GET /posts/{postID}/delete", canDeletePost(http.HandlerFunc(ph.GetDeletePost)))

	mux.Handle("/user/", http.StripPrefix("/user", middleware.AuthMiddleware(authorizedMux)))
}
