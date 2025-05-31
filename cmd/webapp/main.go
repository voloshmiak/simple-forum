package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/application"
	"forum-project/internal/config"
	"forum-project/internal/route"
	"forum-project/pkg/postgres"
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

	// Database connection and migration
	conn, err := postgres.Connect(cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
		cfg.Path.ToMigrations())
	if err != nil {
		return err
	}

	// Application
	app := application.NewApp(conn, cfg)

	// Register routes
	mux := route.RegisterRoutes(app)

	// Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool, 1)

	go gracefulShutdown(done, server, conn)

	log.Printf("Starting server on port %s", server.Addr)

	// Run server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed to start: %w", err)
	}

	// Wait for a graceful shutdown to complete
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
