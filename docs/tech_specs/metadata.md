# NovusPack Technical Specifications - Metadata System

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Per-File Tags System Specification](#1-per-file-tags-system-specification)
  - [1.1 Tag Storage Format](#11-tag-storage-format)
    - [1.1.1 Tag Structure](#111-tag-structure)
  - [1.2 Tag Value Types](#12-tag-value-types)
    - [1.2.1 Basic Types](#121-basic-types)
    - [1.2.2 Structured Data](#122-structured-data)
    - [1.2.3 Identifier Value Types](#123-identifier-value-types)
    - [1.2.4 Time Value Types](#124-time-value-types)
    - [1.2.5 Network and Communication Value Types](#125-network-and-communication-value-types)
    - [1.2.6 File System](#126-file-system)
    - [1.2.7 Localization Value Types](#127-localization-value-types)
    - [1.2.8 NovusPack Special Files](#128-novuspack-special-files)
    - [1.2.9 Reserved Value Types](#129-reserved-value-types)
  - [1.3 PathMetadata System](#13-pathmetadata-system)
    - [1.3.1 PathMetadata File](#131-pathmetadata-file)
    - [1.3.2 PathMetadataEntry Structure](#132-pathmetadataentry-structure)
    - [1.3.3 Tag Inheritance Rules](#133-tag-inheritance-rules)
    - [1.3.4 Inheritance Examples](#134-inheritance-examples)
  - [1.4 Tag Validation](#14-tag-validation)
  - [1.5 Per-File Tags Usage Examples](#15-per-file-tags-usage-examples)
    - [1.5.1 Texture File Tagging](#151-texture-file-tagging)
- [2. Package Metadata File Specification](#2-package-metadata-file-specification)
  - [2.1 Metadata File Requirements](#21-metadata-file-requirements)
  - [2.2 YAML Schema Structure](#22-yaml-schema-structure)
    - [2.2.1 Package Metadata Schema v1.0](#221-package-metadata-schema-v10)
  - [2.3 Metadata File API](#23-metadata-file-api)
  - [2.4 Package Metadata Example](#24-package-metadata-example)
    - [2.4.1 Package Information Example](#241-package-information-example)

---

## 0. Overview

This document defines the per-file tags system and package metadata file specifications for the NovusPack system.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [Package File Format](package_file_format.md) - .nvpk format and FileEntry structure
- [File Types System](file_type_system.md) - Comprehensive file type system
- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Security and Encryption](security.md) - Comprehensive security architecture, encryption implementation, and digital signature requirements
- [File Validation](file_validation.md) - File validation and transparency requirements
- [Testing Requirements](testing.md) - Comprehensive testing requirements and validation
- [Package Metadata API](api_metadata.md) - Comment/AppID/VendorID methods and path metadata APIs

## 1. Per-File Tags System Specification

The per-file tags system provides extensible metadata for individual files within the package, supporting key-value pairs with type validation and inheritance.

### 1.1 Tag Storage Format

This section describes the tag storage format used in packages.

#### 1.1.1 Tag Structure

- **TagCount (2 bytes)**: Number of tags (0 if no tags)
- **Tags (variable)**: Array of tag entries
  - **Each tag entry**: `[KeyLength: 1 bytes][Key: UTF-8 string][ValueType: 1 byte][ValueLength: 2 bytes][Value: variable]`
  - **KeyLength (1 byte)**: Number of bytes in the key string
  - **Key (UTF-8 string)**: UTF-8 encoded string, no null termination
  - **ValueType (1 byte)**: Value type (0x00=string, 0x01=integer, 0x02=float, 0x03=boolean, 0x04=JSON, 0x05=YAML, 0x06=StringList, 0x07=UUID, 0x08=Hash, 0x09=Version, 0x0A=Timestamp, 0x0B=URL, 0x0C=Email, 0x0D=Path, 0x0E=MimeType, 0x0F=Language, 0x10=NovusPackMetadata)
  - **ValueLength (2 bytes)**: Number of bytes in the value
  - **Value (variable)**: UTF-8 encoded value based on ValueType

### 1.2 Tag Value Types

This section describes tag value types supported in the metadata system.

#### 1.2.1 Basic Types

- **0x00 - String**: UTF-8 encoded string value
- **0x01 - Integer**: 64-bit signed integer (stored as UTF-8 string representation)
- **0x02 - Float**: 64-bit floating point number (stored as UTF-8 string representation)
- **0x03 - Boolean**: "true" or "false" (stored as UTF-8 string)

#### 1.2.2 Structured Data

- **0x04 - JSON**: JSON-encoded object or array (stored as UTF-8 string)
- **0x05 - YAML**: YAML-encoded data (stored as UTF-8 string)
- **0x06 - StringList**: Comma-separated list of strings (stored as UTF-8 string)

#### 1.2.3 Identifier Value Types

- **0x07 - UUID**: UUID string (stored as UTF-8 string)
- **0x08 - Hash**: Hash/checksum string (stored as UTF-8 string)
- **0x09 - Version**: Semantic version string (stored as UTF-8 string)

#### 1.2.4 Time Value Types

- **0x0A - Timestamp**: ISO8601 timestamp (stored as UTF-8 string)

#### 1.2.5 Network and Communication Value Types

- **0x0B - URL**: URL string (stored as UTF-8 string)
- **0x0C - Email**: Email address (stored as UTF-8 string)

#### 1.2.6 File System

- **0x0D - Path**: File system path (stored as UTF-8 string)
- **0x0E - MimeType**: MIME type string (stored as UTF-8 string)

#### 1.2.7 Localization Value Types

- **0x0F - Language**: Language code (ISO 639-1) (stored as UTF-8 string)

#### 1.2.8 NovusPack Special Files

- **0x10 - NovusPackMetadata**: NovusPack special metadata file reference (stored as UTF-8 string)

#### 1.2.9 Reserved Value Types

- **0x11-0xFF**: Reserved for future value types

### 1.3 PathMetadata System

Cross-Reference: Operational APIs, structures, and methods for path metadata are defined in [Package Metadata API](api_metadata.md#8-pathmetadata-system).

Since NovusPack uses a flat file structure, path metadata is stored in special metadata files rather than implicit directory relationships.

#### 1.3.1 PathMetadata File

**File Type**: 65001 (Path metadata file - see [File Types System](file_type_system.md#339-special-file-types-65000-65535))
**File Name**: `__NVPK_PATH_65001__.nvpkpath` (case-sensitive)
**Content Format**: YAML syntax (stored as uncompressed or LZ4-compressed data)
**Purpose**: Defines path properties and inheritance rules

Cross-Reference: Storage format, CompressionType, and automatic decompression rules for special metadata files are defined in [Package Metadata API - Special Metadata File Management](api_metadata.md#83-special-metadata-file-management).

#### 1.3.2 PathMetadataEntry Structure

Each path metadata entry in the metadata file contains:

```yaml
paths:
  - path: "/assets/"                    # Path (directory paths must end with /, all paths have leading /)
    type: 1                             # PathMetadataTypeDirectory
    properties:
      category: "texture"               # Path-specific tags
      compression: "lossless"
      mipmaps: true
    inheritance:
      enabled: true                     # Whether this path provides inheritance
      priority: 1                       # Inheritance priority (higher = more specific)
    metadata:
      created: "2024-01-01T00:00:00Z"  # Path creation time
      modified: "2024-01-15T12:30:00Z" # Last modification time
      description: "Asset path"        # Human-readable description
```

#### 1.3.3 Tag Inheritance Rules

**Important**: Tag inheritance only works with `PathMetadataEntry` instances, not `FileEntry` instances directly.
This is because `FileEntry` can have multiple paths, while `PathMetadataEntry` represents a single path.
Inheritance is resolved per-path by walking up the path hierarchy for each `PathMetadataEntry`.

1. **Path-Based Inheritance**: Tags are inherited from path metadata entries in the path hierarchy
    - For a file at `/assets/textures/ui/button.png`, the associated `PathMetadataEntry` inherits from:
        - `/assets/textures/ui/` (if path metadata exists with inheritance enabled)
        - `/assets/textures/` (if path metadata exists with inheritance enabled)
        - `/assets/` (if path metadata exists with inheritance enabled)
        - `/` (root path metadata with inheritance enabled)

2. **Override Priority**: Child path tags override parent path tags
    - PathMetadataEntry tags (path-specific) have highest priority
    - Inherited path tags override based on inheritance priority
    - FileEntry tags are treated as if applied to each associated PathMetadataEntry (included in effective tags)
    - Root path tags have lowest priority

3. **Inheritance Resolution**: When multiple paths could provide tags:
    - Paths with exact path matches take priority
    - Paths with higher priority values override lower ones
    - If priorities are equal, more recently modified paths take priority

4. **Path Matching Rules**:
    - Directory paths must end with `/` in metadata
    - Path matching is case-sensitive
    - All paths are stored with forward slashes (`/`) as separators (Unix-style), regardless of source platform
    - All paths are stored with a leading `/` to indicate the package root (internal storage format)
    - Root path is represented as `/` in metadata files
    - **User Display**: When displaying paths to end users, the leading `/` MUST be stripped
      - Stored: `/assets/textures/ui/button.png` => Displayed: `assets/textures/ui/button.png`
    - When extracting or viewing on Windows, the leading `/` is stripped and paths are converted to Windows-style (backslashes) for file system operations

5. **Inheritance Scope**:
    - Only `PathMetadataEntry` instances with `Inheritance.Enabled = true` participate in inheritance
    - Inheritance walks up the `PathMetadataEntry.ParentPath` chain
    - Each `PathMetadataEntry` represents a single path, making inheritance unambiguous
    - `FileEntry` tags are treated as if applied to each associated `PathMetadataEntry` and are included in effective tags

#### 1.3.4 Inheritance Examples

This section provides examples of tag inheritance patterns.

##### 1.3.4.1 Example 1: Basic PathInheritance

```yaml
# Path metadata file (__NVPK_PATH_65001__.nvpkpath)
paths:
  - path: "/assets/"
    properties:
      - key: "category"
        value_type: "string"
        value: "texture"
      - key: "compression"
        value_type: "string"
        value: "lossless"
    inheritance:
      enabled: true
      priority: 1
  - path: "/assets/textures/"
    properties:
      - key: "format"
        value_type: "string"
        value: "png"
      - key: "mipmaps"
        value_type: "boolean"
        value: "true"
    inheritance:
      enabled: true
      priority: 2
  - path: "/assets/textures/ui/"
    properties:
      - key: "priority"
        value_type: "string"
        value: "high"
    inheritance:
      enabled: true
      priority: 3

# File: /assets/textures/ui/button.png
# PathMetadataEntry for this path inherits: category=texture, compression=lossless, format=png, mipmaps=true, priority=high
# Note: Inheritance is resolved per-path using PathMetadataEntry, not per-FileEntry
```

##### 1.3.4.2 Example 2: Priority-Based Override

```yaml
# Path metadata file
paths:
  - path: "/assets/"
    properties:
      - key: "category"
        value_type: "string"
        value: "texture"
      - key: "compression"
        value_type: "string"
        value: "lossless"
    inheritance:
      enabled: true
      priority: 1
  - path: "/assets/textures/"
    properties:
      - key: "category"
        value_type: "string"
        value: "image"
      - key: "format"
        value_type: "string"
        value: "png"
    inheritance:
      enabled: true
      priority: 2
  - path: "assets/textures/ui/"
    properties:
      - key: "category"
        value_type: "string"
        value: "ui"
      - key: "priority"
        value_type: "string"
        value: "high"
    inheritance:
      enabled: true
      priority: 3

# File: /assets/textures/ui/button.png
# Result: category=ui (highest priority), compression=lossless, format=png, priority=high
```

##### 1.3.4.3 Example 3: Inheritance Disabled

```yaml
# Path metadata file
paths:
  - path: "assets/"
    properties:
      - key: "category"
        value_type: "string"
        value: "texture"
      - key: "compression"
        value_type: "string"
        value: "lossless"
    inheritance:
      enabled: true
      priority: 1
  - path: "assets/temp/"
    properties:
      - key: "category"
        value_type: "string"
        value: "temporary"
      - key: "compression"
        value_type: "string"
        value: "lossy"
    inheritance:
      enabled: false  # This path doesn't provide inheritance
      priority: 2

# File: assets/temp/file.png
# Result: No inherited tags (temp/ inheritance disabled)
```

### 1.4 Tag Validation

- **Key Validation**: Keys must be valid UTF-8, non-empty, max 255 bytes
- **Value Validation**: Values must match their declared type
- **JSON Validation**: JSON values must be valid JSON syntax
- **Integer Validation**: Integer values must be valid 64-bit signed integers

### 1.5 Per-File Tags Usage Examples

This section provides usage examples for per-file tags.

#### 1.5.1 Texture File Tagging

- Set comprehensive tags on texture files including category, type, format, size, compression, priority, and descriptive tags
- Example: UI button texture with PNG format, 1024x1024 size, lossless compression, priority 5, and UI/button/interface tags

##### Audio File Tagging

- Set audio-specific tags including category, type, format, duration, loop settings, and volume
- Example: Ambient forest sound with WAV format, 120-second duration, loop enabled, 0.7 volume

##### Path Tagging

- Set tags on paths that are inherited by child files
- Example: Textures path with texture category, lossless compression, and mipmaps enabled

##### File Search by Tags

- Search for files by specific tag values using Package.FindEntriesByTag() or by iterating FileEntry tags
- Examples: Find all texture files by category, UI files by type, high-priority files by priority level
- Tag management is performed directly on FileEntry or PathMetadataEntry instances using AddTag(), SetTag(), GetTag(), etc.

---

## 2. Package Metadata File Specification

The package metadata file is a special file type (65000 - see [File Types System](file_type_system.md#339-special-file-types-65000-65535)) that contains structured YAML metadata about the package.

### 2.1 Metadata File Requirements

- **File Type**: Must be marked as Type 65000 (Package metadata file)
- **File Name**: Reserved name `__NVPK_META_65000__.nvpkmeta` (case-sensitive)
- **Content Format**: YAML syntax (stored as uncompressed or LZ4-compressed data)
- **Compression**: Optional (FastWrite SHOULD prefer uncompressed, but readers MUST support both uncompressed and LZ4-compressed data)
- **Automatic Decompression**: See [Package Metadata API - Special Metadata File Management](api_metadata.md#83-special-metadata-file-management)
- **Encryption**: Optional (can be encrypted like any other file)
- **Validation**: Must be valid YAML syntax

### 2.2 YAML Schema Structure

This section describes the YAML schema structure for package metadata.

#### 2.2.1. Package Metadata Schema V1.0

##### Package Information

- **name (string)**: Package name
- **version (string)**: Package version
- **description (string)**: Package description
- **author (string)**: Package author
- **license (string)**: Package license
- **created (ISO8601-timestamp)**: Creation timestamp
- **modified (ISO8601-timestamp)**: Last modification timestamp

##### Game-Specific Metadata

- **engine (string)**: Game engine (Unity, Unreal, etc.)
- **platform (array of strings)**: Target platforms
- **genre (string)**: Game genre
- **rating (string)**: Age rating
- **requirements**: System requirements
  - **min_ram (integer)**: Minimum RAM in MB
  - **min_storage (integer)**: Minimum storage in MB
  - **graphics (string)**: Graphics requirements
  - **os (array of strings)**: Supported operating systems

##### Asset Metadata

- **textures (integer)**: Number of texture files
- **sounds (integer)**: Number of sound files
- **models (integer)**: Number of 3D model files
- **scripts (integer)**: Number of script files
- **total_size (integer)**: Total asset size in bytes

##### Security Metadata

- **encryption_level (string)**: Encryption level used
- **signature_type (string, optional)**: Signature type used (deferred to v2)
- **security_scan (boolean)**: Whether security scan was performed
- **trusted_source (boolean)**: Whether from trusted source

##### Custom Metadata

- **custom (object)**: Extensible key-value pairs for additional metadata

### 2.3 Metadata File API

See the authoritative API definitions in [Package Metadata API](api_metadata.md) for package metadata operations and special metadata file management.

### 2.4 Package Metadata Example

**File**: `__NVPK_META_65000__.nvpkmeta`

#### 2.4.1 Package Information Example

- **name**: "MyAwesomeGame"
- **version**: "1.2.0"
- **description**: "An epic adventure game with stunning graphics"
- **author**: "GameStudio Inc."
- **license**: "Commercial"
- **created**: "2024-01-15T10:30:00Z"
- **modified**: "2024-01-20T14:45:00Z"

##### Game-Specific Metadata Example

- **engine**: "Unity 2023.3"
- **platform**: ["Windows", "macOS", "Linux"]
- **genre**: "Action-Adventure"
- **rating**: "T"
- **requirements**:
  - **min_ram**: 8192 MB
  - **min_storage**: 50000 MB
  - **graphics**: "DirectX 11 compatible"
  - **os**: ["Windows 10", "macOS 12", "Ubuntu 20.04"]

##### Asset Metadata Example

- **textures**: 1247 files
- **sounds**: 89 files
- **models**: 156 files
- **scripts**: 23 files
- **total_size**: 15728640 bytes

##### Security Metadata Example

- **encryption_level**: "ML-KEM Level 3"
- **signature_type**: "none"
- **security_scan**: true
- **trusted_source**: true

##### Custom Metadata Example

- **build_number**: 42
- **beta_version**: false
- **dlc_ready**: true
- **achievements**: 25
