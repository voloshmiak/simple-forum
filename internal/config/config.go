package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"APP_ENV" env-default:"development"`
	Database struct {
		User     string `env:"DB_USER" env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:""`
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     string `env:"DB_PORT" env-default:"5432"`
		Name     string `env:"DB_NAME" env-default:"forum_database"`
	}
	Server struct {
		Port         string `env:"SERVER_PORT" env-default:"8080"`
		ReadTimeout  int    `env:"SERVER_READ_TIMEOUT" env-default:"5"`
		WriteTimeout int    `env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
		IdleTimeout  int    `env:"SERVER_IDLE_TIMEOUT" env-default:"15"`
	}
	JWT struct {
		Secret     string `env:"JWT_SECRET" env-default:""`
		Expiration int    `env:"JWT_EXPIRATION_HOURS" env-default:"24"`
	}
	Path struct {
		Migrations string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
		Static     string `env:"STATIC_PATH" env-default:"./web/static"`
		Templates  string `env:"TEMPLATES_PATH" env-default:"./web/templates"`
	}
}

func Load() (*Config, error) {
	config := new(Config)
	// Environment variables
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
