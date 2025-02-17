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

	"github.com/harrydayexe/Omni/internal/auth"
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

	mockedAuthService := auth.StubbedAuthService{
		SignupFn: func(ctx context.Context, password string) ([]byte, error) {
			return []byte("hashed_password"), nil
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

	if strings.HasPrefix(rr.Header().Get("Location"), "test.com/api/user/") {
		t.Errorf("handler did not return Location header")
	}
}

func TestInsertUserNoPassword(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{}

	mockedAuthService := auth.StubbedAuthService{
		SignupFn: func(ctx context.Context, password string) ([]byte, error) {
			return nil, auth.ErrPasswordTooShort
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

func TestInsertUserPassTooLong(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateUserFn: func(ctx context.Context, arg storage.CreateUserParams) error {
			return nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		SignupFn: func(ctx context.Context, password string) ([]byte, error) {
			return nil, auth.ErrPasswordTooLong
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

func TestInsertUserBadFormedJsonRequest(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateUserFn: func(ctx context.Context, arg storage.CreateUserParams) error {
			return nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		SignupFn: func(ctx context.Context, password string) ([]byte, error) {
			return []byte("hashed_password"), nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

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

func TestInsertUserDatabaseError(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		CreateUserFn: func(ctx context.Context, arg storage.CreateUserParams) error {
			return fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		SignupFn: func(ctx context.Context, password string) ([]byte, error) {
			return []byte("hashed_password"), nil
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
		mockedAuthService,
		snowflakeGenerator,
		config,
	)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
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
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
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

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))
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

	if rr.Header().Get("Location") != "test.com:80/api/user/1796290045997481984" {
		t.Errorf("handler did not return Location header, got %v, want %v", rr.Header().Get("Location"), "test.com:80/api/user/1796290045997481984")
	}

	expected := `{"username":"johndoe","id":1796290045997481984}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, sql.ErrNoRows
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
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

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))
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

func TestUpdateUserDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
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

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))
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

func TestUpdateUserDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return fmt.Errorf("database error")
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
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

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))
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

	expected := "failed to update user\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateUserInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
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

	req := httptest.NewRequest("PUT", "/user/hello", bytes.NewBuffer(jsonBody))
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

func TestUpdateUserInvalidJsonBody(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		UpdateUserFn: func(ctx context.Context, arg storage.UpdateUserParams) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	jsonBody := []byte(`{"username"ohndoe"}`)

	req := httptest.NewRequest("PUT", "/user/1796290045997481984", bytes.NewBuffer(jsonBody))
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

func TestDeleteUserValid(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteUserFn: func(ctx context.Context, id int64) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

func TestDeleteUserNotFound(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteUserFn: func(ctx context.Context, id int64) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, sql.ErrNoRows
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

func TestDeleteUserDbErrorOnRead(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteUserFn: func(ctx context.Context, id int64) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{}, fmt.Errorf("database error")
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

func TestDeleteUserDbErrorOnWrite(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteUserFn: func(ctx context.Context, id int64) error {
			return fmt.Errorf("database error")
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

	expected := "failed to delete user\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteUserInvalidId(t *testing.T) {
	mockedQueries := &storage.StubbedQueries{
		DeleteUserFn: func(ctx context.Context, id int64) error {
			return nil
		},
		GetUserByIDFn: func(ctx context.Context, id int64) (storage.GetUserByIDRow, error) {
			return storage.GetUserByIDRow{
				ID:       1796290045997481984,
				Username: "tester",
			}, nil
		},
	}

	mockedAuthService := auth.StubbedAuthService{
		VerifyTokenFn: func(ctx context.Context, token string, id snowflake.Identifier) error {
			return nil
		},
	}

	req := httptest.NewRequest("DELETE", "/user/hello", nil)
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
