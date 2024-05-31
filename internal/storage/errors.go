package storage

import (
	"fmt"
)

// DatabaseError is a general error type that is returned for unknown database
// errors.
// message is the error message
// err is the underlying error
type DatabaseError struct {
	message string
	err     error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err.Error())
}

// NewDatabaseError creates a new DatabaseError instance.
func NewDatabaseError(message string, err error) error {
	return &DatabaseError{message: message, err: err}
}
