package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/kw510/z/pkg/db"
	"github.com/kw510/z/pkg/srv"
	"github.com/kw510/z/pkg/srv/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	// Main running context during app lifecycle
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Initialize database
	if err := db.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	srv, err := srv.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	srv = middleware.WithCORS(srv)

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(
		":8080",
		srv,
	)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
