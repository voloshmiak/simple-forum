package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

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
	Path struct {
		ToMigrations string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
		ToStatic     string `env:"STATIC_PATH" env-default:"./web/static"`
		ToTemplates  string `env:"TEMPLATES_PATH" env-default:"./web/templates"`
	}
}

func New() (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(".env", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
