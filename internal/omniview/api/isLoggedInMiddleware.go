package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// IsLoggedInCtxKey is the key used to store the is-logged-in value in the context.
const IsLoggedInCtxKey string = "is-logged-in"

// UserIdCtxKey is the key used to store the user id in the context.
const UserIdCtxKey string = "user-id"

const AuthTokenCtxKey string = "jwt-token"

// NewIsLoggedInMiddleware returns middleware which checks if the user is logged in
// and saves the result in the context.
// If the user is not logged in, the middleware will also redirect to the login page,
// when restricted routes are accessed.
func newIsLoggedInMiddleware(logger *slog.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id, jwt, prs := hasValidAuthToken(r, logger); prs {
				ctx := context.WithValue(r.Context(), IsLoggedInCtxKey, true)
				ctx2 := context.WithValue(ctx, UserIdCtxKey, id.Id().ToInt())
				ctx3 := context.WithValue(ctx2, AuthTokenCtxKey, jwt)
				newReq := r.WithContext(ctx3)

				if r.URL.Path == "/login" {
					http.Redirect(w, newReq, "/", http.StatusSeeOther)
					return
				}

				next.ServeHTTP(w, newReq)
			} else {
				ctx := context.WithValue(r.Context(), IsLoggedInCtxKey, false)
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

// GetUserIdFromCtx returns the user id from the context if a user is logged in.
// Otherwise, it returns nil.
func GetUserIdFromCtx(ctx context.Context) snowflake.Identifier {
	if ctx.Value(IsLoggedInCtxKey) == true {
		return snowflake.ParseId(ctx.Value(UserIdCtxKey).(uint64))
	}
	return nil
}
