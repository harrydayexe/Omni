package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/storage"
)

func AddRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
	dataConnector connector.Connector,
) {
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	mux.Handle("/", loggingMiddleware(handleIndex(templates, dataConnector)))
}

type Content struct {
	Head struct {
		Title string
	}
	Title string
	Error string
	Posts []storage.GetPostsPagedRow
}

func handleIndex(t *templates.Templates, dataConnector connector.Connector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get posts
		posts, err := dataConnector.GetMostRecentPosts(r.Context(), 0)
		content := Content{
			Head: struct {
				Title string
			}{
				Title: "Omni | Home",
			},
		}
		if err != nil {
			content.Error = "An error occurred while fetching the most recent posts. Try again later."
		}

		// Demo posts data
		content.Posts = posts

		t.Templates.ExecuteTemplate(w, "posts.html", content)
	})
}
