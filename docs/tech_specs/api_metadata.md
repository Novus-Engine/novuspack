# NovusPack Technical Specifications - Package Metadata API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Comment Management](#1-comment-management)
  - [1.1 PackageComment Methods](#11-packagecomment-methods)
  - [1.2 Comment Security Validation](#12-comment-security-validation)
  - [1.3 Signature Comment Security](#13-signature-comment-security)
- [2. AppID Management](#2-appid-management)
- [3. VendorID Management](#3-vendorid-management)
- [4. Combined Management](#4-combined-management)
- [5. Special Metadata File Types](#5-special-metadata-file-types)
  - [5.1 Package Metadata File (Type 65000)](#51-package-metadata-file-type-65000)
  - [5.2 Package Manifest File (Type 65001)](#52-package-manifest-file-type-65001)
  - [5.3 Package Index File (Type 65002)](#53-package-index-file-type-65002)
  - [5.4 Package Signature File (Type 65003)](#54-package-signature-file-type-65003)
  - [5.5 Special File Management](#55-special-file-management)
- [6. Metadata-Only Packages](#6-metadata-only-packages)
  - [6.1 Metadata-Only Package Definition](#61-metadata-only-package-definition)
  - [6.2 Valid Use Cases](#62-valid-use-cases)
  - [6.3 Security Considerations](#63-security-considerations)
  - [6.4 Metadata-Only Package API](#64-metadata-only-package-api)
- [7. Package Information Structures](#7-package-information-structures)
  - [7.1 PackageInfo Structure](#71-packageinfo-structure)
  - [7.2 SignatureInfo Structure](#72-signatureinfo-structure)
  - [7.3 SecurityStatus Structure](#73-securitystatus-structure)
  - [7.4 Package Information Methods](#74-package-information-methods)
- [8. Directory Metadata System](#8-directory-metadata-system)
  - [8.1 Directory Structures](#81-directory-structures)
  - [8.2 Directory Management Methods](#82-directory-management-methods)
  - [8.3 Special Metadata File Management](#83-special-metadata-file-management)
  - [8.4 Directory Association System](#84-directory-association-system)
  - [8.5 File-Directory Association](#85-file-directory-association)

---

## 0. Overview

This document defines the package metadata API for the NovusPack system, including comment management, AppID/VendorID operations, special metadata file types, and security validation for metadata fields.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Metadata System](metadata.md) - Format/schema and special file definitions (Source of Truth)

## 1. Comment Management

```go
// SetComment sets or updates the package comment
func (p *Package) SetComment(comment string) error

// GetComment retrieves the current package comment
func (p *Package) GetComment() string

// ClearComment removes the package comment
func (p *Package) ClearComment() error

// HasComment checks if the package has a comment
func (p *Package) HasComment() bool
```

### 1.1 PackageComment Methods

```go
// Size returns the size of the package comment
func (pc *PackageComment) Size() int

// WriteTo writes the comment to a writer
func (pc *PackageComment) WriteTo(w io.Writer) (int64, error)

// ReadFrom reads the comment from a reader
func (pc *PackageComment) ReadFrom(r io.Reader) (int64, error)

// Validate validates the package comment
func (pc *PackageComment) Validate() error
```

**Purpose**: Provides low-level access to package comment data and serialization.

**Size Returns**: `int` indicating the size of the comment in bytes

#### 1.1.1 WriteTo Parameters

- `w`: Writer to write comment data to

**WriteTo Returns**: Number of bytes written and error

#### 1.1.2 ReadFrom Parameters

- `r`: Reader to read comment data from

**ReadFrom Returns**: Number of bytes read and error

**Validate Returns**: Error if comment is invalid

#### 1.1.3 Error Conditions

- `ErrIOError`: I/O error during read/write operations
- `ErrInvalidComment`: Comment format is invalid
- `ErrCommentTooLarge`: Comment exceeds size limits

#### 1.1.4 Example Usage

```go
comment := &PackageComment{...}

// Get comment size
size := comment.Size()
fmt.Printf("Comment size: %d bytes\n", size)

// Write comment to file
file, err := os.Create("comment.txt")
if err != nil {
    return err
}
defer file.Close()

bytesWritten, err := comment.WriteTo(file)
if err != nil {
    return fmt.Errorf("failed to write comment: %w", err)
}

// Read comment from file
file, err := os.Open("comment.txt")
if err != nil {
    return err
}
defer file.Close()

bytesRead, err := comment.ReadFrom(file)
if err != nil {
    return fmt.Errorf("failed to read comment: %w", err)
}

// Validate comment
err = comment.Validate()
if err != nil {
    return fmt.Errorf("invalid comment: %w", err)
}
```

### 1.2 Comment Security Validation

```go
// ValidateComment validates comment content for security issues
func ValidateComment(comment string) error

// SanitizeComment sanitizes comment content to prevent injection attacks
func SanitizeComment(comment string) (string, error)

// ValidateCommentEncoding validates UTF-8 encoding of comment
func ValidateCommentEncoding(comment []byte) error

// CheckCommentLength validates comment length against limits
func CheckCommentLength(comment string) error

// DetectInjectionPatterns scans comment for malicious patterns
func DetectInjectionPatterns(comment string) ([]string, error)
```

### 1.3 Signature Comment Security

```go
// ValidateSignatureComment validates signature comment for security issues
func ValidateSignatureComment(comment string) error

// SanitizeSignatureComment sanitizes signature comment content
func SanitizeSignatureComment(comment string) (string, error)

// CheckSignatureCommentLength validates signature comment length
func CheckSignatureCommentLength(comment string) error

// AuditSignatureComment logs signature comment for security auditing
func AuditSignatureComment(comment string, signatureIndex int) error
```

## 2. AppID Management

```go
// SetAppID sets or updates the package AppID
func (p *Package) SetAppID(ctx context.Context, appID uint64) error

// GetAppID retrieves the current package AppID
func (p *Package) GetAppID() uint64

// ClearAppID removes the package AppID (set to 0)
func (p *Package) ClearAppID(ctx context.Context) error

// HasAppID checks if the package has an AppID (non-zero)
func (p *Package) HasAppID() bool

// GetAppIDInfo gets detailed AppID information if available
func (p *Package) GetAppIDInfo() AppIDInfo
```

## 3. VendorID Management

```go
// SetVendorID sets or updates the package VendorID
func (p *Package) SetVendorID(ctx context.Context, vendorID uint32) error

// GetVendorID retrieves the current package VendorID
func (p *Package) GetVendorID() uint32

// ClearVendorID removes the package VendorID (set to 0)
func (p *Package) ClearVendorID(ctx context.Context) error

// HasVendorID checks if the package has a VendorID (non-zero)
func (p *Package) HasVendorID() bool

// GetVendorIDInfo gets detailed VendorID information if available
func (p *Package) GetVendorIDInfo() VendorIDInfo
```

## 4. Combined Management

```go
// SetPackageIdentity sets both VendorID and AppID
func (p *Package) SetPackageIdentity(vendorID uint32, appID uint64) error

// GetPackageIdentity gets both VendorID and AppID
func (p *Package) GetPackageIdentity() (uint32, uint64)

// ClearPackageIdentity clears both VendorID and AppID
func (p *Package) ClearPackageIdentity() error

// GetPackageInfo gets comprehensive package information
func (p *Package) GetPackageInfo() PackageInfo
```

## 5. Special Metadata File Types

NovusPack supports special metadata file types (see [File Types System - Special Files](file_type_system.md#special-files-65000-65535)) that provide structured metadata and package management capabilities.

### 5.1 Package Metadata File (Type 65000)

```go
// AddMetadataFile adds a YAML metadata file to the package
func (p *Package) AddMetadataFile(metadata map[string]interface{}) error

// GetMetadataFile retrieves metadata from the special metadata file
func (p *Package) GetMetadataFile() (map[string]interface{}, error)

// UpdateMetadataFile updates the package metadata file
func (p *Package) UpdateMetadataFile(updates map[string]interface{}) error

// RemoveMetadataFile removes the package metadata file
func (p *Package) RemoveMetadataFile() error

// HasMetadataFile checks if package has a metadata file
func (p *Package) HasMetadataFile() bool
```

**Purpose**: Contains structured YAML metadata about the package including:

- Package description and version information
- Author and license details
- Build and compilation metadata
- Custom package-specific data

### 5.2 Package Manifest File (Type 65001)

```go
// AddManifestFile adds a package manifest file
func (p *Package) AddManifestFile(manifest ManifestData) error

// GetManifestFile retrieves the package manifest
func (p *Package) GetManifestFile() (ManifestData, error)

// UpdateManifestFile updates the package manifest
func (p *Package) UpdateManifestFile(updates ManifestData) error

// RemoveManifestFile removes the package manifest
func (p *Package) RemoveManifestFile() error

// HasManifestFile checks if package has a manifest file
func (p *Package) HasManifestFile() bool
```

**Purpose**: Defines the package structure and dependencies including:

- File organization and structure
- Dependency requirements
- Installation instructions
- Package relationships

### 5.3 Package Index File (Type 65002)

```go
// AddIndexFile adds a package index file
func (p *Package) AddIndexFile(index IndexData) error

// GetIndexFile retrieves the package index
func (p *Package) GetIndexFile() (IndexData, error)

// UpdateIndexFile updates the package index
func (p *Package) UpdateIndexFile(updates IndexData) error

// RemoveIndexFile removes the package index
func (p *Package) RemoveIndexFile() error

// HasIndexFile checks if package has an index file
func (p *Package) HasIndexFile() bool
```

**Purpose**: Provides file navigation and indexing including:

- File location mappings
- Content-based indexing
- Search and navigation data
- File relationship mappings

### 5.4 Package Signature File (Type 65003)

```go
// AddSignatureFile adds a digital signature file
func (p *Package) AddSignatureFile(signature SignatureData) error

// GetSignatureFile retrieves the signature file
func (p *Package) GetSignatureFile() (SignatureData, error)

// UpdateSignatureFile updates the signature file
func (p *Package) UpdateSignatureFile(updates SignatureData) error

// RemoveSignatureFile removes the signature file
func (p *Package) RemoveSignatureFile() error

// HasSignatureFile checks if package has a signature file
func (p *Package) HasSignatureFile() bool
```

**Purpose**: Contains digital signature information including:

- Signature metadata and timestamps
- Public key information
- Signature validation data
- Trust chain information

### 5.5 Special File Management

```go
// GetSpecialFiles returns all special files in the package
func (p *Package) GetSpecialFiles() ([]SpecialFileInfo, error)

// GetSpecialFileByType retrieves special file by type
func (p *Package) GetSpecialFileByType(fileType FileType) (SpecialFileInfo, error)

// RemoveSpecialFile removes a special file by type
func (p *Package) RemoveSpecialFile(fileType FileType) error

// ValidateSpecialFiles validates all special files
func (p *Package) ValidateSpecialFiles() error
```

#### 5.5.1 Special File Data Structures

```go
type SpecialFileInfo struct {
    Type        FileType // File type (see [File Types System - Special Files](file_type_system.md#special-files-65000-65535))
    Name        string   // Special file name (e.g., "__NPK_META_240__.npkmeta")
    Size        int64    // File size in bytes
    Offset      int64    // Offset in package
    Data        []byte   // File content
    Valid       bool     // Whether file is valid
    Error       string   // Error message if invalid
}

type ManifestData struct {
    Version     string            // Manifest version
    Package     PackageInfo       // Package information
    Dependencies []Dependency     // Package dependencies
    Structure   []FileStructure   // File organization
    Install     InstallInfo       // Installation instructions
}

type IndexData struct {
    Version     string            // Index version
    Files       []FileIndex       // File index entries
    Navigation  NavigationData    // Navigation structure
    Search      SearchIndex       // Search index data
}

type SignatureData struct {
    Version     string            // Signature format version
    Signatures  []SignatureInfo   // Signature information
    TrustChain  []TrustInfo       // Trust chain data
    Validation  ValidationData    // Validation metadata
}
```

## 6. Metadata-Only Packages

Metadata-only packages are NovusPack packages that contain only special metadata files (see [File Types System - Special Files](file_type_system.md#special-files-65000-65535)) and no regular content files. These packages serve specific purposes in package management and distribution systems.

### 6.1 Metadata-Only Package Definition

A metadata-only package is defined as:

- **FileCount = 0**: No regular content files
- **HasSpecialMetadataFiles = true**: Contains at least one special metadata file
- **TotalSize = 0**: No uncompressed content data
- **CompressedSize > 0**: Contains compressed metadata files

### 6.2 Valid Use Cases

#### 6.2.1 Package Catalogs and Registries

- **Package listings**: Catalogs of available packages with metadata
- **Dependency resolution**: Packages defining dependency trees
- **Package discovery**: Searchable indexes of package repositories

#### 6.2.2 Configuration and Schema Packages

- **Configuration templates**: Packages containing configuration schemas
- **API specifications**: Packages with API definitions and schemas
- **Data structure definitions**: Packages containing only data models

#### 6.2.3 Package Management Operations

- **Update manifests**: Packages describing updates without actual files
- **Installation scripts**: Packages containing installation instructions
- **Package relationships**: Packages defining inter-package relationships

#### 6.2.4 Development and Build Tools

- **Build configurations**: Packages containing build system configurations
- **Development metadata**: Packages with development environment specifications
- **Testing configurations**: Packages containing test specifications

### 6.3 Security Considerations

#### 6.3.1 Signature Validation

- **Metadata integrity**: Signatures must validate all metadata files
- **Empty content handling**: Special validation for packages with no content
- **Signature scope**: Clear definition of what gets signed in metadata-only packages

#### 6.3.2 Trust and Verification

- **Content verification**: No actual content to verify, trust relies on metadata
- **Metadata tampering**: Risk of metadata manipulation without content cross-reference
- **Trust chain**: Enhanced trust requirements for metadata-only packages

#### 6.3.3 Package Integrity

- **Size validation**: Very small packages require enhanced validation
- **Structure validation**: Ensure package structure is valid without content
- **Metadata consistency**: Verify metadata files are internally consistent

#### 6.3.4 Attack Vectors

- **Metadata injection**: Potential for malicious metadata injection
- **Dependency confusion**: Risk of redirecting dependencies maliciously
- **Trust abuse**: Exploiting trust in metadata-only packages

### 6.4 Metadata-Only Package API

```go
// IsMetadataOnlyPackage checks if package contains only metadata files
func (p *Package) IsMetadataOnlyPackage() bool

// ValidateMetadataOnlyPackage validates a metadata-only package
func (p *Package) ValidateMetadataOnlyPackage() error

// CreateMetadataOnlyPackage creates a new metadata-only package
func CreateMetadataOnlyPackage() (*Package, error)

// AddMetadataOnlyFile adds a special metadata file to a metadata-only package
func (p *Package) AddMetadataOnlyFile(fileType FileType, data []byte) error

// GetMetadataOnlyFiles returns all metadata files in the package
func (p *Package) GetMetadataOnlyFiles() ([]SpecialFileInfo, error)

// ValidateMetadataOnlyIntegrity validates metadata-only package integrity
func (p *Package) ValidateMetadataOnlyIntegrity() error
```

#### 6.4.1 Metadata-Only Package Validation

```go
// ValidateMetadataOnlyPackage performs comprehensive validation
func (p *Package) ValidateMetadataOnlyPackage() error {
    // Ensure FileCount == 0
    // Ensure HasSpecialMetadataFiles == true
    // Validate all special metadata files
    // Check for malicious metadata patterns
    // Verify signature scope includes all metadata
    // Ensure metadata consistency
    // Validate package structure
}
```

#### 6.4.2 Enhanced Security Requirements

- **Mandatory signatures**: Metadata-only packages must be signed
- **Enhanced validation**: Stricter validation requirements
- **Trust verification**: Higher trust requirements for metadata-only packages
- **Audit logging**: Enhanced logging for metadata-only package operations

## 7. Package Information Structures

### 7.1 PackageInfo Structure

The PackageInfo structure provides comprehensive package information and metadata:

```go
type PackageInfo struct {
    // Basic Package Information
    FileCount             int       // Number of files in the package
    FilesUncompressedSize int64     // Total uncompressed size of all files
    FilesCompressedSize   int64     // Total compressed size of all files

    // Package Identity
    VendorID       uint32    // Vendor/platform identifier
    AppID          uint64    // Application identifier

    // Package Comment
    HasComment     bool      // Whether package has a comment
    Comment        string    // Actual package comment content

    // Digital Signatures (Multiple Signatures Support)
    HasSignatures  bool      // Whether package has any signatures
    SignatureCount int       // Number of signatures in the package
    Signatures     []SignatureInfo // Detailed signature information

    // Security Information
    SecurityLevel  SecurityLevel // Overall security level
    IsImmutable    bool      // Whether package is immutable (signed)

    // Timestamps
    Created        time.Time // Package creation timestamp
    Modified       time.Time // Package modification timestamp

    // Package Features
    HasMetadataFiles    bool      // Whether package has metadata files
    HasEncryptedData    bool      // Whether package contains encrypted files
    HasCompressedData   bool      // Whether package contains compressed files
    IsMetadataOnly      bool      // Whether package contains only metadata files (no content)

    // Package Compression
    PackageCompression  uint8     // Package compression type (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    IsPackageCompressed bool      // Whether the entire package is compressed
    PackageOriginalSize int64     // Original package size before compression (0 if not compressed)
    PackageCompressedSize int64   // Compressed package size (0 if not compressed)
    PackageCompressionRatio float64 // Compression ratio (0.0-1.0, 0.0 if not compressed)
}
```

### 7.2 SignatureInfo Structure

```go
type SignatureInfo struct {
    Index         int       // Signature index in the package
    Type          uint32    // Signature type (ML-DSA, SLH-DSA, PGP, X.509)
    Size          uint32    // Size of signature data in bytes
    Offset        uint64    // Offset to signature data from start of file
    Flags         uint32    // Signature-specific flags
    Timestamp     uint32    // Unix timestamp when signature was created
    Comment       string    // Signature comment (if any)
    Algorithm     string    // Algorithm name/description
    SecurityLevel int       // Security level (1-5)
    Valid         bool      // Whether signature is valid
    Trusted       bool      // Whether signature is trusted
    Error         string    // Error message if validation failed
}
```

### 7.3 SecurityStatus Structure

```go
type SecurityStatus struct {
    SignatureCount      int                           // Number of signatures
    ValidSignatures     int                           // Number of valid signatures
    TrustedSignatures   int                           // Number of trusted signatures
    SignatureResults    []SignatureValidationResult   // Individual results
    HasChecksums        bool                          // Checksums present
    ChecksumsValid      bool                          // Checksums valid
    SecurityLevel       SecurityLevel                 // Overall security level
    ValidationErrors    []string                      // Validation errors
}
```

### 7.4 Package Information Methods

```go
// GetPackageInfo returns comprehensive package information
func (p *Package) GetPackageInfo(ctx context.Context) PackageInfo

// GetSecurityStatus returns current security status
func (p *Package) GetSecurityStatus() SecurityStatus

// RefreshPackageInfo refreshes package information from current state
func (p *Package) RefreshPackageInfo(ctx context.Context) error
```

## 8. Directory Metadata System

The directory metadata system provides structured directory definitions and inheritance rules for the NovusPack system using special metadata files.

**Cross-Reference**: For file-directory association methods and FileEntry directory operations, see [File Management API](api_file_management.md#fileentry-directory-association-methods).

**Note**: These methods assume the following imports:

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "time"
    "gopkg.in/yaml.v3"
    "github.com/google/uuid"
)
// Tag, TagValueType from api_file_management
```

### 8.1 Directory Structures

```go
// DirectoryEntry represents a directory with metadata and inheritance rules
type DirectoryEntry struct {
    Path        PathEntry         `yaml:"path"`        // Directory path entry (must end with /)
    Properties  []Tag             `yaml:"properties"`  // Directory-specific tags
    Inheritance DirectoryInheritance `yaml:"inheritance"` // Inheritance settings
    Metadata    DirectoryMetadata `yaml:"metadata"`    // Directory metadata
    FileSystem  DirectoryFileSystem `yaml:"filesystem"` // Filesystem properties

    // Directory hierarchy (runtime only, not stored in file)
    ParentDirectory *DirectoryEntry `yaml:"-"` // Pointer to parent directory (nil for root)
}

// DirectoryInheritance controls tag inheritance behavior
type DirectoryInheritance struct {
    Enabled  bool `yaml:"enabled"`  // Whether this directory provides inheritance
    Priority int  `yaml:"priority"` // Inheritance priority (higher = more specific)
}

// DirectoryMetadata contains directory metadata
type DirectoryMetadata struct {
    Created     string `yaml:"created"`     // Directory creation time (ISO8601)
    Modified    string `yaml:"modified"`    // Last modification time (ISO8601)
    Description string `yaml:"description"` // Human-readable description
}

// DirectoryFileSystem contains filesystem-specific properties
type DirectoryFileSystem struct {
    // Unix/Linux properties
    Mode    *uint32            `yaml:"mode,omitempty"`    // Directory permissions (octal)
    UID     *uint32            `yaml:"uid,omitempty"`     // User ID
    GID     *uint32            `yaml:"gid,omitempty"`     // Group ID
    ACL     []ACLEntry         `yaml:"acl,omitempty"`     // Access Control List

    // Windows properties
    WindowsAttrs *uint32       `yaml:"windows_attrs,omitempty"` // Windows attributes

    // Extended attributes
    ExtendedAttrs map[string]string `yaml:"extended_attrs,omitempty"` // Extended attributes

    // Filesystem flags
    Flags   *uint16            `yaml:"flags,omitempty"`   // Filesystem-specific flags
}

// ACLEntry represents an Access Control List entry
type ACLEntry struct {
    Type    string `yaml:"type"`    // "user", "group", "other", "mask"
    ID      *uint32 `yaml:"id,omitempty"` // User/Group ID (nil for "other")
    Perms   string  `yaml:"perms"`  // Permissions (e.g., "rwx", "r--")
}

// DirectoryEntry tag management methods
func (de *DirectoryEntry) GetTags(ctx context.Context) []Tag {
    return de.Properties
}

func (de *DirectoryEntry) SetTags(ctx context.Context, tags []Tag) {
    de.Properties = tags
}

func (de *DirectoryEntry) SetTag(ctx context.Context, key string, valueType TagValueType, value string) {
    for i, tag := range de.Properties {
        if tag.Key == key {
            de.Properties[i] = Tag{Key: key, ValueType: valueType, Value: value}
            return
        }
    }
    de.Properties = append(de.Properties, Tag{Key: key, ValueType: valueType, Value: value})
}

func (de *DirectoryEntry) GetTag(ctx context.Context, key string) (Tag, bool) {
    for _, tag := range de.Properties {
        if tag.Key == key {
            return tag, true
        }
    }
    return Tag{}, false
}

func (de *DirectoryEntry) RemoveTag(ctx context.Context, key string) {
    for i, tag := range de.Properties {
        if tag.Key == key {
            de.Properties = append(de.Properties[:i], de.Properties[i+1:]...)
            return
        }
    }
}

func (de *DirectoryEntry) HasTag(ctx context.Context, key string) bool {
    _, exists := de.GetTag(ctx, key)
    return exists
}

// Convenience methods for common tag types
func (de *DirectoryEntry) SetStringTag(ctx context.Context, key, value string) {
    de.SetTag(ctx, key, TagValueTypeString, value)
}

func (de *DirectoryEntry) GetStringTag(ctx context.Context, key string) (string, bool) {
    tag, exists := de.GetTag(ctx, key)
    if !exists || tag.ValueType != TagValueTypeString {
        return "", false
    }
    return tag.Value, true
}

func (de *DirectoryEntry) SetIntegerTag(ctx context.Context, key string, value int64) {
    de.SetTag(ctx, key, TagValueTypeInteger, strconv.FormatInt(value, 10))
}

func (de *DirectoryEntry) GetIntegerTag(ctx context.Context, key string) (int64, bool) {
    tag, exists := de.GetTag(ctx, key)
    if !exists || tag.ValueType != TagValueTypeInteger {
        return 0, false
    }
    val, err := strconv.ParseInt(tag.Value, 10, 64)
    if err != nil {
        return 0, false
    }
    return val, true
}

func (de *DirectoryEntry) SetBooleanTag(ctx context.Context, key string, value bool) {
    de.SetTag(ctx, key, TagValueTypeBoolean, strconv.FormatBool(value))
}

func (de *DirectoryEntry) GetBooleanTag(ctx context.Context, key string) (bool, bool) {
    tag, exists := de.GetTag(ctx, key)
    if !exists || tag.ValueType != TagValueTypeBoolean {
        return false, false
    }
    val, err := strconv.ParseBool(tag.Value)
    if err != nil {
        return false, false
    }
    return val, true
}

// Path management methods for DirectoryEntry
func (de *DirectoryEntry) SetPath(ctx context.Context, path string) {
    de.Path.Path = path
}

func (de *DirectoryEntry) GetPath(ctx context.Context) string {
    return de.Path.Path
}

func (de *DirectoryEntry) GetPathEntry(ctx context.Context) PathEntry {
    return de.Path
}

// Symbolic link methods for DirectoryEntry
func (de *DirectoryEntry) SetSymlink(ctx context.Context, target string) {
    de.Path.SetSymlink(ctx, target)
}

func (de *DirectoryEntry) ClearSymlink(ctx context.Context) {
    de.Path.ClearSymlink(ctx)
}

func (de *DirectoryEntry) IsSymbolicLink(ctx context.Context) bool {
    return de.Path.IsSymbolicLink(ctx)
}

func (de *DirectoryEntry) GetLinkTarget(ctx context.Context) string {
    return de.Path.GetLinkTarget(ctx)
}

func (de *DirectoryEntry) ResolveSymlink(ctx context.Context) string {
    return de.Path.ResolveSymlink(ctx)
}

// Parent directory management methods for DirectoryEntry
func (de *DirectoryEntry) SetParentDirectory(ctx context.Context, parent *DirectoryEntry) {
    de.ParentDirectory = parent
}

func (de *DirectoryEntry) GetParentDirectory(ctx context.Context) *DirectoryEntry {
    return de.ParentDirectory
}

func (de *DirectoryEntry) GetParentPath(ctx context.Context) string {
    if de.ParentDirectory == nil {
        return ""
    }
    return de.ParentDirectory.GetPath(ctx)
}

func (de *DirectoryEntry) GetDepth(ctx context.Context) int {
    depth := 0
    current := de.ParentDirectory
    for current != nil {
        depth++
        current = current.ParentDirectory
    }
    return depth
}

func (de *DirectoryEntry) IsRoot(ctx context.Context) bool {
    return de.ParentDirectory == nil
}

func (de *DirectoryEntry) GetAncestors(ctx context.Context) []*DirectoryEntry {
    var ancestors []*DirectoryEntry
    current := de.ParentDirectory
    for current != nil {
        ancestors = append(ancestors, current)
        current = current.ParentDirectory
    }
    return ancestors
}

// DirectoryInfo provides runtime directory information
type DirectoryInfo struct {
    Entry      DirectoryEntry // Directory entry data
    FileCount  int           // Number of files in this directory
    SubDirs    []string      // Immediate subdirectories
    ParentPath string        // Parent directory path
    Depth      int           // Directory depth (0 = root)
}

// FileDirectoryAssociation links files to their directory metadata
type FileDirectoryAssociation struct {
    FilePath     string        // File path
    DirectoryPath string       // Parent directory path
    Directory    *DirectoryInfo // Directory information (nil if no directory metadata)
    InheritedTags []Tag        // Tags inherited from directory hierarchy
    EffectiveTags []Tag        // All tags including inheritance
}
```

### 8.2 Directory Management Methods

```go
// Directory metadata management
func (p *Package) GetDirectoryMetadata(ctx context.Context) ([]*DirectoryEntry, error)
func (p *Package) SetDirectoryMetadata(ctx context.Context, entries []*DirectoryEntry) error
func (p *Package) AddDirectory(ctx context.Context, path string, properties map[string]string, inheritance DirectoryInheritance, metadata DirectoryMetadata) error
func (p *Package) RemoveDirectory(ctx context.Context, path string) error
func (p *Package) UpdateDirectory(ctx context.Context, path string, properties map[string]string, inheritance DirectoryInheritance, metadata DirectoryMetadata) error

// Directory information queries
func (p *Package) GetDirectoryInfo(ctx context.Context, path string) (*DirectoryInfo, error)
func (p *Package) ListDirectories(ctx context.Context) ([]DirectoryInfo, error)
func (p *Package) GetDirectoryHierarchy(ctx context.Context) (map[string][]string, error)

// Directory validation
func (p *Package) ValidateDirectoryMetadata(ctx context.Context) error
func (p *Package) GetDirectoryConflicts(ctx context.Context) ([]string, error)

// Special metadata file management
func (p *Package) SaveDirectoryMetadataFile(ctx context.Context) error
func (p *Package) LoadDirectoryMetadataFile(ctx context.Context) error
func (p *Package) UpdateSpecialMetadataFlags(ctx context.Context) error

// Special metadata file creation helpers
func (p *Package) CreateSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string, content []byte) (*FileEntry, error)
func (p *Package) UpdateSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string, content []byte) error
func (p *Package) RemoveSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string) error

// Directory association management
func (p *Package) AssociateFileWithDirectory(ctx context.Context, filePath string, dirPath string) error
func (p *Package) DisassociateFileFromDirectory(ctx context.Context, filePath string) error
func (p *Package) UpdateFileDirectoryAssociations(ctx context.Context) error
func (p *Package) GetFileDirectoryAssociations(ctx context.Context) (map[string]*DirectoryEntry, error)
```

### 8.3 Special Metadata File Management

Special metadata files must be saved with specific flags and file types to ensure proper recognition and processing.

#### 8.3.1 Special File Requirements

##### 8.3.1.1 File Type Requirements

- Must use special file types (see [File Types System - Special Files](file_type_system.md#special-files-65000-65535))
- Must have reserved file names (e.g., `__NPK_DIR_241__.npkdir`)
- Must be uncompressed for FastWrite compatibility
- Must have proper package header flags set

##### 8.3.1.2 Special File Types

- **Type 65000**: Package metadata (`__NPK_PKG_65000__.yaml`)
- **Type 65001**: Directory metadata (`__NPK_DIR_65001__.npkdir`)
- **Type 65002**: Symbolic link metadata (`__NPK_SYMLINK_65002__.npksym`)
- **Type 65003-65535**: Reserved for future use

##### 8.3.1.3 Package Header Flags

- **Bit 6**: Has special metadata files (set to 1 when special files exist)
- **Bit 5**: Has per-file tags (set to 1 if directory metadata provides inheritance)

##### 8.3.1.4 FileEntry Requirements

- `Type` field set to appropriate special file type (65001 for directory metadata)
- `CompressionType` set to 0 (no compression)
- `EncryptionType` set to 0x00 (no encryption) - special files should not be encrypted
- `Tags` should include `file_type=special_metadata` and `metadata_type=directory`

#### 8.3.2 Implementation Details

```go
// SaveDirectoryMetadataFile creates and saves the directory metadata file
func (p *Package) SaveDirectoryMetadataFile(ctx context.Context) error {
    // 1. Get current directory metadata
    entries, err := p.GetDirectoryMetadata(ctx)
    if err != nil {
        return err
    }

    // 2. Marshal to YAML
    yamlData, err := yaml.Marshal(map[string]interface{}{
        "directories": entries,
    })
    if err != nil {
        return err
    }

    // 3. Create special metadata file entry
    fileEntry, err := p.CreateSpecialMetadataFile(ctx, 65001, "__NPK_DIR_65001__.npkdir", yamlData)
    if err != nil {
        return err
    }

    // 4. Set appropriate tags
    fileEntry.SetStringTag("file_type", "special_metadata")
    fileEntry.SetStringTag("metadata_type", "directory")
    fileEntry.SetStringTag("format", "yaml")
    fileEntry.SetIntegerTag("version", 1)

    // 5. Update package flags
    return p.UpdateSpecialMetadataFlags(ctx)
}

// UpdateSpecialMetadataFlags updates package header flags based on special files
func (p *Package) UpdateSpecialMetadataFlags(ctx context.Context) error {
    // Check for special metadata files
    hasSpecialFiles := p.hasSpecialMetadataFiles()
    hasPerFileTags := p.hasPerFileTags()

    // Update package header flags
    return p.updatePackageFlags(hasSpecialFiles, hasPerFileTags)
}
```

### 8.4 Directory Association System

The directory association system links FileEntry objects to their corresponding DirectoryEntry metadata, enabling tag inheritance and filesystem property management.

#### 8.4.1 Association Properties

#### 8.4.1 FileEntry Directory Properties

- `DirectoryEntry` - Pointer to the directory metadata for the file's immediate directory
- `ParentDirectory` - Pointer to the parent directory metadata (for inheritance resolution)
- `InheritedTags` - Cached inherited tags from the directory hierarchy

#### 8.4.2 DirectoryEntry Filesystem Properties

- `Mode` - Unix/Linux directory permissions (octal)
- `UID`/`GID` - User and Group IDs
- `ACL` - Access Control List entries
- `WindowsAttrs` - Windows directory attributes
- `ExtendedAttrs` - Extended attributes map
- `Flags` - Filesystem-specific flags

#### 8.4.3 Association Management

```go
// AssociateFileWithDirectory links a file to its directory metadata
func (p *Package) AssociateFileWithDirectory(ctx context.Context, filePath string, dirPath string) error {
    // 1. Get file entry
    fileEntry, err := p.GetFileByPath(ctx, filePath)
    if err != nil {
        return err
    }

    // 2. Get directory entry
    dirInfo, err := p.GetDirectoryInfo(ctx, dirPath)
    if err != nil {
        return err
    }

    // 3. Set association
    fileEntry.SetDirectoryEntry(ctx, &dirInfo.Entry)

    // 4. Resolve parent directory
    parentPath := filepath.Dir(dirPath)
    if parentPath != dirPath && parentPath != "." {
        parentInfo, err := p.GetDirectoryInfo(ctx, parentPath)
        if err == nil {
            fileEntry.SetParentDirectory(ctx, &parentInfo.Entry)
        }
    }

    // 5. Resolve inherited tags
    inheritedTags, err := p.resolveInheritedTags(ctx, filePath)
    if err == nil {
        fileEntry.UpdateInheritedTags(ctx, inheritedTags)
    }

    return nil
}

// UpdateFileDirectoryAssociations rebuilds all file-directory associations
func (p *Package) UpdateFileDirectoryAssociations(ctx context.Context) error {
    // 1. Get all files
    files, err := p.ListFiles(ctx)
    if err != nil {
        return err
    }

    // 2. Get all directories
    dirs, err := p.GetDirectoryMetadata(ctx)
    if err != nil {
        return err
    }

    // 3. Build directory path map
    dirMap := make(map[string]*DirectoryEntry)
    for _, dir := range dirs {
        dirMap[dir.Path] = &dir
    }

    // 4. Associate each file with its directory
    for _, file := range files {
        for _, path := range file.Paths {
            dirPath := filepath.Dir(path.Path)
            if dirEntry, exists := dirMap[dirPath]; exists {
                file.SetDirectoryEntry(ctx, dirEntry)

                // Set parent directory
                parentPath := filepath.Dir(dirPath)
                if parentEntry, exists := dirMap[parentPath]; exists {
                    file.SetParentDirectory(ctx, parentEntry)
                }
            }
        }
    }

    return nil
}
```

### 8.5 File-Directory Association

```go
// File-directory association methods
func (p *Package) GetFileDirectoryAssociation(ctx context.Context, filePath string) (*FileDirectoryAssociation, error)
func (p *Package) GetFileInheritedTags(ctx context.Context, filePath string) ([]Tag, error)
func (p *Package) GetFileEffectiveTags(ctx context.Context, filePath string) ([]Tag, error)
func (p *Package) GetFilesInDirectory(ctx context.Context, dirPath string) ([]*FileEntry, error)
func (p *Package) GetDirectoryFiles(ctx context.Context, dirPath string) ([]*FileEntry, error)

// Directory hierarchy analysis
func (p *Package) GetDirectoryTree(ctx context.Context) (*DirectoryTree, error)
func (p *Package) GetDirectoryStats(ctx context.Context) (map[string]DirectoryStats, error)

// Symbolic link management methods
func (p *Package) AddSymlink(ctx context.Context, symlink SymlinkEntry) error
func (p *Package) RemoveSymlink(ctx context.Context, sourcePath string) error
func (p *Package) GetSymlink(ctx context.Context, sourcePath string) (*SymlinkEntry, error)
func (p *Package) ListSymlinks(ctx context.Context) ([]SymlinkEntry, error)
func (p *Package) UpdateSymlink(ctx context.Context, sourcePath string, symlink SymlinkEntry) error
func (p *Package) SaveSymlinkMetadataFile(ctx context.Context, symlink SymlinkEntry) error
func (p *Package) LoadSymlinkMetadataFile(ctx context.Context, fileEntry *FileEntry) (*SymlinkEntry, error)

// DirectoryTree represents the complete directory hierarchy
type DirectoryTree struct {
    Root      *DirectoryNode `json:"root"`
    TotalDirs int           `json:"total_dirs"`
    TotalFiles int          `json:"total_files"`
}

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
    Path      string           `json:"path"`
    Directory *DirectoryInfo   `json:"directory,omitempty"`
    Files     []*FileEntry     `json:"files"`
    Children  []*DirectoryNode `json:"children"`
}

// DirectoryStats provides statistics for a directory
type DirectoryStats struct {
    FileCount    int    `json:"file_count"`
    TotalSize    int64  `json:"total_size"`
    CompressedSize int64 `json:"compressed_size"`
    LastModified string `json:"last_modified"`
}

// SymlinkEntry represents a symbolic link with metadata
type SymlinkEntry struct {
    SourcePath    string    `yaml:"source_path"`    // Original symlink path
    TargetPath    string    `yaml:"target_path"`    // Target path (resolved)
    Properties    []Tag     `yaml:"properties"`     // Symlink-specific tags
    Metadata      SymlinkMetadata `yaml:"metadata"` // Symlink metadata
    FileSystem    SymlinkFileSystem `yaml:"filesystem"` // Filesystem properties
}

// SymlinkMetadata contains symlink creation and modification information
type SymlinkMetadata struct {
    Created   time.Time `yaml:"created"`   // When symlink was created
    Modified  time.Time `yaml:"modified"`  // When symlink was last modified
    Description string  `yaml:"description,omitempty"` // Optional description
}

// SymlinkFileSystem contains filesystem-specific properties for symlinks
type SymlinkFileSystem struct {
    Mode    *uint32            `yaml:"mode,omitempty"`    // Symlink permissions (octal)
    UID     *uint32            `yaml:"uid,omitempty"`     // User ID
    GID     *uint32            `yaml:"gid,omitempty"`     // Group ID
    ACL     []ACLEntry         `yaml:"acl,omitempty"`     // Access Control List
    WindowsAttrs *uint32       `yaml:"windows_attrs,omitempty"` // Windows attributes
    ExtendedAttrs map[string]string `yaml:"extended_attrs,omitempty"` // Extended attributes
    Flags   *uint16            `yaml:"flags,omitempty"`   // Filesystem-specific flags
}
```
