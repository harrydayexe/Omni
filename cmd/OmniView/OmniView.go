package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/harrydayexe/Omni/internal/omniview/api"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	templateModule "github.com/harrydayexe/Omni/internal/omniview/templates"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, stdout io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.Default()

	srv := NewServer(
		templates.New(logger),
		logger,
	)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}
	go func() {
		logger.Info(
			"server listening",
			slog.String("address", httpServer.Addr),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func NewServer(
	templates *templates.Templates,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()
	api.AddRoutes(
		mux,
		templates,
		logger,
	)
	templateModule.AddStaticFileRoutes(mux)
	var handler http.Handler = mux
	return handler
}
