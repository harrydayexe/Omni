package storage

import (
	"fmt"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// CouldNotFindEntityError is an error type that is returned when an entity
// could not be retrieved from the database.
// entityType is the type of the entity we are searching for
// id is the id of the entity we are searching for
type CouldNotFindEntityError struct {
	entityType string
	id         snowflake.Snowflake
}

func (e *CouldNotFindEntityError) Error() string {
	return fmt.Sprintf("could not find an entity %s with id %d", e.entityType, e.id.ToInt())
}

// NewCouldNotFindEntityError creates a new CouldNotFindEntityError instance.
func NewCouldNotFindEntityError(entityType string, id snowflake.Snowflake) error {
	return &CouldNotFindEntityError{entityType: entityType, id: id}
}

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
