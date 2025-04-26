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
	adminMux := http.NewServeMux()

	// initialize handlers
	th := handlers.NewTopicHandler(app.logger, app.templates, app.topicService)
	ph := handlers.NewPostHandler(app.logger, app.templates, app.postService)
	ah := handlers.NewAuthHandler(app.logger, app.templates, app.userService)

	// guests routing
	mux.HandleFunc("GET /topics/", th.GetAllTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopicByID)
	mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPostsByTopicID)
	mux.HandleFunc("GET /posts/{id}", ph.GetPostByID)

	// authorization routing
	mux.HandleFunc("GET /login", ah.GetLogin)
	mux.HandleFunc("POST /login", ah.PostLogin)
	mux.HandleFunc("GET /logout", ah.GetLogout)
	mux.HandleFunc("POST /logout", ah.PostLogout)
	mux.HandleFunc("GET /register", ah.GetRegister)
	mux.HandleFunc("POST /register", ah.PostRegister)

	// authorized users routing
	authorizedMux.HandleFunc("POST /posts", ph.CreatePost)
	authorizedMux.HandleFunc("PUT /posts/{id}", ph.UpdatePost)
	authorizedMux.HandleFunc("DELETE /posts/{id}", ph.DeletePost)

	mux.Handle("/", middleware.Auth(authorizedMux))

	// admin routing
	adminMux.HandleFunc("POST /topics", th.CreateTopic)
	adminMux.HandleFunc("PUT /topics/{id}", th.UpdateTopic)
	adminMux.HandleFunc("DELETE /topics/{id}", th.DeleteTopic)

	mux.Handle("/admin", middleware.EnsureAdmin(middleware.Auth(adminMux), app.userService))

	app.mux = mux
}
