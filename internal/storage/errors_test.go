package storage

import (
	"errors"
	"testing"
)

func TestDatabaseError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := NewDatabaseError("test error", wrappedErr)
	if err == nil {
		t.Errorf("NewDatabaseError returned nil")
	}
	expected := "test error: wrapped error"
	if err.Error() != expected {
		t.Errorf(
			"NewDatabaseError returned error with message '%s', expected '%s'",
			err.Error(),
			expected,
		)
	}
}
