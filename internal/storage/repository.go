// Package storage provides interfaces and implementations for data storage
// and access.
package storage

import (
	"context"
	"time"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Repository is an interface for accessing persisted data.
// It provides methods for standard CRUD operations.
//
// Repository has an associated type T which must conform to Identifier.
type Repository[T snowflake.Identifier] interface {
	// Read retrieves an entity from the database.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Read(ctx context.Context, id snowflake.Snowflake) (*T, error)

	// Create adds entity to the database.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Create(ctx context.Context, entity T) error

	// Update modifies an entry to the database.
	// It matches based on the Id() of entity.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Update(ctx context.Context, entity T) error

	// Delete removes an entry from the database.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Delete(ctx context.Context, id snowflake.Snowflake) error
}

// CommentRepository is an interface for accessing persisted comment data.
// It provides methods for standard CRUD operations via the Repository interface.
// It also provides additional methods for querying comments.
type CommentRepository interface {
	// Repository is embedded to provide standard CRUD operations.
	Repository[models.Comment]

	// GetCommentsForPost retrieves comments for a particular post.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	GetCommentsForPost(ctx context.Context, postId snowflake.Snowflake, from time.Time, limit int) ([]models.Comment, error)
}

// PostRepository is an interface for accessing persisted post data.
// It provides methods for standard CRUD operations via the Repository interface.
// It also provides additional methods for querying posts.
type PostRepository interface {
	// Repository is embedded to provide standard CRUD operations.
	Repository[models.Post]

	// GetPostsForUser retrieves posts for a particular user.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	GetPostsForUser(ctx context.Context, userId snowflake.Snowflake, from time.Time, limit int) ([]models.Post, error)
}
