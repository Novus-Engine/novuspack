# NovusPack Technical Specifications - File Addition API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. General Notes](#1-general-notes)
- [2. AddFile Operations](#2-addfile-operations)
  - [2.1 AddFile Package Method](#21-packageaddfile-method)
  - [2.2 AddFileFromMemory Package Method](#22-packageaddfilefrommemory-method)
  - [2.3 AddFileWithEncryption Package Method](#23-packageaddfilewithencryption-method)
  - [2.4 AddFilePattern Package Method](#24-packageaddfilepattern-method)
  - [2.5 AddDirectory Package Method](#25-packageadddirectory-method)
  - [2.6 AddFileOptions: Path Determination](#26-addfileoptions-path-determination)
  - [2.7 PathHandling Type](#27-pathhandling-type)
  - [2.8 AddFileOptions Configuration](#28-addfileoptions-struct)
  - [2.9 Usage Notes](#29-usage-notes)
  - [2.10 Path Normalization and Validation](#210-path-normalization-and-validation)
  - [2.11 Multi-Stage Transformation Pipelines](#211-multi-stage-transformation-pipelines)
- [3. File Addition Implementation Flow](#3-file-addition-implementation-flow)
  - [3.1 Processing Order Requirements](#31-processing-order-requirements)

---

## 0. Overview

This document specifies file addition operations and related options.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [Generic Types and Patterns](api_generics.md)
- [Security Validation API](api_security.md)
- [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md)

## 1. General Notes

This document preserves the original section numbering for AddFile operations (2.x) and implementation flow (3.x).
The multi-stage pipeline system is specified in [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md).

## 2. AddFile Operations

This section describes AddFile operations for adding files to packages.

### 2.1 Package.AddFile Method

```go
// AddFile adds a file to the package.
func (p *Package) AddFile(ctx context.Context, path string, options *AddFileOptions) (*FileEntry, error)
```

#### 2.1.1 AddFile Purpose

Adds a file to the package by reading file data from the filesystem path.
Returns the created FileEntry.

#### 2.1.2 AddFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Filesystem-style input path used to determine the stored package path
- `options`: Optional configuration for file processing (can be nil for defaults)

#### 2.1.3. Filesystem Input Path and Stored Path Derivation

`AddFile` accepts a filesystem-style input path.
It MUST derive a stored package path from that input, unless `AddFileOptions.StoredPath` is set.

The derived stored path MUST follow the package path rules.
See [Package Path Semantics](api_core.md#2-package-path-semantics).

The derived stored path MUST be normalized to NFC before storage to ensure consistent lookups across platforms.

If `path` is a relative path, it is treated as a package-relative path input.
The derived stored path MUST be the normalized stored form of the relative path (leading `/` added, separators normalized, dot segments canonicalized).

If `path` is an absolute path, it is treated as an OS filesystem absolute path input.
The derived stored path MUST be computed using the path determination options described in [2.6 AddFileOptions: Path Determination](#26-addfileoptions-path-determination).
If no explicit `BasePath` is provided, a package-level session base is automatically established (see [Session Base Management](api_basic_operations.md)).

If `path` refers to a directory on disk, `AddFile` MUST return a validation error.
Directory recursion is handled by `AddDirectory` or `AddFilePattern`.

#### 2.1.4 AddFile Returns

- `*FileEntry`: The created FileEntry with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file addition

#### 2.1.5 AddFile Behavior

- Validates the filesystem path exists and is accessible.
- Reads file metadata (size, timestamps).
- Captures execute permission status:
  - Sets `PathFileSystem.IsExecutable` to true if any execute bit is set (user, group, or other)
  - Always captured regardless of `PreservePermissions` setting
- Reads full permission bits if `PreservePermissions` is enabled.
- Determines file type:
  - If `AddFileOptions.FileType` is set, uses the specified type
  - Otherwise, automatically determines type based on extension and content analysis
- Applies early processing based on encryption and compression settings (see [2.1.7 File Data Processing Model](#217-file-data-processing-model)):
  - For encryption (with or without compression): Reads and processes file data, writes to temp file.
    Note: Compression is applied BEFORE encryption when both are selected.
  - For compression-only: Does NOT process file data (deferred to Write operations).
  - For unprocessed files: Does NOT read file data (deferred to Write operations).
- Creates a new FileEntry in the package with complete metadata.
- Creates or updates the PathMetadataEntry for the derived stored path.
- Applies `AddFileOptions.PathMetadataPatch` to the PathMetadataEntry when provided.
- Tracks file data location in runtime fields for subsequent Write operations.

#### 2.1.6 FileEntry Field Effects

On success, AddFile MUST create (or update, in the deduplication case) a FileEntry in the in-memory package state.

AddFile MUST execute operations in the following sequence:

1. **Filesystem Validation and Metadata Read** (first):

   - Validate the filesystem path exists and is accessible
   - **Open source file**: Open a file handle for the filesystem file and stage it as `FileEntry.CurrentSource` using an external `FileSource` with offset `0`
   - Read file metadata from filesystem (size, timestamps if requested)
   - Determine file type based on extension and content analysis (may require reading file header)
   - Calculate `OriginalSize` from the raw file size
   - **Set initial offset**: Set `FileEntry.CurrentSource.Offset` to `0` (start of staged data)
   - Derive stored package path from input path and options
   - Determine effective `CompressionType`, `CompressionLevel`, and `EncryptionType` from options

2. **Deduplication Check** (before any processing):

   - If `AddFileOptions.AllowDuplicate` is true: Skip deduplication, treat as unique, continue to step 3
   - Otherwise, search for existing FileEntry with matching raw content:
     - Match on `OriginalSize` (fast size-based filter)
     - If potential matches found:
       - Read file content and calculate `RawChecksum`
       - Compare `RawChecksum` with candidate FileEntries (matching raw content)
       - Check `EncryptionType` matching:
         - If `EncryptionType` matches: Potential duplicate
         - If `EncryptionType` differs: Not a duplicate, but emit warning that same raw content exists with different encryption (potential unintended duplication)
       - Compression matching rules (for same `EncryptionType` only):
         - If current file has `CompressionType` = None: Match against any `CompressionType` (deduplicate with compressed versions)
         - If current file has `CompressionType` set: Only match same `CompressionType` (allow multiple compression methods for benchmarking)
       - Use SHA-256 hash for final verification if `RawChecksum` and processing options match
   - If duplicate found:
     - Add derived stored package path to existing `Paths` array
     - Update `PathCount` to reflect new path count
     - Increment `MetadataVersion`
     - **Note**: `FileEntry.CurrentSource` remains unchanged because the duplicate shares the existing entry's staged data location
     - Skip to step 5 (runtime field finalization)
   - If unique or no potential matches: Continue to step 3

3. **Conditional Encryption Processing** (only if encryption is required and not a duplicate):

   - Use `RawChecksum` calculated in step 2 (if available)
   - If `RawChecksum` not yet calculated: Read file content and calculate it
   - Apply processing (compress first if both compression and encryption, then encrypt)
   - Calculate `StoredSize` and `StoredChecksum` from processed output
   - Write processed data to a staged output (single temp file for small files, or stage output via `FileEntry.TransformPipeline` for large files)
   - **Update source tracking**: Set `FileEntry.OriginalSource` to the pre-processing source (external file), initialize or extend `FileEntry.TransformPipeline` with the applied stages, and advance `FileEntry.CurrentSource` to the latest stage output (offset `0`) while setting `FileEntry.ProcessingState` accordingly

4. **FileEntry Allocation** (for unique files only):

   - Allocate new FileEntry
   - Assign unique `FileID`
   - Set `Paths` to array containing the derived stored package path
   - Set `PathCount` to 1
   - Set `OriginalSize` from step 1
   - Set `RawChecksum` from step 2 or 3 (may be zero if no deduplication check was performed)
   - For encryption cases: Set `StoredSize` and `StoredChecksum` from step 3
   - For non-encryption cases: Set `StoredSize` and `StoredChecksum` to zero (placeholders, will be calculated during Write)
   - Set `CompressionType`, `CompressionLevel`, `EncryptionType` from effective options
   - Set `Type` from file type detection
   - Set `FileVersion` to initial value
   - Set `MetadataVersion` to initial value

5. **Runtime Field Finalization** (for all cases):

   - Verify `FileEntry.CurrentSource` is set and valid (set in step 1, possibly updated in step 3)
   - Verify `FileEntry.CurrentSource.Offset` is correct (`0` for newly staged data)
   - Set `FileEntry.CurrentSource.Size` to the length of the staged file data in its current state
   - Set `FileEntry.CurrentSource.IsTempFile` to `true` if staged data is stored in a temporary file, or `false` if staged directly from the original filesystem file
   - Set `FileEntry.ProcessingState` to track the current state of data in `FileEntry.CurrentSource` (raw, encrypted, or compressed+encrypted)

The following fields are NOT set by AddFile and are managed by other operations:

- `Data` MUST NOT be populated by AddFile (use `FileEntry.CurrentSource` and its offset/size for Write operations).
- `IsDataLoaded` MUST be false after AddFile completes.
- `EntryOffset` is set during Write operations when the final package file offset is known.
- `PathMetadataEntries` is populated only when filesystem metadata capture is enabled.

#### 2.1.7 File Data Processing Model

AddFile applies different processing based on the requested compression and encryption settings.
The goal is to apply encryption as early as possible for security, while deferring compression to Write operations for flexibility.

The following rules define how AddFile processes file data:

- **Encryption-only (no compression)**:
  AddFile MUST encrypt the file data and write to a temporary file.
  `FileEntry.CurrentSource` MUST point to the temporary file containing encrypted data.
  `FileEntry.CurrentSource.IsTempFile` MUST be true.
  Write operations read encrypted data directly from temp file.

- **Compression-only (no encryption)**:
  AddFile MUST NOT compress the file data.
  `FileEntry.CurrentSource` MUST point to the original source file.
  `FileEntry.CurrentSource.IsTempFile` MUST be false.
  Write operations MUST compress the data when reading from source file.

- **Both compression and encryption**:
  AddFile MUST compress first, then encrypt, writing to a temporary file.
  Compression must happen first because encryption should be applied to the smallest possible data.
  `FileEntry.CurrentSource` MUST point to the temporary file containing compressed and encrypted data.
  `FileEntry.CurrentSource.IsTempFile` MUST be true.
  Write operations read already-processed data directly from temp file.

- **No compression or encryption**:
  AddFile performs no processing.
  `FileEntry.CurrentSource` MUST point to the original source file.
  `FileEntry.CurrentSource.IsTempFile` MUST be false.
  Write operations read raw data directly from source file.

The `ProcessingState` field MUST track the current state of the data in `FileEntry.CurrentSource` to inform Write operations what additional processing is needed:

- `ProcessingStateRaw`: Data in `FileEntry.CurrentSource` is raw (unprocessed)
- `ProcessingStateCompressed`: Data in `FileEntry.CurrentSource` is compressed but not encrypted
- `ProcessingStateEncrypted`: Data in `FileEntry.CurrentSource` is encrypted but not compressed
- `ProcessingStateCompressedAndEncrypted`: Data in `FileEntry.CurrentSource` is both compressed and encrypted

For deduplication to work correctly, `RawChecksum` is the primary deduplication key, calculated from raw file content before any processing.

`StoredSize` and `StoredChecksum` are calculated differently based on processing:

- **Encryption cases**: `StoredSize` is size after encryption, `StoredChecksum` is checksum after encryption. These values are used to verify duplicate encrypted content but NOT for primary deduplication (which uses `RawChecksum`).
- **Compression-only cases**: `StoredSize` and `StoredChecksum` cannot be known until Write operations. Set to zero (placeholder) during AddFile.
- **Unprocessed files**: `StoredSize` equals `OriginalSize`, `StoredChecksum` equals `RawChecksum`.

For compression-only cases, `StoredSize` and `StoredChecksum` are set to zero or placeholder values during AddFile and are calculated and updated during Write operations.

When the same raw content (matching `RawChecksum`) is added with different `EncryptionType` values (e.g., once unencrypted, once encrypted), both versions will be stored as separate FileEntries.
AddFile MUST emit a warning when this occurs, as it indicates potential unintended duplication where the same data exists in both encrypted and unencrypted forms within the package.
The warning SHOULD include the file paths and encryption types involved.

Write operations (`Write`, `SafeWrite`, `FastWrite`) are responsible for:

1. Checking `ProcessingState` to determine what processing is needed.
2. If `ProcessingState` is `ProcessingStateRaw` and `CompressionType` is set, applying compression when reading from `FileEntry.CurrentSource`.
3. If `ProcessingState` is `ProcessingStateCompressed` or `ProcessingStateEncrypted` or `ProcessingStateCompressedAndEncrypted`, reading staged data directly from `FileEntry.CurrentSource` without additional processing.
4. Closing and deleting temporary files when `FileEntry.CurrentSource.IsTempFile` is true after successful package write.
5. Closing but NOT deleting source files when `FileEntry.CurrentSource.IsTempFile` is false.

#### 2.1.8 In-Memory Package State Effects

On success, AddFile MUST update the package-level list/index of FileEntry objects in memory.
The created or updated FileEntry MUST be visible to subsequent in-process package operations (for example, ListFiles, GetFileByPath, and Find operations) without requiring a write to disk.

AddFile MUST update PackageInfo to reflect the new in-memory package state.
AddFile MUST increment `PackageInfo.PackageDataVersion` for any successful addition or content change.

#### 2.1.9 Temporary File Management

AddFile MUST NOT delete temporary files before returning.
Temporary files MUST remain open and accessible until Write operations complete.

If AddFile fails after creating a temporary file, AddFile MUST clean up the temporary file before returning the error.

If the Package is closed before Write operations complete, the Package close operation MUST clean up all temporary files tracked by FileEntry objects.

#### 2.1.10 AddFile Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Invalid or malformed file path
- `ErrTypeValidation`: File already exists at the specified path
- `ErrTypeValidation`: File content exceeds size limits
- `ErrTypeIO`: I/O error during file operations
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 2.2 Package.AddFileFromMemory Method

```go
// AddFileFromMemory adds a file to the package from in-memory data
func (p *Package) AddFileFromMemory(ctx context.Context, path string, data []byte, options *AddFileOptions) (*FileEntry, error)
```

#### 2.2.1 AddFileFromMemory Purpose

Adds a file to the package using provided in-memory byte data instead of reading from the filesystem.
Returns the created FileEntry.

#### 2.2.2 AddFileFromMemory Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Package-relative path where the file will be stored
- `data`: Byte slice containing raw (uncompressed, unencrypted) file content
- `options`: Optional configuration for file processing (can be nil for defaults)

**Important**: The `data` parameter is treated as raw file content (before any compression or encryption).
The package will calculate `RawChecksum` and `OriginalSize` from this data for deduplication purposes.
Compression and encryption are applied during processing according to the options.

#### 2.2.3 AddFileFromMemory PathHandling

Unlike `AddFile`, `AddFileFromMemory` MUST treat the `path` parameter as a package-relative path.
The path MUST be normalized to NFC and follow package path rules (see [Package Path Semantics](api_core.md)).

If `path` starts with `/`, it is used as-is (after normalization).
If `path` does not start with `/`, a leading `/` MUST be prepended.

The `AddFileOptions.StoredPath` field is NOT supported for `AddFileFromMemory` and MUST be ignored if set.
The `AddFileOptions.BasePath` field is NOT supported for `AddFileFromMemory` and MUST be ignored if set.

#### 2.2.4 AddFileFromMemory Returns

- `*FileEntry`: The created FileEntry with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file addition

#### 2.2.5 AddFileFromMemory Behavior

- Validates the path follows package path rules
- Uses provided byte data as raw file content (assumed uncompressed, unencrypted)
- Calculates `OriginalSize` from the data length
- Calculates `RawChecksum` from the raw data for deduplication
- Performs deduplication check (see [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects) step 2):
  - Searches for existing FileEntry with matching `RawChecksum` and `OriginalSize`
  - If duplicate found, adds path to existing entry
  - If unique, creates new FileEntry
- Determines file type:
  - If `AddFileOptions.FileType` is set, uses the specified type
  - Otherwise, automatically determines type based on path extension and content analysis
  - The path extension (e.g., `.json` in `/config/settings.json`) is used for type detection
- Applies processing based on encryption and compression settings (see [2.1.7 File Data Processing Model](#217-file-data-processing-model)):
  - For encryption (with or without compression): Processes data, writes to temp file
  - For compression-only: Defers compression to Write operations
  - For unprocessed files: Uses data directly
- Creates a new FileEntry in the package with complete metadata
- Tracks file data location in runtime fields for subsequent Write operations

#### 2.2.6 AddFileFromMemory FileEntry Field Effects

The FileEntry creation process follows the same sequence as `AddFile` (see [2.1.6](#216-fileentry-field-effects)), with these differences:

1. **Data Source and Validation** (instead of filesystem validation):

   - Validate path follows package path rules
   - Calculate `OriginalSize` from `len(*data)`
   - Write raw data to temporary file for consistent processing pipeline
   - Set `FileEntry.CurrentSource` to a temporary-file `FileSource` with offset `0` and size `len(*data)`
   - Determine file type from path extension and content analysis

2. **Deduplication Check** (same as AddFile):

   - Calculate `RawChecksum` from the raw data in `*data`
   - Search for existing FileEntry with matching `OriginalSize` and `RawChecksum`
   - If `AddFileOptions.AllowDuplicate` is true, skip deduplication
   - If duplicate found, add path to existing entry and skip to step 5
   - If unique or deduplication disabled, continue to step 3

3. **Conditional Encryption Processing** (same as AddFile):

   - Applies encryption/compression as needed based on options
   - For encryption cases: compress first (if requested), then encrypt, write to new temp file
   - Updates `FileEntry.CurrentSource` to point to the processed temp file if encryption is used
   - Updates `StoredSize` and `StoredChecksum` for encrypted cases

4. **FileEntry Allocation** (for unique files, same as AddFile):

   - Allocate new FileEntry with unique `FileID`
   - Set `Paths` array with the package path
   - Set `OriginalSize`, `RawChecksum`, `StoredSize`, `StoredChecksum`
   - Set compression and encryption settings from options

5. **Runtime Field Finalization** (same as AddFile):
   - Verify `FileEntry.CurrentSource` is set and valid
   - Ensure `FileEntry.CurrentSource.IsTempFile` is set appropriately
   - Set `ProcessingState` to indicate data state in `FileEntry.CurrentSource`

#### 2.2.7 AddFileFromMemory Data Management

`AddFileFromMemory` treats the provided `data` as raw (uncompressed, unencrypted) content:

- Calculates `RawChecksum` from the raw data for deduplication
- Calculates `OriginalSize` from `len(data)`
- Writes raw data to a temporary file for consistent processing pipeline with `AddFile`
- The temporary file serves as the `FileEntry.CurrentSource` for Write operations
- For encryption cases, the temporary file is replaced with the encrypted output file (same as `AddFile`)
- Compression and encryption are applied according to options during processing
- Temporary files are managed by the package lifecycle and cleaned up during `Close` or after successful `Write`

**Data Ownership**: The `data` slice parameter references the caller's memory.
Since slices are reference types in Go, only the slice header (24 bytes) is copied, not the underlying data.
The caller must ensure the data remains valid until the Write operation completes or the data is processed (for encryption cases).

#### 2.2.8 AddFileFromMemory Error Conditions

- `ErrTypeValidation`: Invalid path format
- `ErrTypeValidation`: Nil or empty data slice
- `ErrTypeValidation`: Path conflict (duplicate path with different content, and `AllowOverwrite` is false)
- `ErrTypeEncryption`: Encryption failed
- `ErrTypeIO`: Failed to write temporary file
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 2.2.9 AddFileFromMemory Usage Notes

Use `AddFileFromMemory` when:

- File content is generated programmatically
- File content is downloaded from network
- File content is already loaded in memory
- Source is not a filesystem file

For files stored on disk, prefer `AddFile` to avoid unnecessary memory allocation.

#### 2.2.10 AddFileFromMemory Data Requirements

The `data` parameter must contain raw file content (uncompressed, unencrypted):

- The package calculates `RawChecksum` from this data for deduplication
- Compression and encryption are applied during processing according to options
- Do not pre-compress or pre-encrypt the data

#### 2.2.11 AddFileFromMemory Memory Management

The `data` parameter is a slice, which is already a reference type in Go:

```go
data := []byte("file content")
entry, err := pkg.AddFileFromMemory(ctx, "/config/settings.json", data, nil)
```

The caller must ensure the data remains valid until:

- The Write operation completes (for unprocessed or compression-only files), OR
- The encryption processing completes (for encrypted files, which copies data to temp file immediately)

#### 2.2.12 AddFileFromMemory Deduplication

AddFileFromMemory performs the same deduplication as `AddFile`:

- Files with identical raw content (matching `RawChecksum`) share storage
- Set `AddFileOptions.AllowDuplicate = true` to disable deduplication

#### 2.2.13 AddFileFromMemory File Type Specification

File type is automatically detected from the path extension (e.g., `.json` in `/config/settings.json`) and content analysis.

To override automatic detection (recommended), explicitly specify the file type via `AddFileOptions.FileType`:

```go
opts := &AddFileOptions{
    FileType: Option.Some(FileTypeJSON), // Explicitly specify type
}
entry, err := pkg.AddFileFromMemory(ctx, "/config/data", jsonData, opts)
```

This is useful when the path has no extension or when you want to override the extension-based detection.

### 2.3 Package.AddFileWithEncryption Method

```go
// AddFileWithEncryption adds a file with encryption enabled
// This is a convenience wrapper around AddFile that configures encryption options
func (p *Package) AddFileWithEncryption(ctx context.Context, path string, key *EncryptionKey, options *AddFileOptions) (*FileEntry, error)
```

#### 2.3.1 AddFileWithEncryption Purpose

Convenience wrapper that simplifies adding encrypted files without manually constructing `AddFileOptions` with encryption settings.

#### 2.3.2 AddFileWithEncryption Implementation

This function merges the provided `options` with `EncryptionKey` set to the provided key, then calls `AddFile`.

If `options` is `nil`, it creates a new `AddFileOptions` with only the encryption key set.

If `options` already has `EncryptionKey` set, the key parameter takes precedence.

#### 2.3.3 AddFileWithEncryption Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Filesystem-style input path (see [2.1.3](#213-filesystem-input-path-and-stored-path-derivation))
- `key`: Encryption key to use (must be valid and not expired)
- `options`: Optional additional configuration (can be nil)

#### 2.3.4 AddFileWithEncryption Returns

- `*FileEntry`: The created FileEntry with encryption enabled
- `error`: Any error that occurred during file addition or encryption

#### 2.3.5 AddFileWithEncryption Equivalent AddFile Call

```go
opts := options
if opts == nil {
    opts = &AddFileOptions{}
}
opts.EncryptionKey = Option.Some(key)
return package.AddFile(ctx, path, opts)
```

#### 2.3.6 AddFileWithEncryption Usage Notes

- The `options` parameter allows additional configuration (compression, deduplication, etc.) while ensuring encryption is enabled.
- The encryption key must be valid and not expired, or the operation will fail with `ErrTypeEncryption`.
- For full control, use `AddFile` directly with custom `AddFileOptions`.
- See [Generic Encryption Patterns](api_security.md) for key handling requirements.

### 2.4 Package.AddFilePattern Method

```go
// AddFilePattern adds files matching a pattern with options
func (p *Package) AddFilePattern(ctx context.Context, pattern string, options *AddFileOptions) ([]*FileEntry, error)
```

#### 2.4.1 AddFilePattern Purpose

Adds multiple files to the package based on a file system pattern and returns the created FileEntry objects.

AddFilePattern is a convenience wrapper that scans the filesystem for files matching a glob pattern, then internally calls `AddFile` for each matched file.

#### 2.4.2 AddFilePattern Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: Filesystem pattern (glob) used to locate files on disk
- `options`: Configuration options for file processing (can be nil for defaults)

#### 2.4.3 Pattern Base Determination

For each matched filesystem path, `AddFilePattern` MUST derive a stored package path.

The pattern base MUST be derived from the non-wildcard prefix of the pattern.
The base MUST be one directory above the directory that contains the first wildcard.

Example:

- Pattern: `/home/user/project/**/*.png`
- Non-wildcard prefix: `/home/user/project/`
- Pattern directory: `/home/user/project`
- Pattern base: `/home/user`
- Matched: `/home/user/project/assets/button.png`
- Stored: `/project/assets/button.png`

If the pattern contains no wildcard characters, the pattern base MUST be the parent directory of the pattern path.

#### 2.4.4 AddFilePattern Returns

- `[]*FileEntry`: Slice of created FileEntry objects with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file addition (if error occurs, some files may have been added successfully)

#### 2.4.5 AddFilePattern Behavior

AddFilePattern scans the filesystem for files matching the pattern and adds each file to the package.

For each matched file, AddFilePattern:

1. Scan file system for files matching the pattern

   - Apply pattern-specific filters from options (exclude patterns, max file size)
   - Match files based on glob pattern syntax
   - Track discovered directories for metadata capture

2. For each matched file, call `AddFile` with the file path and options

   - Each file addition follows the complete operational sequence defined in [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects)
   - Derive stored package path using the pattern base rules in [2.4.3 Pattern Base Determination](#243-pattern-base-determination)
   - Apply encryption, compression, and deduplication as specified in options

3. Capture directory metadata when metadata preservation is enabled

   - If `AddFileOptions.PreservePermissions` is true, capture directory metadata for each discovered directory
   - For each directory discovered during the pattern scan, call `AddDirectoryMetadata` from the Metadata API
   - Directory metadata includes:
     - Directory permissions (Mode, UID, GID when `PreserveOwnership` is true)
     - Directory timestamps (modification time, access time, creation time)
     - Windows attributes when running on Windows and `PreservePermissions` is true
     - Extended attributes when `PreserveExtendedAttrs` is true
     - ACL data when `PreserveACL` is true
   - Directory metadata is stored in special metadata files (see [Package Metadata API - Path Metadata System](api_metadata.md))
   - Metadata preservation applies to all directories in the pattern match hierarchy

4. Collect results
   - Aggregate all successfully created FileEntry objects
   - Preserve directory structure by default for pattern operations
   - Return aggregated results even if some files failed
   - Directory metadata is persisted separately via the Metadata API and does not appear in the returned FileEntry list

#### 2.4.6 FileEntry Field Effects (AddFilePattern)

For each successfully added file, AddFilePattern MUST apply the same FileEntry field effects defined in [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects) via its internal `AddFile` calls.
Each returned FileEntry MUST reflect the effective options (compression, encryption, file type, and tags) for that file.

#### 2.4.7 File Data Processing Model (AddFilePattern)

AddFilePattern MUST apply the same file data processing model defined in [2.1.7 File Data Processing Model](#217-file-data-processing-model) for each matched file via its internal `AddFile` calls.

#### 2.4.8 In-Memory Package State Effects (AddFilePattern)

AddFilePattern MUST update the package-level list/index of FileEntry objects in memory for each successfully added file.
Each successfully added FileEntry MUST be visible to subsequent in-process package operations without requiring a write to disk.

AddFilePattern MUST update PackageInfo to reflect the new in-memory package state.
AddFilePattern MUST increment `PackageInfo.PackageDataVersion` when any file is successfully added or updated.

#### 2.4.9 AddFilePattern Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Invalid or malformed pattern
- `ErrTypeValidation`: No files match the pattern
- `ErrTypeIO`: I/O error during file operations

#### 2.4.10 AddFilePattern Usage Notes

AddFilePattern is a convenience wrapper that scans the filesystem for files matching a glob pattern, then internally calls `AddFile` for each matched file.
Each file addition follows the complete operational sequence defined in [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects).

**Directory Metadata Preservation**: When `AddFileOptions` includes metadata preservation flags (`PreservePermissions`, `PreserveOwnership`, `PreserveACL`, `PreserveExtendedAttrs`), AddFilePattern automatically captures and stores directory metadata for all discovered directories using `AddDirectoryMetadata` from the Metadata API.

This ensures that when adding files via patterns, directory permissions, timestamps, ownership, and attributes are preserved along with the file metadata.

```go
// Add files via pattern with metadata preservation
options := &AddFileOptions{
    PreservePermissions: Option.Some(true),
    PreserveOwnership:   Option.Some(true),
}
entries, err := pkg.AddFilePattern(ctx, "/project/src/**/*.go", options)
// This will:
// 1. Add all .go files under /project/src
// 2. Capture and store metadata for /project, /project/src, and all subdirectories
// 3. Preserve permissions and ownership for both files and directories
```

See [2.8 AddFileOptions Configuration](#28-addfileoptions-struct) for configuration options.

**Cross-Reference**: For directory metadata details and special metadata file format, see [Package Metadata API - Path Metadata System](api_metadata.md).

### 2.5 Package.AddDirectory Method

```go
// AddDirectory recursively adds files from a filesystem directory into the package
func (p *Package) AddDirectory(ctx context.Context, dirPath string, options *AddFileOptions) ([]*FileEntry, error)
```

#### 2.5.1 AddDirectory Purpose

Recursively adds all files from a filesystem directory to the package and returns the created FileEntry objects.

AddDirectory is a convenience wrapper that recursively scans a directory, then internally calls `AddFile` for each discovered file.

#### 2.5.2 AddDirectory Parameters

- `ctx`: Context for cancellation and timeout handling
- `dirPath`: Filesystem directory path to recursively add
- `options`: Configuration options for file processing (can be nil for defaults)

#### 2.5.3 Stored Path Derivation

For each discovered file, `AddDirectory` MUST derive a stored package path following the same rules as AddFile.

The directory base for path derivation is the parent directory of `dirPath`.

Example:

- Input: `dirPath = /home/user/project`
- Directory base: `/home/user`
- Discovered file: `/home/user/project/src/main.go`
- Stored path: `/project/src/main.go`

If `AddFileOptions.StoredPath` is set, it specifies the base directory within the package where files will be stored, preserving the relative directory structure.

Example with StoredPath:

- Input: `dirPath = /home/user/project`, `StoredPath = /myapp`
- Discovered file: `/home/user/project/src/main.go`
- Stored path: `/myapp/src/main.go`

#### 2.5.4 AddDirectory Returns

- `[]*FileEntry`: Slice of created FileEntry objects for all successfully added files
- `error`: Any error that occurred during directory addition (if error occurs, some files may have been added successfully)

#### 2.5.5 AddDirectory Behavior

AddDirectory recursively scans the filesystem directory and adds each file to the package.

For each discovered file, AddDirectory:

1. Recursively scan `dirPath` for files

   - Follow subdirectories recursively by default
   - Apply pattern-specific filters from options (exclude patterns, max file size)
   - Skip non-regular files (symlinks, devices, etc.) unless options specify otherwise
   - Track discovered directories for metadata capture

2. For each discovered file, call `AddFile` with the file path and options

   - Each file addition follows the complete operational sequence defined in [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects)
   - Derive stored package path using directory base rules above
   - Apply encryption, compression, and deduplication as specified in options

3. Capture directory metadata when metadata preservation is enabled

   - If `AddFileOptions.PreservePermissions` is true, capture directory metadata for each discovered directory
   - For each directory discovered during the recursive scan, call `AddDirectoryMetadata` from the Metadata API
   - Directory metadata includes:
     - Directory permissions (Mode, UID, GID when `PreserveOwnership` is true)
     - Directory timestamps (modification time, access time, creation time)
     - Windows attributes when running on Windows and `PreservePermissions` is true
     - Extended attributes when `PreserveExtendedAttrs` is true
     - ACL data when `PreserveACL` is true
   - Directory metadata is stored in special metadata files (see [Package Metadata API - Path Metadata System](api_metadata.md))
   - Metadata preservation applies to both the root `dirPath` and all subdirectories

4. Collect results
   - Aggregate all successfully created FileEntry objects
   - Return aggregated results even if some files failed
   - Directory metadata is persisted separately via the Metadata API and does not appear in the returned FileEntry list

#### 2.5.6 FileEntry Field Effects (AddDirectory)

For each successfully added file, AddDirectory MUST apply the same FileEntry field effects defined in [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects) via its internal `AddFile` calls.

Each returned FileEntry MUST reflect the effective options (compression, encryption, file type, and tags) for that file.

#### 2.5.7 File Data Processing Model (AddDirectory)

AddDirectory MUST apply the same file data processing model defined in [2.1.7 File Data Processing Model](#217-file-data-processing-model) for each added file via its internal `AddFile` calls.

#### 2.5.8 In-Memory Package State Effects (AddDirectory)

AddDirectory MUST update the package-level list/index of FileEntry objects in memory for each successfully added file.

Each successfully added FileEntry MUST be visible to subsequent in-process package operations without requiring a write to disk.

AddDirectory MUST update PackageInfo to reflect the new in-memory package state.
AddDirectory MUST increment `PackageInfo.PackageDataVersion` for each file successfully added or updated.

#### 2.5.9 AddDirectory Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: `dirPath` does not exist
- `ErrTypeValidation`: `dirPath` is not a directory
- `ErrTypeValidation`: Directory is empty and no files were added
- `ErrTypeIO`: I/O error during directory scan or file operations
- `ErrTypeEncryption`: Encryption error for one or more files
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 2.5.10 AddDirectory Usage Notes

AddDirectory is a convenience method that internally uses `AddFile` for each discovered file.
It simplifies adding entire directory trees to a package.

See [2.8 AddFileOptions Configuration](#28-addfileoptions-struct) for configuration options.

**Directory Metadata Preservation**: When `AddFileOptions` includes metadata preservation flags (`PreservePermissions`, `PreserveOwnership`, `PreserveACL`, `PreserveExtendedAttrs`), AddDirectory automatically captures and stores directory metadata for all discovered directories using `AddDirectoryMetadata` from the Metadata API.

This ensures that directory permissions, timestamps, ownership, and attributes are preserved when extracting the package, maintaining the complete directory structure with all metadata intact.

```go
// Add directory with metadata preservation
options := &AddFileOptions{
    PreservePermissions: Option.Some(true),
    PreserveOwnership:   Option.Some(true),
    PreserveACL:         Option.Some(true),
}
entries, err := pkg.AddDirectory(ctx, "/path/to/project", options)
// This will:
// 1. Add all files under /path/to/project
// 2. Capture and store metadata for /path/to/project and all subdirectories
// 3. Preserve permissions, ownership, and ACLs for both files and directories
```

To add only directory metadata without files, use `AddDirectoryMetadata` from the Metadata API directly.

**Cross-Reference**: For directory metadata details and special metadata file format, see [Package Metadata API - Path Metadata System](api_metadata.md).

### 2.6 AddFileOptions: Path Determination

This section defines how the path determination fields in `AddFileOptions` control the mapping from filesystem-style input paths to stored package paths.

These options apply to all file addition operations: AddFile, AddFilePattern, and AddDirectory.

Stored package paths MUST follow the stored path format rules.
See [Package Path Semantics](api_core.md#2-package-path-semantics).

#### 2.6.1 Path Determination Priority Order

The path determination options MUST be applied in the following priority order:

1. `StoredPath`
2. `BasePath`
3. `FlattenPaths`
4. `PreserveDepth`
5. Auto-detection (default behavior when no path determination options are set)

At most one of the following options may be set: `StoredPath`, `BasePath`, `PreserveDepth`, or `FlattenPaths` (set to true).
If more than one is set, `AddFile` and `AddFilePattern` MUST return a validation error.

#### 2.6.2 StoredPath Option

If `StoredPath` is set, it MUST be treated as the exact stored package path for the file.
The input `path` is used only to provide filesystem context for reading (if applicable) and for error context.
`StoredPath` MUST be validated as a stored package path.

#### 2.6.3 BasePath Option

If `BasePath` is set, it MUST be used as the filesystem base directory to strip from absolute input paths.
The derived stored path MUST be the cleaned, separator-normalized path relative to `BasePath`, with a leading `/` added.

#### 2.6.4 Session Base - Package-Level Automatic BasePath

When `BasePath` is not explicitly set in `AddFileOptions` and an absolute input path is provided, a **session base** is automatically established at the package level.

The session base serves as an implicit package-level base path that persists across multiple file addition operations within the same package construction session.

#### 2.6.5 When Session Base Is Established

A session base MUST be established when ALL of the following conditions are met:

1. `AddFileOptions.BasePath` is not explicitly set
2. An absolute input path is provided to AddFile, AddFilePattern, or AddDirectory
3. No session base has been previously established for this package

The session base MUST be derived from the first absolute input path using the PreserveDepth rules (default PreserveDepth: 1).

#### 2.6.6 Session Base Persistence

Once established, the session base persists for all subsequent file addition operations within the same package construction session.

Subsequent absolute paths MUST be validated against the established session base.
If a subsequent absolute input path is not under the established session base, the operation MUST return a validation error.
The error message MUST direct the user to set `BasePath` explicitly (or to use `StoredPath`).

#### 2.6.7 Session Base Lifecycle

The session base is runtime-only construction context.
It MUST NOT be persisted to disk or exposed through package metadata.
It exists only during package construction and is discarded when the package is written or closed.

#### 2.6.8 Overriding Session Base

To override the session base for a specific operation, set `AddFileOptions.BasePath` explicitly.
The per-operation `BasePath` option takes precedence over the package-level session base (see [2.6.1 Path Determination Priority Order](#261-path-determination-priority-order)).

For the package-level session base management API (SetSessionBase, GetSessionBase, ClearSessionBase, HasSessionBase), see [Basic Operations API - Session Base Management](api_basic_operations.md).

#### 2.6.9 PreserveDepth Option

For an absolute input path and an effective PreserveDepth value `d`:

- If `d` is `1`, the derived stored path MUST preserve exactly one parent directory segment.
- If `d` is `2`, the derived stored path MUST preserve exactly two parent directory segments.
- If `d` is `-1`, the derived stored path MUST preserve all path segments under the filesystem root (or volume root).

If `PreserveDepth` is not set and no `BasePath` is set and no session base is established, the default PreserveDepth MUST be treated as `1`.

#### 2.6.10 FlattenPaths Option

If `FlattenPaths` is set to true, the derived stored path MUST be `"/" + baseName(inputPath)`.
This MAY cause filename conflicts.

#### 2.6.11 Path Length Validation

The derived stored path MUST NOT exceed the on-disk format limits.
See [PathEntry](api_generics.md) for the stored path length constraints.

#### 2.6.12 On-Disk Format Limit

The maximum stored path length is **65,535 bytes** (enforced by `PathEntry.PathLength` being `uint16`).

If a derived stored path exceeds 65,535 bytes, the operation MUST return a validation error with `ErrTypeValidation`.

#### 2.6.13 Platform Extraction Limits

Common platform extraction limits vary by operating system.

Unix/Linux Systems:

- Traditional limit: 4,096 bytes (PATH_MAX on most systems)
- Modern systems support longer paths depending on filesystem

Windows Systems:

- Legacy MAX_PATH limit: 260 characters (approximately 260 bytes for ASCII)
- Extended-length path limit: 32,767 characters (with `\\?\` prefix)
- Default Windows extraction uses extended-length path API automatically for paths exceeding MAX_PATH

Warning Behavior:

If a derived stored path exceeds 4,096 bytes, the operation SHOULD emit a warning (non-fatal) indicating potential extraction issues on some platforms.

The warning message SHOULD include:

- The actual path length in bytes
- The stored path value
- Platform compatibility note

Warnings MUST NOT prevent the file from being added or change the stored path.

#### 2.6.14 Windows Extended-Length PathHandling

When extracting files on Windows, the implementation SHOULD automatically use extended-length path syntax (`\\?\`) for paths that would exceed MAX_PATH (260 characters).

This allows Windows extraction to handle paths up to 32,767 characters without user intervention.

See Windows documentation on [Maximum Path Length Limitation](https://learn.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation) for details.

### 2.7 PathHandling Type

See [PathHandling Type](api_basic_operations.md#93-pathhandling-type) for the complete type definition.

The `PathHandling` type is defined in the Basic Operations API and is used here to configure path handling behavior during file addition operations.

#### 2.7.1 PathHandling Purpose

The `PathHandling` type specifies how the system should handle multiple paths pointing to the same file content during deduplication and file addition operations.

#### 2.7.2 PathHandling Values

- `PathHandlingDefault` (0): Use package default (`Package.DefaultPathHandling`)
- `PathHandlingHardLinks` (1): Store multiple paths as hard links (current behavior, backward compatible)
- `PathHandlingSymlinks` (2): Convert additional paths to symlinks
- `PathHandlingPreserve` (3): Preserve original filesystem behavior (detect and respect symlinks/hardlinks)

### 2.8 AddFileOptions Struct

```go
// AddFileOptions configures file addition behavior for both individual files and patterns
type AddFileOptions struct {
    // Path determination options.
    //
    // These options control how filesystem-style input paths are mapped to stored package paths.
    // Exactly one of StoredPath, BasePath, PreserveDepth, or FlattenPaths may be set.
    // See [2.6 AddFileOptions: Path Determination](#26-addfileoptions-path-determination) for complete rules.
    //
    // If BasePath is not set, a package-level session base is automatically established from
    // the first absolute path (see [Session Base Management](api_basic_operations.md)).
    StoredPath     Option[string] // Explicit stored package path override (stored format)
    BasePath       Option[string] // Filesystem base directory to strip from absolute input paths (overrides package-level session base)
    PreserveDepth  Option[int]    // Parent directory depth to preserve for absolute input paths (-1 = preserve all segments)
    FlattenPaths   Option[bool]   // Store at package root using file name only (equivalent to PreserveDepth = 0)

    // Conflict handling.
    AllowOverwrite Option[bool]   // Allow overwrite when a stored path already exists with different content
    AllowDuplicate Option[bool]   // Skip deduplication checks and always create a new FileEntry (default: false)

    // Symlink handling.
    FollowSymlinks Option[bool]   // Follow symlinks by default when reading from filesystem paths

    // Path handling for duplicate content.
    PathHandling PathHandling     // How to handle multiple paths pointing to the same content (default: PathHandlingDefault)
    PrimaryPathSelector func(paths []string) string // Custom selector for primary path when converting to symlinks (default: nil, uses lexicographic ordering)

    // Filesystem metadata capture options.
    //
    // Note: Execute permissions are always captured via PathFileSystem.IsExecutable
    // regardless of PreservePermissions setting.
    PreservePermissions   Option[bool] // Capture full permission bits (Mode, UID, GID)
    PreserveOwnership     Option[bool] // Capture UID and GID (requires PreservePermissions)
    PreserveACL           Option[bool] // Capture ACL entries when available
    PreserveExtendedAttrs Option[bool] // Capture extended attributes when available

    // Path metadata patch.
    //
    // If set, the add operation MUST create (if missing) or update (if present) the PathMetadataEntry
    // for the derived stored path with the specified fields.
    //
    // The PathMetadataEntry Type MUST be consistent with the stored path format:
    // - "/" and paths ending with "/" are directories.
    // - all other paths are files.
    //
    // If filesystem metadata capture is enabled (for example PreservePermissions),
    // captured filesystem metadata MUST be applied first, then PathMetadataPatch MUST be applied.
    // Fields explicitly set in PathMetadataPatch MUST override captured values.
    PathMetadataPatch Option[*PathMetadataPatch]

    // File processing options
    Compress        Option[bool]            // Whether to compress the file
    CompressionType Option[uint8]           // Compression algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    CompressionLevel Option[int]            // Compression level (1-9, 0 = default)
    FileType        Option[uint16]          // File type identifier
    Tags            Option[[]*Tag[any]] // Per-file tags (typed tags)

    // Encryption options
    EncryptionKey   Option[*EncryptionKey]  // Encryption key (presence enables encryption)

    // Multi-stage transformation pipeline options
    MaxTransformStages      Option[int]  // Maximum transformation stages per pipeline (default: 10)
    ValidateProcessingState Option[bool] // Enable ProcessingState validation (default: false)

    // Pattern-specific options (only used for pattern operations)
    ExcludePatterns Option[[]string]        // Patterns to exclude from processing
    MaxFileSize     Option[int64]           // Maximum file size to include (0 = no limit)
    PreservePaths   Option[bool]            // Whether to preserve directory structure for pattern operations (default: true)
}
```

#### 2.8.1 AddFileOptions Purpose

Unified configuration options for all file addition operations, supporting both individual files and pattern-based operations.

#### 2.8.2 PathMetadataPatch Struct

```go
// PathMetadataPatch specifies persisted PathMetadataEntry fields to create or update at add time.
//
// This patch is applied to the PathMetadataEntry for the derived stored path.
// It does not change the stored path itself.
//
// Cross-Reference:
// - [Package Metadata API - PathMetadataEntry Structure](api_metadata.md#812-pathmetadataentry-structure)
type PathMetadataPatch struct {
    // Tags are typed tags stored on PathMetadataEntry.
    Tags  Option[[]*Tag[any]]

    // Inheritance and Metadata are directory-only fields.
    // If applied to a file path, the implementation MUST return ErrTypeValidation.
    Inheritance Option[*PathInheritance]
    Metadata    Option[*PathMetadata]

    // Destination extraction override paths.
    DestPath    Option[string]
    DestPathWin Option[string]

    // Filesystem metadata fields.
    FileSystem  Option[PathFileSystem]
}
```

#### 2.8.3 Conflict Handling Options

If a derived stored path already exists in the package:

- If the content is identical (deduplication), the path MUST be added as an additional path to the existing FileEntry.
- If the content differs and `AllowOverwrite` is not set to true, the operation MUST return a validation error.
- If the content differs and `AllowOverwrite` is set to true, the operation MUST overwrite the existing stored path.

#### 2.8.4 Deduplication Options

- `AllowDuplicate`: If set to true, skip deduplication checks and always create a new FileEntry even if identical content exists in the package (default: false). This option allows intentional storage of duplicate content for performance (skip deduplication overhead) or testing/benchmarking purposes.

#### 2.8.5 File Processing Options

- `Compress`: Whether to compress the file (default: false)
- `CompressionType`: Compression algorithm identifier (default: 0 = none, 1=Zstd, 2=LZ4, 3=LZMA)
- `CompressionLevel`: Compression level 1-9 (default: 0 = default)
- `FileType`: File type identifier (default: 0 = regular file)
  - If set, overrides automatic file type detection
  - Useful when path has no extension or when overriding extension-based detection
  - Applies to both `AddFile` and `AddFileFromMemory`
- `Tags`: Per-file tags as key-value pairs (default: nil)

#### 2.8.6 Encryption Options

- `EncryptionKey`: Encryption key to use for file encryption (if set, file will be encrypted with this key)

#### 2.8.7 Multi-Stage Transformation Pipeline Options

- `MaxTransformStages`: Maximum number of transformation stages allowed per pipeline (default: 10)

  - Limits the depth of multi-stage transformations (e.g., decrypt => decompress => verify)
  - Prevents memory leaks and runaway resource usage
  - Covers typical 3-stage operations with generous headroom
  - Can be increased for advanced use cases requiring more stages
  - If set to 0 or negative, uses the default value

- `ValidateProcessingState`: Enable ProcessingState validation (default: false)

  - When true, validates that ProcessingState matches the actual data state in CurrentSource
  - Returns `ErrTypeValidation` if mismatch detected between ProcessingState and actual transformations applied
  - Useful for debugging and development environments
  - Disabled by default for performance (validation adds overhead)
  - Recommended for testing and strict validation scenarios

#### 2.8.8 Pattern-Specific Options

- `ExcludePatterns`: Patterns to exclude from processing (default: nil)
- `MaxFileSize`: Maximum file size to include, 0 = no limit (default: 0)
- `PreservePaths`: Whether to preserve directory structure for pattern operations (default: true)

#### 2.8.9 Symlink Options

If `FollowSymlinks` is not set, it MUST default to true.

#### 2.8.10 PathHandling Options

- `PathHandling`: Specifies how to handle multiple paths pointing to the same content during deduplication
  - `PathHandlingDefault` (0): Use package default (`Package.DefaultPathHandling`)
  - `PathHandlingHardLinks` (1): Store multiple paths as hard links (current behavior, backward compatible)
  - `PathHandlingSymlinks` (2): Convert additional paths to symlinks
  - `PathHandlingPreserve` (3): Preserve original filesystem behavior (detect and respect symlinks/hardlinks)
  - Default: `PathHandlingDefault` (uses `Package.DefaultPathHandling`)

- `PrimaryPathSelector`: Custom function to select primary path when converting to symlinks
  - Only used when `PathHandling` is `PathHandlingSymlinks` and `PathHandlingDefault` resolves to symlinks
  - Receives all paths and returns the chosen primary path
  - Default: nil (uses first path lexicographically)
  - Allows custom logic (shortest, shallowest, pattern matching, etc.)

If `FollowSymlinks` is true and the input filesystem path is a symlink, the file content MUST be read from the symlink target.
The stored package path MUST still be derived from the symlink path.

If `FollowSymlinks` is false and the input filesystem path is a symlink, `AddFile` MUST return a validation error.
Symlinks MUST be added through the symlink API.
See [Symlink Metadata](api_metadata.md).

#### 2.8.11 Filesystem Metadata Capture Options

Filesystem metadata capture applies only when adding from filesystem paths.

**Directory Metadata**: When using `AddDirectory`, these options also control whether directory metadata is captured and stored.
If `PreservePermissions` is true, AddDirectory automatically captures directory metadata (permissions, timestamps, ownership, etc.) for all discovered directories using `AddDirectoryMetadata` from the Metadata API.
See [2.5.10 AddDirectory Usage Notes](#2510-adddirectory-usage-notes) for complete directory metadata behavior.

#### 2.8.12 Execute Permissions (Always Captured)

The implementation MUST always capture execute permission status into `PathFileSystem.IsExecutable`, regardless of `PreservePermissions` setting:

- On Unix/Linux: Set to `true` if any execute bit is set in the file mode (user, group, or other: `mode & 0111 != 0`)
- On Windows: Set to `true` if the file extension indicates executable (.exe, .bat, .cmd, .ps1, .com, .vbs, .wsf, etc.)
- For `AddFileFromMemory`: Set based on path extension analysis
- This field is NOT optional and MUST always be set during file addition
- Default value is `false` if detection is not possible

#### 2.8.13 Full Permission Bits (Optional)

If `PreservePermissions` is true, the implementation MUST capture the full permission bits into `PathFileSystem.Mode`.

When using `AddDirectory` with `PreservePermissions` set to true, directory permissions are also captured and stored via `AddDirectoryMetadata`.

If `PreserveOwnership` is true, the implementation MUST capture UID and GID into the path metadata system.
`PreserveOwnership` MUST require `PreservePermissions` to be true.

When using `AddDirectory` with `PreserveOwnership` set to true, directory ownership (UID and GID) is also captured and stored.

If `PreserveACL` is true, the implementation MUST capture ACL entries when available into the path metadata system.

When using `AddDirectory` with `PreserveACL` set to true, directory ACL entries are also captured and stored.

If `PreserveExtendedAttrs` is true, the implementation MUST capture extended attributes when available into the path metadata system.

When using `AddDirectory` with `PreserveExtendedAttrs` set to true, directory extended attributes are also captured and stored.

#### 2.8.14 Windows Attributes

When running on Windows, the implementation MUST capture Windows file attributes into `PathFileSystem.WindowsAttrs` when `PreservePermissions` is true.

When using `AddDirectory` on Windows with `PreservePermissions` set to true, directory attributes are also captured and stored via `AddDirectoryMetadata`.

When running on Windows, `PathFileSystem.Mode`, `PathFileSystem.UID`, and `PathFileSystem.GID` are not applicable.
When running on Windows, these fields MUST be left unset.

When running on Windows and `PreserveACL` is true, the implementation MUST capture the Windows security descriptor using an ASCII SDDL string and store it in `PathFileSystem.ExtendedAttrs`.

When using `AddDirectory` on Windows with `PreserveACL` set to true, directory security descriptors are also captured and stored.

When running on Windows and `PreserveOwnership` is true, the implementation MUST capture the owner and group SIDs and store them as ASCII strings in `PathFileSystem.ExtendedAttrs`.

When using `AddDirectory` on Windows with `PreserveOwnership` set to true, directory owner and group SIDs are also captured and stored.

The following keys MUST be used for Windows security descriptor storage:

- `windows.security_descriptor_sddl`
- `windows.owner_sid`
- `windows.group_sid`

If the security descriptor or SIDs cannot be retrieved on Windows, the implementation MUST return a validation error when `PreserveACL` or `PreserveOwnership` is set.

See [Path Metadata System](api_metadata.md) for the storage structures.

### 2.9 Usage Notes

AddFile reads file data from the filesystem path.
AddFileFromMemory adds file data from memory.
Use AddFileOptions to configure compression, encryption, tags, and path determination behavior.

### 2.10 Path Normalization and Validation

All file paths added to the package MUST undergo normalization and validation to ensure cross-platform compatibility and correctness.

For complete path normalization and validation requirements, see:

- [Unicode Normalization](api_core.md#214-unicode-normalization) - NFC normalization requirements
- [Path Length Limits](api_core.md#215-path-length-limits) - Path length handling and portability warnings
- [Case Sensitivity](api_core.md#221-case-sensitivity) - Case-sensitive storage and extraction behavior
- [Path Normalization Rules](api_core.md#21-path-normalization-rules) - Complete normalization requirements
- [ValidatePackagePath Function](api_core.md#123-validatepackagepath-function) - Path validation requirements

These rules apply to all file addition operations (`AddFile`, `AddFileFromMemory`, `AddFilePattern`, `AddDirectory`).

### 2.11 Multi-Stage Transformation Pipelines

For the pipeline system used for memory-efficient large file processing, see [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md).

## 3. File Addition Implementation Flow

This section describes the implementation flow for file addition operations.

### 3.1 Processing Order Requirements

The file addition process must follow a specific sequence to ensure proper encryption, deferred compression, and deduplication.

For the complete operational sequence with all FileEntry field effects and processing details, see [2.1.6 FileEntry Field Effects](#216-fileentry-field-effects).

#### 3.1.1 High-Level Processing Steps

1. **Filesystem Validation and Metadata Read** - Validate filesystem path, open file handle, read file metadata, set initial `FileEntry.CurrentSource` (external `FileSource`, offset `0`)

2. **Deduplication Check** - Check for existing files based on `RawChecksum` (primary key), `EncryptionType` (must match), and `CompressionType` (special rules). Skip if `AllowDuplicate` is true.

3. **Conditional Encryption Processing** - Only applies if encryption is required and file is not a duplicate. Compression + encryption are applied together during `AddFile`, writing to a temporary file. `FileEntry.CurrentSource` is updated to point to the processed temp file.

4. **FileEntry Allocation** - Create new FileEntry with all metadata, checksums, and processing flags

5. **Runtime Field Finalization** - Set `FileEntry.CurrentSource.Size`, `FileEntry.CurrentSource.IsTempFile`, and `ProcessingState`

**Important**: Compression-only files (no encryption) defer all processing to `Write` operations. Their `StoredSize` and `StoredChecksum` are placeholders until Write.

#### 3.1.2 Error Handling Requirements

- **Compression failures**: Must prevent file addition and return appropriate error during Write operations
- **Encryption failures**: Must prevent file addition and return appropriate error during AddFile
- **Resource cleanup**: Failed operations must properly clean up allocated resources (temp files, file handles)
- **User feedback**: Provide clear error messages explaining failures and recovery options

#### 3.1.3 Performance Requirements

- **Deduplication efficiency**: Use `OriginalSize` as first filter, then conditionally calculate `RawChecksum` only when potential matches exist
- **Memory efficiency**: Stream file content for compression/encryption, avoid loading entire files into memory
- **Temp file management**: Clean up temporary files after successful Write or on error
- **SHA-256 optimization**: Only compute expensive SHA-256 when size and checksum match
- **Memory management**: Handle large files efficiently with streaming when needed
- **No wasteful processing**: Do not compress file data only to discard it
