package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutDownTimeout = 10 * time.Second
	httpReadTimeout = 5 * time.Second
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	mux := http.NewServeMux()

	// Register routes
	mux.Handle("GET /{$}", ipInfo{})
	mux.HandleFunc("/", notFound)

	server := &http.Server{
		Addr:        *addr,
		Handler:     mux,
		ReadTimeout: httpReadTimeout,
	}

	go func() {
		slog.Info("Starting HTTP server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("err", err.Error()))
	}

	slog.Info("Server exiting")
}
