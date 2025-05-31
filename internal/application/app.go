package application

import (
	"database/sql"
	"forum-project/internal/config"
	"forum-project/internal/model"
	"forum-project/internal/repository"
	"forum-project/internal/responder"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"os"
)

type ErrorHandler interface {
	BadRequest(rw http.ResponseWriter, msg string, err error)
	InternalServer(rw http.ResponseWriter, msg string, err error)
	NotFound(rw http.ResponseWriter, msg string, err error)
	Unauthorized(rw http.ResponseWriter, msg string, err error)
}

type Renderer interface {
	AddDefaultData(td *model.Page, r *http.Request) *model.Page
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
	VerifyPostAuthor(post *model.Post, userID int) bool
	VerifyPostAuthorOrAdmin(post *model.Post, userID int, userRole string) bool
}

type UserServicer interface {
	Authenticate(email, password string, jwtSecret string, expiryHours int) (string, error)
	Register(username, email, password1, password2 string) error
	GetUserByID(id int) (*model.User, error)
}

type App struct {
	Config         *config.Config
	Logger         *slog.Logger
	ErrorResponder ErrorHandler
	Templates      Renderer
	TopicService   TopicServicer
	PostService    PostServicer
	UserService    UserServicer
}

func NewApp(conn *sql.DB, config *config.Config) *App {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// error responder
	errorResponder := responder.NewErrorResponder(logger)

	// templates renderer
	templates := template.NewTemplates(config)

	// repositories and services
	postRepository := repository.NewPostRepository(conn)
	postService := service.NewPostService(postRepository)

	topicRepository := repository.NewTopicRepository(conn)
	topicService := service.NewTopicService(topicRepository)

	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)

	return &App{
		Config:         config,
		Logger:         logger,
		ErrorResponder: errorResponder,
		Templates:      templates,
		TopicService:   topicService,
		PostService:    postService,
		UserService:    userService,
	}
}
