package utilities

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubbedDB struct {
	ShouldReturnError bool
}

func (s *StubbedDB) PingContext(ctx context.Context) error {
	if s.ShouldReturnError {
		return errors.New("ping error")
	}
	return nil
}

func TestHealthCheck(t *testing.T) {
	var cases = []struct {
		name              string
		shouldReturnError bool
		expectedStatus    int
	}{
		{
			name:              "should return 200",
			shouldReturnError: false,
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "should return 500",
			shouldReturnError: true,
			expectedStatus:    http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := &StubbedDB{ShouldReturnError: tc.shouldReturnError}
			mux := http.NewServeMux()
			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
			AddHealthCheck(mux, logger, db)

			req := httptest.NewRequest("GET", "/healthz", nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}
		})
	}
}
