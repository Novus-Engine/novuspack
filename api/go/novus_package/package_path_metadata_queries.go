// This file implements path metadata query methods.
// It contains methods for querying path metadata including directory counts.
// Note: ValidatePathMetadata and GetPathConflicts are already implemented in
// package_path_metadata.go
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"strings"
)

// GetDirectoryCount returns the total number of directories in the package.
//
// This counts both explicit directories from PathMetadataEntries and implied
// directories from file paths in FileEntries. Directories are automatically
// deduplicated.
//
// Returns:
//   - int: Number of unique directories
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetDirectoryCount() (int, error) {
	// This is a pure in-memory operation
	directories := make(map[string]bool)

	// Collect explicit directories from PathMetadataEntries
	if p.PathMetadataEntries != nil {
		for _, entry := range p.PathMetadataEntries {
			if entry.IsDirectory() {
				directories[entry.GetPath()] = true
			}
		}
	}

	// Collect implied directories from FileEntries paths
	if p.FileEntries != nil {
		for _, fileEntry := range p.FileEntries {
			for _, pathEntry := range fileEntry.Paths {
				// Extract all parent directories from the path
				extractParentDirectories(pathEntry.Path, directories)
			}
		}
	}

	return len(directories), nil
}

// extractParentDirectories extracts all parent directory paths from a file path.
// For example, "/foo/bar/file.txt" yields "/foo/" and "/foo/bar/"
// Directories are added to the provided map for deduplication.
func extractParentDirectories(path string, directories map[string]bool) {
	// Normalize to ensure leading slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Split path and build parent directories
	parts := strings.Split(path, "/")
	currentPath := ""

	// Iterate through parts, excluding the last part (filename) and empty parts
	for i, part := range parts {
		if part == "" {
			continue
		}

		currentPath += "/" + part

		// Don't count the last segment if it's a file (no trailing slash in original)
		if i < len(parts)-1 {
			// This is a directory segment
			dirPath := currentPath + "/"
			directories[dirPath] = true
		}
	}
}

// Note: ValidatePathMetadata and GetPathConflicts are already implemented
// in package_path_metadata.go (lines 298 and 332 respectively)
