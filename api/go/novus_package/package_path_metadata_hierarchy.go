// This file implements path hierarchy operations for path metadata.
// It contains methods for retrieving path information, listing paths,
// and building path hierarchies. This file should contain only hierarchy
// operations (GetPathInfo, ListPaths, GetPathHierarchy) as specified in
// api_metadata.md Section 8.2.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// GetPathInfo returns path information for a specific path.
//
// Parameters:
//   - path: Path string
//
// Returns:
//   - *PathInfo: Path information, or nil if not found
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetPathInfo(path string) (*PathInfo, error) {
	// This is a pure in-memory operation - path metadata is already loaded
	if p.PathMetadataEntries == nil {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path not found",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "existing path",
			},
		)
	}

	entries := p.PathMetadataEntries

	// Find path
	var targetEntry *metadata.PathMetadataEntry
	for _, entry := range entries {
		if entry.GetPath() == path {
			targetEntry = entry
			break
		}
	}

	if targetEntry == nil {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path not found",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "existing path",
			},
		)
	}

	// Count files associated with this path
	fileCount := len(targetEntry.GetAssociatedFileEntries())

	// Find subdirectories
	subDirs := make([]string, 0)
	parentPath := ""
	if targetEntry.ParentPath != nil {
		parentPath = targetEntry.GetParentPathString()
	}

	// Find immediate subdirectories (paths that have this path as parent)
	for _, entry := range entries {
		if entry.ParentPath == targetEntry && entry.IsDirectory() {
			subDirs = append(subDirs, entry.GetPath())
		}
	}

	return &PathInfo{
		Entry:      targetEntry,
		FileCount:  fileCount,
		SubDirs:    subDirs,
		ParentPath: parentPath,
		Depth:      targetEntry.GetDepth(),
	}, nil
}

// ListPaths returns all path information entries.
//
// This is a pure in-memory operation that does not require context.
//
// Returns:
//   - []PathInfo: All path information entries
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) ListPaths() ([]PathInfo, error) {
	// This is a pure in-memory operation - path metadata is already loaded
	if p.PathMetadataEntries == nil {
		return []PathInfo{}, nil
	}

	entries := p.PathMetadataEntries

	// Convert to PathInfo
	result := make([]PathInfo, 0, len(entries))
	for _, entry := range entries {
		pathInfo, err := p.GetPathInfo(entry.GetPath())
		if err != nil {
			// Skip entries that can't be converted
			continue
		}
		result = append(result, *pathInfo)
	}

	return result, nil
}

// GetPathHierarchy returns the path hierarchy as a map of parent paths to child paths.
//
// This is a pure in-memory operation that does not require context.
//
// Returns:
//   - map[string][]string: Map of parent path -> child paths
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetPathHierarchy() (map[string][]string, error) {
	// This is a pure in-memory operation - path metadata is already loaded
	if p.PathMetadataEntries == nil {
		return make(map[string][]string), nil
	}

	entries := p.PathMetadataEntries

	// Build hierarchy map
	hierarchy := make(map[string][]string)
	for _, entry := range entries {
		if entry.ParentPath != nil {
			parentPath := entry.ParentPath.GetPath()
			hierarchy[parentPath] = append(hierarchy[parentPath], entry.GetPath())
		} else {
			// Root paths
			hierarchy[""] = append(hierarchy[""], entry.GetPath())
		}
	}

	return hierarchy, nil
}
