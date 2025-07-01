// Package storage provides interfaces for persisting and retrieving SBOM data.
// This package defines the repository pattern for data access in our hexagonal architecture.
package storage

import (
	"context"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
)

// Repository defines the contract for storing and retrieving SBOM documents.
// Implementations of this interface handle the persistence layer details
// while keeping the core business logic database-agnostic.
type Repository interface {
	// Store persists an SBOM document to the underlying storage system.
	// Returns an error if the SBOM cannot be stored.
	Store(ctx context.Context, sbom core.SBOM) error
	
	// FindByID retrieves an SBOM document by its unique identifier.
	// Returns nil and no error if the SBOM is not found.
	// Returns an error if there's a problem accessing the storage system.
	FindByID(ctx context.Context, id string) (*core.SBOM, error)
}