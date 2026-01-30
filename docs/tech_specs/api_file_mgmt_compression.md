# NovusPack Technical Specifications - File Compression API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. FileEntry.Compress Method](#1-fileentrycompress-method)
  - [1.1 FileEntry.Compress Method Signature](#11-fileentrycompress-method-signature)
  - [1.2 FileEntry.Compress Parameters](#12-fileentrycompress-parameters)
  - [1.3 FileEntry.Compress Returns](#13-fileentrycompress-returns)
  - [1.4 FileEntry.Compress Behavior](#14-fileentrycompress-behavior)
  - [1.5 FileEntry.Compress Error conditions](#15-fileentrycompress-error-conditions)
  - [1.6 Package.CompressFile Method](#16-packagecompressfile-method)
  - [1.7 Package.CompressFile Parameters](#17-packagecompressfile-parameters)
  - [1.8 Package CompressFile Returns](#18-package-compressfile-returns)
  - [1.9 Package.CompressFile Behavior](#19-packagecompressfile-behavior)
  - [1.10 Package CompressFile Error conditions](#110-package-compressfile-error-conditions)
  - [1.11 CompressFile Usage notes](#111-compressfile-usage-notes)
- [2. FileEntry.Decompress Method](#2-fileentrydecompress-method)
  - [2.1 FileEntry.Decompress Method Signature](#21-fileentrydecompress-method-signature)
  - [2.2 FileEntry.Decompress Parameters](#22-fileentrydecompress-parameters)
  - [2.3 FileEntry.Decompress Returns](#23-fileentrydecompress-returns)
  - [2.4 FileEntry Decompress Behavior](#24-fileentry-decompress-behavior)
  - [2.5 FileEntry.Decompress Error conditions](#25-fileentrydecompress-error-conditions)
  - [2.6 Package.DecompressFile Method](#26-packagedecompressfile-method)
  - [2.7 Package.DecompressFile Parameters](#27-packagedecompressfile-parameters)
  - [2.8 Package DecompressFile Returns](#28-package-decompressfile-returns)
  - [2.9 Package.DecompressFile Behavior](#29-packagedecompressfile-behavior)
  - [2.10 Package DecompressFile Error conditions](#210-package-decompressfile-error-conditions)
  - [2.11 DecompressFile Usage notes](#211-decompressfile-usage-notes)
- [3. FileEntry.GetCompressionInfo Method](#3-fileentrygetcompressioninfo-method)
  - [3.1 FileEntry.GetCompressionInfo Method Signature](#31-fileentrygetcompressioninfo-method-signature)
  - [3.2 FileEntry GetCompressionInfo Parameters](#32-fileentry-getcompressioninfo-parameters)
  - [3.3 FileEntry.GetCompressionInfo Returns](#33-fileentrygetcompressioninfo-returns)
  - [3.4 FileEntry.GetCompressionInfo Error conditions](#34-fileentrygetcompressioninfo-error-conditions)
  - [3.5 Package.GetFileCompressionInfo Method](#35-packagegetfilecompressioninfo-method)
  - [3.6 Package GetFileCompressionInfo Parameters](#36-package-getfilecompressioninfo-parameters)
  - [3.7 Package.GetFileCompressionInfo Returns](#37-packagegetfilecompressioninfo-returns)
  - [3.8 Package GetFileCompressionInfo Error conditions](#38-package-getfilecompressioninfo-error-conditions)
- [4. FileCompressionInfo Purpose](#4-filecompressioninfo-purpose)
  - [4.1 FileCompressionInfo Struct Definition](#41-filecompressioninfo-struct-definition)

---

## 0. Overview

This document specifies file-level compression operations for compressing individual files within a package.
It is extracted from the File Management API specification.

**Note**: This API is for compressing individual files within a package.
For compressing the entire package structure (all files, metadata, and index together), see [Package Compression API](api_package_compression.md).

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [Core Package Interface](api_core.md)
- [Package Compression API](api_package_compression.md) - Package-level compression (entire package structure)
- [FileEntry API](api_file_mgmt_file_entry.md)

## 1. FileEntry.Compress Method

Compresses the content of an existing file within the package using the specified compression type.

### 1.1 FileEntry.Compress Method Signature

See [FileEntry.Compress Method](api_file_mgmt_file_entry.md#81-fileentrycompress-method) for the complete method definition.

### 1.2 FileEntry.Compress Parameters

- `ctx`: Context for cancellation and timeout handling
- `compressionType`: Compression algorithm to use

### 1.3 FileEntry.Compress Returns

- `error`: Returns `*PackageError` on failure

### 1.4 FileEntry.Compress Behavior

- Compress the file content in-place in the in-memory package state.
- Update FileEntry metadata to reflect the new compression state.
- Preserve file integrity and signatures.

### 1.5. FileEntry.Compress Error Conditions

- `ErrTypeValidation`: FileEntry is invalid
- `ErrTypeValidation`: File is already compressed
- `ErrTypeCompression`: Failed to compress file content
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 1.6 Package.CompressFile Method

```go
// CompressFile compresses an existing file in the package by path
// This is a convenience wrapper that looks up the FileEntry and calls Compress
// Returns *PackageError on failure
func (p *Package) CompressFile(ctx context.Context, path string, compressionType uint8) error
```

### 1.7 Package.CompressFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to compress
- `compressionType`: Compression algorithm to use

### 1.8 Package CompressFile Returns

- `error`: Returns `*PackageError` on failure

### 1.9 Package.CompressFile Behavior

- Looks up the FileEntry by path.
- Calls `FileEntry.Compress()` on the found entry.
- Returns errors from the lookup or compression operation.

### 1.10. Package CompressFile Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: File does not exist at the specified path
- `ErrTypeValidation`: Invalid or malformed file path
- All error conditions from [FileEntry.Compress](#15-fileentrycompress-error-conditions)

### 1.11. CompressFile Usage Notes

Compression becomes durable only after Write, SafeWrite, or FastWrite completes successfully.

## 2. FileEntry.Decompress Method

Decompresses the content of an existing file within the package.

### 2.1 FileEntry.Decompress Method Signature

See [FileEntry.Decompress Method](api_file_mgmt_file_entry.md#82-fileentrydecompress-method) for the complete method definition.

### 2.2 FileEntry.Decompress Parameters

- `ctx`: Context for cancellation and timeout handling

### 2.3 FileEntry.Decompress Returns

- `error`: Returns `*PackageError` on failure

### 2.4 FileEntry Decompress Behavior

- Decompress the file content in-place in the in-memory package state.
- Update FileEntry metadata to reflect the new compression state.
- Preserve file integrity and signatures.

### 2.5. FileEntry.Decompress Error Conditions

- `ErrTypeValidation`: FileEntry is invalid
- `ErrTypeValidation`: File is not compressed
- `ErrTypeCompression`: Failed to decompress file content
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 2.6 Package.DecompressFile Method

```go
// DecompressFile decompresses an existing file in the package by path
// This is a convenience wrapper that looks up the FileEntry and calls Decompress
// Returns *PackageError on failure
func (p *Package) DecompressFile(ctx context.Context, path string) error
```

### 2.7 Package.DecompressFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to decompress

### 2.8 Package DecompressFile Returns

- `error`: Returns `*PackageError` on failure

### 2.9 Package.DecompressFile Behavior

- Looks up the FileEntry by path.
- Calls `FileEntry.Decompress()` on the found entry.
- Returns errors from the lookup or decompression operation.

### 2.10. Package DecompressFile Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: File does not exist at the specified path
- `ErrTypeValidation`: Invalid or malformed file path
- All error conditions from [FileEntry.Decompress](#25-fileentrydecompress-error-conditions)

### 2.11. DecompressFile Usage Notes

Decompression becomes durable only after Write, SafeWrite, or FastWrite completes successfully.

## 3. FileEntry.GetCompressionInfo Method

Returns compression information for a file.

### 3.1 FileEntry.GetCompressionInfo Method Signature

See [FileEntry.GetCompressionInfo Method](api_file_mgmt_file_entry.md#83-fileentrygetcompressioninfo-method) for the complete method definition.

### 3.2 FileEntry GetCompressionInfo Parameters

None (operates on the FileEntry instance).

### 3.3 FileEntry.GetCompressionInfo Returns

- `*FileCompressionInfo`: Compression information for the file
- `error`: Any error that occurred during inspection

### 3.4. FileEntry.GetCompressionInfo Error Conditions

- `ErrTypeValidation`: FileEntry is invalid

### 3.5 Package.GetFileCompressionInfo Method

```go
// GetFileCompressionInfo gets compression information for a file by path
// This is a convenience wrapper that looks up the FileEntry and calls GetCompressionInfo
func (p *Package) GetFileCompressionInfo(path string) (*FileCompressionInfo, error)
```

### 3.6 Package GetFileCompressionInfo Parameters

- `path`: Virtual path of the file to inspect

### 3.7 Package.GetFileCompressionInfo Returns

- `*FileCompressionInfo`: Compression information for the file
- `error`: Any error that occurred during inspection

### 3.8. Package GetFileCompressionInfo Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: File does not exist at the specified path
- `ErrTypeValidation`: Invalid or malformed file path
- All error conditions from [FileEntry.GetCompressionInfo](#34-fileentrygetcompressioninfo-error-conditions)

## 4. FileCompressionInfo Purpose

Captures file compression details for inspection and reporting.

### 4.1 FileCompressionInfo Struct Definition

```go
// FileCompressionInfo contains file compression details
type FileCompressionInfo struct {
    IsCompressed     bool
    CompressionType  uint8
    OriginalSize     int64
    CompressedSize   int64
    CompressionRatio float64
}
```
