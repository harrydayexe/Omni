package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	datamodelsread "github.com/harrydayexe/Omni/internal/omniread/datamodels"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
	"github.com/oxtoacart/bpool"
)

func handleGetIndexPage(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
	isHTMXRequest bool,
) http.Handler {
	type Content struct {
		Head     datamodels.Head
		NavBar   datamodels.NavBar
		AllPosts datamodels.AllPosts
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for index", slog.Bool("isHTMXRequest", isHTMXRequest))

		// Attempt to get page url query param
		var pageNum int = 1
		pageQuery := r.URL.Query().Get("page")
		if pageQuery != "" {
			logger.DebugContext(r.Context(), "Page number specified", slog.String("pageQueryParam", pageQuery))
			pageNumTemp, err := strconv.Atoi(pageQuery)
			if err == nil {
				logger.DebugContext(r.Context(), "Page number parsed", slog.Int("pageNum", pageNumTemp))
				pageNum = pageNumTemp
			}
		}

		// Get posts
		var errorMsg string
		posts, err := dataConnector.GetMostRecentPosts(r.Context(), pageNum)
		if err != nil {
			errorMsg = "An error occurred while fetching the most recent posts. Try again later."
		}

		content := Content{
			Head: datamodels.Head{
				Title: "Omni | Home",
			},
			NavBar: datamodels.NewNavBar(r.Context()),
			AllPosts: datamodels.NewAllPosts(
				errorMsg, GetUserIdFromCtx(r.Context()),
				posts, false,
				pageNum-1, pageNum+1,
			),
		}

		// Write template
		if isHTMXRequest {
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "posts", t, bufpool, w, content.AllPosts)
		} else {
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "posts.html", t, bufpool, w, content)
		}
	})
}

func handleGetUserPage(t *templates.Templates, dataConnector connector.Connector, bufpool *bpool.BufferPool, logger *slog.Logger) http.Handler {
	type Content struct {
		Head   datamodels.Head
		NavBar datamodels.NavBar
		User   datamodels.User
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
			NavBar: datamodels.NewNavBar(r.Context()),
			User: datamodels.User{
				Username: "",
				AllPosts: datamodels.NewAllPosts(
					"",
					GetUserIdFromCtx(r.Context()),
					datamodelsread.AllPosts{}, true,
					0, 0,
				),
			},
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
			content.User.Username = user.Username
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
			content.User.AllPosts.Posts = utilities.Map(
				posts,
				func(post storage.Post) datamodels.PostListItem {
					return datamodels.PostListItem{
						PostAndUsername: datamodelsread.PostAndUsername{
							Username: "",
							Post:     post,
						},
						IsDeleteable: true,
						IsEditable:   true,
					}
				})
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
	type Content struct {
		Head   datamodels.Head
		NavBar datamodels.NavBar
		Post   datamodels.Post
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
		var commentsModel datamodels.CommentsModel

		// Create error channel
		const numRoutines = 3
		errChan := make(chan error, numRoutines)
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

		// Get Comments
		go func() {
			commentResp, err := dataConnector.GetPostComments(subctx, postSnowflake, 1)
			commentsModel = datamodels.NewCommentsModel(
				err,
				GetUserIdFromCtx(r.Context()),
				commentResp,
				int64(postSnowflake.ToInt()),
				2,
			)
			if err != nil {
				// Don't need to cancel other routines as comments are not core content
				logger.InfoContext(r.Context(), "An error occurred while fetching comments", slog.String("error", err.Error()))
			} else {
				logger.DebugContext(r.Context(), "Comments data fetched", slog.Int64("id", int64(postSnowflake.ToInt())))
			}
			errChan <- nil
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
		for i := 0; i < numRoutines; i++ {
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
		content.Post = datamodels.NewPost(post, user, html, commentsModel)

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
		content.Form.Values["Title"] = "Login"
		content.Form.Values["HXDest"] = "/login"

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
		content.Form.FormMeta["Title"] = "Create A New Post"
		content.Form.FormMeta["URL"] = "/post/new"

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "newpost.html", templates, bufpool, w, content)
	})
}

func handleGetPostEditPage(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /post/edit", slog.String("id", r.PathValue("id")))

		// Parse post id
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			content := datamodels.NewErrorPageModel(
				"Post not found",
				"The post you are looking for does not exist.",
			)
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "errorpage.html", t, bufpool, w, content)
			return
		}

		// Create content variable
		content := datamodels.NewFormPage(r.Context(), "Edit Post")

		// Get current values
		post, err := dataConnector.GetPost(r.Context(), postSnowflake)
		if err != nil {
			errContent := datamodels.NewErrorPageModel(
				"Post details could not be fetched",
				"An error occurred while fetching the post details.",
			)
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "errorpage.html", t, bufpool, w, errContent)
			return
		}

		content.Form.FormMeta["Title"] = "Edit Post"
		content.Form.FormMeta["URL"] = fmt.Sprintf("/post/%d/edit", postSnowflake.ToInt())
		content.Form.Values["Title"] = post.Title
		content.Form.Values["Description"] = post.Description
		content.Form.Values["URL"] = post.MarkdownUrl

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "newpost.html", t, bufpool, w, content)
	})
}

func handleGetSignupPage(
	t *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for /signup")

		content := datamodels.NewFormPage(r.Context(), "Sign Up")
		content.Form.Values["Title"] = "Sign Up"
		content.Form.Values["HXDest"] = "/signup"

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login.html", t, bufpool, w, content)
	})
}
