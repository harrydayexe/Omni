package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/storage"
)

type stubbedDB struct {
	ShouldReturnError bool
}

func (s *stubbedDB) PingContext(ctx context.Context) error {
	if s.ShouldReturnError {
		return fmt.Errorf("ping error")
	}
	return nil
}

func TestGetUserKnown(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			newUser := storage.GetUserByIDRow{
				ID:       id,
				Username: "johndoe",
			}
			return newUser, nil
		},
	}

	req := httptest.NewRequest("GET", "/user/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"id":1796290045997481984,"username":"johndoe"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserUnknown(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, sql.ErrNoRows
		},
	}

	req := httptest.NewRequest("GET", "/user/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetUserBadFormedId(t *testing.T) {
	req := httptest.NewRequest("GET", "/user/hello", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "Url parameter could not be parsed properly.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserDBError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, fmt.Errorf("database error")
		},
	}

	req := httptest.NewRequest("GET", "/user/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestGetPostKnown(t *testing.T) {
	expectedTime := time.Date(2021, 1, 1, 11, 40, 35, 0, time.UTC)
	expectedURL := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/foo",
	}
	mockedQueries := &storage.StubbedQueries{
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			newPost := storage.Post{
				ID:          id,
				UserID:      1796290045997481985,
				CreatedAt:   expectedTime,
				Title:       "Hello, World!",
				Description: "Foobarbaz",
				MarkdownUrl: expectedURL.String(),
			}
			return newPost, nil
		},
	}

	req := httptest.NewRequest("GET", "/post/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"id":1796290045997481984,"user_id":1796290045997481985,"created_at":"2021-01-01T11:40:35Z","title":"Hello, World!","markdown_url":"https://example.com/foo","description":"Foobarbaz"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostUnknown(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, sql.ErrNoRows
		},
	}

	req := httptest.NewRequest("GET", "/post/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetPostBadFormedId(t *testing.T) {
	req := httptest.NewRequest("GET", "/post/hello", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "Url parameter could not be parsed properly.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostDBError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, fmt.Errorf("database error")
		},
	}

	req := httptest.NewRequest("GET", "/post/1796290045997481984", nil)

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{ShouldReturnError: false},
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestGetPostsForUser(t *testing.T) {
	const userIdNum = 1796290045997481984
	const userIdString = "1796290045997481984"
	const basePostNum = 1796290045997481985
	const mdUrl = "https://example.com"

	expectedPosts := []storage.Post{
		{
			ID:          basePostNum,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			Title:       "Post 0",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 1,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			Title:       "Post 1",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 2,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
			Title:       "Post 2",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 3,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
			Title:       "Post 3",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 4,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC),
			Title:       "Post 4",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 5,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			Title:       "Post 5",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 6,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			Title:       "Post 6",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 7,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC),
			Title:       "Post 7",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 8,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 5, 4, 0, 0, 0, 0, time.UTC),
			Title:       "Post 8",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 9,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 6, 4, 0, 0, 0, 0, time.UTC),
			Title:       "Post 9",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
		{
			ID:          basePostNum + 10,
			UserID:      userIdNum,
			CreatedAt:   time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),
			Title:       "Post 10",
			Description: "Foobarbaz",
			MarkdownUrl: mdUrl,
		},
	}

	tests := []struct {
		name                   string
		urlQuery               string
		errorToReturn          error
		postsToReturn          []int
		expectedStatusCode     int
		expectedJsonResponse   string
		expectedRequestedLimit int
		expectedRequestedFrom  time.Time
	}{
		{
			name:                   "No parameters",
			urlQuery:               userIdString + "/posts",
			errorToReturn:          nil,
			postsToReturn:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481985,"user_id":1796290045997481984,"created_at":"2024-04-04T00:00:00Z","title":"Post 0","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481986,"user_id":1796290045997481984,"created_at":"2024-04-05T00:00:00Z","title":"Post 1","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481987,"user_id":1796290045997481984,"created_at":"2024-04-06T00:00:00Z","title":"Post 2","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481988,"user_id":1796290045997481984,"created_at":"2024-04-07T00:00:00Z","title":"Post 3","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481989,"user_id":1796290045997481984,"created_at":"2024-04-08T00:00:00Z","title":"Post 4","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481990,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 5","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481991,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 6","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481992,"user_id":1796290045997481984,"created_at":"2024-04-10T00:00:00Z","title":"Post 7","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481993,"user_id":1796290045997481984,"created_at":"2024-05-04T00:00:00Z","title":"Post 8","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481994,"user_id":1796290045997481984,"created_at":"2024-06-04T00:00:00Z","title":"Post 9","markdown_url":"https://example.com","description":"Foobarbaz"}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
		{
			name:                   "From date, limit not specified",
			urlQuery:               userIdString + "/posts?from=2024-04-07T00%3A00%3A00Z",
			errorToReturn:          nil,
			postsToReturn:          []int{3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481988,"user_id":1796290045997481984,"created_at":"2024-04-07T00:00:00Z","title":"Post 3","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481989,"user_id":1796290045997481984,"created_at":"2024-04-08T00:00:00Z","title":"Post 4","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481990,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 5","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481991,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 6","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481992,"user_id":1796290045997481984,"created_at":"2024-04-10T00:00:00Z","title":"Post 7","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481993,"user_id":1796290045997481984,"created_at":"2024-05-04T00:00:00Z","title":"Post 8","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481994,"user_id":1796290045997481984,"created_at":"2024-06-04T00:00:00Z","title":"Post 9","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481995,"user_id":1796290045997481984,"created_at":"2024-07-04T00:00:00Z","title":"Post 10","markdown_url":"https://example.com","description":"Foobarbaz"}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified",
			urlQuery:               userIdString + "/posts?from=2024-04-07T00%3A00%3A00Z&limit=2",
			errorToReturn:          nil,
			postsToReturn:          []int{3, 4},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481988,"user_id":1796290045997481984,"created_at":"2024-04-07T00:00:00Z","title":"Post 3","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481989,"user_id":1796290045997481984,"created_at":"2024-04-08T00:00:00Z","title":"Post 4","markdown_url":"https://example.com","description":"Foobarbaz"}]`,
			expectedRequestedLimit: 2,
			expectedRequestedFrom:  time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than number of posts",
			urlQuery:               userIdString + "/posts?from=2024-02-06T00%3A00%3A00Z&limit=100",
			errorToReturn:          nil,
			postsToReturn:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481985,"user_id":1796290045997481984,"created_at":"2024-04-04T00:00:00Z","title":"Post 0","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481986,"user_id":1796290045997481984,"created_at":"2024-04-05T00:00:00Z","title":"Post 1","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481987,"user_id":1796290045997481984,"created_at":"2024-04-06T00:00:00Z","title":"Post 2","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481988,"user_id":1796290045997481984,"created_at":"2024-04-07T00:00:00Z","title":"Post 3","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481989,"user_id":1796290045997481984,"created_at":"2024-04-08T00:00:00Z","title":"Post 4","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481990,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 5","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481991,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 6","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481992,"user_id":1796290045997481984,"created_at":"2024-04-10T00:00:00Z","title":"Post 7","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481993,"user_id":1796290045997481984,"created_at":"2024-05-04T00:00:00Z","title":"Post 8","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481994,"user_id":1796290045997481984,"created_at":"2024-06-04T00:00:00Z","title":"Post 9","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481995,"user_id":1796290045997481984,"created_at":"2024-07-04T00:00:00Z","title":"Post 10","markdown_url":"https://example.com","description":"Foobarbaz"}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than 100",
			urlQuery:               userIdString + "/posts?from=2024-02-06T00%3A00%3A00Z&limit=200",
			errorToReturn:          nil,
			postsToReturn:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481985,"user_id":1796290045997481984,"created_at":"2024-04-04T00:00:00Z","title":"Post 0","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481986,"user_id":1796290045997481984,"created_at":"2024-04-05T00:00:00Z","title":"Post 1","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481987,"user_id":1796290045997481984,"created_at":"2024-04-06T00:00:00Z","title":"Post 2","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481988,"user_id":1796290045997481984,"created_at":"2024-04-07T00:00:00Z","title":"Post 3","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481989,"user_id":1796290045997481984,"created_at":"2024-04-08T00:00:00Z","title":"Post 4","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481990,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 5","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481991,"user_id":1796290045997481984,"created_at":"2024-04-09T00:00:00Z","title":"Post 6","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481992,"user_id":1796290045997481984,"created_at":"2024-04-10T00:00:00Z","title":"Post 7","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481993,"user_id":1796290045997481984,"created_at":"2024-05-04T00:00:00Z","title":"Post 8","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481994,"user_id":1796290045997481984,"created_at":"2024-06-04T00:00:00Z","title":"Post 9","markdown_url":"https://example.com","description":"Foobarbaz"},{"id":1796290045997481995,"user_id":1796290045997481984,"created_at":"2024-07-04T00:00:00Z","title":"Post 10","markdown_url":"https://example.com","description":"Foobarbaz"}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                 "Non number post id string",
			urlQuery:             "hello/posts",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: "Url parameter could not be parsed properly.\n",
		},
		{
			name:                 "Date badly formed",
			urlQuery:             userIdString + "/posts?from=202406T00%3A00%3A00Z",
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedJsonResponse: "Url parameter could not be parsed properly.\n",
		},
		{
			name:                 "Limit non integer",
			urlQuery:             userIdString + "/posts?limit=hello",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: "Url parameter could not be parsed properly.\n",
		},
		{
			name:                   "DB error",
			urlQuery:               userIdString + "/posts?",
			errorToReturn:          fmt.Errorf("database error"),
			expectedStatusCode:     http.StatusInternalServerError,
			expectedJsonResponse:   "internal server error\n",
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
		{
			name:                   "Unknown User",
			urlQuery:               "1/posts",
			errorToReturn:          sql.ErrNoRows,
			expectedStatusCode:     http.StatusNotFound,
			expectedJsonResponse:   "entity not found\n",
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
	}

	var userObj = storage.GetUserByIDRow{
		ID:       userIdNum,
		Username: "johndoe",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockedQueries := &storage.StubbedQueries{
				GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
					if id != userIdNum {
						return storage.GetUserByIDRow{}, sql.ErrNoRows
					}

					return userObj, nil
				},
				GetUserAndPostsByIDPagedFn: func(ctx context.Context, arg storage.GetUserAndPostsByIDPagedParams) ([]storage.GetUserAndPostsByIDPagedRow, error) {
					if int(arg.Limit) != tt.expectedRequestedLimit {
						t.Fatal("Expected limit to be", tt.expectedRequestedLimit, "but got", arg.Limit)
					}

					if arg.CreatedAfter != tt.expectedRequestedFrom {
						t.Fatal("Expected from to be", tt.expectedRequestedFrom, "but got", arg.CreatedAfter)
					}

					if tt.errorToReturn != nil {
						return []storage.GetUserAndPostsByIDPagedRow{}, tt.errorToReturn
					}

					rows := make([]storage.GetUserAndPostsByIDPagedRow, len(tt.postsToReturn))
					for i, idx := range tt.postsToReturn {
						rows[i] = storage.GetUserAndPostsByIDPagedRow{
							ID:       userObj.ID,
							Username: userObj.Username,
							Post:     expectedPosts[idx],
						}
					}
					return rows, nil
				},
			}

			req := httptest.NewRequest("GET", "/user/"+tt.urlQuery, nil)

			rr := httptest.NewRecorder()
			handler := NewHandler(
				slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
				mockedQueries,
				&stubbedDB{ShouldReturnError: false},
			)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

			if rr.Body.String() != tt.expectedJsonResponse {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedJsonResponse)
			}
		})
	}
}

func TestGetCommentsForPost(t *testing.T) {
	const postIdNum = 1796290045997481984
	const postIdString = "1796290045997481984"
	const userIdNum = 1796290045997481985
	const baseCommentNum = 1796290045997481986

	expectedComments := []storage.Comment{
		{
			ID:        baseCommentNum,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 1",
		},
		{
			ID:        baseCommentNum + 1,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 2",
		},
		{
			ID:        baseCommentNum + 2,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 5, 20, 0, 0, 0, time.UTC),
			Content:   "Example Comment 3",
		},
		{
			ID:        baseCommentNum + 3,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 4",
		},
		{
			ID:        baseCommentNum + 4,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 5",
		},
		{
			ID:        baseCommentNum + 5,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 6",
		},
		{
			ID:        baseCommentNum + 6,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 7",
		},
		{
			ID:        baseCommentNum + 7,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 8",
		},
		{
			ID:        baseCommentNum + 8,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 5, 7, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 9",
		},
		{
			ID:        baseCommentNum + 9,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 10",
		},
		{
			ID:        baseCommentNum + 10,
			PostID:    postIdNum,
			UserID:    userIdNum,
			CreatedAt: time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),
			Content:   "Example Comment 11",
		},
	}

	tests := []struct {
		name                    string
		urlQuery                string
		errorToReturn           error
		commentsToReturn        []int
		expectedStatusCode      int
		expectedJsonResponse    string
		expectedRequestedOffset int32
	}{
		{
			name:                    "No parameters",
			urlQuery:                postIdString + "/comments",
			errorToReturn:           nil,
			commentsToReturn:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedStatusCode:      http.StatusOK,
			expectedJsonResponse:    `{"TotalPages":0,"Comments":[{"id":1796290045997481986,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-04T00:00:00Z","content":"Example Comment 1"},{"id":1796290045997481987,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-05T00:00:00Z","content":"Example Comment 2"},{"id":1796290045997481988,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-05T20:00:00Z","content":"Example Comment 3"},{"id":1796290045997481989,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-08T00:00:00Z","content":"Example Comment 10"}]}`,
			expectedRequestedOffset: 0,
		},
		{
			name:                    "Page specified",
			urlQuery:                postIdString + "/comments?page=2",
			errorToReturn:           nil,
			commentsToReturn:        []int{3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:      http.StatusOK,
			expectedJsonResponse:    `{"TotalPages":0,"Comments":[{"id":1796290045997481989,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-08T00:00:00Z","content":"Example Comment 10"},{"id":1796290045997481996,"post_id":1796290045997481984,"user_id":1796290045997481985,"username":"johndoe","created_at":"2024-05-09T00:00:00Z","content":"Example Comment 11"}]}`,
			expectedRequestedOffset: 10,
		},
		{
			name:                 "Non number post id string",
			urlQuery:             "hello/comments",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: "Url parameter could not be parsed properly.\n",
		},
		{
			name:                 "Page non integer",
			urlQuery:             postIdString + "/comments?page=hello",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: "Url parameter could not be parsed properly.\n",
		},
		{
			name:                    "DB error",
			urlQuery:                postIdString + "/comments?",
			errorToReturn:           fmt.Errorf("database error"),
			expectedStatusCode:      http.StatusInternalServerError,
			expectedJsonResponse:    "internal server error\n",
			expectedRequestedOffset: 0,
		},
	}

	var userObj = storage.GetUserByIDRow{
		ID:       userIdNum,
		Username: "johndoe",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockedQueries := &storage.StubbedQueries{
				FindCommentsAndUserByPostIDPagedFn: func(ctx context.Context, arg storage.FindCommentsAndUserByPostIDPagedParams) ([]storage.FindCommentsAndUserByPostIDPagedRow, error) {
					if arg.PostID != postIdNum {
						t.Fatal("Expected post id to be", postIdNum, "but got", arg.PostID)
					}
					if arg.Offset != tt.expectedRequestedOffset {
						t.Fatal("Expected offset to be", tt.expectedRequestedOffset, "but got", arg.Offset)
					}

					if tt.errorToReturn != nil {
						return nil, tt.errorToReturn
					}

					rows := make([]storage.FindCommentsAndUserByPostIDPagedRow, len(tt.commentsToReturn))
					for i, idx := range tt.commentsToReturn {
						rows[i] = storage.FindCommentsAndUserByPostIDPagedRow{
							ID:       userObj.ID,
							Username: userObj.Username,
							Comment:  expectedComments[idx],
						}
					}
					return rows, nil
				},
			}

			req := httptest.NewRequest("GET", "/post/"+tt.urlQuery, nil)

			rr := httptest.NewRecorder()
			handler := NewHandler(
				slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
				mockedQueries,
				&stubbedDB{ShouldReturnError: false},
			)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatusCode)
			}

			if rr.Body.String() != tt.expectedJsonResponse {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedJsonResponse)
			}
		})
	}
}
