package router

import (
	"github.com/justinas/nosurf"
	"net/http"
	"simple-forum/internal/app"
	"simple-forum/internal/handler"
	"simple-forum/internal/middleware"
)

func RegisterRoutes(app *app.App) http.Handler {
	// Handlers
	hh := handler.NewHomeHandler(app)
	ph := handler.NewPostHandler(app)
	th := handler.NewTopicHandler(app)
	uh := handler.NewUserHandler(app)

	// Mux
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// Middleware
	adminMiddleware := middleware.PermissionMiddleware(app, "admin")
	authorMiddleware := middleware.PermissionMiddleware(app, "author")
	sharedMiddleware := middleware.PermissionMiddleware(app, "admin", "author")
	authMiddleware := middleware.AuthMiddleware(app)
	loggingMiddleware := middleware.LoggingMiddleware(app)

	// Static
	fileserver := http.FileServer(http.Dir(app.Config.Path.ToStatic()))
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

	return nosurf.New(loggingMiddleware(mux))
}
