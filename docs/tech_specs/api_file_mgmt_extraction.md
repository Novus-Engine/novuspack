# NovusPack Technical Specifications - File Extraction API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. ExtractPath Package Method](#1-extractpath-package-method)
  - [1.1 ExtractPath Purpose](#11-extractpath-purpose)
  - [1.2 `Package.ExtractPath` Signature](#12-packageextractpath-method)
  - [1.3 ExtractPath Parameters](#13-extractpath-parameters)
  - [1.4 ExtractPath Returns](#14-extractpath-returns)
  - [1.5 ExtractPath Behavior](#15-extractpath-behavior)
    - [1.5.1 ExtractPath Destination Resolution](#151-extractpath-destination-resolution)
    - [1.5.2 ExtractPath Directory Creation and Metadata Application](#152-extractpath-directory-creation-and-metadata-application)
    - [1.5.3 ExtractPath External Destination Handling](#153-extractpath-external-destination-handling)
    - [1.5.4 ExtractPath Concurrency and Ordering](#154-extractpath-concurrency-and-ordering)
  - [1.6 ExtractPath Path Handling](#16-extractpath-path-handling)
    - [1.6.1 ExtractPath Case Sensitivity Handling](#161-extractpath-case-sensitivity-handling)
  - [1.7 ExtractPath Error Conditions](#17-extractpath-error-conditions)
  - [1.8 ExtractPath Usage Notes](#18-extractpath-usage-notes)
- [2. ExtractPathOptions Struct](#2-extractpathoptions-struct)
  - [2.1 ExtractPathOptions Purpose](#21-extractpathoptions-purpose)
  - [2.2 ExtractPathOptions Destination Configuration](#22-extractpathoptions-destination-configuration)
  - [2.3 ExtractPathOptions Symlink Behavior](#23-extractpathoptions-symlink-behavior)
  - [2.4 ExtractPathOptions Security Limits](#24-extractpathoptions-security-limits)
  - [2.5 ExtractPathOptions Concurrency](#25-extractpathoptions-concurrency)
    - [2.5.1 Extraction ordering for key lifetime minimization](#251-extraction-ordering-for-key-lifetime-minimization)
- [3. Extraction Multi-Stage Pipeline Flow](#3-extraction-multi-stage-pipeline-flow)

---

## 0. Overview

This document specifies file extraction operations.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md)
- [Security Validation API](api_security.md)
- [Package Metadata API](api_metadata.md)

## 1. ExtractPath Package Method

This section describes the ExtractPath method for extracting files from packages.

### 1.1 ExtractPath Purpose

Extracts a file or directory subtree from the package to the filesystem.

### 1.2 Package.ExtractPath Method

```go
// ExtractPath extracts a file or directory subtree to the filesystem.
func (p *Package) ExtractPath(ctx context.Context, storedPath string, isWindows bool, opts *ExtractPathOptions) error
```

### 1.3 ExtractPath Parameters

- `ctx`: Context for cancellation and timeout handling
- `storedPath`: Stored package path of the file or directory to extract (stored format with forward slashes and a leading `/`)
- `isWindows`: Target filesystem semantics
- `opts`: Optional extraction configuration (may be nil)

### 1.4 ExtractPath Returns

`error`:

- Returns `*PackageError` on failure.

### 1.5 ExtractPath Behavior

ExtractPath locates the target path in the package path tree, plans an extraction, resolves destination paths, and writes the extracted content to the filesystem.

ExtractPath MUST treat `storedPath` as a package-internal path.
It MUST NOT treat `storedPath` as an OS filesystem path.
See [Package Path Semantics](api_core.md#2-package-path-semantics).

#### 1.5.1 ExtractPath Destination Resolution

Extraction destination resolution MUST follow this precedence order:

1. Call-time destination overrides from `ExtractPathOptions` (most specific wins: file override, then nearest parent directory override, then root override).
2. Stored destination overrides from `PathMetadataEntry` (most specific wins: file destination, then nearest parent directory destination, then root destination).
3. Default destination under session base (treating session base as the package root `/`).

##### 1.5.1.1 DestPathSpec Struct

```go
// DestPathSpec configures a destination path override.
//
// DestPath and DestPathWin may be absolute or relative.
// Relative destinations are resolved relative to the default extraction directory for the path.
type DestPathSpec struct {
    DestPath    Option[string]
    DestPathWin Option[string]
}
```

Destination overrides may resolve to either a directory destination or a file destination.
This depends on the path type:

- For file paths, destination overrides MUST be treated as file destinations.
- For directory paths (including `/`), destination overrides MUST be treated as directory destinations.

Stored path matching rules:

- When resolving destination overrides from options or from PathMetadataEntry, if a stored path string does not begin with `/`, the implementation MUST prefix `/` before matching.
- If multiple overrides apply, the most specific override MUST win as defined by the precedence order above.

The final extracted filesystem target for a file is computed as follows:

- If a file destination override applies, the resolved destination MUST be used as the full file path.
- Otherwise, the resolved destination MUST be treated as a directory, and the file path suffix MUST be joined to it.

The path suffix rules are:

- If a directory destination override applies to `/a/b/` and the file is `/a/b/c.txt`, the suffix is `c.txt`.
- If a directory destination override applies to `/a/` and the file is `/a/b/c.txt`, the suffix is `b/c.txt`.
- If no destination override applies, the suffix is the display path for the full stored path (for `/a/b/c.txt`, suffix is `a/b/c.txt`).

Destination overrides may be absolute or relative.
Relative destinations MAY include `.` and `..` segments.
Relative destinations MUST be resolved relative to the default extraction directory for the path, not relative to the current working directory.

The default extraction directory for a file path `/a/b/c.txt` is `Join(sessionBase, "a/b")`.
The default extraction directory for a directory path `/a/b/` is `Join(sessionBase, "a/b")`.

If the resolved destination path (directory or file path) is an empty string, extraction MUST fail with `ErrTypeValidation`.

Session base requirements:

- If session base is required to compute the destination (default-relative extraction or relative destination resolution), and it is not available (not set on the package and not provided in options), extraction MUST fail with `ErrTypeValidation`.
- If the caller provides an explicit absolute destination for a file at call-time, extraction MAY proceed without session base.
- If the destination is an absolute destination from `PathMetadataEntry` and the caller supplies an explicit allow override for extracting outside of root, extraction MAY proceed without session base.

ExtractPath MUST set the package runtime session base when `opts.SessionBase` is provided.

Stored destinations are sourced from `PathMetadataEntry.DestPath` and `PathMetadataEntry.DestPathWin`.
See [PathMetadataEntry Structure](api_metadata.md#812-pathmetadataentry-structure).

#### 1.5.2 ExtractPath Directory Creation and Metadata Application

ExtractPath MUST create destination directories before extracting files.
For directory extraction, this includes all directory entries in the subtree, including implied directories.

If a directory has a `PathMetadataEntry` with filesystem metadata, extraction SHOULD apply it as best-effort after creating the directory and before extracting files within it.

Directory metadata application includes permissions, timestamps, ownership, and other supported `PathMetadataEntry.FileSystem` fields.

#### 1.5.3 ExtractPath External Destination Handling

ExtractPath enforces a default destination safety policy.
By default, extraction MUST only write to:

- session base, or
- a subdirectory under session base that is part of the package directory set.

The package directory set MUST include:

- explicit directory paths in `PathMetadataEntry` instances, and
- implied directories derived from stored file paths.

If the resolved destination is outside the allowed destinations, ExtractPath MUST behave as follows:

- If `opts.AllowExternalDestinations` is true, extraction MAY proceed.
- Else if `opts.SkipDisallowedExternalDestinations` is true, the extracted target MUST be skipped.
  - If the skipped target is a directory, the entire subtree under that directory MUST be skipped.
- Else extraction MUST fail with `ErrTypeSecurity`.

ExtractPath MUST support ignoring stored destinations.
If `opts.IgnoreStoredDestPaths` is true, extraction MUST ignore `PathMetadataEntry.DestPath` and `PathMetadataEntry.DestPathWin` when resolving destinations.

#### 1.5.4 ExtractPath Concurrency and Ordering

For directory subtree extraction, ExtractPath SHOULD use concurrent extraction by default.
The concurrency configuration is controlled by `ExtractPathOptions`.

ExtractPath MUST follow the ordering rule for key lifetime minimization.
Encrypted files that will be decrypted during extraction MUST be prioritized ahead of unencrypted files.

### 1.6 ExtractPath Path Handling

The `storedPath` parameter MUST be treated as a stored package path.
If it does not begin with `/`, the implementation MUST prefix `/` before lookup and matching.

When extracting to the file system, file paths are converted as follows:

- On Windows targets (`isWindows == true`), forward slashes are converted to backslashes for filesystem operations.
- On Unix-like targets (`isWindows == false`), forward slashes are used.

See [Package Path Semantics](api_core.md#2-package-path-semantics) for complete path format details.

#### 1.6.1 ExtractPath Case Sensitivity Handling

By default, when extracting to a case-insensitive filesystem (Windows, default macOS), if the package contains paths that differ only in case, extraction MUST fail with an error.
See [Case Sensitivity](api_core.md#221-case-sensitivity) for complete case sensitivity behavior and error handling.

### 1.7 ExtractPath Error Conditions

- `ErrTypeValidation`: Package is not currently open
- `ErrTypeValidation`: Target path does not exist in the package
- `ErrTypeValidation`: Invalid or malformed stored path
- `ErrTypeValidation`: Session base is required but not set
- `ErrTypeValidation`: Destination path resolution results in an empty destination directory
- `ErrTypeValidation`: Case conflict detected (case-insensitive filesystem only)
- `ErrTypeIO`: I/O error during filesystem extraction
- `ErrTypeEncryption`: Failed to decrypt encrypted file during extraction
- `ErrTypeCompression`: Failed to decompress file content during extraction
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded
- `ErrTypeSecurity`: Disallowed external destination and skip not enabled

### 1.8 ExtractPath Usage Notes

ExtractPath writes extracted content to the filesystem.
To read file content into memory, use `ReadFile`.
See [ReadFile Method Contract](api_core.md#122-packagereaderreadfile-method).

## 2. ExtractPathOptions Struct

```go
// ExtractPathOptions configures filesystem extraction behavior.
type ExtractPathOptions struct {
    // Session base path for extraction destination resolution.
    //
    // If set, this value MUST be stored as the package runtime session base.
    // See [Package Session Base Management](api_basic_operations.md#19-package-session-base-management) for details.
    SessionBase Option[string]

    // Call-time destination overrides.
    //
    // Precedence:
    // - FileDestOverrides
    // - DirDestOverrides (nearest parent directory match)
    // - RootDestOverride
    RootDestOverride Option[DestPathSpec]
    DirDestOverrides  Option[map[string]DestPathSpec] // Keys are stored directory paths (paths ending with "/"); destination is a directory destination
    FileDestOverrides Option[map[string]DestPathSpec] // Keys are stored file paths (paths not ending with "/"); destination is a file destination

    // Stored destination behavior.
    IgnoreStoredDestPaths Option[bool] // Default: false

    // External destination policy.
    AllowExternalDestinations         Option[bool] // Default: false
    SkipDisallowedExternalDestinations Option[bool] // Default: false

    // Symlink extraction behavior
    PreserveSymlinks Option[bool]   // nil = platform default (true on Unix-like, false on Windows), Some(true/false) = explicit override

    // Symlink security
    MaxSymlinkDepth        Option[int]    // Default: 40 (prevent symlink bombs)
    ValidateSymlinkTargets Option[bool]   // Default: true
    RejectCircularSymlinks Option[bool]   // Default: true

    // Compression security (zip-bomb prevention)
    MaxCompressionRatio    Option[int]    // Default: 1000 (1000:1 ratio max)
    WarnOnSuspiciousRatios Option[bool]   // Default: true (log warnings for high ratios)

    // Resource limits (denial of service prevention)
    MaxTotalExtractedSize  Option[int64]  // Default: 0 (disabled, use disk space checks instead)
    MaxFileSize            Option[int64]  // Default: 0 (disabled, use disk space checks instead)
    MaxFileCount           Option[int64]  // Default: 1,000,000 (prevent "million files" attacks)
    MaxDirectoryDepth      Option[int]    // Default: 250 (prevent deep nesting)

    // Filesystem space validation (PRIMARY PROTECTION)
    RequiredSpaceMargin    Option[float64] // Default: 0.10 (require 10% extra space)
    MinimumFreeSpace       Option[int64]   // Default: 1GB (stop if space drops below this threshold)
    CheckSpaceInterval     Option[int64]   // Default: 5 seconds (check space every N seconds during sequential extraction)

    // Concurrent extraction settings
    MaxConcurrentExtractions Option[int]    // Default: number of CPU cores
    SpaceReservationMargin   Option[float64] // Default: 0.10 (10% safety margin per file reservation)
    SpaceCheckInterval       Option[int64]  // Default: 5 seconds (refresh available space cache interval in seconds for concurrent extraction)
    EnableConcurrentExtraction Option[bool] // Default: true

    // Security enforcement
    EnforceSecurityLimits  Option[bool]   // Default: true (cannot be disabled for untrusted packages)
    LogSecurityEvents      Option[bool]   // Default: true
}
```

### 2.1 ExtractPathOptions Purpose

Configures filesystem extraction behavior including destination resolution, external destination policy, symlink handling, security limits, and concurrency settings.

### 2.2 ExtractPathOptions Destination Configuration

ExtractPath destination configuration is controlled by:

- `SessionBase` (package-level session base path)
- call-time overrides (`RootDestOverride`, `DirDestOverrides`, `FileDestOverrides`)
- stored destinations in `PathMetadataEntry` (`DestPath`, `DestPathWin`), unless ignored
- external destination policy flags

### 2.3 ExtractPathOptions Symlink Behavior

- `PreserveSymlinks`: Controls whether symlinks are extracted as symlinks or file copies
  - `nil`: Platform default (true on Unix-like systems, false on Windows)
  - `Some(true)`: Extract symlinks as symlinks (requires appropriate privileges on Windows)
  - `Some(false)`: Extract symlinks as regular file copies

### 2.4 ExtractPathOptions Security Limits

Security limits prevent denial of service attacks:

- **Symlink Security**: Prevents symlink bombs via depth limits and circular link detection
- **Compression Security**: Prevents zip-bombs via compression ratio limits
- **Resource Limits**: Prevents resource exhaustion via file count and directory depth limits
- **Filesystem Space Validation**: Primary protection using real-time disk space checks

### 2.5 ExtractPathOptions Concurrency

Concurrent extraction settings enable parallel file extraction:

- **Worker Pool**: Configurable number of concurrent workers (default: CPU cores)
- **Thread-Safe Space Tracking**: Coordinated disk space monitoring across workers
- **Cancellation Coordination**: Shared cancellation mechanism for all workers

#### 2.5.1. Extraction Ordering for Key Lifetime Minimization

When extracting a whole package or directory subset to the filesystem, the implementation MUST minimize the time that decryption key material remains live in memory.

If an extraction plan includes a mixture of encrypted files (requiring decryption) and unencrypted files, encrypted files MUST be scheduled ahead of unencrypted files.
For concurrent extraction, the worker queue MUST prioritize encrypted files that will be decrypted during extraction.

After decrypting an individual file, the implementation MUST clear any file-specific key material and intermediate decrypt buffers as soon as they are no longer needed.
This includes derived keys, nonces, and temporary plaintext buffers used only for decryption.

## 3. Extraction Multi-Stage Pipeline Flow

This section is specified canonically in [File Transformation Pipelines](api_file_mgmt_transform_pipelines.md).
The detailed stage-by-stage walkthrough that previously lived here has been consolidated into the pipeline specification to keep a single source of truth.

For extraction-specific pipeline integration points, see:

- [ProcessingState transitions](api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
- [Pipeline execution model](api_file_mgmt_transform_pipelines.md#13-pipeline-execution-model)
- [Disk space management](api_file_mgmt_transform_pipelines.md#14-disk-space-management)
- [Intermediate stage cleanup](api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup)
- [Temporary file security](api_file_mgmt_transform_pipelines.md#17-temporary-file-security)
- [Error recovery and resume](api_file_mgmt_transform_pipelines.md#18-error-recovery-and-resume)
- [Example: Multi-Stage Extraction](api_file_mgmt_transform_pipelines.md#110-example-multi-stage-extraction)
