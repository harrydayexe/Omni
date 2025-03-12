package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

func TestLogin(t *testing.T) {

	var cases = []struct {
		name                string
		username            string
		password            string
		secretKey           string
		GetPasswordByIDFn   func(context.Context, int64) (string, error)
		GetUserByUsernameFn func(context.Context, string) (int64, error)
		expectedErr         error
	}{
		{
			name:      "Valid login",
			username:  "test",
			password:  "password",
			secretKey: "omni-secret",
			GetUserByUsernameFn: func(ctx context.Context, username string) (int64, error) {
				return 1796290045997481984, nil
			},
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6", nil
			},
			expectedErr: nil,
		},
		{
			name:      "Invalid login",
			username:  "test",
			password:  "invalid",
			secretKey: "omni-secret",
			GetUserByUsernameFn: func(ctx context.Context, username string) (int64, error) {
				return 1796290045997481984, nil
			},
			GetPasswordByIDFn: func(ctx context.Context, id int64) (string, error) {
				return "$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6", nil
			},
			expectedErr: ErrUnauthorized,
		},
		{
			name:      "Unknown user",
			username:  "test",
			password:  "invalid",
			secretKey: "omni-secret",
			GetUserByUsernameFn: func(ctx context.Context, username string) (int64, error) {
				return 0, sql.ErrNoRows
			},
			expectedErr: ErrUserNotFound,
		},
		{
			name:      "db error on user lookup",
			username:  "test",
			password:  "invalid",
			secretKey: "omni-secret",
			GetUserByUsernameFn: func(ctx context.Context, username string) (int64, error) {
				return 0, fmt.Errorf("db error")
			},
			expectedErr: ErrDbFailed,
		},
		{
			name:      "db error on hash lookup",
			username:  "test",
			password:  "invalid",
			secretKey: "omni-secret",
			GetUserByUsernameFn: func(ctx context.Context, username string) (int64, error) {
				return 1796290045997481984, nil
			},
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
					GetUserByUsernameFn: c.GetUserByUsernameFn,
					GetPasswordByIDFn:   c.GetPasswordByIDFn,
				},
				testLogger,
			)

			_, err := service.Login(context.Background(), c.username, c.password)
			if err != c.expectedErr {
				t.Errorf("Expected error to be %v, got %v", c.expectedErr, err)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	var cases = []struct {
		name        string
		tokenString string
		id          snowflake.Identifier
		secretKey   string
		expectedErr error
	}{
		{
			name:        "Valid token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjk1NzkzNTg1LCJzdWIiOiIxNzk2MjkwMDQ1OTk3NDgxOTg0In0.RMBAJGkKahsECMiOpDcib__YU1CTCWEf4C_h7m_4HJs",
			id:          snowflake.ParseId(1796290045997481984),
			secretKey:   "omni-secret",
			expectedErr: nil,
		},
		{
			name:        "Invalid token",
			tokenString: "eyJhbGxxxxxxxxiOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjk1NzkzNTg1LCJzdWIiOiIxNzk2MjkwMDQ1OTk3NDgxOTg0In0.RMBAJGkKahsECMiOpDcib__YU1CTCWEf4C_h7m_4HJs",
			id:          snowflake.ParseId(1796290045997481984),
			secretKey:   "omni-secret",
			expectedErr: ErrTokenInvalid,
		},
		{
			name:        "Invalid subject",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjk1NzkzNTg1LCJzdWIiOiIxNzk2MjkwMDQ1OTk3NDgxOTg1In0.fZ4lcr1VcC8iAu45CCPHRAXvERtwE0RzKkdWU3HFvAk",
			id:          snowflake.ParseId(1796290045997481984),
			secretKey:   "omni-secret",
			expectedErr: ErrTokenInvalid,
		},
		{
			name:        "Expired token",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzg1ODQ2ODUsInN1YiI6IjE3OTYyOTAwNDU5OTc0ODE5ODQifQ.Vy565OuUSSOdT9vusvmKNDaWPAcQVS7wlrE537sH2AA",
			id:          snowflake.ParseId(1796290045997481984),
			secretKey:   "omni-secret",
			expectedErr: ErrTokenInvalid,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service := NewAuthService(
				[]byte(c.secretKey),
				nil,
				testLogger,
			)

			err := service.VerifyToken(context.Background(), c.tokenString, c.id)
			if err != c.expectedErr {
				t.Errorf("Expected error to be %v, got %v", c.expectedErr, err)
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	var cases = []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "Valid password",
			password:    "password",
			expectedErr: nil,
		},
		{
			name:        "Short password",
			password:    "pass",
			expectedErr: ErrPasswordTooShort,
		},
		{
			name:        "Long password",
			password:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890",
			expectedErr: ErrPasswordTooLong,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service := NewAuthService(
				[]byte("test-secret"),
				nil,
				testLogger,
			)

			_, err := service.Signup(context.Background(), c.password)
			if err != c.expectedErr {
				t.Errorf("Expected error to be %v, got %v", c.expectedErr, err)
			}
		})
	}
}
