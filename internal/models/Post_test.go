package models

import (
	"net/url"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestPostId(t *testing.T) {
	idGen := snowflake.NewSnowflakeGenerator(0)
	id := idGen.NextID()

	post := NewPost(id, idGen.NextID(), "test", time.Now(), "test", "test", url.URL{}, make([]snowflake.Snowflake, 0), make([]string, 0))

	if post.Id() != id {
		t.Errorf("Post id does not match, got %d", post.Id())
	}
}
