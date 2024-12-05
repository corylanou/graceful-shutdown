package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopherguides/graceful-shutdown/service"
)

// Run starts the service and handles shutdown
func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s := service.New(ctx, ":8080")

	// start the service in a go routine
	errChan := make(chan error, 1)
	go func() {
		if err := s.Run(); err != nil {
			errChan <- fmt.Errorf("service error: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		slog.Error("service failed to start", "error", err)
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Handle shutdown in goroutine to prevent blocking
	shutdownErrCh := make(chan error, 1)
	go func() {
		slog.Info("initiating graceful shutdown")
		shutdownErrCh <- s.Close(shutdownCtx)
	}()

	// Create shutdown monitoring
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigChan:
			slog.Warn("forced shutdown initiated")
			cancel()
		case <-shutdownCtx.Done():
		}
	}()

	// Wait for either shutdown completion or timeout
	select {
	case err := <-shutdownErrCh:
		if err != nil {
			return fmt.Errorf("error during shutdown: %w", err)
		}
		slog.Info("shutdown completed successfully")
	case <-shutdownCtx.Done():
		return fmt.Errorf("shutdown timed out: %w", shutdownCtx.Err())
	}

	return nil
}
