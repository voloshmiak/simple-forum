package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"path/filepath"
)

func Migrate(addr, path string) error {
	migrationsAbsPath, _ := filepath.Abs(path)
	migrationsSlashPath := filepath.ToSlash(migrationsAbsPath)
	migrationsPath := fmt.Sprintf("file://%s", migrationsSlashPath)

	m, err := migrate.New(migrationsPath, addr)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
