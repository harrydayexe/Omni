package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

func TestLogin(t *testing.T) {
	var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	var cases = []struct {
		name         string
		loginFn      func(ctx context.Context, id snowflake.Identifier, password string) (string, error)
		id           uint64
		password     string
		expectedCode int
		expectedBody string
	}{
		{
			name: "valid login",
			loginFn: func(ctx context.Context, id snowflake.Identifier, password string) (string, error) {
				return "$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6", nil
			},
			id:           1796290045997481984,
			password:     "password",
			expectedCode: http.StatusOK,
			expectedBody: `{"access_token":"$2a$10$RV8G09OWcyqjj6n0S/OZaegrth8X24p5ai/pQMbjZlr.v9iu5QKT6","token_type":"Bearer","expires_in":86400}`,
		},
		{
			name: "invalid login",
			loginFn: func(ctx context.Context, id snowflake.Identifier, password string) (string, error) {
				return "", fmt.Errorf("invalid login")
			},
			id:           1796290045997481984,
			password:     "invalid",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var authService = auth.StubbedAuthService{
				LoginFn: tc.loginFn,
			}

			reqBody := map[string]interface{}{
				"id":       tc.id,
				"password": tc.password,
			}
			jsonBody, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))

			rr := httptest.NewRecorder()
			handler := NewHandler(testLogger, authService, nil)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			if rr.Body.String() != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tc.expectedBody)
			}
		})
	}
}
