package models

import (
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Comment is a struct that represents a comment on a post.
type Comment struct {
	ID        snowflake.Identifier // ID is a unique identifier for the comment.
	Author    User                 // Author is the user who wrote the comment.
	Timestamp time.Time            // Timestamp is the time the comment was created.
	Content   string               // Content is the content of the comment.
	LikeCount int                  // LikeCount is the number of likes the comment has.
}

// NewComment creates a new Comment with the given ID, author, timestamp, and content.
func NewComment(
	id snowflake.Identifier,
	author User,
	timestamp time.Time,
	content string,
	likeCount int,
) Comment {
	return Comment{
		ID:        id,
		Author:    author,
		Timestamp: timestamp,
		Content:   content,
		LikeCount: likeCount,
	}
}

// Id returns the ID of the comment.
func (c Comment) Id() int64 {
	return c.ID.Id()
}
