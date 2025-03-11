package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/omniwrite/datamodels"
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
	// Login logs a user in and returns a token
	Login(ctx context.Context, username, password string) (auth.LoginResponse, error)
	// Signup signs a user up and returns the user object
	Signup(ctx context.Context, username, password string) (datamodels.NewUserResponse, error)
	// CreatePost creates a post and returns the new post
	CreatePost(ctx context.Context, newPost datamodels.NewPost) (storage.Post, error)
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

func (c *APIConnector) Login(ctx context.Context, username, password string) (auth.LoginResponse, error) {
	c.logger.InfoContext(ctx, "Login called", slog.String("username", username))
	loginUrl, err := c.cfg.AuthApiUrl.Parse("/login")
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative login url", slog.Any("error", err))
		return auth.LoginResponse{}, NewAPIError(0, err)
	}

	postData := auth.LoginRequest{
		Username: username,
		Password: password,
	}
	postDataBytes, err := json.Marshal(postData)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to marshal login request", slog.Any("error", err))
		return auth.LoginResponse{}, NewAPIError(0, err)
	}

	resp, err := http.Post(loginUrl.String(), "application/json", bytes.NewBuffer(postDataBytes))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to send POST request to backend", slog.Any("error", err))
		return auth.LoginResponse{}, NewAPIError(0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.InfoContext(ctx, "POST request did not return 200", slog.Int("http status", resp.StatusCode))
		return auth.LoginResponse{}, NewAPIError(resp.StatusCode, nil)
	}

	var loginResponse auth.LoginResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&loginResponse)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode login response", slog.Any("error", err))
		return auth.LoginResponse{}, NewAPIError(0, err)
	}

	return loginResponse, nil
}

func (c *APIConnector) Signup(
	ctx context.Context,
	username, password string,
) (datamodels.NewUserResponse, error) {
	c.logger.InfoContext(ctx, "Signup called", slog.String("username", username))
	signupUrl, err := c.cfg.WriteApiUrl.Parse("/user")
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative signup url", slog.Any("error", err))
		return datamodels.NewUserResponse{}, NewAPIError(0, err)
	}

	postData := datamodels.NewUserRequest{
		Username: username,
		Password: password,
	}
	postDataBytes, err := json.Marshal(postData)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to marshal signup request", slog.Any("error", err))
		return datamodels.NewUserResponse{}, NewAPIError(0, err)
	}

	resp, err := http.Post(signupUrl.String(), "application/json", bytes.NewBuffer(postDataBytes))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to send POST request to backend", slog.Any("error", err))
		return datamodels.NewUserResponse{}, NewAPIError(0, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		c.logger.InfoContext(ctx, "POST request did not return 201", slog.Int("http status", resp.StatusCode))
		return datamodels.NewUserResponse{}, NewAPIError(resp.StatusCode, nil)
	}

	var signupResponse datamodels.NewUserResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&signupResponse)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode signup response", slog.Any("error", err))
		return datamodels.NewUserResponse{}, NewAPIError(0, err)
	}

	return signupResponse, nil
}

func (c *APIConnector) CreatePost(
	ctx context.Context,
	newPost datamodels.NewPost,
) (storage.Post, error) {
	c.logger.InfoContext(ctx, "CreatePost called", slog.Any("newPost", newPost))

	if ctx.Value("jwt-token") == nil {
		c.logger.ErrorContext(ctx, "no auth token in context")
		return storage.Post{}, NewAPIError(0, fmt.Errorf("no auth token in context"))
	}

	newPostUrl, err := c.cfg.WriteApiUrl.Parse("/post")
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to parse relative post url", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}
	c.logger.DebugContext(ctx, "created new post url", slog.String("url", newPostUrl.String()))

	bodyData, err := json.Marshal(newPost)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to marshal new post", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}
	c.logger.DebugContext(ctx, "created body data", slog.String("data", string(bodyData)))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, newPostUrl.String(), bytes.NewBuffer(bodyData))
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to create POST request", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}
	req.Header.Add("Authorization", "Bearer "+ctx.Value("jwt-token").(string))
	req.Header.Add("Content-Type", "application/json")
	c.logger.DebugContext(ctx, "created POST request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to send POST request to backend", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}
	defer resp.Body.Close()
	c.logger.DebugContext(ctx, "sent POST request successfully")

	if resp.StatusCode != http.StatusCreated {
		c.logger.InfoContext(ctx, "POST request did not return 201", slog.Int("http status", resp.StatusCode))
		return storage.Post{}, NewAPIError(resp.StatusCode, nil)
	}
	c.logger.DebugContext(ctx, "POST request returned 201")

	var post storage.Post
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&post)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to decode post", slog.Any("error", err))
		return storage.Post{}, NewAPIError(0, err)
	}
	c.logger.DebugContext(ctx, "decoded post", slog.Any("post", post))

	return post, nil
}
