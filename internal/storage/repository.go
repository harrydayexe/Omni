// Package storage provides interfaces and implementations for data storage
// and access.
package storage

import "github.com/harrydayexe/Omni/internal/snowflake"

// Repository is an interface for accessing persisted data.
// It provides methods for standard CRUD operations.
//
// Repository has an associated type T which must conform to Identifier.
type Repository[T snowflake.Identifier] interface {
	// Create adds entity to the database.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Create(entity T) error

	// Update modifies an entry to the database.
	// It matches based on the Id() of entity.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Update(entity T) error

	// Delete removes an entry from the database.
	// It matches based on the Id() of entity.
	// The returned error is nil if the operation is successful, otherwise it
	// contains the error that occurred.
	Delete(entity T) error
}
