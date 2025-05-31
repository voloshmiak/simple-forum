package route

import (
	"forum-project/internal/application"
	"forum-project/internal/handler"
	"forum-project/internal/middleware"
	"github.com/justinas/nosurf"
	"net/http"
)

func RegisterRoutes(app *application.App) http.Handler {
	// Handlers
	hh := handler.NewHomeHandler(app)
	ph := handler.NewPostHandler(app)
	th := handler.NewTopicHandler(app)
	uh := handler.NewUserHandler(app)

	// Mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// Middleware
	auth := middleware.AuthMiddleware(app)
	isAdmin := middleware.IsAdmin
	isPostAuthor := middleware.IsPostAuthor(app)
	isPostAuthorOrAdmin := middleware.IsPostAuthorOrAdmin(app)
	logging := middleware.Logging(app)

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
	authorizedMux.HandleFunc("GET /topics/{topicID}/posts/new", ph.GetCreatePost)
	authorizedMux.HandleFunc("POST /posts", ph.PostCreatePost)
	authorizedMux.Handle("GET /posts/{postID}/edit", isPostAuthor(http.HandlerFunc(ph.GetEditPost)))
	authorizedMux.Handle("POST /posts/{postID}/edit", isPostAuthor(http.HandlerFunc(ph.PostEditPost)))
	authorizedMux.Handle("GET /posts/{postID}/delete", isPostAuthorOrAdmin(http.HandlerFunc(ph.GetDeletePost)))

	mux.Handle("/user/", http.StripPrefix("/user", auth(authorizedMux))) // grouping

	// Topic
	mux.HandleFunc("GET /topics", th.GetTopics)
	mux.HandleFunc("GET /topics/{topicID}", th.GetTopic)
	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{topicID}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/delete", th.GetDeleteTopic)

	mux.Handle("/admin/", http.StripPrefix("/admin", auth(isAdmin(adminMux)))) // grouping

	return nosurf.New(logging(mux))
}
