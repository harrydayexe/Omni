// Package middleware contains common middleware functions for use in APIs
package middleware

import "net/http"

// middleware is a function that wraps http.Handlers
// proving functionality before and after execution
// of the h handler.
type Middleware func(h http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
