package storage

import (
	"fmt"

	"github.com/harrydayexe/Omni/internal/snowflake"
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

// NotFoundError is an error type that is returned when a query returns no rows.
// err is the underlying error
type NotFoundError struct {
	id  snowflake.Snowflake
	err error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find entity with id: %d, error: %s", e.id.ToInt(), e.err.Error())
}

// NewNotFoundError creates a new NotFoundError instance.
func NewNotFoundError(id snowflake.Snowflake, err error) error {
	return &NotFoundError{
		id:  id,
		err: err,
	}
}
