package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/storage"
)

func NewHandler(
	logger *slog.Logger,
	db *storage.Queries,
) http.Handler {
	mux := http.NewServeMux()
	AddReadRoutes(
		mux,
		logger,
		db,
	)
	var handler http.Handler = mux
	return handler
}
