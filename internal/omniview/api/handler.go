package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/omniview/templates"
)

func NewHandler(
	logger *slog.Logger,
	temps *templates.Templates,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		temps,
		logger,
	)
	templates.AddStaticFileRoutes(mux)
	var handler http.Handler = mux
	return handler
}
