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
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/storage"
)

type mockPostRepo struct {
	readFunc   func(ctx context.Context, id snowflake.Snowflake) (*models.Post, error)
	createFunc func(ctx context.Context, entity models.Post) error
	updateFunc func(ctx context.Context, entity models.Post) error
	deleteFunc func(ctx context.Context, id snowflake.Snowflake) error
}

func (m mockPostRepo) Read(ctx context.Context, id snowflake.Snowflake) (*models.Post, error) {
	return m.readFunc(ctx, id)
}
func (m mockPostRepo) Create(ctx context.Context, entity models.Post) error {
	return m.createFunc(ctx, entity)
}
func (m mockPostRepo) Update(ctx context.Context, entity models.Post) error {
	return m.updateFunc(ctx, entity)
}
func (m mockPostRepo) Delete(ctx context.Context, id snowflake.Snowflake) error {
	return m.deleteFunc(ctx, id)
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

func TestDeletePostKnown(t *testing.T) {
	mockedRepo := &mockPostRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return nil
		},
	}

	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
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
}

func TestDeletePostUnknown(t *testing.T) {
	mockedRepo := &mockPostRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return storage.NewNotFoundError(storage.Post, snowflake.ParseId(1796290045997481984))
		},
	}

	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
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

func TestDeletePostBadFormedId(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/post/hello", nil)
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

func TestDeletePostDBError(t *testing.T) {
	mockedRepo := &mockPostRepo{
		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
			return storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
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

func TestCreatePostSuccess(t *testing.T) {
	mockedRepo := &mockPostRepo{
		createFunc: func(ctx context.Context, entity models.Post) error {
			return nil
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewHandler(
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		nil,
		mockedRepo,
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

	expected := `{"id":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2021-01-01T11:40:35Z","title":"Hello, World!","description":"Foobarbaz","contentFileUrl":"https://example.com/foo","comments":[],"tags":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreatePostDuplicate(t *testing.T) {
	mockedRepo := &mockPostRepo{
		createFunc: func(ctx context.Context, entity models.Post) error {
			return storage.NewEntityAlreadyExistsError(snowflake.ParseId(1796290045997481984))
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
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

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"error":"Conflict","message":"Post with that ID already exists."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreatePostBadFormedBody(t *testing.T) {
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
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
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

func TestCreatePostDBError(t *testing.T) {
	mockedRepo := &mockPostRepo{
		createFunc: func(ctx context.Context, entity models.Post) error {
			return storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
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

func TestUpdatePostSuccess(t *testing.T) {
	mockedRepo := &mockPostRepo{
		updateFunc: func(ctx context.Context, entity models.Post) error {
			return nil
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
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

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"id":1796290045997481984,"authorId":"1796290045997481985",authorName:"johndoe","timestamp":"2021-01-01T11:40:35Z","title":"Hello, World!","description":"Foobarbaz","contentFileUrl":"https://example.com/foo","comments":[],"tags":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdatePostNotFound(t *testing.T) {
	mockedRepo := &mockPostRepo{
		updateFunc: func(ctx context.Context, entity models.Post) error {
			return storage.NewNotFoundError(storage.Post, entity.Id())
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
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

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"error":"Not Found","message":"Post with that ID could not be found to update."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdatePostBadFormedBody(t *testing.T) {
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
	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
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

func TestUpdatePostDBError(t *testing.T) {
	mockedRepo := &mockPostRepo{
		updateFunc: func(ctx context.Context, entity models.Post) error {
			return storage.NewDatabaseError("database error", errors.New("database error"))
		},
	}

	body := struct {
		Id          uint64 `json:"id"`
		AuthorId    uint64 `json:"authorId"`
		AuthorName  string `json:"authorName"`
		Timestamp   string `json:"timestamp"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ContentFile string `json:"contentFileUrl"`
	}{
		Id:          1796290045997481984,
		AuthorId:    1796290045997481985,
		AuthorName:  "johndoe",
		Timestamp:   "2021-01-01T11:40:35Z",
		Title:       "Hello, World!",
		Description: "Foobarbaz",
		ContentFile: "https://example.com/foo",
	}
	out, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
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
