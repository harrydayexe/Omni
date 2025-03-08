package middleware

import "net/http"

func NewMaxBytesReader() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.MaxBytesReader(w, r.Body, 1048576)
			next.ServeHTTP(w, r)
		})
	}
}
