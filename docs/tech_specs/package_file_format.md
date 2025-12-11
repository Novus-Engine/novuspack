# NovusPack Technical Specifications - Package File Format

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1 `.npk` File Format Overview](#1-npk-file-format-overview)
  - [1.1 File Layout Order](#11-file-layout-order)
- [2 Package Header](#2-package-header)
  - [2.1 Header Structure](#21-header-structure)
  - [2.2 Package Version Fields Specification](#22-package-version-fields-specification)
  - [2.3 VendorID Field Specification](#23-vendorid-field-specification)
  - [2.4 AppID Field Specification](#24-appid-field-specification)
  - [2.5 Package Features Flags](#25-package-features-flags)
  - [2.6 ArchivePartInfo Field Specification](#26-archivepartinfo-field-specification)
  - [2.7 LocaleID Field Specification](#27-localeid-field-specification)
  - [2.8 Header Initialization](#28-header-initialization)
  - [2.9 Signed Package File Immutability and Incremental Signatures](#29-signed-package-file-immutability-and-incremental-signatures)
- [3 Package Compression](#3-package-compression)
  - [3.1 Compression Scope](#31-compression-scope)
  - [3.2 Compression Behavior](#32-compression-behavior)
- [4 File Entries and Data Section](#4-file-entries-and-data-section)
  - [4.1 File Entry Binary Format Specification](#41-file-entry-binary-format-specification)
  - [4.2 FileEntry Field Specifications](#42-fileentry-field-specifications)
- [5 File Index Section](#5-file-index-section)
- [6 Package Comment Section (Optional)](#6-package-comment-section-optional)
  - [6.1 Package Comment Format Specification](#61-package-comment-format-specification)
- [7 Digital Signatures Section (Optional)](#7-digital-signatures-section-optional)
  - [7.1 Signature Structure](#71-signature-structure)
  - [7.2 Signature Types](#72-signature-types)
  - [7.3 Signature Data Sizes](#73-signature-data-sizes)
  - [7.4 Cross-References](#74-cross-references)

---

## 0. Overview

This document defines the complete .npk package file format structure, including the package header, file entry binary format, and package comment specifications for the NovusPack system.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [Testing Requirements](testing.md) - Comprehensive testing requirements and validation
- [API Signatures Index](api_func_signatures_index.md) - Index of complete function references
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [File Types System](file_type_system.md) - Comprehensive file type system
- [Metadata System](metadata.md) - Package metadata and tags system

---

## 1 `.npk` File Format Overview

The .npk file format is a structured archive format designed for efficient storage, compression, encryption, and digital signing of files and directories.
It provides a modern alternative to traditional archive formats like ZIP, TAR, and 7Z with enhanced security and performance features.
The format consists of several sections arranged in a specific order to optimize both reading and writing operations.

### 1.1 File Layout Order

The .npk file structure follows this ordered layout:

1. **Package Header** (see [Package Header](#2-package-header)) - Fixed-size header with metadata and navigation information (immutable after first signature)
2. **File Entries and Data** (variable length) - Interleaved file entries and their data:
   - File Entry 1 (64-byte binary format + extended data)
   - File Data 1 (compressed/encrypted content)
   - File Entry 2 (64-byte binary format + extended data)
   - File Data 2 (compressed/encrypted content)
   - ... (repeat for each file)
3. **File Index** (variable length) - Index of all files with metadata and offsets
4. **Package Comment** (variable length, optional) - Human-readable package description
5. **Digital Signatures** (variable length, optional) - Multiple digital signatures for package integrity (appended incrementally)

## 2 Package Header

The package header provides comprehensive metadata and navigation information for the entire package.

**Note:** This is the authoritative definition of the package header size. All other references to header size should link to this section.

### 2.1 Header Structure

| Field              | Size    | Description                                                                   |
| ------------------ | ------- | ----------------------------------------------------------------------------- |
| Magic              | 4 bytes | Package identifier (0x4E56504B "NVPK")                                        |
| FormatVersion      | 4 bytes | Format version (current: 1)                                                   |
| Flags              | 4 bytes | Package-level features and options                                            |
| PackageDataVersion | 4 bytes | Package data version (increments on data changes)                             |
| MetadataVersion    | 4 bytes | Package metadata version (increments on metadata changes including comment)   |
| PackageCRC         | 4 bytes | CRC32 of package content (header excluded, signatures excluded, 0 if skipped) |
| CreatedTime        | 8 bytes | Package creation timestamp (Unix nanoseconds)                                 |
| ModifiedTime       | 8 bytes | Package modification timestamp (Unix nanoseconds)                             |
| LocaleID           | 4 bytes | Locale identifier for path encoding                                           |
| Reserved           | 4 bytes | Reserved for future use (must be 0)                                           |
| AppID              | 8 bytes | Application/game identifier (0 if not associated with specific app)           |
| VendorID           | 4 bytes | Storefront/platform identifier (0 if not associated with specific vendor)     |
| CreatorID          | 4 bytes | Creator identifier or reserved for future use                                 |
| IndexStart         | 8 bytes | Offset to file index from start of file                                       |
| IndexSize          | 8 bytes | Size of file index in bytes                                                   |
| ArchiveChainID     | 8 bytes | Archive chain identifier                                                      |
| ArchivePartInfo    | 4 bytes | Combined part number (2 bytes) + total parts (2 bytes)                        |
| CommentSize        | 4 bytes | Size of package comment in bytes (0 if no comment)                            |
| CommentStart       | 8 bytes | Offset to package comment from start of file                                  |
| SignatureOffset    | 8 bytes | Offset to signatures block from start of file                                 |

### 2.2 Package Version Fields Specification

#### 2.2.1 PackageDataVersion Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Tracks changes to package data content (files, file data, file index)
- **Increment**: Incremented whenever files are added, removed, or their data is modified
- **Initial Value**: 1 for new packages
- **Range**: 1 to 4,294,967,295 (4+ billion versions)
- **Usage**: Enables package data change detection and conflict resolution

#### 2.2.2 MetadataVersion Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Tracks changes to package metadata including the package comment
- **Increment**: Incremented whenever package metadata (including comment) is modified
- **Initial Value**: 1 for new packages
- **Range**: 1 to 4,294,967,295 (4+ billion versions)
- **Usage**: Enables package metadata change detection and version tracking

#### 2.2.3 PackageCRC Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Provides integrity validation for the entire package content
- **Calculation Scope**: All data from end of header through end of package file, excluding signatures
- **Algorithm**: CRC32 (same as file-level checksums for consistency)
- **Zero Value**: 0 indicates CRC calculation was skipped for performance
- **Inclusion**: File entries, file data, file index, and package comment
- **Exclusion**: Package header (to avoid circular dependency) and digital signatures
- **Usage**: Enables package integrity verification and corruption detection

#### 2.2.4 PackageCRC Calculation Process

The PackageCRC is calculated over the following data in order:

1. **File Entries and Data**: All file entries and their associated data (compressed/encrypted content)
2. **File Index**: Complete file index section
3. **Package Comment**: Package comment section (if present)

##### 2.2.4.1 Excluded from calculation

The following are excluded from the package-level CRC32 calculation

- Package header (to avoid circular dependency)
- Digital signatures (to allow signature addition without recalculating CRC)

##### 2.2.4.2 Performance Considerations

- CRC calculation can be computationally expensive for large packages
- Can be skipped during write operations for performance (set to 0)
- Can be calculated and updated later using API methods
- Recommended for production packages and integrity-critical scenarios

### 2.3 VendorID Field Specification

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Identify the storefront, platform, or vendor associated with the package
- **Default Value**: 0 (indicates no specific vendor association)
- **Platform Support**: Compatible with major storefront and platform systems

#### 2.3.1 VendorID Example Mappings

Examples using common software distribution methods.

- **Steam**: 0x53544541 (STEAM)
- **Epic Games Store**: 0x45504943 (EPIC)
- **GOG**: 0x474F4720 (GOG )
- **Itch.io**: 0x49544348 (ITCH)
- **Humble Bundle**: 0x48554D42 (HUMB)
- **Microsoft Store**: 0x4D494352 (MICR)
- **PlayStation Store**: 0x50534E59 (PSNY)
- **Xbox Store**: 0x58424F58 (XBOX)
- **Nintendo eShop**: 0x4E54444F (NTDO)
- **Unity Asset Store**: 0x554E4954 (UNIT)
- **Unreal Marketplace**: 0x554E5245 (UNRE)
- **GitHub**: 0x47495448 (GITH)
- **GitLab**: 0x4749544C (GITL)
- **No Vendor**: 0x00000000 (0 - default)

### 2.4 AppID Field Specification

- **Size**: 8 bytes (64-bit unsigned integer)
- **Purpose**: Associate package with specific application, game, or platform
- **Default Value**: 0 (indicates no specific application association)
- **Platform Support**: Compatible with various AppID systems:
  - **Steam**: 32-bit AppIDs (stored in lower 32 bits, upper 32 bits = 0)
  - **Itch.io**: Numeric game IDs
  - **Epic Games**: Store AppIDs
  - **GOG**: Game identifiers
  - **Unity Asset Store**: Package IDs
  - **Unreal Marketplace**: Product IDs
  - **Custom**: Any 64-bit identifier for proprietary systems

#### 2.4.1 AppID Examples

- **Steam AppID (CS:GO)**: 0x00000000000002DA (730)
- **Steam AppID (TF2)**: 0x00000000000001B8 (440)
- **Itch.io Game ID**: 0x0000000012345678 (custom format)
- **Epic Games AppID**: 0x00000000ABCDEF01 (custom format)
- **Generic Game ID**: 0x1234567890ABCDEF (custom 64-bit ID)
- **No Association**: 0x0000000000000000 (default)

#### 2.4.2 VendorID + AppID Combination Examples

- **Steam CS:GO**: VendorID=0x00000001, AppID=0x00000000000002DA
- **Epic Games Fortnite**: VendorID=0x00000002, AppID=0x00000000ABCDEF01
- **GOG Witcher 3**: VendorID=0x00000003, AppID=0x0000000012345678
- **Itch.io Indie Game**: VendorID=0x00000004, AppID=0x0000000056789ABC
- **Unity Asset**: VendorID=0x0000000A, AppID=0x00000000FEDCBA98

### 2.5 Package Features Flags

#### 2.5.1 Flags Field Encoding

- **Bit 31-24**: Reserved for future use (must be 0)
- **Bit 23-16**: Reserved for future use (must be 0)
- **Bit 15-8**: Package compression type (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
- **Bit 7-0**: Package features
  - **Bit 7**: Metadata-only package
  - **Bit 6**: Has special metadata files
  - **Bit 5**: Has per-file tags
  - **Bit 4**: Has package comment
  - **Bit 3**: Has extended attributes
  - **Bit 2**: Has encrypted files
  - **Bit 1**: Has compressed files
  - **Bit 0**: Has signatures (signed)

#### 2.5.2 Metadata-Related Flags

- **Bit 7**: Metadata-only package

  - **Purpose**: Indicates that the package contains only special metadata files and no regular content
  - **Usage**: Set to 1 if package contains only special metadata files (see [File Types System - Special File Types](file_type_system.md#special-file-types-65000-65535))

- **Bit 6**: Has special metadata files

  - **Purpose**: Indicates the presence of special metadata files in the package
  - **Usage**: Set to 1 if package contains special metadata files (see [File Types System - Special File Types](file_type_system.md#special-file-types-65000-65535))
  - **Related**: Corresponds to special metadata file detection

- **Bit 5**: Has per-file tags

  - **Purpose**: Indicates that files in the package have per-file tags
  - **Usage**: Set to 1 if any files have associated tags
  - **Related**: Corresponds to per-file tag system usage

- **Bit 4**: Has package comment

  - **Purpose**: Indicates that the package has a comment
  - **Usage**: Set to 1 if package contains a comment
  - **Related**: Corresponds to `HasComment` in PackageInfo

- **Bit 3**: Has extended attributes
  - **Purpose**: Indicates that files have extended attributes
  - **Usage**: Set to 1 if any files have extended attributes
  - **Future**: Reserved for extended file attributes

#### 2.5.3 Content-Related Flags

- **Bit 2**: Has encrypted files

  - **Purpose**: Indicates that the package contains encrypted files (per-file encryption)
  - **Usage**: Set to 1 if any files in the package are encrypted
  - **Related**: Corresponds to `HasEncryptedData` in PackageInfo
  - **Note**: Per-file encryption, not package-level encryption

- **Bit 1**: Has compressed files

  - **Purpose**: Indicates that the package contains compressed files
  - **Usage**: Set to 1 if any files in the package are compressed
  - **Related**: Corresponds to `HasCompressedData` in PackageInfo

- **Bit 0**: Has signatures (signed)
  - **Purpose**: Indicates that the package has digital signatures
  - **Usage**: Set to 1 if package contains any digital signatures
  - **Related**: Corresponds to `HasSignatures` in PackageInfo
  - **Note**: With incremental signatures, detailed signature data is stored in signature sections
  - **CRITICAL**: This bit must be set before adding the first signature to maintain signature integrity

#### 2.5.4 Package Compression Type

- **Bit 15-8**: Package compression type
  - **0**: No package compression
  - **1**: Zstd compression
  - **2**: LZ4 compression
  - **3**: LZMA compression
  - **4-255**: Reserved for future compression algorithms

### 2.6 ArchivePartInfo Field Specification

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Combined archive part information for split archives
- **Format**:
  - Bits 31-16: Archive part number (0-65535, 0 for single archive)
  - Bits 15-0: Total parts in archive (1-65535, 1 for single archive)
- **Default Value**: 0x00010001 (part 1 of 1 - single archive)
- **Usage**: Enables support for large archives split across multiple files

### 2.7 LocaleID Field Specification

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Locale identifier for path encoding across all files in the package
- **Format**: Standard locale ID format (e.g., 0x0409 for en-US, 0x0411 for ja-JP)
- **Default Value**: 0 (system default locale)
- **Scope**: Package-wide setting that applies to all file paths

### 2.8 Header Initialization

#### 2.8.1 Initial Package Creation

- **Magic**: Always set to 0x4E56504B
- **FormatVersion**: Set to 1 for current format
- **Flags**: Set based on package configuration (encryption, signing, compression)
- **PackageDataVersion**: Set to 1 for new packages, increments on data changes
- **MetadataVersion**: Set to 1 for new packages, increments on metadata changes
- **PackageCRC**: Set to 0 if skipped, or CRC32 of package content (calculated at write time)
- **CreatedTime**: Set to current timestamp when package is created (immutable)
- **ModifiedTime**: Set to current timestamp when package is created or modified
- **LocaleID**: Set to locale identifier for path encoding (0 for system default)
- **Reserved**: Set to 0 (reserved for future use)
- **AppID**: Set to application/game identifier or 0 for generic packages
- **VendorID**: Set to storefront/platform identifier or 0 for generic packages
- **CreatorID**: Set to 0 (reserved for future use)
- **IndexStart**: Set to offset to file index from start of file
- **IndexSize**: Set to size of file index in bytes
- **ArchiveChainID**: Set to unique identifier for archive chain (0 for single archive)
- **ArchivePartInfo**: Set to 0x00010001 for single archive, or part number + total parts for split archives
- **CommentSize**: Set to 0 if no comment, or size of comment including null terminator
- **CommentStart**: Set to 0 if no comment, or offset to comment from start of file
- **SignatureOffset**: Set to 0 if no signatures, or offset to signature index from start of file

### 2.9 Signed Package File Immutability and Incremental Signatures

The entire file becomes immutable after the first signature is added to prevent invalidation of existing signatures.

- The entire file becomes immutable after the first signature is added
- The SignatureOffset field will be set (defaults to 0 otherwise)
- The signature entry will be added at the end of the file
- Each subsequent signature is appended to the end of the file
- Each signature signs all content up to that point (including previous signatures), its metadata, and signature comment
- Signatures are appended sequentially - new signatures are added, never modified in place
- This ensures no existing signatures are invalidated when new ones are added

#### 2.9.1 File Immutability Enforcement

- Any write operation to an existing package must first check if `SignatureOffset` is non-zero
- If `SignatureOffset > 0`, the package is signed and all content modifications are prohibited
- Only reads and signature addition operations are allowed on signed packages
- This includes header, file entries, file data, file index, and package comment
- This prevents accidental invalidation of existing signatures

## 3 Package Compression

Package compression is a file format feature that compresses the package content while preserving the header, package comment, and signatures in an uncompressed state for direct access.

### 3.1 Compression Scope

#### 3.1.1 Compressed Content

- File entries (directory structure)
- File data (actual file contents)
- Package index

#### 3.1.2 Uncompressed Content

- Package header (see [Package Header](#2-package-header))
- Package comment
- Digital signatures

### 3.2 Compression Behavior

Package compression behavior is defined by the compression type specified in the header flags (bits 15-8). The compression process and constraints are detailed in the [Package Compression API](api_package_compression.md).

#### 3.2.1 Key Constraints

- Compressed packages can be signed, but signed packages cannot be compressed
- Package compression is applied after per-file compression/encryption operations
- Package decompression must occur before per-file decompression

For detailed compression methods, types, and implementation details, see the [Package Compression API](api_package_compression.md).

## 4 File Entries and Data Section

This section contains interleaved file entries and their data. Each file entry immediately precedes its related data, allowing for efficient streaming and processing.

- **File Entry Structure**: 64-byte binary format + extended data (paths, hashes, optional data)
- **File Data**: Compressed and/or encrypted file content immediately following each entry
- **Interleaved Layout**: Entry 1 => Data 1 => Entry 2 => Data 2 => ... => Entry N => Data N
- **Variable Length**: Based on content and processing applied

### 4.1 File Entry Binary Format Specification

The file entry binary format consists of a fixed-size header followed by optional extended data. The format version is determined by the package header, not individual file entries.

#### 4.1.1 File Entry Static Section Field Encoding

| Field              | Size    | Description                                                      |
| ------------------ | ------- | ---------------------------------------------------------------- |
| FileID             | 8 bytes | Unique file identifier (64-bit unsigned integer)                 |
| OriginalSize       | 8 bytes | Original file size before processing                             |
| StoredSize         | 8 bytes | Final file size after compression/encryption                     |
| RawChecksum        | 4 bytes | CRC32 of raw file content                                        |
| StoredChecksum     | 4 bytes | CRC32 of processed file content                                  |
| FileVersion        | 4 bytes | File data version (increments on data changes)                   |
| MetadataVersion    | 4 bytes | File metadata version (increments on metadata changes)           |
| PathCount          | 2 bytes | Total number of paths (1 for single path, 2+ for multiple paths) |
| Type               | 2 bytes | File type identifier                                             |
| CompressionType    | 1 byte  | Compression algorithm identifier                                 |
| CompressionLevel   | 1 byte  | Compression level (0-9, 0=default)                               |
| EncryptionType     | 1 byte  | Encryption algorithm identifier                                  |
| HashCount          | 1 byte  | Number of hash entries (0 if no hashes)                          |
| HashDataOffset     | 4 bytes | Offset to hash data from start of variable-length data           |
| HashDataLen        | 2 bytes | Total length of all hash data in bytes (0 if no hashes)          |
| OptionalDataLen    | 2 bytes | Total length of optional data in bytes (0 if no optional data)   |
| OptionalDataOffset | 4 bytes | Offset to optional data from start of variable-length data       |
| Reserved           | 4 bytes | Reserved for future use (must be 0)                              |

##### 4.1.1.1 FileID Field Specification

- **FileID**: 8 bytes (64-bit unsigned integer) - Unique file identifier
- **Purpose**: Provides a stable, unique identifier for each file entry within the package
- **Uniqueness**: Must be unique across all file entries in the package
- **Generation**: Assigned sequentially during file addition (1, 2, 3, ...)
- **Persistence**: FileID remains constant for the lifetime of the file entry
- **Future-Proofing**: 64-bit range supports up to 18,446,744,073,709,551,615 files
- **Usage**: Enables efficient file tracking, referencing, and API operations
- **Zero Value**: FileID 0 is reserved and must not be used for actual files

##### 4.1.1.2 File Version Fields Specification

- **FileVersion**: 4 bytes (32-bit unsigned integer) - File data version

  - **Purpose**: Tracks changes to file content/data
  - **Increment**: Incremented whenever file data is modified
  - **Initial Value**: 1 for new files
  - **Range**: 1 to 4,294,967,295 (4+ billion versions)
  - **Usage**: Enables change detection and incremental operations

- **MetadataVersion**: 4 bytes (32-bit unsigned integer) - File metadata version
  - **Purpose**: Tracks changes to file metadata (paths, tags, compression, encryption, etc.)
  - **Increment**: Incremented whenever file metadata is modified
  - **Initial Value**: 1 for new files
  - **Range**: 1 to 4,294,967,295 (4+ billion versions)
  - **Usage**: Enables file metadata change detection and conflict resolution

##### 4.1.1.3 Compression and Encryption Types

- **CompressionType**: 1 byte - Direct compression algorithm identifier

  - 0: No compression
  - 1: Zstd compression
  - 2: LZ4 compression
  - 3: LZMA compression
  - 4-255: Reserved for future algorithms

- **EncryptionType**: 1 byte - Direct encryption algorithm identifier

  - 0x00: No encryption
  - 0x01: AES-256-GCM encryption
  - 0x02: Quantum-safe encryption (ML-KEM + ML-DSA)
  - 0x03-0xFF: Reserved for future algorithms

- **Type**: 2 bytes - File type identifier (see [File Types System](file_type_system.md) for detailed file type system)

#### 4.1.2 File Entry Structure Requirements

The file entry structure supports unique file identification, version tracking, multiple paths pointing to the same content, hash-based content identification, and comprehensive security metadata.

##### 4.1.2.1 Unique File Identification

Each file entry includes a unique 64-bit FileID that provides stable identification across package operations. The FileID enables efficient file tracking, API operations, and future extensibility without relying on path-based identification.

##### 4.1.2.2 File Version Tracking

Each file entry includes two version fields that track changes independently:

- **FileVersion**: Tracks changes to file content/data
- **MetadataVersion**: Tracks changes to file metadata (paths, tags, compression, encryption, etc.)

This dual versioning enables granular change detection and supports efficient incremental operations.
Note that package-level metadata (including the package comment) is tracked by the package header's `MetadataVersion` field, not individual file metadata versions.

##### 4.1.2.3 Multiple Path Support with Per-Path Metadata

Each file entry can have multiple paths pointing to the same content, with each path having its own metadata (permissions, timestamps, etc.). This enables efficient storage of hard links and symbolic links while maintaining individual path attributes.

##### 4.1.2.4 Hash-based Content Identification

File entries include multiple hash types for different purposes:

- Content hashes for deduplication
- Integrity hashes for verification
- Fast lookup hashes for quick identification

##### 4.1.2.5 Security Metadata

Each file entry includes encryption and compression metadata, allowing per-file security and optimization settings.

#### 4.1.3 Fixed Structure (64 bytes, optimized for 8-byte alignment)

The fixed structure is optimized for 8-byte alignment to minimize padding and improve performance on modern systems.

##### 4.1.3.1 Field Ordering

Fields are ordered by size (largest to smallest) to minimize padding:

1. 8-byte fields (FileID, OriginalSize, Size)
2. 4-byte fields (RawChecksum, StoredChecksum, FileVersion, MetadataVersion, HashDataOffset, OptionalDataOffset, Reserved)
3. 2-byte fields (PathCount, Type, HashDataLen, OptionalDataLen)
4. 1-byte fields (CompressionType, CompressionLevel, EncryptionType, HashCount)

#### 4.1.4 Variable-Length Data (follows fixed structure)

- **Primary path:** The main path for the file entry (stored in the Name field)
- **Additional paths:** Secondary paths that point to the same content (stored in Paths metadata)
- **Path metadata:** Array of additional paths stored as part of the file entry metadata
- **Per-path metadata:** Each path can have its own mode, UID, GID, and timestamps
- **Path validation:** All paths must be normalized and valid according to path validation rules
- **Content consistency:** All paths must resolve to identical content when extracted
- **Metadata flexibility:** Different paths can have different permissions and ownership while sharing content

##### 4.1.4.1 Variable-Length Data Order

The variable-length data section follows this order:

1. Path entries (if PathCount > 0, at offset 0)
2. Hash data (if HashCount > 0, at HashDataOffset)
3. Optional data (if OptionalDataLen > 0, at OptionalDataOffset)

##### 4.1.4.2 Path Entries

All paths: PathCount × path entries with metadata

Each path entry: `[PathLength: 2 bytes][Path: UTF-8][Mode: 4 bytes][UserID: 4 bytes][GroupID: 4 bytes][ModTime: 8 bytes][CreateTime: 8 bytes][AccessTime: 8 bytes]`

- **PathLength**: 2 bytes - Length of path in bytes (UTF-8 encoded)
- **Path**: UTF-8 string - File path (not null-terminated)
- **Mode**: 4 bytes - File permissions and type (Unix-style)
- **UserID**: 4 bytes - User ID (Unix-style)
- **GroupID**: 4 bytes - Group ID (Unix-style)
- **ModTime**: 8 bytes - Modification time (Unix nanoseconds)
- **CreateTime**: 8 bytes - Creation time (Unix nanoseconds)
- **AccessTime**: 8 bytes - Access time (Unix nanoseconds)

- **Primary path**: First path entry (index 0) is the primary path
- **Additional paths**: Secondary paths pointing to the same content with their own metadata

##### 4.1.4.3 Hash Data

Multiple hash entries (if HashCount > 0, located at HashDataOffset from start of variable-length data)

- **HashCount**: 1 byte - Number of hash entries (0 if no hashes)
- **Hash Entries**: Array of hash entries

Each hash entry: `[HashType: 1 byte][HashPurpose: 1 byte][HashLength: 2 bytes][HashData: variable]`

- **HashType**: 1 byte - Hash algorithm type (0x00=SHA-256, 0x01=SHA-512, 0x02=BLAKE3, 0x03=XXH3, 0x04-0xFF=reserved)
- **HashPurpose**: 1 byte - Hash purpose (0x00=content verification, 0x01=deduplication, 0x02=integrity, 0x03-0xFF=reserved)
- **HashLength**: 2 bytes - Length of hash data in bytes
- **HashData**: Variable-length hash data based on HashType

##### 4.1.4.4 Optional Data

Rarely-used file attributes (if OptionalDataLen > 0, located at OptionalDataOffset from start of variable-length data)

- **OptionalDataCount**: 2 bytes - Number of optional data entries (0 if no optional data)
- **Optional Data Entries**: Array of optional attribute entries

Each optional data entry: `[DataType: 1 byte][DataLength: 2 bytes][Data: variable]`

- **DataType**: 1 byte - Optional data type identifier
  - 0x00: TagsData (variable) - Per-file tags data
  - 0x01: PathEncoding (1 byte) - Path encoding type for this file
  - 0x02: PathFlags (1 byte) - Path handling flags for this file
  - 0x03: CompressionDictionaryID (4 bytes) - Dictionary identifier for solid compression
  - 0x04: SolidGroupID (4 bytes) - Solid compression group identifier
  - 0x05: FileSystemFlags (2 bytes) - File system specific flags
  - 0x06: WindowsAttributes (4 bytes) - Windows file attributes
  - 0x07: ExtendedAttributes (variable) - Unix extended attributes
  - 0x08: ACLData (variable) - Access Control List data
  - 0x09-0xFF: Reserved for future optional data types
- **DataLength**: 2 bytes - Length of optional data in bytes
- **Data**: Variable-length optional data (type determined by DataType)

#### 4.1.5 Hash Algorithm Support

- **HashType**: 1 byte - Hash algorithm identifier

  - 0x00: SHA-256 (32 bytes) - Standard cryptographic hash
  - 0x01: SHA-512 (64 bytes) - Stronger cryptographic hash
  - 0x02: BLAKE3 (32 bytes) - Fast cryptographic hash
  - 0x03: XXH3 (8 bytes) - Ultra-fast non-cryptographic hash
  - 0x04: BLAKE2b (64 bytes) - Cryptographic hash with configurable output
  - 0x05: BLAKE2s (32 bytes) - Cryptographic hash optimized for 32-bit systems
  - 0x06: SHA-3-256 (32 bytes) - SHA-3 family hash
  - 0x07: SHA-3-512 (64 bytes) - SHA-3 family hash
  - 0x08: CRC32 (4 bytes) - Fast checksum for error detection
  - 0x09: CRC64 (8 bytes) - Stronger checksum for error detection
  - 0x0A-0xFF: Reserved for future hash algorithms

- **HashPurpose**: 1 byte - Hash purpose identifier
  - 0x00: Content verification - Verify file content integrity
  - 0x01: Deduplication - Identify duplicate content
  - 0x02: Integrity check - General integrity verification
  - 0x03: Fast lookup - Quick content identification
  - 0x04: Error detection - Detect data corruption
  - 0x05-0xFF: Reserved for future purposes

### 4.2 FileEntry Field Specifications

#### 4.2.1 HashCount Field

- **Size**: 1 byte (8-bit unsigned integer)
- **Purpose**: Number of hash entries for this file
- **Range**: 0-255 hash entries
- **Default Value**: 0 (no hashes)

#### 4.2.2 HashDataLen Field

- **Size**: 2 bytes (16-bit unsigned integer)
- **Purpose**: Total length of all hash data in bytes
- **Range**: 0-65535 bytes
- **Default Value**: 0 (no hash data)

#### 4.2.3 HashDataOffset Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Offset to hash data from start of variable-length data section
- **Range**: 0-4294967295 bytes
- **Default Value**: 0 (no hash data)

#### 4.2.4 OptionalDataLen Field

- **Size**: 2 bytes (16-bit unsigned integer)
- **Purpose**: Total length of optional data in bytes
- **Range**: 0-65535 bytes
- **Default Value**: 0 (no optional data)

#### 4.2.5 OptionalDataOffset Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Offset to optional data from start of variable-length data
- **Range**: 0-4294967295 bytes
- **Default Value**: 0 (no optional data)

---

## 5 File Index Section

- FileIndexBinary: 16 bytes + entry references
- Entry references: 16 bytes per entry
- Contains metadata and offsets for all files in the package

## 6 Package Comment Section (Optional)

- UTF-8 string + null terminator
- Human-readable description of package contents
- Variable length, only present if CommentSize > 0

### 6.1 Package Comment Format Specification

The package comment is an optional, variable-length field that provides human-readable metadata about the package contents, purpose, or other descriptive information.

#### 6.1.1 Package Comment Structure

- **CommentLength (4 bytes)**: Length of comment in bytes including null terminator (0 if no comment)
- **Comment (variable)**: UTF-8 encoded string + null terminator
- **Reserved (3 bytes)**: Reserved for future use (must be 0)

##### 6.1.1.1 Field Specifications

- **CommentLength**: 4 bytes, little-endian unsigned integer

  - Value of 0 indicates no comment is present
  - Maximum length: 1,048,575 bytes (1MB - 1 byte) to prevent abuse
  - Must match the actual length of the Comment field including null terminator

- **Comment**: Variable-length UTF-8 string with null termination

  - Must be null-terminated (ends with 0x00)
  - Must be valid UTF-8 encoding
  - Can contain newlines, tabs, and other whitespace characters
  - Should be human-readable text describing the package
  - Cannot contain embedded null characters (except at the end)

- **Reserved**: 3 bytes, must be set to 0
  - Reserved for future extensions
  - Must be initialized to 0 when writing
  - Should be ignored when reading

##### 6.1.1.2 Implementation Requirements

- **Write behavior**: If no comment is provided, write CommentLength as 0 and skip Comment field
- **Read behavior**: If CommentLength is 0, skip reading Comment field
- **Validation**: Verify CommentLength matches actual Comment field size including null terminator
- **Error handling**: Return error for invalid UTF-8 encoding, length mismatches, or missing null terminator
- **Null termination**: Always append null byte (0x00) when writing comments

## 7 Digital Signatures Section (Optional)

Digital signatures provide package integrity and authenticity verification.
With incremental signing, signatures are appended sequentially without requiring a separate index.

**Note:** This section defines the binary format for signatures.
For signature implementation details, see [Digital Signature API](api_signatures.md).

### 7.1 Signature Structure

Each signature consists of a metadata header, optional comment, and signature data.
The signature validates all content up to its creation point, including its own metadata and comment.

| Field              | Size     | Description                                                          |
| ------------------ | -------- | -------------------------------------------------------------------- |
| SignatureType      | 4 bytes  | Signature type (0x01=ML-DSA, 0x02=SLH-DSA, 0x03=PGP, 0x04=X.509)     |
| SignatureSize      | 4 bytes  | Size of this signature data in bytes                                 |
| SignatureFlags     | 4 bytes  | Signature-specific flags                                             |
| SignatureTimestamp | 4 bytes  | When this signature was created (Unix nanoseconds)                   |
| CommentLength      | 2 bytes  | Length of signature comment in bytes (0 if no comment)               |
| SignatureComment   | Variable | Human-readable comment about this signature (UTF-8, null-terminated) |
| SignatureData      | Variable | Raw signature data                                                   |

### 7.2 Signature Types

#### 7.2.1 SignatureType Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Identify the signature algorithm used
- **Values**:
  - 0x01: ML-DSA (Module-Lattice Digital Signature Algorithm)
  - 0x02: SLH-DSA (Stateless Hash-based Digital Signature Algorithm)
  - 0x03: PGP (Pretty Good Privacy)
  - 0x04: X.509 (X.509 Certificate-based signature)
  - 0x05-0xFFFFFFFF: Reserved for future signature types

#### 7.2.2 SignatureFlags Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Signature-specific metadata and options
- **Bit encoding**:
  - Bit 31-16: Reserved for future use (must be 0)
  - Bit 15-8: Signature features (bit 7=has timestamp, bit 6=has metadata, bit 5=has chain validation, bit 4=has revocation, bit 3=has expiration, bit 2-0=reserved)
  - Bit 7-0: Signature status (bit 7=valid, bit 6=verified, bit 5=trusted, bit 4-0=reserved)

#### 7.2.3 SignatureTimestamp Field

- **Size**: 4 bytes (32-bit unsigned integer)
- **Purpose**: Timestamp when the signature was created
- **Format**: Unix timestamp in nanoseconds
- **Range**: 0-4294967295 (Unix nanoseconds)

#### 7.2.4 CommentLength Field

- **Size**: 2 bytes (16-bit unsigned integer)
- **Purpose**: Length of the signature comment in bytes
- **Range**: 0-65535 bytes
- **Default Value**: 0 (no comment)
- **Security Note**: Comments are included in signature validation along with all signature metadata, so they cannot be modified without invalidating the signature

### 7.3 Signature Data Sizes

- **ML-DSA**: ~2,420-4,595 bytes (depending on security level)
- **SLH-DSA**: ~7,856-17,088 bytes (depending on security level)
- **PGP**: Variable size (typically 256-512 bytes)
- **X.509**: Variable size (typically 256-4096 bytes)

### 7.4 Cross-References

For signature implementation details, see:

- [Digital Signature API](api_signatures.md) - Complete signature management API
- [Incremental Signing Process](api_signatures.md#12-incremental-signing-implementation) - Implementation details
- [Immutability Check](api_signatures.md#13-immutability-check) - Signature immutability requirements
