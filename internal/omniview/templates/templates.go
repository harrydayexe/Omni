package templates

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
)

//go:embed *.html
var embeddedTemplates embed.FS

//go:embed static/*
var staticFiles embed.FS

func AddStaticFileRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
}

type Templates struct {
	Templates *template.Template
}

// New initializes and loads all templates
func New(logger *slog.Logger) *Templates {
	tmpls, err := template.ParseFS(embeddedTemplates, "*.html")
	if err != nil {
		logger.Error("Error parsing templates", slog.Any("error", err))
	}
	return &Templates{Templates: tmpls}
}
