package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniwrite/api"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func main() {
	ctx := context.Background()
	cfg, err := env.ParseAs[config.WriteConfig]()
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

	nodeId, err := utilities.GetNodeIDFromDeployment(logger, cfg.NodeName)
	if err != nil {
		logger.Error("Failed to get node id", slog.Any("error", err))
		panic(fmt.Errorf("failed to get node id: %w", err))
	}

	// Create snowflake generator
	snowflakeGenerator := snowflake.NewSnowflakeGenerator(uint16(nodeId))

	db, err := cmd.GetDBConnection(cfg.DatabaseConfig)
	if err != nil {
		logger.Error("failed to connect to database: %v", slog.Any("error", err))
		panic(err)
	}

	queries := storage.New(db)
	authService := auth.NewAuthService([]byte(cfg.JWTSecret), queries, logger)

	if err := cmd.Run(ctx, api.NewHandler(logger, queries, db, authService, snowflakeGenerator, &cfg.Config), os.Stdout, cfg.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
