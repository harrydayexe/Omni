package storage

import (
	"fmt"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

type EntityType int

const (
	User EntityType = iota
	Post
	Comment
)

var entityTypeNameMap = map[EntityType]string{User: "User", Post: "Post", Comment: "Comment"}

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

// NotFoundError is an error type that is returned when an entity cannot be found
type NotFoundError struct {
	id     snowflake.Snowflake
	entity EntityType
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find %s with id: %d", entityTypeNameMap[e.entity], e.id.ToInt())
}

// NewNotFoundError creates a new NotFoundError instance.
func NewNotFoundError(entity EntityType, id snowflake.Snowflake) error {
	return &NotFoundError{
		id:     id,
		entity: entity,
	}
}

// EntityAlreadyExistsError is an error type that is returned when a query returns no rows.
// err is the underlying error
type EntityAlreadyExistsError struct {
	id snowflake.Snowflake
}

func (e *EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("any entity with id: %d, already exists", e.id.ToInt())
}

// NewEntityAlreadyExistsError creates a new EntityAlreadyExistsError instance.
func NewEntityAlreadyExistsError(id snowflake.Snowflake) error {
	return &EntityAlreadyExistsError{
		id: id,
	}
}

// RequiredEntityDoesNotExistError is an error type that is returned when there is something
// wrong with the data of a request. For example trying to create a comment from
// a user that does not exist.
type RequiredEntityDoesNotExistError struct {
	entity EntityType
	id     snowflake.Snowflake
}

func (e *RequiredEntityDoesNotExistError) Error() string {
	return fmt.Sprintf("%s with id: %d", entityTypeNameMap[e.entity], e.id.ToInt())
}

// NewRequiredEntityDoesNotExist creates a new RequiredEntityDoesNotExist instance.
func NewRequiredEntityDoesNotExist(entity EntityType, id snowflake.Snowflake) error {
	return &RequiredEntityDoesNotExistError{
		entity: entity,
		id:     id,
	}
}
