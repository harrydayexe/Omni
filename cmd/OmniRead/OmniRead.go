package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/omniread/api"
	"github.com/harrydayexe/Omni/internal/omniread/config"
	"github.com/harrydayexe/Omni/internal/storage"
)

func main() {
	ctx := context.Background()
	verbose := flag.Bool("v", false, "verbose")
	fptr := flag.String("config", "config/OmniRead/dev.yml", "file path to read the config from")
	flag.Parse()

	var logLevel slog.Leveler
	if *verbose {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	cfg, err := config.NewConfig(*fptr)
	if err != nil {
		logger.Error("failed to load config: %v", err)
		panic(err)
	}

	db, err := GetDBConnection(*cfg)
	if err != nil {
		logger.Error("failed to connect to database: %v", err)
		panic(err)
	}

	ur := storage.NewUserRepo(db, logger)
	pr := storage.NewPostRepo(db)
	cr := storage.NewCommentRepo(db, logger)

	if err := cmd.Run(ctx, api.NewHandler(logger, ur, pr, cr), os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func GetDBConnection(config config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.Database.DriverName, config.Database.DataSourceName)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.Database.ConnMaxLifetime))
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
