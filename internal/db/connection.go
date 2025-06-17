package db

import (
	"context"
	"database/sql"
	"time"
)

func NewConnection(addr string) (*sql.DB, error) {
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
