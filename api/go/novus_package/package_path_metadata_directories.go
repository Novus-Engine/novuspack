// This file implements directory-specific path metadata operations.
// It contains convenience wrapper methods for directory paths that delegate
// to core path operations. This file should contain only directory-specific
// operations (AddDirectoryMetadata, RemoveDirectoryMetadata, UpdateDirectoryMetadata, ListDirectories)
// as specified in api_metadata.md Section 8.2.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"
	"strings"

	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
)

// AddDirectoryMetadata adds a new directory path metadata entry to the package.
//
// This is a convenience wrapper around AddPathMetadata for directory paths.
// This is a metadata-only method for path metadata management, distinct from the
// file management AddDirectory method which recursively adds directory files.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Directory path string (should end with /)
//   - properties: Map of tag properties
//   - inheritance: Inheritance settings
//   - metadata: Path metadata
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) AddDirectoryMetadata(ctx context.Context, path string, properties map[string]string, inheritance *metadata.PathInheritance, pathMetadata *metadata.PathMetadata) error {
	if err := internal.CheckContext(ctx, "AddDirectoryMetadata"); err != nil {
		return err
	}

	// Ensure path ends with /
	if path != "" && path[len(path)-1] != '/' {
		path += "/"
	}
	return p.AddPathMetadata(
		ctx,
		path,
		metadata.PathMetadataTypeDirectory,
		properties,
		inheritance,
		pathMetadata,
	)
}

// RemoveDirectoryMetadata removes a directory path metadata entry from the package.
//
// This is a convenience wrapper around RemovePathMetadata for directory paths.
// This is a metadata-only method for path metadata management, distinct from the
// file management RemoveDirectory method which removes directory files.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Directory path string
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) RemoveDirectoryMetadata(ctx context.Context, path string) error {
	if err := internal.CheckContext(ctx, "RemoveDirectoryMetadata"); err != nil {
		return err
	}

	// Ensure path ends with /
	if path != "" && path[len(path)-1] != '/' {
		path += "/"
	}
	return p.RemovePathMetadata(ctx, path)
}

// UpdateDirectoryMetadata updates an existing directory path metadata entry in the package.
//
// This is a convenience wrapper around UpdatePathMetadata for directory paths.
// This is a metadata-only method for path metadata management, distinct from the
// file management operations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Directory path string
//   - properties: Map of tag properties
//   - inheritance: Inheritance settings
//   - metadata: Path metadata
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) UpdateDirectoryMetadata(ctx context.Context, path string, properties map[string]string, inheritance *metadata.PathInheritance, pathMetadata *metadata.PathMetadata) error {
	if err := internal.CheckContext(ctx, "UpdateDirectoryMetadata"); err != nil {
		return err
	}

	// Ensure path ends with /
	if path != "" && path[len(path)-1] != '/' {
		path += "/"
	}
	return p.UpdatePathMetadata(ctx, path, properties, inheritance, pathMetadata)
}

// ListDirectories returns all directory path information entries.
//
// This returns both explicit directories from PathMetadataEntries and implied
// directories from file paths in FileEntries. Implied directories have synthetic
// PathInfo entries with minimal metadata.
//
// Returns:
//   - []PathInfo: Directory path information entries
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
//
//nolint:gocognit // hierarchy walk branches
func (p *filePackage) ListDirectories() ([]PathInfo, error) {
	// This is an in-memory operation
	directories := make(map[string]*PathInfo)

	// Collect explicit directories from PathMetadataEntries
	if p.PathMetadataEntries != nil {
		for _, entry := range p.PathMetadataEntries {
			if entry.IsDirectory() {
				pathInfo, err := p.GetPathInfo(entry.GetPath())
				if err != nil {
					// Skip entries that can't be converted
					continue
				}
				directories[entry.GetPath()] = pathInfo
			}
		}
	}

	// Collect implied directories from FileEntries paths
	if p.FileEntries != nil {
		for _, fileEntry := range p.FileEntries {
			for _, pathEntry := range fileEntry.Paths {
				// Extract all parent directories from the path
				extractImpliedDirectories(pathEntry.Path, directories)
			}
		}
	}

	// Convert map to slice
	result := make([]PathInfo, 0, len(directories))
	for _, pathInfo := range directories {
		result = append(result, *pathInfo)
	}

	return result, nil
}

// extractImpliedDirectories extracts implied directory PathInfo from a file path.
// Only creates synthetic PathInfo entries for directories not already in the map.
func extractImpliedDirectories(path string, directories map[string]*PathInfo) {
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

		// Don't process the last segment if it's a file
		if i < len(parts)-1 {
			// This is a directory segment
			dirPath := currentPath + "/"

			// Only create synthetic entry if not already present
			if _, exists := directories[dirPath]; !exists {
				// Create synthetic PathInfo for implied directory
				directories[dirPath] = &PathInfo{
					Entry:      nil, // No explicit metadata entry
					FileCount:  0,   // Will be calculated if needed
					SubDirs:    []string{},
					ParentPath: "",
					Depth:      strings.Count(dirPath, "/") - 1,
				}
			}
		}
	}
}
