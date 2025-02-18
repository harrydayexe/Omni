package templates

import (
	"embed"
	"html/template"
	"log/slog"
)

//go:embed *.html
var embeddedTemplates embed.FS

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
