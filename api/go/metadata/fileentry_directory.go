// This file implements directory-related operations for FileEntry structures.
// It contains methods for working with directory entries, checking if entries
// are directories, and managing directory-specific metadata. This file should
// contain directory-related methods for FileEntry as needed for directory
// handling in the package system.
//
// Specification: api_file_mgmt_index.md: 0 Overview

package metadata

import (
	"path/filepath"
)

// GetParentPath returns the path of the parent path metadata entry.
//
// Returns the parent path string from the first associated PathMetadataEntry.
// If no PathMetadataEntry is associated, returns empty string.
// Uses only interface methods to avoid circular dependency.
//
// Returns:
//   - string: Parent path string (empty if root-relative or not set)
//
// Specification: api_file_mgmt_file_entry.md: 1 FileEntry Structure
func (f *FileEntry) GetParentPath() string {
	// Get parent path by deriving it from the path string
	for _, pme := range f.PathMetadataEntries {
		if pme == nil {
			continue
		}
		path := pme.Path.Path
		if path != "" {
			// Derive parent path from the path string
			parentPath := filepath.Dir(path)
			if parentPath != "." && parentPath != path {
				return parentPath
			}
		}
	}

	return ""
}

// GetDirectoryDepth returns the depth of this file in the path hierarchy.
//
// Returns the depth from the first associated PathMetadataEntry.
// If no PathMetadataEntry is associated, returns 0.
//
// Returns:
//   - int: Path depth (0 for root-relative files)
//
// Specification: api_metadata.md: 8.1.8 Path Management Methods
func (f *FileEntry) GetDirectoryDepth() int {
	// Get depth from the first associated PathMetadataEntry
	for _, pme := range f.PathMetadataEntries {
		if pme != nil {
			return pme.GetDepth()
		}
	}

	return 0
}

// IsRootRelative checks if the file is root-relative (no parent path).
//
// Returns true if no PathMetadataEntry is associated or all associated
// PathMetadataEntry instances have no parent path.
//
// Returns:
//   - bool: True if file is root-relative, false otherwise
//
// Specification: api_file_mgmt_file_entry.md: 1 FileEntry Structure
func (f *FileEntry) IsRootRelative() bool {
	// Check if any associated PathMetadataEntry has a parent path
	for _, pme := range f.PathMetadataEntries {
		if pme != nil && pme.ParentPath != nil {
			return false
		}
	}

	return true
}
