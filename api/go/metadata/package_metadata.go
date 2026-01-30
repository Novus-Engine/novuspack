// This file implements the PackageMetadata structure providing comprehensive
// package metadata including all package information plus detailed file and
// metadata file contents. It contains the PackageMetadata type definition
// and NewPackageMetadata constructor as specified in api_core.md Section 1.1.6
// and api_metadata.md Section 7.5.
//
// Specification: api_core.md: 1.1.6 GetMetadata Method Contract
// Specification: api_metadata.md: 1. Comment Management

// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This file contains the PackageMetadata structure which provides comprehensive
// metadata including all package information plus detailed file and metadata file contents.
package metadata

// PackageMetadata represents comprehensive package metadata including all package
// information plus detailed file and metadata file contents.
//
// PackageMetadata extends PackageInfo with detailed file entries, path metadata entries,
// and special metadata files. This structure is returned by GetMetadata() and provides
// a complete view of all package metadata that was eagerly loaded by OpenPackage.
//
// This method MUST NOT perform additional disk I/O or parsing beyond what OpenPackage
// already loaded. All data in PackageMetadata comes from already-loaded package state.
//
// Specification: api_core.md: 1.1.6.5 PackageMetadata Contents
// Specification: api_metadata.md: 1. Comment Management
type PackageMetadata struct {
	// Embed PackageInfo for basic package information
	*PackageInfo

	// Detailed file entries (all files in the package)
	FileEntries []*FileEntry

	// Path metadata entries (path hierarchy and metadata)
	PathMetadataEntries []*PathMetadataEntry

	// Special metadata files (type 65000-65535)
	// Maps file type to FileEntry for special metadata files
	SpecialFiles map[uint16]*FileEntry
}

// NewPackageMetadata creates a new PackageMetadata with default values.
//
// NewPackageMetadata initializes all fields including creating a new PackageInfo
// via NewPackageInfo(). This factory function ensures consistent initialization
// of PackageMetadata instances.
//
// Returns:
//   - *PackageMetadata: A new PackageMetadata instance with default values
//
// Specification: api_metadata.md: 1. Comment Management
func NewPackageMetadata() *PackageMetadata {
	return &PackageMetadata{
		PackageInfo:         NewPackageInfo(),
		FileEntries:         make([]*FileEntry, 0),
		PathMetadataEntries: make([]*PathMetadataEntry, 0),
		SpecialFiles:        make(map[uint16]*FileEntry),
	}
}
