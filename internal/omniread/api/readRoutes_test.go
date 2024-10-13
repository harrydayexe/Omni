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

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

type mockPostRepo struct {
	readFunc func(ctx context.Context, id snowflake.Snowflake) (*models.Post, error)
}

func (m mockPostRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
	return m.readFunc(ctx, id)
}
func (m mockPostRepo) Create(ctx context.Context, entity models.Post) error {
	return errors.New("not implemented")
}
func (m mockPostRepo) Update(ctx context.Context, entity models.Post) error {
	return errors.New("not implemented")
}
func (m mockPostRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	return errors.New("not implemented")
}

type mockUserRepo struct {
	readFunc func(ctx context.Context, id snowflake.Snowflake) (*models.User, error)
}

func (m mockUserRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
	return m.readFunc(ctx, id)
}
func (m mockUserRepo) Create(ctx context.Context, entity models.User) error {
	return errors.New("not implemented")
}
func (m mockUserRepo) Update(ctx context.Context, entity models.User) error {
	return errors.New("not implemented")
}
func (m mockUserRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	return errors.New("not implemented")
}

type mockCommentRepo struct {
	getCommentsForPostFunc func(ctx context.Context, postId snowflake.Snowflake, from time.Time, limit int) ([]models.Comment, error)
}

func (m mockCommentRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.Comment, error) {
	return nil, errors.New("not implemented")
}
func (m mockCommentRepo) Create(ctx context.Context, entity models.Comment) error {
	return errors.New("not implemented")
}
func (m mockCommentRepo) Update(ctx context.Context, entity models.Comment) error {
	return errors.New("not implemented")
}
func (m mockCommentRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	return errors.New("not implemented")
}
func (m mockCommentRepo) GetCommentsForPost(ctx context.Context, postId snowflake.Snowflake, from time.Time, limit int) ([]models.Comment, error) {
	return m.getCommentsForPostFunc(ctx, postId, from, limit)
}

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
		name                 string
		urlQuery             string
		commentsToReturn     []int
		expectedStatusCode   int
		expectedJsonResponse string
	}{
		{
			name:                 "No parameters",
			urlQuery:             "",
			commentsToReturn:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedStatusCode:   http.StatusOK,
			expectedJsonResponse: `[{"id":1796290045997481986,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-04T00:00:00Z","content":"Example Comment 1"},{"id":1796290045997481987,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T00:00:00Z","content":"Example Comment 2"},{"id":1796290045997481988,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-05T20:00:00Z","content":"Example Comment 3"},{"id":1796290045997481989,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-06T00:00:00Z","content":"Example Comment 4"},{"id":1796290045997481990,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-07T00:00:00Z","content":"Example Comment 5"},{"id":1796290045997481991,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-08T00:00:00Z","content":"Example Comment 6"},{"id":1796290045997481992,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-04-09T00:00:00Z","content":"Example Comment 7"},{"id":1796290045997481993,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-06T00:00:00Z","content":"Example Comment 8"},{"id":1796290045997481994,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-07T00:00:00Z","content":"Example Comment 9"},{"id":1796290045997481995,"postId":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2024-05-08T00:00:00Z","content":"Example Comment 10"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockedRepo := &mockCommentRepo{
				getCommentsForPostFunc: func(ctx context.Context, postId snowflake.Snowflake, from time.Time, limit int) ([]models.Comment, error) {
					comments := make([]models.Comment, len(tt.commentsToReturn))
					for i, idx := range tt.commentsToReturn {
						comments[i] = expectedComments[idx]
					}
					return comments, nil
				},
			}

			req, err := http.NewRequest("GET", "/post/"+postIdString+"/comments?"+tt.urlQuery, nil)
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
