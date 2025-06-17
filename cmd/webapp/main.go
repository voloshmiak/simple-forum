package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

	a.Logger.Info(fmt.Sprintf("Starting server on port %s", server.Addr))

	signs := make(chan os.Signal)

	// Run server
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Error("Server failed to start", "error", err)
			signs <- syscall.SIGTERM
		}
	}()

	signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)

	<-signs
	a.Logger.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.Logger.Error("Server forced to shutdown", "error", err.Error())
	}

	err = conn.Close()
	if err != nil {
		a.Logger.Error("Failed to close database", "error", err.Error())
	}

	a.Logger.Info("Graceful shutdown complete")

	return nil
}
