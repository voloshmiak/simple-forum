package config

import (
	"database/sql"
	"forum-project/internal/database"
	"forum-project/internal/httperror"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
	"os"
)

type AppConfig struct {
	Logger         *slog.Logger
	Database       *sql.DB
	Server         *http.Server
	Templates      *template.Templates
	TopicService   *service.TopicService
	PostService    *service.PostService
	UserService    *service.UserService
	ErrorResponder *httperror.ErrorResponder
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
		Logger:         logger,
		Database:       conn,
		Templates:      templates,
		ErrorResponder: responder,
	}, nil
}
