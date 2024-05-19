package models

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestMarshalCommentJSON(t *testing.T) {
	g := snowflake.NewSnowflakeGenerator(1)
	commentId := g.NextID()
	authorId := g.NextID()
	timestamp := time.Now()

	c := NewComment(commentId, authorId, timestamp, "Hello, world!", 0)

	commentJson, err := c.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal comment to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"authorId":%d,"timestamp":"%s","content":"Hello, world!","likeCount":0}`, commentId.Id(), authorId.Id(), timestamp.Format(time.RFC3339))

	if string(commentJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(commentJson))
	}
}

func TestMarshalPostJSON(t *testing.T) {
	g := snowflake.NewSnowflakeGenerator(1)
	postId := g.NextID()
	authorId := g.NextID()
	commentId := g.NextID()
	timestamp := time.Now()
	url := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/example",
	}
	comments := []Comment{NewComment(commentId, authorId, timestamp, "Hello, world!", 0)}
	tags := []string{"tag1", "tag2"}

	p := NewPost(postId, authorId, timestamp, "Hello, world!", "Lorem Ipsum Dolar", url, 10, comments, tags)

	postJson, err := p.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal post to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"authorId":%d,"timestamp":"%s","title":"Hello, world!","description":"Lorem Ipsum Dolar","contentFileUrl":"https://example.com/example","likeCount":10,"comments":[%d],"tags":["tag1","tag2"]}`, postId.Id(), authorId.Id(), timestamp.Format(time.RFC3339), commentId.Id())

	if string(postJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(postJson))
	}
}

func TestMarshallUserJson(t *testing.T) {
	g := snowflake.NewSnowflakeGenerator(1)
	postId := g.NextID()
	authorId := g.NextID()
	timestamp := time.Now()
	url := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/example",
	}
	comments := []Comment{}
	tags := []string{}
	userId := g.NextID()
	username := "testuser"
	posts := []Post{NewPost(postId, authorId, timestamp, "Hello, world!", "Lorem Ipsum Dolar", url, 10, comments, tags)}

	u := NewUser(userId, username)
	u.Posts = posts

	userJson, err := u.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal user to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"username":"testuser","posts":[%d]}`, userId.Id(), postId.Id())

	if string(userJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(userJson))
	}
}
