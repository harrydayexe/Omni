package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

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
	r.logger.DebugContext(ctx, "Reading comment from database", slog.Int64("id", int64(id.ToInt())))
	var postId int64
	var userId int64
	var username string
	var content string
	var createdAt time.Time
	r.logger.DebugContext(ctx, "Querying comment information")
	row := r.db.QueryRowContext(ctx, `SELECT Comments.post_id, Comments.user_id, Users.username, Comments.content, Comments.created_at FROM Comments INNER JOIN Users ON Comments.user_id = Users.id WHERE Comments.id = ?;`, id.ToInt())
	err := row.Scan(&postId, &userId, &username, &content, &createdAt)
	switch {
	case err == sql.ErrNoRows:
		r.logger.DebugContext(ctx, "Could not find comment", slog.Int64("id", int64(id.ToInt())))
		return nil, nil
	case err != nil:
		r.logger.ErrorContext(ctx, "An unknown database error occurred when reading the comment", slog.Any("error", err))
		return nil, NewDatabaseError("an unknown database error occurred when reading the comment", err)
	}

	user := models.NewComment(id, snowflake.ParseId(uint64(postId)), snowflake.ParseId(uint64(userId)), username, createdAt, content)
	r.logger.DebugContext(ctx, "Successfully read user from database", slog.Any("user", user))
	return &user, nil
}

func (r *CommentRepo) Create(ctx context.Context, comment models.Comment) error {
	r.logger.DebugContext(ctx, "Creating comment in database", slog.Any("comment", comment))

	result, err := r.db.ExecContext(ctx, "INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?);", comment.Id().ToInt(), comment.PostId.ToInt(), comment.AuthorId.ToInt(), comment.Content, comment.Timestamp)
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when creating the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when creating the comment", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when creating the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when creating the comment", err)
	}
	if rows != 1 {
		r.logger.DebugContext(ctx, "Expected to affect 1 row", slog.Int64("affected", rows))
		return NewEntityAlreadyExistsError(comment.Id(), nil)
	}
	return nil
}

func (r *CommentRepo) Update(ctx context.Context, comment models.Comment) error {
	r.logger.DebugContext(ctx, "Updating comment in database", slog.Any("comment", comment))

	result, err := r.db.ExecContext(ctx, "UPDATE Comments SET post_id=?, user_id=?, content=?, created_at=? WHERE id=?;", comment.PostId.ToInt(), comment.AuthorId.ToInt(), comment.Content, comment.Timestamp, comment.Id().ToInt())
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when updating the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when updating the comment", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when updating the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when updating the comment", err)
	}
	if rows != 1 {
		r.logger.DebugContext(ctx, "Expected to affect 1 row", slog.Int64("affected", rows))
		return NewNotFoundError(comment.Id(), nil)
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	r.logger.DebugContext(ctx, "Deleting comment from database", slog.Int("id", int(id.ToInt())))

	result, err := r.db.ExecContext(ctx, "DELETE FROM Comments WHERE id = ?", id.ToInt())
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when deleting the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when deleting the comment", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when deleting the comment", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when deleting the comment", err)
	}
	if rows != 1 {
		r.logger.DebugContext(ctx, "Expected to affect 1 row", slog.Int64("affected", rows))
		return NewNotFoundError(id, nil)
	}
	return nil
}
