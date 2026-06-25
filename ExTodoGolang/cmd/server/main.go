package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/http"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/storage"
	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/logger"
)

func main() {
	log := logger.New()

	cfg := config.Load()
	log.Info("Starting Todo App", "version", "0.1.0", "env", cfg.Environment)

	// Initialize repository
	repo, err := storage.NewSQLiteRepository(cfg.DBPath)
	if err != nil {
		log.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}
	defer repo.Close()

	// HTTP server setup
	server := http.NewServer(log, cfg, repo)

	// Graceful shutdown (SPEC-007)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("HTTP server starting", "addr", cfg.HTTPAddr)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && err.Error() != "http: Server closed" {
			log.Error("Server error", "err", err)
			os.Exit(1)
		}
	case sig := <-sigChan:
		log.Info("Received signal", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Shutdown error", "err", err)
			os.Exit(1)
		}
	}

	log.Info("Server shut down gracefully")
}
