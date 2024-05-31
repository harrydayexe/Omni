package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// CommentRepo is an implementation of the Repository interface focused on
// providing access to the Comment table in the database.
type CommentRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewCommentRepo creates a new CommentRepo instance.
func NewCommentRepo(db *sql.DB, logger *slog.Logger) *CommentRepo {
	return &CommentRepo{
		db:     db,
		logger: logger,
	}
}

func (r *CommentRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.Comment, error) {
	// TODO: Implement this method
	return nil, nil
}

func (r *CommentRepo) Create(ctx context.Context, Comment models.Comment) error {
	// TODO: Implement this method
	return nil
}

func (r *CommentRepo) Update(ctx context.Context, Comment models.Comment) error {
	// TODO: Implement this method
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	// TODO: Implement this method
	return nil
}
