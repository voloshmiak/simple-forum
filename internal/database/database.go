package database

import (
	"database/sql"
	"fmt"
	"os"
)

func New() (*sql.DB, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		name     = os.Getenv("DB_NAME")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		url      = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, name, user, password)
	)
	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
