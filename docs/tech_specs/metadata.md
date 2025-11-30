# NovusPack Technical Specifications - Metadata System

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1 Per-File Tags System Specification](#1-per-file-tags-system-specification)
  - [1.1 Tag Storage Format](#11-tag-storage-format)
  - [1.2 Tag Value Types](#12-tag-value-types)
  - [1.3 Directory Metadata System](#13-directory-metadata-system)
  - [1.4 Tag Validation](#14-tag-validation)
  - [1.5 Per-File Tags Usage Examples](#15-per-file-tags-usage-examples)
- [2 Package Metadata File Specification](#2-package-metadata-file-specification)
  - [2.1 Metadata File Requirements](#21-metadata-file-requirements)
  - [2.2 YAML Schema Structure](#22-yaml-schema-structure)
  - [2.3 Metadata File API](#23-metadata-file-api)
  - [2.4 Package Metadata Example](#24-package-metadata-example)

---

## 0. Overview

This document defines the per-file tags system and package metadata file specifications for the NovusPack system.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [Package File Format](package_file_format.md) - .npk format and file entry structure
- [File Types System](file_type_system.md) - Comprehensive file type system
- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Security and Encryption](security.md) - Comprehensive security architecture, encryption implementation, and digital signature requirements
- [File Validation](file_validation.md) - File validation and transparency requirements
- [Testing Requirements](testing.md) - Comprehensive testing requirements and validation
- [Package Metadata API](api_metadata.md) - Comment/AppID/VendorID methods and directory metadata APIs

## 1 Per-File Tags System Specification

The per-file tags system provides extensible metadata for individual files within the package, supporting key-value pairs with type validation and inheritance.

### 1.1 Tag Storage Format

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

#### 1.2.1 Basic Types

- **0x00 - String**: UTF-8 encoded string value
- **0x01 - Integer**: 64-bit signed integer (stored as UTF-8 string representation)
- **0x02 - Float**: 64-bit floating point number (stored as UTF-8 string representation)
- **0x03 - Boolean**: "true" or "false" (stored as UTF-8 string)

#### 1.2.2 Structured Data

- **0x04 - JSON**: JSON-encoded object or array (stored as UTF-8 string)
- **0x05 - YAML**: YAML-encoded data (stored as UTF-8 string)
- **0x06 - StringList**: Comma-separated list of strings (stored as UTF-8 string)

#### 1.2.3 Identifiers

- **0x07 - UUID**: UUID string (stored as UTF-8 string)
- **0x08 - Hash**: Hash/checksum string (stored as UTF-8 string)
- **0x09 - Version**: Semantic version string (stored as UTF-8 string)

#### 1.2.4 Time

- **0x0A - Timestamp**: ISO8601 timestamp (stored as UTF-8 string)

#### 1.2.5 Network/Communication

- **0x0B - URL**: URL string (stored as UTF-8 string)
- **0x0C - Email**: Email address (stored as UTF-8 string)

#### 1.2.6 File System

- **0x0D - Path**: File system path (stored as UTF-8 string)
- **0x0E - MimeType**: MIME type string (stored as UTF-8 string)

#### 1.2.7 Localization

- **0x0F - Language**: Language code (ISO 639-1) (stored as UTF-8 string)

#### 1.2.8 NovusPack Special Files

- **0x10 - NovusPackMetadata**: NovusPack special metadata file reference (stored as UTF-8 string)

#### 1.2.9 Reserved

- **0x11-0xFF**: Reserved for future value types

### 1.3 Directory Metadata System

Cross-Reference: Operational APIs, structures, and methods for directory metadata are defined in [Package Metadata API](api_metadata.md#8-directory-metadata-system).

Since NovusPack uses a flat file structure, directory metadata is stored in special metadata files rather than implicit directory relationships.

#### 1.3.1 Directory Metadata File

**File Type**: 65001 (Directory metadata file - see [File Types System](file_type_system.md#special-files-65000-65535))
**File Name**: `__NPK_DIR_65001__.npkdir` (case-sensitive)
**Content Format**: YAML syntax (uncompressed for FastWrite compatibility)
**Purpose**: Defines directory properties and inheritance rules

#### 1.3.2 Directory Entry Structure

Each directory entry in the metadata file contains:

```yaml
directories:
  - path: "/assets/"                    # Directory path (must end with /)
    properties:
      category: "texture"               # Directory-specific tags
      compression: "lossless"
      mipmaps: true
    inheritance:
      enabled: true                     # Whether this directory provides inheritance
      priority: 1                       # Inheritance priority (higher = more specific)
    metadata:
      created: "2024-01-01T00:00:00Z"  # Directory creation time
      modified: "2024-01-15T12:30:00Z" # Last modification time
      description: "Asset directory"    # Human-readable description
```

#### 1.3.3 Tag Inheritance Rules

1. **Directory-Based Inheritance**: Tags are inherited from directory metadata files
    - File `/assets/textures/ui/button.png` inherits from:
        - `/assets/textures/ui/` (if directory metadata exists)
        - `/assets/textures/` (if directory metadata exists)
        - `/assets/` (if directory metadata exists)
        - `/` (root directory metadata)

2. **Override Priority**: Child directory tags override parent directory tags
    - Direct file tags have highest priority
    - Directory tags override based on inheritance priority
    - Root directory tags have lowest priority

3. **Inheritance Resolution**: When multiple directories could provide tags:
    - Directories with exact path matches take priority
    - Directories with higher priority values override lower ones
    - If priorities are equal, more recently modified directories take priority

4. **Path Matching Rules**:
    - Directory paths must end with `/` in metadata
    - Path matching is case-sensitive
    - Path separators must match exactly (`/` on Unix, `\` on Windows)
    - Root directory is represented as `/` in metadata

#### 1.3.4 Inheritance Examples

##### 1.3.4.1 Example 1: Basic Directory Inheritance

```yaml
# Directory metadata file (__NPK_DIR_65001__.npkdir)
directories:
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
# Result: Inherits category=texture, compression=lossless, format=png, mipmaps=true, priority=high
```

##### 1.3.4.2 Example 2: Priority-Based Override

```yaml
# Directory metadata file
directories:
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
  - path: "/assets/textures/ui/"
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
# Directory metadata file
directories:
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
  - path: "/assets/temp/"
    properties:
      - key: "category"
        value_type: "string"
        value: "temporary"
      - key: "compression"
        value_type: "string"
        value: "lossy"
    inheritance:
      enabled: false  # This directory doesn't provide inheritance
      priority: 2

# File: /assets/temp/file.png
# Result: No inherited tags (temp/ inheritance disabled)
```

### 1.4 Tag Validation

- **Key Validation**: Keys must be valid UTF-8, non-empty, max 255 bytes
- **Value Validation**: Values must match their declared type
- **JSON Validation**: JSON values must be valid JSON syntax
- **Integer Validation**: Integer values must be valid 64-bit signed integers

### 1.5 Per-File Tags Usage Examples

#### 1.5.1 Texture File Tagging

- Set comprehensive tags on texture files including category, type, format, size, compression, priority, and descriptive tags
- Example: UI button texture with PNG format, 1024x1024 size, lossless compression, priority 5, and UI/button/interface tags

##### Audio File Tagging

- Set audio-specific tags including category, type, format, duration, loop settings, and volume
- Example: Ambient forest sound with WAV format, 120-second duration, loop enabled, 0.7 volume

###### Directory Tagging

- Set tags on directories that are inherited by child files
- Example: Textures directory with texture category, lossless compression, and mipmaps enabled

###### File Search by Tags

- Search for files by specific tag values using GetFilesByTag()
- Examples: Find all texture files by category, UI files by type, high-priority files by priority level

---

## 2 Package Metadata File Specification

The package metadata file is a special file type (65000 - see [File Types System](file_type_system.md#special-files-65000-65535)) that contains structured YAML metadata about the package.

### 2.1 Metadata File Requirements

- **File Type**: Must be marked as Type 65000 (Package metadata file)
- **File Name**: Reserved name `__NPK_META_65000__.npkmeta` (case-sensitive)
- **Content Format**: YAML syntax (uncompressed for FastWrite compatibility)
- **Compression**: Optional (disabled by default for FastWrite performance)
- **Encryption**: Optional (can be encrypted like any other file)
- **Validation**: Must be valid YAML syntax

### 2.2 YAML Schema Structure

#### 2.2.1 Package Metadata Schema v1.0

##### Package Information

- **name (string)**: Package name
- **version (string)**: Package version
- **description (string)**: Package description
- **author (string)**: Package author
- **license (string)**: Package license
- **created (ISO8601-timestamp)**: Creation timestamp
- **modified (ISO8601-timestamp)**: Last modification timestamp

###### Game-Specific Metadata

- **engine (string)**: Game engine (Unity, Unreal, etc.)
- **platform (array of strings)**: Target platforms
- **genre (string)**: Game genre
- **rating (string)**: Age rating
- **requirements**: System requirements
  - **min_ram (integer)**: Minimum RAM in MB
  - **min_storage (integer)**: Minimum storage in MB
  - **graphics (string)**: Graphics requirements
  - **os (array of strings)**: Supported operating systems

###### Asset Metadata

- **textures (integer)**: Number of texture files
- **sounds (integer)**: Number of sound files
- **models (integer)**: Number of 3D model files
- **scripts (integer)**: Number of script files
- **total_size (integer)**: Total asset size in bytes

###### Security Metadata

- **encryption_level (string)**: Encryption level used
- **signature_type (string)**: Signature type used
- **security_scan (boolean)**: Whether security scan was performed
- **trusted_source (boolean)**: Whether from trusted source

###### Custom Metadata

- **custom (object)**: Extensible key-value pairs for additional metadata

### 2.3 Metadata File API

See the authoritative API definitions in [Package Metadata API](api_metadata.md) for package metadata operations and special metadata file management.

### 2.4 Package Metadata Example

**File**: `__NPK_META_65000__.npkmeta`

#### 2.2.2 Package Information Example

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

###### Asset Metadata Example

- **textures**: 1247 files
- **sounds**: 89 files
- **models**: 156 files
- **scripts**: 23 files
- **total_size**: 15728640 bytes

###### Security Metadata Example

- **encryption_level**: "ML-KEM Level 3"
- **signature_type**: "ML-DSA Level 3"
- **security_scan**: true
- **trusted_source**: true

###### Custom Metadata Example

- **build_number**: 42
- **beta_version**: false
- **dlc_ready**: true
- **achievements**: 25
