package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func NewHandler(
	logger *slog.Logger,
	db storage.Querier,
	dbconn utilities.PingableDB,
	snowflakeGenerator *snowflake.SnowflakeGenerator,
) http.Handler {
	mux := http.NewServeMux()

	// TODO: Add routes here
	AddUserRoutes(mux, logger, db, snowflakeGenerator)

	var handler http.Handler = mux
	return handler
}
