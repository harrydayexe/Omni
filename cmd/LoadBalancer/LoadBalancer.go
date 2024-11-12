package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"net/url"
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

	config, err := loadbalancer.ReadConfig(*fptr)
	if err != nil {
		logger.Error("Could not read config", slog.String("configFile", *fptr))
		panic("could not read config")
	}

	err = config.IsValid()
	if err != nil {
		logger.Error(err.Error())
		panic("config is not valid")
	}

	router := http.NewServeMux()
	for _, l := range config.Location {
		urls, err := convertToUrl(l.ProxyPass)
		if err != nil {
			logger.Error(err.Error())
			panic("could not convert locations to urls")
		}

		httpProxy, err := loadbalancer.NewLoadBalancerProxy(config.Algorithm, urls)
		if err != nil {
			logger.Error(err.Error())
			panic("could not create load balancer proxy")
		}
		httpProxy.StartHealthCheck(context.Background(), 10)
		router.Handle(l.Pattern, httpProxy)
	}

	server := &http.Server{
		Handler: router,
	}

	logger.Info("Starting server")
	err = server.ListenAndServe()
}

func convertToUrl(locations []string) ([]*url.URL, error) {
	returnSet := make([]*url.URL, len(locations))
	for i, v := range locations {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		returnSet[i] = u
	}

	return returnSet, nil
}
