package app

import (
	"forum-project/internal/handlers"
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

	// guests routing
	mux.HandleFunc("GET /topics/", th.GetAllTopics)
	mux.HandleFunc("GET /topics/{id}", th.GetTopicByID)
	mux.HandleFunc("GET /topics/{id}/posts/", ph.GetPostsByTopicID)
	mux.HandleFunc("GET /posts/{id}", ph.GetPostByID)

	// authorized users routing
	authorizedMux.HandleFunc("POST /posts", ph.CreatePost)
	authorizedMux.HandleFunc("PUT /posts/{id}", ph.UpdatePost)
	authorizedMux.HandleFunc("DELETE /posts/{id}", ph.DeletePost)

	mux.Handle("/", authorizedMux)

	// admin routing
	adminMux.HandleFunc("POST /topics", th.CreateTopic)
	adminMux.HandleFunc("PUT /topics/{id}", th.UpdateTopic)
	adminMux.HandleFunc("DELETE /topics/{id}", th.DeleteTopic)

	mux.Handle("/admin", adminMux)

	app.mux = mux
}
