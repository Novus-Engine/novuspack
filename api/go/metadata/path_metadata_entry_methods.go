// This file implements path management and hierarchy methods for PathMetadataEntry
// structures. It contains methods for managing parent paths, path depth, ancestors,
// and path hierarchy operations. This file should contain path hierarchy and
// filesystem property methods for PathMetadataEntry as specified in api_metadata.md
// Section 8.1.2.
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure

package metadata

import (
	"path/filepath"
	"strings"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// SetPath sets the path for this PathMetadataEntry.
//
// Parameters:
//   - path: Path string
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) SetPath(path string) {
	pme.Path.Path = path
	pme.Path.PathLength = uint16(len(path))
}

// GetPath returns the path string.
//
// Returns:
//   - string: Path string
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetPath() string {
	return pme.Path.Path
}

// GetPathForPlatform returns the path string converted for the specified platform.
//
// This method implements the display format conversion per api_generics 1.3.3.7:
// 1. Strips the leading "/" (package root indicator, not filesystem root)
// 2. On Windows: converts forward slashes to backslashes
// 3. On Unix/Linux: returns path with forward slashes (no leading "/")
//
// Parameters:
//   - isWindows: True for Windows path format, false for Unix/Linux
//
// Returns:
//   - string: Path in display format for the specified platform
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetPathForPlatform(isWindows bool) string {
	// Strip leading "/" for display (per 1.3.3.7)
	display := strings.TrimPrefix(pme.Path.Path, "/")

	if isWindows {
		// Convert forward slashes to backslashes for Windows
		return strings.ReplaceAll(display, "/", "\\")
	}

	// Return with forward slashes for Unix/Linux (no leading "/")
	return display
}

// GetPathEntry returns the PathEntry for this PathMetadataEntry.
//
// Returns:
//   - generics.PathEntry: Path entry
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetPathEntry() generics.PathEntry {
	return pme.Path
}

// GetType returns the type of this path entry.
//
// Returns:
//   - PathMetadataType: Path type
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetType() PathMetadataType {
	return pme.Type
}

// IsDirectory checks if this path is a directory.
//
// Returns:
//   - bool: True if this is a directory
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) IsDirectory() bool {
	return pme.Type == PathMetadataTypeDirectory
}

// IsFile checks if this path is a file.
//
// Returns:
//   - bool: True if this is a file
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) IsFile() bool {
	return pme.Type == PathMetadataTypeFile
}

// IsSymlink checks if this path is a symbolic link.
//
// Returns:
//   - bool: True if this is a symlink
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) IsSymlink() bool {
	return pme.Type == PathMetadataTypeFileSymlink || pme.Type == PathMetadataTypeDirectorySymlink || pme.Path.IsSymlink
}

// GetLinkTarget returns the symbolic link target.
//
// Returns:
//   - string: Link target path (empty if not a symlink)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetLinkTarget() string {
	if pme.FileSystem.LinkTarget != "" {
		return pme.FileSystem.LinkTarget
	}
	return pme.Path.LinkTarget
}

// ResolveSymlink resolves the symbolic link target.
//
// Returns:
//   - string: Resolved path
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) ResolveSymlink() string {
	target := pme.GetLinkTarget()
	if target == "" {
		return pme.Path.Path
	}

	// If target is absolute, return it
	if filepath.IsAbs(target) {
		return target
	}

	// Resolve relative to the path's parent
	dir := filepath.Dir(pme.Path.Path)
	resolved := filepath.Join(dir, target)
	return filepath.Clean(resolved)
}

// SetParentPath sets the parent path for this PathMetadataEntry.
//
// Parameters:
//   - parent: Pointer to parent path (nil for root)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) SetParentPath(parent *PathMetadataEntry) {
	pme.ParentPath = parent
}

// GetParentPath returns the parent path for this PathMetadataEntry.
//
// Returns:
//   - *PathMetadataEntry: Pointer to parent path (nil for root)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetParentPath() *PathMetadataEntry {
	return pme.ParentPath
}

// GetParentPathString returns the path string of the parent path.
//
// Returns:
//   - string: Parent path string (empty for root)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetParentPathString() string {
	if pme.ParentPath == nil {
		return ""
	}
	return pme.ParentPath.Path.Path
}

// GetDepth returns the depth of this path in the hierarchy.
//
// Returns:
//   - int: Path depth (0 for root)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetDepth() int {
	depth := 0
	current := pme.ParentPath
	for current != nil {
		depth++
		current = current.ParentPath
	}
	return depth
}

// IsRoot checks if this path is the root path.
//
// Returns:
//   - bool: True if this is the root path (no parent)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) IsRoot() bool {
	return pme.ParentPath == nil
}

// GetAncestors returns all ancestor paths up to the root.
//
// Returns:
//   - []*PathMetadataEntry: Slice of ancestor paths (empty for root)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetAncestors() []*PathMetadataEntry {
	var ancestors []*PathMetadataEntry
	current := pme.ParentPath
	for current != nil {
		ancestors = append(ancestors, current)
		current = current.ParentPath
	}
	return ancestors
}

// GetInheritedTags returns tags inherited from parent paths.
//
// Walks up the ParentPath chain and collects tags from paths with inheritance enabled,
// sorted by priority (higher priority first).
//
// Returns:
//   - []*Tag[any]: Inherited tags
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetInheritedTags() ([]*generics.Tag[any], error) {
	if pme.ParentPath == nil {
		return []*generics.Tag[any]{}, nil
	}

	// Collect all paths with inheritance enabled
	type pathWithTags struct {
		path     *PathMetadataEntry
		priority int
		tags     []*generics.Tag[any]
	}

	var pathsWithTags []pathWithTags
	processedPaths := make(map[*PathMetadataEntry]bool)

	// Walk up the path hierarchy and collect paths
	current := pme.ParentPath
	for current != nil {
		if processedPaths[current] {
			break // Prevent cycles
		}
		processedPaths[current] = true

		// Only collect tags if inheritance is enabled
		if current.Inheritance != nil && current.Inheritance.Enabled {
			pathTags, err := GetPathMetaTags(current)
			if err != nil {
				// Continue on error, don't fail entire inheritance
				current = current.ParentPath
				continue
			}

			if len(pathTags) > 0 {
				pathsWithTags = append(pathsWithTags, pathWithTags{
					path:     current,
					priority: current.Inheritance.Priority,
					tags:     pathTags,
				})
			}
		}

		current = current.ParentPath
	}

	// Sort by priority (higher priority first)
	// Use a simple insertion sort for small lists
	for i := 1; i < len(pathsWithTags); i++ {
		key := pathsWithTags[i]
		j := i - 1
		for j >= 0 && pathsWithTags[j].priority < key.priority {
			pathsWithTags[j+1] = pathsWithTags[j]
			j--
		}
		pathsWithTags[j+1] = key
	}

	// Apply tags in priority order (higher priority overwrites lower priority)
	tagMap := make(map[string]*generics.Tag[any])
	for _, pwt := range pathsWithTags {
		for _, tag := range pwt.tags {
			tagMap[tag.Key] = tag
		}
	}

	// Convert map to slice
	result := make([]*generics.Tag[any], 0, len(tagMap))
	for _, tag := range tagMap {
		result = append(result, tag)
	}

	return result, nil
}

// GetEffectiveTags returns all tags for this PathMetadataEntry, including:
// 1. Tags directly on this PathMetadataEntry
// 2. Tags inherited from parent PathMetadataEntry instances (path hierarchy)
// 3. Tags from associated FileEntry instances (treated as if applied to this PathMetadataEntry)
//
// Returns:
//   - []*Tag[any]: All effective tags
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetEffectiveTags() ([]*generics.Tag[any], error) {
	tagMap := make(map[string]*generics.Tag[any])

	// 1. Add tags directly on this PathMetadataEntry
	directTags, err := GetPathMetaTags(pme)
	if err != nil {
		return nil, err
	}
	for _, tag := range directTags {
		tagMap[tag.Key] = tag
	}

	// 2. Add inherited tags from parent paths
	inheritedTags, err := pme.GetInheritedTags()
	if err != nil {
		// Log error but continue with direct tags
	} else {
		for _, tag := range inheritedTags {
			// Direct tags take precedence over inherited tags
			if _, exists := tagMap[tag.Key]; !exists {
				tagMap[tag.Key] = tag
			}
		}
	}

	// 3. Add tags from associated FileEntry instances
	for _, fe := range pme.AssociatedFileEntries {
		if fe == nil {
			continue
		}

		fileTags, err := GetFileEntryTags(fe)
		if err != nil {
			// Continue on error, don't fail entire operation
			continue
		}

		for _, tag := range fileTags {
			// Direct and inherited tags take precedence over FileEntry tags
			if _, exists := tagMap[tag.Key]; !exists {
				tagMap[tag.Key] = tag
			}
		}
	}

	// Convert map to slice
	result := make([]*generics.Tag[any], 0, len(tagMap))
	for _, tag := range tagMap {
		result = append(result, tag)
	}

	return result, nil
}

// AssociateWithFileEntry associates this PathMetadataEntry with a FileEntry.
//
// The association is established if the PathMetadataEntry.Path.Path matches one of the FileEntry.Paths.
// This method also calls fe.AssociateWithPathMetadata() to establish bidirectional association.
//
// Parameters:
//   - fe: FileEntry to associate with
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) AssociateWithFileEntry(fe *FileEntry) error {
	if fe == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "FileEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "FileEntry",
			Value: nil,
		})
	}

	// Check if path matches
	pathStr := pme.Path.Path
	matched := false
	for _, pe := range fe.Paths {
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

	// Check if already associated
	for _, feRef := range pme.AssociatedFileEntries {
		if feRef == fe {
			// Already associated, nothing to do
			return nil
		}
	}

	// Add to AssociatedFileEntries
	pme.AssociatedFileEntries = append(pme.AssociatedFileEntries, fe)

	// Establish bidirectional association by calling FileEntry method
	if err := fe.AssociateWithPathMetadata(pme); err != nil {
		// Remove from AssociatedFileEntries if association fails
		for i, feRef := range pme.AssociatedFileEntries {
			if feRef == fe {
				pme.AssociatedFileEntries = append(pme.AssociatedFileEntries[:i], pme.AssociatedFileEntries[i+1:]...)
				break
			}
		}
		return err
	}

	return nil
}

// GetAssociatedFileEntries returns all FileEntry instances associated with this PathMetadataEntry.
//
// Returns:
//   - []*FileEntry: Slice of associated FileEntry instances (empty if none)
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
func (pme *PathMetadataEntry) GetAssociatedFileEntries() []*FileEntry {
	return pme.AssociatedFileEntries
}
