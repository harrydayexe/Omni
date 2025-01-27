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

	// TODO: Add routes here

	var handler http.Handler = mux
	return handler
}
