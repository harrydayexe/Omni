package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func AddAuthRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	authService auth.Authable,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
	)

	// Get the details of a post by id
	mux.Handle("POST /login", stack(handleLogin(logger, authService)))
}

func handleLogin(logger *slog.Logger, authService auth.Authable) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "login POST request received")

		var u auth.LoginRequest
		err := utilities.DecodeJsonBody(r.Context(), logger, w, r, &u)
		if err != nil {
			return
		}

		token, err := authService.Login(r.Context(), u.Username, u.Password)
		if errors.Is(err, auth.ErrUserNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if errors.Is(err, auth.ErrUnauthorized) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		expiresIn := time.Hour * 24

		utilities.MarshallToResponse(r.Context(), logger, w,
			auth.LoginResponse{
				Token:   token,
				Type:    "Bearer",
				Expires: int(expiresIn.Seconds()),
			},
		)
	})
}
