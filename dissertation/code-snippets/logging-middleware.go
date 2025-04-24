package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// middleware is a function that wraps http.Handlers
// proving functionality before and after execution
// of the h handler.
type Middleware func(h http.Handler) http.Handler

type WrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// NewLoggingMiddleware returns middleware which logs incoming requests
func NewLoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				wrapped := &WrappedWriter{
					ResponseWriter: w,
					statusCode:     http.StatusOK,
				}

				logger.InfoContext(r.Context(),
					"handling incoming request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)

				next.ServeHTTP(wrapped, r)

				logger.InfoContext(r.Context(),
					"finished handling request",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status_code", wrapped.statusCode),
					slog.Duration("duration", time.Since(start)),
				)
			})
	}
}
