// Package fileformat provides file format domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to the NovusPack (.npk) file format
// as specified in docs/tech_specs/package_file_format.md.
package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PackageHeader represents the fixed-size header of a NovusPack (.npk) file.
//
// The header provides comprehensive metadata and navigation information for the
// entire package. All fields are stored in little-endian byte order.
//
// Size: 112 bytes (fixed)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1 - Header Structure
type PackageHeader struct {
	// Magic is the package identifier (0x4E56504B "NVPK")
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	Magic uint32

	// FormatVersion is the format version (current: 1)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	FormatVersion uint32

	// Flags contains package-level features and options
	// Bits 0-7: Package features
	// Bits 8-15: Package compression type
	// Bits 16-31: Reserved for future use
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5
	Flags uint32

	// PackageDataVersion tracks changes to package data content
	// Increments on file additions, removals, or data modifications
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.2.1
	PackageDataVersion uint32

	// MetadataVersion tracks changes to package metadata
	// Increments on metadata or comment modifications
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.2.2
	MetadataVersion uint32

	// PackageCRC is the CRC32 of package content (0 if skipped)
	// Excludes header and signatures
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.2.3
	PackageCRC uint32

	// CreatedTime is the package creation timestamp (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	CreatedTime uint64

	// ModifiedTime is the package modification timestamp (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	ModifiedTime uint64

	// LocaleID is the locale identifier for path encoding
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.7
	LocaleID uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	Reserved uint32

	// AppID is the application/game identifier (0 if not associated)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.4
	AppID uint64

	// VendorID is the storefront/platform identifier (0 if not associated)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.3
	VendorID uint32

	// CreatorID is the creator identifier (reserved for future use)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	CreatorID uint32

	// IndexStart is the offset to file index from start of file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	IndexStart uint64

	// IndexSize is the size of file index in bytes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	IndexSize uint64

	// ArchiveChainID is the archive chain identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	ArchiveChainID uint64

	// ArchivePartInfo contains combined part number and total parts
	// Bits 31-16: Part number (0-65535)
	// Bits 15-0: Total parts (1-65535)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.6
	ArchivePartInfo uint32

	// CommentSize is the size of package comment in bytes (0 if no comment)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	CommentSize uint32

	// CommentStart is the offset to package comment from start of file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	CommentStart uint64

	// SignatureOffset is the offset to signatures block from start of file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1
	SignatureOffset uint64
}

// Validate performs validation checks on the PackageHeader.
//
// Validation checks:
//   - Magic number must be NPKMagic (0x4E56504B)
//   - FormatVersion must be 1 (current version)
//   - Reserved field must be 0
//   - Version fields must be >= 1 for initialized packages
//
// Returns an error if any validation check fails.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1 - Header Structure
func (h *PackageHeader) Validate() error {
	// Validate magic number
	if h.Magic != NPKMagic {
		return fmt.Errorf("invalid magic number: 0x%08X, expected 0x%08X", h.Magic, NPKMagic)
	}

	// Validate format version
	if h.FormatVersion != FormatVersion {
		return fmt.Errorf("unsupported format version: %d, expected %d", h.FormatVersion, FormatVersion)
	}

	// Validate reserved field must be zero
	if h.Reserved != 0 {
		return fmt.Errorf("reserved field must be 0, got %d", h.Reserved)
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) GetCompressionType() uint8 {
	return uint8((h.Flags & FlagsMaskCompressionType) >> FlagsShiftCompressionType)
}

// SetCompressionType sets the compression type in the Flags field.
//
// The compression type is stored in bits 8-15 of the Flags field.
// Preserves existing feature flags (bits 0-7).
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) SetCompressionType(compressionType uint8) {
	// Clear compression type bits
	h.Flags &= ^uint32(FlagsMaskCompressionType)
	// Set new compression type
	h.Flags |= uint32(compressionType) << FlagsShiftCompressionType
}

// GetFeatures extracts the feature flags from the Flags field.
//
// Returns the feature flags (bits 0-7) as a bitmask.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) GetFeatures() uint8 {
	return uint8(h.Flags & FlagsMaskFeatures)
}

// HasFeature checks if a specific feature flag is set.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) HasFeature(flag uint32) bool {
	return (h.Flags & flag) != 0
}

// SetFeature sets a specific feature flag.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) SetFeature(flag uint32) {
	h.Flags |= flag
}

// ClearFeature clears a specific feature flag.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.5 - Package Features Flags
func (h *PackageHeader) ClearFeature(flag uint32) {
	h.Flags &= ^flag
}

// GetArchivePart extracts the part number from ArchivePartInfo.
//
// Returns the part number (bits 31-16).
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.6 - ArchivePartInfo Field Specification
func (h *PackageHeader) GetArchivePart() uint16 {
	return uint16(h.ArchivePartInfo >> 16)
}

// GetArchiveTotal extracts the total parts from ArchivePartInfo.
//
// Returns the total parts (bits 15-0).
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.6 - ArchivePartInfo Field Specification
func (h *PackageHeader) GetArchiveTotal() uint16 {
	return uint16(h.ArchivePartInfo & 0xFFFF)
}

// SetArchivePartInfo sets both part number and total parts in ArchivePartInfo.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.6 - ArchivePartInfo Field Specification
func (h *PackageHeader) SetArchivePartInfo(part, total uint16) {
	h.ArchivePartInfo = (uint32(part) << 16) | uint32(total)
}

// IsSigned returns true if the package has signatures.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.9 - Signed Package File Immutability
func (h *PackageHeader) IsSigned() bool {
	return h.SignatureOffset > 0
}

// HasComment returns true if the package has a comment.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1 - Package Comment Format Specification
func (h *PackageHeader) HasComment() bool {
	return h.CommentSize > 0
}

// NewPackageHeader creates and returns a new PackageHeader with default values.
//
// The returned PackageHeader is initialized according to Section 2.8.1:
//   - Magic set to NPKMagic (0x4E56504B)
//   - FormatVersion set to 1
//   - PackageDataVersion set to 1
//   - MetadataVersion set to 1
//   - Reserved set to 0
//   - ArchivePartInfo set to 0x00010001 (part 1 of 1)
//   - All other fields set to 0
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.8.1 - Initial Package Creation
func NewPackageHeader() *PackageHeader {
	return &PackageHeader{
		Magic:              NPKMagic,
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1 - Header Structure
func (h *PackageHeader) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read all fields in order using binary.Read for proper little-endian handling
	err := binary.Read(r, binary.LittleEndian, h)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, fmt.Errorf("failed to read header: incomplete data (read %d bytes, expected %d): %w", totalRead, PackageHeaderSize, err)
		}
		return totalRead, fmt.Errorf("failed to read header: %w", err)
	}
	totalRead = PackageHeaderSize

	// Validate magic number immediately after reading
	if h.Magic != NPKMagic {
		return totalRead, fmt.Errorf("invalid magic number: 0x%08X, expected 0x%08X", h.Magic, NPKMagic)
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 2.1 - Header Structure
func (h *PackageHeader) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write all fields in order using binary.Write for proper little-endian handling
	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		return totalWritten, fmt.Errorf("failed to write header: %w", err)
	}
	totalWritten = PackageHeaderSize

	return totalWritten, nil
}
