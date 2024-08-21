package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

// AddRoutes adds all api routes to the provided http.ServeMux.
// It also adds logging middleware to each route.
func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	postRepo storage.Repository[models.Post],
	userRepo storage.Repository[models.User],
	commentRepo storage.Repository[models.Comment],
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("GET /user/{id}", stack(handleReadUser(logger, userRepo)))
	mux.Handle("POST /user", stack(handleCreateUser(logger, userRepo)))
	mux.Handle("PUT /user/{id}", stack(handleUpdateUser(logger, userRepo)))
	mux.Handle("DELETE /user/{id}", stack(handleDeleteUser(logger, userRepo)))
	// mux.Handle("GET /post/{id}", stack(handleReadPost(logger, postRepo)))
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

func handleCreateUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func handleUpdateUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func handleDeleteUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		logger.InfoContext(r.Context(), "read user DELETE request received", slog.String("id", idString))
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
			errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		id := snowflake.ParseId(idInt)

		err = userRepo.Delete(r.Context(), id)
		var e *storage.NotFoundError
		if errors.As(err, &e) {
			logger.DebugContext(r.Context(), "user not found", slog.Any("id", id))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to delete user from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

// route: GET /post{id}
// return the details of a post by it's id
// func handleReadPost(logger *slog.Logger, postRepo storage.Repository[models.Post]) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		idString := r.PathValue("id")
// 		logger.InfoContext(r.Context(), "read post GET request received", slog.String("id", idString))
// 		idInt, err := strconv.ParseUint(idString, 10, 64)
// 		if err != nil {
// 			logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
//
// 		id := snowflake.ParseId(idInt)
//
// 		post, err := postRepo.Read(r.Context(), id)
// 		if err != nil {
// 			logger.ErrorContext(r.Context(), "failed to read post", slog.Any("error", err))
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
//
// 		if post == nil {
// 			logger.DebugContext(r.Context(), "post not found in db", slog.Any("id", id))
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}
//
// 		b, err := json.Marshal(post)
// 		if err != nil {
// 			logger.ErrorContext(r.Context(), "failed to serialize post to json", slog.Any("error", err))
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Write(b)
// 	})
// }
