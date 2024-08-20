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
	postId := g.NextID()
	timestamp := time.Now()

	c := NewComment(commentId, postId, authorId, "Test Name", timestamp, "Hello, world!")

	commentJson, err := c.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal comment to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"postId":%d,"authorId":%d,"authorName":"Test Name","timestamp":"%s","content":"Hello, world!"}`, commentId.ToInt(), postId.ToInt(), authorId.ToInt(), timestamp.UTC().Format(time.RFC3339))

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
	comments := []snowflake.Snowflake{commentId}
	tags := []string{"tag1", "tag2"}

	p := NewPost(postId, authorId, "Test Name", timestamp, "Hello, world!", "Lorem Ipsum Dolar", url, comments, tags)

	postJson, err := p.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal post to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"authorId":%d,"authorName":"Test Name","timestamp":"%s","title":"Hello, world!","description":"Lorem Ipsum Dolar","contentFileUrl":"https://example.com/example","comments":[%d],"tags":["tag1","tag2"]}`, postId.ToInt(), authorId.ToInt(), timestamp.UTC().Format(time.RFC3339), commentId.ToInt())

	if string(postJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(postJson))
	}
}

func TestMarshallUserJson(t *testing.T) {
	g := snowflake.NewSnowflakeGenerator(1)
	postId := g.NextID()
	userId := g.NextID()
	username := "testuser"

	u := NewUser(userId, username, []snowflake.Snowflake{postId})

	userJson, err := u.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal user to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"username":"testuser","posts":[%d]}`, userId.ToInt(), postId.ToInt())

	if string(userJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(userJson))
	}
}

func TestMarshallUserJsonEmptyPosts(t *testing.T) {
	g := snowflake.NewSnowflakeGenerator(1)
	userId := g.NextID()
	username := "testuser"

	u := NewUser(userId, username, []snowflake.Snowflake{})

	userJson, err := u.MarshalJSON()
	if err != nil {
		t.Errorf("Failed to marshal user to json: %v", err)
	}

	expected := fmt.Sprintf(`{"id":%d,"username":"testuser","posts":[]}`, userId.ToInt())

	if string(userJson) != expected {
		t.Errorf("Expected %s, got %s", expected, string(userJson))
	}
}
