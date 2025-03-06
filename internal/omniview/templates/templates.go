package templates

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
)

//go:embed partials/*.html
var partialEmbeddedTemplates embed.FS

//go:embed pages/*.html
var pageEmbeddedTemplates embed.FS

//go:embed static/*
var staticFiles embed.FS

func AddStaticFileRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("GET /static/", http.FileServer(http.FS(staticFiles)))
}

type Templates struct {
	Templates *template.Template
}

// New initializes and loads all templates
func New(logger *slog.Logger) (*Templates, error) {
	tmpls, err := template.ParseFS(pageEmbeddedTemplates, "pages/*.html")
	if err != nil {
		logger.Error("Error parsing page templates", slog.Any("error", err))
		return nil, errors.Wrap(err, "template.ParseFS")
	}
	tmpls, err = tmpls.ParseFS(partialEmbeddedTemplates, "partials/*.html")
	if err != nil {
		logger.Error("Error parsing partial templates", slog.Any("error", err))
		return nil, errors.Wrap(err, "template.ParseFS")
	}

	return &Templates{Templates: tmpls}, nil
}
