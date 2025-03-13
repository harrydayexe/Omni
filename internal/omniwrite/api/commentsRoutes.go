package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniwrite/datamodels"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func AddCommentsRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	db storage.Querier,
	snowflakeGenerator *snowflake.SnowflakeGenerator,
	authService auth.Authable,
	config *config.Config,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewSetContentTypeJson(),
		middleware.NewMaxBytesReader(),
	)

	mux.Handle("POST /post/{id}/comment", stack(handleInsertComment(logger, db, snowflakeGenerator, authService, config)))
	mux.Handle("PUT /comment/{id}", stack(handleUpdateComment(logger, db, authService, config)))
	mux.Handle("DELETE /comment/{id}", stack(handleDeleteComment(logger, authService, db)))
}

// route: POST /post/{id}/comment
// insert a new comment into the database
func handleInsertComment(logger *slog.Logger, db storage.Querier, gen *snowflake.SnowflakeGenerator, authService auth.Authable, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "insert comment POST request received")

		post_id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		var c datamodels.NewComment
		err = utilities.DecodeJsonBody(r.Context(), logger, w, r, &c)
		if err != nil {
			return
		}

		err = utilities.CheckBearerAuth(snowflake.ParseId(uint64(c.UserID)), authService, logger, w, r)
		if err != nil {
			return
		}

		newComment := storage.Comment{
			ID:        int64(gen.NextID().ToInt()),
			PostID:    int64(post_id.ToInt()),
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		}

		err = db.CreateComment(r.Context(), storage.CreateCommentParams(newComment))
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to insert comment", slog.Any("error", err))
			http.Error(w, "failed to create comment", http.StatusInternalServerError)
			return
		}

		strPostId := strconv.Itoa(int(newComment.PostID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/api/post/"+strPostId+"/comments")
		w.WriteHeader(http.StatusCreated)
		utilities.MarshallToResponse(r.Context(), logger, w, newComment)
	})
}

// route: PUT /comment/{id}
// update a comment by id
func handleUpdateComment(logger *slog.Logger, db storage.Querier, authService auth.Authable, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "update comment PUT request received")

		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		var c struct {
			Content string `json:"content"`
		}
		err = utilities.DecodeJsonBody(r.Context(), logger, w, r, &c)
		if err != nil {
			return
		}

		// Check comment exists
		currentCommentAndUser, err := db.FindCommentAndUserByID(r.Context(), int64(id.ToInt()))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.InfoContext(r.Context(), "entity not found", slog.Any("id", id))
				http.Error(w, "entity not found", http.StatusNotFound)
				return
			}
			logger.ErrorContext(r.Context(), "failed to read entity from db", slog.Any("error", err))
			http.Error(w, "failed to read entity from db", http.StatusInternalServerError)
			return
		}

		err = utilities.CheckBearerAuth(snowflake.ParseId(uint64(currentCommentAndUser.Comment.UserID)), authService, logger, w, r)
		if err != nil {
			return
		}

		updatedComment := storage.Comment{
			ID:        currentCommentAndUser.Comment.ID,
			PostID:    currentCommentAndUser.Comment.PostID,
			UserID:    currentCommentAndUser.Comment.UserID,
			Content:   c.Content,
			CreatedAt: currentCommentAndUser.Comment.CreatedAt,
		}

		err = db.UpdateComment(r.Context(), storage.UpdateCommentParams{
			ID:      updatedComment.ID,
			Content: updatedComment.Content,
		})
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to update comment", slog.Any("error", err))
			http.Error(w, "failed to update comment", http.StatusInternalServerError)
			return
		}

		strPostId := strconv.Itoa(int(updatedComment.PostID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/api/post/"+strPostId+"/comments")
		w.WriteHeader(http.StatusOK)
		utilities.MarshallToResponse(r.Context(), logger, w, updatedComment)
	})
}

// route: DELETE /comment/{id}
// delete a comment by id
func handleDeleteComment(logger *slog.Logger, authService auth.Authable, db storage.Querier) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "delete comment DELETE request received")

		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		// Check post exists
		currentCommentAndUser, err := db.FindCommentAndUserByID(r.Context(), int64(id.ToInt()))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.InfoContext(r.Context(), "entity not found", slog.Any("id", id))
				http.Error(w, "entity not found", http.StatusNotFound)
				return
			}
			logger.ErrorContext(r.Context(), "failed to read entity from db", slog.Any("error", err))
			http.Error(w, "failed to read entity from db", http.StatusInternalServerError)
			return
		}

		err = utilities.CheckBearerAuth(snowflake.ParseId(uint64(currentCommentAndUser.Comment.UserID)), authService, logger, w, r)
		if err != nil {
			return
		}

		err = db.DeleteComment(r.Context(), int64(id.ToInt()))
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to delete comment", slog.Any("error", err))
			http.Error(w, "failed to delete comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
