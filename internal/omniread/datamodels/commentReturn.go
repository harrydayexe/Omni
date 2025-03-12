package datamodels

import "time"

type CommentReturn struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type CommentsForPostReturn struct {
	TotalPages int
	Comments   []CommentReturn
}
