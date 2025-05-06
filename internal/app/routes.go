package app

import (
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func registerRoutes(hh *handlers.HomeHandler, th *handlers.TopicHandler, ph *handlers.PostHandler, uh *handlers.UserHandler) *http.ServeMux {
	// initialize mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// guests routing
	mux.HandleFunc("GET /about", hh.GetAbout)
	mux.HandleFunc("GET /topics", th.GetTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopic)
	mux.HandleFunc("GET /topics/{topicID}/posts/{postID}", ph.GetPost)

	// authorization routing
	mux.HandleFunc("GET /login", uh.GetLogin)
	mux.HandleFunc("POST /login", uh.PostLogin)
	mux.HandleFunc("GET /logout", uh.GetLogout)
	mux.HandleFunc("GET /register", uh.GetRegister)
	mux.HandleFunc("POST /register", uh.PostRegister)

	// authorized users routing
	authorizedMux.HandleFunc("GET /topics/{id}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.HandleFunc("GET /posts/{id}/edit", ph.GetEditPost)
	authorizedMux.HandleFunc("POST /posts/{id}/edit", ph.PostEditPost)
	authorizedMux.HandleFunc("GET /posts/{id}/delete", ph.GetDeletePost)

	mux.Handle("/user/", http.StripPrefix("/user", middleware.AuthMiddleware(authorizedMux)))

	// admin routing
	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{id}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{id}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{id}/delete", th.GetDeleteTopic)

	mux.Handle("/admin/", http.StripPrefix("/admin", middleware.AuthMiddleware(middleware.IsAdmin(adminMux))))

	return mux
}
