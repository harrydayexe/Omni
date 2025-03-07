package middleware

import (
	"context"
	"net/http"
)

func NewJwtSecret(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "jwt-secret", secret)
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		})
	}
}
