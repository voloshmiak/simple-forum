package app

import (
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
)

func initRouter(logger *slog.Logger, templates *template.Manager, topicService *service.TopicService, postService *service.PostService, userService *service.UserService) *http.ServeMux {
	// initialize mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(logger, templates, topicService)
	ph := handlers.NewPostHandler(logger, templates, postService, topicService)
	uh := handlers.NewUserHandler(logger, templates, userService)

	// guests routing
	mux.HandleFunc("GET /topics", th.GetTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopic)
	mux.HandleFunc("GET /topics/{id}/posts", ph.GetPosts)
	mux.HandleFunc("GET /posts/{id}", ph.GetPost)

	// authorization routing
	mux.HandleFunc("GET /login", uh.GetLogin)
	mux.HandleFunc("POST /login", uh.PostLogin)
	mux.HandleFunc("GET /logout", uh.GetLogout)
	mux.HandleFunc("POST /logout", uh.PostLogout)
	mux.HandleFunc("GET /register", uh.GetRegister)
	mux.HandleFunc("POST /register", uh.PostRegister)

	// authorized users routing
	authorizedMux.HandleFunc("GET /topics/{id}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.HandleFunc("GET /posts/{id}/edit", ph.GetEditPost)
	authorizedMux.HandleFunc("POST /posts/{id}/edit", ph.PostEditPost)
	authorizedMux.HandleFunc("GET /posts/{id}/delete", ph.GetDeletePost)
	authorizedMux.HandleFunc("POST /posts/{id}/delete", ph.PostDeletePost)

	mux.Handle("/", middleware.UserAuthorization(authorizedMux))

	// admin routing
	mux.Handle("GET /topics/new", middleware.AdminAuthorization(http.HandlerFunc(th.GetCreateTopic)))
	mux.Handle("POST /topics", middleware.AdminAuthorization(http.HandlerFunc(th.PostCreateTopic)))
	mux.Handle("GET /topics/{id}/edit", middleware.AdminAuthorization(http.HandlerFunc(th.GetEditTopic)))
	mux.Handle("POST /topics/{id}/edit", middleware.AdminAuthorization(http.HandlerFunc(th.PutEditTopic)))
	mux.Handle("GET /topics/{id}/delete", middleware.AdminAuthorization(http.HandlerFunc(th.GetDeleteTopic)))
	mux.Handle("POST /topics/{id}/delete", middleware.AdminAuthorization(http.HandlerFunc(th.PostDeleteTopic)))

	return mux
}
