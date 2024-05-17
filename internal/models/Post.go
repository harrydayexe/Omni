package models

import (
	"net/url"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Post represents a blog post.
type Post struct {
	ID          snowflake.Identifier // ID is a unique identifier for the post.
	Author      User                 // Author is the user who wrote the post.
	Timestamp   time.Time            // Timestamp is the time the post was created.
	Title       string               // Title is the title of the post.
	Description string               // Description is a short description of the post.
	ContentFile url.URL              // ContentFile is the URL of the markdown file containing the post's content.
	LikeCount   int                  // LikeCount is the number of likes the post has.
	Comments    []Comment            // Comments is a list of comments on the post.
	Tags        []string             // Tags is a list of tags associated with the post.
}
