package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrDbFailed = errors.New("failed to read from db")
var ErrUnauthorized = errors.New("unauthorized")
var ErrPasswordTooLong = errors.New("password too long")
var ErrPasswordGen = errors.New("failed to generate password hash")

type Authable interface {
	// Login checks if the password for a given user id matches the stored hash
	Login(context.Context, snowflake.Identifier, string) (bool, error)
	// Signup creates a hash for the given password
	Signup(context.Context, string) ([]byte, error)
}

type AuthService struct {
	db     storage.Querier
	logger *slog.Logger
}

func (a *AuthService) Login(
	ctx context.Context,
	id snowflake.Identifier,
	password string,
) (bool, error) {
	// Check if user with id exists and retrieve their password hash
	hash, err := a.db.GetPasswordByID(ctx, int64(id.Id().ToInt()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.InfoContext(ctx, "user not found", slog.Any("id", id))
			return false, ErrUserNotFound
		}
		a.logger.ErrorContext(ctx, "failed to read user from db", slog.Any("error", err))
		return false, ErrDbFailed
	}

	// Compare given password against hash
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		a.logger.InfoContext(ctx, "incorrect password", slog.Any("id", id))
		return false, ErrUnauthorized
	}

	return true, nil
}

func (a *AuthService) Signup(ctx context.Context, password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		a.logger.InfoContext(ctx, "password too long")
		return nil, ErrPasswordTooLong
	}
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to generate password hash", slog.Any("error", err))
		return nil, ErrPasswordGen
	}

	return hash, nil
}
