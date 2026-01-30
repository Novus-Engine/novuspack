// This file implements the PathMetadataEntry structure representing path metadata
// with inheritance, filesystem properties, and tag management. It contains the
// PathMetadataEntry type definition, path hierarchy management, and basic path
// operations. This file should contain the core PathMetadataEntry struct and
// path management methods as specified in api_metadata.md Section 8.
//
// Specification: api_metadata.md: 1. Comment Management

// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and api_metadata.md.
package metadata

import (
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// PathMetadataType represents the type of path entry.
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type PathMetadataType uint8

const (
	// PathMetadataTypeFile represents a regular file.
	PathMetadataTypeFile PathMetadataType = 0

	// PathMetadataTypeDirectory represents a regular directory.
	PathMetadataTypeDirectory PathMetadataType = 1

	// PathMetadataTypeFileSymlink represents a symlink to a file.
	PathMetadataTypeFileSymlink PathMetadataType = 2

	// PathMetadataTypeDirectorySymlink represents a symlink to a directory.
	PathMetadataTypeDirectorySymlink PathMetadataType = 3
)

// PathInheritance controls tag inheritance behavior (for directories only).
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type PathInheritance struct {
	Enabled  bool `yaml:"enabled"`  // Whether this path provides inheritance
	Priority int  `yaml:"priority"` // Inheritance priority (higher = more specific)
}

// PathMetadata contains path metadata (for directories only).
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type PathMetadata struct {
	Created     string `yaml:"created"`     // Path creation time (ISO8601)
	Modified    string `yaml:"modified"`    // Last modification time (ISO8601)
	Description string `yaml:"description"` // Human-readable description
}

// ACLEntry represents an Access Control List entry.
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type ACLEntry struct {
	Type  string  `yaml:"type"`         // "user", "group", "other", "mask"
	ID    *uint32 `yaml:"id,omitempty"` // User/Group ID (nil for "other")
	Perms string  `yaml:"perms"`        // Permissions (e.g., "rwx", "r--")
}

// PathMetadataEntry represents a path (file, directory, or symlink) with metadata, inheritance rules, and filesystem properties.
//
// PathMetadataEntry allows the same file content to have different permissions, ownership,
// and timestamps at different paths. This is stored separately from PathEntry to minimize
// the size of the core path structure.
//
// Path metadata is stored in special metadata files (similar to DirectoryEntry) and can be
// linked to PathEntry instances by matching the path string.
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type PathMetadataEntry struct {
	// Path is the path entry (minimal PathEntry from generics package)
	Path generics.PathEntry `yaml:"path"`

	// Type indicates whether this is a file, directory, file symlink, or directory symlink
	Type PathMetadataType `yaml:"type"`

	// Properties are path-specific tags (typed tags)
	Properties []*generics.Tag[any] `yaml:"properties"`

	// Inheritance controls tag inheritance behavior (optional, only for directories)
	Inheritance *PathInheritance `yaml:"inheritance,omitempty"` // Inheritance settings (nil for files)

	// Metadata contains path metadata (optional, only for directories)
	Metadata *PathMetadata `yaml:"metadata,omitempty"` // Path metadata (nil for files)

	// FileSystem contains filesystem-specific properties
	FileSystem PathFileSystem `yaml:"filesystem"`

	// Path hierarchy (runtime only, not stored in file)
	ParentPath *PathMetadataEntry `yaml:"-"` // Pointer to parent path (nil for root)

	// FileEntry associations (runtime only, not stored in file)
	// Since FileEntry and PathMetadataEntry are in the same package, direct references are used.
	AssociatedFileEntries []*FileEntry `yaml:"-"` // FileEntry instances associated with this path
}

// PathFileSystem contains filesystem-specific properties for a path.
//
// Specification: api_metadata.md: 8.1 PathMetadata Structures
type PathFileSystem struct {
	// Execute permissions (always captured)
	IsExecutable bool `yaml:"is_executable"` // Whether file has any execute permission bits set (tracked by default)

	// Unix/Linux properties (optional, captured when PreservePermissions is enabled)
	Mode *uint32    `yaml:"mode,omitempty"` // File/directory permissions and type (Unix-style)
	UID  *uint32    `yaml:"uid,omitempty"`  // User ID
	GID  *uint32    `yaml:"gid,omitempty"`  // Group ID
	ACL  []ACLEntry `yaml:"acl,omitempty"`  // Access Control List

	// Timestamps (Unix nanoseconds since epoch)
	ModTime    uint64 `yaml:"mod_time,omitempty"`    // Modification time
	CreateTime uint64 `yaml:"create_time,omitempty"` // Creation time
	AccessTime uint64 `yaml:"access_time,omitempty"` // Access time

	// Symbolic link support
	LinkTarget string `yaml:"link_target,omitempty"` // Target path for symbolic links (empty if not a symlink)

	// Windows properties (optional)
	WindowsAttrs *uint32 `yaml:"windows_attrs,omitempty"` // Windows attributes

	// Extended attributes (optional)
	ExtendedAttrs map[string]string `yaml:"extended_attrs,omitempty"` // Extended attributes

	// Filesystem flags (optional)
	Flags *uint16 `yaml:"flags,omitempty"` // Filesystem-specific flags
}

// Validate performs validation checks on the PathMetadataEntry.
//
// Validation checks:
//   - PathEntry must be valid (call PathEntry.Validate())
//   - Type must be a valid enum value (0-3)
//   - Inheritance and Metadata must be nil for files (Type indicates file)
//   - Inheritance and Metadata can be non-nil for directories
//
// Returns an error if any validation check fails.
func (pme *PathMetadataEntry) Validate() error {
	// Validate PathEntry
	if err := pme.Path.Validate(); err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, "path entry validation failed", pkgerrors.ValidationErrorContext{
			Field: "Path",
			Value: pme.Path.Path,
		})
	}

	// Validate Type is valid enum value
	if pme.Type > PathMetadataTypeDirectorySymlink {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid path metadata type", nil, pkgerrors.ValidationErrorContext{
			Field:    "Type",
			Value:    pme.Type,
			Expected: "0-3 (PathMetadataTypeFile, PathMetadataTypeDirectory, PathMetadataTypeFileSymlink, PathMetadataTypeDirectorySymlink)",
		})
	}

	// Validate Inheritance and Metadata are nil for files
	isFile := pme.Type == PathMetadataTypeFile || pme.Type == PathMetadataTypeFileSymlink
	if isFile {
		if pme.Inheritance != nil {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "inheritance must be nil for file paths", nil, pkgerrors.ValidationErrorContext{
				Field:    "Inheritance",
				Value:    pme.Inheritance,
				Expected: "nil for file paths",
			})
		}
		if pme.Metadata != nil {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "metadata must be nil for file paths", nil, pkgerrors.ValidationErrorContext{
				Field:    "Metadata",
				Value:    pme.Metadata,
				Expected: "nil for file paths",
			})
		}
	}

	return nil
}
