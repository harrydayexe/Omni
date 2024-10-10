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
	readFunc   func(ctx context.Context, id snowflake.Snowflake) (*models.User, error)
	createFunc func(ctx context.Context, entity models.User) error
	updateFunc func(ctx context.Context, entity models.User) error
	deleteFunc func(ctx context.Context, id snowflake.Snowflake) error
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
