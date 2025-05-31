package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

func Connect(user, password, host, port, name, migrationsPath string) (*sql.DB, error) {
	// Connecting to a database
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	conn, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	// Ensuring a connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Migrate
	m, err := migrate.New(migrationsPath, url)
	if err != nil {
		return nil, err
	}

	// Apply migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply.")
	} else {
		log.Println("ToMigrations applied successfully!")
	}

	return conn, nil
}
