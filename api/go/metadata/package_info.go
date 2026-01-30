// This file implements the PackageInfo structure providing comprehensive package
// information and metadata. It contains the PackageInfo type definition, NewPackageInfo
// constructor, and methods for retrieving package information. This file should
// contain all code related to package information structures as specified in
// api_metadata.md Section 7.
//
// Specification: api_metadata.md: 1. Comment Management

// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and api_metadata.md.
package metadata

import (
	"time"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/signatures"
)

// PackageInfo represents comprehensive metadata about a package.
//
// PackageInfo contains detailed information about a package including file counts,
// compression data, signatures, security information, and timestamps. This structure
// is returned by GetInfo() and provides a complete snapshot of package metadata.
//
// Specification: api_metadata.md: 7.1 PackageInfo Structure
type PackageInfo struct {
	// Basic Package Information
	FormatVersion         uint32 // Package format version
	FileCount             int    // Number of files in the package
	FilesUncompressedSize int64  // Total uncompressed size of all files
	FilesCompressedSize   int64  // Total compressed size of all files

	// Package Identity
	VendorID uint32 // Vendor/platform identifier
	AppID    uint64 // Application identifier

	// Package Comment
	HasComment bool   // Whether package has a comment
	Comment    string // Actual package comment content

	// Digital Signatures (Multiple Signatures Support)
	HasSignatures  bool                       // Whether package has any signatures
	SignatureCount int                        // Number of signatures in the package
	Signatures     []signatures.SignatureInfo // Detailed signature information

	// Security Information
	SecurityLevel SecurityLevel // Overall security level
	IsImmutable   bool          // Whether package is immutable (signed)

	// Timestamps
	Created  time.Time // Package creation timestamp
	Modified time.Time // Package modification timestamp

	// Version Tracking
	PackageDataVersion uint32 // Tracks changes to package data content (file additions, removals, modifications)
	MetadataVersion    uint32 // Tracks changes to package metadata (comment, identity changes)

	// Package Features
	HasMetadataFiles  bool // Whether package has metadata files
	HasPerFileTags    bool // Whether package has per-file tags (path metadata with properties)
	HasExtendedAttrs  bool // Whether package has extended attributes (filesystem metadata)
	HasEncryptedData  bool // Whether package contains encrypted files
	HasCompressedData bool // Whether package contains compressed files
	IsMetadataOnly    bool // Whether package contains only metadata files (no content)

	// Package Compression
	PackageCompression      uint8   // Package compression type (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
	IsPackageCompressed     bool    // Whether the entire package is compressed
	PackageOriginalSize     int64   // Original package size before compression (0 if not compressed)
	PackageCompressedSize   int64   // Compressed package size (0 if not compressed)
	PackageCompressionRatio float64 // Compression ratio (0.0-1.0, 0.0 if not compressed)
}

// NewPackageInfo creates a new PackageInfo with default values.
//
// NewPackageInfo initializes all fields to their zero values or spec defaults,
// providing a consistent starting point for package metadata. This helper
// ensures that all PackageInfo initialization is centralized and aligned
// with specification defaults.
//
// Returns:
//   - *PackageInfo: A new PackageInfo instance with default values
//
// Specification: api_metadata.md: 7.1 PackageInfo Structure
func NewPackageInfo() *PackageInfo {
	return &PackageInfo{
		FormatVersion:           fileformat.FormatVersion,
		FileCount:               0,
		FilesUncompressedSize:   0,
		FilesCompressedSize:     0,
		VendorID:                0,
		AppID:                   0,
		HasComment:              false,
		Comment:                 "",
		HasSignatures:           false,
		SignatureCount:          0,
		Signatures:              []signatures.SignatureInfo{},
		SecurityLevel:           SecurityLevelNone,
		IsImmutable:             false,
		Created:                 time.Time{}, // Zero time, will be set on Create
		Modified:                time.Time{}, // Zero time, will be set on Create
		PackageDataVersion:      0,
		MetadataVersion:         0,
		HasMetadataFiles:        false,
		HasEncryptedData:        false,
		HasCompressedData:       false,
		IsMetadataOnly:          false,
		PackageCompression:      0,
		IsPackageCompressed:     false,
		PackageOriginalSize:     0,
		PackageCompressedSize:   0,
		PackageCompressionRatio: 0.0,
	}
}
