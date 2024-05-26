package models

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Post represents a blog post.
type Post struct {
	id          snowflake.Snowflake   // ID is a unique identifier for the post.
	AuthorId    snowflake.Snowflake   // AuthorId is the id of the user who wrote the post.
	AuthorName  string                // AuthorName is the name of the user who wrote the post.
	Timestamp   time.Time             // Timestamp is the time the post was created.
	Title       string                // Title is the title of the post.
	Description string                // Description is a short description of the post.
	ContentFile url.URL               // ContentFile is the URL of the markdown file containing the post's content.
	LikeCount   int                   // LikeCount is the number of likes the post has.
	Comments    []snowflake.Snowflake // Comments is a list of comment id's on the post.
	Tags        []string              // Tags is a list of tags associated with the post.
}

// NewPost creates a new Post with the given ID, author, timestamp, title, description, content file URL, like count, comments, and tags.
func NewPost(
	id snowflake.Snowflake,
	authorId snowflake.Snowflake,
	authorName string,
	timestamp time.Time,
	title string,
	description string,
	contentFile url.URL,
	likeCount int,
	comments []snowflake.Snowflake,
	tags []string,
) Post {
	return Post{
		id:         id,
		AuthorId:   authorId,
		AuthorName: authorName,
		Timestamp:  timestamp,
		Title:      title,

		Description: description,
		ContentFile: contentFile,
		LikeCount:   likeCount,
		Comments:    comments,
		Tags:        tags,
	}
}

// Id returns the ID of the post.
func (p Post) Id() snowflake.Snowflake {
	return p.id
}

func (p Post) MarshalJSON() ([]byte, error) {
	var comments []uint64
	for _, comment := range p.Comments {
		comments = append(comments, comment.ToInt())
	}

	postAltered := struct {
		Id          uint64   `json:"id"`
		AuthorId    uint64   `json:"authorId"`
		AuthorName  string   `json:"authorName"`
		Timestamp   string   `json:"timestamp"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		ContentFile string   `json:"contentFileUrl"`
		LikeCount   int      `json:"likeCount"`
		Comments    []uint64 `json:"comments"`
		Tags        []string `json:"tags"`
	}{
		Id:          p.Id().ToInt(),
		AuthorId:    p.AuthorId.ToInt(),
		AuthorName:  p.AuthorName,
		Timestamp:   p.Timestamp.Format(time.RFC3339),
		Title:       p.Title,
		Description: p.Description,
		ContentFile: p.ContentFile.String(),
		LikeCount:   p.LikeCount,
		Comments:    comments,
		Tags:        p.Tags,
	}
	return json.Marshal(postAltered)
}
