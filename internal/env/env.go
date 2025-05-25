package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Environment variable %s not set, using fallback '%s'", key, fallback)
	return fallback
}

func GetMigrationPath() string {
	migrationsPath := GetEnv("MIGRATIONS_PATH", "./migrations")
	migrationsAbsPath, _ := filepath.Abs(migrationsPath)
	migrationsSlashPath := filepath.ToSlash(migrationsAbsPath)
	migrationURL := fmt.Sprintf("file://%s", migrationsSlashPath)

	return migrationURL
}

func GetStaticPath() string {
	staticPath := GetEnv("STATIC_PATH", "./web/static")
	staticAbsPath, _ := filepath.Abs(staticPath)
	staticSlashPath := filepath.ToSlash(staticAbsPath)

	return staticSlashPath
}

func GetTemplatePath() string {
	templatePath := GetEnv("TEMPLATES_PATH", "./web/templates")
	templateAbsPath, _ := filepath.Abs(templatePath)
	templateSlashPath := filepath.ToSlash(templateAbsPath)

	return templateSlashPath
}

func GetDataSource() string {
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	name := GetEnv("DB_NAME", "forum_database")

	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)

	return dataSource
}
