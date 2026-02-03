// This file implements the PackageHeader structure representing the fixed-size
// header of a NovusPack (.nvpk) file. It contains the PackageHeader type definition,
// marshaling/unmarshaling methods, validation, and header initialization functions.
// This file should contain all code related to reading, writing, and validating
// the 112-byte package header as specified in package_file_format.md Section 2.
//
// Specification: package_file_format.md: 2 Package Header

// Package fileformat provides file format domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to the NovusPack (.nvpk) file format
// as specified in package_file_format.md.
package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// PackageHeader represents the fixed-size header of a NovusPack (.nvpk) file.
//
// The header provides comprehensive metadata and navigation information for the
// entire package. All fields are stored in little-endian byte order.
//
// Size: 112 bytes (fixed)
//
// Specification: package_file_format.md: 2.1 Header Structure
type PackageHeader struct {
	// Magic is the package identifier (0x4E56504B "NVPK")
	// Specification: package_file_format.md: 2.1 Header Structure
	Magic uint32

	// FormatVersion is the format version (current: 1)
	// Specification: package_file_format.md: 2.1 Header Structure
	FormatVersion uint32

	// Flags contains package-level features and options
	// Bits 0-7: Package features
	// Bits 8-15: Package compression type
	// Bits 16-31: Reserved for future use
	// Specification: package_file_format.md: 2.5 Package Features Flags
	Flags uint32

	// PackageDataVersion tracks changes to package data content
	// Increments on file additions, removals, or data modifications
	// Specification: package_file_format.md: 2.2.1 PackageDataVersion Field
	PackageDataVersion uint32

	// MetadataVersion tracks changes to package metadata
	// Increments on metadata or comment modifications
	// Specification: package_file_format.md: 2.2.2 MetadataVersion Field
	MetadataVersion uint32

	// PackageCRC is the CRC32 of package content (0 if skipped)
	// Excludes header and signatures
	// Specification: package_file_format.md: 2.2.3 PackageCRC Field
	PackageCRC uint32

	// CreatedTime is the package creation timestamp (Unix nanoseconds)
	// Specification: package_file_format.md: 2.1 Header Structure
	CreatedTime uint64

	// ModifiedTime is the package modification timestamp (Unix nanoseconds)
	// Specification: package_file_format.md: 2.1 Header Structure
	ModifiedTime uint64

	// LocaleID is the locale identifier for path encoding
	// Specification: package_file_format.md: 2.7 LocaleID Field Specification
	LocaleID uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: package_file_format.md: 2.1 Header Structure
	Reserved uint32

	// AppID is the application/game identifier (0 if not associated)
	// Specification: package_file_format.md: 2.4 AppID Field Specification
	AppID uint64

	// VendorID is the storefront/platform identifier (0 if not associated)
	// Specification: package_file_format.md: 2.3 VendorID Field Specification
	VendorID uint32

	// CreatorID is the creator identifier (reserved for future use)
	// Specification: package_file_format.md: 2.1 Header Structure
	CreatorID uint32

	// IndexStart is the offset to file index from start of file
	// Specification: package_file_format.md: 2.1 Header Structure
	IndexStart uint64

	// IndexSize is the size of file index in bytes
	// Specification: package_file_format.md: 2.1 Header Structure
	IndexSize uint64

	// ArchiveChainID is the archive chain identifier
	// Specification: package_file_format.md: 2.1 Header Structure
	ArchiveChainID uint64

	// ArchivePartInfo contains combined part number and total parts
	// Bits 31-16: Part number (0-65535)
	// Bits 15-0: Total parts (1-65535)
	// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
	ArchivePartInfo uint32

	// CommentSize is the size of package comment in bytes (0 if no comment)
	// Specification: package_file_format.md: 2.1 Header Structure
	CommentSize uint32

	// CommentStart is the offset to package comment from start of file
	// Specification: package_file_format.md: 2.1 Header Structure
	CommentStart uint64

	// SignatureOffset is the offset to signatures block from start of file
	// Specification: package_file_format.md: 2.1 Header Structure
	SignatureOffset uint64
}

// Validate performs validation checks on the PackageHeader.
//
// Validation checks:
//   - Magic number must be NVPKMagic (0x4E56504B)
//   - FormatVersion must be 1 (current version)
//   - Reserved field must be 0
//   - Version fields must be >= 1 for initialized packages
//
// Returns an error if any validation check fails.
//
// Specification: package_file_format.md: 2.1 Header Structure
func (h *PackageHeader) validate() error {
	// Validate magic number
	if h.Magic != NVPKMagic {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid magic number", nil, pkgerrors.ValidationErrorContext{
			Field:    "Magic",
			Value:    fmt.Sprintf("0x%08X", h.Magic),
			Expected: fmt.Sprintf("0x%08X", NVPKMagic),
		})
	}

	// Validate format version
	if h.FormatVersion != FormatVersion {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "unsupported format version", nil, pkgerrors.ValidationErrorContext{
			Field:    "FormatVersion",
			Value:    h.FormatVersion,
			Expected: fmt.Sprintf("%d", FormatVersion),
		})
	}

	// Validate reserved field must be zero
	if h.Reserved != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "reserved field must be 0", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    h.Reserved,
			Expected: "0",
		})
	}

	return nil
}

// GetCompressionType extracts the compression type from the Flags field.
//
// Returns the compression type as defined in the constants:
//   - 0: No compression
//   - 1: Zstd compression
//   - 2: LZ4 compression
//   - 3: LZMA compression
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) getCompressionType() uint8 {
	return uint8((h.Flags & FlagsMaskCompressionType) >> FlagsShiftCompressionType)
}

// SetCompressionType sets the compression type in the Flags field.
//
// The compression type is stored in bits 8-15 of the Flags field.
// Preserves existing feature flags (bits 0-7).
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) setCompressionType(compressionType uint8) {
	// Clear compression type bits
	h.Flags &= ^uint32(FlagsMaskCompressionType)
	// Set new compression type
	h.Flags |= uint32(compressionType) << FlagsShiftCompressionType
}

// GetFeatures extracts the feature flags from the Flags field.
//
// Returns the feature flags (bits 0-7) as a bitmask.
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) getFeatures() uint8 {
	return uint8(h.Flags & FlagsMaskFeatures)
}

// HasFeature checks if a specific feature flag is set.
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) hasFeature(flag uint32) bool {
	return (h.Flags & flag) != 0
}

// SetFeature sets a specific feature flag.
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) setFeature(flag uint32) {
	h.Flags |= flag
}

// ClearFeature clears a specific feature flag.
//
// Specification: package_file_format.md: 2.5 Package Features Flags
func (h *PackageHeader) clearFeature(flag uint32) {
	h.Flags &= ^flag
}

// GetArchivePart extracts the part number from ArchivePartInfo.
//
// Returns the part number (bits 31-16).
//
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func (h *PackageHeader) getArchivePart() uint16 {
	return uint16(h.ArchivePartInfo >> 16)
}

// GetArchiveTotal extracts the total parts from ArchivePartInfo.
//
// Returns the total parts (bits 15-0).
//
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func (h *PackageHeader) getArchiveTotal() uint16 {
	return uint16(h.ArchivePartInfo & 0xFFFF)
}

// SetArchivePartInfo sets both part number and total parts in ArchivePartInfo.
//
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func (h *PackageHeader) setArchivePartInfo(part, total uint16) {
	h.ArchivePartInfo = (uint32(part) << 16) | uint32(total)
}

// IsSigned returns true if the package has signatures.
//
// Specification: package_file_format.md: 2.9 Signed Package File Immutability and Incremental Signatures
func (h *PackageHeader) isSigned() bool {
	return h.SignatureOffset > 0
}

// HasComment returns true if the package has a comment.
//
// Specification: package_file_format.md: 7.1 Package Comment Format Specification
func (h *PackageHeader) hasComment() bool {
	return h.CommentSize > 0
}

// NewPackageHeader creates and returns a new PackageHeader with default values.
//
// The returned PackageHeader is initialized according to Section 2.8.1:
//   - Magic set to NVPKMagic (0x4E56504B)
//   - FormatVersion set to 1
//   - PackageDataVersion set to 1
//   - MetadataVersion set to 1
//   - Reserved set to 0
//   - ArchivePartInfo set to 0x00010001 (part 1 of 1)
//   - All other fields set to 0
//
// Specification: package_file_format.md: 2.8.1 Initial Package Creation
func NewPackageHeader() *PackageHeader {
	return &PackageHeader{
		Magic:              NVPKMagic,
		FormatVersion:      FormatVersion,
		PackageDataVersion: 1,
		MetadataVersion:    1,
		Reserved:           0,
		ArchivePartInfo:    0x00010001, // Part 1 of 1
	}
}

// ReadFrom reads a PackageHeader from the provided io.Reader.
//
// The binary format is 112 bytes in little-endian byte order:
//   - Magic (4 bytes)
//   - FormatVersion (4 bytes)
//   - Flags (4 bytes)
//   - PackageDataVersion (4 bytes)
//   - MetadataVersion (4 bytes)
//   - PackageCRC (4 bytes)
//   - CreatedTime (8 bytes)
//   - ModifiedTime (8 bytes)
//   - LocaleID (4 bytes)
//   - Reserved (4 bytes)
//   - AppID (8 bytes)
//   - VendorID (4 bytes)
//   - CreatorID (4 bytes)
//   - IndexStart (8 bytes)
//   - IndexSize (8 bytes)
//   - ArchiveChainID (8 bytes)
//   - ArchivePartInfo (4 bytes)
//   - CommentSize (4 bytes)
//   - CommentStart (8 bytes)
//   - SignatureOffset (8 bytes)
//
// Returns the number of bytes read and any error encountered.
// If the magic number is invalid, returns a validation error.
//
// Specification: package_file_format.md: 2.1 Header Structure
func (h *PackageHeader) readFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read all fields in order using binary.Read for proper little-endian handling
	err := binary.Read(r, binary.LittleEndian, h)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, fmt.Sprintf("failed to read header: incomplete data (read %d bytes, expected %d)", totalRead, PackageHeaderSize), pkgerrors.ValidationErrorContext{
				Field:    "Header",
				Value:    totalRead,
				Expected: fmt.Sprintf("%d bytes", PackageHeaderSize),
			})
		}
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read header", pkgerrors.ValidationErrorContext{
			Field:    "Header",
			Value:    nil,
			Expected: fmt.Sprintf("%d bytes", PackageHeaderSize),
		})
	}
	totalRead = PackageHeaderSize

	// Validate magic number immediately after reading
	if h.Magic != NVPKMagic {
		return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid magic number", nil, pkgerrors.ValidationErrorContext{
			Field:    "Magic",
			Value:    fmt.Sprintf("0x%08X", h.Magic),
			Expected: fmt.Sprintf("0x%08X", NVPKMagic),
		})
	}

	return totalRead, nil
}

// WriteTo writes a PackageHeader to the provided io.Writer.
//
// The binary format is 112 bytes in little-endian byte order.
// All fields are written in the order specified in Section 2.1.
//
// Returns the number of bytes written and any error encountered.
//
// Specification: package_file_format.md: 2.1 Header Structure
func (h *PackageHeader) writeTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write all fields in order using binary.Write for proper little-endian handling
	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write header", pkgerrors.ValidationErrorContext{
			Field:    "Header",
			Value:    nil,
			Expected: fmt.Sprintf("%d bytes", PackageHeaderSize),
		})
	}
	totalWritten = PackageHeaderSize

	return totalWritten, nil
}
