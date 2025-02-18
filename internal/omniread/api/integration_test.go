package api

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func TestGetUser(t *testing.T) {
	var cases = []struct {
		name                string
		id                  int
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "Get user by id success",
			id:                  1796290045997481984,
			expectedCode:        200,
			expectedBody:        `{"id":1796290045997481984,"username":"johndoe"}`,
			expectedContentType: "application/json",
		},
		{
			name:                "Get user by id not found",
			id:                  1,
			expectedCode:        404,
			expectedBody:        "entity not found\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	testLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, cleanup, err := utilities.SetupTestContainer("../../../db/migrations/", "testdata.sql")
			if err != nil {
				t.Fatalf("failed to setup test container: %v", err)
			}
			defer cleanup()

			queries := storage.New(db)
			handler := NewHandler(testLogger, queries, db)

			path := "/user/" + strconv.Itoa(tc.id)
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			if body := rr.Body.String(); body != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tc.expectedBody)
			}

			if contentType := rr.Header().Get("Content-Type"); contentType != tc.expectedContentType {
				t.Errorf("handler returned unexpected content type: got %v want %v", contentType, tc.expectedContentType)
			}
		})
	}
}
