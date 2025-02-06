package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	db storage.Querier,
	dbconn utilities.PingableDB,
	cfg config.Config,
) http.Handler {
	mux := http.NewServeMux()
	AddAuthRoutes(
		mux,
		logger,
		db,
		auth.NewAuthService([]byte(cfg.JWTSecret), db, logger),
	)
	utilities.AddHealthCheck(
		mux,
		logger,
		dbconn,
	)
	var handler http.Handler = mux
	return handler
}
