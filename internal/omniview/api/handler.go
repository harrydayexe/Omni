package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
)

func NewHandler(
	logger *slog.Logger,
	temps *templates.Templates,
	dataConnector connector.Connector,
	config config.ViewConfig,
) http.Handler {
	mux := http.NewServeMux()
	AddPageRoutes(
		mux,
		temps,
		logger,
		dataConnector,
		config,
	)
	templates.AddStaticFileRoutes(mux)
	var handler http.Handler = mux
	return handler
}
