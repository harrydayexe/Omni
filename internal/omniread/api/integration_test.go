package api

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

func TestGetUserAndGetPost(t *testing.T) {
	var cases = []struct {
		name                string
		path                string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "Get user by id success",
			path:                "/user/1796290045997481984",
			expectedCode:        200,
			expectedBody:        `{"id":1796290045997481984,"username":"johndoe"}`,
			expectedContentType: "application/json",
		},
		{
			name:                "Get user by id not found",
			path:                "/user/1",
			expectedCode:        404,
			expectedBody:        "entity not found\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name:                "Get post by id success",
			path:                "/post/1796290045997481995",
			expectedCode:        200,
			expectedBody:        `{"id":1796290045997481995,"user_id":1796290045997481984,"created_at":"2024-04-04T00:00:00Z","title":"My first post","markdown_url":"https://example.com/johndoe-first-post","description":"First post description"}`,
			expectedContentType: "application/json",
		},
		{
			name:                "Get post by id not found",
			path:                "/post/1",
			expectedCode:        404,
			expectedBody:        "entity not found\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	testLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run the test case in parallel

			db, cleanup, err := utilities.SetupTestContainer("../../../db/migrations/", "testdata.sql")
			if err != nil {
				t.Fatalf("failed to setup test container: %v", err)
			}
			defer cleanup()

			queries := storage.New(db)
			handler := NewHandler(testLogger, queries, db)

			req := httptest.NewRequest("GET", tc.path, nil)
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

func TestGetMostRecentPosts(t *testing.T) {
	var cases = []struct {
		name                string
		path                string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "Get most recent posts success",
			path:                "/posts",
			expectedCode:        200,
			expectedBody:        `[{"username":"harrydayexe","post":{"id":1796290045997482998,"user_id":1796290045997481986,"created_at":"2024-12-07T00:00:00Z","title":"My twelfth post","markdown_url":"https://example.com/harrydayexe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482997,"user_id":1796290045997481985,"created_at":"2024-12-06T00:00:00Z","title":"My eleventh post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482996,"user_id":1796290045997481985,"created_at":"2024-12-05T00:00:00Z","title":"My tenth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482995,"user_id":1796290045997481985,"created_at":"2024-12-04T00:00:00Z","title":"My ninth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482994,"user_id":1796290045997481985,"created_at":"2024-11-04T00:00:00Z","title":"My eighth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482993,"user_id":1796290045997481985,"created_at":"2024-10-04T00:00:00Z","title":"My seventh post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482992,"user_id":1796290045997481985,"created_at":"2024-09-04T00:00:00Z","title":"My sixth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482991,"user_id":1796290045997481985,"created_at":"2024-08-04T00:00:00Z","title":"My fifth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997482990,"user_id":1796290045997481985,"created_at":"2024-07-04T00:00:00Z","title":"My fourth post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997481999,"user_id":1796290045997481985,"created_at":"2024-06-04T02:00:00Z","title":"My third post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}}]`,
			expectedContentType: "application/json",
		},
		{
			name:                "Get second page results success",
			path:                "/posts?page=2",
			expectedCode:        200,
			expectedBody:        `[{"username":"janedoe","post":{"id":1796290045997481998,"user_id":1796290045997481985,"created_at":"2024-06-04T01:00:00Z","title":"My second post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"janedoe","post":{"id":1796290045997481997,"user_id":1796290045997481985,"created_at":"2024-06-04T00:00:00Z","title":"My first post","markdown_url":"https://example.com/janedoe-first-post","description":"First post description"}},{"username":"johndoe","post":{"id":1796290045997481996,"user_id":1796290045997481984,"created_at":"2024-05-04T00:00:00Z","title":"My second post","markdown_url":"https://example.com/johndoe-second-post","description":"Second post description"}},{"username":"johndoe","post":{"id":1796290045997481995,"user_id":1796290045997481984,"created_at":"2024-04-04T00:00:00Z","title":"My first post","markdown_url":"https://example.com/johndoe-first-post","description":"First post description"}}]`,
			expectedContentType: "application/json",
		},
		{
			name:                "Get third page results success",
			path:                "/posts?page=3",
			expectedCode:        200,
			expectedBody:        `[]`,
			expectedContentType: "application/json",
		},
	}

	testLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run the test case in parallel

			db, cleanup, err := utilities.SetupTestContainer("../../../db/migrations/", "testdata.sql")
			if err != nil {
				t.Fatalf("failed to setup test container: %v", err)
			}
			defer cleanup()

			queries := storage.New(db)
			handler := NewHandler(testLogger, queries, db)

			req := httptest.NewRequest("GET", tc.path, nil)
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

// TODO: Write integration tests for testing limits and pages etc
// func TestGetPostsOfUser(t *testing.T) {
// 	var cases = []struct {
// 		name                string
// 		path                string
// 		expectedCode        int
// 		expectedBody        string
// 		expectedContentType string
// 	}{}
// }
