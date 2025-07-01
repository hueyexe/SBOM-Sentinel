// Package ingestion provides interfaces and implementations for parsing SBOM documents
// from various formats into our core domain models.
package ingestion

import (
	"io"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
)

// Parser defines the contract for parsing SBOM documents from various formats.
// Implementations of this interface should handle specific SBOM formats like CycloneDX or SPDX.
type Parser interface {
	// Parse reads an SBOM document from the provided io.Reader and converts it
	// into our core SBOM domain model.
	// Returns an error if the document cannot be parsed or is invalid.
	Parse(r io.Reader) (*core.SBOM, error)
}