package utilities

import (
	"context"
	"log/slog"
	"net/http"
)

type PingableDB interface {
	PingContext(ctx context.Context) error
}

func AddHealthCheck(
	mux *http.ServeMux,
	logger *slog.Logger,
	db PingableDB,
) {
	mux.Handle("GET /healthz", handleHealthCheck(logger, db))
}

func handleHealthCheck(logger *slog.Logger, db PingableDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "health check GET request received")

		if err := db.PingContext(r.Context()); err != nil {
			logger.ErrorContext(r.Context(), "failed health check", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
