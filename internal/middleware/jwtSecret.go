package middleware

import (
	"context"
	"net/http"
)

// JWTCtxKey is the key used to store the jwt secret in the context.
const JWTCtxKey string = "jwt-secret"

// NewJwtSecret returns middleware which sets the jwt secret in the context.
func NewJwtSecret(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), JWTCtxKey, secret)
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		})
	}
}
