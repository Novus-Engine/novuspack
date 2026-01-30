// This file contains helper functions for path metadata operations.
// It provides common functionality used across path metadata file operations
// and file-path association operations to reduce code duplication.
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"path"

	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// findFileEntryByPath searches for a FileEntry by matching any of its paths.
//
// Parameters:
//   - path: The path string to search for
//
// Returns:
//   - *metadata.FileEntry: The found FileEntry, or nil if not found
//   - error: *PackageError if file not found
func (p *filePackage) findFileEntryByPath(path string) (*metadata.FileEntry, error) {
	for _, fe := range p.FileEntries {
		if fe == nil {
			continue
		}
		// Check if any of FileEntry's paths match the search path
		for _, pe := range fe.Paths {
			if pe.Path == path {
				return fe, nil
			}
		}
	}

	return nil, pkgerrors.NewPackageError(
		pkgerrors.ErrTypeValidation,
		"file not found",
		nil,
		pkgerrors.ValidationErrorContext{
			Field:    "FilePath",
			Value:    path,
			Expected: "existing file path",
		},
	)
}

// findPathMetadataByPath searches for a PathMetadataEntry by path.
//
// Parameters:
//   - path: The path string to search for
//
// Returns:
//   - *metadata.PathMetadataEntry: The found PathMetadataEntry, or nil if not found
//   - error: *PackageError if path metadata not found
func (p *filePackage) findPathMetadataByPath(path string) (*metadata.PathMetadataEntry, error) {
	for _, pme := range p.PathMetadataEntries {
		if pme != nil && pme.Path.Path == path {
			return pme, nil
		}
	}

	return nil, pkgerrors.NewPackageError(
		pkgerrors.ErrTypeValidation,
		"path metadata not found",
		nil,
		pkgerrors.ValidationErrorContext{
			Field:    "Path",
			Value:    path,
			Expected: "existing path metadata entry",
		},
	)
}

// isValidParentPath checks if a parent path is valid for a child path.
//
// A parent path is valid if:
//   - It is not "." (current directory)
//   - It is different from the child path
//
// Parameters:
//   - childPath: The child path
//   - parentPath: The parent path to validate
//
// Returns:
//   - bool: true if parent path is valid, false otherwise
func isValidParentPath(childPath, parentPath string) bool {
	return parentPath != "." && parentPath != childPath
}

// setParentPathAssociation sets the parent path association for a PathMetadataEntry.
//
// This function computes the parent path and looks up the parent PathMetadataEntry
// in the package. If found, it sets the ParentPath field. If the parent is not found
// or the entry is nil, the function simply returns without setting the ParentPath.
//
// This operation never fails - missing parents are valid (root paths, optional metadata).
// ParentPath is a runtime-only convenience field for hierarchy traversal.
//
// Note: Directory paths may have trailing slashes. This function normalizes paths
// by stripping trailing slashes before calling filepath.Dir(), then tries to match
// both with and without trailing slashes.
//
// Parameters:
//   - pme: The PathMetadataEntry to set parent for (nil is safe)
func (p *filePackage) setParentPathAssociation(pme *metadata.PathMetadataEntry) {
	if pme == nil {
		return
	}

	currentPath := pme.Path.Path

	// Strip trailing slash for directory paths before computing parent
	pathForDir := currentPath
	if len(pathForDir) > 0 && pathForDir[len(pathForDir)-1] == '/' {
		pathForDir = pathForDir[:len(pathForDir)-1]
	}

	// Compute parent path using path.Dir for platform-stable forward slashes
	parentPath := path.Dir(pathForDir)

	// Check if parent is valid
	if !isValidParentPath(currentPath, parentPath) {
		return
	}

	// Look up parent PathMetadataEntry (try both with and without trailing slash)
	for _, parentPME := range p.PathMetadataEntries {
		if parentPME == nil {
			continue
		}
		if parentPME.Path.Path == parentPath || parentPME.Path.Path == parentPath+"/" {
			pme.ParentPath = parentPME
			return
		}
	}

	// Parent not found - this is normal for root paths or optional metadata
}
