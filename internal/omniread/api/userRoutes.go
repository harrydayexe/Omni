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
func AddUserRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	userRepo storage.Repository[models.User],
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
		logger.InfoContext(r.Context(), "create user POST request received")

		var u struct {
			Id       uint64 `json:"id"`
			Username string `json:"username"`
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&u)
		if err != nil {
			var errorMessage = `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		newUser := models.NewUser(snowflake.ParseId(u.Id), u.Username, []snowflake.Snowflake{})

		err = userRepo.Create(r.Context(), newUser)
		var e *storage.EntityAlreadyExistsError
		if errors.As(err, &e) {
			logger.DebugContext(r.Context(), "user already exists", slog.Any("id", newUser.Id()))
			var errorMessage = `{"error":"Conflict","message":"User with that ID already exists."}`
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(errorMessage))
			return
		}
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to create user in db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		b, err := json.Marshal(newUser)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to serialize user to json", slog.Any("error", err))
			return
		}

		w.Write(b)
	})
}

func handleUpdateUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "update user PUT request received")

		var u struct {
			Id       uint64 `json:"id"`
			Username string `json:"username"`
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&u)
		if err != nil {
			var errorMessage = `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		newUser := models.NewUser(snowflake.ParseId(u.Id), u.Username, []snowflake.Snowflake{})

		err = userRepo.Update(r.Context(), newUser)
		var e *storage.NotFoundError
		if errors.As(err, &e) {
			logger.DebugContext(r.Context(), "user does not exist", slog.Any("id", newUser.Id()))
			var errorMessage = `{"error":"Not Found","message":"User with that ID could not be found to update."}`
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to update user in db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		b, err := json.Marshal(newUser)
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to serialize user to json", slog.Any("error", err))
			return
		}

		w.Write(b)
	})
}

func handleDeleteUser(logger *slog.Logger, userRepo storage.Repository[models.User]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		logger.InfoContext(r.Context(), "delete user DELETE request received", slog.String("id", idString))
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
