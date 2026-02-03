// This file implements the FileIndex structure representing the file index section
// of a NovusPack package. It contains the FileIndex and IndexEntry type definitions,
// marshaling/unmarshaling methods, and index management functions. This file should
// contain all code related to reading, writing, and managing the file index as
// specified in package_file_format.md Section 6.
//
// Specification: package_file_format.md: 6 File Index Section

package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"
	"runtime"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
	"github.com/samber/lo"
)

// IndexEntry represents a single file index entry.
//
// Size: 16 bytes (8 + 8)
//
// Specification: package_file_format.md: 6 File Index Section
type IndexEntry struct {
	// FileID is the unique file identifier
	// Specification: package_file_format.md: 6 File Index Section
	FileID uint64

	// Offset is the file entry offset from start of file
	// Specification: package_file_format.md: 6 File Index Section
	Offset uint64
}

// FileIndex represents the file index section of a package.
//
// Size: 16 bytes + (16 * entry_count) bytes
//
// Specification: package_file_format.md: 6 File Index Section
type FileIndex struct {
	// EntryCount is the number of file entries
	// Specification: package_file_format.md: 6 File Index Section
	EntryCount uint32

	// Reserved is reserved for future use (must be 0)
	// Specification: package_file_format.md: 6 File Index Section
	Reserved uint32

	// FirstEntryOffset is the offset to the first file entry
	// Specification: package_file_format.md: 6 File Index Section
	FirstEntryOffset uint64

	// Entries contains all index entries
	// Specification: package_file_format.md: 6 File Index Section
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
func (f *FileIndex) validate() error {
	if f.Reserved != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "reserved field must be zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    f.Reserved,
			Expected: "0",
		})
	}

	if f.EntryCount != uint32(len(f.Entries)) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "entry count mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "EntryCount",
			Value:    f.EntryCount,
			Expected: fmt.Sprintf("%d", len(f.Entries)),
		})
	}

	// Check for zero FileIDs (needs index for error message)
	for i, entry := range f.Entries {
		if entry.FileID == 0 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("file ID at index %d cannot be zero", i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Entries",
				Value:    i,
				Expected: "non-zero FileID",
			})
		}
	}

	// Check for duplicate FileIDs using lo.UniqBy
	unique := lo.UniqBy(f.Entries, func(e IndexEntry) uint64 { return e.FileID })
	if len(unique) != len(f.Entries) {
		// Find and report specific duplicate with indices
		seen := make(map[uint64]int)
		for i, entry := range f.Entries {
			if prev, exists := seen[entry.FileID]; exists {
				return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("duplicate file ID %d at indices %d and %d", entry.FileID, prev, i), nil, pkgerrors.ValidationErrorContext{
					Field:    "Entries",
					Value:    entry.FileID,
					Expected: "unique FileID",
				})
			}
			seen[entry.FileID] = i
		}
	}

	return nil
}

// Size returns the total size of the FileIndex in bytes.
//
// Specification: package_file_format.md: 6 File Index Section
func (f *FileIndex) size() int {
	return 16 + (IndexEntrySize * len(f.Entries))
}

// NewFileIndex creates and returns a new FileIndex with zero values.
//
// The returned FileIndex is initialized with all fields set to zero or empty.
//
// Specification: package_file_format.md: 6 File Index Section
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
// readFileIndexHeader reads the 16-byte header and returns entryCount, reserved, firstEntryOffset, totalRead, error.
func readFileIndexHeader(r io.Reader) (entryCount, reserved uint32, firstEntryOffset uint64, totalRead int64, err error) {
	if err = binary.Read(r, binary.LittleEndian, &entryCount); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return 0, 0, 0, 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to read entry count: incomplete data", pkgerrors.ValidationErrorContext{
				Field: "EntryCount", Value: int64(0), Expected: "4 bytes",
			})
		}
		return 0, 0, 0, 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read entry count", pkgerrors.ValidationErrorContext{
			Field: "EntryCount", Value: nil, Expected: "4 bytes",
		})
	}
	totalRead = 4
	if err = binary.Read(r, binary.LittleEndian, &reserved); err != nil {
		return 0, 0, 0, totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read reserved", pkgerrors.ValidationErrorContext{
			Field: "Reserved", Value: nil, Expected: "4 bytes",
		})
	}
	totalRead += 4
	if err = binary.Read(r, binary.LittleEndian, &firstEntryOffset); err != nil {
		return 0, 0, 0, totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read first entry offset", pkgerrors.ValidationErrorContext{
			Field: "FirstEntryOffset", Value: nil, Expected: "8 bytes",
		})
	}
	totalRead += 8
	return entryCount, reserved, firstEntryOffset, totalRead, nil
}

// validateEntryCountAllocation checks entryCount against allocation limits; returns an error if allocation would be unsafe.
func validateEntryCountAllocation(entryCount uint32, totalRead int64) error {
	const maxInt = int(^uint(0) >> 1)
	if int(entryCount) > maxInt {
		return pkgerrors.WrapErrorWithContext(
			fmt.Errorf("entry count %d exceeds maximum slice size %d", entryCount, maxInt),
			pkgerrors.ErrTypeValidation, "entry count exceeds system allocation limits",
			pkgerrors.ValidationErrorContext{Field: "EntryCount", Value: entryCount, Expected: fmt.Sprintf("value <= %d", maxInt)},
		)
	}
	requiredBytes := uint64(entryCount) * uint64(IndexEntrySize)
	if entryCount > 0 && int(entryCount) > maxInt/int(IndexEntrySize) {
		return pkgerrors.WrapErrorWithContext(
			fmt.Errorf("entry count %d would require allocation exceeding maximum slice size", entryCount),
			pkgerrors.ErrTypeValidation, "entry count exceeds maximum allocation size",
			pkgerrors.ValidationErrorContext{Field: "EntryCount", Value: entryCount, Expected: fmt.Sprintf("value <= %d", maxInt/int(IndexEntrySize))},
		)
	}
	if requiredBytes <= 1024*1024*1024 {
		return nil
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	maxReasonableAllocation := uint64(10 * 1024 * 1024 * 1024)
	if memStats.Sys > 0 && memStats.Sys < maxReasonableAllocation*2 {
		maxReasonableAllocation = memStats.Sys / 2
	}
	if requiredBytes > maxReasonableAllocation {
		return pkgerrors.WrapErrorWithContext(
			fmt.Errorf("entry count %d would require %d bytes (%d GB), exceeding available system memory", entryCount, requiredBytes, requiredBytes/(1024*1024*1024)),
			pkgerrors.ErrTypeValidation, "entry count exceeds available system memory",
			pkgerrors.ValidationErrorContext{Field: "EntryCount", Value: entryCount, Expected: "value within available system memory constraints"},
		)
	}
	return nil
}

// Specification: package_file_format.md: 6 File Index Section
func (f *FileIndex) readFrom(r io.Reader) (int64, error) {
	entryCount, reserved, firstEntryOffset, totalRead, err := readFileIndexHeader(r)
	if err != nil {
		return totalRead, err
	}
	f.EntryCount = entryCount
	f.Reserved = reserved
	f.FirstEntryOffset = firstEntryOffset
	if err := validateEntryCountAllocation(entryCount, totalRead); err != nil {
		return totalRead, err
	}
	f.Entries = make([]IndexEntry, 0, entryCount)
	for i := uint32(0); i < entryCount; i++ {
		var entry IndexEntry
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to read entry %d", i), pkgerrors.ValidationErrorContext{
				Field: "Entries", Value: i, Expected: "valid index entry",
			})
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
// Specification: package_file_format.md: 6 File Index Section
func (f *FileIndex) writeTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Update EntryCount to match actual entries
	f.EntryCount = uint32(len(f.Entries))

	// Write header (16 bytes)
	if err := binary.Write(w, binary.LittleEndian, f.EntryCount); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write entry count", pkgerrors.ValidationErrorContext{
			Field:    "EntryCount",
			Value:    f.EntryCount,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, f.Reserved); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write reserved", pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    f.Reserved,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, f.FirstEntryOffset); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write first entry offset", pkgerrors.ValidationErrorContext{
			Field:    "FirstEntryOffset",
			Value:    f.FirstEntryOffset,
			Expected: "written successfully",
		})
	}
	totalWritten += 8

	// Write entries
	for i, entry := range f.Entries {
		if err := binary.Write(w, binary.LittleEndian, entry); err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to write entry %d", i), pkgerrors.ValidationErrorContext{
				Field:    "Entries",
				Value:    i,
				Expected: "written successfully",
			})
		}
		totalWritten += IndexEntrySize
	}

	return totalWritten, nil
}
