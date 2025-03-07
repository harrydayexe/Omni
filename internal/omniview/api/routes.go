package api

import (
	"context"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
	"github.com/microcosm-cc/bluemonday"
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
	mux.Handle("GET /post/{id}", loggingMiddleware(handleGetPost(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /login", loggingMiddleware(handleGetLogin(templates, dataConnector, bufpool, logger)))
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
		Head       datamodels.Head
		Error      string
		IsUserPage bool
		Posts      []storage.GetPostsPagedRow
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for index")
		// Get posts
		posts, err := dataConnector.GetMostRecentPosts(r.Context(), 0)
		content := Content{
			Head: datamodels.Head{
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
		Head       datamodels.Head
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
			Head: datamodels.Head{
				Title: "Omni | User",
			},
			IsUserPage: true,
		}

		// Create error channel
		errChan := make(chan error, 2)
		// Create sub-context
		subctx, cancel := context.WithCancel(r.Context())

		// Get user
		go func() {
			user, err := dataConnector.GetUser(subctx, snowflake)
			if err != nil {
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}
			logger.InfoContext(r.Context(), "User data fetched", slog.Int64("id", int64(snowflake.ToInt())))
			content.Head.Title = "Omni | " + user.Username
			content.Username = user.Username
			errChan <- nil
		}()

		// Get user posts
		go func() {
			posts, err := dataConnector.GetUserPosts(subctx, snowflake)
			if err != nil {
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}

			logger.InfoContext(r.Context(), "User posts data fetched", slog.Int64("id", int64(snowflake.ToInt())))
			content.Posts = posts
			errChan <- nil
		}()

		// Wait for goroutines to finish
		var firstErr error
		for i := 0; i < 2; i++ {
			// If the error channel has an error and it is the first error...
			if err := <-errChan; err != nil && firstErr == nil {
				logger.InfoContext(r.Context(), "An error occurred while fetching data", slog.String("error", err.Error()))
				firstErr = err
			}
		}
		logger.InfoContext(r.Context(), "DB Calls completed", slog.String("Handler", "handleGetUser"))

		// Handle firstErr
		var ae *connector.APIError
		if errors.As(firstErr, &ae) {
			if ae.StatusCode == http.StatusNotFound {
				// User not found
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("User not found", "The user you are looking for does not exist."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while fetching user data."),
				)
			}
			return
		}

		writeTemplateWithBuffer(r.Context(), logger, "user.html", t, bufpool, w, content)
	})
}

func handleGetPost(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Post struct {
		Title       string
		Description string
		CreatedAt   string
		Author      string
		Content     template.HTML
	}
	type Content struct {
		Head datamodels.Head
		Post Post
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /post", slog.String("id", r.PathValue("id")))
		// Parse user id
		snowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			content := datamodels.NewErrorPageModel(
				"Post not found",
				"The post you are looking for does not exist.",
			)
			writeTemplateWithBuffer(r.Context(), logger, "errorpage.html", t, bufpool, w, content)
			return
		}

		content := Content{
			Head: datamodels.Head{
				Title: "Omni | Post",
			},
		}
		var post storage.Post

		// Create error channel
		errChan := make(chan error, 1)
		// Create sub-context
		subctx, cancel := context.WithCancel(r.Context())

		// Get user
		go func() {
			postResp, err := dataConnector.GetPost(subctx, snowflake)
			if err != nil {
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}
			logger.InfoContext(r.Context(), "Post data fetched", slog.Int64("id", int64(snowflake.ToInt())))
			post = postResp
			// Finally return with no error
			errChan <- nil
		}()

		// Wait for goroutines to finish
		var firstErr error
		for i := 0; i < 1; i++ {
			// If the error channel has an error and it is the first error...
			if err := <-errChan; err != nil && firstErr == nil {
				logger.InfoContext(r.Context(), "An error occurred while fetching data", slog.String("error", err.Error()))
				firstErr = err
			}
		}
		logger.InfoContext(r.Context(), "DB Calls completed", slog.String("Handler", "handleGetPost"))

		// Handle firstErr
		var ae *connector.APIError
		if errors.As(firstErr, &ae) {
			if ae.StatusCode == http.StatusNotFound {
				// User not found
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Post not found", "The post you are looking for does not exist."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while fetching post data."),
				)
			}
			return
		}

		// Get markdown data
		html, err := fetchMarkdownData(r.Context(), logger, post.MarkdownUrl)
		if errors.As(err, &ae) {
			if ae.StatusCode == http.StatusNotFound {
				// User not found
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Markdown not found", "The markdown file for this post could not be read."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while processing markdown data."),
				)
			}
			return
		}

		logger.InfoContext(r.Context(), "Setting page content")
		content.Post = Post{
			Title:       post.Title,
			Description: post.Description,
			CreatedAt:   post.CreatedAt.Format(time.DateTime),
			Author:      "Author", // TODO: Do a call to get this info
			Content:     template.HTML(html),
		}

		writeTemplateWithBuffer(r.Context(), logger, "post.html", t, bufpool, w, content)
	})
}

func fetchMarkdownData(ctx context.Context, logger *slog.Logger, url string) (string, error) {
	logger.InfoContext(ctx, "Fetching markdown data", slog.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		logger.InfoContext(ctx, "Failed to fetch markdown data", slog.String("error", err.Error()))
		return "", connector.NewAPIError(0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.InfoContext(ctx, "Failed to fetch markdown data", slog.Int("status", resp.StatusCode))
		return "", connector.NewAPIError(resp.StatusCode, nil)
	}

	rawMdBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.InfoContext(ctx, "Failed to read markdown data", slog.String("error", err.Error()))
		return "", connector.NewAPIError(0, err)
	}

	logger.InfoContext(ctx, "Markdown data fetched", slog.Int("length", len(rawMdBytes)))
	maybeUnsafeHTML := markdown.ToHTML(rawMdBytes, nil, nil)
	html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)

	return string(html), nil
}

func handleGetLogin(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head datamodels.Head
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /login")
		content := Content{
			Head: datamodels.Head{
				Title: "Omni | Login",
			},
		}

		writeTemplateWithBuffer(r.Context(), logger, "login.html", t, bufpool, w, content)
	})
}
