package api

import (
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/storage"
)

func NewHandler(
	logger *slog.Logger,
	userRepo storage.Repository[models.User],
	postRepo storage.Repository[models.Post],
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
		mux,
		logger,
		postRepo,
		userRepo,
	)
	var handler http.Handler = mux
	return handler
}
