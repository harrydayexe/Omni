package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
)

func AddRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
) {
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	// mux.Handle("/images", http.FileServer(http.Dir("../../web/static/images")))
	// mux.Handle("/styles.css", http.FileServer(http.Dir("../../web/static/styles.css")))
	mux.Handle("/", loggingMiddleware(handleIndex(templates)))
	mux.Handle("/hello", loggingMiddleware(handleHello(templates)))
}

func handleIndex(t *templates.Templates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Templates.ExecuteTemplate(w, "index.html", nil)
	})
}

func handleHello(t *templates.Templates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Templates.ExecuteTemplate(w, "hello.html", nil)
	})
}
