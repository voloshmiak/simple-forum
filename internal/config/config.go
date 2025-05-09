package config

import (
	"database/sql"
	"forum-project/internal/helpers"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"log/slog"
	"net/http"
)

type AppConfig struct {
	Logger       *slog.Logger
	Database     *sql.DB
	Server       *http.Server
	Templates    *template.Manager
	TopicService *service.TopicService
	PostService  *service.PostService
	UserService  *service.UserService
	Errors       *helpers.ErrorHandler
}
