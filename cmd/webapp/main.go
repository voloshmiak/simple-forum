package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"simple-forum/internal/app"
	"simple-forum/internal/config"
	"simple-forum/internal/database"
	"simple-forum/internal/router"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Config
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// DB connection and migration
	conn, err := database.Connect(cfg.DB.User, cfg.DB.Password, cfg.DB.Host,
		cfg.DB.Port, cfg.DB.Name, cfg.Path.ToMigrations())
	if err != nil {
		return err
	}

	// App
	a := app.New(conn, cfg)

	// Router
	r := router.RegisterRoutes(a)

	// Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool, 1)

	go gracefulShutdown(done, server, conn)

	a.Logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	// Run server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed to start: %w", err)
	}

	// Wait for a graceful shutdown to complete
	<-done

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
		log.Println("Server forced to shutdown: " + err.Error())
	}

	err := conn.Close()
	if err != nil {
		log.Println("Failed to close database: " + err.Error())
	}

	log.Println("Server shutdown complete")

	done <- true
}
