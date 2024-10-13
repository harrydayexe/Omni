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
	commentRepo storage.CommentRepository,
) http.Handler {
	mux := http.NewServeMux()
	AddReadRoutes(
		mux,
		logger,
		userRepo,
		postRepo,
		commentRepo,
	)
	var handler http.Handler = mux
	return handler
}
