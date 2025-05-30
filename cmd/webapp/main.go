package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/application"
	"forum-project/internal/config"
	"forum-project/internal/env"
	"forum-project/internal/middleware"
	"forum-project/internal/route"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load environment variables
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Database connection
	dataSource := env.GetDataSource()
	conn, err := sql.Open("pgx", dataSource)
	if err != nil {
		return err
	}

	// Ensuring connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		return err
	}

	// Migrate
	m, err := migrate.New(env.GetMigrationPath(), dataSource)
	if err != nil {
		return err
	}

	// Apply migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply.")
	} else {
		log.Println("Migrations applied successfully!")
	}

	// Application
	app := application.NewApp(conn, cfg)

	// Register routes
	mux := route.RegisterRoutes(app)

	// Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      middleware.Logging(mux, app.Logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool, 1)

	go gracefulShutdown(done, server, conn)

	log.Printf("Starting server on port %s", server.Addr)

	// Run server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed to start: %w", err)
	}

	// Wait for graceful shutdown to complete
	<-done
	log.Println("Server shutdown complete")

	return nil
}

func gracefulShutdown(done chan bool, server *http.Server, conn *sql.DB) {
	// Wait for interruption
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	err := conn.Close()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to close database: %s", err))
	}

	done <- true
}
