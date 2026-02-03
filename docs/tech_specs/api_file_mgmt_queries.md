# NovusPack Technical Specifications - File Queries API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Existence and Listing](#1-file-existence-and-listing)
  - [1.1 Package Methods](#11-package-methods)
  - [1.2 Purpose](#12-purpose)
  - [1.3 FileEntry Access](#13-fileentry-access)
  - [1.4 Usage Notes](#14-usage-notes)
- [2. Single-Entry Lookups](#2-single-entry-lookups)
  - [2.1 GetFileByPath](#21-getfilebypath)
    - [2.1.1 Package GetFileByPath Method](#211-packagegetfilebypath-method)
    - [2.1.2 GetFileByPath Purpose](#212-getfilebypath-purpose)
    - [2.1.3 GetFileByPath Parameters](#213-getfilebypath-parameters)
    - [2.1.4 GetFileByPath Returns](#214-getfilebypath-returns)
    - [2.1.5 GetFileByPath Use Cases](#215-getfilebypath-use-cases)
  - [2.2 GetFileByOffset](#22-getfilebyoffset)
    - [2.2.1 Package GetFileByOffset Method](#221-packagegetfilebyoffset-method)
    - [2.2.2 GetFileByOffset Purpose](#222-getfilebyoffset-purpose)
    - [2.2.3 GetFileByOffset Parameters](#223-getfilebyoffset-parameters)
    - [2.2.4 GetFileByOffset Returns](#224-getfilebyoffset-returns)
    - [2.2.5 GetFileByOffset Use Cases](#225-getfilebyoffset-use-cases)
  - [2.3 GetFileByFileID](#23-getfilebyfileid)
    - [2.3.1 Package GetFileByFileID Method](#231-packagegetfilebyfileid-method)
    - [2.3.2 GetFileByFileID Purpose](#232-getfilebyfileid-purpose)
    - [2.3.3 GetFileByFileID Parameters](#233-getfilebyfileid-parameters)
    - [2.3.4 GetFileByFileID Returns](#234-getfilebyfileid-returns)
    - [2.3.5 GetFileByFileID Use Cases](#235-getfilebyfileid-use-cases)
  - [2.4 GetFileByHash](#24-getfilebyhash)
    - [2.4.1 Package GetFileByHash Method](#241-packagegetfilebyhash-method)
    - [2.4.2 GetFileByHash Purpose](#242-getfilebyhash-purpose)
    - [2.4.3 GetFileByHash Parameters](#243-getfilebyhash-parameters)
    - [2.4.4 GetFileByHash Returns](#244-getfilebyhash-returns)
    - [2.4.5 GetFileByHash Use Cases](#245-getfilebyhash-use-cases)
  - [2.5 GetFileByChecksum](#25-getfilebychecksum)
    - [2.5.1 Package GetFileByChecksum Method](#251-packagegetfilebychecksum-method)
    - [2.5.2 GetFileByChecksum Purpose](#252-getfilebychecksum-purpose)
    - [2.5.3 GetFileByChecksum Parameters](#253-getfilebychecksum-parameters)
    - [2.5.4 GetFileByChecksum Returns](#254-getfilebychecksum-returns)
    - [2.5.5 GetFileByChecksum Use Cases](#255-getfilebychecksum-use-cases)
- [3. Multi-Entry Queries](#3-multi-entry-queries)
  - [3.1 FindEntriesByTag](#31-findentriesbytag)
    - [3.1.1 Package FindEntriesByTag Method](#311-packagefindentriesbytag-method)
    - [3.1.2 FindEntriesByTag Purpose](#312-findentriesbytag-purpose)
    - [3.1.3 FindEntriesByTag Parameters](#313-findentriesbytag-parameters)
    - [3.1.4 FindEntriesByTag Returns](#314-findentriesbytag-returns)
    - [3.1.5 FindEntriesByTag Use Cases](#315-findentriesbytag-use-cases)
  - [3.2 FindEntriesByType](#32-findentriesbytype)
    - [3.2.1 Package FindEntriesByType Method](#321-packagefindentriesbytype-method)
    - [3.2.2 FindEntriesByType Purpose](#322-findentriesbytype-purpose)
    - [3.2.3 FindEntriesByType Parameters](#323-findentriesbytype-parameters)
    - [3.2.4 FindEntriesByType Returns](#324-findentriesbytype-returns)
    - [3.2.5 FindEntriesByType Use Cases](#325-findentriesbytype-use-cases)
  - [3.3 FindEntriesByPathPatterns](#33-findentriesbypathpatterns)
    - [3.3.1 Package FindEntriesByPathPatterns Method](#331-packagefindentriesbypathpatterns-method)
    - [3.3.2 FindEntriesByPathPatterns Purpose](#332-findentriesbypathpatterns-purpose)
    - [3.3.3 FindEntriesByPathPatterns Parameters](#333-findentriesbypathpatterns-parameters)
    - [3.3.4 FindEntriesByPathPatterns Returns](#334-findentriesbypathpatterns-returns)
    - [3.3.5 FindEntriesByPathPatterns Use Cases](#335-findentriesbypathpatterns-use-cases)
- [4. Aggregate Queries and Filtered Lists](#4-aggregate-queries-and-filtered-lists)
  - [4.1 GetFileCount](#41-getfilecount)
    - [4.1.1 Package GetFileCount Method](#411-packagegetfilecount-method)
    - [4.1.2 GetFileCount Purpose](#412-getfilecount-purpose)
    - [4.1.3 GetFileCount Parameters](#413-getfilecount-parameters)
    - [4.1.4 GetFileCount Returns](#414-getfilecount-returns)
    - [4.1.5 GetFileCount Use Cases](#415-getfilecount-use-cases)
  - [4.2 ListCompressedFiles](#42-listcompressedfiles)
    - [4.2.1 Package ListCompressedFiles Method](#421-packagelistcompressedfiles-method)
    - [4.2.2 ListCompressedFiles Purpose](#422-listcompressedfiles-purpose)
    - [4.2.3 ListCompressedFiles Parameters](#423-listcompressedfiles-parameters)
    - [4.2.4 ListCompressedFiles Returns](#424-listcompressedfiles-returns)
    - [4.2.5 ListCompressedFiles Use Cases](#425-listcompressedfiles-use-cases)
  - [4.3 ListEncryptedFiles](#43-listencryptedfiles)
    - [4.3.1 Package ListEncryptedFiles Method](#431-packagelistencryptedfiles-method)
    - [4.3.2 ListEncryptedFiles Purpose](#432-listencryptedfiles-purpose)
    - [4.3.3 ListEncryptedFiles Parameters](#433-listencryptedfiles-parameters)
    - [4.3.4 ListEncryptedFiles Returns](#434-listencryptedfiles-returns)
    - [4.3.5 ListEncryptedFiles Use Cases](#435-listencryptedfiles-use-cases)

---

## 0. Overview

This document specifies file query and inspection operations.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [Core Package Interface](api_core.md)

## 1. File Existence and Listing

This section describes methods for checking file existence and listing files.

### 1.1 Package Methods

This section describes package-level methods for file existence and listing.

#### 1.1.1 Package.FileExists Method

```go
// FileExists checks if a file with the given path exists in the package
func (p *Package) FileExists(path string) (bool, error)
```

#### 1.1.2 Package.ListFiles Method

```go
// ListFiles returns lightweight file info for all files in the package
func (p *Package) ListFiles() ([]FileInfo, error)
```

### 1.2 Purpose

Defines basic file existence checks and listing operations.

### 1.3 FileEntry Access

Some query functions return full `*FileEntry` objects (or `[]*FileEntry` arrays).
These objects provide comprehensive file information including metadata, compression status, encryption details, checksums, and timestamps.

`ListFiles()` returns `[]FileInfo` for lightweight listing.
Use `GetFileByPath()` or other single-entry lookups when full `*FileEntry` details are required.

### 1.4 Usage Notes

To list directories in the package, use `ListDirectories()` from the [Metadata API](api_metadata.md).
This returns `[]PathInfo` containing all directory paths with their metadata.
See [Package Metadata API - Path Information Queries](api_metadata.md).

## 2. Single-Entry Lookups

Single-entry lookups retrieve one `*FileEntry` by a specific identifier (path, offset, file ID, or content hash).
When a file is not found, these methods return a `*PackageError`.

### 2.1 GetFileByPath

This section describes the GetFileByPath method for retrieving files by path.

#### 2.1.1 Package.GetFileByPath Method

```go
// GetFileByPath gets a FileEntry by path
// Returns *PackageError if file not found
func (p *Package) GetFileByPath(path string) (*FileEntry, error)
```

#### 2.1.2 GetFileByPath Purpose

Finds a FileEntry by its virtual path in the package.

#### 2.1.3 GetFileByPath Parameters

- `path`: Virtual path to look up

#### 2.1.4 GetFileByPath Returns

- `*FileEntry`: The FileEntry with the specified path
- `error`: `*PackageError` if file not found

#### 2.1.5 GetFileByPath Use Cases

- Lookup by path before extraction
- Path-based metadata inspection
- Path-based integrity checks

### 2.2 GetFileByOffset

This section describes the GetFileByOffset method for retrieving files by offset.

#### 2.2.1 Package.GetFileByOffset Method

```go
// GetFileByOffset gets a FileEntry by offset
// Returns *PackageError if file not found
func (p *Package) GetFileByOffset(offset int64) (*FileEntry, error)
```

#### 2.2.2 GetFileByOffset Purpose

Finds a FileEntry by its offset in the package file.

#### 2.2.3 GetFileByOffset Parameters

- `offset`: Package file offset to look up

#### 2.2.4 GetFileByOffset Returns

- `*FileEntry`: The FileEntry at the specified offset
- `error`: `*PackageError` if file not found

#### 2.2.5 GetFileByOffset Use Cases

- Low-level tooling and diagnostics
- Offset-based validation
- Debugging serialization or index issues

### 2.3 GetFileByFileID

This section describes the GetFileByFileID method for retrieving files by file ID.

#### 2.3.1 Package.GetFileByFileID Method

```go
// GetFileByFileID gets a FileEntry by its unique FileID
// Returns *PackageError if file not found
func (p *Package) GetFileByFileID(fileID uint64) (*FileEntry, error)
```

#### 2.3.2 GetFileByFileID Purpose

Finds a FileEntry by its unique 64-bit FileID.

#### 2.3.3 GetFileByFileID Parameters

- `fileID`: Unique 64-bit file identifier

#### 2.3.4 GetFileByFileID Returns

- `*FileEntry`: The FileEntry with the specified FileID
- `error`: `*PackageError` if file not found

#### 2.3.5 GetFileByFileID Use Cases

- Stable file references across package modifications
- Database-style lookups by primary key
- File tracking and management systems

### 2.4 GetFileByHash

This section describes the GetFileByHash method for retrieving files by hash.

#### 2.4.1 Package.GetFileByHash Method

```go
// GetFileByHash gets a FileEntry by content hash
// Returns *PackageError if file not found
func (p *Package) GetFileByHash(hashType HashType, hashData []byte) (*FileEntry, error)
```

#### 2.4.2 GetFileByHash Purpose

Finds a FileEntry by its content hash (SHA-256, SHA-512, BLAKE3, or XXH3).

#### 2.4.3 GetFileByHash Parameters

- `hashType`: Type of hash algorithm used
- `hashData`: Hash data to search for

#### 2.4.4 GetFileByHash Returns

- `*FileEntry`: The FileEntry with matching hash
- `error`: `*PackageError` if file not found

#### 2.4.5 GetFileByHash Use Cases

- Content deduplication
- Integrity verification
- Finding files by content rather than path

### 2.5 GetFileByChecksum

This section describes the GetFileByChecksum method for retrieving files by checksum.

#### 2.5.1 Package.GetFileByChecksum Method

```go
// GetFileByChecksum gets a FileEntry by CRC32 checksum
// Returns *PackageError if file not found
func (p *Package) GetFileByChecksum(checksum uint32) (*FileEntry, error)
```

#### 2.5.2 GetFileByChecksum Purpose

Finds a FileEntry by its CRC32 checksum (fast lookup).

#### 2.5.3 GetFileByChecksum Parameters

- `checksum`: CRC32 checksum value to search for

#### 2.5.4 GetFileByChecksum Returns

- `*FileEntry`: The FileEntry with matching checksum
- `error`: `*PackageError` if file not found

#### 2.5.5 GetFileByChecksum Use Cases

- Fast content identification
- Quick duplicate detection
- Lightweight file matching

## 3. Multi-Entry Queries

Multi-entry queries return a slice of `*FileEntry` values for tag, type, or pattern-based queries.

### 3.1 FindEntriesByTag

This section describes the FindEntriesByTag method for finding entries by tag.

#### 3.1.1 Package.FindEntriesByTag Method

```go
// FindEntriesByTag finds all FileEntry objects with a specific tag
func (p *Package) FindEntriesByTag(tag string) ([]*FileEntry, error)
```

#### 3.1.2 FindEntriesByTag Purpose

Finds all FileEntry objects that have a specific tag.

#### 3.1.3 FindEntriesByTag Parameters

- `tag`: Tag string to search for

#### 3.1.4 FindEntriesByTag Returns

- `[]*FileEntry`: All FileEntry objects with the specified tag
- `error`: Any error that occurred during the search

#### 3.1.5 FindEntriesByTag Use Cases

- Finding all files with a specific label
- Organizing files by category
- Tag-based file management

### 3.2 FindEntriesByType

This section describes the FindEntriesByType method for finding entries by type.

#### 3.2.1 Package.FindEntriesByType Method

```go
// FindEntriesByType finds all FileEntry objects of a specific type
func (p *Package) FindEntriesByType(fileType uint16) ([]*FileEntry, error)
```

#### 3.2.2 FindEntriesByType Purpose

Finds all FileEntry objects of a specific file type.

#### 3.2.3 FindEntriesByType Parameters

- `fileType`: File type identifier (0-65535)

#### 3.2.4 FindEntriesByType Returns

- `[]*FileEntry`: All FileEntry objects of the specified type
- `error`: Any error that occurred during the search

#### 3.2.5 FindEntriesByType Use Cases

- Finding all files of a specific format
- Type-based file processing
- File organization by category

### 3.3 FindEntriesByPathPatterns

This section describes the FindEntriesByPathPatterns method for finding entries by path patterns.

#### 3.3.1 Package.FindEntriesByPathPatterns Method

```go
// FindEntriesByPathPatterns gets files matching patterns from the package
func (p *Package) FindEntriesByPathPatterns(patterns []string) ([]*FileEntry, error)
```

#### 3.3.2 FindEntriesByPathPatterns Purpose

Finds FileEntry objects by matching their virtual paths against one or more patterns.

#### 3.3.3 FindEntriesByPathPatterns Parameters

- `patterns`: List of patterns used for matching file paths

#### 3.3.4 FindEntriesByPathPatterns Returns

- `[]*FileEntry`: All FileEntry objects matching any provided pattern
- `error`: Any error that occurred during pattern matching or search

#### 3.3.5 FindEntriesByPathPatterns Use Cases

- Bulk selection for extraction or removal
- Bulk selection for transformation pipelines
- Report generation by path pattern

## 4. Aggregate Queries and Filtered Lists

Aggregate queries return non-`FileEntry` values, or return filtered lists for common criteria.

### 4.1 GetFileCount

This section describes the GetFileCount method for getting the total file count.

#### 4.1.1 Package.GetFileCount Method

```go
// GetFileCount returns the total number of regular content files in the package
// Excludes special metadata files (types 65000-65535)
func (p *Package) GetFileCount() (int, error)
```

#### 4.1.2 GetFileCount Purpose

Returns the total number of regular content files in the package (file types 0-64999).
It does not include special metadata files (types 65000-65535).

#### 4.1.3 GetFileCount Parameters

None.

#### 4.1.4 GetFileCount Returns

- `int`: Total number of regular content files in the package (excludes special metadata files)
- `error`: Any error that occurred

#### 4.1.5 GetFileCount Use Cases

- Package statistics
- Progress tracking
- Validation and bounds checking

To get the total number of directories in the package, use `GetDirectoryCount()` from the Metadata API.
This returns the count of all directory paths with metadata.
See [Package Metadata API - Path Information Queries](api_metadata.md).

### 4.2 ListCompressedFiles

This section describes the ListCompressedFiles method for listing compressed files.

#### 4.2.1 Package.ListCompressedFiles Method

```go
// ListCompressedFiles returns all compressed file entries
func (p *Package) ListCompressedFiles() ([]*FileEntry, error)
```

#### 4.2.2 ListCompressedFiles Purpose

Returns FileEntry objects that are currently marked as compressed in the package.

#### 4.2.3 ListCompressedFiles Parameters

None.

#### 4.2.4 ListCompressedFiles Returns

- `[]*FileEntry`: All compressed FileEntry objects
- `error`: Any error that occurred during the query

#### 4.2.5 ListCompressedFiles Use Cases

- Reporting and inspection
- Verification after compression operations
- Selecting a subset for decompression

### 4.3 ListEncryptedFiles

This section describes the ListEncryptedFiles method for listing encrypted files.

#### 4.3.1 Package.ListEncryptedFiles Method

```go
// ListEncryptedFiles returns all encrypted file entries
func (p *Package) ListEncryptedFiles() ([]*FileEntry, error)
```

#### 4.3.2 ListEncryptedFiles Purpose

Returns FileEntry objects that are currently marked as encrypted in the package.

#### 4.3.3 ListEncryptedFiles Parameters

None.

#### 4.3.4 ListEncryptedFiles Returns

- `[]*FileEntry`: All encrypted FileEntry objects
- `error`: Any error that occurred during the query

#### 4.3.5 ListEncryptedFiles Use Cases

- Reporting and inspection
- Verification after encryption operations
- Selecting a subset for decryption during extraction
