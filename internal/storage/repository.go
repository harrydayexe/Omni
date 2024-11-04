package storage

import (
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

type Repository interface {
	GetUserById(id snowflake.Snowflake) (*User, error)
	GetPostById(id snowflake.Snowflake) (*Post, error)
	GetPostsOfUser(id snowflake.Snowflake, createdAfter time.Time, limit int) ([]*Post, error)
	GetCommentsOfPost(id snowflake.Snowflake, createdAfter time.Time, limit int) ([]*Comment, error)
}

// TODO: Implement the Repository interface
// TODO: Create custom errors for the repository
