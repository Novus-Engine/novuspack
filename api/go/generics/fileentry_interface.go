// This file defines a minimal FileEntryRef interface to avoid circular
// dependencies between packages. It provides a reference interface for
// FileEntry structures used in generic code. This file should contain only
// the interface definition for breaking circular dependencies.
//
// Specification: api_generics.md: 1. Core Generic Types

package generics

// FileEntryRef provides a minimal interface for FileEntry references
// to avoid circular dependencies between fileformat and metadata packages.
//
// This interface exposes only the methods needed for path metadata operations.
// The actual FileEntry type in the fileformat package implements this interface.
//
// Specification: api_generics.md: 1. Core Generic Types
type FileEntryRef interface {
	// GetPaths returns all paths associated with this file entry.
	//
	// Returns:
	//   - []PathEntry: Slice of path entries
	GetPaths() []PathEntry

	// GetFileID returns the unique file identifier.
	//
	// Returns:
	//   - uint64: File identifier
	GetFileID() uint64
}
