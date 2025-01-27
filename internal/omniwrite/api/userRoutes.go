package api

import (
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
}

// route: POST /post/{id}
// return the details of a user by it's id
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
