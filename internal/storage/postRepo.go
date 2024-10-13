package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// PostRepo is an implementation of the Repository interface focused on
// providing access to the Post table in the database.
type PostRepo struct {
	db *sql.DB
}

// NewPostRepo creates a new PostRepo instance.
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
	// TODO: Implement this method
	return nil, nil
}

func (r *PostRepo) Create(ctx context.Context, user models.Post) error {
	// TODO: Implement this method
	return nil
}

func (r *PostRepo) Update(ctx context.Context, user models.Post) error {
	// TODO: Implement this method
	return nil
}

func (r *PostRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	// TODO: Implement this method
	return nil
}

func (r *PostRepo) GetPostsForUser(ctx context.Context, userId snowflake.Snowflake, from time.Time, limit int) ([]models.Post, error) {
	return nil, nil
}
