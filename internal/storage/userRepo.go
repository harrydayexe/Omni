package storage

import (
	"context"
	"database/sql"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// UserRepo is an implementation of the Repository interface focused on
// providing access to the User table in the database.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new UserRepo instance.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
	var username string
	row := r.db.QueryRowContext(ctx, `SELECT username FROM Users WHERE id = ?;`, id.ToInt())
	err := row.Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		return nil, NewCouldNotFindEntityError("User", id)
	case err != nil:
		return nil, NewDatabaseError("an unknown database error occurred when reading the user", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id FROM Posts WHERE user_id = ?;`, id.ToInt())
	defer rows.Close()

	postIDs := make([]snowflake.Snowflake, 0)
	for rows.Next() {
		var postID int64
		err := rows.Scan(&postID)
		switch {
		case err == sql.ErrNoRows:
			return nil, NewCouldNotFindEntityError("User", id)
		case err != nil:
			return nil, NewDatabaseError("an unknown database error occurred when reading user post ids", err)
		}
		postIDs = append(postIDs, snowflake.ParseId(uint64(postID)))
	}
	if err := rows.Err(); err != nil {
		return nil, NewDatabaseError("an unknown database error occurred when reading user post ids", err)
	}

	user := models.NewUser(id, username, postIDs)
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	// TODO: Implement this method
	return nil
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	// TODO: Implement this method
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	// TODO: Implement this method
	return nil
}
