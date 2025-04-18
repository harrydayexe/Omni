package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	db storage.Querier,
	dbconn utilities.PingableDB,
) http.Handler {
	mux := http.NewServeMux()
	AddReadRoutes(
		mux,
		logger,
		db,
	)
	utilities.AddHealthCheck(
		mux,
		logger,
		dbconn,
	)
	var handler http.Handler = mux
	return handler
}
