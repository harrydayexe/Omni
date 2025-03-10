package auth

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// Checks if a given token is valid and returns the expiration time as a
// formatted RFC1123 string.
func IsValidToken(
	ctx context.Context,
	tokenString string,
	logger *slog.Logger,
) (snowflake.Identifier, error) {
	logger.DebugContext(ctx, "validating token", slog.String("token", tokenString))

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			jwtSecret, ok := ctx.Value(middleware.JWTCtxKey).(string)
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
		return snowflake.Snowflake{}, ErrTokenInvalid
	}

	if !token.Valid {
		logger.DebugContext(ctx, "token is not valid")
		return snowflake.Snowflake{}, ErrTokenInvalid
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		logger.ErrorContext(ctx, "failed to get subject after verifying token", slog.Any("error", err))
		return snowflake.Snowflake{}, ErrTokenInvalid
	}

	idNum, err := strconv.Atoi(sub)

	id := snowflake.ParseId(uint64(idNum))

	return id, nil
}
