package datamodels

import (
	"html/template"
	"time"

	"github.com/harrydayexe/Omni/internal/omniread/datamodels"
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
	Comments    CommentsModel
}

// CommentModel is the data model for the "comments" partial template
type CommentsModel struct {
	Error    string
	Comments []datamodels.CommentReturn
}

// NewCommentsModel creates the data model from a list of comments and an error
func NewCommentsModel(err error, comments []datamodels.CommentReturn) CommentsModel {
	if err != nil {
		return CommentsModel{
			Error: err.Error(),
		}
	} else {
		return CommentsModel{
			Comments: comments,
		}
	}
}

// NewPost creates the data model from a post, user and content
func NewPost(post storage.Post, user storage.User, content string, comments CommentsModel) Post {
	return Post{
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt.Format(time.DateTime),
		Author:      user.Username,
		AuthorID:    user.ID,
		Content:     template.HTML(content),
		Comments:    comments,
	}
}
