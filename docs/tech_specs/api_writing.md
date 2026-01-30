# NovusPack Technical Specifications - Package Writing Operations API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. SafeWrite - Atomic Package Writing](#1-safewrite---atomic-package-writing)
  - [1.1 SafeWrite Method Signature](#11-packagesafewrite-method)
  - [1.2 SafeWrite Implementation Strategy](#12-safewrite-implementation-strategy)
  - [1.3 SafeWrite Use Cases](#13-safewrite-use-cases)
  - [1.4 SafeWrite Performance Characteristics](#14-safewrite-performance-characteristics)
  - [1.5 Streaming Implementation](#15-streaming-implementation)
  - [1.6 SafeWrite Error Handling](#16-safewrite-error-handling)
- [2. FastWrite - In-Place Package Updates](#2-fastwrite---in-place-package-updates)
  - [2.1 FastWrite Method Signature](#21-packagefastwrite-method)
  - [2.2 FastWrite Implementation Strategy](#22-fastwrite-implementation-strategy)
  - [2.3 FastWrite Use Cases](#23-fastwrite-use-cases)
  - [2.4 FastWrite Performance Characteristics](#24-fastwrite-performance-characteristics)
  - [2.5 FastWrite Error Handling](#25-fastwrite-error-handling)
  - [2.6 FastWrite Recovery Capabilities](#26-fastwrite-recovery-capabilities)
    - [2.6.1 COW Filesystem Backup Recovery](#261-cow-filesystem-backup-recovery)
    - [2.6.2 Automatic Recovery File Creation](#262-automatic-recovery-file-creation)
    - [2.6.3 Recovery Cross-Reference](#263-recovery-cross-reference)
  - [2.7 Recovery File Format](#27-recovery-file-format)
    - [2.7.1 Recovery File Name](#271-recovery-file-name)
    - [2.7.2 Recovery File Structure](#272-recovery-file-structure)
    - [2.7.3 OpenBrokenPackage Integration](#273-openbrokenpackage-integration)
    - [2.7.4 Recovery File Lifecycle](#274-recovery-file-lifecycle)
- [3. Write Strategy Selection](#3-write-strategy-selection)
  - [3.1 Automatic Selection Logic](#31-automatic-selection-logic)
  - [3.2 Selection Criteria](#32-selection-criteria)
  - [3.3 Performance Comparison](#33-performance-comparison)
- [4. Signed File Write Operations](#4-signed-file-write-operations)
  - [4.1 Signed File Protection](#41-signed-file-protection)
  - [4.2 Signed Package Writing Behavior](#42-signed-package-writing-behavior)
    - [4.2.1 Signed Package Protection](#421-signed-package-protection)
    - [4.2.2 Compression Configuration](#422-compression-configuration)
  - [4.3 Writing Signed Package Content to New Path](#43-writing-signed-package-content-to-new-path)
  - [4.4 Signed Package Writing Error Conditions](#44-signed-package-writing-error-conditions)
  - [4.5 Signed Package Use Cases](#45-signed-package-use-cases)
  - [4.6 Security Considerations](#46-security-considerations)
- [5. Compressed Package Writing Operations](#5-compressed-package-writing-operations)
  - [5.1 Compressed Package Detection](#51-compressed-package-detection)
  - [5.2 Write Operations on Compressed Packages](#52-write-operations-on-compressed-packages)
    - [5.2.1 SafeWrite with Compressed Packages](#521-safewrite-with-compressed-packages)
    - [5.2.2 FastWrite with Compressed Packages](#522-fastwrite-with-compressed-packages)
  - [5.3 Compression Operations](#53-compression-operations)
    - [5.3.1 In-Memory Compression Methods](#531-in-memory-compression-methods)
    - [5.3.2 File-Based Compression Methods](#532-file-based-compression-methods)
    - [5.3.3 Write Method Compression Handling](#533-packagewrite-method)
  - [5.4 Compression and Signing Relationship](#54-compression-and-signing-relationship)
    - [5.4.1 Signing Compressed Packages](#541-signing-compressed-packages)
    - [5.4.2 Compressing Signed Packages](#542-compressing-signed-packages)
  - [5.5 Compression Strategy Selection](#55-compression-strategy-selection)
    - [5.5.1 Automatic Compression Detection](#551-automatic-compression-detection)
    - [5.5.2 Compression Workflow Options](#552-compression-workflow-options)
  - [5.6 Performance Considerations](#56-performance-considerations)
    - [5.6.1 Compressed Package Performance](#561-compressed-package-performance)
    - [5.6.2 Compression Decision Factors](#562-compression-decision-factors)
  - [5.7 Error Handling](#57-error-handling)
    - [5.7.1 Compression Errors](#571-compression-errors)
    - [5.7.2 Write Strategy Errors](#572-write-strategy-errors)
  - [5.8 Compression Use Cases](#58-compression-use-cases)
    - [5.8.1 When to Use Compressed Packages](#581-when-to-use-compressed-packages)
    - [5.8.2 When to Use Uncompressed Packages](#582-when-to-use-uncompressed-packages)

---

## 0. Overview

This document defines the package writing operations for the NovusPack system, including atomic writing, in-place updates, strategy selection, and handling of signed packages.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation

## 1. SafeWrite - Atomic Package Writing

**Purpose**: Provides atomic package writing with guaranteed consistency and rollback capability.

### 1.1 Package.SafeWrite Method

```go
// Returns *PackageError on failure
func (p *Package) SafeWrite(ctx context.Context, overwrite bool) error
```

### 1.2 SafeWrite Implementation Strategy

- **PackageInfo Synchronization**: Before writing, synchronizes package header flags and metadata from PackageInfo fields (PackageInfo is the source of truth during in-memory operations)
- **Metadata-Only Packages**: Packages with FileCount = 0 are valid and MUST have the metadata-only flag (Bit 7) set in the header
- **Temp File Creation**: Creates temporary file in same directory as target
- **Same-Directory Requirement**: `SafeWrite` MUST create the temporary file in the same directory as the target file to ensure atomic rename operations
- **Directory Writable Check**: If the target directory is not writable, `SafeWrite` must return an error rather than falling back to a system temp directory (which would break atomicity guarantees)
- **Streaming Write**: Streams data from source package or temp files for large content
- **Memory Management**: Uses in-memory data for small files, streaming for large files
- **Atomic Rename**: Atomically renames temp file to target path
- **Atomic Replace Guarantee**: `SafeWrite` MUST guarantee atomic replace on the same filesystem when overwriting an existing file
  - This means the target file is either completely replaced with the new content or remains unchanged (no partial writes)
  - This guarantee applies when both the temporary file and target file are on the same filesystem
- **Cross-Filesystem Restriction**: `SafeWrite` MUST NOT support cross-filesystem operations
  - If the temp file and target are on different filesystems, `SafeWrite` MUST return an error
  - This restriction is required to maintain atomicity guarantees, as cross-filesystem rename operations are not atomic
- **Rollback**: Automatically cleans up temp file on failure
- **Compression**: Compression is controlled by the package state (header and/or FileEntries), not passed as a parameter
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

### 2.1 Package.FastWrite Method

```go
// Returns *PackageError on failure
func (p *Package) FastWrite(ctx context.Context) error
```

### 2.2 FastWrite Implementation Strategy

- **PackageInfo Synchronization**: Before writing, synchronizes package header flags and metadata from PackageInfo fields (PackageInfo is the source of truth during in-memory operations)
- **Metadata-Only Packages**: Packages with FileCount = 0 are valid and MUST have the metadata-only flag (Bit 7) set in the header
- **Target Path Restriction**: `FastWrite` requires the target path to be the same as the currently opened package path, because it performs in-place updates on the existing file
  - If a different path is specified, `FastWrite` must return an error
- **Missing Target Behavior**: If the target file does not exist, `FastWrite` must return an error indicating the file must exist for in-place updates
- **COW Filesystem Detection**: `FastWrite` detects Copy-On-Write (COW) filesystems before executing writes
  - Detection is performed by checking filesystem type and characteristics
  - Common COW filesystems include Btrfs, ZFS, and APFS
  - Detection may use filesystem-specific features or heuristics to identify COW behavior
- **Automatic Backup on COW Systems**: On COW filesystems, `FastWrite` automatically creates a backup copy of the target file before executing writes
  - Backup is created in the same directory as the target file with a `.nvpk.backup` extension (for example, `package.nvpk.backup`)
  - Backup is created atomically before any in-place modifications begin
  - If `FastWrite` completes successfully, the backup file is automatically removed
  - If `FastWrite` fails or is interrupted, the backup file remains for recovery purposes
  - Backup creation failure on COW systems causes `FastWrite` to return an error before any modifications
- **Entry Comparison**: Compares existing vs new file entries
- **Change Detection**: Identifies modified, added, and unchanged entries
- **In-Place Updates**: Updates only changed entries in existing file
- **Append New Data**: Appends new entries to end of file
- **Metadata Updates**: Updates file index and offsets
- **Concurrent Writes**: Uses multiple goroutines to write different portions of the package concurrently, providing additional speedup for large packages or multiple file updates
  - Different file entries can be written in parallel when they don't overlap
  - Coordination ensures proper ordering and consistency
  - Theoretical speedup scales with the number of independent write operations
- **Compression**: Compression is controlled by the package state (header and/or FileEntries), not passed as a parameter

### 2.3 FastWrite Use Cases

- **Incremental Updates**: Adding/modifying individual files
- **Large Package Modifications**: When SafeWrite would be too slow
- **Frequent Updates**: When performance is critical
- **Existing Packages**: When target package already exists
- **Unsigned Packages Only**: Cannot be used with signed packages (SignatureOffset > 0)
- **Uncompressed Packages Only**: Cannot be used with compressed packages

### 2.4 FastWrite Performance Characteristics

- **Speed**: Much faster than SafeWrite for updates
  - Concurrent goroutine writes provide additional speedup for large packages or multiple file updates
  - Theoretical speedup scales with the number of independent write operations that can be performed in parallel
- **Memory**: Lower memory usage (only changed data in memory)
- **Disk I/O**: Minimal disk usage (only changed data written)
- **Concurrency**: Multiple goroutines write different portions of the package concurrently
  - Parallel writes for non-overlapping file entries
  - Coordination ensures proper ordering and consistency
- **Safety**: Good safety with partial update recovery
- **Crash Safety Risk**: `FastWrite` CAN corrupt the target file on crash or interruption
  - `FastWrite` operates by overwriting parts of an existing package in-place on disk, then updating footer information (index, metadata, etc.) afterwards
  - If the operation is interrupted during the write phase (before footer updates complete), the file will be in a corrupted state
  - **NOTE:**`SafeWrite` should be used for critical operations where data integrity is paramount

### 2.5 FastWrite Error Handling

- **Entry Validation**: Validates existing package before modification
- **COW Backup Creation**: On COW filesystems, returns error if backup file creation fails before any modifications
- **Partial Recovery**: Can recover from partial update failures using recovery files or COW backup files
- **Change Tracking**: Tracks what was successfully updated
- **Fallback**: Falls back to SafeWrite if FastWrite fails
- **Signed File Check**: Returns error if attempting to use FastWrite on signed packages
- **Compressed File Check**: Returns error if attempting to use FastWrite on compressed packages
- **Target Path Mismatch**: Returns error if target path differs from opened package path
- **Missing Target File**: Returns error if target file does not exist
- **Backup Cleanup**: On successful completion, automatically removes COW backup file if created

### 2.6 FastWrite Recovery Capabilities

FastWrite provides partial recovery capabilities for cases where an in-place update is interrupted and leaves the target file in a corrupted state.

FastWrite recovery is best-effort.
Recovery may be incomplete depending on when the interruption occurred, and which data was already overwritten.

TODO: This spec needs to be fleshed out more; right now it's still an initial draft.

#### 2.6.1 COW Filesystem Backup Recovery

On COW filesystems, `FastWrite` creates an automatic backup copy (`.nvpk.backup`) before executing writes.
If `FastWrite` fails or is interrupted, the backup file can be used to restore the original package state.
The backup file is automatically removed on successful completion, but remains available for manual recovery if the operation fails.
Recovery from COW backup is simpler and more reliable than recovery from partial writes, as the backup contains the complete original package state.

#### 2.6.2 Automatic Recovery File Creation

When `FastWrite` detects a failure or interruption, the API MUST attempt to dump relevant in-memory data to a recovery file (`filename.nvpk.recovery`) in the same directory as the original target file.
This recovery file creation must execute quickly and should resist cancellation attempts (context cancellation should be ignored during recovery file write to maximize recovery chances).
The recovery file only exists if the API detected the `FastWrite` failure and was able to continue executing long enough to write the recovery file.
On COW filesystems, both the backup file (if created) and the recovery file may be available for recovery.

#### 2.6.3 Recovery Cross-Reference

See [OpenBrokenPackage](api_basic_operations.md#12-openbrokenpackage-function) for repair workflow details.

### 2.7 Recovery File Format

The recovery file is a length-prefixed binary format intended for fast writes during failure handling.
All integer fields are encoded using little-endian, consistent with the main package file format.

#### 2.7.1 Recovery File Name

- `{original-package-name}.nvpk.recovery`

#### 2.7.2 Recovery File Structure

1. RecoveryFileHeader (fixed size)
2. PackageRecoveryState (length-prefixed)
3. FileEntryRecoveryState list (count + length-prefixed entries)
4. RecoveryMetadata (length-prefixed)

##### 2.7.2.1 RecoveryFileHeader Structure

```go
// RecoveryFileHeader contains header information for recovery files.
type RecoveryFileHeader struct {
    Magic          [4]byte  // "NVPR"
    Version        uint16   // Recovery format version (1)
    Flags          uint16   // Reserved for future use
    OriginalOffset uint64   // Last known good offset in the target file
    StateSize      uint32   // PackageRecoveryState size in bytes
    EntryCount     uint32   // Number of FileEntryRecoveryState entries
    MetadataSize   uint32   // RecoveryMetadata size in bytes
    Timestamp      uint64   // Unix nanoseconds
    Checksum       uint32   // CRC32 of the entire file excluding this field
    Reserved       [20]byte // Reserved for future use
}
```

##### 2.7.2.2 PackageRecoveryState (Length-Prefixed)

- PackageInfo snapshot.
- Target path.
- Write phase identifier.
- Last written offset.
- Staged file entries not yet persisted.
- List of package paths successfully written before failure.

##### 2.7.2.3 FileEntryRecoveryState List

For each FileEntry, store the FileEntry plus a status string:

- not_started
- in_progress
- completed
- failed

##### 2.7.2.4 RecoveryMetadata (Length-Prefixed)

- Error type and message.
- Stack trace (if available).
- Best-effort goroutine identifier (if available).
- Operation parameters.
- System state (disk space, timestamp, permissions, filesystem type).

##### 2.7.2.5 Checksum Calculation

The header checksum is CRC32 over the entire recovery file contents excluding the Checksum field itself.

#### 2.7.3 OpenBrokenPackage Integration

Recovery file consumption is handled by [OpenBrokenPackage](api_basic_operations.md#12-openbrokenpackage-function).

OpenBrokenPackage SHOULD:

- Prefer COW backup files when available.
- Use recovery file data to reconstruct in-memory state and attempt repair when no backup is available.
- Leave the recovery file in place if recovery fails.

#### 2.7.4 Recovery File Lifecycle

- The recovery file is written only when FastWrite detects a failure and has enough time to write it.
- The recovery file SHOULD be deleted only after a successful recovery and successful write of a repaired package.
- If recovery fails, the recovery file MUST remain for manual inspection and debugging.

## 3. Write Strategy Selection

This section describes write strategy selection for package operations.

### 3.1 Automatic Selection Logic

The `Write() error` method automatically selects the appropriate writing strategy:

- **New package**: If target file doesn't exist, uses SafeWrite for new package creation
- **Complete rewrite**: If package requires complete rewrite, uses SafeWrite
- **In-place updates**: Attempts FastWrite first for existing unsigned packages (when target path matches opened path)
- **Signed packages**: Refuses write operations unless package has been reconfigured to a new target path (which clears signatures)
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

### 4.2 Signed Package Writing Behavior

This section describes writing behavior for signed packages.

#### 4.2.1 Signed Package Protection

When attempting to write a signed package, write operations are refused by default.
To write a new package file derived from a signed package:

1. Reconfigure the Package to a new target path using [`SetTargetPath`](api_basic_operations.md#8-packagesettargetpath-method) (this automatically clears signature information from the in-memory Package)
2. Call `Write()` or `SafeWrite()` to write the new unsigned package file

`FastWrite()` MUST NOT be used for signed packages.
V1 enforces immutability based on signature presence (for example, `SignatureOffset > 0`) and does not validate signature contents.

#### 4.2.2 Compression Configuration

Compression is controlled by the package state (header and/or FileEntries), not passed as a parameter.
Compression settings are determined by the in-memory package configuration.

### 4.3 Writing Signed Package Content to New Path

When reconfiguring a signed package to a new target path using [`SetTargetPath`](api_basic_operations.md#8-packagesettargetpath-method):

1. **Signature Clearing**: `SetTargetPath` MUST automatically clear signature information from the in-memory Package when the new path differs from the current path
2. **New File Creation**: Creates a new, unsigned package file using SafeWrite (complete rewrite)
3. **Filename Requirement**: New filename must be different from the current signed file
4. **Signature Removal**: All signatures are stripped from the new file
5. **Content Preservation**: All package content (files, metadata, comments) is preserved
6. **Immutability Reset**: The new file can be modified normally (not immutable)
7. **Write Strategy**: Always uses SafeWrite since signed files cannot be modified in-place

**Important**: Signature clearing only occurs when the new path differs from the current path.
If `SetTargetPath` is called with the same path as the current path, signatures are NOT cleared.

### 4.4 Signed Package Writing Error Conditions

- **SignedFileError**: Returned when attempting to write to signed file without reconfiguring to a new path
- **SameFilenameError**: Returned when new target path matches current signed file path
- **ValidationError**: Returned if the signed file is corrupted or invalid
- **ErrTypeSecurity**: Returned when attempting to overwrite a signed package

### 4.5 Signed Package Use Cases

- **Development Workflow**: Reconfigure to new path to continue development on a previously signed package
- **Package Modification**: Make changes to a signed package while preserving content
- **Signature Clearing**: Remove signatures by writing an unsigned copy to a new path
- **Testing**: Create unsigned copies of signed packages for testing

### 4.6 Security Considerations

- **Explicit Intent**: Requires reconfiguring to a new path to prevent accidental signature clearing
- **Filename Protection**: Prevents accidental overwrite of signed files
- **Audit Trail**: Path reconfiguration and signature clearing operations should be logged for security auditing
- **Backup Recommendation**: Users should backup signed files before reconfiguring to a new path

## 5. Compressed Package Writing Operations

**Purpose**: Defines behavior for write operations on package-compressed packages to ensure proper handling of compression state.

### 5.1 Compressed Package Detection

When writing packages, the system must check for package compression:

- **PackageCompression Field**: Check Bit 15-8 in header flags for compression type
- **IsPackageCompressed**: Boolean flag indicating if package is compressed
- **Compression Type**: 0=none, 1=Zstd, 2=LZ4, 3=LZMA

### 5.2 Write Operations on Compressed Packages

This section describes write operations on compressed packages.

#### 5.2.1 SafeWrite with Compressed Packages

See [SafeWrite Method Signature](api_writing.md#11-packagesafewrite-method) for the complete method definition.

The following describes behavior specific to compressed packages:

##### 5.2.1.1 Behavior for Compressed Packages

- **Decompression Required**: Package must be decompressed before writing
- **Compression Preservation**: Original compression settings are preserved in header
- **Recompression**: Package is recompressed after writing if it was originally compressed
- **Memory Management**: Uses streaming for large compressed packages
- **Header/Comment/Signature Access**: Header, comment, and signatures remain uncompressed for direct access

#### 5.2.2 FastWrite with Compressed Packages

See [FastWrite Method Signature](api_writing.md#21-packagefastwrite-method) for the complete method definition.

The following describes behavior specific to compressed packages:

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

#### 5.3.3 Package.Write Method

```go
// Returns *PackageError on failure
func (p *Package) Write(ctx context.Context) error
```

- **Compression Control**: Compression is controlled by the package state (header and/or FileEntries), not passed as a parameter
- **Internal Process**: Uses internal compression methods before writing based on package state
- **Method Selection**: Uses SafeWrite for compressed packages, FastWrite for uncompressed packages

### 5.4 Compression and Signing Relationship

Signature generation and signature validation are deferred to v2.
V1 only enforces signed package immutability.

#### 5.4.1 Signing Compressed Packages

**Deferred to v2**: Signing compressed packages.

**V2 Supported Operation**: Compressed packages can be signed

- **Process**: Compress package first, then add signatures
- **Header Access**: Header remains uncompressed for signature validation
- **Comment Access**: Comment remains uncompressed for easy reading
- **Signature Access**: Signatures remain uncompressed for validation
- **Validation**: Signatures validate the compressed content

#### 5.4.2 Compressing Signed Packages

Signed packages cannot be compressed in place.
Compression changes package content and is disallowed by signed package immutability constraints.

To compress content from a signed package in v1:

1. Clear signatures by writing an unsigned, compressed copy to a new path.
2. Signing the compressed package is deferred to v2.

### 5.5 Compression Strategy Selection

This section describes compression strategy selection for write operations.

#### 5.5.1 Automatic Compression Detection

The `Write() error` method automatically handles compression:

- **Preserve State**: Maintains current compression state by default
- **Compressed Input**: If reading compressed package, writes compressed package
- **Uncompressed Input**: If reading uncompressed package, writes uncompressed package
- **New Package**: New packages are uncompressed by default
- **Signed Package Check**: Refuses compression if package is signed

#### 5.5.2 Compression Workflow Options

Multiple approaches for handling compression:

- **In-Memory Workflow**: Use `CompressPackage()`/`DecompressPackage()` then `Write()`
- **File-Based Workflow**: Use `CompressPackageFile()`/`DecompressPackageFile()` directly
- **Write with Compression**: Use `CompressPackage()` to set compression state, then `Write()` to write with compression
- **State Management**: In-memory methods update package state, file methods don't affect in-memory state
- **Signed Package Check**: All compression methods return error if package is signed

### 5.6 Performance Considerations

This section describes performance considerations for write operations.

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

This section describes error handling for write operations.

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

### 5.8 Compression Use Cases

This section describes use cases for compression in write operations.

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
