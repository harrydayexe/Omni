package api

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
	"github.com/oxtoacart/bpool"
)

func AddRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
	dataConnector connector.Connector,
) {
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	var bufpool *bpool.BufferPool = bpool.NewBufferPool(64)

	mux.Handle("GET /", loggingMiddleware(handleGetIndex(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /user/{id}", loggingMiddleware(handleGetUser(templates, dataConnector, bufpool, logger)))
}

func writeTemplateWithBuffer(ctx context.Context, logger *slog.Logger, name string, t *templates.Templates, bufpool *bpool.BufferPool, w http.ResponseWriter, content interface{}) {
	// Get buffer
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := t.Templates.ExecuteTemplate(buf, name, content)
	if err != nil {
		logger.ErrorContext(ctx, "Error executing template", slog.String("error message", err.Error()))
		// NOTE: We are assuming here that this error page won't fail to render
		_ = t.Templates.ExecuteTemplate(w, "errorpage.html", datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while rendering the page."))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func handleGetIndex(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head struct {
			Title string
		}
		Error      string
		IsUserPage bool
		Posts      []storage.GetPostsPagedRow
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for index")
		// Get posts
		posts, err := dataConnector.GetMostRecentPosts(r.Context(), 0)
		content := Content{
			Head: struct {
				Title string
			}{
				Title: "Omni | Home",
			},
			IsUserPage: false,
		}
		if err != nil {
			content.Error = "An error occurred while fetching the most recent posts. Try again later."
		}

		// Demo posts data
		content.Posts = posts

		// Write template
		writeTemplateWithBuffer(r.Context(), logger, "posts.html", t, bufpool, w, content)
	})
}

func handleGetUser(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head struct {
			Title string
		}
		Error      string
		Username   string
		IsUserPage bool
		Posts      []storage.Post
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /user", slog.String("id", r.PathValue("id")))
		// Parse user id
		snowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			content := datamodels.NewErrorPageModel(
				"User not found",
				"The user you are looking for does not exist.",
			)
			writeTemplateWithBuffer(r.Context(), logger, "errorpage.html", t, bufpool, w, content)
			return
		}

		// Create content variable
		content := Content{
			Head: struct {
				Title string
			}{
				Title: "Omni | User",
			},
			IsUserPage: true,
		}

		// Create waitGroup
		var wg sync.WaitGroup

		// Get user
		wg.Add(1)
		go func() {
			defer wg.Done()
			user, err := dataConnector.GetUser(r.Context(), snowflake)
			if err != nil {
				// TODO: Handle this error
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			content.Head.Title = "Omni | " + user.Username
			content.Username = user.Username
		}()

		// Get user posts
		wg.Add(1)
		go func() {
			defer wg.Done()
			posts, err := dataConnector.GetUserPosts(r.Context(), snowflake)
			if err != nil {
				// TODO: Handle this error
				http.Error(w, "No posts found for user", http.StatusNotFound)
				return
			}

			content.Posts = posts
		}()

		// Wait for goroutines to finish
		wg.Wait()
		logger.InfoContext(r.Context(), "User data fetched", slog.Any("content", content))

		writeTemplateWithBuffer(r.Context(), logger, "user.html", t, bufpool, w, content)
	})
}
