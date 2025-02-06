package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/storage"
)

func AddAuthRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db storage.Querier,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("POST /login", stack(handleLogin(logger)))
	mux.Handle("POST /signup", stack(handleSignup(logger)))
}

func handleLogin(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "login POST request received")
	})
}

func handleSignup(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "signup POST request received")
	})
}
