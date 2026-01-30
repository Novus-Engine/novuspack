// This file defines a minimal PathMetadataEntryRef interface to avoid circular
// dependencies between packages. It provides a reference interface for
// PathMetadataEntry structures used in generic code. This file should contain
// only the interface definition for breaking circular dependencies.
//
// Specification: api_generics.md: 1. Core Generic Types

package generics

// PathMetadataEntryRef provides a minimal interface for PathMetadataEntry references
// to avoid circular dependencies between fileformat and metadata packages.
//
// This interface exposes only the methods needed for file entry operations.
// The actual PathMetadataEntry type in the metadata package implements this interface.
//
// Specification: api_generics.md: 1. Core Generic Types
type PathMetadataEntryRef interface {
	// GetPath returns the path string for this path metadata entry.
	//
	// Returns:
	//   - string: Path string
	GetPath() string

	// GetPathEntry returns the PathEntry for this path metadata entry.
	//
	// Returns:
	//   - PathEntry: Path entry
	GetPathEntry() PathEntry
}
