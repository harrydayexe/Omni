package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniread/api"
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

	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		logger.Error("failed to parse config", slog.Any("error", err))
	}
	logger.Info("config", slog.Any("config", cfg))

	db, err := GetDBConnection(cfg)
	if err != nil {
		logger.Error("failed to connect to database: %v", slog.Any("error", err))
		panic(err)
	}

	queries := storage.New(db)

	if err := cmd.Run(ctx, api.NewHandler(logger, queries, db), os.Stdout, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func GetDBConnection(config config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DataSourceName)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnMaxLifetime))
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
