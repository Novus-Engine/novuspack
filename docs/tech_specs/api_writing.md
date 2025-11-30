# NovusPack Technical Specifications - Package Writing Operations API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. SafeWrite - Atomic Package Writing](#1-safewrite---atomic-package-writing)
  - [1.1 SafeWrite Method Signature](#11-safewrite-method-signature)
  - [1.2 SafeWrite Implementation Strategy](#12-safewrite-implementation-strategy)
  - [1.3 SafeWrite Use Cases](#13-safewrite-use-cases)
  - [1.4 SafeWrite Performance Characteristics](#14-safewrite-performance-characteristics)
  - [1.5 Streaming Implementation](#15-streaming-implementation)
  - [1.6 SafeWrite Error Handling](#16-safewrite-error-handling)
- [2. FastWrite - In-Place Package Updates](#2-fastwrite---in-place-package-updates)
  - [2.1 FastWrite Method Signature](#21-fastwrite-method-signature)
  - [2.2 FastWrite Implementation Strategy](#22-fastwrite-implementation-strategy)
  - [2.3 FastWrite Use Cases](#23-fastwrite-use-cases)
  - [2.4 FastWrite Performance Characteristics](#24-fastwrite-performance-characteristics)
  - [2.5 FastWrite Error Handling](#25-fastwrite-error-handling)
- [3. Write Strategy Selection](#3-write-strategy-selection)
  - [3.1 Automatic Selection Logic](#31-automatic-selection-logic)
  - [3.2 Selection Criteria](#32-selection-criteria)
  - [3.3 Performance Comparison](#33-performance-comparison)
- [4. Signed File Write Operations](#4-signed-file-write-operations)
  - [4.1 Signed File Protection](#41-signed-file-protection)
  - [4.2 Clear-Signatures Flag](#42-clear-signatures-flag)
  - [4.3 Clear-Signatures Behavior](#43-clear-signatures-behavior)
  - [4.4 Error Conditions](#44-error-conditions)
  - [4.5 Use Cases](#45-use-cases)
  - [4.6 Security Considerations](#46-security-considerations)
- [5. Compressed Package Write Operations](#5-compressed-package-write-operations)
  - [5.1 Compressed Package Detection](#51-compressed-package-detection)
  - [5.2 Write Operations on Compressed Packages](#52-write-operations-on-compressed-packages)
  - [5.3 Compression Operations](#53-compression-operations)
  - [5.4 Compression and Signing Relationship](#54-compression-and-signing-relationship)
  - [5.5 Compression Strategy Selection](#55-compression-strategy-selection)
  - [5.6 Performance Considerations](#56-performance-considerations)
  - [5.7 Error Handling](#57-error-handling)
  - [5.8 Use Cases](#58-use-cases)

---

## 0. Overview

This document defines the package writing operations for the NovusPack system, including atomic writing, in-place updates, strategy selection, and handling of signed packages.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation

## 1. SafeWrite - Atomic Package Writing

**Purpose**: Provides atomic package writing with guaranteed consistency and rollback capability.

### 1.1 SafeWrite Method Signature

```go
func (p *Package) SafeWrite(ctx context.Context, path string, compressionType uint8) error
```

### 1.2 SafeWrite Implementation Strategy

- **Temp File Creation**: Creates temporary file in same directory as target
- **Streaming Write**: Streams data from source package or temp files for large content
- **Memory Management**: Uses in-memory data for small files, streaming for large files
- **Atomic Rename**: Atomically renames temp file to target path
- **Rollback**: Automatically cleans up temp file on failure
- **Compression Parameter**: compressionType specifies whether and how to compress the package
- **Compression Scope**: Compresses file entries + data + index (NOT header, comment, or signatures)

### 1.3 SafeWrite Use Cases

- **New Packages**: Creating packages from scratch
- **Complete Rewrites**: When entire package content changes
- **Critical Operations**: When data integrity is paramount
- **Large Packages**: Handles packages of any size through intelligent streaming
- **Defragmentation**: Complete package reorganization with guaranteed atomicity
- **Compressed Packages**: All compressed packages (SafeWrite required for compression)
- **Signed Packages**: All signed packages with clearSignatures flag

### 1.4 SafeWrite Performance Characteristics

- **Speed**: Slower than FastWrite (complete rewrite required)
- **Memory**: Intelligent memory usage (streaming for large files, in-memory for small files)
- **Disk I/O**: Higher disk usage (temporary file + final file)
- **Safety**: Highest safety level with guaranteed atomicity
- **Scalability**: Handles packages of any size through streaming

### 1.5 Streaming Implementation

- **Small Files (<10MB)**: Written directly from memory to temp file
- **Medium Files (10MB-100MB)**: Streamed from source package file if available
- **Large Files (>100MB)**: Streamed from temp files or source with chunked processing
- **Memory Thresholds**: Configurable thresholds for memory vs streaming decisions
- **Buffer Management**: Uses buffer pool for efficient streaming operations
- **Source Detection**: Automatically detects if data can be streamed from existing package

### 1.6 SafeWrite Error Handling

- **Directory Validation**: Ensures target directory exists
- **Temp File Cleanup**: Automatic cleanup on failure
- **Rollback**: No partial writes possible
- **Error Propagation**: Clear error messages for debugging
- **Streaming Errors**: Handles streaming failures gracefully with cleanup

## 2. FastWrite - In-Place Package Updates

**Purpose**: Provides efficient in-place updates for existing packages with minimal I/O.

### 2.1 FastWrite Method Signature

```go
func (p *Package) FastWrite(ctx context.Context, path string, compressionType uint8) error
```

### 2.2 FastWrite Implementation Strategy

- **Entry Comparison**: Compares existing vs new file entries
- **Change Detection**: Identifies modified, added, and unchanged entries
- **In-Place Updates**: Updates only changed entries in existing file
- **Append New Data**: Appends new entries to end of file
- **Metadata Updates**: Updates file index and offsets

### 2.3 FastWrite Use Cases

- **Incremental Updates**: Adding/modifying individual files
- **Large Package Modifications**: When SafeWrite would be too slow
- **Frequent Updates**: When performance is critical
- **Existing Packages**: When target package already exists
- **Unsigned Packages Only**: Cannot be used with signed packages (SignatureOffset > 0)
- **Uncompressed Packages Only**: Cannot be used with compressed packages

### 2.4 FastWrite Performance Characteristics

- **Speed**: Much faster than SafeWrite for updates
- **Memory**: Lower memory usage (only changed data in memory)
- **Disk I/O**: Minimal disk usage (only changed data written)
- **Safety**: Good safety with partial update recovery

### 2.5 FastWrite Error Handling

- **Entry Validation**: Validates existing package before modification
- **Partial Recovery**: Can recover from partial update failures
- **Change Tracking**: Tracks what was successfully updated
- **Fallback**: Falls back to SafeWrite if FastWrite fails
- **Signed File Check**: Returns error if attempting to use FastWrite on signed packages
- **Compressed File Check**: Returns error if attempting to use FastWrite on compressed packages

## 3. Write Strategy Selection

### 3.1 Automatic Selection Logic

The `Write(path string) error` method automatically selects the appropriate writing strategy:

- **New package**: If file doesn't exist, uses SafeWrite for new package creation
- **Complete rewrite**: If package requires complete rewrite, uses SafeWrite
- **In-place updates**: Attempts FastWrite first for existing unsigned packages
- **Signed packages**: Refuses write operations unless clearSignatures flag is provided (which will write a new file without the signatures)
- **Compressed packages**: Always uses SafeWrite (FastWrite not supported for compressed packages)
- **Fallback strategy**: Falls back to SafeWrite if FastWrite fails
- **Success**: Returns nil on successful write operation

### 3.2 Selection Criteria

| Scenario              | Recommended Method               | Reason                                               |
| --------------------- | -------------------------------- | ---------------------------------------------------- |
| New Package           | SafeWrite                        | No existing file to update                           |
| Complete Rewrite      | SafeWrite                        | Simpler and safer for full replacement               |
| Single File Addition  | FastWrite                        | Minimal I/O overhead                                 |
| Multiple File Changes | FastWrite                        | Efficient for incremental updates                    |
| Large Package (>1GB)  | FastWrite                        | Memory and I/O efficiency                            |
| Critical Data         | SafeWrite                        | Maximum safety and atomicity                         |
| Frequent Updates      | FastWrite                        | Performance optimization                             |
| Signed Package        | SafeWrite (with clearSignatures) | Signed files are immutable, require complete rewrite |
| Compressed Package    | SafeWrite                        | FastWrite not supported for compressed packages      |

### 3.3 Performance Comparison

| Metric                | SafeWrite     | FastWrite        |
| --------------------- | ------------- | ---------------- |
| New Package           | Fast          | N/A              |
| Single File Update    | Slow          | Very Fast        |
| Multiple File Updates | Slow          | Fast             |
| Complete Rewrite      | Fast          | Slow             |
| Memory Usage          | Intelligent   | Low              |
| Disk I/O              | High          | Low              |
| Safety Level          | Maximum       | Good             |
| Recovery              | Full Rollback | Partial Recovery |
| Scalability           | Excellent     | Good             |

## 4. Signed File Write Operations

**Purpose**: Defines behavior for write operations on signed packages to prevent accidental signature invalidation.

### 4.1 Signed File Protection

When attempting to write to a signed package (SignatureOffset > 0), write operations are refused by default to prevent signature invalidation.

### 4.2 Clear-Signatures Flag

```go
func (p *Package) Write(ctx context.Context, path string, compressionType uint8, clearSignatures bool) error
```

- **clearSignatures = false**: Write operation is refused for signed packages
- **clearSignatures = true**: Allows writing to signed packages, but creates new unsigned file
- **compressionType = 0**: No compression (uncompressed package)
- **compressionType = 1-3**: Apply specified compression type (1=Zstd, 2=LZ4, 3=LZMA)
- **Compression Handling**: Compresses file entries + data + index (NOT header, comment, or signatures)

### 4.3 Clear-Signatures Behavior

When `clearSignatures = true` is passed:

1. **New File Creation**: Creates a new, unsigned package file using SafeWrite (complete rewrite)
2. **Filename Requirement**: New filename must be different from the current signed file
3. **Signature Removal**: All signatures are stripped from the new file
4. **Content Preservation**: All package content (files, metadata, comments) is preserved
5. **Immutability Reset**: The new file can be modified normally (not immutable)
6. **Write Strategy**: Always uses SafeWrite since signed files cannot be modified in-place

### 4.4 Error Conditions

- **SignedFileError**: Returned when attempting to write to signed file without clearSignatures flag
- **SameFilenameError**: Returned when clearSignatures=true but new filename matches current filename
- **ValidationError**: Returned if the signed file is corrupted or invalid

### 4.5 Use Cases

- **Development Workflow**: Clear signatures to continue development on a previously signed package
- **Package Modification**: Make changes to a signed package while preserving content
- **Signature Management**: Remove signatures before re-signing with different keys
- **Testing**: Create unsigned copies of signed packages for testing

### 4.6 Security Considerations

- **Explicit Intent**: Requires explicit flag to prevent accidental signature clearing
- **Filename Protection**: Prevents accidental overwrite of signed files
- **Audit Trail**: Clear-signatures operations should be logged for security auditing
- **Backup Recommendation**: Users should backup signed files before clearing signatures

## 5. Compressed Package Write Operations

**Purpose**: Defines behavior for write operations on package-compressed packages to ensure proper handling of compression state.

### 5.1 Compressed Package Detection

When writing packages, the system must check for package compression:

- **PackageCompression Field**: Check Bit 15-8 in header flags for compression type
- **IsPackageCompressed**: Boolean flag indicating if package is compressed
- **Compression Type**: 0=none, 1=Zstd, 2=LZ4, 3=LZMA

### 5.2 Write Operations on Compressed Packages

#### 5.2.1 SafeWrite with Compressed Packages

```go
func (p *Package) SafeWrite(ctx context.Context, path string, compressionType uint8) error
```

##### 5.2.1.1 Behavior for Compressed Packages

- **Decompression Required**: Package must be decompressed before writing
- **Compression Preservation**: Original compression settings are preserved in header
- **Recompression**: Package is recompressed after writing if it was originally compressed
- **Memory Management**: Uses streaming for large compressed packages
- **Header/Comment/Signature Access**: Header, comment, and signatures remain uncompressed for direct access

#### 5.2.2 FastWrite with Compressed Packages

```go
func (p *Package) FastWrite(ctx context.Context, path string, compressionType uint8) error
```

##### 5.2.2.1 FastWrite Behavior for Compressed Packages

- **Not Supported**: FastWrite cannot be used with compressed packages
- **Automatic Fallback**: Falls back to SafeWrite for compressed packages
- **Error Prevention**: Returns error if attempting FastWrite on compressed package

### 5.3 Compression Operations

See [Package Compression API](api_package_compression.md) for detailed compression method signatures and implementation details.

#### 5.3.1 In-Memory Compression Methods

- **CompressPackage**: Compresses package content in memory
- **DecompressPackage**: Decompresses the package in memory

#### 5.3.2 File-Based Compression Methods

- **CompressPackageFile**: Compresses package content and writes to specified path
- **DecompressPackageFile**: Decompresses the package and writes to specified path

#### 5.3.3 Write Method Compression Handling

```go
func (p *Package) Write(ctx context.Context, path string, compressionType uint8, clearSignatures bool) error
```

- **Compression Parameter**: compressionType specifies compression type (0 = none, 1-3 = specific types)
- **Internal Process**: Uses internal compression methods before writing
- **Method Selection**: Uses SafeWrite for compressed packages, FastWrite for uncompressed packages

### 5.4 Compression and Signing Relationship

#### 5.4.1 Signing Compressed Packages

**Supported Operation**: Compressed packages can be signed

- **Process**: Compress package first, then add signatures
- **Header Access**: Header remains uncompressed for signature validation
- **Comment Access**: Comment remains uncompressed for easy reading
- **Signature Access**: Signatures remain uncompressed for validation
- **Validation**: Signatures validate the compressed content

#### 5.4.2 Compressing Signed Packages

**Unsupported Operation**: Signed packages cannot be compressed

- **Reason**: Would require decompression to access signatures for validation
- **Error**: Returns error if attempting to compress signed package
- **Workflow**: Must clear signatures first, then compress, then re-sign

### 5.5 Compression Strategy Selection

#### 5.5.1 Automatic Compression Detection

The `Write(path string) error` method automatically handles compression:

- **Preserve State**: Maintains current compression state by default
- **Compressed Input**: If reading compressed package, writes compressed package
- **Uncompressed Input**: If reading uncompressed package, writes uncompressed package
- **New Package**: New packages are uncompressed by default
- **Signed Package Check**: Refuses compression if package is signed

#### 5.5.2 Compression Workflow Options

Multiple approaches for handling compression:

- **In-Memory Workflow**: Use `CompressPackage()`/`DecompressPackage()` then `Write()`
- **File-Based Workflow**: Use `CompressPackageFile()`/`DecompressPackageFile()` directly
- **Write with Compression**: Use `Write()` with compressionType parameter
- **State Management**: In-memory methods update package state, file methods don't affect in-memory state
- **Signed Package Check**: All compression methods return error if package is signed

### 5.6 Performance Considerations

#### 5.6.1 Compressed Package Performance

| Operation        | Compressed Package           | Uncompressed Package   |
| ---------------- | ---------------------------- | ---------------------- |
| Read Speed       | Slower (decompression)       | Faster (direct access) |
| Write Speed      | Slower (compression)         | Faster (direct write)  |
| Disk Usage       | Lower                        | Higher                 |
| Memory Usage     | Higher (compression buffers) | Lower                  |
| Network Transfer | Faster                       | Slower                 |

#### 5.6.2 Compression Decision Factors

- **Package Size**: Small packages may not benefit from compression
- **File Count**: Many small files benefit from package compression
- **Content Type**: Text and structured data compress better than binary data
- **Use Case**: Archival vs. frequent access scenarios
- **Network Transfer**: Compressed packages transfer faster over networks

### 5.7 Error Handling

#### 5.7.1 Compression Errors

- **CompressionFailure**: Returned when compression operation fails during write
- **DecompressionFailure**: Returned when decompression operation fails during write
- **UnsupportedCompression**: Returned for unsupported compression types
- **CorruptedCompressedData**: Returned when compressed data is corrupted
- **CompressSignedPackageError**: Returned when attempting to write compressed signed package

#### 5.7.2 Write Strategy Errors

- **FastWriteOnCompressed**: Returned when attempting FastWrite on compressed package
- **CompressionMismatch**: Returned when compression type doesn't match expectations
- **MemoryInsufficient**: Returned when insufficient memory for compression operations

### 5.8 Use Cases

#### 5.8.1 When to Use Compressed Packages

- **Archival Storage**: Long-term storage where space is more important than speed
- **Network Distribution**: Packages distributed over networks
- **Small File Collections**: Packages with many small files
- **Text-Heavy Content**: Packages containing primarily text or structured data

#### 5.8.2 When to Use Uncompressed Packages

- **Frequent Access**: Packages accessed frequently where speed is critical
- **Large Binary Files**: Packages with large binary files that don't compress well
- **Development Workflow**: Packages being modified frequently during development
- **Memory Constraints**: Systems with limited memory for compression operations

---

_This document defines the package writing operations for NovusPack. For core package operations, see the Core Package Interface API._
