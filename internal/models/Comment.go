package models

import (
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

type Comment struct {
	ID        snowflake.Identifier // ID is a unique identifier for the comment.
	Author    User                 // Author is the user who wrote the comment.
	Timestamp time.Time            // Timestamp is the time the comment was created.
	Content   string               // Content is the content of the comment.
	LikeCount int                  // LikeCount is the number of likes the comment has.
}
