package models

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Post represents a blog post.
type Post struct {
	id          snowflake.Identifier // ID is a unique identifier for the post.
	AuthorId    snowflake.Identifier // AuthorId is the id of the user who wrote the post.
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
	authorId snowflake.Identifier,
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
		AuthorId:    authorId,
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

func (p Post) MarshalJSON() ([]byte, error) {
	postAltered := struct {
		Id          uint64    `json:"id"`
		AuthorId    uint64    `json:"authorId"`
		Timestamp   string    `json:"timestamp"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		ContentFile string    `json:"contentFile"`
		LikeCount   int       `json:"likeCount"`
		Comments    []Comment `json:"comments"`
		Tags        []string  `json:"tags"`
	}{
		Id:          p.id.Id(),
		AuthorId:    p.AuthorId.Id(),
		Timestamp:   p.Timestamp.Format(time.RFC3339),
		Title:       p.Title,
		Description: p.Description,
		ContentFile: p.ContentFile.String(),
		LikeCount:   p.LikeCount,
		Comments:    p.Comments,
		Tags:        p.Tags,
	}
	return json.Marshal(postAltered)
}
