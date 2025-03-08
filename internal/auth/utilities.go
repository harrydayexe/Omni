package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Checks if a given token is valid and returns the expiration time as a
// formatted RFC1123 string.
func IsValidToken(
	ctx context.Context,
	tokenString string,
	logger *slog.Logger,
) (string, error) {
	logger.DebugContext(ctx, "validating token", slog.String("token", tokenString))

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			jwtSecret, ok := ctx.Value("jwt-secret").(string)
			if !ok {
				// handle the error, e.g., log or return an error
				panic("jwt-secret could not be cast to a string")
			}
			byteSecret := []byte(jwtSecret)
			return byteSecret, nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		logger.DebugContext(ctx, "token is not valid", slog.Any("error", err))
		return "", ErrTokenInvalid
	}

	if !token.Valid {
		logger.DebugContext(ctx, "token is not valid")
		return "", ErrTokenInvalid
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		logger.ErrorContext(ctx, "failed to get expiration time after verifying token", slog.Any("error", err))
		return "", ErrTokenInvalid
	}

	return exp.Format(time.RFC1123), nil
}
