package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

type MockUserRepo struct {
	readFunc   func(ctx context.Context, id snowflake.Snowflake) (*models.User, error)
	createFunc func(ctx context.Context, entity models.User) error
	updateFunc func(ctx context.Context, entity models.User) error
	deleteFunc func(ctx context.Context, id snowflake.Snowflake) error
}

func (m MockUserRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.User, error) {
	return m.readFunc(ctx, id)
}
func (m MockUserRepo) Create(ctx context.Context, entity models.User) error {
	return m.createFunc(ctx, entity)
}
func (m MockUserRepo) Update(ctx context.Context, entity models.User) error {
	return m.updateFunc(ctx, entity)
}
func (m MockUserRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	return m.deleteFunc(ctx, id)
}

func TestGetUserKnown(t *testing.T) {
	mockedRepo := &MockUserRepo{
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
	mockedRepo := &MockUserRepo{
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
	mockedRepo := &MockUserRepo{
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

func TestDeleteUserKnown(t *testing.T) {
	mockedRepo := &MockUserRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return nil
		},
	}

	req, err := http.NewRequest("DELETE", "/user/1796290045997481984", nil)
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
}

func TestDeleteUserUnknown(t *testing.T) {
	mockedRepo := &MockUserRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return storage.NewNotFoundError(storage.User, snowflake.ParseId(1796290045997481984))
		},
	}

	req, err := http.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

func TestDeleteUserBadFormedId(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/user/hello", nil)
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

func TestDeleteUserDBError(t *testing.T) {
	mockedRepo := &MockUserRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	req, err := http.NewRequest("DELETE", "/user/1796290045997481984", nil)
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

func TestCreateUserSuccess(t *testing.T) {
	mockedRepo := &MockUserRepo{
		createFunc: func(ctx context.Context, entity models.User) error {
			return nil
		},
	}

	body := struct {
		Id       uint64 `json:"id"`
		Username string `json:"username"`
	}{
		Id:       1796290045997481984,
		Username: "johndoe",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(out))
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

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
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

func TestCreateUserDuplicate(t *testing.T) {
	mockedRepo := &MockUserRepo{
		createFunc: func(ctx context.Context, entity models.User) error {
			return storage.NewEntityAlreadyExistsError(snowflake.ParseId(1796290045997481984))
		},
	}

	body := struct {
		Id       uint64 `json:"id"`
		Username string `json:"username"`
	}{
		Id:       1796290045997481984,
		Username: "johndoe",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(out))
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

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"error":"Conflict","message":"User with that ID already exists."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserBadFormedBody(t *testing.T) {
	body := struct {
		Foo uint64 `json:"foo"`
		Bar string `json:"bar"`
		Baz string `json:"baz"`
	}{
		Foo: 1796290045997481984,
		Bar: "johndoe",
		Baz: "foobar",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(out))
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

	expected := `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserDBError(t *testing.T) {
	mockedRepo := &MockUserRepo{
		createFunc: func(ctx context.Context, entity models.User) error {
			return storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	body := struct {
		Id       uint64 `json:"id"`
		Username string `json:"username"`
	}{
		Id:       1796290045997481984,
		Username: "johndoe",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(out))
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
