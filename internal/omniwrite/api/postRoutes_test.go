package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

func TestInsertPostValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreatePostFn: func(ctx context.Context, arg storage.CreatePostParams) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"user_id":      1796290045997481985,
		"created_at":   "2025-01-01T02:30:00Z",
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/post", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
	}

	if strings.HasPrefix(rr.Header().Get("Location"), "test.com/api/post/") {
		t.Errorf("handler did not return Location header")
	}
}

func TestInsertPostBadFormedJsonRequest(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreatePostFn: func(ctx context.Context, arg storage.CreatePostParams) error {
			return nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

	req := httptest.NewRequest("POST", "/post", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}
}

func TestInsertPostDatabaseError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreatePostFn: func(ctx context.Context, arg storage.CreatePostParams) error {
			return fmt.Errorf("database error")
		},
	}

	requestBody := map[string]interface{}{
		"user_id":      1796290045997481985,
		"created_at":   "2025-01-01T02:30:00Z",
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/post", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if rr.Body.String() != "failed to create post\n" {
		t.Errorf("handler did not return expected error message")
	}
}

func TestUpdatePostValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	requestBody := map[string]string{
		"title":        "test title updated",
		"description":  "test description updated",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
	}

	if rr.Header().Get("Location") != "test.com:80/api/post/1796290045997481984" {
		t.Errorf("handler did not return Location header, got %v, want %v", rr.Header().Get("Location"), "test.com:80/api/post/1796290045997481984")
	}

	expected := `{"id":1796290045997481984,"user_id":1796290045997481985,"created_at":"2025-01-01T02:30:00Z","title":"test title updated","description":"test description updated","markdown_url":"https://test.com/post1.md"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePostNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, sql.ErrNoRows
		},
	}

	requestBody := map[string]string{
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "entity not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePostDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, fmt.Errorf("database error")
		},
	}

	requestBody := map[string]string{
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "failed to read entity from db\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePostDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return fmt.Errorf("database error")
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	requestBody := map[string]string{
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "failed to update post\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePostInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	requestBody := map[string]string{
		"title":        "test title",
		"description":  "test description",
		"markdown_url": "https://test.com/post1.md",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/post/hello", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "Url parameter could not be parsed properly.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePostInvalidJsonBody(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdatePostFn: func(ctx context.Context, arg storage.UpdatePostParams) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

	req := httptest.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(jsonBody))

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}
}

func TestDeletePostValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeletePostFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	req := httptest.NewRequest("DELETE", "/post/1796290045997481984", nil)

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
	}
}

func TestDeletePostNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeletePostFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, sql.ErrNoRows
		},
	}

	req := httptest.NewRequest("DELETE", "/post/1796290045997481984", nil)

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "entity not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeletePostDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeletePostFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{}, fmt.Errorf("database error")
		},
	}

	req := httptest.NewRequest("DELETE", "/post/1796290045997481984", nil)

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "failed to read entity from db\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeletePostDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeletePostFn: func(ctx context.Context, id int64) error {
			return fmt.Errorf("database error")
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	req := httptest.NewRequest("DELETE", "/post/1796290045997481984", nil)

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "failed to delete post\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeletePostInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeletePostFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindPostByIDFn: func(ctx context.Context, id int64) (storage.Post, error) {
			return storage.Post{
				ID:          1796290045997481984,
				UserID:      1796290045997481985,
				CreatedAt:   time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				Title:       "test title",
				Description: "test description",
				MarkdownUrl: "https://test.com/post1.md",
			}, nil
		},
	}

	req := httptest.NewRequest("DELETE", "/post/hello", nil)

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "Url parameter could not be parsed properly.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
