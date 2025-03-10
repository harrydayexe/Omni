package datamodels

import "time"

type NewPost struct {
	UserID      uint64    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MarkdownUrl string    `json:"markdown_url"`
}
