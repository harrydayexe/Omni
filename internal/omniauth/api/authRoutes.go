package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/snowflake"
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

		var u struct {
			ID       uint64 `json:"id"`
			Password string `json:"password"`
		}
		err := utilities.DecodeJsonBody(r.Context(), logger, w, r, &u)
		if err != nil {
			return
		}

		token, err := authService.Login(r.Context(), snowflake.ParseId(u.ID), u.Password)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		expiresIn := time.Hour * 24

		utilities.MarshallToResponse(r.Context(), logger, w,
			struct {
				Token   string `json:"access_token"`
				Type    string `json:"token_type"`
				Expires int    `json:"expires_in"`
			}{
				Token:   token,
				Type:    "Bearer",
				Expires: int(expiresIn.Seconds()),
			},
		)
	})
}
