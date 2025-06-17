package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
)

type Path struct {
	Migrations string `env:"MIGRATIONS_PATH" env-default:"./pkg/postgres/migrations"`
	Static     string `env:"STATIC_PATH" env-default:"./web/static"`
	Templates  string `env:"TEMPLATES_PATH" env-default:"./web/templates"`
}

func (p *Path) ToMigrations() string {
	migrationsAbsPath, _ := filepath.Abs(p.Migrations)
	migrationsSlashPath := filepath.ToSlash(migrationsAbsPath)
	return fmt.Sprintf("file://%s", migrationsSlashPath)
}

func (p *Path) ToStatic() string {
	staticAbsPath, _ := filepath.Abs(p.Static)
	staticSlashPath := filepath.ToSlash(staticAbsPath)
	return staticSlashPath
}

func (p *Path) ToTemplates() string {
	templateAbsPath, _ := filepath.Abs(p.Templates)
	templateSlashPath := filepath.ToSlash(templateAbsPath)
	return templateSlashPath
}

type Config struct {
	Env string `env:"APP_ENV" env-default:"development"`
	DB  struct {
		Addr string `env:"DB_ADDR" env-default:"postgres://postgres:your_password_here@localhost:5432/forum-database?sslmode=disable"`
	}
	Server struct {
		Port         string `env:"SERVER_PORT" env-default:"8080"`
		ReadTimeout  int    `env:"SERVER_READ_TIMEOUT" env-default:"5"`
		WriteTimeout int    `env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
		IdleTimeout  int    `env:"SERVER_IDLE_TIMEOUT" env-default:"15"`
	}
	JWT struct {
		Secret     string `env:"JWT_SECRET" env-default:"your_secret_key_here"`
		Expiration int    `env:"JWT_EXPIRATION_HOURS" env-default:"24"`
	}
	Path Path
}

func New() (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(".env", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
