// This file implements package lifecycle operations: Create, Open, Close, and
// related state management methods. It contains the implementation of basic
// package operations as specified in api_basic_operations.md. This file should
// contain methods for creating new packages, opening existing packages, closing
// packages, and managing package state transitions.
//
// Specification: api_basic_operations.md: 1. Context Integration

// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains package lifecycle operations: Create, Open, Close, and related methods.
package novus_package

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
	"github.com/novus-engine/novuspack/api/go/signatures"
)

// IsOpen returns true if the package is currently open.
//
// This is a pure data method that does not perform I/O, so it does not
// accept a context parameter.
//
// Specification: api_basic_operations.md: 1. Context Integration
func (p *filePackage) IsOpen() bool {
	return p.isOpen
}

// IsReadOnly returns false for filePackage (writable package).
//
// filePackage instances are writable packages. For read-only packages,
// use OpenPackageReadOnly which returns a readOnlyPackage wrapper.
//
// Returns:
//   - bool: Always false for filePackage
//
// Specification: api_basic_operations.md: 3.3 Package Implementation Details
func (p *filePackage) IsReadOnly() bool {
	return false
}

// GetPath returns the current package file path.
//
// This is a pure data accessor method that returns the configured file path
// for the package. The path is set during Create/CreateWithOptions or OpenPackage.
//
// Returns:
//   - string: The package file path, or empty string if not configured
//
// Specification: api_basic_operations.md: 3.3 Package Implementation Details
func (p *filePackage) GetPath() string {
	return p.FilePath
}

// Create configures the package for creation at the specified path.
//
// This method validates the path and initializes the package structure in memory.
// **This function does not write to disk** - it only prepares the package structure
// and sets the target path. The package file is only written to disk when one of
// the Write functions (Write, SafeWrite, or FastWrite) is called.
//
// This method validates the path, initializes the package header with creation
// timestamps and default values in memory. The package transitions from the "New"
// state to the "Created" state.
//
// State Transition: New → Created
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: File path where the package will be written (when Write is called)
//
// Returns:
//   - error: Error if path is invalid or context is cancelled
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Path is empty or invalid
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
//	// Package is configured in memory, not yet written to disk
//	// Call Write() to actually write to disk
//
// Specification: api_basic_operations.md: 7.2 NewPackageWithOptions Behavior
func (p *filePackage) Create(ctx context.Context, path string) error {
	// Validate context
	if err := internal.CheckContext(ctx, "Create"); err != nil {
		return err
	}

	// Validate path
	if err := internal.ValidatePath(ctx, path); err != nil {
		return err
	}

	// Trim and normalize path
	path = strings.TrimSpace(path)
	p.FilePath = path

	// Initialize timestamps
	now := time.Now()
	p.header.CreatedTime = uint64(now.UnixNano())
	p.header.ModifiedTime = uint64(now.UnixNano())

	// Update package info
	p.Info.Created = now
	p.Info.Modified = now

	// Initialize index to empty (in memory only)
	if p.index == nil {
		p.index = fileformat.NewFileIndex()
	}
	p.index.EntryCount = 0
	p.index.FirstEntryOffset = uint64(fileformat.PackageHeaderSize)

	// Calculate index position (right after header) - for future Write operations
	indexStart := uint64(fileformat.PackageHeaderSize)
	indexSize := fileIndexSize(p.index)

	// Update header with index location (in memory only)
	p.header.IndexStart = indexStart
	p.header.IndexSize = indexSize

	// Package is configured and ready for operations (in memory)
	// Even though not written to disk, the package is "open" for metadata and file operations
	p.isOpen = true

	return nil
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
//	info, err := pkg.GetInfo(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Files: %d\n", info.FileCount)
//
// Specification: api_basic_operations.md: 10. OpenPackage Function
//
//nolint:gocognit,gocyclo // open/read/validate branches
func OpenPackage(ctx context.Context, path string) (Package, error) {
	// Validate context
	if err := internal.CheckContext(ctx, "OpenPackage"); err != nil {
		return nil, err
	}

	// Validate and normalize path
	if err := internal.ValidatePath(ctx, path); err != nil {
		return nil, err
	}
	path = strings.TrimSpace(path)

	// Open file using helper function
	file, err := internal.OpenFileForReading(path)
	if err != nil {
		return nil, err
	}

	// Read and validate header using helper function
	header, err := internal.ReadAndValidateHeader(ctx, file)
	if err != nil {
		_ = file.Close() // Ignore error on cleanup path
		return nil, err
	}

	// Load file index if it exists
	index := fileformat.NewFileIndex()
	if header.IndexStart > 0 && header.IndexSize > 0 {
		// Seek to index start
		if _, err = file.Seek(int64(header.IndexStart), 0); err != nil {
			_ = file.Close() // Ignore error on cleanup path
			return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "failed to seek to file index", err, struct{}{})
		}

		// Read the index
		if _, err = readFileIndexFrom(file, index); err != nil {
			_ = file.Close() // Ignore error on cleanup path
			return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "failed to read file index", err, struct{}{})
		}

		// Validate the index
		if err := validateFileIndex(index); err != nil {
			_ = file.Close() // Ignore error on cleanup path
			return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid file index", err, struct{}{})
		}
	}

	// Create package instance
	pkg := &filePackage{
		FileEntries:  make([]*metadata.FileEntry, 0),
		SpecialFiles: make(map[uint16]*metadata.FileEntry),
		FilePath:     path,
		header:       header,
		index:        index,
		fileHandle:   file,
		isOpen:       true,
	}

	// Initialize package info and sync from header (header is source on disk, PackageInfo is source in memory)
	pkg.Info = metadata.NewPackageInfo()
	pkg.Info.Created = time.Unix(0, int64(header.CreatedTime))
	pkg.Info.Modified = time.Unix(0, int64(header.ModifiedTime))
	pkg.Info.VendorID = header.VendorID
	pkg.Info.AppID = header.AppID
	pkg.Info.PackageDataVersion = header.PackageDataVersion
	pkg.Info.MetadataVersion = header.MetadataVersion

	// Load all FileEntry metadata from index (eager loading)
	if index.EntryCount > 0 {
		pkg.FileEntries = make([]*metadata.FileEntry, 0, index.EntryCount)
		var totalUncompressedSize int64
		var totalCompressedSize int64
		hasEncryptedData := false
		hasCompressedData := false
		hasMetadataFiles := false

		for _, indexEntry := range index.Entries {
			// Check context cancellation
			select {
			case <-ctx.Done():
				_ = file.Close()
				return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled during file entry loading", ctx.Err(), struct{}{})
			default:
			}

			// Load FileEntry from offset
			entry, err := internal.LoadFileEntry(file, indexEntry.Offset)
			if err != nil {
				_ = file.Close()
				return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to load file entry for FileID %d", indexEntry.FileID), pkgerrors.ValidationErrorContext{
					Field:    "FileID",
					Value:    indexEntry.FileID,
					Expected: "valid file entry",
				})
			}

			// Populate runtime-only offset fields for efficient file data access
			entry.EntryOffset = indexEntry.Offset
			entry.SourceFile = file
			entry.SourceOffset = int64(indexEntry.Offset + uint64(entry.TotalSize()))
			entry.SourceSize = int64(entry.StoredSize)

			// Add to FileEntries
			pkg.FileEntries = append(pkg.FileEntries, entry)

			// Accumulate size statistics
			totalUncompressedSize += int64(entry.OriginalSize)
			totalCompressedSize += int64(entry.StoredSize)

			// Check for encrypted/compressed data
			if entry.EncryptionType != 0 {
				hasEncryptedData = true
			}
			if entry.CompressionType != 0 {
				hasCompressedData = true
			}

			// Check for special metadata files (type 65000-65535)
			if entry.Type >= 65000 {
				hasMetadataFiles = true
				pkg.SpecialFiles[uint16(entry.Type)] = entry
			}
		}

		// Update package info with loaded data
		pkg.Info.FileCount = len(pkg.FileEntries)
		pkg.Info.FilesUncompressedSize = totalUncompressedSize
		pkg.Info.FilesCompressedSize = totalCompressedSize
		pkg.Info.HasEncryptedData = hasEncryptedData
		pkg.Info.HasCompressedData = hasCompressedData
		pkg.Info.HasMetadataFiles = hasMetadataFiles
		pkg.Info.IsMetadataOnly = (pkg.Info.FileCount > 0 && hasMetadataFiles && pkg.Info.FileCount == len(pkg.SpecialFiles))
	} else {
		pkg.Info.FileCount = 0
	}

	// Load package comment if it exists
	if header.CommentSize > 0 {
		if _, err := file.Seek(int64(header.CommentStart), 0); err != nil {
			_ = file.Close()
			return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek to comment", pkgerrors.ValidationErrorContext{
				Field:    "CommentStart",
				Value:    header.CommentStart,
				Expected: "seek successful",
			})
		}

		comment := metadata.NewPackageComment()
		_, err := comment.ReadFrom(file)
		if err != nil {
			_ = file.Close()
			return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read package comment", pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    nil,
				Expected: "valid comment",
			})
		}

		if err := comment.Validate(); err != nil {
			_ = file.Close()
			return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, "invalid package comment", pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    nil,
				Expected: "valid comment",
			})
		}

		pkg.Info.HasComment = true
		pkg.Info.Comment = extractCommentText(comment) // Strip null terminator
	}

	// Load signature metadata if it exists
	if header.SignatureOffset > 0 {
		// TODO: Implement full signature loading (incremental signing support)
		// For now, just mark that signatures exist
		pkg.Info.HasSignatures = true
		pkg.Info.IsImmutable = true
		// SignatureCount and Signatures will be populated when signature loading is fully implemented
		pkg.Info.SignatureCount = 0
		pkg.Info.Signatures = []signatures.SignatureInfo{}
		pkg.Info.SecurityLevel = metadata.SecurityLevelNone // Will be calculated from signatures
	}

	// Load path metadata from special metadata files
	if err := pkg.LoadPathMetadataFile(ctx); err != nil {
		_ = file.Close()
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to load path metadata", struct{}{})
	}

	// Build file-path associations
	if err := pkg.UpdateFilePathAssociations(ctx); err != nil {
		_ = file.Close()
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to build file-path associations", struct{}{})
	}

	// Update compression info from header
	pkg.Info.PackageCompression = extractCompressionType(header)
	pkg.Info.IsPackageCompressed = (pkg.Info.PackageCompression != 0)
	// TODO: Calculate PackageOriginalSize and PackageCompressedSize when package compression is implemented

	return pkg, nil
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
// Specification: api_basic_operations.md: 1. Context Integration
func OpenPackageReadOnly(ctx context.Context, path string) (Package, error) {
	pkg, err := OpenPackage(ctx, path)
	if err != nil {
		return nil, err
	}

	return &readOnlyPackage{inner: pkg}, nil
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
// Specification: api_basic_operations.md: 1. Context Integration
func OpenBrokenPackage(ctx context.Context, path string) (Package, error) {
	// Validate context
	if err := internal.CheckContext(ctx, "OpenBrokenPackage"); err != nil {
		return nil, err
	}

	// Validate and normalize path
	if err := internal.ValidatePath(ctx, path); err != nil {
		return nil, err
	}
	path = strings.TrimSpace(path)

	// Open file using helper function
	file, err := internal.OpenFileForReading(path)
	if err != nil {
		return nil, err
	}

	// Read and validate header - if this fails, we cannot proceed
	header, err := internal.ReadAndValidateHeader(ctx, file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	// Create package structure
	pkg := &filePackage{
		FilePath:     path,
		header:       header,
		fileHandle:   file,
		isOpen:       true,
		FileEntries:  make([]*metadata.FileEntry, 0),
		SpecialFiles: make(map[uint16]*metadata.FileEntry),
		Info:         metadata.NewPackageInfo(),
	}

	// Try to read index - if this fails, continue with empty index
	pkg.index = fileformat.NewFileIndex()
	if header.IndexStart > 0 && header.IndexSize > 0 {
		if _, err := file.Seek(int64(header.IndexStart), 0); err == nil {
			index := fileformat.NewFileIndex()
			if _, err := readFileIndexFrom(file, index); err == nil {
				// Only use the index if it validates successfully
				if err := validateFileIndex(index); err == nil {
					pkg.index = index
				}
			}
		}
	}

	// Populate basic Info from header
	pkg.Info.FormatVersion = header.FormatVersion
	pkg.Info.FileCount = int(pkg.index.EntryCount)
	pkg.Info.PackageCompression = extractCompressionType(header)
	pkg.Info.IsPackageCompressed = (pkg.Info.PackageCompression != 0)

	return pkg, nil
}

// readOnlyPackage is a wrapper that enforces read-only behavior for a Package.
//
// This wrapper must be the dynamic type stored behind the Package interface
// returned by OpenPackageReadOnly. This prevents callers from type-asserting
// to the writable implementation type.
type readOnlyPackage struct {
	inner Package
}

var _ Package = (*readOnlyPackage)(nil)

// ReadOnlyErrorContext provides typed context for read-only enforcement errors.
type ReadOnlyErrorContext struct {
	Operation string
}

func (p *readOnlyPackage) readOnlyError(operation string) error {
	return pkgerrors.NewPackageError(pkgerrors.ErrTypeSecurity, "package is read-only", nil, ReadOnlyErrorContext{
		Operation: operation,
	})
}

// Read operations delegate to inner package.
func (p *readOnlyPackage) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return p.inner.ReadFile(ctx, path)
}

func (p *readOnlyPackage) ListFiles() ([]FileInfo, error) {
	return p.inner.ListFiles()
}

func (p *readOnlyPackage) GetMetadata() (*metadata.PackageMetadata, error) {
	return p.inner.GetMetadata()
}

func (p *readOnlyPackage) GetInfo() (*metadata.PackageInfo, error) {
	return p.inner.GetInfo()
}

func (p *readOnlyPackage) Validate(ctx context.Context) error {
	return p.inner.Validate(ctx)
}

func (p *readOnlyPackage) Close() error {
	return p.inner.Close()
}

func (p *readOnlyPackage) IsOpen() bool {
	return p.inner.IsOpen()
}

func (p *readOnlyPackage) GetComment() string {
	return p.inner.GetComment()
}

func (p *readOnlyPackage) HasComment() bool {
	return p.inner.HasComment()
}

func (p *readOnlyPackage) GetAppID() uint64 {
	return p.inner.GetAppID()
}

func (p *readOnlyPackage) HasAppID() bool {
	return p.inner.HasAppID()
}

func (p *readOnlyPackage) GetVendorID() uint32 {
	return p.inner.GetVendorID()
}

func (p *readOnlyPackage) HasVendorID() bool {
	return p.inner.HasVendorID()
}

func (p *readOnlyPackage) GetPackageIdentity() (vendorID uint32, appID uint64) {
	return p.inner.GetPackageIdentity()
}

// Mutating operations are rejected.
func (p *readOnlyPackage) Create(ctx context.Context, path string) error {
	return p.readOnlyError("Create")
}

func (p *readOnlyPackage) Defragment(ctx context.Context) error {
	return p.readOnlyError("Defragment")
}

func (p *readOnlyPackage) Write(ctx context.Context) error {
	return p.readOnlyError("Write")
}

func (p *readOnlyPackage) SafeWrite(ctx context.Context, overwrite bool) error {
	return p.readOnlyError("SafeWrite")
}

func (p *readOnlyPackage) FastWrite(ctx context.Context) error {
	return p.readOnlyError("FastWrite")
}

func (p *readOnlyPackage) SetComment(comment string) error {
	return p.readOnlyError("SetComment")
}

func (p *readOnlyPackage) ClearComment() error {
	return p.readOnlyError("ClearComment")
}

func (p *readOnlyPackage) SetAppID(appID uint64) error {
	return p.readOnlyError("SetAppID")
}

func (p *readOnlyPackage) ClearAppID() error {
	return p.readOnlyError("ClearAppID")
}

func (p *readOnlyPackage) SetVendorID(vendorID uint32) error {
	return p.readOnlyError("SetVendorID")
}

func (p *readOnlyPackage) ClearVendorID() error {
	return p.readOnlyError("ClearVendorID")
}

func (p *readOnlyPackage) SetPackageIdentity(vendorID uint32, appID uint64) error {
	return p.readOnlyError("SetPackageIdentity")
}

func (p *readOnlyPackage) ClearPackageIdentity() error {
	return p.readOnlyError("ClearPackageIdentity")
}

// File management operations are rejected.
func (p *readOnlyPackage) AddFile(ctx context.Context, filesystemPath string, options *AddFileOptions) (*metadata.FileEntry, error) {
	return nil, p.readOnlyError("AddFile")
}

func (p *readOnlyPackage) AddFileFromMemory(ctx context.Context, path string, data []byte, options *AddFileOptions) (*metadata.FileEntry, error) {
	return nil, p.readOnlyError("AddFileFromMemory")
}

func (p *readOnlyPackage) AddFilePattern(ctx context.Context, pattern string, options *AddFileOptions) ([]*metadata.FileEntry, error) {
	return nil, p.readOnlyError("AddFilePattern")
}

func (p *readOnlyPackage) AddDirectory(ctx context.Context, dirPath string, options *AddFileOptions) ([]*metadata.FileEntry, error) {
	return nil, p.readOnlyError("AddDirectory")
}

func (p *readOnlyPackage) RemoveFile(ctx context.Context, path string) error {
	return p.readOnlyError("RemoveFile")
}

func (p *readOnlyPackage) RemoveFilePattern(ctx context.Context, pattern string) ([]string, error) {
	return nil, p.readOnlyError("RemoveFilePattern")
}

func (p *readOnlyPackage) RemoveDirectory(ctx context.Context, dirPath string, options *RemoveDirectoryOptions) ([]string, error) {
	return nil, p.readOnlyError("RemoveDirectory")
}

// Target path management is rejected.
func (p *readOnlyPackage) SetTargetPath(ctx context.Context, path string) error {
	return p.readOnlyError("SetTargetPath")
}

// Session base management delegates to inner package for reads, rejects writes.
func (p *readOnlyPackage) SetSessionBase(basePath string) error {
	return p.readOnlyError("SetSessionBase")
}

func (p *readOnlyPackage) GetSessionBase() string {
	return p.inner.GetSessionBase()
}

func (p *readOnlyPackage) ClearSessionBase() {
	// Cannot clear on read-only package - silently ignore
}

func (p *readOnlyPackage) HasSessionBase() bool {
	return p.inner.HasSessionBase()
}

func (p *readOnlyPackage) CreateWithOptions(ctx context.Context, path string, options *CreateOptions) error {
	return p.readOnlyError("CreateWithOptions")
}

func (p *readOnlyPackage) CloseWithCleanup(ctx context.Context) error {
	return p.inner.CloseWithCleanup(ctx)
}

func (p *readOnlyPackage) IsReadOnly() bool {
	return true
}

func (p *readOnlyPackage) GetPath() string {
	return p.inner.GetPath()
}

func extractCompressionType(header *fileformat.PackageHeader) uint8 {
	if header == nil {
		return 0
	}
	return uint8((header.Flags & fileformat.FlagsMaskCompressionType) >> fileformat.FlagsShiftCompressionType)
}

func fileIndexSize(index *fileformat.FileIndex) uint64 {
	if index == nil {
		return 0
	}
	return uint64(16 + len(index.Entries)*fileformat.IndexEntrySize)
}

func readFileIndexFrom(r io.Reader, index *fileformat.FileIndex) (int64, error) {
	if index == nil {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file index is nil", nil, struct{}{})
	}

	entryCount, reserved, firstEntryOffset, totalRead, err := readFileIndexHeader(r)
	if err != nil {
		return totalRead, err
	}
	index.EntryCount = entryCount
	index.Reserved = reserved
	index.FirstEntryOffset = firstEntryOffset
	if err := validateEntryCountAllocation(entryCount); err != nil {
		return totalRead, err
	}
	index.Entries = make([]fileformat.IndexEntry, 0, entryCount)
	for i := uint32(0); i < entryCount; i++ {
		var entry fileformat.IndexEntry
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to read entry %d", i), pkgerrors.ValidationErrorContext{
				Field: "Entries", Value: i, Expected: "valid index entry",
			})
		}
		totalRead += fileformat.IndexEntrySize
		index.Entries = append(index.Entries, entry)
	}
	return totalRead, nil
}

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

func validateEntryCountAllocation(entryCount uint32) error {
	const maxInt = int(^uint(0) >> 1)
	if int(entryCount) > maxInt {
		return pkgerrors.WrapErrorWithContext(
			fmt.Errorf("entry count %d exceeds maximum slice size %d", entryCount, maxInt),
			pkgerrors.ErrTypeValidation, "entry count exceeds system allocation limits",
			pkgerrors.ValidationErrorContext{Field: "EntryCount", Value: entryCount, Expected: fmt.Sprintf("value <= %d", maxInt)},
		)
	}
	requiredBytes := uint64(entryCount) * uint64(fileformat.IndexEntrySize)
	if entryCount > 0 && int(entryCount) > maxInt/int(fileformat.IndexEntrySize) {
		return pkgerrors.WrapErrorWithContext(
			fmt.Errorf("entry count %d would require allocation exceeding maximum slice size", entryCount),
			pkgerrors.ErrTypeValidation, "entry count exceeds maximum allocation size",
			pkgerrors.ValidationErrorContext{Field: "EntryCount", Value: entryCount, Expected: fmt.Sprintf("value <= %d", maxInt/int(fileformat.IndexEntrySize))},
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

func validateFileIndex(index *fileformat.FileIndex) error {
	if index == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file index is nil", nil, struct{}{})
	}
	if index.Reserved != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "reserved field must be zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    index.Reserved,
			Expected: "0",
		})
	}
	if index.EntryCount != uint32(len(index.Entries)) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "entry count mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "EntryCount",
			Value:    index.EntryCount,
			Expected: fmt.Sprintf("%d", len(index.Entries)),
		})
	}
	for i, entry := range index.Entries {
		if entry.FileID == 0 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("file ID at index %d cannot be zero", i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Entries",
				Value:    i,
				Expected: "non-zero FileID",
			})
		}
	}
	seen := make(map[uint64]int)
	for i, entry := range index.Entries {
		if prev, exists := seen[entry.FileID]; exists {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("duplicate file ID %d at indices %d and %d", entry.FileID, prev, i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Entries",
				Value:    entry.FileID,
				Expected: "unique FileID",
			})
		}
		seen[entry.FileID] = i
	}
	return nil
}

func validatePackageHeader(header *fileformat.PackageHeader) error {
	if header == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package header is nil", nil, struct{}{})
	}
	if header.Magic != fileformat.NVPKMagic {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid magic number", nil, pkgerrors.ValidationErrorContext{
			Field:    "Magic",
			Value:    fmt.Sprintf("0x%08X", header.Magic),
			Expected: fmt.Sprintf("0x%08X", fileformat.NVPKMagic),
		})
	}
	if header.FormatVersion != fileformat.FormatVersion {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "unsupported format version", nil, pkgerrors.ValidationErrorContext{
			Field:    "FormatVersion",
			Value:    header.FormatVersion,
			Expected: fmt.Sprintf("%d", fileformat.FormatVersion),
		})
	}
	if header.Reserved != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "reserved field must be 0", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    header.Reserved,
			Expected: "0",
		})
	}
	return nil
}

// Close closes the package and releases all resources.
//
// This method closes the file handle (if open), releases system resources,
// and transitions the package to the "Closed" state. After Close() is called,
// no operations can be performed on the package except Close() itself (which
// is idempotent).
//
// State Transition: Any State → Closed
//
// Idempotency:
//   - Safe to call multiple times
//   - Subsequent calls have no effect and return nil
//   - Always call Close() using defer after opening a package
//
// Returns:
//   - error: Error if closing the file handle fails (rare)
//
// Error Conditions:
//   - ErrTypeIO: Failed to close file handle (very rare, typically OS error)
//
// Example:
//
//	pkg, err := novuspack.OpenPackage(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close() // Always close to release resources
//
// Specification: api_basic_operations.md: 13. Package.Close Method
func (p *filePackage) Close() error {
	// If already closed, this is a no-op (idempotent)
	if !p.isOpen && p.fileHandle == nil {
		return nil
	}

	// Close file handle if it exists
	if p.fileHandle != nil {
		err := p.fileHandle.Close()
		p.fileHandle = nil
		if err != nil {
			// Mark as closed even on error to prevent resource leaks
			p.isOpen = false
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "failed to close file handle", err, struct{}{})
		}
	}

	// Mark as closed
	p.isOpen = false

	return nil
}

// CreateWithOptions configures a package with options for creation at the specified path.
//
// This function configures an existing package (typically created with NewPackage) with
// initial metadata and target path. **This function does not write to disk** - it only
// configures the package structure in memory. The package file is only written to disk
// when one of the Write functions (Write, SafeWrite, or FastWrite) is called.
//
// This function calls Create internally to validate the path and set up the basic package
// structure, then applies the provided options.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - path: File system path where the package will be written (when Write is called)
//   - options: Initial package configuration and metadata
//
// Returns:
//   - error: Error if path is invalid, context is cancelled, or options are invalid
//
// Specification: api_basic_operations.md: 1. Context Integration
func (p *filePackage) CreateWithOptions(ctx context.Context, path string, options *CreateOptions) error {
	// Validate context
	if err := internal.CheckContext(ctx, "CreateWithOptions"); err != nil {
		return err
	}

	// Call Create to validate path and set up basic structure
	if err := p.Create(ctx, path); err != nil {
		return err
	}

	// Apply options if provided
	if options != nil {
		if options.Comment != "" {
			p.Info.Comment = options.Comment
			p.Info.HasComment = true
			// Note: header comment is set when package is written to disk
		}
		// Info is the single source of truth
		// Header will be synced from Info when package is written to disk
		if options.VendorID != 0 {
			p.Info.VendorID = options.VendorID
		}
		if options.AppID != 0 {
			p.Info.AppID = options.AppID
		}
		// Permissions are stored for later use during Write operations
		// TODO: Store permissions in package state
	}

	return nil
}

// CloseWithCleanup closes the package and performs cleanup operations.
//
// This function closes the package file and releases all associated resources, then
// performs additional cleanup operations such as clearing caches and resetting state.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns:
//   - error: Error if closing or cleanup fails
//
// Specification: api_basic_operations.md: 14. Package.CloseWithCleanup Method
func (p *filePackage) CloseWithCleanup(ctx context.Context) error {
	// Validate context
	if err := internal.CheckContext(ctx, "CloseWithCleanup"); err != nil {
		return err
	}

	// Close the package
	if err := p.Close(); err != nil {
		return err
	}

	// Perform additional cleanup
	// Clear file entries
	p.FileEntries = nil
	p.SpecialFiles = nil

	// Reset state
	p.header = nil
	p.index = nil
	p.Info = nil

	return nil
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
	// Validate context
	if err := internal.CheckContext(ctx, "ReadHeader"); err != nil {
		return nil, err
	}

	// Read and validate header using helper function
	header, err := internal.ReadAndValidateHeader(ctx, reader)
	if err != nil {
		return nil, err
	}

	return header, nil
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
//   - Only reads the header (typically first 512 bytes)
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
	// Validate context
	if err := internal.CheckContext(ctx, "ReadHeaderFromPath"); err != nil {
		return nil, err
	}

	// Validate and normalize path
	if err := internal.ValidatePath(ctx, path); err != nil {
		return nil, err
	}
	path = strings.TrimSpace(path)

	// Open file using helper function
	file, err := internal.OpenFileForReading(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Read header from file
	return ReadHeader(ctx, file)
}
