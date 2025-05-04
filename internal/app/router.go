package app

import (
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func (app *App) initRouter() {
	// initialize mux
	mux := http.NewServeMux()
	authorizedMux := http.NewServeMux()
	//adminMux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(app.logger, app.templates, app.topicService)
	ph := handlers.NewPostHandler(app.logger, app.templates, app.postService)
	ah := handlers.NewAuthHandler(app.logger, app.templates, app.userService)

	// guests routing
	mux.HandleFunc("GET /topics", th.GetAllTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopicByID)
	mux.HandleFunc("GET /topics/{id}/posts", ph.GetPostsByTopicID)
	mux.HandleFunc("GET /posts/{id}", ph.GetPostByID)

	// authorization routing
	mux.HandleFunc("GET /login", ah.GetLogin)
	mux.HandleFunc("POST /login", ah.PostLogin)
	mux.HandleFunc("GET /logout", ah.GetLogout)
	mux.HandleFunc("POST /logout", ah.PostLogout)
	mux.HandleFunc("GET /register", ah.GetRegister)
	mux.HandleFunc("POST /register", ah.PostRegister)

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

	app.mux = mux
}
