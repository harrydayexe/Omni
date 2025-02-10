package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	authService auth.Authable,
	dbconn utilities.PingableDB,
) http.Handler {
	mux := http.NewServeMux()
	AddAuthRoutes(
		mux,
		logger,
		authService,
	)
	utilities.AddHealthCheck(
		mux,
		logger,
		dbconn,
	)
	var handler http.Handler = mux
	return handler
}
