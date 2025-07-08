package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(addr string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", addr)
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

	return conn, nil
}

func Migrate(addr, path string) error {
	migrationsAbsPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	migrationsSlashPath := filepath.ToSlash(migrationsAbsPath)
	migrationsPath := fmt.Sprintf("file://%s", migrationsSlashPath)

	m, err := migrate.New(migrationsPath, addr)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
