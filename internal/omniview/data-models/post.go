package datamodels

import (
	"html/template"
	"time"

	"github.com/harrydayexe/Omni/internal/omniread/datamodels"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

type Comment struct {
	datamodels.CommentReturn
	IsDeleteable bool
	IsEditable   bool
}

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
	Error          string
	Comments       []Comment
	PostID         int64
	NextPageNumber int
	IsLoggedIn     bool
}

// NewCommentsModel creates the data model from a list of comments and an error
func NewCommentsModel(
	err error,
	userID snowflake.Identifier,
	comments datamodels.CommentsForPostReturn,
	postID int64,
	nextPage int,
) CommentsModel {
	var np int = 0
	if nextPage <= comments.TotalPages {
		np = nextPage
	}
	var commentsList []Comment = utilities.Map(
		comments.Comments,
		convertCommentReturnToComment(userID),
	)
	var isLoggedIn bool = false
	if userID != nil {
		isLoggedIn = true
	}
	if err != nil {
		return CommentsModel{
			Error:          "An error occurred while retrieving the comments",
			PostID:         postID,
			NextPageNumber: np,
			IsLoggedIn:     isLoggedIn,
		}
	} else {
		return CommentsModel{
			Comments:       commentsList,
			PostID:         postID,
			NextPageNumber: np,
			IsLoggedIn:     isLoggedIn,
		}
	}
}

// convertCommentReturnToComment returns a function which maps a CommentReturn to a Comment
func convertCommentReturnToComment(
	userID snowflake.Identifier,
) func(datamodels.CommentReturn) Comment {
	if userID == nil {
		return func(cr datamodels.CommentReturn) Comment {
			return Comment{
				CommentReturn: cr,
				IsDeleteable:  false,
				IsEditable:    false,
			}
		}
	}
	return func(cr datamodels.CommentReturn) Comment {
		return Comment{
			CommentReturn: cr,
			IsDeleteable:  cr.UserID == int64(userID.Id().ToInt()),
			IsEditable:    cr.UserID == int64(userID.Id().ToInt()),
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
