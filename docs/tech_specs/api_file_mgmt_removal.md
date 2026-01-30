# NovusPack Technical Specifications - File Removal API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Removal Semantics And Multi-Path Files](#1-file-removal-semantics-and-multi-path-files)
- [2. RemoveFile Package Method](#2-removefile-package-method)
  - [2.1 RemoveFile Purpose](#21-removefile-purpose)
  - [2.2 RemoveFile Signature](#22-packageremovefile-method)
  - [2.3 RemoveFile Parameters](#23-removefile-parameters)
  - [2.4 RemoveFile Returns](#24-removefile-returns)
  - [2.5 RemoveFile Behavior](#25-removefile-behavior)
  - [2.6 RemoveFile Error conditions](#26-removefile-error-conditions)
  - [2.7 RemoveFile Usage notes](#27-removefile-usage-notes)
- [3. RemoveFilePattern Package Method](#3-removefilepattern-package-method)
  - [3.1 RemoveFilePattern Purpose](#31-removefilepattern-purpose)
  - [3.2 RemoveFilePattern Signature](#32-packageremovefilepattern-method)
  - [3.3 RemoveFilePattern Parameters](#33-removefilepattern-parameters)
  - [3.4 RemoveFilePattern Returns](#34-removefilepattern-returns)
  - [3.5 RemoveFilePattern Behavior](#35-removefilepattern-behavior)
    - [3.5.1 Pattern Matching](#351-pattern-matching)
    - [3.5.2 File Removal](#352-file-removal)
    - [3.5.3 Directory Metadata Cleanup](#353-directory-metadata-cleanup)
    - [3.5.4 Package State Updates](#354-package-state-updates)
  - [3.6 RemoveFilePattern Error conditions](#36-removefilepattern-error-conditions)
    - [3.6.1 Partial Failure Handling](#361-partial-failure-handling)
  - [3.7 RemoveFilePattern Usage notes](#37-removefilepattern-usage-notes)
- [4. RemoveDirectory Package Method](#4-removedirectory-package-method)
  - [4.1 RemoveDirectory Purpose](#41-removedirectory-purpose)
  - [4.2 RemoveDirectory Signature](#42-packageremovedirectory-method)
  - [4.3 RemoveDirectory Parameters](#43-removedirectory-parameters)
  - [4.4 RemoveDirectoryOptions Struct](#44-removedirectoryoptions-struct)
  - [4.5 RemoveDirectory Returns](#45-removedirectory-returns)
  - [4.6 RemoveDirectory Behavior](#46-removedirectory-behavior)
    - [4.6.1 Directory Path Validation](#461-directory-path-validation)
    - [4.6.2 File Discovery](#462-file-discovery)
    - [4.6.3 RemoveDirectory File Removal](#463-removedirectory-file-removal)
    - [4.6.4 RemoveDirectory Directory Metadata Cleanup](#464-removedirectory-directory-metadata-cleanup)
    - [4.6.5 Package State Update](#465-package-state-update)
  - [4.7 RemoveDirectory Error conditions](#47-removedirectory-error-conditions)
  - [4.8 RemoveDirectory Usage notes](#48-removedirectory-usage-notes)

---

## 0. Overview

This document specifies file removal operations.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [Core Package Interface](api_core.md)

## 1. File Removal Semantics and Multi-Path Files

This document specifies high-level file removal methods that are counterparts to the `AddFile` family of methods.

NovusPack supports multiple paths (aliases) for the same file content through a single FileEntry.
When removing a path, the following rules apply.

- If the FileEntry has multiple paths, only the specified path is removed from the `Paths` array.
- If the removed path is the last path for that FileEntry, the entire FileEntry is removed from the package.
- This ensures that file content is only removed when all paths referencing it have been removed.

## 2. RemoveFile Package Method

This section describes the RemoveFile method for removing files from packages.

### 2.1 RemoveFile Purpose

Removes a file from the package by its package-internal path.
This is the high-level counterpart to AddFile.

This method provides a symmetric API with AddFile: where AddFile adds files, RemoveFile removes them.

### 2.2 Package.RemoveFile Method

```go
// RemoveFile removes a file from the package.
// High-level counterpart to AddFile.
// Returns *PackageError on failure
func (p *Package) RemoveFile(ctx context.Context, path string) error
```

### 2.3 RemoveFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Package-internal path of the file to remove

### 2.4 RemoveFile Returns

- `error`: Returns `*PackageError` on failure

### 2.5 RemoveFile Behavior

- Remove the file from the in-memory package state.
- If the target FileEntry has multiple paths, only the matching path is removed.
- If the removed path is the last path for that entry, the entire FileEntry is removed.
- Update `PackageInfo` to reflect the new in-memory package state.
- Increment `PackageInfo.PackageDataVersion`.
- Changes become durable only after Write, SafeWrite, or FastWrite completes successfully.

### 2.6. RemoveFile Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: File does not exist at the specified path
- `ErrTypeValidation`: Invalid or malformed file path
- `ErrTypeValidation`: Package is in read-only mode
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 2.7. RemoveFile Usage Notes

RemoveFile provides naming symmetry with AddFile for intuitive file management operations.

```go
// Add and remove files with symmetric API
pkg.AddFile(ctx, "/path/to/file.txt", source, nil)
pkg.RemoveFile(ctx, "/file.txt")

// Memory-based operations
pkg.AddFileFromMemory(ctx, "/file.txt", data, nil)
pkg.RemoveFile(ctx, "/file.txt")
```

## 3. RemoveFilePattern Package Method

This section describes the RemoveFilePattern method for removing files by pattern.

### 3.1 RemoveFilePattern Purpose

Removes multiple files from the package based on a file system pattern.
This is the high-level counterpart to AddFilePattern.

This method provides a symmetric API with AddFilePattern: where AddFilePattern adds files by pattern, RemoveFilePattern removes them.

### 3.2 Package.RemoveFilePattern Method

```go
// RemoveFilePattern removes files matching a pattern from the package.
// High-level counterpart to AddFilePattern.
// Returns *PackageError on failure
func (p *Package) RemoveFilePattern(ctx context.Context, pattern string) ([]string, error)
```

### 3.3 RemoveFilePattern Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: File system pattern (e.g., "*.txt", "documents/**/*.pdf")

### 3.4 RemoveFilePattern Returns

- `[]string`: Slice of removed file paths (stored path format).
  Contains all paths that were successfully removed, even if an error occurred during processing.
- `error`: Any error that occurred during file removal.
  If an error occurs, some files may have been removed successfully before the error.

### 3.5 RemoveFilePattern Behavior

RemoveFilePattern removes files matching a pattern from the package.

#### 3.5.1 Pattern Matching

- Scan package for files matching the pattern.
- Process pattern matching using the same rules as AddFilePattern.
- Identify all files to be removed.

#### 3.5.2 File Removal

- Remove each matching file from the in-memory package state.
- If a FileEntry has multiple paths, only remove the matching path.
- If a removed path is the last path for an entry, remove the entire FileEntry.

#### 3.5.3 Directory Metadata Cleanup

- After removing files, check for directory metadata entries that no longer contain any files.
- Remove directory metadata entries for empty directories using RemoveDirectoryMetadata from the Metadata API.
- This ensures directory metadata is cleaned up when all files in a directory are removed via pattern.

#### 3.5.4 Package State Updates

- Update package metadata and file count in memory.
- Increment `PackageInfo.PackageDataVersion`.
- Changes become durable only after Write, SafeWrite, or FastWrite completes successfully.

### 3.6. RemoveFilePattern Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Invalid or malformed pattern
- `ErrTypeValidation`: No files match the pattern
- `ErrTypeValidation`: Package is in read-only mode
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

#### 3.6.1 Partial Failure Handling

If an error occurs during pattern-based removal (e.g., context cancellation, validation error on a specific file), some files may have been successfully removed before the error occurred.
In this case:

- The method returns both a non-empty `[]string` (containing all successfully removed paths) and a non-nil `error`
- Callers can inspect the returned paths to determine which removals completed successfully
- The error indicates why processing stopped, but does not invalidate the successful removals that occurred before the error

### 3.7. RemoveFilePattern Usage Notes

RemoveFilePattern is the high-level batch file removal method that provides naming symmetry with AddFilePattern.

Returns a slice of removed file paths (as strings in stored path format) for logging and verification purposes.

When removing files via pattern causes directories to become empty, RemoveFilePattern automatically removes directory metadata entries using RemoveDirectoryMetadata from the Metadata API.
This ensures that directory metadata is kept synchronized with the file structure.

```go
// Symmetric batch operations
entries, _ := pkg.AddFilePattern(ctx, "/path/to/*.txt", nil)
removedPaths, err := pkg.RemoveFilePattern(ctx, "*.txt")
if err != nil {
    // Some files may have been removed before the error
    // removedPaths contains all successfully removed paths
    log.Printf("Removed %d files before error: %v", len(removedPaths), removedPaths)
    log.Printf("Error: %v", err)
} else {
    // All files removed successfully
    log.Printf("Successfully removed %d files", len(removedPaths))
}
// If removing *.txt files empties any directories, their metadata is also removed
```

For directory metadata operations, see [Package Metadata API - Path Metadata System](api_metadata.md).

## 4. RemoveDirectory Package Method

This section describes the RemoveDirectory method for removing directories.

### 4.1 RemoveDirectory Purpose

Removes all files within a directory path from the package.
This is the high-level counterpart to AddDirectory.

This method provides a symmetric API with AddDirectory: where AddDirectory adds entire directory trees, RemoveDirectory removes them.

### 4.2 Package.RemoveDirectory Method

```go
// RemoveDirectory removes all files within a directory path from the package.
// High-level counterpart to AddDirectory.
// Returns *PackageError on failure
func (p *Package) RemoveDirectory(ctx context.Context, dirPath string, options *RemoveDirectoryOptions) ([]string, error)
```

### 4.3 RemoveDirectory Parameters

- `ctx`: Context for cancellation and timeout handling
- `dirPath`: Package-internal directory path (e.g., "/project/src")
- `options`: Configuration options for directory removal (can be nil for defaults)

### 4.4 RemoveDirectoryOptions Struct

```go
// RemoveDirectoryOptions configures directory removal behavior
type RemoveDirectoryOptions struct {
    // Recursive controls whether to remove files in subdirectories
    Recursive Option[bool] // Default: true

    // Pattern filters which files to remove (e.g., "*.txt")
    Pattern Option[string] // Default: all files

    // RemoveEmptyDirs controls whether to remove directory metadata entries
    // when all files in a directory are removed
    RemoveEmptyDirs Option[bool] // Default: true
}
```

### 4.5 RemoveDirectory Returns

- `[]string`: Slice of removed file paths (stored path format)
- `error`: Any error that occurred during directory removal (if error occurs, some files may have been removed successfully)

### 4.6 RemoveDirectory Behavior

RemoveDirectory recursively removes files under the specified directory path from the package.

#### 4.6.1 Directory Path Validation

- Validate that the package is open and writable.
- Normalize `dirPath` to stored path format (leading "/", forward slashes).
- Validate that `dirPath` is a valid directory path.

#### 4.6.2 File Discovery

- Scan package for files under `dirPath`.
- If Recursive is true (default), include files in all subdirectories.
- If Recursive is false, only include files directly in `dirPath`.
- If Pattern is set, filter files by pattern match.

#### 4.6.3 RemoveDirectory File Removal

- For each discovered file, remove it from the in-memory package state.
- If a FileEntry has multiple paths, only remove paths under `dirPath`.
- If a removed path is the last path for an entry, remove the entire FileEntry.
- Process removals in a consistent order (lexicographic by path).

#### 4.6.4 RemoveDirectory Directory Metadata Cleanup

- If RemoveEmptyDirs is true (default):
  - After removing files, check for directory metadata entries under `dirPath`.
  - Remove directory metadata entries that no longer contain any files.
  - Use RemoveDirectoryMetadata from the Metadata API for each empty directory.
- If RemoveEmptyDirs is false:
  - Preserve all directory metadata entries even if directories become empty.

#### 4.6.5 Package State Update

- Update package metadata and file count in memory.
- Increment `PackageInfo.PackageDataVersion`.
- Changes become durable only after Write, SafeWrite, or FastWrite completes successfully.

### 4.7. RemoveDirectory Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Invalid or malformed directory path
- `ErrTypeValidation`: Directory path does not exist in package
- `ErrTypeValidation`: No files found under directory path
- `ErrTypeValidation`: Package is in read-only mode
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 4.8. RemoveDirectory Usage Notes

RemoveDirectory provides naming symmetry with AddDirectory for intuitive directory management operations.

Returns a slice of removed file paths (as strings in stored path format) for logging and verification purposes.

```go
// Symmetric directory operations
entries, _ := pkg.AddDirectory(ctx, "/path/to/project", nil)
removedPaths, _ := pkg.RemoveDirectory(ctx, "/project", nil)
// removedPaths contains the list of removed file paths

// Remove with pattern filter
options := &RemoveDirectoryOptions{
    Pattern: Option.Some("*.tmp"),
    Recursive: Option.Some(true),
}
removedPaths, _ := pkg.RemoveDirectory(ctx, "/cache", options)

// Remove files but preserve directory metadata
options := &RemoveDirectoryOptions{
    RemoveEmptyDirs: Option.Some(false),
}
removedPaths, _ := pkg.RemoveDirectory(ctx, "/templates", options)
```

For directory metadata operations, see [Package Metadata API - Path Metadata System](api_metadata.md).
