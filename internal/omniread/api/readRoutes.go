package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

const (
	epoch int64 = 1704067200000
)

// AddReadRoutes adds all api routes to the provided http.ServeMux.
// It also adds logging middleware to each route.
func AddReadRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db storage.Querier,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("GET /post/{id}", stack(handleReadPost(logger, db)))
	mux.Handle("GET /user/{id}", stack(handleReadUser(logger, db)))
	mux.Handle("GET /user/{id}/posts", stack(handleReadUserPosts(logger, db)))
	mux.Handle("GET /post/{id}/comments", stack(handleReadPostComments(logger, db)))
	mux.Handle("GET /posts", stack(handleMostRecentPosts(logger, db)))
	// TODO: Add new routes for things like getting a user and their posts together
}

// route: GET /post/{id}
// return the details of a post by it's id
func handleReadPost(logger *slog.Logger, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read post GET request received")
		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		post, err := db.FindPostByID(r.Context(), int64(id.ToInt()))
		if utilities.IsDbError(r.Context(), logger, w, id, err) {
			return
		}

		utilities.MarshallToResponse(r.Context(), logger, w, post)
	})
}

// route: GET /user/{id}
// return the details of a user by it's id
func handleReadUser(logger *slog.Logger, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read user GET request received")
		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		user, err := db.GetUserByID(r.Context(), int64(id.ToInt()))
		if utilities.IsDbError(r.Context(), logger, w, id, err) {
			return
		}

		utilities.MarshallToResponse(r.Context(), logger, w, user)
	})
}

// TODO: Probably should change this to be paginated instead of time based

// extract the from and limit query parameters from the request
func extractPaginationParams(logger *slog.Logger, r *http.Request, w http.ResponseWriter) (time.Time, int, error) {
	var fromDate time.Time
	var limit int

	fromVal := r.URL.Query().Get("from")
	if fromVal == "" {
		fromDate = time.UnixMilli(epoch)
		// TODO: In this case we may want to return the most recent comments instead of the oldest
	} else {
		var err error
		fromDate, err = time.Parse(time.RFC3339, fromVal)
		if err != nil {
			logger.InfoContext(r.Context(), "failed to parse date from url query param", slog.Any("error", err))
			errorMessage := "Url parameter could not be parsed properly."
			http.Error(w, errorMessage, http.StatusBadRequest)
			return time.Time{}, 0, err
		}
	}
	limitVal := r.URL.Query().Get("limit")
	if limitVal == "" {
		limit = 10
	} else {
		var err error
		limit, err = strconv.Atoi(limitVal)
		if err != nil {
			logger.InfoContext(r.Context(), "failed to parse comment limit from url query param", slog.Any("error", err))
			errorMessage := "Url parameter could not be parsed properly."
			http.Error(w, errorMessage, http.StatusBadRequest)
			return time.Time{}, 0, err
		}
		if limit > 100 {
			limit = 100
		}
	}

	return fromDate, limit, nil
}

// route: GET /user/{id}/posts?from=2006-01-02T15%3A04%3A05Z07%3A00&limit=10
// there are two query parameters from and limit
// from is the date and time to retrieve posts since in RFC3339 format
// limit is the number of posts to return (default to 10, max is 100)
func handleReadUserPosts(logger *slog.Logger, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read posts GET request received")
		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		fromDate, limit, err := extractPaginationParams(logger, r, w)
		if err != nil {
			return
		}

		rows, err := db.GetUserAndPostsByIDPaged(r.Context(), storage.GetUserAndPostsByIDPagedParams{
			ID:           int64(id.ToInt()),
			CreatedAfter: fromDate,
			Limit:        int32(limit),
		})
		if utilities.IsDbError(r.Context(), logger, w, id, err) {
			return
		}

		// Extract the posts from the rows
		posts := make([]storage.Post, len(rows))
		for i, row := range rows {
			posts[i] = row.Post
		}

		utilities.MarshallToResponse(r.Context(), logger, w, posts)
	})
}

// route: GET /post/{id}/comments?from=2006-01-02T15%3A04%3A05Z07%3A00&limit=10
// there are two query parameters from and limit
// from is the date and time to retrieve comments since in RFC3339 format
// limit is the number of comments to return (default to 10, max is 100)
func handleReadPostComments(logger *slog.Logger, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read comments GET request received")
		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		fromDate, limit, err := extractPaginationParams(logger, r, w)
		if err != nil {
			return
		}

		rows, err := db.FindCommentsAndUserByPostIDPaged(r.Context(), storage.FindCommentsAndUserByPostIDPagedParams{
			PostID:       int64(id.ToInt()),
			CreatedAfter: fromDate,
			Limit:        int32(limit),
		})
		if utilities.IsDbError(r.Context(), logger, w, id, err) {
			return
		}

		type CommentReturn struct {
			ID        int64     `json:"id"`
			PostID    int64     `json:"post_id"`
			UserID    int64     `json:"user_id"`
			Username  string    `json:"username"`
			CreatedAt time.Time `json:"created_at"`
			Content   string    `json:"content"`
		}

		// Extract the comments from the rows with usernames
		comments := make([]CommentReturn, len(rows))
		for i, row := range rows {
			comments[i] = CommentReturn{
				ID:        row.Comment.ID,
				PostID:    row.Comment.PostID,
				UserID:    row.Comment.UserID,
				Username:  row.Username,
				CreatedAt: row.Comment.CreatedAt,
				Content:   row.Comment.Content,
			}
		}
		utilities.MarshallToResponse(r.Context(), logger, w, comments)
	})
}

// route: GET /posts
// return the most recent posts (paged by 10 at a time)
func handleMostRecentPosts(logger *slog.Logger, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "get most recent posts GET request received")

		var pageNum int32

		pageVal := r.URL.Query().Get("page")
		if pageVal == "" {
			pageNum = 0
		} else {
			var err error
			pageNumBig, err := strconv.Atoi(pageVal)
			if err != nil {
				logger.InfoContext(r.Context(), "failed to parse page from url query param", slog.Any("error", err))
				errorMessage := "Url parameter could not be parsed properly."
				http.Error(w, errorMessage, http.StatusBadRequest)
				return
			}
			if pageNumBig < 1 {
				pageNum = 0
			} else {
				pageNum = int32(pageNumBig - 1)
			}
		}

		var offset = pageNum * 10

		logger.InfoContext(r.Context(), "getting posts from db", slog.Int("offset", int(offset)))

		rows, err := db.GetPostsPaged(r.Context(), offset)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to find more posts from db", slog.Int("pageNum", int(pageNum)), slog.Any("error", err))
			w.Write([]byte("[]"))
			return
		} else if len(rows) == 0 || rows == nil {
			logger.InfoContext(r.Context(), "no posts found")
			w.Write([]byte("[]"))
			return
		}

		utilities.MarshallToResponse(r.Context(), logger, w, rows)
	})
}
