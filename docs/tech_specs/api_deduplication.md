# NovusPack Technical Specifications - Multi-Layer Deduplication System API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Deduplication Strategy](#1-deduplication-strategy)
  - [1.1 Deduplication Layers](#11-deduplication-layers)
  - [1.2 Deduplication Implementation Strategy](#12-deduplication-implementation-strategy)
  - [1.3 Deduplication Performance Characteristics](#13-deduplication-performance-characteristics)
  - [1.4 Deduplication Use Cases](#14-deduplication-use-cases)
  - [1.5 Encryption and Deduplication](#15-encryption-and-deduplication)
- [2. Deduplication at Different Processing Levels](#2-deduplication-at-different-processing-levels)
  - [2.1 Raw Content Deduplication](#21-raw-content-deduplication)
  - [2.2 Processed Content Deduplication](#22-processed-content-deduplication)
  - [2.3 Final Content Deduplication](#23-final-content-deduplication)
  - [2.4 Deduplication Level Selection](#24-deduplication-level-selection)

---

## 0. Overview

This document defines the multi-layer deduplication system for the NovusPack system, providing efficient content deduplication using a layered approach for optimal performance.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Streaming and Buffer Management](api_streaming.md) - File streaming interface and buffer management system
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation

## 1. Deduplication Strategy

**Purpose**: Provides efficient content deduplication using a multi-layer approach for optimal performance.

### 1.1 Deduplication Layers

1. **Layer 1: Size Check (Instant Elimination)**

    - Compares file sizes for instant elimination of different files
    - Zero computational cost, maximum efficiency
    - Eliminates 99%+ of non-matches instantly

2. **Layer 2: CRC32 Check (Fast Elimination)**

    - Uses existing CRC32 checksums for fast comparison
    - Leverages existing infrastructure
    - Provides good collision resistance for most use cases

3. **Layer 3: SHA256 Check (Hash-on-Demand)**
    - Computes SHA256 only when size + CRC32 match
    - Provides cryptographic collision resistance
    - Minimal computational overhead (only for potential matches)

### 1.2 Deduplication Implementation Strategy

#### 1.2.1 findExistingEntry(originalSize int64, rawChecksum uint32, contentHash []byte) *FileEntry

- **Layer 1 - Size check**: Instant elimination of files with different sizes
- **Layer 2 - CRC32 check**: Fast comparison using existing CRC32 checksums
- **Layer 3 - SHA256 check**: Hash-on-demand comparison for collision prevention
- **Return**: Matching FileEntry if found, nil if no match

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

#### 2.4.1 selectDeduplicationLevel(entry *FileEntry) DeduplicationLevel

- **Encrypted files**: Returns DeduplicationLevelProcessed (deduplicate before encryption)
- **Compressed files**: Returns DeduplicationLevelProcessed (deduplicate before compression)
- **Raw files**: Returns DeduplicationLevelRaw (deduplicate raw content)

---

*This document defines the multi-layer deduplication system for NovusPack. For core package operations, see the Core Package Interface API.*
