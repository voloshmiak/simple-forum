package application

import (
	"database/sql"
	"forum-project/internal/httperror"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"os"
)

type App struct {
	Logger         *slog.Logger
	ErrorResponder httperror.ErrorHandler
	Templates      template.Renderer
	TopicService   service.TopicServicer
	PostService    service.PostServicer
	UserService    service.UserServicer
}

func NewApp(conn *sql.DB) *App {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// error responder
	responder := httperror.NewErrorResponder(logger)

	// templates renderer
	templates := template.NewTemplates()

	// repositories and services
	postRepository := repository.NewPostRepository(conn)
	postService := service.NewPostService(postRepository)

	topicRepository := repository.NewTopicRepository(conn)
	topicService := service.NewTopicService(topicRepository)

	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)

	return &App{
		Logger:         logger,
		ErrorResponder: responder,
		Templates:      templates,
		TopicService:   topicService,
		PostService:    postService,
		UserService:    userService,
	}
}
