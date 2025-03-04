package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniauth/api"
	"github.com/harrydayexe/Omni/internal/storage"
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

	cfg, err := env.ParseAs[config.AuthConfig]()
	if err != nil {
		logger.Error("failed to parse config", slog.Any("error", err))
		panic(err)
	}
	logger.Info("config", slog.Any("config", cfg))

	db, err := cmd.GetDBConnection(cfg.DatabaseConfig)
	if err != nil {
		logger.Error("failed to connect to database: %v", slog.Any("error", err))
		panic(err)
	}

	queries := storage.New(db)
	authService := auth.NewAuthService([]byte(cfg.JWTSecret), queries, logger)

	if err := cmd.Run(ctx, api.NewHandler(logger, authService, db), os.Stdout, cfg.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
