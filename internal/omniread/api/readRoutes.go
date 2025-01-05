package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/harrydayexe/Omni/internal/middleware"
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
	db *storage.Queries,
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
}

// extract the id parameter from the http request
// if the parameter cannot be parsed then an error is written to the http response
func extractIdParam(r *http.Request, w http.ResponseWriter, logger *slog.Logger) (snowflake.Snowflake, error) {
	idString := r.PathValue("id")
	logger.InfoContext(r.Context(), "extracting id from request", slog.String("idString", idString))
	idInt, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
		errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return snowflake.ParseId(0), fmt.Errorf("failed to parse id to int", err)
	}

	return snowflake.ParseId(idInt), nil
}

// check if an error is present and handle the http response if it is
func isDbError(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, id snowflake.Snowflake, err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		logger.ErrorContext(ctx, "entity not found", slog.Any("id", id))
		w.WriteHeader(http.StatusNotFound)
		return true
	}
	if err != nil {
		logger.ErrorContext(ctx, "failed to read entity from db", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}

	return false
}

// marshall an entity to a json object and write to http response
func marshallToResponse(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.ErrorContext(ctx, "failed to serialize entity to json", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// route: GET /post/{id}
// return the details of a user by it's id
func handleReadPost(logger *slog.Logger, db *storage.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read post GET request received")
		id, err := extractIdParam(r, w, logger)
		if err != nil {
			return
		}

		post, err := db.FindPostByID(r.Context(), int64(id.ToInt()))
		if isDbError(r.Context(), logger, w, id, err) {
			return
		}

		marshallToResponse(r.Context(), logger, w, post)
	})
}

// route: GET /user/{id}
// return the details of a user by it's id
func handleReadUser(logger *slog.Logger, db *storage.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read user GET request received")
		id, err := extractIdParam(r, w, logger)
		if err != nil {
			return
		}

		user, err := db.GetUserByID(r.Context(), int64(id.ToInt()))
		if isDbError(r.Context(), logger, w, id, err) {
			return
		}

		marshallToResponse(r.Context(), logger, w, user)
	})
}

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
			logger.ErrorContext(r.Context(), "failed to parse date from url query param", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
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
			logger.ErrorContext(r.Context(), "failed to parse comment limit from url query param", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
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
// from is the date and time to retrieve comments since in RFC3339 format
// limit is the number of comments to return (default to 10, max is 100)
func handleReadUserPosts(logger *slog.Logger, db *storage.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read posts GET request received")
		id, err := extractIdParam(r, w, logger)
		if err != nil {
			return
		}

		fromDate, limit, err := extractPaginationParams(logger, r, w)
		if err != nil {
			return
		}

		comments, err := db.GetUserAndPostsByIDPaged(r.Context(), storage.GetUserAndPostsByIDPagedParams{
			ID:           int64(id.ToInt()),
			CreatedAfter: fromDate,
			Limit:        int32(limit),
		})
		if isDbError(r.Context(), logger, w, id, err) {
			return
		}

		marshallToResponse(r.Context(), logger, w, comments)
	})
}

// route: GET /post/{id}/comments?from=2006-01-02T15%3A04%3A05Z07%3A00&limit=10
// there are two query parameters from and limit
// from is the date and time to retrieve comments since in RFC3339 format
// limit is the number of comments to return (default to 10, max is 100)
func handleReadPostComments(logger *slog.Logger, db *storage.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "read comments GET request received")
		id, err := extractIdParam(r, w, logger)
		if err != nil {
			return
		}

		fromDate, limit, err := extractPaginationParams(logger, r, w)
		if err != nil {
			return
		}

		comments, err := db.FindCommentsAndUserByPostIDPaged(r.Context(), storage.FindCommentsAndUserByPostIDPagedParams{
			PostID:       int64(id.ToInt()),
			CreatedAfter: fromDate,
			Limit:        int32(limit),
		})
		if isDbError(r.Context(), logger, w, id, err) {
			return
		}

		marshallToResponse(r.Context(), logger, w, comments)
	})
}
