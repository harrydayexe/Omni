// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MarkdownUrl string    `json:"markdown_url"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
