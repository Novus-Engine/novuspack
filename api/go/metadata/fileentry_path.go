// This file implements path management operations for FileEntry structures.
// It contains methods for managing multiple paths per file entry, resolving
// symlinks, and working with path metadata associations. This file should
// contain path-related methods for FileEntry as specified in api_file_mgmt_index.md
// Section 1.5.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure

package metadata

import (
	"path/filepath"
	"strings"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// HasSymlinks checks if the file entry has any symbolic link paths.
//
// Returns:
//   - bool: True if any path is a symbolic link, false otherwise
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) HasSymlinks() bool {
	for _, path := range f.Paths {
		if path.IsSymlink {
			return true
		}
	}
	return false
}

// GetSymlinkPaths returns all paths that are symbolic links.
//
// Returns:
//   - []PathEntry: Slice of path entries that are symbolic links
//
// Specification: api_file_mgmt_file_entry.md: 5. Path Management
func (f *FileEntry) GetSymlinkPaths() []generics.PathEntry {
	var symlinks []generics.PathEntry
	for _, path := range f.Paths {
		if path.IsSymlink {
			symlinks = append(symlinks, path)
		}
	}
	return symlinks
}

// GetPrimaryPath returns the primary path in display format (no leading slash).
// Returns the first path from the Paths slice, with the leading "/" stripped
// for user-facing display per api_generics 1.3.3.7.
//
// Returns:
//   - string: Primary path in display format (empty if no paths)
//
// Specification: api_file_mgmt_file_entry.md: 5. Path Management
func (f *FileEntry) GetPrimaryPath() string {
	if len(f.Paths) == 0 {
		return ""
	}
	return strings.TrimPrefix(f.Paths[0].Path, "/")
}

// ResolveAllSymlinks resolves all symlink targets for this file entry.
//
// Returns all resolved symlink targets.
// Non-symlink paths are returned as-is.
//
// Returns:
//   - []string: Slice of resolved paths
//
// Specification: api_file_mgmt_file_entry.md: 5. Path Management
func (f *FileEntry) ResolveAllSymlinks() []string {
	var resolved []string
	for _, path := range f.Paths {
		if path.IsSymlink && path.LinkTarget != "" {
			// Resolve the symlink target
			// If LinkTarget is relative, resolve it relative to the path's directory
			if filepath.IsAbs(path.LinkTarget) {
				resolved = append(resolved, path.LinkTarget)
			} else {
				// Resolve relative to the path's directory
				dir := filepath.Dir(path.Path)
				resolvedPath := filepath.Join(dir, path.LinkTarget)
				resolved = append(resolved, filepath.Clean(resolvedPath))
			}
		} else {
			// Not a symlink, return path as-is
			resolved = append(resolved, path.Path)
		}
	}
	return resolved
}

// AssociateWithPathMetadata associates this FileEntry with a PathMetadataEntry.
//
// The association is established by matching the PathMetadataEntry.Path.Path with one of the FileEntry.Paths.
// This method also calls PathMetadataEntry.AssociateWithFileEntry() to establish bidirectional association.
//
// Parameters:
//   - pme: PathMetadataEntry to associate with
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) AssociateWithPathMetadata(pme *PathMetadataEntry) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Initialize PathMetadataEntries map if needed
	if f.PathMetadataEntries == nil {
		f.PathMetadataEntries = make(map[string]*PathMetadataEntry)
	}

	// Match pme.Path.Path with FileEntry.Paths entries
	pathStr := pme.Path.Path
	matched := false
	for _, pe := range f.Paths {
		if pe.Path == pathStr {
			matched = true
			break
		}
	}

	if !matched {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path does not match any FileEntry path", nil, pkgerrors.ValidationErrorContext{
			Field:    "Path",
			Value:    pathStr,
			Expected: "path matching one of FileEntry.Paths",
		})
	}

	// Add to PathMetadataEntries map using path string as key
	f.PathMetadataEntries[pathStr] = pme

	// Establish bidirectional association by calling PathMetadataEntry method
	if err := pme.AssociateWithFileEntry(f); err != nil {
		// Remove from map if association fails
		delete(f.PathMetadataEntries, pathStr)
		return err
	}

	return nil
}

// GetPathMetadataForPath returns the PathMetadataEntry associated with a specific path.
//
// Parameters:
//   - path: Path string to look up
//
// Returns:
//   - *PathMetadataEntry: Associated PathMetadataEntry, or nil if not found
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) GetPathMetadataForPath(path string) *PathMetadataEntry {
	if f.PathMetadataEntries == nil {
		return nil
	}
	return f.PathMetadataEntries[path]
}

// GetPaths returns all paths associated with this file entry.
//
// This method implements the generics.FileEntryRef interface for external use.
//
// Returns:
//   - []generics.PathEntry: Slice of path entries
//
// Specification: api_file_mgmt_file_entry.md: 5. Path Management
func (f *FileEntry) GetPaths() []generics.PathEntry {
	return f.Paths
}

// GetFileID returns the unique file identifier.
//
// This method implements the generics.FileEntryRef interface for external use.
//
// Returns:
//   - uint64: File identifier
//
// Specification: api_file_mgmt_file_entry.md: 5. Path Management
func (f *FileEntry) GetFileID() uint64 {
	return f.FileID
}
