package app

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func registerRoutes(hh *handlers.HomeHandler, th *handlers.TopicHandler, ph *handlers.PostHandler, uh *handlers.UserHandler, config *config.AppConfig) *http.ServeMux {
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
	mux.HandleFunc("GET /signup", uh.GetRegister)
	mux.HandleFunc("POST /signup", uh.PostRegister)

	authRequired := middleware.AuthMiddleware
	isAdmin := middleware.IsAdmin
	canEditPost := middleware.IsPostAuthor(config)
	canDeletePost := middleware.IsPostAuthorOrAdmin(config)

	// authorized users routing
	authorizedMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.Handle("GET /posts/{postID}/edit", canEditPost(http.HandlerFunc(ph.GetEditPost)))
	authorizedMux.Handle("POST /posts/{postID}/edit", canEditPost(http.HandlerFunc(ph.PostEditPost)))
	authorizedMux.Handle("GET /posts/{postID}/delete", canDeletePost(http.HandlerFunc(ph.GetDeletePost)))

	mux.Handle("/user/", http.StripPrefix("/user", authRequired(authorizedMux)))

	// admin routing
	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{topicID}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/delete", th.GetDeleteTopic)

	adminware := middleware.CreateStack(
		authRequired,
		isAdmin,
	)

	mux.Handle("/admin/", http.StripPrefix("/admin", adminware(adminMux)))

	return mux
}
