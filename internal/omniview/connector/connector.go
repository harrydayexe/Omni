package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

// Connector is an interface that defines the methods that a data provider for
// OmniView must implement
type Connector interface {
	// GetPost returns a post by its id
	GetPost(ctx context.Context, id snowflake.Identifier) (storage.Post, error)
	// GetUser returns a user by its id
	GetUser(ctx context.Context, id snowflake.Identifier) (storage.User, error)
	// GetUserPosts returns all posts by a user
	GetUserPosts(ctx context.Context, id snowflake.Identifier) ([]storage.Post, error)
	// GetPostComments returns all comments on a post
	// GetPostComments(ctx context.Context, id snowflake.Identifier) ([]storage.Comment, error)
	// GetMostRecentPosts returns the most recent posts from the page
	GetMostRecentPosts(ctx context.Context, page int) ([]storage.GetPostsPagedRow, error)
}

// APIConnector is a struct that implements the Connector interface
type APIConnector struct {
	cfg    config.ViewConfig
	logger *slog.Logger
}

// NewAPIConnector creates a new APIConnector
func NewAPIConnector(cfg config.ViewConfig, logger *slog.Logger) *APIConnector {
	return &APIConnector{
		cfg:    cfg,
		logger: logger,
	}
}

// APIError is an error type that is returned when an API request fails.
// It contains the status code of the response or the underlying error.
// Only one is expected to be set per error.
type APIError struct {
	StatusCode int
	Underlying error
}

func (e *APIError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Errorf("API returned non-200 status code: %d", e.StatusCode).Error()
	} else {
		return fmt.Errorf("Connector error: %w", e.Underlying).Error()
	}
}

func NewAPIError(statusCode int, underlying error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Underlying: underlying,
	}
}

func (c *APIConnector) GetRequest(ctx context.Context, url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to send GET request to backend", slog.Any("error", err))
		return nil, NewAPIError(0, err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "GET request did not return 200", slog.Int("http status", resp.StatusCode))
		return nil, NewAPIError(resp.StatusCode, nil)
	}

	return resp, nil
}

func (c *APIConnector) GetPost(ctx context.Context, id snowflake.Identifier) (storage.Post, error) {
	c.logger.InfoContext(ctx, "GetPost called", slog.Int64("id", int64(id.Id().ToInt())))
	postUrl, err := c.cfg.ReadApiUrl.Parse("/post/" + strconv.FormatUint(id.Id().ToInt(), 10))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get post url", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}

	resp, err := c.GetRequest(ctx, postUrl.String())
	if err != nil {
		return storage.Post{}, err
	}
	defer resp.Body.Close()

	var post storage.Post
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&post)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode post", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}

	return post, nil
}

func (c *APIConnector) GetUser(ctx context.Context, id snowflake.Identifier) (storage.User, error) {
	c.logger.InfoContext(ctx, "GetUser called", slog.Int64("id", int64(id.Id().ToInt())))
	userUrl, err := c.cfg.ReadApiUrl.Parse("/user/" + strconv.FormatUint(id.Id().ToInt(), 10))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get user url", slog.Any("error", err))
		return storage.User{}, NewAPIError(0, err)
	}

	resp, err := c.GetRequest(ctx, userUrl.String())
	if err != nil {
		return storage.User{}, err
	}
	defer resp.Body.Close()

	var user storage.User
	userDecoder := json.NewDecoder(resp.Body)
	userDecoder.DisallowUnknownFields()

	err = userDecoder.Decode(&user)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode user", slog.Any("error", err))
		return storage.User{}, NewAPIError(0, err)
	}

	return user, nil
}

func (c *APIConnector) GetUserPosts(ctx context.Context, id snowflake.Identifier) ([]storage.Post, error) {
	c.logger.InfoContext(ctx, "GetUserPosts called", slog.Int64("id", int64(id.Id().ToInt())))
	userPostsUrl, err := c.cfg.ReadApiUrl.Parse("/user/" + strconv.FormatUint(id.Id().ToInt(), 10) + "/posts")
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get user posts url", slog.Any("error", err))
		return nil, NewAPIError(0, err)
	}

	resp, err := c.GetRequest(ctx, userPostsUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []storage.Post
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&posts)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode posts", slog.Any("error", err))
		return nil, NewAPIError(0, err)
	}

	return posts, nil
}

func (c *APIConnector) GetMostRecentPosts(ctx context.Context, page int) ([]storage.GetPostsPagedRow, error) {
	c.logger.InfoContext(ctx, "GetMostRecentPosts called", slog.Int("page num", page))
	postsUrl, err := c.cfg.ReadApiUrl.Parse("/posts?page=" + strconv.Itoa(page))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get posts url", slog.Any("error", err))
		return nil, NewAPIError(0, err)
	}

	resp, err := c.GetRequest(ctx, postsUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []storage.GetPostsPagedRow
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&posts)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode posts", slog.Any("error", err))
		return nil, NewAPIError(0, err)
	}

	return posts, nil
}
