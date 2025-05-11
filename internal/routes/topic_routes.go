package routes

import (
	"forum-project/internal/config"
	"forum-project/internal/handlers"
	"forum-project/internal/middleware"
	"net/http"
)

func RegisterTopicRoutes(appConfig *config.AppConfig, th *handlers.TopicHandler) {
	adminMux := http.NewServeMux()
	adminware := middleware.CreateStack(
		middleware.AuthMiddleware,
		middleware.IsAdmin,
	)

	appConfig.Mux.HandleFunc("GET /topics", th.GetTopics)
	appConfig.Mux.HandleFunc("GET /topics/{topicID}", th.GetTopic)

	adminMux.HandleFunc("GET /topics/new", th.GetCreateTopic)
	adminMux.HandleFunc("POST /topics", th.PostCreateTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/edit", th.GetEditTopic)
	adminMux.HandleFunc("POST /topics/{topicID}/edit", th.PostEditTopic)
	adminMux.HandleFunc("GET /topics/{topicID}/delete", th.GetDeleteTopic)

	appConfig.Mux.Handle("/admin/", http.StripPrefix("/admin", adminware(adminMux)))
}
