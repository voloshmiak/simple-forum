package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple-forum/internal/app"
	"simple-forum/internal/config"
	"simple-forum/internal/db"
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

	// Migrate
	err = db.Migrate(cfg.DB.Addr, cfg.Path.ToMigrations)
	if err != nil {
		return err
	}

	// Connection
	conn, err := db.NewConnection(cfg.DB.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

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

	// Serve

	a.Logger.Info("Starting server on port: " + server.Addr)

	done := make(chan bool)

	go func() {
		signs := make(chan os.Signal)
		signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)

		<-signs
		a.Logger.Info("Shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err = server.Shutdown(shutdownCtx); err != nil {
			a.Logger.Error("Server forced to shutdown", "error", err)
		}

		done <- true
	}()

	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-done
	a.Logger.Info("Graceful shutdown complete")

	return nil
}
