package api

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
	"github.com/oxtoacart/bpool"
)

func handleGetIndexPage(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head       datamodels.Head
		NavBar     datamodels.NavBar
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
			NavBar:     datamodels.NewNavBar(r.Context()),
			IsUserPage: false,
		}
		if err != nil {
			content.Error = "An error occurred while fetching the most recent posts. Try again later."
		}

		// Demo posts data
		content.Posts = posts

		// Write template
		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "posts.html", t, bufpool, w, content)
	})
}

func handleGetUserPage(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head       datamodels.Head
		NavBar     datamodels.NavBar
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
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "errorpage.html", t, bufpool, w, content)
			return
		}

		// Create content variable
		content := Content{
			Head: datamodels.Head{
				Title: "Omni | User",
			},
			NavBar:     datamodels.NewNavBar(r.Context()),
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
			logger.DebugContext(r.Context(), "User data fetched", slog.Int64("id", int64(snowflake.ToInt())))
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

			logger.DebugContext(r.Context(), "User posts data fetched", slog.Int64("id", int64(snowflake.ToInt())))
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
		logger.DebugContext(r.Context(), "DB Calls completed", slog.String("Handler", "handleGetUser"))

		// Handle firstErr
		var ae *connector.APIError
		if errors.As(firstErr, &ae) {
			if ae.StatusCode == http.StatusNotFound {
				// User not found
				writeTemplateWithBuffer(
					r.Context(), logger, http.StatusOK,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("User not found", "The user you are looking for does not exist."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger, http.StatusOK,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while fetching user data."),
				)
			}
			return
		}

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "user.html", t, bufpool, w, content)
	})
}

func handleGetPostPage(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Post struct {
		Title       string
		Description string
		CreatedAt   string
		Author      string
		AuthorID    int64
		Content     template.HTML
	}
	type Content struct {
		Head   datamodels.Head
		NavBar datamodels.NavBar
		Post   Post
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /post", slog.String("id", r.PathValue("id")))
		// Parse user id
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			content := datamodels.NewErrorPageModel(
				"Post not found",
				"The post you are looking for does not exist.",
			)
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "errorpage.html", t, bufpool, w, content)
			return
		}

		content := Content{
			Head: datamodels.Head{
				Title: "Omni | Post",
			},
			NavBar: datamodels.NewNavBar(r.Context()),
		}
		var post storage.Post
		var user storage.User

		// Create error channel
		errChan := make(chan error, 2)
		userIdChan := make(chan int64)
		// Create sub-context
		subctx, cancel := context.WithCancel(r.Context())

		// Get post
		go func() {
			postResp, err := dataConnector.GetPost(subctx, postSnowflake)
			if err != nil {
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}
			logger.DebugContext(r.Context(), "Post data fetched", slog.Int64("id", int64(postSnowflake.ToInt())))
			post = postResp
			// Finally return with no error
			errChan <- nil
			userIdChan <- post.UserID
		}()

		// Get user
		go func() {
			select {
			case <-subctx.Done():
				errChan <- subctx.Err()
				return
			case userId := <-userIdChan:
				userResp, err := dataConnector.GetUser(subctx, snowflake.ParseId(uint64(userId)))
				if err != nil {
					// Cancel other routines
					cancel()
					errChan <- err
					return
				}
				logger.DebugContext(r.Context(), "User data fetched", slog.Int64("id", userId))
				user = userResp
				errChan <- nil
			}
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
		logger.DebugContext(r.Context(), "DB Calls completed", slog.String("Handler", "handleGetPost"))

		// Handle firstErr
		var ae *connector.APIError
		if errors.As(firstErr, &ae) {
			if ae.StatusCode == http.StatusNotFound {
				// User not found
				writeTemplateWithBuffer(
					r.Context(), logger, http.StatusOK,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Post not found", "The post you are looking for does not exist."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger, http.StatusOK,
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
					r.Context(), logger, http.StatusOK,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Markdown not found", "The markdown file for this post could not be read."),
				)
			} else {
				// Error in getting backend data
				writeTemplateWithBuffer(
					r.Context(), logger, http.StatusOK,
					"errorpage.html", t, bufpool, w,
					datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while processing markdown data."),
				)
			}
			return
		}

		logger.DebugContext(r.Context(), "Setting page content")
		content.Post = Post{
			Title:       post.Title,
			Description: post.Description,
			CreatedAt:   post.CreatedAt.Format(time.DateTime),
			Author:      user.Username,
			AuthorID:    user.ID,
			Content:     template.HTML(html),
		}

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "post.html", t, bufpool, w, content)
	})
}

func handleGetLoginPage(
	t *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /login")

		content := datamodels.NewFormPage(r.Context(), "Login")

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login.html", t, bufpool, w, content)
	})
}

func handleGetCreatePostPage(
	templates *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /post/new")

		content := datamodels.NewFormPage(r.Context(), "New Post")
		content.NavBar.IsLoggedIn = true

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "newpost.html", templates, bufpool, w, content)
	})
}
