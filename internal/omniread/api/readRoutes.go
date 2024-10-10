package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

// AddReadRoutes adds all api routes to the provided http.ServeMux.
// It also adds logging middleware to each route.
func AddReadRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	userRepo storage.Repository[models.User],
	postRepo storage.Repository[models.Post],
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("GET /post/{id}", stack(handleReadPost(logger, postRepo)))
	mux.Handle("GET /user/{id}", stack(handleReadUser(logger, userRepo)))
	mux.Handle("GET /user/{id}/posts", stack(handleReadUserPosts(logger, userRepo, postRepo)))
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
