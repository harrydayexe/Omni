package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniview/api"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
)

func main() {
	ctx := context.Background()
	verbose := flag.Bool("v", false, "verbose")

	var logLevel slog.Leveler
	if *verbose {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	cfg, err := env.ParseAs[config.ViewConfig]()
	if err != nil {
		logger.Error("failed to parse config", slog.Any("error", err))
		panic(err)
	}
	logger.Info("config", slog.Any("config", cfg))

	dataConnector := connector.NewAPIConnector(cfg, logger)

	tmpls, err := templates.New(logger)
	if err != nil {
		panic(err)
	}

	if err := cmd.Run(ctx, api.NewHandler(logger, tmpls, dataConnector, cfg), os.Stdout, cfg.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
