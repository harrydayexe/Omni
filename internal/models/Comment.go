package models

import (
	"encoding/json"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Comment is a struct that represents a comment on a post.
type Comment struct {
	id         snowflake.Snowflake // ID is a unique identifier for the comment.
	PostId     snowflake.Snowflake // PostId is the id of the post the comment is on.
	AuthorId   snowflake.Snowflake // AuthorId is the id of the user who wrote the comment.
	AuthorName string              // AuthorName is the name of the user who wrote the comment.
	Timestamp  time.Time           // Timestamp is the time the comment was created.
	Content    string              // Content is the content of the comment.
}

// NewComment creates a new Comment with the given ID, author, timestamp, and content.
func NewComment(
	id snowflake.Snowflake,
	postId snowflake.Snowflake,
	authorId snowflake.Snowflake,
	authorName string,
	timestamp time.Time,
	content string,
) Comment {
	return Comment{
		id:         id,
		PostId:     postId,
		AuthorId:   authorId,
		AuthorName: authorName,
		Timestamp:  timestamp.UTC(),
		Content:    content,
	}
}

// Id returns the ID of the comment.
func (c Comment) Id() snowflake.Snowflake {
	return c.id
}

func (c Comment) MarshalJSON() ([]byte, error) {
	commentAltered := struct {
		Id         uint64 `json:"id"`
		PostId     uint64 `json:"postId"`
		AuthorId   uint64 `json:"authorId"`
		AuthorName string `json:"authorName"`
		Timestamp  string `json:"timestamp"`
		Content    string `json:"content"`
	}{
		Id:         c.Id().ToInt(),
		PostId:     c.PostId.ToInt(),
		AuthorId:   c.AuthorId.ToInt(),
		AuthorName: c.AuthorName,
		Timestamp:  c.Timestamp.Format(time.RFC3339),
		Content:    c.Content,
	}
	return json.Marshal(commentAltered)
}
