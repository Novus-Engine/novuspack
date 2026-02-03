// This file contains type definitions and option structures used by the
// Package interfaces. It includes FileInfo, AddFileOptions, CompressionType,
// and CreateOptions structures. These types are API-level types that are not
// domain-specific and remain in the main package. This file should contain
// only type definitions and constants used across the package API.
//
// Specification: api_file_mgmt_index.md: 1. File Management Document Map

// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains placeholder types and options structures used by the
// Package interfaces. These types are not domain-specific and remain in
// the main package.
package novus_package

import (
	"os"

	"github.com/novus-engine/novuspack/api/go/generics"
)

// FileInfo represents information about a file in the package.
//
// FileInfo is returned by ListFiles() and provides lightweight file information for
// package contents, optimized for listing operations and filtering without requiring
// full FileEntry access.
//
// All fields are populated from FileEntry static section (no variable-length data).
// Paths have leading `/` stripped via internal.ToDisplayPath for user presentation.
//
// Specification: api_core.md: 1.2.4 FileInfo Structure
type FileInfo struct {
	// Basic Identification
	PrimaryPath  string   // Primary display path (leading '/' removed, first path lexicographically)
	Paths        []string // All paths for this file (aliases/hard links, leading '/' removed)
	FileID       uint64   // Unique file identifier
	FileType     uint16   // File type identifier (0-64999: content, 65000-65535: special metadata)
	FileTypeName string   // Human-readable file type name (e.g., "Texture", "Audio", "Unknown")

	// Size Information
	Size       int64 // Original file size in bytes (before compression/encryption)
	StoredSize int64 // Actual stored size in bytes (after compression/encryption)

	// Processing Status
	IsCompressed    bool  // Whether file is compressed
	IsEncrypted     bool  // Whether file is encrypted
	CompressionType uint8 // Compression algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)

	// Content Verification
	RawChecksum    uint32 // CRC32 checksum of original content
	StoredChecksum uint32 // CRC32 checksum of stored content (after compression/encryption)

	// Multi-Path Support
	PathCount uint16 // Number of paths (aliases) for this file entry

	// Version Tracking
	FileVersion     uint32 // File content version
	MetadataVersion uint32 // File metadata version

	// Metadata Indicators
	HasTags bool // Whether file has custom tags/metadata
}

// PathHandling specifies how to handle multiple paths pointing to the same content.
//
// This type is used by AddFileOptions to control behavior during deduplication
// and file addition operations.
//
// Specification: api_file_mgmt_addition.md: 2.7 PathHandling Type
type PathHandling uint8

const (
	// PathHandlingDefault uses the package default (Package.DefaultPathHandling).
	PathHandlingDefault PathHandling = 0

	// PathHandlingHardLinks stores multiple paths as hard links (current behavior, backward compatible).
	PathHandlingHardLinks PathHandling = 1

	// PathHandlingSymlinks converts additional paths to symlinks.
	PathHandlingSymlinks PathHandling = 2

	// PathHandlingPreserve preserves original filesystem behavior (detect and respect symlinks/hardlinks).
	PathHandlingPreserve PathHandling = 3
)

// AddFileOptions represents options for adding files to a package.
//
// AddFileOptions provides comprehensive configuration for file addition operations,
// including path determination, conflict handling, filesystem metadata preservation,
// and compression/encryption settings.
//
// All fields use Option[T] types for optional configuration. Unset options use
// implementation-defined defaults.
//
// Specification: api_file_mgmt_addition.md: 2.8 AddFileOptions Struct
type AddFileOptions struct {
	// Path determination options
	StoredPath    generics.Option[string] // Explicit path override for storage location
	BasePath      generics.Option[string] // Filesystem base path to strip from source paths
	PreserveDepth generics.Option[int]    // Number of directory levels to preserve (default: 1)
	FlattenPaths  generics.Option[bool]   // Store all files at package root (default: false)

	// Conflict and deduplication options
	AllowOverwrite generics.Option[bool] // Allow overwriting existing files at same path (default: false)
	AllowDuplicate generics.Option[bool] // Skip deduplication, always create new FileEntry (default: false)
	FollowSymlinks generics.Option[bool] // Follow symbolic links (default: true)

	// Path handling for duplicate content
	PathHandling        PathHandling                           // How to handle multiple paths pointing to the same content (default: PathHandlingDefault)
	PrimaryPathSelector generics.Option[func([]string) string] // Custom selector for primary path when converting to symlinks (default: nil, uses lexicographic ordering)

	// Filesystem metadata preservation options
	PreservePermissions   generics.Option[bool] // Preserve Unix permissions (default: true)
	PreserveOwnership     generics.Option[bool] // Preserve UID/GID (default: false)
	PreserveACL           generics.Option[bool] // Preserve ACLs (default: false)
	PreserveExtendedAttrs generics.Option[bool] // Preserve extended attributes (default: false)

	// File processing options
	Compress         generics.Option[bool]                 // Enable compression (default: false)
	CompressionType  generics.Option[uint8]                // Compression algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
	CompressionLevel generics.Option[int]                  // Compression level 1-9 (default: 6)
	FileType         generics.Option[uint16]               // File type identifier (override auto-detection)
	Tags             generics.Option[[]*generics.Tag[any]] // Per-file tags

	// Encryption options (deferred to Priority 6)
	EncryptionKey generics.Option[*EncryptionKey] // Encryption key (enables encryption when set)

	// Multi-stage transformation pipeline options
	MaxTransformStages      generics.Option[int]  // Maximum transformation stages per pipeline (default: 10)
	ValidateProcessingState generics.Option[bool] // Enable ProcessingState validation (default: false)

	// Pattern-specific options (for AddFilePattern and AddDirectory)
	ExcludePatterns  generics.Option[[]string]           // Patterns to exclude from pattern operations
	MaxFileSize      generics.Option[int64]              // Maximum file size for pattern operations (0=no limit)
	PreservePaths    generics.Option[bool]               // Preserve directory structure in pattern operations (default: true)
	ProgressCallback generics.Option[func(int64, int64)] // Progress callback (bytesProcessed, totalBytes)
}

// RemoveDirectoryOptions configures directory removal behavior.
//
// Specification: api_file_mgmt_removal.md: 4.4 RemoveDirectoryOptions Struct
type RemoveDirectoryOptions struct {
	// Recursive controls whether to remove files in subdirectories (default: true).
	Recursive generics.Option[bool]

	// Pattern filters which files to remove (default: all files).
	Pattern generics.Option[string]

	// RemoveEmptyDirs controls whether to remove directory metadata entries
	// when all files in a directory are removed (default: true).
	RemoveEmptyDirs generics.Option[bool]
}

// CompressionType represents the type of compression to use.
//
// This type is used by package write operations.
// Use constants from fileformat package (CompressionNone, CompressionZstd, etc.)
// via the novuspack package re-exports.
type CompressionType uint8

// EncryptionAlgorithm represents the encryption algorithm identifier.
//
// Specification: api_security.md: 3.1.1 EncryptionAlgorithm Type
type EncryptionAlgorithm int

const (
	EncryptionAlgorithmNone EncryptionAlgorithm = iota
	EncryptionAlgorithmAES256GCM
	EncryptionAlgorithmChaCha20Poly1305
	EncryptionAlgorithmMLKEM512
	EncryptionAlgorithmMLKEM768
	EncryptionAlgorithmMLKEM1024
)

// EncryptionType is a v1 alias of EncryptionAlgorithm.
//
// Specification: api_security.md: 3.1.2 EncryptionType Alias
type EncryptionType = EncryptionAlgorithm

const (
	EncryptionNone             EncryptionType = EncryptionAlgorithmNone
	EncryptionAES256GCM        EncryptionType = EncryptionAlgorithmAES256GCM
	EncryptionChaCha20Poly1305 EncryptionType = EncryptionAlgorithmChaCha20Poly1305
	EncryptionMLKEM512         EncryptionType = EncryptionAlgorithmMLKEM512
	EncryptionMLKEM768         EncryptionType = EncryptionAlgorithmMLKEM768
	EncryptionMLKEM1024        EncryptionType = EncryptionAlgorithmMLKEM1024
)

// EncryptionKey represents an encryption key for file encryption.
//
// This is a placeholder type for Priority 6 (Encryption Support).
// The actual implementation will be defined when encryption features are added.
//
// TODO: Priority 6 - Define complete EncryptionKey structure with:
//   - Key data
//   - Algorithm identifier
//   - Key metadata (creation time, expiry, etc.)
type EncryptionKey struct {
	// Placeholder - will be implemented in Priority 6
}

// CreateOptions represents options for creating a package.
//
// CreateOptions allows configuring package creation with metadata,
// comments, and identifiers.
//
// Specification: api_basic_operations.md: 1. Context Integration
type CreateOptions struct {
	Comment     string      // Initial package comment
	VendorID    uint32      // Vendor identifier
	AppID       uint64      // Application identifier
	Permissions os.FileMode // File permissions (default: 0644)
}
