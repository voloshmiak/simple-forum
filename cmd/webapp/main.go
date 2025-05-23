package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-project/internal/application"
	"forum-project/internal/database"
	"forum-project/internal/middleware"
	"forum-project/internal/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	// load env variables
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	// connect to database
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	// Application
	app := application.NewApp(conn)

	// Register routes
	mux := routes.Register(app)

	// Server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
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
