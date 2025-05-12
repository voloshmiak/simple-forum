package config

import (
	"database/sql"
	"forum-project/internal/database"
	"forum-project/internal/httperror"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"github.com/joho/godotenv"
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

func LoadEnv() error {
	// load environment variables
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

func NewAppConfig() (*AppConfig, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	responder := httperror.NewErrorResponder(logger)

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}

	templates, err := template.NewTemplates()
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Logger:    logger,
		Database:  conn,
		Templates: templates,
		Errors:    responder,
	}, nil
}
