package auth

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

func TestLogin(t *testing.T) {
	var testLogger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	var cases = []struct {
		name              string
		id                snowflake.Identifier
		password          string
		secretKey         string
		GetPasswordByIDFn func(context.Context, int64) (string, error)
		expectedErr       error
	}{
		{
			name:      "Valid login",
			id:        snowflake.ParseId(1796290045997481984),
			password:  "password",
			secretKey: "omni-secret",
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6", nil
			},
			expectedErr: nil,
		},
		{
			name:      "Invalid login",
			id:        snowflake.ParseId(1796290045997481984),
			password:  "invalid",
			secretKey: "omni-secret",
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6", nil
			},
			expectedErr: ErrUnauthorized,
		},
		{
			name:      "Unknown user",
			id:        snowflake.ParseId(1796290045997481984),
			password:  "invalid",
			secretKey: "omni-secret",
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "", sql.ErrNoRows
			},
			expectedErr: ErrUserNotFound,
		},
		{
			name:      "db error",
			id:        snowflake.ParseId(1796290045997481984),
			password:  "invalid",
			secretKey: "omni-secret",
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "", fmt.Errorf("db error")
			},
			expectedErr: ErrDbFailed,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service := NewAuthService(
				[]byte(c.secretKey),
				&storage.StubbedQueries{
					GetPasswordByIDFn: c.GetPasswordByIDFn,
				},
				testLogger,
			)

			_, err := service.Login(context.Background(), c.id, c.password)
			if err != c.expectedErr {
				t.Errorf("Expected error to be %v, got %v", c.expectedErr, err)
			}
		})
	}
}
