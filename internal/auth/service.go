package auth

import "github.com/harrydayexe/Omni/internal/snowflake"

type Authable interface {
	Login(snowflake.Identifier, string) (bool, error)
	Signup(snowflake.Identifier, string) (bool, error)
}

type AuthService struct {
}
