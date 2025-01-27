package api

import (
	"log/slog"
	"net/http"

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
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	mux.Handle("POST /user", stack(handleInsertUser(logger, db, snowflakeGenerator)))
}

// route: POST /post/{id}
// return the details of a user by it's id
func handleInsertUser(logger *slog.Logger, db storage.Querier, gen *snowflake.SnowflakeGenerator) http.Handler {
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
	})
}
