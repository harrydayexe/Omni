package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

// Notes:
// POST for insert, PUT for update
// post is not idempotent (multiple requests = multiple new users)
// POST should return Location header with the URL of the new resource
// POST should return 201 for creation
// PUT should return 200/204 for success

func AddUserRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db storage.Querier,
	snowflakeGenerator *snowflake.SnowflakeGenerator,
	config *config.Config,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	mux.Handle("POST /user", stack(handleInsertUser(logger, db, snowflakeGenerator, config)))
	mux.Handle("PUT /user/{id}", stack(handleUpdateUser(logger, db, config)))
}

// route: POST /user/
// insert a new user into the database
func handleInsertUser(logger *slog.Logger, db storage.Querier, gen *snowflake.SnowflakeGenerator, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "insert user POST request received")

		var u struct {
			Username string `json:"username"`
		}
		err := utilities.DecodeJsonBody(r.Context(), logger, w, r, &u)
		if err != nil {
			return
		}

		newUser := storage.User{
			Username: u.Username,
			ID:       int64(gen.NextID().ToInt()),
		}

		err = db.CreateUser(r.Context(), storage.CreateUserParams(newUser))
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to insert user", slog.Any("error", err))
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		strId := strconv.Itoa(int(newUser.ID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/user/"+strId)
		w.WriteHeader(http.StatusCreated)
	})
}

// route: PUT /user/{id}
// update a user by id
func handleUpdateUser(logger *slog.Logger, db storage.Querier, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "update user PUT request received")

		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		// Check user exists
		_, err = db.GetUserByID(r.Context(), int64(id.ToInt()))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.ErrorContext(r.Context(), "entity not found", slog.Any("id", id))
				http.Error(w, "entity not found", http.StatusNotFound)
				return
			}
			logger.ErrorContext(r.Context(), "failed to read entity from db", slog.Any("error", err))
			http.Error(w, "failed to read entity from db", http.StatusInternalServerError)
			return
		}

		var u struct {
			Username string `json:"username"`
		}
		err = utilities.DecodeJsonBody(r.Context(), logger, w, r, &u)
		if err != nil {
			return
		}

		updatedUser := storage.User{
			Username: u.Username,
			ID:       int64(id.ToInt()),
		}

		err = db.UpdateUser(r.Context(), storage.UpdateUserParams{
			ID:       updatedUser.ID,
			Username: updatedUser.Username,
		})
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to update user", slog.Any("error", err))
			http.Error(w, "failed to update user", http.StatusInternalServerError)
			return
		}

		strId := strconv.Itoa(int(updatedUser.ID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/user/"+strId)
		utilities.MarshallToResponse(r.Context(), logger, w, updatedUser)
	})
}
