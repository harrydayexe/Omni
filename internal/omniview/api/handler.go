package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/oxtoacart/bpool"
)

func NewHandler(
	logger *slog.Logger,
	temps *templates.Templates,
	dataConnector connector.Connector,
	config config.ViewConfig,
) http.Handler {
	var bufpool *bpool.BufferPool = bpool.NewBufferPool(64)
	mux := http.NewServeMux()
	addRoutes(
		mux,
		temps,
		logger,
		dataConnector,
		bufpool,
		config,
	)
	templates.AddStaticFileRoutes(mux)
	var handler http.Handler = mux
	return handler
}
