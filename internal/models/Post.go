package models

import (
	"net/url"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Post represents a blog post.
type Post struct {
	id          snowflake.Identifier // ID is a unique identifier for the post.
	Author      User                 // Author is the user who wrote the post.
	Timestamp   time.Time            // Timestamp is the time the post was created.
	Title       string               // Title is the title of the post.
	Description string               // Description is a short description of the post.
	ContentFile url.URL              // ContentFile is the URL of the markdown file containing the post's content.
	LikeCount   int                  // LikeCount is the number of likes the post has.
	Comments    []Comment            // Comments is a list of comments on the post.
	Tags        []string             // Tags is a list of tags associated with the post.
}

// NewPost creates a new Post with the given ID, author, timestamp, title, description, content file URL, like count, comments, and tags.
func NewPost(
	id snowflake.Identifier,
	author User,
	timestamp time.Time,
	title string,
	description string,
	contentFile url.URL,
	likeCount int,
	comments []Comment,
	tags []string,
) Post {
	return Post{
		id:          id,
		Author:      author,
		Timestamp:   timestamp,
		Title:       title,
		Description: description,
		ContentFile: contentFile,
		LikeCount:   likeCount,
		Comments:    comments,
		Tags:        tags,
	}
}

// Id returns the ID of the post.
func (p Post) Id() uint64 {
	return p.id.Id()
}
