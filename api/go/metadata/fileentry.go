// This file implements the FileEntry structure representing a file entry in a
// NovusPack package. It contains the FileEntry type definition with fixed 64-byte
// header and variable-length data sections, along with basic file entry operations.
// This file should contain the core FileEntry struct, NewFileEntry constructor,
// and basic property access methods as specified in api_file_mgmt_file_entry.md Section 1.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure

package metadata

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
	"github.com/samber/lo"
)

// FileEntryFixedSize is the fixed size of the FileEntry structure in bytes.
const FileEntryFixedSize = 64

// FileEntry represents a file entry with fixed and variable-length sections.
//
// The fixed section is 64 bytes, followed by variable-length data containing
// paths, hashes, and optional data.
//
// Specification: package_file_format.md: 4.1 FileEntry Binary Format Specification
type FileEntry struct {
	// Fixed section (64 bytes) - Section 4.1.1

	// FileID is the unique file identifier
	// Specification: package_file_format.md: 4.1.1.1 FileID Field Specification
	FileID uint64

	// OriginalSize is the file size before compression/encryption
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	OriginalSize uint64

	// StoredSize is the file size after compression/encryption
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	StoredSize uint64

	// RawChecksum is the CRC32 of raw file content
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	RawChecksum uint32

	// StoredChecksum is the CRC32 of processed file content
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	StoredChecksum uint32

	// FileVersion tracks file data version
	// Specification: package_file_format.md: 4.1.1.2 File Version Fields Specification
	FileVersion uint32

	// MetadataVersion tracks file metadata version
	// Specification: package_file_format.md: 4.1.1.2 File Version Fields Specification
	MetadataVersion uint32

	// PathCount is the total number of paths
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	PathCount uint16

	// Type is the file type identifier
	// Specification: package_file_format.md: 4.1.1.3 Compression and Encryption Types
	Type uint16

	// CompressionType is the compression algorithm identifier
	// Specification: package_file_format.md: 4.1.1.3 Compression and Encryption Types
	CompressionType uint8

	// CompressionLevel is the compression level (0-9)
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	CompressionLevel uint8

	// EncryptionType is the encryption algorithm identifier
	// Specification: package_file_format.md: 4.1.1.3 Compression and Encryption Types
	EncryptionType uint8

	// HashCount is the number of hash entries
	// Specification: package_file_format.md: 4.2.1 HashCount Field
	HashCount uint8

	// HashDataOffset is the offset to hash data from start of variable-length data
	// Specification: package_file_format.md: 4.2.3 HashDataOffset Field
	HashDataOffset uint32

	// HashDataLen is the total length of all hash data
	// Specification: package_file_format.md: 4.2.2 HashDataLen Field
	HashDataLen uint16

	// OptionalDataLen is the total length of optional data
	// Specification: package_file_format.md: 4.2.4 OptionalDataLen Field
	OptionalDataLen uint16

	// OptionalDataOffset is the offset to optional data from start of variable-length data
	// Specification: package_file_format.md: 4.2.5 OptionalDataOffset Field
	OptionalDataOffset uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
	Reserved uint32

	// Variable-length data (not part of fixed structure)

	// Paths contains all path entries for this file
	// Specification: package_file_format.md: 4.1.4.2 Path Entries
	Paths []generics.PathEntry

	// Hashes contains all hash entries for this file
	// Specification: package_file_format.md: 4.1.4.3 Hash Data
	Hashes []HashEntry

	// OptionalData contains optional file attributes
	// Specification: package_file_format.md: 4.1.4.4 Optional Data
	OptionalData []OptionalDataEntry

	// Data management (runtime only, not stored in file)
	// Specification: api_file_mgmt_file_entry.md: 1.1 FileEntry Structure Definition
	EntryOffset     uint64          // Absolute offset to the FileEntry metadata start in the package file
	Data            []byte          // File content in memory (only for small files being processed)
	SourceFile      *os.File        // Source file handle for streaming
	SourceOffset    int64           // Offset in source file
	SourceSize      int64           // Size of data to read from source
	TempFilePath    string          // Path to temp file for large files being processed
	IsDataLoaded    bool            // Whether data is currently loaded in memory
	IsTempFile      bool            // Whether file is stored as temp file during processing
	ProcessingState ProcessingState // Current processing state of the file

	// PathMetadataEntry associations (runtime only, not stored in file)
	// Maps each path in FileEntry.Paths to its corresponding PathMetadataEntry
	// Inheritance is now handled only on PathMetadataEntry, not FileEntry.
	// To access inherited tags, use the associated PathMetadataEntry's GetInheritedTags() method.
	// Specification: api_file_mgmt_file_entry.md: 1.1 FileEntry Structure Definition
	PathMetadataEntries map[string]*PathMetadataEntry // Path -> PathMetadataEntry mapping
}

// FileEntryFixed is the 64-byte fixed section used for binary read/write.
// Specification: package_file_format.md: 4.1 FileEntry Binary Format Specification
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

// skipReaderToOffset skips bytes in r from currentOffset to targetOffset when targetOffset > currentOffset.
// Returns bytes skipped and nil, or 0 and a wrapped error on failure.
func skipReaderToOffset(r io.Reader, _, currentOffset, targetOffset int64, fieldName string, fieldValue uint32) (int64, error) {
	if targetOffset <= 0 || targetOffset <= currentOffset {
		return 0, nil
	}
	skip := targetOffset - currentOffset
	_, err := io.CopyN(io.Discard, r, skip)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to skip to "+fieldName, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    fieldValue,
			Expected: "skip successful",
		})
	}
	return skip, nil
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
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file ID cannot be zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "FileID",
			Value:    f.FileID,
			Expected: "non-zero value",
		})
	}

	if f.Reserved != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "reserved field must be zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    f.Reserved,
			Expected: "0",
		})
	}

	if f.PathCount != uint16(len(f.Paths)) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path count mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "PathCount",
			Value:    f.PathCount,
			Expected: fmt.Sprintf("%d", len(f.Paths)),
		})
	}

	if f.HashCount != uint8(len(f.Hashes)) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "hash count mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "HashCount",
			Value:    f.HashCount,
			Expected: fmt.Sprintf("%d", len(f.Hashes)),
		})
	}

	// Validate all paths
	for i, path := range f.Paths {
		if err := path.Validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, fmt.Sprintf("invalid path at index %d", i), pkgerrors.ValidationErrorContext{
				Field:    "Paths",
				Value:    i,
				Expected: "valid path entry",
			})
		}
	}

	// Validate all hashes
	for i, hash := range f.Hashes {
		if err := hash.validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, fmt.Sprintf("invalid hash at index %d", i), pkgerrors.ValidationErrorContext{
				Field:    "Hashes",
				Value:    i,
				Expected: "valid hash entry",
			})
		}
	}

	// Validate all optional data
	for i, opt := range f.OptionalData {
		if err := opt.validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, fmt.Sprintf("invalid optional data at index %d", i), pkgerrors.ValidationErrorContext{
				Field:    "OptionalData",
				Value:    i,
				Expected: "valid optional data entry",
			})
		}
	}

	return nil
}

// FixedSize returns the size of the fixed section.
//
// Specification: package_file_format.md: 4.1.1 FileEntry Static Section Field Encoding
func (f *FileEntry) FixedSize() int {
	return FileEntryFixedSize
}

// VariableSize returns the size of the variable-length data section.
//
// Specification: package_file_format.md: 4.1.4 Variable-Length Data (follows fixed structure)
func (f *FileEntry) VariableSize() int {
	pathSize := lo.SumBy(f.Paths, func(p generics.PathEntry) int { return p.Size() })
	hashSize := lo.SumBy(f.Hashes, func(h HashEntry) int { return h.size() })
	optSize := lo.SumBy(f.OptionalData, func(o OptionalDataEntry) int { return o.size() })
	return pathSize + hashSize + optSize
}

// TotalSize returns the total size of the FileEntry (fixed + variable).
//
// Specification: package_file_format.md: 4.1 FileEntry Binary Format Specification
func (f *FileEntry) TotalSize() int {
	return f.FixedSize() + f.VariableSize()
}

// NewFileEntry creates and returns a new FileEntry with zero values.
//
// The returned FileEntry is initialized with all fields set to zero or empty.
// ProcessingState is initialized to ProcessingStateIdle.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func NewFileEntry() *FileEntry {
	return &FileEntry{
		ProcessingState: ProcessingStateIdle,
	}
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
// readFileEntryFixed reads the 64-byte fixed section, copies to f, and initializes slices; returns bytes read and error.
func readFileEntryFixed(r io.Reader, f *FileEntry) (int64, error) {
	var fixed FileEntryFixed
	if err := binary.Read(r, binary.LittleEndian, &fixed); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, fmt.Sprintf("failed to read fixed section: incomplete data (expected %d bytes)", FileEntryFixedSize), pkgerrors.ValidationErrorContext{
				Field: "FixedSection", Value: int64(0), Expected: fmt.Sprintf("%d bytes", FileEntryFixedSize),
			})
		}
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read fixed section", pkgerrors.ValidationErrorContext{
			Field: "FixedSection", Value: nil, Expected: fmt.Sprintf("%d bytes", FileEntryFixedSize),
		})
	}
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
	f.Paths = make([]generics.PathEntry, 0, f.PathCount)
	f.Hashes = make([]HashEntry, 0, f.HashCount)
	f.OptionalData = make([]OptionalDataEntry, 0)
	return FileEntryFixedSize, nil
}

// readFileEntryPaths reads PathCount path entries into f.Paths; returns total bytes read (including existing totalRead) and error.
func readFileEntryPaths(r io.Reader, f *FileEntry, totalRead int64) (int64, error) {
	for i := uint16(0); i < f.PathCount; i++ {
		var path generics.PathEntry
		n, err := path.ReadFrom(r)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to read path entry %d", i), pkgerrors.ValidationErrorContext{
				Field: "Paths", Value: i, Expected: "valid path entry",
			})
		}
		totalRead += n
		f.Paths = append(f.Paths, path)
	}
	return totalRead, nil
}

// readFileEntryHashes skips to HashDataOffset if needed, then reads HashCount hash entries; returns total bytes read and error.
func readFileEntryHashes(r io.Reader, f *FileEntry, totalRead, pathsSize int64) (int64, error) {
	n, err := skipReaderToOffset(r, totalRead, pathsSize, int64(f.HashDataOffset), "HashDataOffset", f.HashDataOffset)
	if err != nil {
		return totalRead, err
	}
	totalRead += n
	for i := uint8(0); i < f.HashCount; i++ {
		var hash HashEntry
		n, err := hash.readFrom(r)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to read hash entry %d", i), pkgerrors.ValidationErrorContext{
				Field: "Hashes", Value: i, Expected: "valid hash entry",
			})
		}
		totalRead += n
		f.Hashes = append(f.Hashes, hash)
	}
	return totalRead, nil
}

// readFileEntryOptionalData skips to OptionalDataOffset if needed, then reads optional data entries until OptionalDataLen bytes consumed.
func readFileEntryOptionalData(r io.Reader, f *FileEntry, totalRead, currentOffset int64) (int64, error) {
	n, err := skipReaderToOffset(r, totalRead, currentOffset, int64(f.OptionalDataOffset), "OptionalDataOffset", f.OptionalDataOffset)
	if err != nil {
		return totalRead, err
	}
	totalRead += n
	optionalDataRead := int64(0)
	for optionalDataRead < int64(f.OptionalDataLen) {
		var opt OptionalDataEntry
		n, err := opt.readFrom(r)
		if err != nil {
			if err == io.EOF && optionalDataRead > 0 {
				break
			}
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read optional data entry", pkgerrors.ValidationErrorContext{
				Field: "OptionalData", Value: optionalDataRead, Expected: fmt.Sprintf("%d bytes", f.OptionalDataLen),
			})
		}
		totalRead += n
		optionalDataRead += n
		f.OptionalData = append(f.OptionalData, opt)
	}
	return totalRead, nil
}

// Returns the number of bytes read and any error encountered.
//
// Specification: package_file_format.md: 4.1 FileEntry Binary Format Specification
func (f *FileEntry) ReadFrom(r io.Reader) (int64, error) {
	totalRead, err := readFileEntryFixed(r, f)
	if err != nil {
		return totalRead, err
	}
	totalRead, err = readFileEntryPaths(r, f, totalRead)
	if err != nil {
		return totalRead, err
	}
	pathsSize := int64(lo.SumBy(f.Paths, func(p generics.PathEntry) int { return p.Size() }))
	totalRead, err = readFileEntryHashes(r, f, totalRead, pathsSize)
	if err != nil {
		return totalRead, err
	}
	hashSize := int64(lo.SumBy(f.Hashes, func(h HashEntry) int { return h.size() }))
	return readFileEntryOptionalData(r, f, totalRead, pathsSize+hashSize)
}

// WriteTo writes both metadata and data to a writer.
//
// Writes both metadata and data to a writer.
// Implements the io.WriterTo interface.
// Provides memory-efficient marshaling for large files.
//
// Parameters:
//   - w: Writer to write both metadata and data to
//
// Returns total number of bytes written and error.
//
// Error conditions:
//   - ErrTypeValidation: Invalid FileEntry state or data not available
//   - ErrTypeIO: I/O error during writing
//
// Note: WriteTo requires data to be available. If data is not available,
// WriteDataTo will return an error which is passed through by WriteTo.
// Use WriteMetaTo if you only need to write metadata without data.
//
// Use this method when you need to stream both metadata and data together.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) WriteTo(w io.Writer) (int64, error) {
	// Write metadata first
	metaWritten, err := f.WriteMetaTo(w)
	if err != nil {
		return metaWritten, err
	}

	// Write data second
	dataWritten, err := f.WriteDataTo(w)
	if err != nil {
		// Pass through all errors from WriteDataTo
		return metaWritten + dataWritten, err
	}

	return metaWritten + dataWritten, nil
}
