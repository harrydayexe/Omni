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

// AddRoutes adds all api routes to the provided http.ServeMux.
// It also adds logging middleware to each route.
func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	postUser storage.Repository[models.Post],
	userRepo storage.Repository[models.User],
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	mux.Handle("GET /post/{id}", stack(handleReadPost(logger, postUser)))
	mux.Handle("GET /user/{id}", stack(handleReadUser(logger, userRepo)))
}

func handleReadPost(logger *slog.Logger, postRepo storage.Repository[models.Post]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.Error("failed to parse id: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := snowflake.ParseId(idInt)

		post, err := postRepo.Read(id)
		if err != nil {
			logger.Error("failed to read post: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(post)
		if err != nil {
			logger.Error("failed to serialize post to json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	})
}

func handleReadUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		idInt, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			logger.Error("failed to parse id: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := snowflake.ParseId(idInt)

		user, err := userRepo.Read(id)
		if err != nil {
			logger.Error("failed to read user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(user)
		if err != nil {
			logger.Error("failed to serialize user to json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	})
}
