package config

import (
	"database/sql"
	"forum-project/internal/database"
	"forum-project/internal/httperror"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type AppConfig struct {
	Logger       *slog.Logger
	Database     *sql.DB
	Server       *http.Server
	Templates    *template.Templates
	TopicService *service.TopicService
	PostService  *service.PostService
	UserService  *service.UserService
	Errors       *httperror.ErrorResponder
}

func LoadEnv() {
	// load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load .env file", "error", err)
	}
}

func NewAppConfig() *AppConfig {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	conn, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	templates, err := template.NewTemplates()
	if err != nil {
		log.Fatal("Failed to create templates", "error", err)
	}

	responder := httperror.NewErrorResponder(logger)

	return &AppConfig{
		Logger:    logger,
		Database:  conn,
		Templates: templates,
		Errors:    responder,
	}
}
