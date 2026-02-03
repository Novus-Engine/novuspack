// This file provides type and constant re-exports from subpackages to create
// a unified API. It re-exports types from fileformat, metadata, signatures,
// generics, and pkgerrors packages so users can access everything through a
// single import. This file should contain only re-export declarations and
// convenience wrapper functions, not implementation code.
//
// Specification: api_core.md: 0 Overview

// Package novuspack provides the core NovusPack file format implementation.
//
// This package implements the NovusPack (.nvpk) file format as specified in
// package_file_format.md.
//
// This root package re-exports all types, constants, and functions from the
// domain-specific subpackages (fileformat, metadata, signatures) to provide
// a clean, unified API accessible through a single import:
//
//	import "github.com/novus-engine/novuspack/api/go"
package novuspack

import (
	"context"
	"io"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/novus_package"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
	"github.com/novus-engine/novuspack/api/go/signatures"
)

// Re-export core package interfaces from novus_package
type (
	Package        = novus_package.Package
	PackageBuilder = novus_package.PackageBuilder
)

// Re-export package operation types from novus_package
type (
	FileInfo               = novus_package.FileInfo
	AddFileOptions         = novus_package.AddFileOptions
	RemoveDirectoryOptions = novus_package.RemoveDirectoryOptions
	CreateOptions          = novus_package.CreateOptions
	CompressionType        = novus_package.CompressionType
	EncryptionType         = novus_package.EncryptionType
)

// Re-export types from pkgerrors
type (
	ErrorType    = pkgerrors.ErrorType
	PackageError = pkgerrors.PackageError
)

// Re-export types from fileformat
type (
	PackageHeader = fileformat.PackageHeader
	FileIndex     = fileformat.FileIndex
	IndexEntry    = fileformat.IndexEntry
)

// Re-export types from metadata
type (
	FileEntry         = metadata.FileEntry
	HashEntry         = metadata.HashEntry
	OptionalDataEntry = metadata.OptionalDataEntry
	ProcessingState   = metadata.ProcessingState
)

// Re-export types from generics (shared types)
type (
	PathEntry = generics.PathEntry
)

// Re-export types from metadata
type (
	PackageComment = metadata.PackageComment
	PackageInfo    = metadata.PackageInfo
	SecurityLevel  = metadata.SecurityLevel
)

// Re-export types from signatures
type (
	Signature     = signatures.Signature
	SignatureInfo = signatures.SignatureInfo
)

// Re-export generic types from generics package
//
// Generic types provide type-safe abstractions for optional values, error handling,
// strategy patterns, and validation. These types use Go generics to provide
// compile-time type safety.
// See api_generics.md for complete documentation.
type (
	// Option provides type-safe optional configuration values.
	// Use Option[T] for optional fields instead of pointer types (*int, *string, etc.).
	Option[T any] = generics.Option[T]

	// Result provides type-safe error handling for operations that may fail.
	// Use Result[T] to encapsulate either a successful value or an error.
	// Create Result values using Ok[T] for success or Err[T] for failure.
	Result[T any] = generics.Result[T]

	// Strategy defines a generic strategy pattern for processing different data types.
	// Use Strategy[T, U] for pluggable behavior patterns (compression, encryption, etc.).
	Strategy[T any, U any] = generics.Strategy[T, U]

	// Validator defines a generic validation interface for type-safe validation.
	// Use Validator[T] for validation operations that can be composed.
	Validator[T any] = generics.Validator[T]

	// ValidationRule represents a single validation rule with a predicate function.
	// Use ValidationRule[T] to create validators from simple predicate functions.
	ValidationRule[T any] = generics.ValidationRule[T]

	// Tag represents a type-safe tag with a specific value type.
	// Use Tag[T] for type-safe tag operations where T is the type of the tag value.
	// All tags are stored and accessed as typed tags for compile-time type safety.
	Tag[T any] = generics.Tag[T]

	// TagValueType represents the type of a tag value.
	TagValueType = generics.TagValueType
)

// Re-export constants from errors
const (
	ErrTypeValidation  = pkgerrors.ErrTypeValidation
	ErrTypeIO          = pkgerrors.ErrTypeIO
	ErrTypeSecurity    = pkgerrors.ErrTypeSecurity
	ErrTypeUnsupported = pkgerrors.ErrTypeUnsupported
	ErrTypeContext     = pkgerrors.ErrTypeContext
	ErrTypeCorruption  = pkgerrors.ErrTypeCorruption
)

// Re-export constants from metadata
const (
	MaxCommentLength = metadata.MaxCommentLength
)

// Re-export constants from fileformat
const (
	// Package magic and version
	NVPKMagic         = fileformat.NVPKMagic
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
	FileEntryFixedSize = metadata.FileEntryFixedSize
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

// Re-export constants from generics (tag value types)
//
// Tag value types define the data type of tag values for type-safe tag operations.
// These constants are used when creating tags with NewTag to specify the value type.
const (
	// TagValueTypeString represents a string value (0x00).
	TagValueTypeString = generics.TagValueTypeString

	// TagValueTypeInteger represents a 64-bit signed integer value (0x01).
	TagValueTypeInteger = generics.TagValueTypeInteger

	// TagValueTypeFloat represents a 64-bit floating point number value (0x02).
	TagValueTypeFloat = generics.TagValueTypeFloat

	// TagValueTypeBoolean represents a boolean value (0x03).
	TagValueTypeBoolean = generics.TagValueTypeBoolean

	// TagValueTypeJSON represents a JSON-encoded object or array value (0x04).
	TagValueTypeJSON = generics.TagValueTypeJSON

	// TagValueTypeYAML represents a YAML-encoded data value (0x05).
	TagValueTypeYAML = generics.TagValueTypeYAML

	// TagValueTypeStringList represents a comma-separated list of strings value (0x06).
	TagValueTypeStringList = generics.TagValueTypeStringList

	// TagValueTypeUUID represents a UUID string value (0x07).
	TagValueTypeUUID = generics.TagValueTypeUUID

	// TagValueTypeHash represents a hash or checksum string value (0x08).
	TagValueTypeHash = generics.TagValueTypeHash

	// TagValueTypeVersion represents a semantic version string value (0x09).
	TagValueTypeVersion = generics.TagValueTypeVersion

	// TagValueTypeTimestamp represents an ISO8601 timestamp value (0x0A).
	TagValueTypeTimestamp = generics.TagValueTypeTimestamp

	// TagValueTypeURL represents a URL string value (0x0B).
	TagValueTypeURL = generics.TagValueTypeURL

	// TagValueTypeEmail represents an email address value (0x0C).
	TagValueTypeEmail = generics.TagValueTypeEmail

	// TagValueTypePath represents a file system path value (0x0D).
	TagValueTypePath = generics.TagValueTypePath

	// TagValueTypeMimeType represents a MIME type string value (0x0E).
	TagValueTypeMimeType = generics.TagValueTypeMimeType

	// TagValueTypeLanguage represents a language code (ISO 639-1) value (0x0F).
	TagValueTypeLanguage = generics.TagValueTypeLanguage

	// TagValueTypeNovusPackMetadata represents a NovusPack special metadata file reference value (0x10).
	TagValueTypeNovusPackMetadata = generics.TagValueTypeNovusPackMetadata
)

// Re-export functions from metadata
var (
	NewPackageComment = metadata.NewPackageComment
	NewFileEntry      = metadata.NewFileEntry
)

// Re-export functions from fileformat
var (
	NewPackageHeader = fileformat.NewPackageHeader
	NewFileIndex     = fileformat.NewFileIndex
)

// Re-export functions from signatures
var (
	NewSignature = signatures.NewSignature
)

// Re-export functions from pkgerrors
var (
	NewPackageError = pkgerrors.NewPackageError[struct{}]
	WrapError       = pkgerrors.WrapError
	IsPackageError  = pkgerrors.IsPackageError
	GetErrorType    = pkgerrors.GetErrorType
)

// Re-export generic functions from errors
// Note: Generic functions cannot be re-exported via variables, so they must be
// accessed directly from the errors package: pkgerrors.AddErrorContext[T](...)
// or imported as: import "github.com/novus-engine/novuspack/api/go/pkgerrors"

// ValidateWith validates a single value using a validator.
//
// ValidateWith[T] applies the given validator to the value and returns
// any validation error that occurs. Errors are returned as *PackageError
// with ValidationErrorContext for structured error handling.
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	validator := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) > 0 },
//	    Message:   "string cannot be empty",
//	}
//	err := ValidateWith(ctx, "", validator)
func ValidateWith[T any](ctx context.Context, value T, validator Validator[T]) error {
	return generics.ValidateWith(ctx, value, validator)
}

// ValidateAll validates multiple values using a validator.
//
// ValidateAll[T] applies the given validator to each value and returns
// a slice of errors, one for each validation failure. If all validations
// pass, returns an empty slice. Errors are returned as *PackageError
// with ValidationErrorContext for structured error handling.
//
// Type Parameters:
//   - T: The type of values to validate
//
// Example:
//
//	validator := &ValidationRule[int]{
//	    Predicate: func(n int) bool { return n > 0 },
//	    Message:   "value must be positive",
//	}
//	errors := ValidateAll(ctx, []int{1, -1, 2, -2}, validator)
func ValidateAll[T any](ctx context.Context, values []T, validator Validator[T]) []error {
	return generics.ValidateAll(ctx, values, validator)
}

// ComposeValidators creates a validator that runs multiple validators.
//
// ComposeValidators[T] creates a composite validator that runs all provided
// validators in sequence. If any validator fails, the composite validator
// returns that error immediately (short-circuit evaluation).
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	validator1 := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) > 0 },
//	    Message:   "string cannot be empty",
//	}
//	validator2 := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) < 100 },
//	    Message:   "string too long",
//	}
//	composite := ComposeValidators(validator1, validator2)
//	err := composite.Validate("test")
func ComposeValidators[T any](validators ...Validator[T]) Validator[T] {
	return generics.ComposeValidators(validators...)
}

// Ok creates a successful Result with the given value.
//
// Ok[T] returns a Result[T] that represents a successful operation with the given value.
//
// Type Parameters:
//   - T: The type of the successful value
//
// Example:
//
//	result := Ok("success")
//	if result.IsOk() {
//	    value, _ := result.Unwrap()
//	    fmt.Println(value)
//	}
func Ok[T any](value T) Result[T] {
	return generics.Ok(value)
}

// Err creates a failed Result with the given error.
//
// Err[T] returns a Result[T] that represents a failed operation with the given error.
//
// Type Parameters:
//   - T: The type of the value (zero value will be used)
//
// Example:
//
//	result := Err[string](errors.New("failed"))
//	if result.IsErr() {
//	    _, err := result.Unwrap()
//	    fmt.Println(err)
//	}
func Err[T any](err error) Result[T] {
	return generics.Err[T](err)
}

// =============================================================================
// PATH AND VALIDATION UTILITIES
// =============================================================================
//
// These functions are exported from the Go API for path normalization, display
// conversion, and validation. They implement the rules in api_core 1.1.2,
// api_generics 1.3.3, and file_validation.

// NormalizePackagePath normalizes a package-internal path for consistent
// comparison and storage. Applies separator normalization, dot-segment
// canonicalization, leading slash, NFC, and path-length checks.
//
// Returns the normalized path with leading "/" or an error if the path is
// invalid or would escape the package root.
//
// Specification: api_core.md: 12.1 NormalizePackagePath Function
func NormalizePackagePath(path string) (string, error) {
	return internal.NormalizePackagePath(path)
}

// ToDisplayPath converts a stored package path (with leading "/") to display
// format by stripping the leading slash. Use for user-facing path display.
//
// Specification: api_core.md: 12.2 ToDisplayPath Function
func ToDisplayPath(storedPath string) string {
	return internal.ToDisplayPath(storedPath)
}

// ValidatePathLength checks path length against platform limits and returns
// portability warnings. Fails only if path exceeds 32,767 bytes.
//
// Returns warnings for paths over 260 (Windows), 1024 (macOS), 4096 (Linux)
// bytes; error only when over the absolute maximum.
//
// Specification: api_core.md: 12.4 ValidatePathLength Function
func ValidatePathLength(path string) ([]string, error) {
	return internal.ValidatePathLength(path)
}

// ValidatePackagePath validates a package-internal path: non-empty, no
// root-escape via dot segments, valid format. Delegates to NormalizePackagePath.
//
// Specification: api_core.md: 12.3 ValidatePackagePath Function
func ValidatePackagePath(path string) error {
	return internal.ValidatePackagePath(path)
}

// Collection operations (filtering, mapping, searching, aggregation) should use
// samber/lo directly. See docs/implementations/go/samber_lo_usage.md for guidelines.

// =============================================================================
// PACKAGE LIFECYCLE OPERATIONS
// =============================================================================

// NewPackage creates a new in-memory Package instance.
//
// This constructor creates a package in memory only and does not write to disk.
// The returned Package is in the "New" state and must be configured using Create()
// before it can be used for file operations.
//
// Use Cases:
//   - Create a new package: NewPackage() → Create()
//   - Alternative: Use OpenPackage() to open an existing package
//
// Returns:
//   - Package: A new Package instance in the "New" state
//   - error: Always returns nil for the error (reserved for future use)
//
// Example:
//
//	pkg, err := novuspack.NewPackage()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = pkg.Create(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close()
//
// Specification: api_basic_operations.md: 6.1 NewPackage Behavior
func NewPackage() (Package, error) {
	return novus_package.NewPackage()
}

// NewBuilder creates a new PackageBuilder for building packages with a fluent API.
//
// PackageBuilder provides a fluent interface for creating packages with complex
// configurations, improving code readability and reducing parameter complexity.
//
// Returns:
//   - PackageBuilder: A new builder instance with default values
//
// Example:
//
//	pkg, err := novuspack.NewBuilder().
//	    WithCompression(novuspack.CompressionZstd).
//	    WithComment("My application package").
//	    WithVendorID(0x12345678).
//	    WithAppID(0x87654321).
//	    Build(ctx)
func NewBuilder() PackageBuilder {
	return novus_package.NewBuilder()
}

// OpenPackage opens an existing package from the specified path.
//
// This function reads the package header, validates the format, and loads the
// file index. The returned Package is in the "Open" state and ready for read
// operations. The caller is responsible for calling Close() to release resources.
//
// State Transition: None → Open
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: File path to the package to open
//
// Returns:
//   - Package: The opened Package instance (interface)
//   - error: Error if file doesn't exist, format is invalid, or context is cancelled
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Path is empty, format is invalid, or magic number doesn't match
//   - ErrTypeIO: File not found, cannot open file, or read error
//
// Example:
//
//	pkg, err := novuspack.OpenPackage(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close()
//	info, err := pkg.GetInfo()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Files: %d\n", info.FileCount)
//
// Specification: api_basic_operations.md: 10. OpenPackage Function
func OpenPackage(ctx context.Context, path string) (Package, error) {
	return novus_package.OpenPackage(ctx, path)
}

// OpenPackageReadOnly opens a package in read-only mode.
//
// This function opens an existing NovusPack package file for reading only.
// It reuses OpenPackage parsing logic and returns a wrapper Package that
// enforces read-only behavior by rejecting all mutating operations.
//
// The returned Package is a wrapper type that prevents callers from
// type-asserting to the writable implementation type. All mutating
// operations return structured security errors.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: File path to the package to open
//
// Returns:
//   - Package: The opened Package instance in read-only mode
//   - error: *PackageError on failure
//
// Error Conditions:
//   - All errors from OpenPackage
//   - ErrTypeSecurity: A write or mutation operation is attempted on a read-only package
//
// Example:
//
//	pkg, err := novuspack.OpenPackageReadOnly(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close()
//	// Read operations work normally
//	data, err := pkg.ReadFile(ctx, "file.txt")
//	// Mutating operations return security errors
//	err = pkg.StageFile(ctx, "new.txt", []byte("data"), nil)
//	// err is *PackageError with Type == ErrTypeSecurity
//
// Specification: api_basic_operations.md: 11.2 OpenPackageReadOnly Function
func OpenPackageReadOnly(ctx context.Context, path string) (Package, error) {
	return novus_package.OpenPackageReadOnly(ctx, path)
}

// OpenBrokenPackage opens a package that may be invalid or partially corrupted.
//
// This function is intended for repair workflows and forensic inspection.
// It attempts to open packages even when validation fails, providing best-effort
// access to whatever data can be recovered.
//
// Behavior:
//   - If the header cannot be read, returns an error
//   - If the header is valid but the index cannot be read or is invalid:
//   - Returns a Package with an empty index
//   - The Package can be inspected and potentially repaired
//   - Read operations will fail gracefully rather than panicking
//   - Does NOT enforce the same validation guarantees as OpenPackage
//
// Use Cases:
//   - Package repair workflows
//   - Forensic inspection of corrupted packages
//   - Data recovery from partially readable packages
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: File path to the potentially broken package
//
// Returns:
//   - Package: Package instance with whatever data could be loaded
//   - error: Error if the file is completely unreadable or header invalid
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Path is empty or invalid
//   - ErrTypeIO: File not found or cannot be opened
//   - ErrTypeValidation: Header is invalid or cannot be read
//
// Example:
//
//	pkg, err := novuspack.OpenBrokenPackage(ctx, "corrupted.nvpk")
//	if err != nil {
//	    log.Fatal("Cannot open even as broken package:", err)
//	}
//	defer pkg.Close()
//	// Attempt to extract whatever data is accessible
//
// Specification: api_basic_operations.md: 12. OpenBrokenPackage Function
func OpenBrokenPackage(ctx context.Context, path string) (Package, error) {
	return novus_package.OpenBrokenPackage(ctx, path)
}

// ReadHeader reads the package header from an io.Reader.
//
// This function reads and validates the package header from any reader source,
// providing flexibility for header inspection from various sources (files,
// network streams, memory buffers, etc.).
//
// The reader must provide at least PackageHeaderSize bytes for successful reading.
// The header is validated after reading to ensure it conforms to the NovusPack format.
//
// Use Cases:
//   - Read header from an already-open file
//   - Read header from a network stream
//   - Read header from in-memory data
//   - Parse header without filesystem operations
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - reader: Input stream to read header from
//
// Returns:
//   - *fileformat.PackageHeader: The package header structure
//   - error: Error if reader doesn't provide enough data, format is invalid, or context is cancelled
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Magic number is invalid or header is malformed
//   - ErrTypeIO: Read error from the reader
//
// Example:
//
//	file, _ := os.Open("mypackage.nvpk")
//	defer file.Close()
//	header, err := novuspack.ReadHeader(ctx, file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Format Version: %d\n", header.FormatVersion)
//
// Specification: api_basic_operations.md: 18. Header Inspection
func ReadHeader(ctx context.Context, reader io.Reader) (*fileformat.PackageHeader, error) {
	return novus_package.ReadHeader(ctx, reader)
}

// ReadHeaderFromPath reads the package header from a file path.
//
// This is a convenience function that opens a file, reads the header, and closes
// the file automatically. For more control over the file handle or to read from
// other sources, use ReadHeader with an io.Reader.
//
// Use Cases:
//   - Quick format version check from a file path
//   - Header inspection without full package loading
//   - Validation of package file format
//   - Reading metadata without resource allocation
//
// Performance:
//   - Only reads the header (typically first 112 bytes)
//   - Does not load file index or entries
//   - File is automatically opened and closed
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: File path to the package
//
// Returns:
//   - *fileformat.PackageHeader: The package header structure
//   - error: Error if file doesn't exist, format is invalid, or context is cancelled
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Path is empty, magic number is invalid, or header is malformed
//   - ErrTypeIO: File not found, cannot open file, or read error
//
// Example:
//
//	header, err := novuspack.ReadHeaderFromPath(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Format Version: %d\n", header.FormatVersion)
//	fmt.Printf("Magic: 0x%08X\n", header.Magic)
//	fmt.Printf("Index Start: %d\n", header.IndexStart)
//
// Specification: api_basic_operations.md: 18. Header Inspection
func ReadHeaderFromPath(ctx context.Context, path string) (*fileformat.PackageHeader, error) {
	return novus_package.ReadHeaderFromPath(ctx, path)
}
