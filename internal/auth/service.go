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
var ErrPasswordGen = errors.New("failed to generate password hash")

type Authable interface {
	// VerifyToken checks if the given token is valid
	VerifyToken(context.Context, string) error
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

func (a *AuthService) VerifyToken(ctx context.Context, tokenString string) error {
	return a.verifyToken(tokenString)
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
	token, err := a.createToken(id)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed to create token", slog.Any("error", err))
		return "", ErrTokenGenFail
	}

	return token, nil
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

func (a *AuthService) createToken(id snowflake.Identifier) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   fmt.Sprintf("%d", id.Id().ToInt()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
