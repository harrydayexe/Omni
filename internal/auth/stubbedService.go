package auth

import (
	"context"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Implement the auth.Authable interface
type StubbedAuthService struct {
	VerifyTokenFn func(ctx context.Context, token string, id snowflake.Identifier) error
	LoginFn       func(ctx context.Context, id snowflake.Identifier, password string) (string, error)
	SignupFn      func(ctx context.Context, password string) ([]byte, error)
}

func (m StubbedAuthService) VerifyToken(ctx context.Context, token string, id snowflake.Identifier) error {
	return m.VerifyTokenFn(ctx, token, id)
}

func (m StubbedAuthService) Login(ctx context.Context, id snowflake.Identifier, password string) (string, error) {
	return m.LoginFn(ctx, id, password)
}

func (m StubbedAuthService) Signup(ctx context.Context, password string) ([]byte, error) {
	return m.SignupFn(ctx, password)
}
