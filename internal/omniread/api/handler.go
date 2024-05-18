package api

import (
	"log/slog"
	"net/http"
)

func NewHandler(
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		logger,
	)
	var handler http.Handler = mux
	return handler
}
