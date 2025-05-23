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

func AddPostRoutes(
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

	mux.Handle("POST /post", stack(handleInsertPost(logger, db, snowflakeGenerator, authService, config)))
	mux.Handle("PUT /post/{id}", stack(handleUpdatePost(logger, db, authService, config)))
	mux.Handle("DELETE /post/{id}", stack(handleDeletePost(logger, db, authService)))
}

// route: POST /post/
// insert a new post into the database
func handleInsertPost(logger *slog.Logger, db storage.Querier, gen *snowflake.SnowflakeGenerator, authService auth.Authable, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "insert post POST request received")

		var p datamodels.NewPost
		err := utilities.DecodeJsonBody(r.Context(), logger, w, r, &p)
		if err != nil {
			return
		}
		logger.DebugContext(r.Context(), "decoded json body", slog.Any("body", p))

		err = utilities.CheckBearerAuth(snowflake.ParseId(p.UserID), authService, logger, w, r)
		if err != nil {
			return
		}

		newPost := storage.Post{
			ID:          int64(gen.NextID().ToInt()),
			UserID:      int64(p.UserID),
			CreatedAt:   p.CreatedAt,
			Title:       p.Title,
			Description: p.Description,
			MarkdownUrl: p.MarkdownUrl,
		}

		err = db.CreatePost(r.Context(), storage.CreatePostParams{
			ID:          newPost.ID,
			UserID:      newPost.UserID,
			CreatedAt:   newPost.CreatedAt,
			Title:       newPost.Title,
			Description: newPost.Description,
			MarkdownUrl: newPost.MarkdownUrl,
		})
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to insert post", slog.Any("error", err))
			http.Error(w, "failed to create post", http.StatusInternalServerError)
			return
		}
		logger.DebugContext(r.Context(), "inserted post into db", slog.Any("post", newPost))

		strId := strconv.Itoa(int(newPost.ID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/api/post/"+strId)
		w.WriteHeader(http.StatusCreated)
		utilities.MarshallToResponse(r.Context(), logger, w, newPost)
	})
}

// route: PUT /post/{id}
// update a post by id
func handleUpdatePost(logger *slog.Logger, db storage.Querier, authService auth.Authable, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "update post PUT request received")

		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		var p datamodels.UpdatedPost
		err = utilities.DecodeJsonBody(r.Context(), logger, w, r, &p)
		if err != nil {
			return
		}

		// Check post exists
		currentPost, err := db.FindPostByID(r.Context(), int64(id.ToInt()))
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

		// Check user is authorized to update post
		err = utilities.CheckBearerAuth(snowflake.ParseId(uint64(currentPost.UserID)), authService, logger, w, r)
		if err != nil {
			return
		}

		updatedPost := storage.Post{
			ID:          int64(id.ToInt()),
			UserID:      currentPost.UserID,
			CreatedAt:   currentPost.CreatedAt,
			Title:       p.Title,
			Description: p.Description,
			MarkdownUrl: p.MarkdownUrl,
		}

		err = db.UpdatePost(r.Context(), storage.UpdatePostParams{
			ID:          updatedPost.ID,
			Title:       updatedPost.Title,
			Description: updatedPost.Description,
			MarkdownUrl: updatedPost.MarkdownUrl,
		})
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to update post", slog.Any("error", err))
			http.Error(w, "failed to update post", http.StatusInternalServerError)
			return
		}

		strId := strconv.Itoa(int(updatedPost.ID))
		strPort := strconv.Itoa(config.Port)
		w.Header().Set("Location", config.Host+":"+strPort+"/api/post/"+strId)
		utilities.MarshallToResponse(r.Context(), logger, w, updatedPost)
	})
}

// route: DELETE /post/{id}
// delete a post by id
func handleDeletePost(logger *slog.Logger, db storage.Querier, authService auth.Authable) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "delete post DELETE request received")

		id, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		// Check post exists
		post, err := db.FindPostByID(r.Context(), int64(id.ToInt()))
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

		// Check user is authorized to delete post
		err = utilities.CheckBearerAuth(snowflake.ParseId(uint64(post.UserID)), authService, logger, w, r)
		if err != nil {
			return
		}

		err = db.DeletePost(r.Context(), int64(id.ToInt()))
		if err != nil {
			logger.ErrorContext(r.Context(), "failed to delete post", slog.Any("error", err))
			http.Error(w, "failed to delete post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
