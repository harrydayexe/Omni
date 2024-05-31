package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// UserRepo is an implementation of the Repository interface focused on
// providing access to the User table in the database.
// db is the database connection that the repository will use to perform queries.
// logger is the logger that the repository will use to perform logging.
type UserRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewUserRepo creates a new UserRepo instance.
func NewUserRepo(db *sql.DB, logger *slog.Logger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
	r.logger.DebugContext(ctx, "Reading user from database", slog.Int64("id", int64(id.ToInt())))
	var username string
	r.logger.DebugContext(ctx, "Querying user information")
	row := r.db.QueryRowContext(ctx, `SELECT username FROM Users WHERE id = ?;`, id.ToInt())
	err := row.Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		r.logger.DebugContext(ctx, "Could not find user", slog.Int64("id", int64(id.ToInt())))
		return nil, nil
	case err != nil:
		r.logger.ErrorContext(ctx, "An unknown database error occurred when reading the user", slog.Any("error", err))
		return nil, NewDatabaseError("an unknown database error occurred when reading the user", err)
	}

	r.logger.DebugContext(ctx, "Querying post ids")
	rows, err := r.db.QueryContext(ctx, `SELECT id FROM Posts WHERE user_id = ?;`, id.ToInt())
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when querying user post ids", slog.Any("error", err))
		return nil, NewDatabaseError("an unknown database error occurred when reading user post ids", err)
	}
	defer rows.Close()

	postIDs := make([]snowflake.Snowflake, 0)
	for rows.Next() {
		r.logger.DebugContext(ctx, "Reading next row from results")
		var postID int64
		err := rows.Scan(&postID)
		// I don't think this can ever actually occur. As far as I can see online .Scan() only ever throws a ErrNoRows
		// And this can't occur here because rows.Next() would return false if there are no more rows
		if err != nil {
			r.logger.ErrorContext(ctx, "An unknown database error occurred when indexing through post ids", slog.Any("error", err))
			return nil, NewDatabaseError("an unknown database error occurred when reading user post ids", err)
		}
		postIDs = append(postIDs, snowflake.ParseId(uint64(postID)))
	}
	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when querying user post ids", slog.Any("error", err))
		return nil, NewDatabaseError("an unknown database error occurred when reading user post ids", err)
	}

	user := models.NewUser(id, username, postIDs)
	r.logger.DebugContext(ctx, "Successfully read user from database", slog.Any("user", user))
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	// TODO: Implement this method
	return nil
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	r.logger.DebugContext(ctx, "Updating user in database", slog.Any("user", user))

	result, err := r.db.ExecContext(ctx, "UPDATE Users SET username = ? WHERE id = ?", user.Username, user.Id().ToInt())
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when updating the user", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when updating the user", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "An unknown database error occurred when updating the user", slog.Any("error", err))
		return NewDatabaseError("an unknown database error occurred when updating the user", err)
	}
	if rows != 1 {
		r.logger.DebugContext(ctx, "Expected to affect 1 row", slog.Int64("affected", rows))
		return NewNotFoundError(user.Id(), nil)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	// TODO: Implement this method
	return nil
}
