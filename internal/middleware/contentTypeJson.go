// Package middleware contains common middleware functions
package middleware

import "net/http"

// NewSetContentTypeJson returns middleware which sets the Content-Type header
// to application/json
func NewSetContentTypeJson() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
