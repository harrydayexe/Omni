package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniread/api"
	"github.com/harrydayexe/Omni/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg, err := env.ParseAs[config.DatabaseConfig]()
	if err != nil {
		panic(err)
	}

	var logLevel slog.Leveler
	if cfg.VerboseMode {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	logger.Info("Config", slog.Any("config", cfg))

	db, err := cmd.GetDBConnection(cfg)
	if err != nil {
		logger.Error("failed to connect to database: %v", slog.Any("error", err))
		panic(err)
	}

	queries := storage.New(db)

	if err := cmd.Run(ctx, api.NewHandler(logger, queries, db), os.Stdout, cfg.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
