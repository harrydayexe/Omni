package connector

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/pkg/errors"
)

// Connector is an interface that defines the methods that a data provider for
// OmniView must implement
type Connector interface {
	// GetPost returns a post by its id
	GetPost(ctx context.Context, id snowflake.Identifier) (storage.Post, error)
	// GetUser returns a user by its id
	// GetUser(ctx context.Context, id snowflake.Identifier) (storage.User, error)
	// GetUserPosts returns all posts by a user
	// GetUserPosts(ctx context.Context, id snowflake.Identifier) ([]storage.Post, error)
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

func (c *APIConnector) GetPost(ctx context.Context, id snowflake.Identifier) (storage.Post, error) {
	c.logger.InfoContext(ctx, "GetPost called", slog.Int64("id", int64(id.Id().ToInt())))
	postUrl, err := c.cfg.ReadApiUrl.Parse("/post/" + strconv.FormatUint(id.Id().ToInt(), 10))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get post url", slog.Any("error", err))
		return storage.Post{}, errors.Wrap(err, "c.cfg.ReadApiUrl.Parse")
	}

	resp, err := http.Get(postUrl.String())
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to get post", slog.Any("error", err))
		return storage.Post{}, errors.Wrap(err, "http.Get")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "failed to get post", slog.Int("http status", resp.StatusCode))
		return storage.Post{}, errors.New("api returned non-200 status")
	}

	var post storage.Post
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&post)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode post", slog.Any("error", err))
		return storage.Post{}, errors.Wrap(err, "Decode")
	}

	return post, nil
}

func (c *APIConnector) GetMostRecentPosts(ctx context.Context, page int) ([]storage.GetPostsPagedRow, error) {
	c.logger.InfoContext(ctx, "GetMostRecentPosts called", slog.Int("page num", page))
	postsUrl, err := c.cfg.ReadApiUrl.Parse("/posts?page=" + strconv.Itoa(page))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative get posts url", slog.Any("error", err))
		return nil, errors.Wrap(err, "c.cfg.ReadApiUrl.Parse")
	}

	resp, err := http.Get(postsUrl.String())
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to get posts", slog.Any("error", err))
		return nil, errors.Wrap(err, "http.Get")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "failed to get posts", slog.Int("http status", resp.StatusCode))
		return nil, errors.New("api returned non-200 status")
	}

	var posts []storage.GetPostsPagedRow
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&posts)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode posts", slog.Any("error", err))
		return nil, errors.Wrap(err, "Decode")
	}

	return posts, nil
}
