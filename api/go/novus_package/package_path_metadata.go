// This file implements core package-level path metadata management operations.
// It contains methods for managing PathMetadataEntry instances at the package
// level, including GetPathMetadata, SetPathMetadata, AddPathMetadata, RemovePathMetadata,
// UpdatePathMetadata, ValidatePathMetadata, and GetPathConflicts. This file should
// contain only core path metadata operations as specified in api_metadata.md
// Section 8.2. Directory operations, hierarchy operations, file-path associations,
// and special metadata file operations are in separate files.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// PathInfo provides runtime path metadata information.
//
// Specification: api_metadata.md: 8.1.2 PathMetadataEntry Structure
type PathInfo struct {
	Entry      *metadata.PathMetadataEntry // Path metadata entry data
	FileCount  int                         // Number of files in this path
	SubDirs    []string                    // Immediate subdirectories
	ParentPath string                      // Parent path
	Depth      int                         // Path depth in hierarchy
}

// GetPathMetadata returns all path metadata entries from the package.
//
// Loads path metadata from special metadata file (type 65001).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - []*PathMetadataEntry: All path metadata entries
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetPathMetadata(ctx context.Context) ([]*metadata.PathMetadataEntry, error) {
	if err := internal.CheckContext(ctx, "GetPathMetadata"); err != nil {
		return nil, err
	}

	// Return cached entries if available
	if p.PathMetadataEntries != nil {
		return p.PathMetadataEntries, nil
	}

	// Load from special metadata file
	// Note: LoadPathMetadataFile is currently a stub that always returns an error
	// When fully implemented, it will return nil on success and populate PathMetadataEntries
	// Since the stub always returns an error, we return it directly
	// This will be updated when LoadPathMetadataFile is fully implemented
	return nil, p.LoadPathMetadataFile(ctx)
}

// SetPathMetadata sets all path metadata entries for the package.
//
// Saves path metadata to special metadata file (type 65001).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - entries: Path metadata entries to set
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) SetPathMetadata(ctx context.Context, entries []*metadata.PathMetadataEntry) error {
	if err := internal.CheckContext(ctx, "SetPathMetadata"); err != nil {
		return err
	}

	// Validate all entries
	for i, entry := range entries {
		if entry == nil {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path metadata entry is nil", nil, pkgerrors.ValidationErrorContext{
				Field: "entries",
				Value: i,
			})
		}
		if err := entry.Validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, "path metadata entry validation failed", pkgerrors.ValidationErrorContext{
				Field: "entries",
				Value: i,
			})
		}
	}

	// Set cached entries
	p.PathMetadataEntries = entries

	// Save to special metadata file
	return p.SavePathMetadataFile(ctx)
}

// AddPathMetadata adds a new path metadata entry to the package.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Path string
//   - pathType: Path type (file, directory, symlink)
//   - properties: Map of tag properties (key-value pairs)
//   - inheritance: Inheritance settings (optional, for directories)
//   - metadata: Path metadata (optional, for directories)
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) AddPathMetadata(ctx context.Context, path string, pathType metadata.PathMetadataType, properties map[string]string, inheritance *metadata.PathInheritance, pathMetadata *metadata.PathMetadata) error {
	if err := internal.CheckContext(ctx, "AddPathMetadata"); err != nil {
		return err
	}

	// Load existing path metadata
	entries, err := p.GetPathMetadata(ctx)
	if err != nil {
		return err
	}

	// Check if path already exists
	for _, entry := range entries {
		if entry.GetPath() == path {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path already exists", nil, pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "non-existent path",
			})
		}
	}

	// Create new PathEntry
	pathEntry := generics.PathEntry{
		Path:       path,
		PathLength: uint16(len(path)),
	}

	// Convert properties map to tags
	tags := make([]*generics.Tag[any], 0, len(properties))
	for key, value := range properties {
		tag := generics.NewTag[any](key, value, generics.TagValueTypeString)
		tags = append(tags, tag)
	}

	// Create new PathMetadataEntry
	newEntry := &metadata.PathMetadataEntry{
		Path:        pathEntry,
		Type:        pathType,
		Properties:  tags,
		Inheritance: inheritance,
		Metadata:    pathMetadata,
		FileSystem:  metadata.PathFileSystem{},
	}

	// Validate entry
	if err := newEntry.Validate(); err != nil {
		return err
	}

	// Add to entries
	entries = append(entries, newEntry)

	// Save
	return p.SetPathMetadata(ctx, entries)
}

// RemovePathMetadata removes a path metadata entry from the package.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Path string to remove
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) RemovePathMetadata(ctx context.Context, path string) error {
	if err := internal.CheckContext(ctx, "RemovePathMetadata"); err != nil {
		return err
	}

	// Load existing path metadata
	entries, err := p.GetPathMetadata(ctx)
	if err != nil {
		return err
	}

	// Find and remove path
	found := false
	newEntries := make([]*metadata.PathMetadataEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.GetPath() == path {
			found = true
			continue
		}
		newEntries = append(newEntries, entry)
	}

	if !found {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path not found", nil, pkgerrors.ValidationErrorContext{
			Field:    "path",
			Value:    path,
			Expected: "existing path",
		})
	}

	// Save
	return p.SetPathMetadata(ctx, newEntries)
}

// UpdatePathMetadata updates an existing path metadata entry in the package.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Path string to update
//   - properties: Map of tag properties (key-value pairs, nil to keep existing)
//   - inheritance: Inheritance settings (nil to keep existing)
//   - metadata: Path metadata (nil to keep existing)
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) UpdatePathMetadata(ctx context.Context, path string, properties map[string]string, inheritance *metadata.PathInheritance, pathMetadata *metadata.PathMetadata) error {
	if err := internal.CheckContext(ctx, "UpdatePathMetadata"); err != nil {
		return err
	}

	// Load existing path metadata
	entries, err := p.GetPathMetadata(ctx)
	if err != nil {
		return err
	}

	// Find path
	var targetEntry *metadata.PathMetadataEntry
	for _, entry := range entries {
		if entry.GetPath() == path {
			targetEntry = entry
			break
		}
	}

	if targetEntry == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path not found", nil, pkgerrors.ValidationErrorContext{
			Field:    "path",
			Value:    path,
			Expected: "existing path",
		})
	}

	// Update properties if provided
	if properties != nil {
		// Convert properties map to tags
		tags := make([]*generics.Tag[any], 0, len(properties))
		for key, value := range properties {
			tag := generics.NewTag[any](key, value, generics.TagValueTypeString)
			tags = append(tags, tag)
		}
		targetEntry.Properties = tags
	}

	// Update inheritance if provided
	if inheritance != nil {
		targetEntry.Inheritance = inheritance
	}

	// Update metadata if provided
	if pathMetadata != nil {
		targetEntry.Metadata = pathMetadata
	}

	// Validate entry
	if err := targetEntry.Validate(); err != nil {
		return err
	}

	// Save
	return p.SetPathMetadata(ctx, entries)
}

// ValidatePathMetadata validates all path metadata entries in the package.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) ValidatePathMetadata(ctx context.Context) error {
	if err := internal.CheckContext(ctx, "ValidatePathMetadata"); err != nil {
		return err
	}

	// Load path metadata
	entries, err := p.GetPathMetadata(ctx)
	if err != nil {
		return err
	}

	// Validate each entry
	for i, entry := range entries {
		if err := entry.Validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, "path metadata entry validation failed", pkgerrors.ValidationErrorContext{
				Field: "entries",
				Value: i,
			})
		}
	}

	return nil
}

// GetPathConflicts returns a list of conflicting paths.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - []string: List of conflicting paths
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetPathConflicts(ctx context.Context) ([]string, error) {
	if err := internal.CheckContext(ctx, "GetPathConflicts"); err != nil {
		return nil, err
	}

	// Load path metadata
	entries, err := p.GetPathMetadata(ctx)
	if err != nil {
		return nil, err
	}

	// Find duplicate paths
	pathMap := make(map[string]int)
	conflicts := make([]string, 0)

	for _, entry := range entries {
		path := entry.GetPath()
		pathMap[path]++
		if pathMap[path] > 1 {
			conflicts = append(conflicts, path)
		}
	}

	return conflicts, nil
}
