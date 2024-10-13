package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

const (
	epoch int64 = 1704067200000
)

// AddReadRoutes adds all api routes to the provided http.ServeMux.
// It also adds logging middleware to each route.
func AddReadRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	userRepo storage.Repository[models.User],
	postRepo storage.Repository[models.Post],
	commentRepo storage.CommentRepository,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("GET /post/{id}", stack(handleReadPost(logger, postRepo)))
	mux.Handle("GET /user/{id}", stack(handleReadUser(logger, userRepo)))
	mux.Handle("GET /user/{id}/posts", stack(handleReadUserPosts(logger, userRepo, postRepo)))
	mux.Handle("GET /post/{id}/comments", stack(handleReadPostComments(logger, commentRepo)))
}

// route: GET /post/{id}
// return the details of a user by it's id
func handleReadPost(logger *slog.Logger, postRepo storage.Repository[models.Post]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		logger.InfoContext(r.Context(), "read post GET request received", slog.String("id", idString))
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		id := snowflake.ParseId(idInt)

		post, err := postRepo.Read(r.Context(), id)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to read post from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post == nil {
			logger.DebugContext(r.Context(), "post not found", slog.Any("id", id))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := json.Marshal(post)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to serialize post to json", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})
}

// route: GET /user/{id}
// return the details of a user by it's id
func handleReadUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		logger.InfoContext(r.Context(), "read user GET request received", slog.String("id", idString))
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		id := snowflake.ParseId(idInt)

		user, err := userRepo.Read(r.Context(), id)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to read user from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if user == nil {
			logger.DebugContext(r.Context(), "user not found", slog.Any("id", id))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := json.Marshal(user)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to serialize user to json", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})
}

// route: GET /user/{id}/posts
// optional query parameter ?from= which represents the post number to return from
// return the posts of a user by it's id. limit is 50 posts
func handleReadUserPosts(logger *slog.Logger, userRepo storage.Repository[models.User], postRepo storage.Repository[models.Post]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// route: GET /post/{id}/comments?from=2006-01-02T15%3A04%3A05Z07%3A00&limit=10
// there are two query parameters from and limit
// from is the date and time to retrieve comments since in RFC3339 format
// limit is the number of comments to return (default to 10, max is 100)
func handleReadPostComments(logger *slog.Logger, commentRepo storage.CommentRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		logger.InfoContext(r.Context(), "read comments GET request received", slog.String("postId", idString))

		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		id := snowflake.ParseId(idInt)

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
				logger.ErrorContext(r.Context(), "failed to parse date from url query param", slog.Any("error", err))
				errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errorMessage))
				return
			}
		}
		limitVal := r.URL.Query().Get("limit")
		if limitVal == "" {
			limit = 10
		} else {
			var err error
			limit, err = strconv.Atoi(limitVal)
			if err != nil {
				logger.ErrorContext(r.Context(), "failed to parse comment limit from url query param", slog.Any("error", err))
				errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errorMessage))
				return
			}
			if limit > 100 {
				limit = 100
			}
		}

		comments, err := commentRepo.GetCommentsForPost(r.Context(), id, fromDate, limit)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to get comments from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(comments)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to serialize comments to json", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	})
}
