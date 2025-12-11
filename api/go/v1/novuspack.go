// Package novuspack provides the core NovusPack file format implementation.
//
// This package implements the NovusPack (.npk) file format as specified in
// docs/tech_specs/package_file_format.md.
//
// This root package re-exports all types, constants, and functions from the
// domain-specific subpackages (fileformat, metadata, signatures) to provide
// a clean, unified API accessible through a single import:
//
//	import "github.com/novus-engine/novuspack/api/go/v1"
package novuspack

import (
	"github.com/novus-engine/novuspack/api/go/v1/fileformat"
	"github.com/novus-engine/novuspack/api/go/v1/metadata"
	"github.com/novus-engine/novuspack/api/go/v1/signatures"
)

// Re-export types from fileformat
type (
	PackageHeader     = fileformat.PackageHeader
	FileIndex         = fileformat.FileIndex
	IndexEntry        = fileformat.IndexEntry
	FileEntry         = fileformat.FileEntry
	PathEntry         = fileformat.PathEntry
	HashEntry         = fileformat.HashEntry
	OptionalDataEntry = fileformat.OptionalDataEntry
)

// Re-export types from metadata
type (
	PackageComment = metadata.PackageComment
)

// Re-export types from signatures
type (
	Signature = signatures.Signature
)

// Re-export constants from fileformat
const (
	// Package magic and version
	NPKMagic          = fileformat.NPKMagic
	FormatVersion     = fileformat.FormatVersion
	PackageHeaderSize = fileformat.PackageHeaderSize

	// Compression types
	CompressionNone = fileformat.CompressionNone
	CompressionZstd = fileformat.CompressionZstd
	CompressionLZ4  = fileformat.CompressionLZ4
	CompressionLZMA = fileformat.CompressionLZMA

	// Encryption types
	EncryptionNone        = fileformat.EncryptionNone
	EncryptionAES256GCM   = fileformat.EncryptionAES256GCM
	EncryptionQuantumSafe = fileformat.EncryptionQuantumSafe

	// Hash types
	HashTypeSHA256   = fileformat.HashTypeSHA256
	HashTypeSHA512   = fileformat.HashTypeSHA512
	HashTypeBLAKE3   = fileformat.HashTypeBLAKE3
	HashTypeXXH3     = fileformat.HashTypeXXH3
	HashTypeBLAKE2b  = fileformat.HashTypeBLAKE2b
	HashTypeBLAKE2s  = fileformat.HashTypeBLAKE2s
	HashTypeSHA3_256 = fileformat.HashTypeSHA3_256
	HashTypeSHA3_512 = fileformat.HashTypeSHA3_512
	HashTypeCRC32    = fileformat.HashTypeCRC32
	HashTypeCRC64    = fileformat.HashTypeCRC64

	// Hash purposes
	HashPurposeContentVerification = fileformat.HashPurposeContentVerification
	HashPurposeDeduplication       = fileformat.HashPurposeDeduplication
	HashPurposeIntegrity           = fileformat.HashPurposeIntegrity
	HashPurposeFastLookup          = fileformat.HashPurposeFastLookup
	HashPurposeErrorDetection      = fileformat.HashPurposeErrorDetection

	// Optional data types
	OptionalDataTagsData              = fileformat.OptionalDataTagsData
	OptionalDataPathEncoding          = fileformat.OptionalDataPathEncoding
	OptionalDataPathFlags             = fileformat.OptionalDataPathFlags
	OptionalDataCompressionDictionary = fileformat.OptionalDataCompressionDictionary
	OptionalDataSolidGroupID          = fileformat.OptionalDataSolidGroupID
	OptionalDataFileSystemFlags       = fileformat.OptionalDataFileSystemFlags
	OptionalDataWindowsAttributes     = fileformat.OptionalDataWindowsAttributes
	OptionalDataExtendedAttributes    = fileformat.OptionalDataExtendedAttributes
	OptionalDataACL                   = fileformat.OptionalDataACL

	// Package feature flags
	FlagHasSignatures      = fileformat.FlagHasSignatures
	FlagHasCompressedFiles = fileformat.FlagHasCompressedFiles
	FlagHasEncryptedFiles  = fileformat.FlagHasEncryptedFiles
	FlagHasExtendedAttrs   = fileformat.FlagHasExtendedAttrs
	FlagHasPackageComment  = fileformat.FlagHasPackageComment
	FlagHasPerFileTags     = fileformat.FlagHasPerFileTags
	FlagHasSpecialMetadata = fileformat.FlagHasSpecialMetadata
	FlagMetadataOnly       = fileformat.FlagMetadataOnly

	// Flags field bit masks
	FlagsMaskFeatures        = fileformat.FlagsMaskFeatures
	FlagsMaskCompressionType = fileformat.FlagsMaskCompressionType
	FlagsMaskReserved1       = fileformat.FlagsMaskReserved1
	FlagsMaskReserved2       = fileformat.FlagsMaskReserved2

	// Flags field bit shifts
	FlagsShiftCompressionType = fileformat.FlagsShiftCompressionType
	FlagsShiftReserved1       = fileformat.FlagsShiftReserved1
	FlagsShiftReserved2       = fileformat.FlagsShiftReserved2

	// Size constants
	FileEntryFixedSize = fileformat.FileEntryFixedSize
	IndexEntrySize     = fileformat.IndexEntrySize

	// VendorID constants
	VendorIDNone        = fileformat.VendorIDNone
	VendorIDSteam       = fileformat.VendorIDSteam
	VendorIDEpic        = fileformat.VendorIDEpic
	VendorIDGOG         = fileformat.VendorIDGOG
	VendorIDItch        = fileformat.VendorIDItch
	VendorIDHumble      = fileformat.VendorIDHumble
	VendorIDMicrosoft   = fileformat.VendorIDMicrosoft
	VendorIDPlayStation = fileformat.VendorIDPlayStation
	VendorIDXbox        = fileformat.VendorIDXbox
	VendorIDNintendo    = fileformat.VendorIDNintendo
	VendorIDUnity       = fileformat.VendorIDUnity
	VendorIDUnreal      = fileformat.VendorIDUnreal
	VendorIDGitHub      = fileformat.VendorIDGitHub
	VendorIDGitLab      = fileformat.VendorIDGitLab
)

// Re-export constants from signatures
const (
	SignatureTypeMLDSA  = signatures.SignatureTypeMLDSA
	SignatureTypeSLHDSA = signatures.SignatureTypeSLHDSA
	SignatureTypePGP    = signatures.SignatureTypePGP
	SignatureTypeX509   = signatures.SignatureTypeX509
)

// Re-export constants from metadata
const (
	MaxCommentLength = metadata.MaxCommentLength
)

// Re-export functions from metadata
var (
	NewPackageComment = metadata.NewPackageComment
)

// Re-export functions from fileformat
var (
	NewPackageHeader = fileformat.NewPackageHeader
	NewFileEntry     = fileformat.NewFileEntry
	NewFileIndex     = fileformat.NewFileIndex
)

// Re-export functions from signatures
var (
	NewSignature = signatures.NewSignature
)
