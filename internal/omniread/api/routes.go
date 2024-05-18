package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
)

func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	mux.Handle("/", stack(handleIndex()))
}

func handleIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{'response': 'Hello, World!'}"))
	})
}
