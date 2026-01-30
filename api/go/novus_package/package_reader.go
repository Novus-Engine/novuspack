// This file implements PackageReader interface methods: ReadFile, ListFiles,
// GetMetadata, Validate, and GetInfo. It contains all read-only operations
// for accessing package contents and metadata as specified in api_core.md.
// This file should contain methods for reading files, listing package contents,
// retrieving package information, and validating package integrity.
//
// Specification: api_core.md: 1. Core Interfaces

// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains package reader operations: GetInfo, GetMetadata, ReadFile, ListFiles, and Validate.
package novus_package

import (
	"context"
	"fmt"
	"io"
	"maps"
	"sort"

	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// GetInfo returns metadata information about the package.
//
// This method returns a snapshot of package metadata including format version,
// file count, total size, and timestamps. The package must not be in the
// "Closed" state.
//
// This is a pure in-memory operation that does not perform I/O, so it does not
// accept a context parameter. All data comes from already-loaded package state.
//
// Available States: New, Created, Open (not Closed)
//
// Returned Information:
//   - Version: NovusPack format version
//   - FileCount: Total number of file entries
//   - TotalSize: Sum of all file sizes (in bytes)
//   - CreatedAt: Package creation timestamp
//   - ModifiedAt: Last modification timestamp
//   - Header: Reference to underlying package header
//
// Returns:
//   - *PackageInfo: Package metadata information
//   - error: Error if package is closed or Info is nil
//
// Error Conditions:
//   - ErrTypeValidation: Package is closed or Info is nil
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
//	fmt.Printf("Uncompressed Size: %d bytes\n", info.FilesUncompressedSize)
//	fmt.Printf("Compressed Size: %d bytes\n", info.FilesCompressedSize)
//	fmt.Printf("Created: %v\n", info.Created)
//
// Specification: api_core.md: 1.1.4 ListFiles Method Contract
func (p *filePackage) GetInfo() (*metadata.PackageInfo, error) {
	// GetInfo is an in-memory operation that is allowed after Close() as long as
	// metadata remains available, but it should not work after CloseWithCleanup().
	// Note: FilePath may be empty for newly created packages not yet written to disk.
	if p.Info == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package metadata is not loaded", nil, struct{}{})
	}

	// Info is already populated by NewPackage or OpenPackage eager loading
	// Return the package info as-is (pointer is returned, caller should not modify)
	return p.Info, nil
}

// GetMetadata retrieves comprehensive package metadata.
//
// This method returns comprehensive metadata including all package information
// plus detailed file and metadata file contents. This is a pure in-memory
// operation that does not perform I/O, so it does not accept a context parameter.
//
// This method MUST NOT perform additional disk I/O or parsing beyond what
// OpenPackage already loaded. All data comes from already-loaded package state.
//
// Returns:
//   - *PackageMetadata: Comprehensive package metadata
//   - error: Error if package is closed or metadata not available
//
// Error Conditions:
//   - ErrTypeValidation: Package is closed or metadata not loaded
//
// Specification: api_core.md: 1.1.6 GetMetadata Method Contract
func (p *filePackage) GetMetadata() (*metadata.PackageMetadata, error) {
	// GetMetadata is an in-memory operation that is allowed after Close() as long
	// as metadata remains available, but it should not work for a package that has
	// never been created or opened, or after CloseWithCleanup().
	if p.Info == nil || p.FilePath == "" {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package metadata is not loaded", nil, struct{}{})
	}

	// Construct comprehensive PackageMetadata from loaded state
	pm := metadata.NewPackageMetadata()
	pm.PackageInfo = p.Info

	// Copy FileEntries (shallow copy is sufficient - entries are already loaded)
	pm.FileEntries = make([]*metadata.FileEntry, len(p.FileEntries))
	copy(pm.FileEntries, p.FileEntries)

	// Copy PathMetadataEntries (shallow copy is sufficient)
	pm.PathMetadataEntries = make([]*metadata.PathMetadataEntry, len(p.PathMetadataEntries))
	copy(pm.PathMetadataEntries, p.PathMetadataEntries)

	// Copy SpecialFiles map
	pm.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	maps.Copy(pm.SpecialFiles, p.SpecialFiles)

	return pm, nil
}

// ReadFile reads a file from the package by path.
//
// This method reads file content from the package, applying decryption and
// decompression as needed. The path must be a valid package-internal path.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - path: Package-internal path to the file
//
// Returns:
//   - []byte: File content (decrypted and decompressed)
//   - error: Error if path is invalid, file not found, or read fails
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Path is invalid or file not found
//   - ErrTypeIO: Failed to read file data
//
// Specification: api_core.md: 1.1.3 ReadFile Method Contract
func (p *filePackage) ReadFile(ctx context.Context, path string) ([]byte, error) {
	// Validate context
	if err := internal.CheckContext(ctx, "ReadFile"); err != nil {
		return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeContext, "error during ReadFile: context validation failed")
	}

	// Check if package is open
	if !p.isOpen {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package is not open", nil, struct{}{})
	}

	// Validate path
	if err := internal.ValidatePackagePath(path); err != nil {
		return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "error during ReadFile: path validation failed")
	}

	// Normalize path for comparison
	normalizedPath, err := internal.NormalizePackagePath(path)
	if err != nil {
		return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "error during ReadFile: path normalization failed")
	}

	// Find FileEntry by normalized path
	fileEntry, err := p.findFileEntryByPath(normalizedPath)
	if err != nil {
		return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "error during ReadFile: file not found")
	}

	// Check if data is already loaded in memory (from StageFile or not-yet-written entries)
	if fileEntry.IsDataLoaded {
		// Return in-memory data directly (no need for decryption/decompression in baseline)
		// This works even when fileHandle is nil (for newly written files)
		return fileEntry.Data, nil
	}

	// Check context cancellation before I/O
	select {
	case <-ctx.Done():
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), struct{}{})
	default:
	}

	// Use the stored SourceFile and SourceOffset to locate file data
	// For opened packages, SourceFile points to the package file and SourceOffset is the file data offset
	if fileEntry.SourceFile == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file source is not available", nil, pkgerrors.ValidationErrorContext{
			Field:    "SourceFile",
			Value:    "nil",
			Expected: "valid file handle",
		})
	}

	if fileEntry.SourceOffset == 0 {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file source offset is not set", nil, pkgerrors.ValidationErrorContext{
			Field:    "SourceOffset",
			Value:    0,
			Expected: "valid file offset",
		})
	}

	// Seek to file data using stored offset
	if _, err := fileEntry.SourceFile.Seek(fileEntry.SourceOffset, 0); err != nil {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to seek to file data", pkgerrors.ValidationErrorContext{
			Field:    "SourceOffset",
			Value:    fileEntry.SourceOffset,
			Expected: "seek successful",
		})
	}

	// Check context cancellation during I/O
	select {
	case <-ctx.Done():
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), struct{}{})
	default:
	}

	// Read file data (stored size)
	data := make([]byte, fileEntry.StoredSize)
	n, err := fileEntry.SourceFile.Read(data)
	if err != nil && err != io.EOF {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read file data", pkgerrors.ValidationErrorContext{
			Field:    "StoredSize",
			Value:    fileEntry.StoredSize,
			Expected: "read successful",
		})
	}

	if uint64(n) != fileEntry.StoredSize {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete file data read", nil, pkgerrors.ValidationErrorContext{
			Field:    "Data",
			Value:    n,
			Expected: fmt.Sprintf("%d bytes", fileEntry.StoredSize),
		})
	}

	// TODO: Apply decryption if file is encrypted
	if fileEntry.EncryptionType != 0 {
		// Decryption not yet implemented
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeUnsupported, "file decryption not yet implemented", nil, pkgerrors.ValidationErrorContext{
			Field:    "EncryptionType",
			Value:    fileEntry.EncryptionType,
			Expected: "decryption support",
		})
	}

	// TODO: Apply decompression if file is compressed
	if fileEntry.CompressionType != 0 {
		// Decompression not yet implemented
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeUnsupported, "file decompression not yet implemented", nil, pkgerrors.ValidationErrorContext{
			Field:    "CompressionType",
			Value:    fileEntry.CompressionType,
			Expected: "decompression support",
		})
	}

	return data, nil
}

// ListFiles returns a list of all files in the package.
//
// This method returns a list of all files in the package, sorted by normalized
// package path alphabetically. This is a pure in-memory operation that does not
// perform I/O, so it does not accept a context parameter.
//
// The results are stable across calls when package state is unchanged.
// Files removed from in-memory package (via UnstageFile) are excluded.
//
// Returns:
//   - []FileInfo: Sorted list of file information
//   - error: Error if package is closed or path normalization fails
//
// Error Conditions:
//   - ErrTypeValidation: Package is closed or path normalization fails
//
// Specification: api_core.md: 1.1.3 ReadFile Method Contract
func (p *filePackage) ListFiles() ([]FileInfo, error) {
	// CloseWithCleanup clears in-memory state.
	// ListFiles is an in-memory operation and is allowed after Close as long as
	// the metadata cache is still available.
	// Note: FilePath may be empty for newly created packages not yet written to disk.
	if p.Info == nil || p.FileEntries == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package metadata is not loaded", nil, struct{}{})
	}

	// Build list of FileInfo from loaded FileEntries
	// Each FileEntry becomes one FileInfo (not one per path)
	fileInfos := make([]FileInfo, 0, len(p.FileEntries))

	// Iterate through all FileEntries
	for _, entry := range p.FileEntries {
		if entry == nil {
			continue
		}

		// Collect and convert all paths for this entry
		displayPaths := make([]string, 0, len(entry.Paths))

		for _, pathEntry := range entry.Paths {
			// Normalize path
			normalizedPath, err := internal.NormalizePackagePath(pathEntry.Path)
			if err != nil {
				// Skip invalid paths
				continue
			}

			displayPaths = append(displayPaths, internal.ToDisplayPath(normalizedPath))
		}

		// Skip entries with no valid paths
		if len(displayPaths) == 0 {
			continue
		}

		// Sort paths lexicographically to determine primary path
		sort.Strings(displayPaths)

		// Determine FileTypeName
		// TODO: Implement type system lookup when file type system is available
		fileTypeName := "Unknown"

		// Determine HasTags - check if any OptionalDataEntry contains tags data
		hasTags := false
		for _, opt := range entry.OptionalData {
			if opt.DataType == metadata.OptionalDataTagsData && len(opt.Data) > 0 {
				hasTags = true
				break
			}
		}

		// Create FileInfo with all fields
		fileInfos = append(fileInfos, FileInfo{
			// Basic Identification
			PrimaryPath:  displayPaths[0], // First path lexicographically
			Paths:        displayPaths,
			FileID:       entry.FileID,
			FileType:     entry.Type,
			FileTypeName: fileTypeName,

			// Size Information
			Size:       int64(entry.OriginalSize),
			StoredSize: int64(entry.StoredSize),

			// Processing Status
			IsCompressed:    entry.CompressionType != 0,
			IsEncrypted:     entry.EncryptionType != 0,
			CompressionType: entry.CompressionType,

			// Content Verification
			RawChecksum:    entry.RawChecksum,
			StoredChecksum: entry.StoredChecksum,

			// Multi-Path Support
			PathCount: entry.PathCount,

			// Version Tracking
			FileVersion:     entry.FileVersion,
			MetadataVersion: entry.MetadataVersion,

			// Metadata Indicators
			HasTags: hasTags,
		})
	}

	// Sort by primary path alphabetically
	// Use stable sort to ensure consistent ordering
	sort.SliceStable(fileInfos, func(i, j int) bool {
		return fileInfos[i].PrimaryPath < fileInfos[j].PrimaryPath
	})

	return fileInfos, nil
}

// Validate checks the package format and verifies integrity.
//
// This method performs comprehensive validation of the package structure,
// ensuring that all components (header, index, entries) are valid and
// consistent. The package must be in the "Open" state for validation.
//
// Validation Checks:
//   - Package must be in "Open" state
//   - Header structure and magic number validation
//   - File index structure validation
//   - Entry integrity checks (if loaded)
//
// Prerequisites:
//   - Package must be opened via OpenPackage()
//   - Cannot validate closed packages
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Error if validation fails or context is cancelled, nil if valid
//
// Error Conditions:
//   - ErrTypeContext: Context is cancelled or has deadline exceeded
//   - ErrTypeValidation: Package is closed, header is invalid, or index is malformed
//
// Example:
//
//	pkg, err := novuspack.OpenPackage(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close()
//	if err := pkg.Validate(ctx); err != nil {
//	    log.Fatalf("Package validation failed: %v", err)
//	}
//	fmt.Println("Package is valid")
//
// Specification: api_basic_operations.md: 9.1 Package Validation
func (p *filePackage) Validate(ctx context.Context) error {
	// Validate context
	if err := internal.CheckContext(ctx, "Validate"); err != nil {
		return err
	}

	// Check if package is open
	if !p.isOpen {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "cannot validate closed package", nil, struct{}{})
	}

	// Validate header
	if p.header == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package header is nil", nil, struct{}{})
	}

	if err := p.header.Validate(); err != nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid package header", err, struct{}{})
	}

	// Validate file index if present
	if p.index != nil {
		if err := p.index.Validate(); err != nil {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "invalid file index", err, struct{}{})
		}
	}

	return nil
}
