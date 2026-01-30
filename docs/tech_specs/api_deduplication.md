# NovusPack Technical Specifications - Multi-Layer Deduplication System API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Deduplication Strategy](#1-deduplication-strategy)
  - [1.1 Deduplication Layers](#11-deduplication-layers)
    - [1.1.1 Layer 1: Size Check (Instant Elimination)](#111-layer-1-size-check-instant-elimination)
    - [1.1.2 Layer 2: CRC32 Check (Fast Elimination)](#112-layer-2-crc32-check-fast-elimination)
    - [1.1.3 Layer 3: SHA256 Check (Hash-on-Demand)](#113-layer-3-sha256-check-hash-on-demand)
  - [1.2 Deduplication Implementation Strategy](#12-deduplication-implementation-strategy)
    - [1.2.1 findExistingEntry Function](#121-findexistingentry-function)
    - [1.2.2 PathHandling Integration](#122-pathhandling-integration)
  - [1.3 Deduplication Performance Characteristics](#13-deduplication-performance-characteristics)
  - [1.4 Deduplication Use Cases](#14-deduplication-use-cases)
  - [1.5 Encryption and Deduplication](#15-encryption-and-deduplication)
- [2. Deduplication at Different Processing Levels](#2-deduplication-at-different-processing-levels)
  - [2.1 Raw Content Deduplication](#21-raw-content-deduplication)
  - [2.2 Processed Content Deduplication](#22-processed-content-deduplication)
  - [2.3 Final Content Deduplication](#23-final-content-deduplication)
  - [2.4 Deduplication Level Selection](#24-deduplication-level-selection)
    - [2.4.1 selectDeduplicationLevel(entry \*FileEntry) DeduplicationLevel](#241-selectdeduplicationlevelentry-fileentry-deduplicationlevel)
- [3. Deduplication API Methods](#3-deduplication-api-methods)
  - [3.1 File Deduplication](#31-file-deduplication)
    - [3.1.4 File Deduplication Purpose](#314-file-deduplication-purpose)
    - [3.1.5 FindExistingEntryByCRC32 Parameters](#315-findexistingentrybycrc32-parameters)
    - [3.1.6 FindExistingEntryMultiLayer Parameters](#316-findexistingentrymultilayer-parameters)
    - [3.1.7 AddPathToExistingEntry Parameters](#317-addpathtoexistingentry-parameters)
    - [3.1.8 File Deduplication Behavior](#318-file-deduplication-behavior)
    - [3.1.9 File Deduplication Usage Notes](#319-file-deduplication-usage-notes)

---

## 0. Overview

This document defines the multi-layer deduplication system for the NovusPack system, providing efficient content deduplication using a layered approach for optimal performance.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Streaming and Buffer Management](api_streaming.md) - File streaming interface and buffer management system
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation

## 1. Deduplication Strategy

**Purpose**: Provides efficient content deduplication using a multi-layer approach for optimal performance.

### 1.1 Deduplication Layers

This section describes the different layers of deduplication in the system.

#### 1.1.1 Layer 1: Size Check (Instant Elimination)

- Compares file sizes for instant elimination of different files
- Zero computational cost, maximum efficiency
- Eliminates 99%+ of non-matches instantly

#### 1.1.2 Layer 2: CRC32 Check (Fast Elimination)

- Uses existing CRC32 checksums for fast comparison
- Leverages existing infrastructure
- Provides good collision resistance for most use cases

#### 1.1.3 Layer 3: SHA256 Check (Hash-on-Demand)

- Computes SHA256 only when size + CRC32 match
- Provides cryptographic collision resistance
- Minimal computational overhead (only for potential matches)

### 1.2 Deduplication Implementation Strategy

This section describes the implementation strategy for deduplication.

#### 1.2.1. FindExistingEntry Function

- **Layer 1 - Size check**: Instant elimination of files with different sizes
- **Layer 2 - CRC32 check**: Fast comparison using existing CRC32 checksums
- **Layer 3 - SHA256 check**: Hash-on-demand comparison for collision prevention
- **Return**: Matching FileEntry if found, nil if no match

#### 1.2.2 PathHandling Integration

Deduplication integrates with the `PathHandling` option from `AddFileOptions` to control how duplicate content is handled:

- **PathHandlingHardLinks** (default): When duplicate content is found, add the path to the existing FileEntry (current behavior)
- **PathHandlingSymlinks**: When duplicate content is found, create a symlink pointing to the primary path instead of adding the path
- **PathHandlingDefault**: Use package default (`Package.DefaultPathHandling`)
- **PathHandlingPreserve**: Preserve original filesystem behavior

##### 1.2.2.1 Automatic Symlink Creation

When `Package.AutoConvertToSymlinks` is enabled, deduplication automatically creates symlinks instead of adding paths to existing FileEntry objects.

##### 1.2.2.2 Integration with AddFile

The `AddFileOptions.PathHandling` field controls behavior during file addition:

- If `PathHandlingSymlinks` is set and duplicate content is found, a symlink is created
- If `PathHandlingHardLinks` is set, the path is added to the existing FileEntry (default behavior)
- The `PrimaryPathSelector` function can be used to choose which path becomes the primary when creating symlinks

See [File Addition API - PathHandling Type](api_file_mgmt_addition.md#27-pathhandling-type) for complete PathHandling documentation.

### 1.3 Deduplication Performance Characteristics

- **Layer 1**: O(1) per comparison, eliminates 99%+ of non-matches
- **Layer 2**: O(1) per comparison, uses existing CRC32 infrastructure
- **Layer 3**: O(n) hash computation, only for potential matches
- **Overall**: Near-optimal performance with cryptographic security

### 1.4 Deduplication Use Cases

- **Content Deduplication**: Eliminates duplicate file content
- **Storage Optimization**: Reduces package size significantly
- **Performance**: Fast deduplication with minimal overhead
- **Security**: Cryptographic collision resistance when needed

### 1.5 Encryption and Deduplication

**Important**: Files that are encrypted separately (even if they have identical original content) will **not** be deduplicated because:

- Each encryption operation produces different encrypted content due to:
  - Random initialization vectors (IVs)
  - Different encryption keys
  - Non-deterministic encryption algorithms
- Files with different paths that are encrypted separately are treated as distinct files
- Only files that are encrypted with the same key and parameters can potentially be deduplicated at the encrypted level

## 2. Deduplication at Different Processing Levels

This section describes how deduplication works at different processing levels.

### 2.1 Raw Content Deduplication

- Compares original file content before any processing
- Uses `OriginalSize`, `RawChecksum`, and raw `ContentHash`
- Eliminates exact duplicate files

### 2.2 Processed Content Deduplication

- Compares content after compression but before encryption
- Uses processed size, processed checksum, and processed hash
- Eliminates files that compress to identical content

### 2.3 Final Content Deduplication

- Compares final stored content (after compression + encryption)
- Uses `Size`, `StoredChecksum`, and final content hash
- Eliminates files that result in identical stored data

### 2.4 Deduplication Level Selection

This section describes how deduplication levels are selected.

#### 2.4.1. SelectDeduplicationLevel(entry \*FileEntry) DeduplicationLevel

- **Encrypted files**: Returns DeduplicationLevelProcessed (deduplicate before encryption)
- **Compressed files**: Returns DeduplicationLevelProcessed (deduplicate before compression)
- **Raw files**: Returns DeduplicationLevelRaw (deduplicate raw content)

## 3. Deduplication API Methods

The deduplication system exposes helper methods used by file addition operations.
These methods operate on FileEntry metadata and checksums to identify duplicate content.

See [File Addition API](api_file_mgmt_addition.md) for how these methods are applied during AddFile flows.

### 3.1 File Deduplication

This section describes file-level deduplication operations.

#### 3.1.1 Package.FindExistingEntryByCRC32 Method

```go
// FindExistingEntryByCRC32 finds existing entry by CRC32 checksum
func (p *Package) FindExistingEntryByCRC32(rawChecksum uint32) *FileEntry
```

#### 3.1.2 Package.FindExistingEntryMultiLayer Method

```go
// FindExistingEntryMultiLayer performs multi-layer deduplication
func (p *Package) FindExistingEntryMultiLayer(originalSize int64, rawChecksum uint32, content []byte) (*FileEntry, []byte)
```

#### 3.1.3 Package.AddPathToExistingEntry Method

```go
// AddPathToExistingEntry adds a path to an existing entry
func (p *Package) AddPathToExistingEntry(existingEntry *FileEntry, newPath string)
```

#### 3.1.4 File Deduplication Purpose

Provides deduplication functionality to avoid storing duplicate content.

#### 3.1.5 FindExistingEntryByCRC32 Parameters

- `rawChecksum`: CRC32 checksum to search for

#### 3.1.6 FindExistingEntryMultiLayer Parameters

- `originalSize`: Original file size
- `rawChecksum`: CRC32 checksum
- `content`: File content for verification

#### 3.1.7 AddPathToExistingEntry Parameters

- `existingEntry`: Existing FileEntry
- `newPath`: New path to add

#### 3.1.8 File Deduplication Behavior

- Searches for existing files with matching checksums
- Performs multi-layer verification (CRC32 + content hash)
- Adds new paths to existing entries when duplicates found
- Reduces storage space by sharing content

#### 3.1.9 File Deduplication Usage Notes

Deduplication functions support both simple CRC32-based lookup and multi-layer verification for accurate duplicate detection.
