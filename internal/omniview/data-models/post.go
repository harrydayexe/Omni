package datamodels

import (
	"html/template"
	"time"

	"github.com/harrydayexe/Omni/internal/storage"
)

// Post is the data model for the "post" partial template
type Post struct {
	Title       string
	Description string
	CreatedAt   string
	Author      string
	AuthorID    int64
	Content     template.HTML
}

// NewPost creates the data model from a post, user and content
func NewPost(post storage.Post, user storage.User, content string) Post {
	return Post{
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt.Format(time.DateTime),
		Author:      user.Username,
		AuthorID:    user.ID,
		Content:     template.HTML(content),
	}
}
