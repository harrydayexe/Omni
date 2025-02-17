package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrDbFailed = errors.New("failed to read from db")
var ErrTokenGenFail = errors.New("failed to generate token")
var ErrUnauthorized = errors.New("unauthorized")
var ErrPasswordTooLong = errors.New("password too long")
var ErrPasswordTooShort = errors.New("password too short")
var ErrPasswordGen = errors.New("failed to generate password hash")
var ErrTokenInvalid = errors.New("invalid token")

type Authable interface {
	// VerifyToken checks if the given token is valid for a given id
	VerifyToken(context.Context, string, snowflake.Identifier) error
	// Login checks if the password for a given user id matches the stored hash
	Login(context.Context, snowflake.Identifier, string) (string, error)
	// Signup creates a hash for the given password
	Signup(context.Context, string) ([]byte, error)
}

type AuthService struct {
	secretKey []byte
	db        storage.Querier
	logger    *slog.Logger
}

func NewAuthService(secretKey []byte, db storage.Querier, logger *slog.Logger) *AuthService {
	return &AuthService{
		secretKey: secretKey,
		db:        db,
		logger:    logger,
	}
}

func (a *AuthService) VerifyToken(ctx context.Context, tokenString string, id snowflake.Identifier) error {
	a.logger.DebugContext(ctx, "verifying token", slog.String("token", tokenString))

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return a.secretKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithSubject(fmt.Sprintf("%d", id.Id().ToInt())),
	)
	if err != nil {
		a.logger.InfoContext(ctx, "invalid token", slog.Any("error", err))
		return ErrTokenInvalid
	}

	if !token.Valid {
		a.logger.InfoContext(ctx, "invalid token")
		return ErrTokenInvalid
	}

	return nil
}

func (a *AuthService) Login(
	ctx context.Context,
	id snowflake.Identifier,
	password string,
) (string, error) {
	// Check if user with id exists and retrieve their password hash
	hash, err := a.db.GetPasswordByID(ctx, int64(id.Id().ToInt()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.InfoContext(ctx, "user not found", slog.Any("id", id))
			return "", ErrUserNotFound
		}
		a.logger.ErrorContext(ctx, "failed to read user from db", slog.Any("error", err))
		return "", ErrDbFailed
	}

	// Compare given password against hash
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		a.logger.InfoContext(ctx, "incorrect password", slog.Any("id", id))
		return "", ErrUnauthorized
	}

	// Create a token for the user
	token, err := a.createToken(ctx, id)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to create token", slog.Any("error", err))
		return "", ErrTokenGenFail
	}

	return token, nil
}

func (a *AuthService) Signup(ctx context.Context, password string) ([]byte, error) {
	if len(password) < 7 {
		return nil, ErrPasswordTooShort
	}

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

// createToken generates a token for a given id
func (a *AuthService) createToken(ctx context.Context, id snowflake.Identifier) (string, error) {
	a.logger.DebugContext(ctx, "creating token", slog.Any("id", id))

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   fmt.Sprintf("%d", id.Id().ToInt()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		a.logger.InfoContext(ctx, "failed to generate token", slog.Any("error", err))
		return "", err
	}

	return tokenString, nil
}
