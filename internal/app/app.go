package app

import (
	"database/sql"
	"forum-project/internal/mylogger"
	"forum-project/internal/service"
	"forum-project/internal/template"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Logger       *mylogger.WrappedLogger
	Database     *sql.DB
	Server       *http.Server
	Templates    *template.Manager
	TopicService *service.TopicService
	PostService  *service.PostService
	UserService  *service.UserService
}
