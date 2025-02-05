package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	db storage.Querier,
	dbconn utilities.PingableDB,
	snowflakeGenerator *snowflake.SnowflakeGenerator,
	config *config.Config,
) http.Handler {
	mux := http.NewServeMux()

	AddUserRoutes(mux, logger, db, snowflakeGenerator, config)
	AddPostRoutes(mux, logger, db, snowflakeGenerator, config)

	var handler http.Handler = mux
	return handler
}
