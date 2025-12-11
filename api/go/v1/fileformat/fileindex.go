package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/samber/lo"
)

// IndexEntry represents a single file index entry.
//
// Size: 16 bytes (8 + 8)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
type IndexEntry struct {
	// FileID is the unique file identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	FileID uint64

	// Offset is the file entry offset from start of file
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	Offset uint64
}

// FileIndex represents the file index section of a package.
//
// Size: 16 bytes + (16 * entry_count) bytes
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
type FileIndex struct {
	// EntryCount is the number of file entries
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	EntryCount uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	Reserved uint32

	// FirstEntryOffset is the offset to the first file entry
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	FirstEntryOffset uint64

	// Entries contains all index entries
	// Specification: ../../docs/tech_specs/package_file_format.md Section 5
	Entries []IndexEntry
}

// Validate performs validation checks on the FileIndex.
//
// Validation checks:
//   - Reserved field must be zero
//   - EntryCount must match actual Entries length
//   - All FileIDs must be unique and non-zero
//
// Returns an error if any validation check fails.
func (f *FileIndex) Validate() error {
	if f.Reserved != 0 {
		return fmt.Errorf("reserved field must be zero")
	}

	if f.EntryCount != uint32(len(f.Entries)) {
		return fmt.Errorf("entry count mismatch: specified %d, actual %d", f.EntryCount, len(f.Entries))
	}

	// Check for zero FileIDs (needs index for error message)
	for i, entry := range f.Entries {
		if entry.FileID == 0 {
			return fmt.Errorf("file ID at index %d cannot be zero", i)
		}
	}

	// Check for duplicate FileIDs using lo.UniqBy
	unique := lo.UniqBy(f.Entries, func(e IndexEntry) uint64 { return e.FileID })
	if len(unique) != len(f.Entries) {
		// Find and report specific duplicate with indices
		seen := make(map[uint64]int)
		for i, entry := range f.Entries {
			if prev, exists := seen[entry.FileID]; exists {
				return fmt.Errorf("duplicate file ID %d at indices %d and %d", entry.FileID, prev, i)
			}
			seen[entry.FileID] = i
		}
	}

	return nil
}

// Size returns the total size of the FileIndex in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5
func (f *FileIndex) Size() int {
	return 16 + (IndexEntrySize * len(f.Entries))
}

// NewFileIndex creates and returns a new FileIndex with zero values.
//
// The returned FileIndex is initialized with all fields set to zero or empty.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
func NewFileIndex() *FileIndex {
	return &FileIndex{}
}

// ReadFrom reads a FileIndex from the provided io.Reader.
//
// The binary format is:
//   - EntryCount (4 bytes, little-endian uint32)
//   - Reserved (4 bytes, little-endian uint32, must be 0)
//   - FirstEntryOffset (8 bytes, little-endian uint64)
//   - Entries (EntryCount × 16 bytes, each entry: FileID(8) + Offset(8))
//
// Returns the number of bytes read and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
func (f *FileIndex) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read header (16 bytes)
	var entryCount uint32
	if err := binary.Read(r, binary.LittleEndian, &entryCount); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, fmt.Errorf("failed to read entry count: incomplete data: %w", err)
		}
		return totalRead, fmt.Errorf("failed to read entry count: %w", err)
	}
	totalRead += 4
	f.EntryCount = entryCount

	var reserved uint32
	if err := binary.Read(r, binary.LittleEndian, &reserved); err != nil {
		return totalRead, fmt.Errorf("failed to read reserved: %w", err)
	}
	totalRead += 4
	f.Reserved = reserved

	var firstEntryOffset uint64
	if err := binary.Read(r, binary.LittleEndian, &firstEntryOffset); err != nil {
		return totalRead, fmt.Errorf("failed to read first entry offset: %w", err)
	}
	totalRead += 8
	f.FirstEntryOffset = firstEntryOffset

	// Read entries
	f.Entries = make([]IndexEntry, 0, entryCount)
	for i := uint32(0); i < entryCount; i++ {
		var entry IndexEntry
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return totalRead, fmt.Errorf("failed to read entry %d: %w", i, err)
		}
		totalRead += IndexEntrySize
		f.Entries = append(f.Entries, entry)
	}

	return totalRead, nil
}

// WriteTo writes a FileIndex to the provided io.Writer.
//
// The binary format is:
//   - EntryCount (4 bytes, little-endian uint32)
//   - Reserved (4 bytes, little-endian uint32, must be 0)
//   - FirstEntryOffset (8 bytes, little-endian uint64)
//   - Entries (EntryCount × 16 bytes, each entry: FileID(8) + Offset(8))
//
// Before writing, the method updates EntryCount to match the actual Entries length.
//
// Returns the number of bytes written and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
func (f *FileIndex) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Update EntryCount to match actual entries
	f.EntryCount = uint32(len(f.Entries))

	// Write header (16 bytes)
	if err := binary.Write(w, binary.LittleEndian, f.EntryCount); err != nil {
		return totalWritten, fmt.Errorf("failed to write entry count: %w", err)
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, f.Reserved); err != nil {
		return totalWritten, fmt.Errorf("failed to write reserved: %w", err)
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, f.FirstEntryOffset); err != nil {
		return totalWritten, fmt.Errorf("failed to write first entry offset: %w", err)
	}
	totalWritten += 8

	// Write entries
	for i, entry := range f.Entries {
		if err := binary.Write(w, binary.LittleEndian, entry); err != nil {
			return totalWritten, fmt.Errorf("failed to write entry %d: %w", i, err)
		}
		totalWritten += IndexEntrySize
	}

	return totalWritten, nil
}
