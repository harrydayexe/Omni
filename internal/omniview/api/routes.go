package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/oxtoacart/bpool"
)

func addRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	cfg config.ViewConfig,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewJwtSecret(cfg.JWTSecret),
	)

	mux.Handle("GET /", stack(handleGetIndex(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /user/{id}", stack(handleGetUser(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /post/{id}", stack(handleGetPost(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /login", stack(handleGetLogin(templates, bufpool, logger)))
}

func handleGetIndex(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetIndexPage(templates, dataConnector, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handleGetUser(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetUserPage(templates, dataConnector, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handleGetPost(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetPostPage(templates, dataConnector, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handleGetLogin(
	templates *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetLogin(templates, bufpool, logger).ServeHTTP(w, r)
		}
	})
}
