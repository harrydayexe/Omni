package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/microcosm-cc/bluemonday"
)

// FetchMarkdownData fetches markdown data from a given url and returns the sanitized html
func FetchMarkdownData(ctx context.Context, logger *slog.Logger, url string) (string, error) {
	logger.DebugContext(ctx, "Fetching markdown data", slog.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		logger.InfoContext(ctx, "Failed to fetch markdown data", slog.String("error", err.Error()))
		return "", connector.NewAPIError(0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.InfoContext(ctx, "Failed to fetch markdown data", slog.Int("status", resp.StatusCode))
		return "", connector.NewAPIError(resp.StatusCode, nil)
	}

	rawMdBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.InfoContext(ctx, "Failed to read markdown data", slog.String("error", err.Error()))
		return "", connector.NewAPIError(0, err)
	}

	logger.DebugContext(ctx, "Markdown data fetched", slog.Int("length", len(rawMdBytes)))
	maybeUnsafeHTML := markdown.ToHTML(rawMdBytes, nil, nil)
	html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)

	return string(html), nil
}

// HasValidAuthHeader checks if the request has a valid Authorization header.
// It does not check if the given header is valid for a given user id.
func HasValidAuthHeader(r *http.Request, logger *slog.Logger) (string, bool) {
	logger.DebugContext(r.Context(), "Checking for valid auth header")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg := "Authorization header is missing"
		logger.InfoContext(r.Context(), msg)
		return "", false
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		msg := "Authorization header format must be Bearer {token}"
		logger.InfoContext(r.Context(), msg)
		return "", false
	}

	time, err := auth.IsValidToken(r.Context(), authHeaderParts[1], logger)
	if err != nil {
		return "", false
	}

	return time, true
}
