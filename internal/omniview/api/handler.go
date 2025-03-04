package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
)

func NewHandler(
	logger *slog.Logger,
	temps *templates.Templates,
	dataConnector connector.Connector,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		temps,
		logger,
		dataConnector,
	)
	templates.AddStaticFileRoutes(mux)
	var handler http.Handler = mux
	return handler
}
