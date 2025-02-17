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

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

func TestInsertCommentValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateCommentFn: func(ctx context.Context, arg storage.CreateCommentParams) error {
			return nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"user_id":    1796290045997481984,
		"content":    "test comment",
		"created_at": "2025-01-01T02:30:00Z",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/post/1796290045997481985/comment", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

	if strings.HasPrefix(rr.Header().Get("Location"), "test.com/api/post/1796290045997481985/comments") {
		t.Errorf("handler did not return Location header")
	}
}

func TestInsertCommentUnauthorised(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return auth.ErrUnauthorized
		},
	}

	requestBody := map[string]interface{}{
		"user_id":    1796290045997481984,
		"content":    "test comment",
		"created_at": "2025-01-01T02:30:00Z",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/post/1796290045997481985/comment", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$HHHHHHHHsuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}
}

func TestInsertCommentBadFormedJsonRequest(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateCommentFn: func(ctx context.Context, arg storage.CreateCommentParams) error {
			return nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

	req := httptest.NewRequest("POST", "/post/1796290045997481985/comment", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestInsertCommentDatabaseError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateCommentFn: func(ctx context.Context, arg storage.CreateCommentParams) error {
			return fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"user_id":    1796290045997481984,
		"content":    "test comment",
		"created_at": "2025-01-01T02:30:00Z",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/post/1796290045997481985/comment", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if rr.Body.String() != "failed to create comment\n" {
		t.Errorf("handler did not return expected error message")
	}
}

func TestUpdateCommentUnauthorised(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return auth.ErrUnauthorized
		},
	}

	requestBody := map[string]interface{}{
		"content": "test updated comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CHHHHHHH4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}
}

func TestUpdateCommentValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateCommentFn: func(ctx context.Context, arg storage.UpdateCommentParams) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"content": "test updated comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

	expectedLocation := "test.com:80/api/post/1796290045997481985/comments"
	if rr.Header().Get("Location") != expectedLocation {
		t.Errorf("handler did not return Location header, got %v, want %v", rr.Header().Get("Location"), expectedLocation)
	}

	expected := `{"id":1796290045997481986,"post_id":1796290045997481985,"user_id":1796290045997481984,"content":"test updated comment","created_at":"2025-01-01T02:30:00Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateCommentNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateCommentFn: func(ctx context.Context, arg storage.UpdateCommentParams) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{}, sql.ErrNoRows
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"content": "test comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestUpdateCommentDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateCommentFn: func(ctx context.Context, arg storage.UpdateCommentParams) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{}, fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"content": "test comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestUpdateCommentDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateCommentFn: func(ctx context.Context, arg storage.UpdateCommentParams) error {
			return fmt.Errorf("database error")
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"content": "test comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

	expected := "failed to update comment\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateCommentInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	requestBody := map[string]interface{}{
		"content": "test comment",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/comment/hello", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestUpdateCommentInvalidJsonBody(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

	req := httptest.NewRequest("PUT", "/comment/1796290045997481986", bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestDeleteCommentUnauthorised(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return auth.ErrUnauthorized
		},
	}

	req := httptest.NewRequest("DELETE", "/comment/1796290045997481986", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CXXXXXXXXXXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}
}

func TestDeleteCommentValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteCommentFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/comment/1796290045997481986", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestDeleteCommentNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteCommentFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{}, sql.ErrNoRows
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/comment/1796290045997481986", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestDeleteCommentDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteCommentFn: func(ctx context.Context, id int64) error {
			return nil
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{}, fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/comment/1796290045997481986", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

func TestDeleteCommentDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteCommentFn: func(ctx context.Context, id int64) error {
			return fmt.Errorf("database error")
		},
		FindCommentAndUserByIDFn: func(ctx context.Context, id int64) (storage.FindCommentAndUserByIDRow, error) {
			return storage.FindCommentAndUserByIDRow{
				Comment: storage.Comment{
					ID:        1796290045997481986,
					PostID:    1796290045997481985,
					UserID:    1796290045997481984,
					Content:   "test comment",
					CreatedAt: time.Date(2025, 01, 01, 02, 30, 00, 00, time.UTC),
				},
				ID:       1796290045997481984,
				Username: "testuser",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/comment/1796290045997481986", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

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
		mockedAuthService,
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

	expected := "failed to delete comment\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteCommentInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{}

	req := httptest.NewRequest("DELETE", "/comment/hello", nil)
	req.Header.Add("Authorization", "Bearer $2a$10$L00CK5Aasuv4UXgXH36hj.xG00iiuDWTza1O8hiC7MdoBsKkDNm9y")

	snowflakeGenerator := snowflake.NewSnowflakeGenerator(0)
	config := &config.Config{
		Host: "test.com",
		Port: 80,
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
		mockedQueries,
		&stubbedDB{},
		mockedAuthService,
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
