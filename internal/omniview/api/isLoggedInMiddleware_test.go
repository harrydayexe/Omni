package api

import (
	"context"
	"testing"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestGetUserIdFromCtx(t *testing.T) {
	loggedInCtx := context.WithValue(context.Background(), IsLoggedInCtxKey, true)
	loggedInCtx = context.WithValue(loggedInCtx, UserIdCtxKey, snowflake.ParseId(1).ToInt())
	loggedOutCtx := context.WithValue(context.Background(), IsLoggedInCtxKey, false)
	var cases = []struct {
		name     string
		ctx      context.Context
		expected snowflake.Identifier
	}{
		{
			name:     "User is logged in",
			ctx:      loggedInCtx,
			expected: snowflake.ParseId(1),
		},
		{
			name:     "User is not logged in",
			ctx:      loggedOutCtx,
			expected: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := GetUserIdFromCtx(c.ctx)
			if result != c.expected {
				t.Errorf("Expected %v, got %v", c.expected, result)
			}
		})
	}
}
