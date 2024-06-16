package models

import (
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestCommentId(t *testing.T) {
	idGen := snowflake.NewSnowflakeGenerator(0)
	id := idGen.NextID()

	comment := NewComment(id, idGen.NextID(), idGen.NextID(), "test", time.Now(), "test")

	if comment.Id() != id {
		t.Errorf("Comment id does not match, got %d", comment.Id())
	}
}
