package models

import (
	"encoding/json"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Comment is a struct that represents a comment on a post.
type Comment struct {
	id        snowflake.Identifier // ID is a unique identifier for the comment.
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
		id:        id,
		Author:    author,
		Timestamp: timestamp,
		Content:   content,
		LikeCount: likeCount,
	}
}

// Id returns the ID of the comment.
func (c Comment) Id() uint64 {
	return c.id.Id()
}

func (c Comment) MarshalJSON() ([]byte, error) {
	commentAltered := struct {
		Id        uint64 `json:"id"`
		AuthorId  uint64 `json:"authorId"`
		Timestamp string `json:"timestamp"`
		Content   string `json:"content"`
		LikeCount int    `json:"likeCount"`
	}{
		Id:        c.Id(),
		AuthorId:  c.Author.Id(),
		Timestamp: c.Timestamp.Format(time.RFC3339),
		Content:   c.Content,
		LikeCount: c.LikeCount,
	}
	return json.Marshal(commentAltered)
}
