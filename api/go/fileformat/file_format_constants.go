// This file contains all constants related to the NovusPack file format including
// magic numbers, version numbers, compression types, encryption types, hash types,
// optional data types, package flags, and vendor IDs. This file should contain
// only constant definitions used throughout the fileformat package and re-exported
// by the main novuspack package.
//
// Specification: package_file_format.md: 1 `.nvpk` File Format Overview

// Package novuspack provides the core NovusPack file format implementation.
//
// This package implements the NovusPack (.nvpk) file format as specified in
// package_file_format.md.
package fileformat

// Package magic number and version constants
const (
	// NVPKMagic is the magic number identifying a NovusPack file ("NVPK")
	// Specification: package_file_format.md: 2.1 Header Structure
	NVPKMagic = 0x4E56504B

	// FormatVersion is the current format version
	// Specification: package_file_format.md: 2.1 Header Structure
	FormatVersion = 1

	// PackageHeaderSize is the fixed size of the package header in bytes
	// Specification: package_file_format.md: 2.1 Header Structure
	PackageHeaderSize = 112
)

// Compression type constants
// Specification: package_file_format.md: 4.1.1.3 Compression and Encryption Types
const (
	CompressionNone = 0 // No compression
	CompressionZstd = 1 // Zstandard compression
	CompressionLZ4  = 2 // LZ4 compression
	CompressionLZMA = 3 // LZMA compression
)

// Encryption type constants
// Specification: package_file_format.md: 4.1.1.3 Compression and Encryption Types
const (
	EncryptionNone        = 0x00 // No encryption
	EncryptionAES256GCM   = 0x01 // AES-256-GCM encryption
	EncryptionQuantumSafe = 0x02 // Quantum-safe encryption (ML-KEM + ML-DSA)
)

// Hash algorithm type constants
// Specification: package_file_format.md: 4.1.5 Hash Algorithm Support
const (
	HashTypeSHA256   = 0x00 // SHA-256 (32 bytes)
	HashTypeSHA512   = 0x01 // SHA-512 (64 bytes)
	HashTypeBLAKE3   = 0x02 // BLAKE3 (32 bytes)
	HashTypeXXH3     = 0x03 // XXH3 (8 bytes)
	HashTypeBLAKE2b  = 0x04 // BLAKE2b (64 bytes)
	HashTypeBLAKE2s  = 0x05 // BLAKE2s (32 bytes)
	HashTypeSHA3_256 = 0x06 // SHA-3-256 (32 bytes)
	HashTypeSHA3_512 = 0x07 // SHA-3-512 (64 bytes)
	HashTypeCRC32    = 0x08 // CRC32 (4 bytes)
	HashTypeCRC64    = 0x09 // CRC64 (8 bytes)
)

// Hash purpose constants
// Specification: package_file_format.md: 4.1.5 Hash Algorithm Support
const (
	HashPurposeContentVerification = 0x00 // Verify file content integrity
	HashPurposeDeduplication       = 0x01 // Identify duplicate content
	HashPurposeIntegrity           = 0x02 // General integrity verification
	HashPurposeFastLookup          = 0x03 // Quick content identification
	HashPurposeErrorDetection      = 0x04 // Detect data corruption
)

// Optional data type constants
// Specification: package_file_format.md: 4.1.4.4 Optional Data
const (
	OptionalDataTagsData              = 0x00 // Per-file tags data
	OptionalDataPathEncoding          = 0x01 // Path encoding type for this file
	OptionalDataPathFlags             = 0x02 // Path handling flags for this file
	OptionalDataCompressionDictionary = 0x03 // Dictionary identifier for solid compression
	OptionalDataSolidGroupID          = 0x04 // Solid compression group identifier
	OptionalDataFileSystemFlags       = 0x05 // File system specific flags
	OptionalDataWindowsAttributes     = 0x06 // Windows file attributes
	OptionalDataExtendedAttributes    = 0x07 // Unix extended attributes
	OptionalDataACL                   = 0x08 // Access Control List data
)

// Package feature flags (Flags field bit positions)
// Specification: package_file_format.md: 2.5 Package Features Flags
const (
	FlagHasSignatures      = 1 << 0 // Bit 0: Has digital signatures
	FlagHasCompressedFiles = 1 << 1 // Bit 1: Has compressed files
	FlagHasEncryptedFiles  = 1 << 2 // Bit 2: Has encrypted files
	FlagHasExtendedAttrs   = 1 << 3 // Bit 3: Has extended attributes
	FlagHasPackageComment  = 1 << 4 // Bit 4: Has package comment
	FlagHasPerFileTags     = 1 << 5 // Bit 5: Has per-file tags
	FlagHasSpecialMetadata = 1 << 6 // Bit 6: Has special metadata files
	FlagMetadataOnly       = 1 << 7 // Bit 7: Metadata-only package
)

// Flags field bit masks
// Specification: package_file_format.md: 2.5.1 Flags Field Encoding
const (
	FlagsMaskFeatures        = 0x000000FF // Bits 0-7: Package features
	FlagsMaskCompressionType = 0x0000FF00 // Bits 8-15: Package compression type
	FlagsMaskReserved1       = 0x00FF0000 // Bits 16-23: Reserved for future use
	FlagsMaskReserved2       = 0xFF000000 // Bits 24-31: Reserved for future use
)

// Flags field bit shift values
const (
	FlagsShiftCompressionType = 8  // Shift for compression type bits
	FlagsShiftReserved1       = 16 // Shift for reserved bits 16-23
	FlagsShiftReserved2       = 24 // Shift for reserved bits 24-31
)

// Size constants for common data structures
const (
	FileEntryFixedSize = 64 // Fixed portion of FileEntry
	IndexEntrySize     = 16 // Size of FileIndex entry (FileID + Offset)
)

// VendorID example constants
// Specification: package_file_format.md: 2.3.1 VendorID Example Mappings
const (
	VendorIDNone        = 0x00000000 // No vendor association
	VendorIDSteam       = 0x53544541 // STEA
	VendorIDEpic        = 0x45504943 // EPIC
	VendorIDGOG         = 0x474F4720 // GOG (with space)
	VendorIDItch        = 0x49544348 // ITCH
	VendorIDHumble      = 0x48554D42 // HUMB
	VendorIDMicrosoft   = 0x4D494352 // MICR
	VendorIDPlayStation = 0x50534E59 // PSNY
	VendorIDXbox        = 0x58424F58 // XBOX
	VendorIDNintendo    = 0x4E54444F // NTDO
	VendorIDUnity       = 0x554E4954 // UNIT
	VendorIDUnreal      = 0x554E5245 // UNRE
	VendorIDGitHub      = 0x47495448 // GITH
	VendorIDGitLab      = 0x4749544C // GITL
)
