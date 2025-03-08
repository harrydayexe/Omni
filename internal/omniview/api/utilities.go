package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/microcosm-cc/bluemonday"
	"github.com/oxtoacart/bpool"
)

func writeTemplateWithBuffer(ctx context.Context, logger *slog.Logger, name string, t *templates.Templates, bufpool *bpool.BufferPool, w http.ResponseWriter, content interface{}) {
	// Get buffer
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := t.Templates.ExecuteTemplate(buf, name, content)
	if err != nil {
		logger.ErrorContext(ctx, "Error executing template", slog.String("error message", err.Error()))
		// NOTE: We are assuming here that this error page won't fail to render
		_ = t.Templates.ExecuteTemplate(w, "errorpage.html", datamodels.NewErrorPageModel("Internal Server Error", "An error occurred while rendering the page."))
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// FetchMarkdownData fetches markdown data from a given url and returns the sanitized html
func fetchMarkdownData(ctx context.Context, logger *slog.Logger, url string) (string, error) {
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
func hasValidAuthHeader(r *http.Request, logger *slog.Logger) (string, bool) {
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

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// formUrlEncoded is the content type for data sent from a form
const formUrlEncoded = "application/x-www-form-urlencoded"

// check that the content type header is application/json
// adapted from https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func checkContentTypeHeader(logger *slog.Logger, r *http.Request, expected string) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != expected {
			msg := fmt.Sprintf("Content-Type header is not %s", expected)
			logger.InfoContext(r.Context(), msg)
			return fmt.Errorf("Content-Type header is not %s", expected)
		}
	}
	return nil
}
