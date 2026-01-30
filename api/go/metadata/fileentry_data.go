// This file implements data management operations for FileEntry structures.
// It contains methods for loading, unloading, and managing file data in memory
// and temporary files. This file should contain data management methods
// (LoadData, UnloadData, GetData, SetData, temp file operations) as specified
// in api_file_mgmt_file_entry.md Section 1.4.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions

package metadata

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// LoadData loads the file data into memory.
//
// Loads file content from package file into memory.
// Prepares data for access and processing.
// May trigger decompression or decryption.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns *PackageError on failure.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) LoadData(ctx context.Context) error {
	// Check if already loaded
	if f.IsDataLoaded && len(f.Data) > 0 {
		return nil
	}

	// Check if we have a source file
	if f.SourceFile == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "no source file available for loading data", nil, pkgerrors.ValidationErrorContext{
			Field:    "SourceFile",
			Value:    nil,
			Expected: "source file available",
		})
	}

	// Check context
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	// Seek to source offset
	if _, err := f.SourceFile.Seek(f.SourceOffset, io.SeekStart); err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek to source offset", pkgerrors.ValidationErrorContext{
			Field:    "SourceOffset",
			Value:    f.SourceOffset,
			Expected: "seek successful",
		})
	}

	// Read data
	f.ProcessingState = ProcessingStateLoading
	data := make([]byte, f.SourceSize)
	n, err := io.ReadFull(f.SourceFile, data)
	if err != nil && err != io.EOF {
		f.ProcessingState = ProcessingStateError
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read data from source file", pkgerrors.ValidationErrorContext{
			Field:    "SourceFile",
			Value:    n,
			Expected: fmt.Sprintf("%d bytes", f.SourceSize),
		})
	}

	if int64(n) != f.SourceSize {
		f.ProcessingState = ProcessingStateError
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete data read", nil, pkgerrors.ValidationErrorContext{
			Field:    "Data",
			Value:    n,
			Expected: fmt.Sprintf("%d bytes", f.SourceSize),
		})
	}

	f.Data = data
	f.IsDataLoaded = true
	f.ProcessingState = ProcessingStateComplete

	return nil
}

// UnloadData unloads the file data from memory.
//
// Clears file content from memory.
// Releases memory resources.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) UnloadData() {
	f.Data = nil
	f.IsDataLoaded = false
	f.ProcessingState = ProcessingStateIdle
}

// GetData returns the file data.
//
// Returns the file content as a byte slice.
// Loads data if not already loaded.
//
// Returns:
//   - []byte: File content
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) GetData() ([]byte, error) {
	// If data is already loaded, return it
	if f.IsDataLoaded && len(f.Data) > 0 {
		result := make([]byte, len(f.Data))
		copy(result, f.Data)
		return result, nil
	}

	// Try to load data if source file is available
	if f.SourceFile != nil {
		if err := f.LoadData(context.Background()); err != nil {
			return nil, err
		}
		result := make([]byte, len(f.Data))
		copy(result, f.Data)
		return result, nil
	}

	return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file entry data not available", nil, pkgerrors.ValidationErrorContext{
		Field:    "Data",
		Value:    nil,
		Expected: "data loaded or source file available",
	})
}

// SetData sets the file data.
//
// Sets the file content in memory.
// Updates IsDataLoaded flag.
//
// Parameters:
//   - data: File content to set
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) SetData(data []byte) {
	f.Data = data
	f.IsDataLoaded = len(data) > 0
	if f.IsDataLoaded {
		f.ProcessingState = ProcessingStateComplete
	} else {
		f.ProcessingState = ProcessingStateIdle
	}
}

// GetProcessingState returns the current processing state.
//
// Returns:
//   - ProcessingState: Current processing state
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) GetProcessingState() ProcessingState {
	return f.ProcessingState
}

// SetProcessingState sets the processing state.
//
// Parameters:
//   - state: Processing state to set
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) SetProcessingState(state ProcessingState) {
	f.ProcessingState = state
}

// SetSourceFile sets the source file handle and offset information.
//
// Parameters:
//   - file: Source file handle
//   - offset: Offset in source file
//   - size: Size of data to read from source
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) SetSourceFile(file *os.File, offset, size int64) {
	f.SourceFile = file
	f.SourceOffset = offset
	f.SourceSize = size
}

// GetSourceFile returns the source file handle and offset information.
//
// Returns:
//   - *os.File: Source file handle
//   - int64: Offset in source file
//   - int64: Size of data to read from source
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) GetSourceFile() (*os.File, int64, int64) {
	return f.SourceFile, f.SourceOffset, f.SourceSize
}

// SetTempPath sets the temporary file path.
//
// Parameters:
//   - path: Path to temporary file
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) SetTempPath(path string) {
	f.TempFilePath = path
	f.IsTempFile = path != ""
}

// GetTempPath returns the temporary file path.
//
// Returns:
//   - string: Path to temporary file (empty if not set)
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) GetTempPath() string {
	return f.TempFilePath
}

// CreateTempFile creates a temporary file for this file entry.
//
// Creates a temporary file for storing file data during processing.
// Sets TempFilePath and IsTempFile flag.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns *PackageError on failure.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) CreateTempFile(ctx context.Context) error {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	tmpFile, err := os.CreateTemp("", "novuspack-fileentry-*")
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to create temporary file", pkgerrors.ValidationErrorContext{
			Field:    "TempFile",
			Value:    nil,
			Expected: "temp file created",
		})
	}

	f.TempFilePath = tmpFile.Name()
	f.IsTempFile = true
	_ = tmpFile.Close() //nolint:errcheck // Close after getting name - error is non-critical

	return nil
}

// StreamToTempFile streams data to a temporary file.
//
// Streams file data from source to temporary file.
// Updates TempFilePath and IsTempFile flag.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns *PackageError on failure.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) StreamToTempFile(ctx context.Context) error {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	if f.SourceFile == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "no source file available for streaming", nil, pkgerrors.ValidationErrorContext{
			Field:    "SourceFile",
			Value:    nil,
			Expected: "source file available",
		})
	}

	// Create temp file if not exists
	if f.TempFilePath == "" {
		if err := f.CreateTempFile(ctx); err != nil {
			return err
		}
	}

	tmpFile, err := os.OpenFile(f.TempFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to open temporary file for writing", pkgerrors.ValidationErrorContext{
			Field:    "TempFilePath",
			Value:    f.TempFilePath,
			Expected: "file opened successfully",
		})
	}
	defer func() { _ = tmpFile.Close() }() // Close on exit - error is non-critical

	// Seek to source offset
	if _, err := f.SourceFile.Seek(f.SourceOffset, io.SeekStart); err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek to source offset", pkgerrors.ValidationErrorContext{
			Field:    "SourceOffset",
			Value:    f.SourceOffset,
			Expected: "seek successful",
		})
	}

	// Copy data
	f.ProcessingState = ProcessingStateProcessing
	_, err = io.CopyN(tmpFile, f.SourceFile, f.SourceSize)
	if err != nil {
		f.ProcessingState = ProcessingStateError
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to copy data to temporary file", pkgerrors.ValidationErrorContext{
			Field:    "TempFile",
			Value:    f.TempFilePath,
			Expected: "copy successful",
		})
	}

	f.IsTempFile = true
	f.ProcessingState = ProcessingStateComplete

	return nil
}

// WriteToTempFile writes data to a temporary file.
//
// Writes the provided data to a temporary file.
// Creates temp file if needed.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - data: Data to write to temporary file
//
// Returns *PackageError on failure.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) WriteToTempFile(ctx context.Context, data []byte) error {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	// Create temp file if not exists
	if f.TempFilePath == "" {
		if err := f.CreateTempFile(ctx); err != nil {
			return err
		}
	}

	tmpFile, err := os.OpenFile(f.TempFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to open temporary file for writing", pkgerrors.ValidationErrorContext{
			Field:    "TempFilePath",
			Value:    f.TempFilePath,
			Expected: "file opened successfully",
		})
	}
	defer func() { _ = tmpFile.Close() }() // Close on exit - error is non-critical

	f.ProcessingState = ProcessingStateWriting
	_, err = tmpFile.Write(data)
	if err != nil {
		f.ProcessingState = ProcessingStateError
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write data to temporary file", pkgerrors.ValidationErrorContext{
			Field:    "TempFile",
			Value:    f.TempFilePath,
			Expected: "write successful",
		})
	}

	f.IsTempFile = true
	f.ProcessingState = ProcessingStateComplete

	return nil
}

// ReadFromTempFile reads data from a temporary file.
//
// Reads data from the temporary file at the specified offset and size.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - offset: Offset in temporary file
//   - size: Size of data to read
//
// Returns:
//   - []byte: Data read from temporary file
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) ReadFromTempFile(ctx context.Context, offset, size int64) ([]byte, error) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	if f.TempFilePath == "" {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "no temporary file available", nil, pkgerrors.ValidationErrorContext{
			Field:    "TempFilePath",
			Value:    nil,
			Expected: "temp file path set",
		})
	}

	tmpFile, err := os.Open(f.TempFilePath)
	if err != nil {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to open temporary file for reading", pkgerrors.ValidationErrorContext{
			Field:    "TempFilePath",
			Value:    f.TempFilePath,
			Expected: "file opened successfully",
		})
	}
	defer func() { _ = tmpFile.Close() }() // Close on exit - error is non-critical

	// Seek to offset
	if _, err := tmpFile.Seek(offset, io.SeekStart); err != nil {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek in temporary file", pkgerrors.ValidationErrorContext{
			Field:    "Offset",
			Value:    offset,
			Expected: "seek successful",
		})
	}

	// Read data
	data := make([]byte, size)
	n, err := io.ReadFull(tmpFile, data)
	if err != nil && err != io.EOF {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read from temporary file", pkgerrors.ValidationErrorContext{
			Field:    "TempFile",
			Value:    n,
			Expected: fmt.Sprintf("%d bytes", size),
		})
	}

	return data[:n], nil
}

// CleanupTempFile removes the temporary file.
//
// Removes the temporary file and clears related fields.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns *PackageError on failure.
//
// Specification: api_file_mgmt_file_entry.md: 1.4 Helper Functions
func (f *FileEntry) CleanupTempFile(ctx context.Context) error {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), pkgerrors.ValidationErrorContext{
				Field:    "Context",
				Value:    nil,
				Expected: "context not cancelled",
			})
		default:
		}
	}

	if f.TempFilePath == "" {
		return nil // Nothing to clean up
	}

	err := os.Remove(f.TempFilePath)
	if err != nil && !os.IsNotExist(err) {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to remove temporary file", pkgerrors.ValidationErrorContext{
			Field:    "TempFilePath",
			Value:    f.TempFilePath,
			Expected: "file removed successfully",
		})
	}

	f.TempFilePath = ""
	f.IsTempFile = false

	return nil
}
