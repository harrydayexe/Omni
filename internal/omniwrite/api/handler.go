package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	db storage.Querier,
	dbconn utilities.PingableDB,
	authService auth.Authable,
	snowflakeGenerator *snowflake.SnowflakeGenerator,
	config *config.Config,
) http.Handler {
	mux := http.NewServeMux()

	AddUserRoutes(mux, logger, db, snowflakeGenerator, authService, config)
	AddPostRoutes(mux, logger, db, snowflakeGenerator, authService, config)
	AddCommentsRoutes(mux, logger, db, snowflakeGenerator, authService, config)

	utilities.AddHealthCheck(
		mux,
		logger,
		dbconn,
	)

	var handler http.Handler = mux
	return handler
}
