package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/harrydayexe/Omni/internal/loadbalancer"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	fptr := flag.String("config", "config.yaml", "file path to read the config from")
	flag.Parse()

	var logLevel slog.Leveler
	if *verbose {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	config, err := loadbalancer.ReadConfig(*fptr, logger)
	if err != nil {
		logger.Error("Could not read config", slog.String("configFile", *fptr))
		panic("could not read config")
	}

	err = config.IsValid()
	if err != nil {
		logger.Error(err.Error())
		panic("config is not valid")
	}

	router, err := loadbalancer.New(config, logger)
	if err != nil {
		panic("could not create router")
	}

	server := &http.Server{
		Handler: router,
	}

	logger.Info("Starting server")
	err = server.ListenAndServe()
}
