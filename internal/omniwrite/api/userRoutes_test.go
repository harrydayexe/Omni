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

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/snowflake"
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

func TestInsertUserValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateUserFn: func(ctx context.Context, arg storage.CreateUserParams) error {
			return nil
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonBody))

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

	if status := rr.Code; status != 201 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 201)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "application/json")
	}

	if strings.HasPrefix(rr.Header().Get("Location"), "test.com/user/") {
		t.Errorf("handler did not return Location header")
	}
}

func TestInsertUserDatabaseError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateUserFn: func(ctx context.Context, arg storage.CreateUserParams) error {
			return fmt.Errorf("database error")
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonBody))

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
			status, 201)
	}

	if rr.Body.String() != "failed to create user\n" {
		t.Errorf("handler did not return expected error message")
	}
}

func TestUpdateUserValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.User, error) {
			return storage.User{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))

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

	if rr.Header().Get("Location") != "test.com:80/user/1796290045997481984" {
		t.Errorf("handler did not return Location header, got %v, want %v", rr.Header().Get("Location"), "test.com/user/1796290045997481984")
	}

	expected := `{"id":1796290045997481984,"username":"johndoe"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.User, error) {
			return storage.User{}, sql.ErrNoRows
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))

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
			status, http.StatusOK)
	}

	if rr.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("handler did not return Content-Type header: got %v want %v", rr.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	}

	expected := "entity not found\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateUserDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.User, error) {
			return storage.User{}, fmt.Errorf("database error")
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))

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

func TestUpdateUserDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return fmt.Errorf("database error")
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.User, error) {
			return storage.User{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	requestBody := map[string]string{
		"username": "johndoe",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))

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

	expected := "failed to update user\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
