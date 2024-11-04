package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

func TestGetUserKnown(t *testing.T) {
	mockedRepo := &mockUserRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
			newUser := models.NewUser(id, "johndoe", make([]snowflake.Snowflake, 0))
			return &newUser, nil
		},
	}

	req, err := http.NewRequest("GET", "/user/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedRepo,
		nil,
		nil,
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

	expected := `{"id":1796290045997481984,"username":"johndoe","posts":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserUnknown(t *testing.T) {
	mockedRepo := &mockUserRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
			return nil, nil
		},
	}

	req, err := http.NewRequest("GET", "/user/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedRepo,
		nil,
		nil,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetUserBadFormedId(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		nil,
		nil,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserDBError(t *testing.T) {
	mockedRepo := &mockUserRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
			return nil, storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	req, err := http.NewRequest("GET", "/user/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedRepo,
		nil,
		nil,
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
	mockedRepo := &mockPostRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
			newPost := models.NewPost(
				id,
				snowflake.ParseId(1796290045997481985),
				"johndoe",
				expectedTime,
				"Hello, World!",
				"Foobarbaz",
				expectedURL,
				[]snowflake.Snowflake{snowflake.ParseId(1796290045997481986)},
				[]string{"foo", "bar", "baz"},
			)
			return &newPost, nil
		},
	}

	req, err := http.NewRequest("GET", "/post/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		mockedRepo,
		nil,
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

	expected := `{"id":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2021-01-01T11:40:35Z","title":"Hello, World!","description":"Foobarbaz","contentFileUrl":"https://example.com/foo","comments":[1796290045997481986],"tags":["foo","bar","baz"]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostUnknown(t *testing.T) {
	mockedRepo := &mockPostRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
			return nil, nil
		},
	}

	req, err := http.NewRequest("GET", "/post/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		mockedRepo,
		nil,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetPostBadFormedId(t *testing.T) {
	req, err := http.NewRequest("GET", "/post/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		nil,
		nil,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostDBError(t *testing.T) {
	mockedRepo := &mockPostRepo{
		readFunc: func(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
			return nil, storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	req, err := http.NewRequest("GET", "/post/1796290045997481984", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		nil,
		mockedRepo,
		nil,
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
	userId := snowflake.ParseId(userIdNum)

	expectedPosts := []models.Post{
		models.NewPost(
			snowflake.ParseId(basePostNum+0),
			userId,
			"johndoe",
			time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			"Post 0",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+1),
			userId,
			"johndoe",
			time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			"Post 1",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+2),
			userId,
			"johndoe",
			time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
			"Post 2",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+3),
			userId,
			"johndoe",
			time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
			"Post 3",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+4),
			userId,
			"johndoe",
			time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC),
			"Post 4",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+5),
			userId,
			"johndoe",
			time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			"Post 5",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+6),
			userId,
			"johndoe",
			time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			"Post 6",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+7),
			userId,
			"johndoe",
			time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC),
			"Post 7",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+8),
			userId,
			"johndoe",
			time.Date(2024, 5, 4, 0, 0, 0, 0, time.UTC),
			"Post 8",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+9),
			userId,
			"johndoe",
			time.Date(2024, 6, 4, 0, 0, 0, 0, time.UTC),
			"Post 9",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
		models.NewPost(
			snowflake.ParseId(basePostNum+10),
			userId,
			"johndoe",
			time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),
			"Post 10",
			"Foobarbaz",
			url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			[]snowflake.Snowflake{},
			[]string{},
		),
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
			expectedJsonResponse:   `[{"id":1796290045997481985,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","title":"Post 0","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481986,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","title":"Post 1","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481987,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","title":"Post 2","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481988,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","title":"Post 3","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481989,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","title":"Post 4","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481990,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 5","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481991,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 6","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481992,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-10T00:00:00Z","title":"Post 7","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481993,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-05-04T00:00:00Z","title":"Post 8","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481994,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-06-04T00:00:00Z","title":"Post 9","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
		{
			name:                   "From date, limit not specified",
			urlQuery:               userIdString + "/posts?from=2024-04-07T00%3A00%3A00Z",
			errorToReturn:          nil,
			postsToReturn:          []int{3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481988,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","title":"Post 3","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481989,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","title":"Post 4","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481990,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 5","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481991,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 6","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481992,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-10T00:00:00Z","title":"Post 7","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481993,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-05-04T00:00:00Z","title":"Post 8","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481994,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-06-04T00:00:00Z","title":"Post 9","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481995,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-07-04T00:00:00Z","title":"Post 10","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified",
			urlQuery:               userIdString + "/posts?from=2024-04-07T00%3A00%3A00Z&limit=2",
			errorToReturn:          nil,
			postsToReturn:          []int{3, 4},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481988,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","title":"Post 3","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481989,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","title":"Post 4","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]}]`,
			expectedRequestedLimit: 2,
			expectedRequestedFrom:  time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than number of posts",
			urlQuery:               userIdString + "/posts?from=2024-02-06T00%3A00%3A00Z&limit=100",
			errorToReturn:          nil,
			postsToReturn:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481985,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","title":"Post 0","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481986,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","title":"Post 1","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481987,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","title":"Post 2","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481988,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","title":"Post 3","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481989,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","title":"Post 4","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481990,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 5","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481991,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 6","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481992,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-10T00:00:00Z","title":"Post 7","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481993,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-05-04T00:00:00Z","title":"Post 8","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481994,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-06-04T00:00:00Z","title":"Post 9","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481995,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-07-04T00:00:00Z","title":"Post 10","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than 100",
			urlQuery:               userIdString + "/posts?from=2024-02-06T00%3A00%3A00Z&limit=200",
			errorToReturn:          nil,
			postsToReturn:          []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481985,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","title":"Post 0","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481986,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","title":"Post 1","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481987,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","title":"Post 2","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481988,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","title":"Post 3","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481989,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","title":"Post 4","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481990,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 5","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481991,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","title":"Post 6","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481992,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-04-10T00:00:00Z","title":"Post 7","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481993,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-05-04T00:00:00Z","title":"Post 8","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481994,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-06-04T00:00:00Z","title":"Post 9","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]},{"id":1796290045997481995,"authorId":1796290045997481984,"authorName":"johndoe","timestamp":"2024-07-04T00:00:00Z","title":"Post 10","description":"Foobarbaz","contentFileUrl":"https://example.com","comments":[],"tags":[]}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                 "Non number post id string",
			urlQuery:             "hello/posts",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                 "Date badly formed",
			urlQuery:             userIdString + "/posts?from=202406T00%3A00%3A00Z",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                 "Limit non integer",
			urlQuery:             userIdString + "/posts?limit=hello",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                   "DB error",
			urlQuery:               userIdString + "/posts?",
			errorToReturn:          storage.NewDatabaseError("database error", errors.New("database error")),
			expectedStatusCode:     http.StatusInternalServerError,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockedRepo := &mockPostRepo{
				getPostsForUser: func(ctx context.Context, userId snowflake.Snowflake, from time.Time, limit int) ([]models.Post, error) {
					if limit != tt.expectedRequestedLimit {
						t.Fatal("Expected limit to be", tt.expectedRequestedLimit, "but got", limit)
					}

					if from != tt.expectedRequestedFrom {
						t.Fatal("Expected from to be", tt.expectedRequestedFrom, "but got", from)
					}

					if tt.errorToReturn != nil {
						return nil, tt.errorToReturn
					}

					posts := make([]models.Post, len(tt.postsToReturn))
					for i, idx := range tt.postsToReturn {
						posts[i] = expectedPosts[idx]
					}
					return posts, nil

				},
			}

			req, err := http.NewRequest("GET", "/user/"+tt.urlQuery, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := NewHandler(
				slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
				nil,
				mockedRepo,
				nil,
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
	postId := snowflake.ParseId(postIdNum)
	userId := snowflake.ParseId(userIdNum)

	expectedComments := []models.Comment{
		models.NewComment(
			snowflake.ParseId(baseCommentNum),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			"Example Comment 1",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+1),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			"Example Comment 2",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+2),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 5, 20, 0, 0, 0, time.UTC),
			"Example Comment 3",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+3),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
			"Example Comment 4",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+4),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 7, 0, 0, 0, 0, time.UTC),
			"Example Comment 5",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+5),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC),
			"Example Comment 6",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+6),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			"Example Comment 7",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+7),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC),
			"Example Comment 8",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+8),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 5, 7, 0, 0, 0, 0, time.UTC),
			"Example Comment 9",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+9),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			"Example Comment 10",
		),
		models.NewComment(
			snowflake.ParseId(baseCommentNum+10),
			postId,
			userId,
			"johndoe",
			time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),
			"Example Comment 11",
		),
	}

	tests := []struct {
		name                   string
		urlQuery               string
		errorToReturn          error
		commentsToReturn       []int
		expectedStatusCode     int
		expectedJsonResponse   string
		expectedRequestedLimit int
		expectedRequestedFrom  time.Time
	}{
		{
			name:                   "No parameters",
			urlQuery:               postIdString + "/comments",
			errorToReturn:          nil,
			commentsToReturn:       []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481986,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","content":"Example Comment 1"},{"id":1796290045997481987,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","content":"Example Comment 2"},{"id":1796290045997481988,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T20:00:00Z","content":"Example Comment 3"},{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-08T00:00:00Z","content":"Example Comment 10"}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
		{
			name:                   "From date, limit not specified",
			urlQuery:               postIdString + "/comments?from=2024-04-06T00%3A00%3A00Z",
			errorToReturn:          nil,
			commentsToReturn:       []int{3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-08T00:00:00Z","content":"Example Comment 10"},{"id":1796290045997481996,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-09T00:00:00Z","content":"Example Comment 11"}]`,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified",
			urlQuery:               postIdString + "/comments?from=2024-04-06T00%3A00%3A00Z&limit=2",
			errorToReturn:          nil,
			commentsToReturn:       []int{3, 4},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"}]`,
			expectedRequestedLimit: 2,
			expectedRequestedFrom:  time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than number of comments",
			urlQuery:               postIdString + "/comments?from=2024-04-06T00%3A00%3A00Z&limit=100",
			errorToReturn:          nil,
			commentsToReturn:       []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481986,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","content":"Example Comment 1"},{"id":1796290045997481987,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","content":"Example Comment 2"},{"id":1796290045997481988,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T20:00:00Z","content":"Example Comment 3"},{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-08T00:00:00Z","content":"Example Comment 10"},{"id":1796290045997481996,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-09T00:00:00Z","content":"Example Comment 11"}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                   "From date, limit specified, limit is greater than 100",
			urlQuery:               postIdString + "/comments?from=2024-04-06T00%3A00%3A00Z&limit=200",
			errorToReturn:          nil,
			commentsToReturn:       []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedStatusCode:     http.StatusOK,
			expectedJsonResponse:   `[{"id":1796290045997481986,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","content":"Example Comment 1"},{"id":1796290045997481987,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","content":"Example Comment 2"},{"id":1796290045997481988,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T20:00:00Z","content":"Example Comment 3"},{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-08T00:00:00Z","content":"Example Comment 10"},{"id":1796290045997481996,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-09T00:00:00Z","content":"Example Comment 11"}]`,
			expectedRequestedLimit: 100,
			expectedRequestedFrom:  time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                 "Non number post id string",
			urlQuery:             "hello/comments",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                 "Date badly formed",
			urlQuery:             postIdString + "/comments?from=202406T00%3A00%3A00Z",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                 "Limit non integer",
			urlQuery:             postIdString + "/comments?limit=hello",
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`,
		},
		{
			name:                   "DB error",
			urlQuery:               postIdString + "/comments?",
			errorToReturn:          storage.NewDatabaseError("database error", errors.New("database error")),
			expectedStatusCode:     http.StatusInternalServerError,
			expectedRequestedLimit: 10,
			expectedRequestedFrom:  time.UnixMilli(1704067200000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockedRepo := &mockCommentRepo{
				getCommentsForPostFunc: func(ctx context.Context, postId snowflake.Snowflake, from time.Time, limit int) ([]models.Comment, error) {
					if limit != tt.expectedRequestedLimit {
						t.Fatal("Expected limit to be", tt.expectedRequestedLimit, "but got", limit)
					}

					if from != tt.expectedRequestedFrom {
						t.Fatal("Expected from to be", tt.expectedRequestedFrom, "but got", from)
					}

					if tt.errorToReturn != nil {
						return nil, tt.errorToReturn
					}

					comments := make([]models.Comment, len(tt.commentsToReturn))
					for i, idx := range tt.commentsToReturn {
						comments[i] = expectedComments[idx]
					}
					return comments, nil
				},
			}

			req, err := http.NewRequest("GET", "/post/"+tt.urlQuery, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := NewHandler(
				slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
				nil,
				nil,
				mockedRepo,
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
