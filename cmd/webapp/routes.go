package main

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

	// serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// guests routing
	mux.HandleFunc("GET /home", hh.GetHome)
	mux.HandleFunc("GET /topics", th.GetTopics)
	mux.HandleFunc("GET /topics/{topicID}", th.GetTopic)
	mux.HandleFunc("GET /topics/{topicID}/posts/{postID}", ph.GetPost)

	// authorization routing
	mux.HandleFunc("GET /login", uh.GetLogin)
	mux.HandleFunc("POST /login", uh.PostLogin)
	mux.HandleFunc("GET /logout", uh.GetLogout)
	mux.HandleFunc("GET /register", uh.GetRegister)
	mux.HandleFunc("POST /register", uh.PostRegister)

	// authorized users routing
	authorizedMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.HandleFunc("GET /posts/{postID}/edit", ph.GetEditPost)
	authorizedMux.HandleFunc("POST /posts/{postID}/edit", ph.PostEditPost)
	authorizedMux.HandleFunc("GET /posts/{postID}/delete", ph.GetDeletePost)

	authware := middleware.AuthMiddleware

	mux.Handle("/user/", http.StripPrefix("/user", authware(authorizedMux)))

	// admin routing
	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{topicID}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/delete", th.GetDeleteTopic)

	adminware := middleware.CreateStack(
		middleware.AuthMiddleware,
		middleware.IsAdmin,
	)

	mux.Handle("/admin/", http.StripPrefix("/admin", adminware(adminMux)))

	return mux
}
