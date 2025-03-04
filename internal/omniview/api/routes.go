package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/storage"
)

func AddRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
) {
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	mux.Handle("/", loggingMiddleware(handleIndex(templates)))
}

type Content struct {
	Head struct {
		Title string
	}
	Title string
	Posts []storage.GetPostsPagedRow
}

func handleIndex(t *templates.Templates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Demo posts data
		content := Content{
			Head: struct {
				Title string
			}{
				Title: "Omni | Home",
			},
			Posts: []storage.GetPostsPagedRow{
				{
					Username: "harrydayexe",
					Post: storage.Post{
						Title:       "Test Post",
						Description: "This is a test post",
						CreatedAt:   time.Now(),
					},
				},
				{
					Username: "smellysprite",
					Post: storage.Post{
						Title:       "Test Post 2",
						Description: "This is a second test post",
						CreatedAt:   time.Now(),
					},
				},
			},
		}

		t.Templates.ExecuteTemplate(w, "outline.html", content)
	})
}
