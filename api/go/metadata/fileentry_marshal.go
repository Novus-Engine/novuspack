// This file implements marshaling and unmarshaling operations for FileEntry
// structures. It contains methods for serializing FileEntry metadata and data
// to binary format and reading them back. This file should contain Marshal,
// MarshalMeta, MarshalData, UnmarshalFileEntry, and streaming methods as
// specified in api_file_mgmt_file_entry.md Section 1.6.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure

package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// writeSliceToWriter writes slice entries to w using writeAt(i) for each index; updates totalWritten.
func writeSliceToWriter(w io.Writer, totalWritten int64, n int, fieldName string, writeAt func(i int) (int64, error)) (int64, error) {
	for i := 0; i < n; i++ {
		written, err := writeAt(i)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to write %s entry %d", fieldName, i), pkgerrors.ValidationErrorContext{
				Field:    fieldName,
				Value:    i,
				Expected: "written successfully",
			})
		}
		totalWritten += written
	}
	return totalWritten, nil
}

// MarshalMeta marshals the FileEntry metadata (header + variable data) to bytes.
//
// Marshals the complete FileEntry metadata structure including:
//   - Fixed 64-byte header
//   - Path entries
//   - Hash data
//   - Optional data (including tags)
//
// Returns complete binary representation of metadata ready for writing to package file.
//
// Error conditions:
//   - ErrTypeValidation: Invalid FileEntry state
//   - ErrTypeIO: I/O error during marshaling
//
// Use this method when you need the metadata as a byte slice for in-memory operations or when metadata size is small.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) MarshalMeta() ([]byte, error) {
	var buf bytes.Buffer
	_, err := f.WriteMetaTo(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalData marshals the FileEntry data (file content) to bytes.
//
// Marshals the file data content (already processed with compression/encryption).
//
// Returns binary representation of file data ready for writing to package file.
//
// Note: Data should already be processed (compressed/encrypted) via ProcessData() before marshaling.
//
// Error conditions:
//   - ErrTypeValidation: FileEntry data not available or not processed
//   - ErrTypeIO: I/O error during marshaling
//
// Use this method when you need the data as a byte slice for in-memory operations or when data size is small.
//
// Specification: api_file_mgmt_file_entry.md: 6. Marshaling
func (f *FileEntry) MarshalData() ([]byte, error) {
	if !f.IsDataLoaded {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file entry data not available", nil, pkgerrors.ValidationErrorContext{
			Field:    "Data",
			Value:    nil,
			Expected: "data loaded or available",
		})
	}

	// Return a copy of the data
	result := make([]byte, len(f.Data))
	copy(result, f.Data)
	return result, nil
}

// Marshal marshals both FileEntry metadata and data.
//
// Returns both metadata and data as separate byte slices for flexible writing.
//
// Convenience method that calls MarshalMeta() and MarshalData() internally.
// Use when you need both metadata and data marshaled together.
//
// Error conditions:
//   - ErrTypeValidation: Invalid FileEntry state or data not available
//   - ErrTypeIO: I/O error during marshaling
//
// Specification: api_file_mgmt_file_entry.md: 6. Marshaling
func (f *FileEntry) Marshal() (metadata, data []byte, err error) {
	metadata, err = f.MarshalMeta()
	if err != nil {
		return nil, nil, err
	}

	data, err = f.MarshalData()
	if err != nil {
		return nil, nil, err
	}

	return metadata, data, nil
}

// WriteMetaTo writes the file entry metadata to a writer using streaming.
//
// Writes file entry metadata to a writer using streaming.
// Implements efficient streaming for large metadata.
//
// Parameters:
//   - w: Writer to write metadata to
//
// Returns number of bytes written and error.
//
// Error conditions:
//   - ErrTypeValidation: Invalid FileEntry state
//   - ErrTypeIO: I/O error during writing
//
// Use this method for memory-efficient streaming of large metadata.
// Follows Go's standard io.WriterTo pattern.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func (f *FileEntry) WriteMetaTo(w io.Writer) (int64, error) {
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
		hashSize += h.size()
	}
	f.HashDataOffset = uint32(pathsSize)
	f.HashDataLen = uint16(hashSize)

	optionalDataSize := 0
	for _, o := range f.OptionalData {
		optionalDataSize += o.size()
	}
	f.OptionalDataOffset = uint32(pathsSize + hashSize)
	f.OptionalDataLen = uint16(optionalDataSize)

	// Write fixed section (64 bytes)
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
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write fixed section", pkgerrors.ValidationErrorContext{
			Field:    "FixedSection",
			Value:    nil,
			Expected: "written successfully",
		})
	}
	totalWritten += FileEntryFixedSize

	// Write path entries (starting at offset 0)
	var err error
	totalWritten, err = writeSliceToWriter(w, totalWritten, len(f.Paths), "Paths", func(i int) (int64, error) { return f.Paths[i].WriteTo(w) })
	if err != nil {
		return totalWritten, err
	}

	// Write hash entries (starting at HashDataOffset)
	totalWritten, err = writeSliceToWriter(w, totalWritten, len(f.Hashes), "Hashes", func(i int) (int64, error) { return f.Hashes[i].writeTo(w) })
	if err != nil {
		return totalWritten, err
	}

	// Write optional data entries (starting at OptionalDataOffset)
	totalWritten, err = writeSliceToWriter(w, totalWritten, len(f.OptionalData), "OptionalData", func(i int) (int64, error) { return f.OptionalData[i].writeTo(w) })
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

// WriteDataTo writes the file entry data to a writer using streaming.
//
// Writes file entry data to a writer using streaming.
// Implements efficient streaming for large files.
//
// Parameters:
//   - w: Writer to write data to
//
// Returns number of bytes written and error.
//
// Error conditions:
//   - ErrTypeValidation: FileEntry data not available or not processed
//   - ErrTypeIO: I/O error during writing
//
// Use this method for memory-efficient streaming of large files.
// Follows Go's standard io.WriterTo pattern.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
//nolint:gocognit // branch count from data source and validation paths
func (f *FileEntry) WriteDataTo(w io.Writer) (int64, error) {
	// If data is in memory, write it directly
	if f.IsDataLoaded {
		n, err := w.Write(f.Data)
		if err != nil {
			return int64(n), pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write data", pkgerrors.ValidationErrorContext{
				Field:    "Data",
				Value:    len(f.Data),
				Expected: "written successfully",
			})
		}
		return int64(n), nil
	}

	// If data is in a source file, stream from it
	if f.SourceFile != nil && f.SourceSize > 0 {
		// Seek to the source offset
		if _, err := f.SourceFile.Seek(f.SourceOffset, io.SeekStart); err != nil {
			return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek to source offset", pkgerrors.ValidationErrorContext{
				Field:    "SourceOffset",
				Value:    f.SourceOffset,
				Expected: "seek successful",
			})
		}

		// Copy the data
		n, err := io.CopyN(w, f.SourceFile, f.SourceSize)
		if err != nil {
			return n, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to copy data from source file", pkgerrors.ValidationErrorContext{
				Field:    "SourceFile",
				Value:    f.SourceSize,
				Expected: "copy successful",
			})
		}
		return n, nil
	}

	// If data is in a temp file, read and write it
	if f.IsTempFile && f.TempFilePath != "" {
		tmpFile, err := os.Open(f.TempFilePath)
		if err != nil {
			return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to open temporary file for reading", pkgerrors.ValidationErrorContext{
				Field:    "TempFilePath",
				Value:    f.TempFilePath,
				Expected: "file opened successfully",
			})
		}
		defer func() { _ = tmpFile.Close() }() // Close on exit - error is non-critical

		n, err := io.Copy(w, tmpFile)
		if err != nil {
			return n, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to copy data from temporary file", pkgerrors.ValidationErrorContext{
				Field:    "TempFile",
				Value:    f.TempFilePath,
				Expected: "copy successful",
			})
		}
		return n, nil
	}

	return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file entry data not available", nil, pkgerrors.ValidationErrorContext{
		Field:    "Data",
		Value:    nil,
		Expected: "data loaded, source file available, or temp file available",
	})
}

// UnmarshalFileEntry unmarshals a FileEntry from binary data with proper tag synchronization.
//
// Unmarshals a FileEntry from binary data with proper tag synchronization.
// Unmarshaling always handles tag synchronization appropriately, ensuring tags are properly loaded and associated with the FileEntry.
//
// Parameters:
//   - data: Binary data containing the FileEntry structure
//
// Returns the unmarshaled FileEntry and error.
//
// Note: Unmarshaling always handles tag synchronization appropriately, ensuring tags are properly loaded and associated with the FileEntry.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
func UnmarshalFileEntry(data []byte) (*FileEntry, error) {
	fe := NewFileEntry()
	reader := bytes.NewReader(data)

	_, err := fe.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	// Ensure tag synchronization
	// This triggers tag loading from OptionalData
	// Tag loading errors are non-fatal - tags may not exist, so we ignore the error
	_, _ = GetFileEntryTags(fe)

	return fe, nil
}
