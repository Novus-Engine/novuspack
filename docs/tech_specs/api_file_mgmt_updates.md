# NovusPack Technical Specifications - File Update API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. UpdateFile Operations](#1-updatefile-operations)
  - [1.1 Package.UpdateFile Method](#11-packageupdatefile-method)
    - [1.1.1 UpdateFile Purpose](#111-updatefile-purpose)
    - [1.1.2 UpdateFile Parameters](#112-updatefile-parameters)
    - [1.1.3 UpdateFile Returns](#113-updatefile-returns)
    - [1.1.4 UpdateFile Behavior](#114-updatefile-behavior)
    - [1.1.5 UpdateFile FileEntry Field Effects](#115-updatefile-fileentry-field-effects)
    - [1.1.6 File Data Processing Model (UpdateFile)](#116-file-data-processing-model-updatefile)
    - [1.1.7 In-Memory Package State Effects (UpdateFile)](#117-in-memory-package-state-effects-updatefile)
    - [1.1.8 UpdateFile Error Conditions](#118-updatefile-error-conditions)
    - [1.1.9 UpdateFile Usage Notes](#119-updatefile-usage-notes)
  - [1.2 Package.UpdateFilePattern Method](#12-packageupdatefilepattern-method)
    - [1.2.1 UpdateFilePattern Purpose](#121-updatefilepattern-purpose)
    - [1.2.2 UpdateFilePattern Parameters](#122-updatefilepattern-parameters)
    - [1.2.3 UpdateFilePattern Returns](#123-updatefilepattern-returns)
    - [1.2.4 UpdateFilePattern Behavior](#124-updatefilepattern-behavior)
    - [1.2.5 UpdateFilePattern Error Conditions](#125-updatefilepattern-error-conditions)
    - [1.2.6 UpdateFilePattern Usage Notes](#126-updatefilepattern-usage-notes)
  - [1.3 Package.UpdateFileMetadata Method](#13-packageupdatefilemetadata-method)
    - [1.3.1 UpdateFileMetadata Purpose](#131-updatefilemetadata-purpose)
    - [1.3.2 UpdateFileMetadata Parameters](#132-updatefilemetadata-parameters)
    - [1.3.3 FileMetadataUpdate Structure](#133-filemetadataupdate-structure)
    - [1.3.4 UpdateFileMetadata Returns](#134-updatefilemetadata-returns)
    - [1.3.5 UpdateFileMetadata Behavior](#135-updatefilemetadata-behavior)
    - [1.3.6 UpdateFileMetadata Error Conditions](#136-updatefilemetadata-error-conditions)
    - [1.3.7 UpdateFileMetadata Usage Notes](#137-updatefilemetadata-usage-notes)
  - [1.4 AddFilePath Package Method Package.AddFilePath](#14-packageaddfilepath-method)
    - [1.4.1 AddFilePath Purpose](#141-addfilepath-purpose)
    - [1.4.2 AddFilePath Parameters](#142-addfilepath-parameters)
    - [1.4.3 AddFilePath Returns](#143-addfilepath-returns)
  - [1.5 RemoveFilePath Package Method Package.RemoveFilePath](#15-packageremovefilepath-method)
    - [1.5.1 RemoveFilePath Purpose](#151-removefilepath-purpose)
    - [1.5.2 RemoveFilePath Parameters](#152-removefilepath-parameters)
    - [1.5.3 RemoveFilePath Returns](#153-removefilepath-returns)
  - [1.6 AddFileHash Package Method Package.AddFileHash](#16-packageaddfilehash-method)
    - [1.6.1 AddFileHash Purpose](#161-addfilehash-purpose)
    - [1.6.2 AddFileHash Parameters](#162-addfilehash-parameters)
    - [1.6.3 AddFileHash Returns](#163-addfilehash-returns)
  - [1.7 ConvertPathsToSymlinks Package Method SymlinkConvertOptions](#17-symlinkconvertoptions-struct)
    - [1.7.1 ConvertPathsToSymlinks Methods](#171-convertpathstosymlinks-methods)
    - [1.7.2 ConvertSymlinksToHardLinks Methods](#172-convertsymlinkstohardlinks-methods)
    - [1.7.3 Multi-Path Query Methods](#173-multi-path-query-methods)
    - [1.7.4 ConvertPathsToSymlinks Purpose](#174-convertpathstosymlinks-purpose)
    - [1.7.5 ConvertPathsToSymlinks Parameters](#175-convertpathstosymlinks-parameters)
    - [1.7.6 ConvertPathsToSymlinks Returns](#176-convertpathstosymlinks-returns)
    - [1.7.7 ConvertPathsToSymlinks Behavior](#177-convertpathstosymlinks-behavior)
    - [1.7.8 ConvertPathsToSymlinks Error Conditions](#178-convertpathstosymlinks-error-conditions)
    - [1.7.9 ConvertAllPathsToSymlinks Purpose](#179-convertallpathstosymlinks-purpose)
    - [1.7.10 ConvertAllPathsToSymlinks Behavior](#1710-convertallpathstosymlinks-behavior)
    - [1.7.11 ConvertSymlinksToHardLinks Purpose](#1711-convertsymlinkstohardlinks-purpose)
    - [1.7.12 ConvertSymlinksToHardLinks Behavior](#1712-convertsymlinkstohardlinks-behavior)

---

## 0. Overview

This document specifies file update operations.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [Core Package Interface](api_core.md)
- [Package Metadata API](api_metadata.md)

## 1. UpdateFile Operations

This section describes UpdateFile operations for updating files in packages.

### 1.1 Package.UpdateFile Method

```go
// UpdateFile updates file content and metadata in the package
// The new file data is read from the sourceFilePath on the filesystem.
// The storedPath identifies which file in the package to update.
func (p *Package) UpdateFile(ctx context.Context, storedPath string, sourceFilePath string, options *AddFileOptions) (*FileEntry, error)
```

#### 1.1.1 UpdateFile Purpose

Updates an existing file's content and metadata in the package, returning the updated FileEntry.

#### 1.1.2 UpdateFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `storedPath`: The package-internal path of the file to update (must already exist in package)
- `sourceFilePath`: Filesystem path to read the new file content from
- `options`: Configuration options for file update (can be nil for defaults)

#### 1.1.3 UpdateFile Returns

- `*FileEntry`: The updated FileEntry with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file update

#### 1.1.4 UpdateFile Behavior

UpdateFile updates an existing file in the package by replacing its content with new data from the filesystem.

UpdateFile MUST execute the following operational sequence:

1. Package File Lookup

   - Validate that the package is open and writable
   - Locate the existing FileEntry using the `storedPath` parameter
   - If not found, return `ErrTypeValidation` error

2. Filesystem Validation and Metadata Read

   - Open the `sourceFilePath` file handle for reading
   - Validate that `sourceFilePath` exists and is not a directory
   - Read file metadata: size, permissions, timestamps
   - Set `FileEntry.CurrentSource` to an external `FileSource` that references `sourceFilePath` with offset `0` and size equal to the raw file size
   - Set `FileEntry.OriginalSource` to the same `FileSource` when multi-stage processing is used (see [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md))
   - Calculate `OriginalSize` from file metadata
   - Determine file type based on extension and content analysis (may require reading file header)
   - Determine effective processing options (compression, encryption) from the provided `AddFileOptions`

3. Deduplication Check

   - UpdateFile MUST perform the same deduplication check defined in [2.1.6 FileEntry Field Effects](api_file_mgmt_addition.md#216-fileentry-field-effects), Step 2
   - If `AddFileOptions.AllowDuplicate` is true, skip deduplication
   - Otherwise, check `OriginalSize` first, then conditionally calculate `RawChecksum` if potential matches exist
   - `EncryptionType` must match for deduplication
   - `CompressionType` matching follows the same rules as AddFile
   - If duplicate found: Clean up any newly opened external file handle(s) and return reference to existing FileEntry (update may be skipped or handled per implementation policy)
   - Emit warning if same `RawChecksum` but different `EncryptionType`

4. Conditional Encryption Processing

   - UpdateFile MUST apply the same conditional encryption processing defined in [2.1.6 FileEntry Field Effects](api_file_mgmt_addition.md#216-fileentry-field-effects), Step 3
   - Only runs if encryption is required and not a duplicate
   - Calculates `RawChecksum` if not already done
   - Applies compression then encryption for dual-processing case
   - Writes processed data to a temporary file (or to pipeline stage output for large files)
   - Updates `FileEntry.CurrentSource` to point to the processed temporary file (and updates `FileEntry.TransformPipeline` when multi-stage processing is used)
   - Sets `FileEntry.ProcessingState` to match the data state model (see [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md#12-processingstate-transitions))

5. FileEntry Update

   - UpdateFile MUST close and clean up any previous `FileEntry.CurrentSource` file handle and any temporary files tracked by `FileEntry.TransformPipeline` before applying updates
   - UpdateFile MUST update all fields as specified in [1.1.5 UpdateFile FileEntry Field Effects](#115-updatefile-fileentry-field-effects)

6. Runtime Field Finalization

   - UpdateFile MUST set the same runtime-only fields defined in [2.1.6 FileEntry Field Effects](api_file_mgmt_addition.md#216-fileentry-field-effects), Step 5
   - Verify `FileEntry.CurrentSource` is set and valid
   - Verify `FileEntry.CurrentSource.Offset` is correct (typically `0` for staged content)
   - Set `FileEntry.CurrentSource.Size` to the length of the file data in its current state
   - Ensure `FileEntry.CurrentSource.FilePath` and `FileEntry.CurrentSource.IsTempFile` match the chosen staging mechanism
   - Set `ProcessingState` to indicate current state of file data (raw, compressed, encrypted, or both)

#### 1.1.5 UpdateFile FileEntry Field Effects

On success, UpdateFile MUST update the existing FileEntry in the in-memory package state.
The following fields MUST be updated to reflect the new content.

- `OriginalSize` MUST be updated to the size of the new raw input bytes.
- `StoredSize` MUST be updated to the size of the new stored bytes after optional compression and optional encryption.
- `RawChecksum` MUST be updated to the checksum of the new raw input bytes.
- `StoredChecksum` MUST be updated to the checksum of the new stored bytes after optional compression and optional encryption.
- `CompressionType` and `CompressionLevel` MUST be updated based on effective options.
- `EncryptionType` MUST be updated based on effective options.

UpdateFile MUST increment `FileVersion` when file content changes.
UpdateFile MUST increment `MetadataVersion` when file metadata changes as a result of the update.

The stored package path is preserved and cannot be changed via UpdateFile.
To change a file's stored path, use `AddFilePath` and `RemoveFilePath`.

The following runtime-only fields MUST be updated using the same rules as AddFile (see [2.1.6 FileEntry Field Effects](api_file_mgmt_addition.md#216-fileentry-field-effects)):

- `FileEntry.CurrentSource` MUST be updated to reference the new staged data location.
- `FileEntry.CurrentSource.Offset` MUST be updated for the new staged data.
- `FileEntry.CurrentSource.Size` MUST be updated to equal the length of the staged data.
- `FileEntry.CurrentSource.FilePath` MUST be set when the source is file-backed.
- `FileEntry.CurrentSource.IsTempFile` MUST be set to reflect whether the staged data is stored in a temporary file.
- `FileEntry.OriginalSource` and `FileEntry.TransformPipeline` MUST be updated consistently when multi-stage processing is used.

UpdateFile MUST close and clean up any previous `FileEntry.CurrentSource` handle and any temporary files tracked by `FileEntry.TransformPipeline` before setting new runtime-only fields.

#### 1.1.6 File Data Processing Model (UpdateFile)

UpdateFile MUST apply the same file data processing model defined in [2.1.7 File Data Processing Model](api_file_mgmt_addition.md#217-file-data-processing-model) for the updated file content.

#### 1.1.7 In-Memory Package State Effects (UpdateFile)

On success, UpdateFile MUST update the existing FileEntry in the package-level list/index of FileEntry objects in memory.
The updated FileEntry MUST be visible to subsequent in-process package operations without requiring a write to disk.

UpdateFile MUST update PackageInfo to reflect the new in-memory package state.
UpdateFile MUST increment `PackageInfo.PackageDataVersion` for any successful update that changes stored bytes.

#### 1.1.8 UpdateFile Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: File does not exist at the specified `storedPath`
- `ErrTypeValidation`: `sourceFilePath` does not exist or is not accessible
- `ErrTypeValidation`: `sourceFilePath` refers to a directory
- `ErrTypeValidation`: File content exceeds size limits
- `ErrTypeIO`: I/O error during file read or update
- `ErrTypeEncryption`: Unsupported encryption type
- `ErrTypeEncryption`: Failed to encrypt file content
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 1.1.9 UpdateFile Usage Notes

UpdateFile updates an existing package file's content by reading new data from the filesystem.
The `storedPath` parameter identifies which file in the package to update.
The `sourceFilePath` parameter specifies where to read the new content from.
AddFileOptions controls compression and encryption settings for the updated content.

### 1.2 Package.UpdateFilePattern Method

```go
// UpdateFilePattern updates files matching a pattern in the package
func (p *Package) UpdateFilePattern(ctx context.Context, pattern string, sourceDir string, options *AddFileOptions) ([]*FileEntry, error)
```

#### 1.2.1 UpdateFilePattern Purpose

Updates multiple files in the package based on a file system pattern and returns the updated FileEntry objects.

#### 1.2.2 UpdateFilePattern Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: File system pattern (e.g., "\*.txt", "documents/\*\*/\*.pdf")
- `sourceDir`: Base directory to search for updated files
- `options`: Configuration options for file processing (can be nil for defaults)

#### 1.2.3 UpdateFilePattern Returns

- `[]*FileEntry`: Slice of updated FileEntry objects with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file update (if error occurs, some files may have been updated successfully)

#### 1.2.4 UpdateFilePattern Behavior

UpdateFilePattern scans the filesystem for files matching the pattern and updates corresponding files in the package.

1. Pattern Matching

   - Scans file system for files matching the pattern in `sourceDir`
   - Applies pattern-specific filters (exclude patterns, max file size)
   - Tracks discovered directories for metadata updates

2. File Updates

   - Finds corresponding files in the package by matching stored paths
   - Updates each matching file with new content from filesystem
   - Applies deduplication checks for updated content
   - Preserves directory structure if requested
   - Reports progress for large file sets

3. Directory Metadata Updates

   - If `AddFileOptions.PreservePermissions` is true, updates directory metadata for discovered directories
   - For each directory discovered during the pattern scan, updates or creates directory metadata using `AddDirectoryMetadata` or `UpdateDirectoryMetadata` from the Metadata API
   - Updates include:
     - Directory permissions (Mode, UID, GID when `PreserveOwnership` is true)
     - Directory timestamps (modification time, access time)
     - Windows attributes when running on Windows and `PreservePermissions` is true
     - Extended attributes when `PreserveExtendedAttrs` is true
     - ACL data when `PreserveACL` is true
   - Directory metadata updates apply to all directories in the pattern match hierarchy

4. Package State Updates

   - Updates package metadata and file count in memory
   - Increments `PackageInfo.PackageDataVersion`
   - Changes become durable only after `Write`, `SafeWrite`, or `FastWrite` completes successfully

#### 1.2.5 UpdateFilePattern Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Invalid or malformed pattern
- `ErrTypeValidation`: No files match the pattern
- `ErrTypeIO`: I/O error during file operations
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 1.2.6 UpdateFilePattern Usage Notes

UpdateFilePattern updates files matching a pattern, allowing specification of a source directory and AddFileOptions.

**Directory Metadata Updates**: When `AddFileOptions` includes metadata preservation flags (`PreservePermissions`, `PreserveOwnership`, `PreserveACL`, `PreserveExtendedAttrs`), UpdateFilePattern automatically updates directory metadata for all discovered directories using `UpdateDirectoryMetadata` from the Metadata API.

This ensures that directory permissions, timestamps, ownership, and attributes remain synchronized with the filesystem when updating files via patterns.

```go
// Update files via pattern with metadata preservation
options := &AddFileOptions{
    PreservePermissions: Option.Some(true),
    PreserveOwnership:   Option.Some(true),
}
entries, err := pkg.UpdateFilePattern(ctx, "**/*.conf", "/etc/myapp", options)
// This will:
// 1. Update all .conf files in the package with content from /etc/myapp
// 2. Update directory metadata for all discovered directories
// 3. Synchronize permissions and ownership for both files and directories
```

**Cross-Reference**: For directory metadata operations, see [Package Metadata API - Path Metadata System](api_metadata.md).

### 1.3 Package.UpdateFileMetadata Method

```go
// UpdateFileMetadata updates file metadata without changing content
func (p *Package) UpdateFileMetadata(ctx context.Context, entry *FileEntry, metadata *FileMetadataUpdate) (*FileEntry, error)
```

#### 1.3.1 UpdateFileMetadata Purpose

Updates file metadata (tags, attributes, compression settings, etc.) without modifying file content.

#### 1.3.2 UpdateFileMetadata Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file to update
- `metadata`: FileMetadataUpdate structure containing new metadata

#### 1.3.3 FileMetadataUpdate Structure

```go
// FileMetadataUpdate contains metadata updates for a FileEntry.
type FileMetadataUpdate struct {
    // Basic metadata
    Tags            []string        // File tags
    CompressionType uint8           // New compression type
    CompressionLevel uint8          // New compression level
    EncryptionType  uint8           // New encryption type

    // Path management
    AddPaths        []generics.PathEntry     // Additional paths to add
    RemovePaths     []string        // Paths to remove (by path string)
    UpdatePaths     []generics.PathEntry     // Paths to update

    // Hash management
    AddHashes       []HashEntry     // Additional hashes to add
    RemoveHashes    []HashType      // Hash types to remove
    UpdateHashes    []HashEntry     // Hashes to update

    // Optional data
    OptionalData    OptionalData   // Structured optional data
}
```

#### 1.3.4 UpdateFileMetadata Returns

- `*FileEntry`: The updated FileEntry with new metadata
- `error`: Any error that occurred during metadata update

#### 1.3.5 UpdateFileMetadata Behavior

- Updates file metadata fields without changing content
- Increments MetadataVersion field
- Preserves file content and data integrity
- Updates package metadata if needed
- Validates new compression/encryption settings
- Applies new compression/encryption to existing content if settings changed
- Manages multiple paths:
  - Adds new paths with per-path metadata (permissions, ownership, timestamps)
  - Removes specified paths (if removing the last path, the entire FileEntry is removed)
  - Updates existing path metadata
- Manages hash data:
  - Adds new hash entries with specified type and purpose
  - Removes hash entries by type
  - Updates existing hash entries
- Manages optional data:
  - Updates tags data (DataType 0x00)
  - Sets path encoding and flags (DataTypes 0x01-0x02)
  - Updates compression dictionary and solid group IDs (DataTypes 0x03-0x04)
  - Sets file system flags and Windows attributes (DataTypes 0x05-0x06)
  - Updates extended attributes and ACL data (DataTypes 0x07-0x08)

#### 1.3.6 UpdateFileMetadata Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: FileEntry does not exist or is invalid
- `ErrTypeCompression`: Unsupported compression type
- `ErrTypeEncryption`: Unsupported encryption type
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 1.3.7 UpdateFileMetadata Usage Notes

UpdateFileMetadata accepts a FileMetadataUpdate structure to modify tags, attributes, compression settings, and optional data without changing file content.

### 1.4 Package.AddFilePath Method

```go
// AddFilePath adds an additional path to an existing FileEntry
func (p *Package) AddFilePath(entry *FileEntry, path generics.PathEntry) (*FileEntry, error)
```

#### 1.4.1 AddFilePath Purpose

Adds an additional path to an existing FileEntry, enabling multiple paths (aliases) to point to the same content.

This allows the same file content to be accessible via different paths within the package.
The additional path is added to the FileEntry's `Paths` array and `PathCount` is incremented.

#### 1.4.2 AddFilePath Parameters

- `entry`: FileEntry reference to the file
- `path`: generics.PathEntry with path (minimal path structure, no metadata)

#### 1.4.3 AddFilePath Returns

- `*FileEntry`: The updated FileEntry with additional path
- `error`: Any error that occurred during path addition

### 1.5 Package.RemoveFilePath Method

```go
// RemoveFilePath removes a path from an existing FileEntry
func (p *Package) RemoveFilePath(entry *FileEntry, path string) (*FileEntry, error)
```

#### 1.5.1 RemoveFilePath Purpose

Removes a specific path from an existing FileEntry.

**Important**: If the removed path is the last path in the FileEntry (PathCount becomes 0), the entire FileEntry is removed from the package and the function returns `nil` for the FileEntry (with no error).
This ensures that file content is only removed when all paths referencing it have been removed.

#### 1.5.2 RemoveFilePath Parameters

- `entry`: FileEntry reference to the file
- `path`: Path string to remove

#### 1.5.3 RemoveFilePath Returns

- `*FileEntry`: The updated FileEntry with path removed, or `nil` if the FileEntry was removed (when last path is removed)
- `error`: Any error that occurred during path removal

### 1.6 Package.AddFileHash Method

```go
// AddFileHash adds a hash entry to an existing FileEntry
func (p *Package) AddFileHash(entry *FileEntry, hash HashEntry) (*FileEntry, error)
```

#### 1.6.1 AddFileHash Purpose

Adds a hash entry to an existing FileEntry for content verification, deduplication, or integrity checking.

#### 1.6.2 AddFileHash Parameters

- `entry`: FileEntry reference to the file
- `hash`: HashEntry with type, purpose, and data

#### 1.6.3 AddFileHash Returns

- `*FileEntry`: The updated FileEntry with additional hash
- `error`: Any error that occurred during hash addition

### 1.7 SymlinkConvertOptions Struct

```go
// SymlinkConvertOptions configures path-to-symlink conversion.
type SymlinkConvertOptions struct {
    // PrimaryPath explicitly specifies which path should be the primary (canonical) path
    // If set, this takes precedence over PrimaryPathSelector
    // If the path is not in FileEntry.Paths, it will be added as the new primary path
    // and all existing paths will be converted to symlinks pointing to it
    // Default: "" (empty, uses PrimaryPathSelector or lexicographic ordering)
    PrimaryPath Option[string]

    // PrimaryPathSelector chooses which path becomes the primary path
    // Only used if PrimaryPath is empty
    // Default: nil (uses first path lexicographically)
    PrimaryPathSelector Option[func(paths []string) string]

    // PreservePathMetadata indicates whether to preserve per-path metadata
    // Default: true
    PreservePathMetadata Option[bool]

    // SymlinkMetadata provides metadata for created symlinks
    // Default: nil (auto-generated from FileEntry metadata)
    SymlinkMetadata Option[func(sourcePath string, targetPath string) SymlinkMetadata]

    // ValidateTargetExists verifies symlink targets exist before creation
    // Default: true (recommended for security and consistency)
    ValidateTargetExists Option[bool]

    // RejectExternalPaths rejects paths pointing outside package root
    // Default: true (required for security, cannot be disabled)
    RejectExternalPaths Option[bool]
}
```

#### 1.7.1 ConvertPathsToSymlinks Methods

This section describes methods for converting paths to symlinks.

##### 1.7.1.1 Package.ConvertPathsToSymlinks Method

```go
// ConvertPathsToSymlinks converts duplicate paths on a FileEntry to symlinks
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - entry: FileEntry with multiple paths (PathCount > 1)
//   - options: Path-to-symlink conversion options (primary path selection, metadata preservation)
// Returns:
//   - Updated FileEntry with single path
//   - Slice of created SymlinkEntry objects
//   - Error if conversion fails
// Validation:
//   - All paths must be within package root (no external references)
//   - Primary path target must exist as FileEntry or PathMetadataEntry
//   - Symlinks will not be created if they would point outside package root
//   - Returns ErrTypeValidation if paths are invalid or outside package root
//   - Returns ErrTypeNotFound if target does not exist
//   - Returns ErrTypeSecurity if symlink target escapes package root
//   - Returns ErrTypePackageState if package is signed (immutable)
func (p *Package) ConvertPathsToSymlinks(ctx context.Context, entry *FileEntry, options *SymlinkConvertOptions) (*FileEntry, []SymlinkEntry, error)
```

##### 1.7.1.2 Package.ConvertAllPathsToSymlinks Method

```go
// ConvertAllPathsToSymlinks converts all multi-path FileEntry objects to symlinks
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - options: Path-to-symlink conversion options
//   - progressCallback: Optional callback for progress reporting (current, total)
// Returns:
//   - Number of FileEntry objects converted
//   - Number of SymlinkEntry objects created
//   - Error if conversion fails
func (p *Package) ConvertAllPathsToSymlinks(ctx context.Context, options *SymlinkConvertOptions, progressCallback func(current, total int)) (int, int, error)
```

#### 1.7.2 ConvertSymlinksToHardLinks Methods

This section describes methods for converting symlinks to hard links.

##### 1.7.2.1 Package.ConvertSymlinksToHardLinks Method

```go
// ConvertSymlinksToHardLinks converts symlinks back to hard links (reverse operation)
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - symlinkEntry: SymlinkEntry to convert back to hard link
// Returns:
//   - Updated FileEntry with additional path added
//   - Error if conversion fails
// Behavior:
//   - Removes SymlinkEntry from package
//   - Adds symlink source path to target FileEntry.Paths array
//   - Preserves all metadata from symlink
//   - Returns ErrTypeNotFound if symlink or target does not exist
//   - Returns ErrTypeValidation if symlink source path already exists in target FileEntry
func (p *Package) ConvertSymlinksToHardLinks(ctx context.Context, symlinkEntry *SymlinkEntry) (*FileEntry, error)
```

##### 1.7.2.2 Package.ConvertAllSymlinksToHardLinks Method

```go
// ConvertAllSymlinksToHardLinks converts all symlinks pointing to FileEntry objects back to hard links
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - progressCallback: Optional callback for progress reporting (current, total)
// Returns:
//   - Number of SymlinkEntry objects converted
//   - Number of FileEntry objects updated
//   - Error if conversion fails
// Behavior:
//   - Only converts symlinks pointing to FileEntry objects (not directory symlinks)
//   - Adds symlink source paths to target FileEntry.Paths arrays
//   - Removes all converted SymlinkEntry objects
//   - Skips symlinks with invalid targets
func (p *Package) ConvertAllSymlinksToHardLinks(ctx context.Context, progressCallback func(current, total int)) (int, int, error)
```

#### 1.7.3 Multi-Path Query Methods

This section describes methods for querying multi-path entries.

##### 1.7.3.1 Package.GetMultiPathEntries Method

```go
// GetMultiPathEntries returns all FileEntry objects with PathCount > 1
// Parameters:
//   - ctx: Context for cancellation and timeout
// Returns:
//   - Slice of FileEntry objects with multiple paths
//   - Error if query fails
func (p *Package) GetMultiPathEntries(ctx context.Context) ([]*FileEntry, error)
```

##### 1.7.3.2 Package.GetMultiPathCount Method

```go
// GetMultiPathCount returns count of multi-path FileEntry objects
// Parameters:
//   - ctx: Context for cancellation and timeout
// Returns:
//   - Count of FileEntry objects with PathCount > 1
//   - Error if query fails
func (p *Package) GetMultiPathCount(ctx context.Context) (int, error)
```

#### 1.7.4 ConvertPathsToSymlinks Purpose

Converts duplicate path entries on a FileEntry into symlinks, transforming a FileEntry with multiple paths into a FileEntry with a single primary path and SymlinkEntry objects pointing to it.

#### 1.7.5 ConvertPathsToSymlinks Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry with multiple paths (PathCount > 1) OR FileEntry with PathCount >= 1 if adding new PrimaryPath
- `options`: Path-to-symlink conversion options (can be nil for defaults)

#### 1.7.6 ConvertPathsToSymlinks Returns

- `*FileEntry`: Updated FileEntry with single primary path (PathCount = 1)
- `[]SymlinkEntry`: Slice of created SymlinkEntry objects (one for each converted path)
- `error`: Any error that occurred during conversion

#### 1.7.7 ConvertPathsToSymlinks Behavior

1. **Pre-Conversion Validation**:
   - Verify `FileEntry.PathCount > 1` OR (`PathCount >= 1` if adding new PrimaryPath)
   - Verify package is not signed (`SignatureOffset == 0` in header)
   - Verify no path conflicts exist (no existing symlinks with same SourcePath)
   - Verify all paths are valid and normalized
   - Verify all paths are within package root (no ".." escapes, no absolute external paths)
   - If `PrimaryPath` is specified, validate it is a valid package-relative path
   - If `PrimaryPath` is specified and NOT in `FileEntry.Paths`, verify it doesn't conflict with existing files
   - Verify no symlinks will be created pointing outside package root

2. **Primary Path Selection** (if needed):
   - If `PrimaryPath` is set and not in `FileEntry.Paths`, add it as new primary path
   - If `PrimaryPath` is set and in `FileEntry.Paths`, use it as primary
   - If `PrimaryPathSelector` is provided, use it to select primary path
   - Otherwise, use first path lexicographically

3. **Symlink Creation**:
   - For each non-primary path, create SymlinkEntry pointing to primary path
   - Create corresponding PathMetadataEntry with `Type: PathMetadataTypeFileSymlink`
   - Preserve path metadata if `PreservePathMetadata` is true

4. **FileEntry Update**:
   - Remove non-primary paths from `FileEntry.Paths` array
   - Update `FileEntry.PathCount` to 1
   - Increment `FileEntry.MetadataVersion` (metadata changed)
   - Preserve all other FileEntry fields (content unchanged)

#### 1.7.8 ConvertPathsToSymlinks Error Conditions

- `ErrTypeValidation`: Invalid FileEntry (PathCount < 1 when adding new path, nil entry, invalid paths, paths outside package root, invalid PrimaryPath format)
- `ErrTypePackageState`: Package is signed (immutable)
- `ErrTypeConflict`: Path conflicts with existing symlinks or files
- `ErrTypeNotFound`: Symlink target does not exist (no FileEntry or PathMetadataEntry for target)
- `ErrTypeIO`: Failure writing symlink metadata files
- `ErrTypeContext`: Context cancellation or timeout
- `ErrTypeSecurity`: Symlink target points outside package root (security violation)

#### 1.7.9 ConvertAllPathsToSymlinks Purpose

Performs batch conversion of all multi-path FileEntry objects to symlinks, useful for optimizing existing packages.

#### 1.7.10 ConvertAllPathsToSymlinks Behavior

- Collects all FileEntry objects with PathCount > 1
- Processes in batches to allow progress reporting
- Calls `ConvertPathsToSymlinks` for each entry
- Returns total counts of converted entries and created symlinks

#### 1.7.11 ConvertSymlinksToHardLinks Purpose

Converts symlinks back to hard links, providing reverse operation for flexibility in package structure.

#### 1.7.12 ConvertSymlinksToHardLinks Behavior

- Removes SymlinkEntry from package
- Adds symlink source path to target FileEntry.Paths array
- Preserves all metadata from symlink
- Updates FileEntry.PathCount and MetadataVersion
