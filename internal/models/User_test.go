package models

import (
	"testing"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestUserId(t *testing.T) {
	idGen := snowflake.NewSnowflakeGenerator(0)
	id := idGen.NextID()

	user := NewUser(id, "test", make([]snowflake.Snowflake, 0))

	if user.Id() != id {
		t.Errorf("User id does not match, got %d", user.Id())
	}
}
