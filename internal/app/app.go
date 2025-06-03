package app

import (
	"database/sql"
	"forum-project/internal/config"
	"forum-project/internal/model"
	"forum-project/internal/repository"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"os"
)

type Renderer interface {
	Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *model.Page) error
}

type TopicServicer interface {
	GetAllTopics() ([]*model.Topic, error)
	GetTopicByID(id int) (*model.Topic, error)
	GetTopicByPostID(id int) (*model.Topic, error)
	CreateTopic(name, description string, authorID int) error
	EditTopic(id int, name, description string) error
	DeleteTopic(id int) error
}

type PostServicer interface {
	GetPostByID(userID int) (*model.Post, error)
	GetPostsByTopicID(topicID int) ([]*model.Post, error)
	CreatePost(title, content string, topicID, authorID int, authorName string) error
	EditPost(title, content string, postID int) error
	DeletePost(postID int) error
}

type Authenticator interface {
	Login(email, password string) (string, error)
	Register(username, email, password1, password2 string) error
}

type App struct {
	Config        *config.Config
	Logger        *slog.Logger
	Templates     Renderer
	TopicService  TopicServicer
	PostService   PostServicer
	Authenticator Authenticator
}

func New(conn *sql.DB, config *config.Config) *App {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// templates
	templates := template.NewTemplates(config.JWT.Secret, config.Env, config.Path.ToTemplates())

	// repositories and services
	postRepository := repository.NewPostRepository(conn)
	postService := service.NewPostService(postRepository)

	topicRepository := repository.NewTopicRepository(conn)
	topicService := service.NewTopicService(topicRepository)

	userRepository := repository.NewUserRepository(conn)
	authService := service.NewAuthService(userRepository, config.JWT.Secret, config.JWT.Expiration)

	return &App{
		Config:        config,
		Logger:        logger,
		Templates:     templates,
		TopicService:  topicService,
		PostService:   postService,
		Authenticator: authService,
	}
}
