package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
)

// isLoggedInCtxKey is the key used to store the is-logged-in value in the context.
const isLoggedInCtxKey string = "is-logged-in"

// NewIsLoggedInMiddleware returns middleware which checks if the user is logged in
// and saves the result in the context.
// If the user is not logged in, the middleware will also redirect to the login page,
// when restricted routes are accessed.
func newIsLoggedInMiddleware(logger *slog.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, prs := hasValidAuthToken(r, logger); prs {
				ctx := context.WithValue(r.Context(), isLoggedInCtxKey, true)
				newReq := r.WithContext(ctx)

				if r.URL.Path == "/login" {
					http.Redirect(w, newReq, "/", http.StatusSeeOther)
					return
				}

				next.ServeHTTP(w, newReq)
			} else {
				ctx := context.WithValue(r.Context(), isLoggedInCtxKey, false)
				newReq := r.WithContext(ctx)
				if r.URL.Path == "/post/new" {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				} else if r.URL.Path == "/logout" {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				next.ServeHTTP(w, newReq)
			}
		})
	}
}
