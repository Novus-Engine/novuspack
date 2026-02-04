// This file implements file management operations for the Package interface.
//
// Specification: api_file_mgmt_index.md: 1. File Management Document Map

package novus_package

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// AddFile adds a file to the package by reading from the filesystem.
//
// Per spec, this method opens a file handle and tracks the source location, deferring actual
// data reading to Write operations (except for encryption cases). This enables memory-efficient
// processing of large files.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - path: Filesystem-style input path used to determine the stored package path
//   - options: Optional configuration for file processing (can be nil for defaults)
//
// Returns:
//   - *metadata.FileEntry: The created file entry with complete metadata
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_addition.md: 2.1 Package.AddFile Method
//
//nolint:gocognit,gocyclo // validation and path-determination branches
func (p *filePackage) AddFile(ctx context.Context, path string, options *AddFileOptions) (*metadata.FileEntry, error) {
	// Validate context
	if err := internal.CheckContext(ctx, "AddFile"); err != nil {
		return nil, err
	}

	// Trim and validate path is not empty
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path cannot be empty or whitespace-only",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "non-empty file path",
			},
		)
	}

	// =========================================================================
	// STEP 1: Filesystem Validation and Metadata Read
	// =========================================================================

	// Determine symlink handling behavior (default: true per spec section 2.7.2.6)
	followSymlinks := true
	if options != nil {
		followSymlinks = options.FollowSymlinks.GetOrDefault(true)
	}

	// First check if path is a symlink using Lstat (doesn't follow symlinks)
	lstatInfo, err := os.Lstat(path)
	if err != nil {
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeIO,
			"AddFile: failed to stat file",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "valid, accessible file",
			},
		)
	}

	// Handle symlink behavior
	var statInfo os.FileInfo
	if lstatInfo.Mode()&os.ModeSymlink != 0 {
		// Path is a symlink
		if !followSymlinks {
			return nil, pkgerrors.NewPackageError(
				pkgerrors.ErrTypeValidation,
				"path is a symlink and FollowSymlinks is false (use symlink API to add symlinks)",
				nil,
				pkgerrors.ValidationErrorContext{
					Field:    "path",
					Value:    path,
					Expected: "non-symlink file path, or set FollowSymlinks to true",
				},
			)
		}
		// Follow the symlink to get target metadata
		statInfo, err = os.Stat(path)
		if err != nil {
			return nil, pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeIO,
				"AddFile: failed to stat symlink target",
				pkgerrors.ValidationErrorContext{
					Field:    "path",
					Value:    path,
					Expected: "valid symlink with accessible target",
				},
			)
		}
	} else {
		// Not a symlink, use lstat result
		statInfo = lstatInfo
	}

	// Verify it's a regular file
	if statInfo.IsDir() {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path is a directory, not a file (use AddDirectory for directory operations)",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "file path, not directory",
			},
		)
	}

	// Open source file handle
	sourceFile, err := os.Open(path)
	if err != nil {
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeIO,
			"AddFile: failed to open file",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "readable file",
			},
		)
	}

	// Calculated original file size from stat info
	originalSize := uint64(statInfo.Size())

	// Set initial offset (start of file)
	sourceOffset := int64(0)

	// Derive stored package path (part of step 1 per spec)
	storedPath, err := p.determineStoredPath(path, options)
	if err != nil {
		_ = sourceFile.Close()
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"AddFile: failed to determine stored package path",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "valid path with proper options configuration",
			},
		)
	}

	// Determine effective compression/encryption from options
	compressionType := uint8(0)
	encryptionType := uint8(0)
	// TODO: Extract from options when compression/encryption are implemented (Priority 5-6)

	// =========================================================================
	// STEP 2: Deduplication Check
	// =========================================================================

	var targetEntry *metadata.FileEntry
	var rawChecksum uint32

	// Check if we should skip deduplication entirely (per spec section 2.1.4.1 step 2)
	// If AllowDuplicate is true, skip deduplication and treat as unique
	allowDuplicate := false
	if options != nil {
		allowDuplicate = options.AllowDuplicate.GetOrDefault(false)
	}

	if !allowDuplicate && p.FileEntries != nil {
		// Fast filter by OriginalSize
		var potentialMatches []*metadata.FileEntry
		for _, entry := range p.FileEntries {
			if entry.OriginalSize == originalSize {
				potentialMatches = append(potentialMatches, entry)
			}
		}

		// If potential matches found, read content and calculate checksum
		if len(potentialMatches) > 0 {
			// Read from the already-open file handle
			data := make([]byte, originalSize)
			_, err := sourceFile.ReadAt(data, 0)
			if err != nil {
				_ = sourceFile.Close()
				return nil, pkgerrors.WrapErrorWithContext(
					err,
					pkgerrors.ErrTypeIO,
					"AddFile: failed to read file for deduplication check",
					pkgerrors.ValidationErrorContext{
						Field:    "path",
						Value:    path,
						Expected: "readable file",
					},
				)
			}
			rawChecksum = internal.CalculateCRC32(data)

			// Compare with potential matches
			for _, entry := range potentialMatches {
				// If existing entry doesn't have RawChecksum calculated, calculate it now
				entryChecksum := entry.RawChecksum
				if entryChecksum == 0 && entry.SourceFile != nil {
					// Read from existing entry's source file
					entryData := make([]byte, entry.OriginalSize)
					_, err := entry.SourceFile.ReadAt(entryData, entry.SourceOffset)
					if err != nil {
						// Can't read existing entry's file - skip this match
						continue
					}
					entryChecksum = internal.CalculateCRC32(entryData)
					// Store it for future comparisons
					entry.RawChecksum = entryChecksum
				}

				if entryChecksum == rawChecksum && entryChecksum != 0 {
					// Found duplicate content
					// Check if path already exists
					pathExists := false
					for _, existingPath := range entry.Paths {
						if existingPath.Path == storedPath {
							pathExists = true
							break
						}
					}

					if pathExists {
						// Path exists with same content
						allowOverwrite := false
						if options != nil {
							allowOverwrite = options.AllowOverwrite.GetOrDefault(false)
						}
						if !allowOverwrite {
							_ = sourceFile.Close()
							return nil, pkgerrors.NewPackageError(
								pkgerrors.ErrTypeValidation,
								"file already exists at specified path",
								nil,
								pkgerrors.ValidationErrorContext{
									Field:    "path",
									Value:    storedPath,
									Expected: "non-existing path or AllowOverwrite=true",
								},
							)
						}
						// Overwrite allowed - use existing entry
						targetEntry = entry
					} else {
						// Add new path to existing entry (multi-path/alias)
						targetEntry = entry
						entry.Paths = append(entry.Paths, generics.PathEntry{PathLength: uint16(len(storedPath)), Path: storedPath})
						entry.PathCount++
						entry.MetadataVersion++
					}

					// Duplicate found - skip to step 5
					break
				}
			}
		}
	}

	// =========================================================================
	// STEP 3: Conditional Encryption Processing
	// =========================================================================

	// TODO: Priority 6 - When encryption is implemented, this step will:
	// - Read file data
	// - Apply compression (if both compression + encryption)
	// - Apply encryption
	// - Write to temp file
	// - Update sourceFile to point to temp file
	// - Update sourceOffset to 0
	// - Calculate StoredSize and StoredChecksum

	// =========================================================================
	// STEP 4: FileEntry Allocation (for unique files only)
	// =========================================================================

	if targetEntry == nil {
		// No duplicate found - create new FileEntry
		newFileID := p.allocateNextFileID()
		targetEntry = metadata.NewFileEntry()
		targetEntry.FileID = newFileID
		targetEntry.Type = 0 // TODO: Determine file type from extension/content
		targetEntry.Paths = []generics.PathEntry{{PathLength: uint16(len(storedPath)), Path: storedPath}}
		targetEntry.PathCount = 1
		targetEntry.OriginalSize = originalSize
		targetEntry.RawChecksum = rawChecksum // May be 0 if no deduplication check
		targetEntry.FileVersion = 1
		targetEntry.MetadataVersion = 1

		// Set compression/encryption
		targetEntry.CompressionType = compressionType
		targetEntry.EncryptionType = encryptionType

		// For non-encryption cases, StoredSize/StoredChecksum are placeholders
		// They'll be calculated during Write operations
		if encryptionType == 0 {
			targetEntry.StoredSize = 0     // Placeholder
			targetEntry.StoredChecksum = 0 // Placeholder
		}

		// Add to package
		if p.FileEntries == nil {
			p.FileEntries = make([]*metadata.FileEntry, 0, 1)
		}
		p.FileEntries = append(p.FileEntries, targetEntry)

		if p.Info == nil {
			p.Info = metadata.NewPackageInfo()
		}
		p.Info.FileCount++
	}

	// =========================================================================
	// STEP 5: Runtime Field Finalization
	// =========================================================================

	// Set SourceFile (already opened in step 1, or updated in step 3 for encryption)
	targetEntry.SourceFile = sourceFile

	// Set SourceOffset (0 for original file, or 0 for temp file after encryption)
	targetEntry.SourceOffset = sourceOffset

	// Set SourceSize (size of data in current state)
	targetEntry.SourceSize = int64(originalSize)

	// Set TempFilePath (empty for now, would be set in step 3 for encryption)
	targetEntry.TempFilePath = "" // TODO: Set in step 3 for encryption

	// Set IsTempFile (false for now, would be true if sourceFile points to temp)
	targetEntry.IsTempFile = false // TODO: Set in step 3 for encryption

	// Set ProcessingState to track what processing has been done
	// For v1 with no compression/encryption: StateRaw
	targetEntry.ProcessingState = metadata.ProcessingStateRaw

	// Data MUST NOT be loaded (per spec)
	targetEntry.IsDataLoaded = false

	// Capture and store filesystem metadata
	if err := p.captureFilesystemMetadata(storedPath, statInfo, targetEntry, options); err != nil {
		_ = sourceFile.Close()
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeIO,
			"AddFile: failed to capture filesystem metadata",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "file with accessible metadata",
			},
		)
	}

	// Ensure path metadata entry exists
	if err := p.ensurePathMetadata(storedPath, targetEntry); err != nil {
		_ = sourceFile.Close()
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"AddFile: failed to create path metadata entry",
			pkgerrors.ValidationErrorContext{
				Field:    "storedPath",
				Value:    storedPath,
				Expected: "valid path for metadata entry",
			},
		)
	}

	return targetEntry, nil
}

// AddFileFromMemory adds a file to the package from in-memory byte data.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - path: Package-relative path for the file
//   - data: Byte slice containing the file data
//   - options: Optional configuration for file processing (can be nil for defaults)
//
// Returns:
//   - *metadata.FileEntry: The created file entry
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_addition.md: 2.2 Package.AddFileFromMemory Method
//
//nolint:gocognit,gocyclo // validation and path branches
func (p *filePackage) AddFileFromMemory(ctx context.Context, path string, data []byte, options *AddFileOptions) (*metadata.FileEntry, error) {
	// Validate context
	if err := internal.CheckContext(ctx, "AddFileFromMemory"); err != nil {
		return nil, err
	}

	// Validate path is not empty or whitespace-only
	if strings.TrimSpace(path) == "" {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path cannot be empty or whitespace-only",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "non-empty package path",
			},
		)
	}

	// Use data directly (slice header is passed by value, but references same underlying array)
	actualData := data

	// Normalize path (ensure leading /)
	normalizedPath, err := internal.NormalizePackagePath(path)
	if err != nil {
		return nil, pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"failed to normalize path for AddFileFromMemory",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "valid normalizable path",
			},
		)
	}

	// Calculate file metadata
	originalSize := uint64(len(actualData))
	rawChecksum := internal.CalculateCRC32(actualData)

	// Check for existing file with same content (deduplication)
	var targetEntry *metadata.FileEntry
	allowOverwrite := false
	if options != nil && options.AllowOverwrite.IsSet() {
		allowOverwrite = options.AllowOverwrite.GetOrDefault(false)
	}

	// Search for duplicate content
	if p.FileEntries != nil {
		for _, entry := range p.FileEntries {
			if entry.OriginalSize == originalSize && entry.RawChecksum == rawChecksum {
				// Found duplicate content - check if path already exists
				pathExists := false
				for _, existingPath := range entry.Paths {
					if existingPath.Path == normalizedPath {
						pathExists = true
						break
					}
				}

				if pathExists {
					if !allowOverwrite {
						return nil, pkgerrors.NewPackageError(
							pkgerrors.ErrTypeValidation,
							"file already exists at specified path",
							nil,
							pkgerrors.ValidationErrorContext{
								Field:    "path",
								Value:    normalizedPath,
								Expected: "non-existing path or AllowOverwrite=true",
							},
						)
					}
					// Path already exists and overwrite is allowed - just return existing entry
					targetEntry = entry
				} else {
					// Add path to existing entry (multi-path support)
					targetEntry = entry
					entry.Paths = append(entry.Paths, generics.PathEntry{PathLength: uint16(len(normalizedPath)), Path: normalizedPath})
					entry.PathCount++
				}
				break
			}
		}
	}

	// If no duplicate found, create new FileEntry
	if targetEntry == nil {
		// Allocate new FileID
		newFileID := p.allocateNextFileID()

		// Create new FileEntry
		targetEntry = metadata.NewFileEntry()
		targetEntry.FileID = newFileID
		targetEntry.Type = 0 // Default type (could be enhanced with content detection)
		targetEntry.Paths = []generics.PathEntry{{PathLength: uint16(len(normalizedPath)), Path: normalizedPath}}
		targetEntry.PathCount = 1
		targetEntry.OriginalSize = originalSize
		targetEntry.RawChecksum = rawChecksum
		targetEntry.StoredSize = originalSize    // No compression yet (Priority 5)
		targetEntry.StoredChecksum = rawChecksum // No compression yet
		targetEntry.CompressionType = 0          // No compression
		targetEntry.EncryptionType = 0           // No encryption (Priority 6)

		// Store data in memory for later write
		targetEntry.SetData(actualData)

		// Add to package's FileEntries
		if p.FileEntries == nil {
			p.FileEntries = make([]*metadata.FileEntry, 0, 1)
		}
		p.FileEntries = append(p.FileEntries, targetEntry)

		// Update package info
		if p.Info == nil {
			p.Info = metadata.NewPackageInfo()
		}
		p.Info.FileCount++
	}

	// Create or update path metadata entry
	if err := p.ensurePathMetadata(normalizedPath, targetEntry); err != nil {
		return nil, err
	}

	return targetEntry, nil
}

// AddFilePattern adds files matching a pattern to the package.
//
// STUB IMPLEMENTATION: This method validates inputs but returns ErrTypeUnsupported.
// Full implementation is deferred to Priority 2.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - pattern: Glob pattern to match files
//   - options: Optional configuration for file processing (can be nil for defaults)
//
// Returns:
//   - []*metadata.FileEntry: Slice of created file entries (stub returns nil)
//   - error: *PackageError with ErrTypeUnsupported
//
// Specification: api_file_mgmt_addition.md: 2.4 Package.AddFilePattern Method
func (p *filePackage) AddFilePattern(ctx context.Context, pattern string, options *AddFileOptions) ([]*metadata.FileEntry, error) {
	return p.addStubWithContextAndNonEmpty(ctx, "AddFilePattern", pattern, "pattern", "non-empty glob pattern", "pattern cannot be empty", "AddFilePattern full implementation deferred to Priority 2")
}

// addStubWithContextAndNonEmpty validates context and non-empty value, then returns ErrTypeUnsupported stub.
func (p *filePackage) addStubWithContextAndNonEmpty(ctx context.Context, opName, value, fieldName, expectedMsg, emptyErrMsg, stubMsg string) ([]*metadata.FileEntry, error) {
	if err := internal.CheckContext(ctx, opName); err != nil {
		return nil, err
	}
	if value == "" {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			emptyErrMsg,
			nil,
			pkgerrors.ValidationErrorContext{Field: fieldName, Value: value, Expected: expectedMsg},
		)
	}
	return nil, pkgerrors.NewPackageError[struct{}](pkgerrors.ErrTypeUnsupported, stubMsg, nil, struct{}{})
}

// AddDirectory recursively adds files from a directory to the package.
//
// STUB IMPLEMENTATION: This method validates inputs but returns ErrTypeUnsupported.
// Full implementation is deferred to Priority 2.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - dirPath: Filesystem path to the directory to add
//   - options: Optional configuration for file processing (can be nil for defaults)
//
// Returns:
//   - []*metadata.FileEntry: Slice of created file entries (stub returns nil)
//   - error: *PackageError with ErrTypeUnsupported
//
// Specification: api_file_mgmt_addition.md: 2.5 Package.AddDirectory Method
func (p *filePackage) AddDirectory(ctx context.Context, dirPath string, options *AddFileOptions) ([]*metadata.FileEntry, error) {
	return p.addStubWithContextAndNonEmpty(ctx, "AddDirectory", dirPath, "dirPath", "non-empty directory path", "directory path cannot be empty", "AddDirectory full implementation deferred to Priority 2")
}

// RemoveFile removes a file from the package.
//
// This method removes the specified path from the package. If the file entry
// has multiple paths, only the specified path is removed. If it's the last
// path for a file entry, the entire file entry is removed.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - path: Package path of the file to remove
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_removal.md: 2. RemoveFile Package Method
//
//nolint:gocognit,gocyclo // validation and removal branches
func (p *filePackage) RemoveFile(ctx context.Context, path string) error {
	// Validate context
	if err := internal.CheckContext(ctx, "RemoveFile"); err != nil {
		return err
	}

	// Validate path is not empty or whitespace-only
	if strings.TrimSpace(path) == "" {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path cannot be empty or whitespace-only",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "non-empty package path",
			},
		)
	}

	// Normalize path (add leading / if missing)
	normalizedPath, err := internal.NormalizePackagePath(path)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"failed to normalize path for RemoveFile",
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "valid normalizable path",
			},
		)
	}

	// Find the file entry with this path
	var targetEntry *metadata.FileEntry
	pathIndex := -1

	if p.FileEntries != nil {
		for _, entry := range p.FileEntries {
			for i, entryPath := range entry.Paths {
				if entryPath.Path == normalizedPath {
					targetEntry = entry
					pathIndex = i
					break
				}
			}
			if targetEntry != nil {
				break
			}
		}
	}

	// Return error if path not found
	if targetEntry == nil {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"file not found at specified path",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "existing file path",
			},
		)
	}

	// Remove the path from the Paths array
	targetEntry.Paths = append(targetEntry.Paths[:pathIndex], targetEntry.Paths[pathIndex+1:]...)
	targetEntry.PathCount--

	// If this was the last path, remove the entire file entry
	if targetEntry.PathCount == 0 {
		newFileEntries := make([]*metadata.FileEntry, 0, len(p.FileEntries)-1)
		for _, entry := range p.FileEntries {
			if entry != targetEntry {
				newFileEntries = append(newFileEntries, entry)
			}
		}
		p.FileEntries = newFileEntries
	}

	// Update path metadata associations
	if p.PathMetadataEntries != nil {
		for _, pathEntry := range p.PathMetadataEntries {
			if pathEntry.GetPath() == normalizedPath {
				// Remove association with this file entry
				newAssociations := make([]*metadata.FileEntry, 0)
				for _, assoc := range pathEntry.AssociatedFileEntries {
					if assoc != targetEntry {
						newAssociations = append(newAssociations, assoc)
					}
				}
				pathEntry.AssociatedFileEntries = newAssociations
				break
			}
		}
	}

	return nil
}

// allocateNextFileID generates the next available FileID.
func (p *filePackage) allocateNextFileID() uint64 {
	maxID := uint64(0)
	if p.FileEntries != nil {
		for _, entry := range p.FileEntries {
			if entry.FileID > maxID {
				maxID = entry.FileID
			}
		}
	}
	return maxID + 1
}

// ensurePathMetadata ensures a path metadata entry exists for the given path.
func (p *filePackage) ensurePathMetadata(path string, fileEntry *metadata.FileEntry) error {
	// Initialize PathMetadataEntries if needed
	if p.PathMetadataEntries == nil {
		p.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	}

	// Check if path metadata already exists
	for _, entry := range p.PathMetadataEntries {
		if entry.GetPath() == path {
			// Already exists - ensure fileEntry is associated
			alreadyAssociated := false
			for _, assoc := range entry.AssociatedFileEntries {
				if assoc == fileEntry {
					alreadyAssociated = true
					break
				}
			}
			if !alreadyAssociated {
				entry.AssociatedFileEntries = append(entry.AssociatedFileEntries, fileEntry)
			}
			return nil
		}
	}

	// Create new path metadata entry
	pathEntry := &metadata.PathMetadataEntry{
		Path:                  generics.PathEntry{PathLength: uint16(len(path)), Path: path},
		Type:                  metadata.PathMetadataTypeFile,
		AssociatedFileEntries: []*metadata.FileEntry{fileEntry},
		ParentPath:            nil, // Could be enhanced to find parent
		Properties:            []*generics.Tag[any]{},
		Inheritance:           nil, // nil for files
		Metadata:              nil, // nil for files
	}

	p.PathMetadataEntries = append(p.PathMetadataEntries, pathEntry)
	return nil
}

// determineStoredPath determines the stored package path from the filesystem path.
// Implements the complete path determination logic per api_file_mgmt_addition.md Section 2.6 (Path Determination Rules).
//
//nolint:gocognit,gocyclo // path-determination branches
func (p *filePackage) determineStoredPath(filesystemPath string, options *AddFileOptions) (string, error) {
	// Validate that at most one path determination option is set
	optionsSet := 0
	if options != nil {
		if options.StoredPath.IsSet() {
			optionsSet++
		}
		if options.BasePath.IsSet() {
			optionsSet++
		}
		if options.PreserveDepth.IsSet() {
			optionsSet++
		}
		if options.FlattenPaths.IsSet() && options.FlattenPaths.GetOrDefault(false) {
			optionsSet++
		}
	}
	if optionsSet > 1 {
		return "", pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"at most one path determination option may be set (StoredPath, BasePath, PreserveDepth, or FlattenPaths)",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "options",
				Value:    "multiple path options",
				Expected: "at most one of StoredPath, BasePath, PreserveDepth, or FlattenPaths",
			},
		)
	}

	// Priority 1: Explicit StoredPath option
	if options != nil && options.StoredPath.IsSet() {
		explicitPath := options.StoredPath.GetOrDefault("")
		normalizedPath, err := internal.NormalizePackagePath(explicitPath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to normalize explicit stored path",
				pkgerrors.ValidationErrorContext{
					Field:    "StoredPath",
					Value:    explicitPath,
					Expected: "valid normalizable path",
				},
			)
		}
		return normalizedPath, nil
	}

	// Priority 4: FlattenPaths option (before relative path handling)
	if options != nil && options.FlattenPaths.IsSet() && options.FlattenPaths.GetOrDefault(false) {
		// Use just the filename
		baseName := filepath.Base(filesystemPath)
		normalizedPath, err := internal.NormalizePackagePath(baseName)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to normalize flattened path",
				pkgerrors.ValidationErrorContext{
					Field:    "FlattenPaths",
					Value:    baseName,
					Expected: "valid normalizable filename",
				},
			)
		}
		return normalizedPath, nil
	}

	// Handle relative paths (not absolute filesystem paths)
	if !filepath.IsAbs(filesystemPath) {
		// Relative path - treat as package path
		normalizedPath, err := internal.NormalizePackagePath(filesystemPath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to normalize relative path",
				pkgerrors.ValidationErrorContext{
					Field:    "filesystemPath",
					Value:    filesystemPath,
					Expected: "valid normalizable path",
				},
			)
		}
		return normalizedPath, nil
	}

	// Absolute path handling
	// Priority 2: Explicit BasePath option
	if options != nil && options.BasePath.IsSet() {
		basePath := options.BasePath.GetOrDefault("")
		relativePath, err := p.stripBasePath(filesystemPath, basePath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to strip base path from filesystem path",
				pkgerrors.ValidationErrorContext{
					Field:    "BasePath",
					Value:    basePath,
					Expected: "base path that is a prefix of filesystemPath",
				},
			)
		}
		normalizedPath, err := internal.NormalizePackagePath(relativePath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to normalize derived path",
				pkgerrors.ValidationErrorContext{
					Field:    "derivedPath",
					Value:    relativePath,
					Expected: "valid normalizable path",
				},
			)
		}
		return normalizedPath, nil
	}

	// Priority 3: PreserveDepth option
	if options != nil && options.PreserveDepth.IsSet() {
		preserveDepth := options.PreserveDepth.GetOrDefault(1)
		basePath := p.deriveSessionBase(filesystemPath, preserveDepth)
		relativePath, err := p.stripBasePath(filesystemPath, basePath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to derive path with PreserveDepth",
				pkgerrors.ValidationErrorContext{
					Field:    "PreserveDepth",
					Value:    preserveDepth,
					Expected: "valid depth value",
				},
			)
		}
		normalizedPath, err := internal.NormalizePackagePath(relativePath)
		if err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to normalize derived path",
				pkgerrors.ValidationErrorContext{
					Field:    "derivedPath",
					Value:    relativePath,
					Expected: "valid normalizable path",
				},
			)
		}
		return normalizedPath, nil
	}

	// Priority 5: Auto-detection with session base
	sessionBase := p.GetSessionBase()
	if sessionBase == "" {
		// Establish new session base from this path (PreserveDepth=1 by default)
		newBase := p.deriveSessionBase(filesystemPath, 1)
		if err := p.SetSessionBase(newBase); err != nil {
			return "", pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"failed to establish session base",
				pkgerrors.ValidationErrorContext{
					Field:    "filesystemPath",
					Value:    filesystemPath,
					Expected: "valid absolute path for session base establishment",
				},
			)
		}
		sessionBase = newBase
	}

	// Strip session base from path
	relativePath, err := p.stripBasePath(filesystemPath, sessionBase)
	if err != nil {
		return "", pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"filesystem path is not under established session base",
			pkgerrors.ValidationErrorContext{
				Field:    "filesystemPath",
				Value:    filesystemPath,
				Expected: "path under session base: " + sessionBase,
			},
		)
	}

	normalizedPath, err := internal.NormalizePackagePath(relativePath)
	if err != nil {
		return "", pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeValidation,
			"failed to normalize derived path",
			pkgerrors.ValidationErrorContext{
				Field:    "derivedPath",
				Value:    relativePath,
				Expected: "valid normalizable path",
			},
		)
	}

	return normalizedPath, nil
}

// captureFilesystemMetadata captures filesystem metadata and stores it in the path metadata entry.
func (p *filePackage) captureFilesystemMetadata(storedPath string, fileInfo os.FileInfo, fileEntry *metadata.FileEntry, options *AddFileOptions) error {
	// Find or create path metadata entry
	var pathMetadata *metadata.PathMetadataEntry
	for _, entry := range p.PathMetadataEntries {
		if entry.GetPath() == storedPath {
			pathMetadata = entry
			break
		}
	}

	if pathMetadata == nil {
		// This shouldn't happen as ensurePathMetadata is called after this
		// But we'll handle it gracefully
		return nil
	}

	// Always capture IsExecutable (required)
	pathMetadata.FileSystem.IsExecutable = (fileInfo.Mode() & 0o111) != 0

	// Capture additional metadata if requested
	preservePermissions := false
	if options != nil && options.PreservePermissions.IsSet() {
		preservePermissions = options.PreservePermissions.GetOrDefault(false)
	}

	if preservePermissions {
		// Capture full permission bits
		mode := uint32(fileInfo.Mode())
		pathMetadata.FileSystem.Mode = &mode
		pathMetadata.FileSystem.ModTime = uint64(fileInfo.ModTime().Unix())

		// Capture ownership if requested (Unix-specific)
		preserveOwnership := false
		if options != nil && options.PreserveOwnership.IsSet() {
			preserveOwnership = options.PreserveOwnership.GetOrDefault(false)
		}
		if preserveOwnership {
			// FIXME: syscall.Stat_t is Unix-only; prevents cross-compiling for Windows (e.g. nvpkg build-windows-amd64). Use build tags or a Windows-specific path.
			if sys, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
				uid := uint32(sys.Uid)
				gid := uint32(sys.Gid)
				pathMetadata.FileSystem.UID = &uid
				pathMetadata.FileSystem.GID = &gid
			}
		}

		// Additional metadata (ACLs, extended attrs) can be added in Priority 8
	}

	return nil
}

// stripBasePath strips the base path from a filesystem path and returns the relative portion.
func (p *filePackage) stripBasePath(filesystemPath, basePath string) (string, error) {
	cleanFilePath := filepath.Clean(filesystemPath)
	cleanBasePath := filepath.Clean(basePath)

	relPath, err := filepath.Rel(cleanBasePath, cleanFilePath)
	if err != nil {
		return "", err
	}

	// If the relative path starts with "..", it's not under basePath
	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("path %s is not under base path %s", filesystemPath, basePath)
	}

	return relPath, nil
}

// deriveSessionBase derives a session base path from a filesystem path using PreserveDepth.
// PreserveDepth of 1 means preserve one parent directory.
func (p *filePackage) deriveSessionBase(filesystemPath string, preserveDepth int) string {
	cleanPath := filepath.Clean(filesystemPath)
	dir := filepath.Dir(cleanPath)

	if preserveDepth == -1 {
		// Preserve all segments - return root
		return filepath.VolumeName(cleanPath) + string(filepath.Separator)
	}

	if preserveDepth <= 0 {
		// No depth preservation - return the file's directory
		return dir
	}

	// Walk up preserveDepth directories
	basePath := dir
	for i := 0; i < preserveDepth && basePath != filepath.Dir(basePath); i++ {
		basePath = filepath.Dir(basePath)
	}

	return basePath
}

// validateStubContextAndNonEmpty checks context and that value is non-empty; used by stub methods.
func (p *filePackage) validateStubContextAndNonEmpty(ctx context.Context, opName, value, emptyErrMsg, expectedMsg, fieldName string) error {
	if err := internal.CheckContext(ctx, opName); err != nil {
		return err
	}
	if value == "" {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			emptyErrMsg,
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    fieldName,
				Value:    value,
				Expected: expectedMsg,
			},
		)
	}
	return nil
}

// RemoveFilePattern removes files matching a pattern from the package.
//
// STUB IMPLEMENTATION: This method validates inputs but returns ErrTypeUnsupported.
// Full implementation is deferred to Priority 2.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - pattern: Glob pattern to match files
//
// Returns:
//   - []string: Nil slice (stub implementation)
//   - error: *PackageError with ErrTypeUnsupported
//
// Specification: api_file_mgmt_removal.md: 3. RemoveFilePattern Package Method
func (p *filePackage) RemoveFilePattern(ctx context.Context, pattern string) ([]string, error) {
	if err := p.validateStubContextAndNonEmpty(ctx, "RemoveFilePattern", pattern, "pattern cannot be empty", "non-empty glob pattern", "pattern"); err != nil {
		return nil, err
	}
	return nil, p.returnStubUnsupported("RemoveFilePattern full implementation deferred to Priority 2")
}

// returnStubUnsupported returns ErrTypeUnsupported for stub methods.
func (p *filePackage) returnStubUnsupported(msg string) error {
	return pkgerrors.NewPackageError[struct{}](pkgerrors.ErrTypeUnsupported, msg, nil, struct{}{})
}

// RemoveDirectory removes files from a directory path in the package.
//
// STUB IMPLEMENTATION: This method validates inputs but returns ErrTypeUnsupported.
// Full implementation is deferred to Priority 2.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - dirPath: Package directory path to remove files from
//
// Returns:
//   - []string: Nil slice (stub implementation)
//   - error: *PackageError with ErrTypeUnsupported
//
// Specification: api_file_mgmt_removal.md: 4. RemoveDirectory Package Method
func (p *filePackage) RemoveDirectory(ctx context.Context, dirPath string, options *RemoveDirectoryOptions) ([]string, error) {
	_ = options
	if err := p.validateStubContextAndNonEmpty(ctx, "RemoveDirectory", dirPath, "directory path cannot be empty", "non-empty directory path", "dirPath"); err != nil {
		return nil, err
	}
	return nil, p.returnStubUnsupported("RemoveDirectory full implementation deferred to Priority 2")
}
