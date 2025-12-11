package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/samber/lo"
)

// FileEntry represents a file entry with fixed and variable-length sections.
//
// The fixed section is 64 bytes, followed by variable-length data containing
// paths, hashes, and optional data.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1 - File Entry Binary Format Specification
type FileEntry struct {
	// Fixed section (64 bytes) - Section 4.1.1

	// FileID is the unique file identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.1
	FileID uint64

	// OriginalSize is the file size before compression/encryption
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	OriginalSize uint64

	// StoredSize is the file size after compression/encryption
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	StoredSize uint64

	// RawChecksum is the CRC32 of raw file content
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	RawChecksum uint32

	// StoredChecksum is the CRC32 of processed file content
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	StoredChecksum uint32

	// FileVersion tracks file data version
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.2
	FileVersion uint32

	// MetadataVersion tracks file metadata version
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.2
	MetadataVersion uint32

	// PathCount is the total number of paths
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	PathCount uint16

	// Type is the file type identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.3
	Type uint16

	// CompressionType is the compression algorithm identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.3
	CompressionType uint8

	// CompressionLevel is the compression level (0-9)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	CompressionLevel uint8

	// EncryptionType is the encryption algorithm identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1.3
	EncryptionType uint8

	// HashCount is the number of hash entries
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.2.1
	HashCount uint8

	// HashDataOffset is the offset to hash data from start of variable-length data
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.2.3
	HashDataOffset uint32

	// HashDataLen is the total length of all hash data
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.2.2
	HashDataLen uint16

	// OptionalDataLen is the total length of optional data
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.2.4
	OptionalDataLen uint16

	// OptionalDataOffset is the offset to optional data from start of variable-length data
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.2.5
	OptionalDataOffset uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
	Reserved uint32

	// Variable-length data (not part of fixed structure)

	// Paths contains all path entries for this file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	Paths []PathEntry

	// Hashes contains all hash entries for this file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3
	Hashes []HashEntry

	// OptionalData contains optional file attributes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4
	OptionalData []OptionalDataEntry
}

// Validate performs validation checks on the FileEntry.
//
// Validation checks:
//   - FileID must not be zero
//   - Reserved field must be zero
//   - PathCount must match actual Paths length
//   - HashCount must match actual Hashes length
//   - All paths must be valid
//   - All hashes must be valid
//   - All optional data must be valid
//
// Returns an error if any validation check fails.
func (f *FileEntry) Validate() error {
	if f.FileID == 0 {
		return fmt.Errorf("file ID cannot be zero")
	}

	if f.Reserved != 0 {
		return fmt.Errorf("reserved field must be zero")
	}

	if f.PathCount != uint16(len(f.Paths)) {
		return fmt.Errorf("path count mismatch: specified %d, actual %d", f.PathCount, len(f.Paths))
	}

	if f.HashCount != uint8(len(f.Hashes)) {
		return fmt.Errorf("hash count mismatch: specified %d, actual %d", f.HashCount, len(f.Hashes))
	}

	// Validate all paths
	for i, path := range f.Paths {
		if err := path.Validate(); err != nil {
			return fmt.Errorf("invalid path at index %d: %w", i, err)
		}
	}

	// Validate all hashes
	for i, hash := range f.Hashes {
		if err := hash.Validate(); err != nil {
			return fmt.Errorf("invalid hash at index %d: %w", i, err)
		}
	}

	// Validate all optional data
	for i, opt := range f.OptionalData {
		if err := opt.Validate(); err != nil {
			return fmt.Errorf("invalid optional data at index %d: %w", i, err)
		}
	}

	return nil
}

// FixedSize returns the size of the fixed section.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.1
func (f *FileEntry) FixedSize() int {
	return FileEntryFixedSize
}

// VariableSize returns the size of the variable-length data section.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4
func (f *FileEntry) VariableSize() int {
	pathSize := lo.SumBy(f.Paths, func(p PathEntry) int { return p.Size() })
	hashSize := lo.SumBy(f.Hashes, func(h HashEntry) int { return h.Size() })
	optSize := lo.SumBy(f.OptionalData, func(o OptionalDataEntry) int { return o.Size() })
	return pathSize + hashSize + optSize
}

// TotalSize returns the total size of the FileEntry (fixed + variable).
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1
func (f *FileEntry) TotalSize() int {
	return f.FixedSize() + f.VariableSize()
}

// NewFileEntry creates and returns a new FileEntry with zero values.
//
// The returned FileEntry is initialized with all fields set to zero or empty.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1 - File Entry Binary Format Specification
func NewFileEntry() *FileEntry {
	return &FileEntry{}
}

// ReadFrom reads a FileEntry from the provided io.Reader.
//
// The binary format is:
//   - Fixed section (64 bytes, little-endian)
//   - Variable-length data:
//   - Path entries (PathCount entries, starting at offset 0)
//   - Hash entries (HashCount entries, starting at HashDataOffset)
//   - Optional data entries (starting at OptionalDataOffset)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1 - File Entry Binary Format Specification
func (f *FileEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read fixed section (64 bytes)
	// Create a temporary struct to read the fixed fields
	type FileEntryFixed struct {
		FileID             uint64
		OriginalSize       uint64
		StoredSize         uint64
		RawChecksum        uint32
		StoredChecksum     uint32
		FileVersion        uint32
		MetadataVersion    uint32
		PathCount          uint16
		Type               uint16
		CompressionType    uint8
		CompressionLevel   uint8
		EncryptionType     uint8
		HashCount          uint8
		HashDataOffset     uint32
		HashDataLen        uint16
		OptionalDataLen    uint16
		OptionalDataOffset uint32
		Reserved           uint32
	}

	var fixed FileEntryFixed
	if err := binary.Read(r, binary.LittleEndian, &fixed); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, fmt.Errorf("failed to read fixed section: incomplete data (read %d bytes, expected %d): %w", totalRead, FileEntryFixedSize, err)
		}
		return totalRead, fmt.Errorf("failed to read fixed section: %w", err)
	}
	totalRead += FileEntryFixedSize

	// Copy fixed fields to FileEntry
	f.FileID = fixed.FileID
	f.OriginalSize = fixed.OriginalSize
	f.StoredSize = fixed.StoredSize
	f.RawChecksum = fixed.RawChecksum
	f.StoredChecksum = fixed.StoredChecksum
	f.FileVersion = fixed.FileVersion
	f.MetadataVersion = fixed.MetadataVersion
	f.PathCount = fixed.PathCount
	f.Type = fixed.Type
	f.CompressionType = fixed.CompressionType
	f.CompressionLevel = fixed.CompressionLevel
	f.EncryptionType = fixed.EncryptionType
	f.HashCount = fixed.HashCount
	f.HashDataOffset = fixed.HashDataOffset
	f.HashDataLen = fixed.HashDataLen
	f.OptionalDataLen = fixed.OptionalDataLen
	f.OptionalDataOffset = fixed.OptionalDataOffset
	f.Reserved = fixed.Reserved

	// Initialize slices
	f.Paths = make([]PathEntry, 0, f.PathCount)
	f.Hashes = make([]HashEntry, 0, f.HashCount)
	f.OptionalData = make([]OptionalDataEntry, 0)

	// Read path entries (starting at offset 0)
	for i := uint16(0); i < f.PathCount; i++ {
		var path PathEntry
		n, err := path.ReadFrom(r)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read path entry %d: %w", i, err)
		}
		totalRead += n
		f.Paths = append(f.Paths, path)
	}

	// Read hash entries (starting at HashDataOffset)
	// Calculate how many bytes we've read so far (paths)
	pathsSize := int64(lo.SumBy(f.Paths, func(p PathEntry) int { return p.Size() }))

	// If HashDataOffset is set, we may need to skip some bytes
	if f.HashDataOffset > 0 && int64(f.HashDataOffset) > pathsSize {
		skip := int64(f.HashDataOffset) - pathsSize
		_, err := io.CopyN(io.Discard, r, skip)
		if err != nil {
			return totalRead, fmt.Errorf("failed to skip to hash data offset: %w", err)
		}
		totalRead += skip
	}

	// Read hash entries
	for i := uint8(0); i < f.HashCount; i++ {
		var hash HashEntry
		n, err := hash.ReadFrom(r)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read hash entry %d: %w", i, err)
		}
		totalRead += n
		f.Hashes = append(f.Hashes, hash)
	}

	// Read optional data entries (starting at OptionalDataOffset)
	// Calculate current position after paths and hashes
	hashSize := int64(lo.SumBy(f.Hashes, func(h HashEntry) int { return h.Size() }))
	currentOffset := pathsSize + hashSize

	if f.OptionalDataOffset > 0 && int64(f.OptionalDataOffset) > currentOffset {
		skip := int64(f.OptionalDataOffset) - currentOffset
		_, err := io.CopyN(io.Discard, r, skip)
		if err != nil {
			return totalRead, fmt.Errorf("failed to skip to optional data offset: %w", err)
		}
		totalRead += skip
	}

	// Read optional data entries
	// We need to read until we've consumed OptionalDataLen bytes
	optionalDataRead := int64(0)
	for optionalDataRead < int64(f.OptionalDataLen) {
		var opt OptionalDataEntry
		n, err := opt.ReadFrom(r)
		if err != nil {
			if err == io.EOF && optionalDataRead > 0 {
				// We've read some optional data, this might be acceptable
				break
			}
			return totalRead, fmt.Errorf("failed to read optional data entry: %w", err)
		}
		totalRead += n
		optionalDataRead += n
		f.OptionalData = append(f.OptionalData, opt)
	}

	return totalRead, nil
}

// WriteTo writes a FileEntry to the provided io.Writer.
//
// The binary format is:
//   - Fixed section (64 bytes, little-endian)
//   - Variable-length data:
//   - Path entries (PathCount entries, starting at offset 0)
//   - Hash entries (HashCount entries, starting at HashDataOffset)
//   - Optional data entries (starting at OptionalDataOffset)
//
// Before writing, the method calculates and updates HashDataOffset, OptionalDataOffset,
// HashDataLen, and OptionalDataLen based on the actual variable-length data.
//
// Returns the number of bytes written and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1 - File Entry Binary Format Specification
func (f *FileEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Update counts to match actual data
	f.PathCount = uint16(len(f.Paths))
	f.HashCount = uint8(len(f.Hashes))

	// Calculate sizes and offsets
	pathsSize := 0
	for _, p := range f.Paths {
		pathsSize += p.Size()
	}

	hashSize := 0
	for _, h := range f.Hashes {
		hashSize += h.Size()
	}
	f.HashDataOffset = uint32(pathsSize)
	f.HashDataLen = uint16(hashSize)

	optionalDataSize := 0
	for _, o := range f.OptionalData {
		optionalDataSize += o.Size()
	}
	f.OptionalDataOffset = uint32(pathsSize + hashSize)
	f.OptionalDataLen = uint16(optionalDataSize)

	// Write fixed section (64 bytes)
	type FileEntryFixed struct {
		FileID             uint64
		OriginalSize       uint64
		StoredSize         uint64
		RawChecksum        uint32
		StoredChecksum     uint32
		FileVersion        uint32
		MetadataVersion    uint32
		PathCount          uint16
		Type               uint16
		CompressionType    uint8
		CompressionLevel   uint8
		EncryptionType     uint8
		HashCount          uint8
		HashDataOffset     uint32
		HashDataLen        uint16
		OptionalDataLen    uint16
		OptionalDataOffset uint32
		Reserved           uint32
	}

	fixed := FileEntryFixed{
		FileID:             f.FileID,
		OriginalSize:       f.OriginalSize,
		StoredSize:         f.StoredSize,
		RawChecksum:        f.RawChecksum,
		StoredChecksum:     f.StoredChecksum,
		FileVersion:        f.FileVersion,
		MetadataVersion:    f.MetadataVersion,
		PathCount:          f.PathCount,
		Type:               f.Type,
		CompressionType:    f.CompressionType,
		CompressionLevel:   f.CompressionLevel,
		EncryptionType:     f.EncryptionType,
		HashCount:          f.HashCount,
		HashDataOffset:     f.HashDataOffset,
		HashDataLen:        f.HashDataLen,
		OptionalDataLen:    f.OptionalDataLen,
		OptionalDataOffset: f.OptionalDataOffset,
		Reserved:           f.Reserved,
	}

	if err := binary.Write(w, binary.LittleEndian, fixed); err != nil {
		return totalWritten, fmt.Errorf("failed to write fixed section: %w", err)
	}
	totalWritten += FileEntryFixedSize

	// Write path entries (starting at offset 0)
	for i, path := range f.Paths {
		n, err := path.WriteTo(w)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write path entry %d: %w", i, err)
		}
		totalWritten += n
	}

	// Write hash entries (starting at HashDataOffset)
	for i, hash := range f.Hashes {
		n, err := hash.WriteTo(w)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write hash entry %d: %w", i, err)
		}
		totalWritten += n
	}

	// Write optional data entries (starting at OptionalDataOffset)
	for i, opt := range f.OptionalData {
		n, err := opt.WriteTo(w)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write optional data entry %d: %w", i, err)
		}
		totalWritten += n
	}

	return totalWritten, nil
}
