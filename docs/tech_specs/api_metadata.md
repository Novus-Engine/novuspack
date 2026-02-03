# NovusPack Technical Specifications - Package Metadata API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Comment Management](#1-comment-management)
  - [1.1 Package-Level Comment Methods](#11-package-level-comment-methods)
    - [1.1.1 Package SetComment Method](#111-packagesetcomment-method)
    - [1.1.2 Package GetComment Method](#112-packagegetcomment-method)
    - [1.1.3 Package ClearComment Method](#113-packageclearcomment-method)
    - [1.1.4 Package HasComment Method](#114-packagehascomment-method)
  - [1.2 PackageComment Structure](#12-packagecomment-structure)
  - [1.3 PackageComment Methods](#13-packagecomment-methods)
    - [1.3.1 PackageComment Size Method](#131-packagecommentsize-method)
    - [1.3.2 PackageComment WriteTo Method](#132-packagecommentwriteto-method)
    - [1.3.3 PackageComment ReadFrom Method](#133-packagecommentreadfrom-method)
    - [1.3.4 PackageComment Validate Method](#134-packagecommentvalidate-method)
    - [1.3.5 NewPackageComment Function](#135-newpackagecomment-function)
  - [1.4 Comment Security Validation](#14-comment-security-validation)
    - [1.4.1 ValidateComment Function](#141-validatecomment-function)
    - [1.4.2 SanitizeComment Function](#142-sanitizecomment-function)
    - [1.4.3 ValidateCommentEncoding Function](#143-validatecommentencoding-function)
    - [1.4.4 CheckCommentLength Function](#144-checkcommentlength-function)
    - [1.4.5 DetectInjectionPatterns Function](#145-detectinjectionpatterns-function)
  - [1.5 Signature Comment Security](#15-signature-comment-security)
    - [1.5.1 ValidateSignatureComment Function](#151-validatesignaturecomment-function)
    - [1.5.2 SanitizeSignatureComment Function](#152-sanitizesignaturecomment-function)
    - [1.5.3 CheckSignatureCommentLength Function](#153-checksignaturecommentlength-function)
    - [1.5.4 AuditSignatureComment Function](#154-auditsignaturecomment-function)
- [2. AppID Management](#2-appid-management)
  - [2.1 SetAppID Method](#21-packagesetappid-method)
    - [2.1.1 Package GetAppID Method](#211-packagegetappid-method)
    - [2.1.2 Package ClearAppID Method](#212-packageclearappid-method)
    - [2.1.3 Package HasAppID Method](#213-packagehasappid-method)
    - [2.1.4 Package GetAppIDInfo Method](#214-packagegetappidinfo-method)
- [3. VendorID Management](#3-vendorid-management)
  - [3.1 SetVendorID Method](#31-packagesetvendorid-method)
    - [3.1.1 Package GetVendorID Method](#311-packagegetvendorid-method)
    - [3.1.2 Package ClearVendorID Method](#312-packageclearvendorid-method)
    - [3.1.3 Package HasVendorID Method](#313-packagehasvendorid-method)
    - [3.1.4 Package GetVendorIDInfo Method](#314-packagegetvendoridinfo-method)
- [4. Combined Management](#4-combined-management)
  - [4.1 SetPackageIdentity Method](#41-packagesetpackageidentity-method)
    - [4.1.1 Package GetPackageIdentity Method](#411-packagegetpackageidentity-method)
    - [4.1.2 Package ClearPackageIdentity Method](#412-packageclearpackageidentity-method)
- [5. Special Metadata File Types](#5-special-metadata-file-types)
  - [5.1 Package Metadata File (Type 65000)](#51-package-metadata-file-type-65000)
    - [5.1.1 Package AddMetadataFile Method](#511-packageaddmetadatafile-method)
    - [5.1.2 Package GetMetadataFile Method](#512-packagegetmetadatafile-method)
    - [5.1.3 Package UpdateMetadataFile Method](#513-packageupdatemetadatafile-method)
    - [5.1.4 Package RemoveMetadataFile Method](#514-packageremovemetadatafile-method)
    - [5.1.5 Package HasMetadataFile Method](#515-packagehasmetadatafile-method)
  - [5.2 Package Manifest File (Type 65001)](#52-package-manifest-file-type-65001)
    - [5.2.1 Package AddManifestFile Method](#521-packageaddmanifestfile-method)
    - [5.2.2 Package GetManifestFile Method](#522-packagegetmanifestfile-method)
    - [5.2.3 Package UpdateManifestFile Method](#523-packageupdatemanifestfile-method)
    - [5.2.4 Package RemoveManifestFile Method](#524-packageremovemanifestfile-method)
    - [5.2.5 Package HasManifestFile Method](#525-packagehasmanifestfile-method)
  - [5.3 Package Index File (Type 65002)](#53-package-index-file-type-65002)
    - [5.3.1 Package AddIndexFile Method](#531-packageaddindexfile-method)
    - [5.3.2 Package GetIndexFile Method](#532-packagegetindexfile-method)
    - [5.3.3 Package UpdateIndexFile Method](#533-packageupdateindexfile-method)
    - [5.3.4 Package RemoveIndexFile Method](#534-packageremoveindexfile-method)
    - [5.3.5 Package HasIndexFile Method](#535-packagehasindexfile-method)
  - [5.4 Package Signature File (Type 65003)](#54-package-signature-file-type-65003)
    - [5.4.1 Package AddSignatureFile Method](#541-packageaddsignaturefile-method)
    - [5.4.2 Package GetSignatureFile Method](#542-packagegetsignaturefile-method)
    - [5.4.3 Package UpdateSignatureFile Method](#543-packageupdatesignaturefile-method)
    - [5.4.4 Package RemoveSignatureFile Method](#544-packageremovesignaturefile-method)
    - [5.4.5 Package HasSignatureFile Method](#545-packagehassignaturefile-method)
  - [5.5 Special File Management](#55-special-file-management)
    - [5.5.1 Package GetSpecialFiles Method](#551-packagegetspecialfiles-method)
    - [5.5.2 Package GetSpecialFileByType Method](#552-packagegetspecialfilebytype-method)
    - [5.5.3 Package RemoveSpecialFile Method](#553-packageremovespecialfile-method)
    - [5.5.4 Package ValidateSpecialFiles Method](#554-packagevalidatespecialfiles-method)
    - [5.5.5 Special File Data Structures](#555-special-file-data-structures)
- [6. Metadata-Only Packages](#6-metadata-only-packages)
  - [6.1 Metadata-Only Package Definition](#61-metadata-only-package-definition)
  - [6.2 Valid Use Cases](#62-valid-use-cases)
    - [6.2.1 Package Catalogs and Registries](#621-package-catalogs-and-registries)
    - [6.2.2 Configuration and Schema Packages](#622-configuration-and-schema-packages)
    - [6.2.3 Package Management Operations](#623-package-management-operations)
    - [6.2.4 Development and Build Tools](#624-development-and-build-tools)
    - [6.2.5 Empty and Placeholder Packages](#625-empty-and-placeholder-packages)
  - [6.3 Security Considerations](#63-security-considerations)
    - [6.3.1 Signature Validation](#631-signature-validation)
    - [6.3.2 Trust and Verification](#632-trust-and-verification)
    - [6.3.3 Package Integrity](#633-package-integrity)
    - [6.3.4 Attack Vectors](#634-attack-vectors)
  - [6.4 Metadata-Only Package API](#64-metadata-only-package-api)
    - [6.4.1 Package IsMetadataOnlyPackage Method](#641-packageismetadataonlypackage-method)
    - [6.4.2 Package AddMetadataOnlyFile Method](#642-packageaddmetadataonlyfile-method)
    - [6.4.3 Package GetMetadataOnlyFiles Method](#643-packagegetmetadataonlyfiles-method)
    - [6.4.4 Package ValidateMetadataOnlyIntegrity Method](#644-packagevalidatemetadataonlyintegrity-method)
    - [6.4.5 Metadata-Only Package Validation](#645-packagevalidatemetadataonlypackage-method)
    - [6.4.6 Write Operation Requirements](#646-write-operation-requirements)
    - [6.4.7 Security Considerations for Metadata-Only Packages](#647-security-considerations-for-metadata-only-packages)
- [7. Package Information Structures](#7-package-information-structures)
  - [7.1 PackageInfo Structure](#71-packageinfo-structure)
    - [7.1.1 PackageInfo Scope and Exclusions](#711-packageinfo-scope-and-exclusions)
    - [7.1.2 NewPackageInfo Function](#712-newpackageinfo-function)
    - [7.1.3. PackageInfo As Source of Truth](#713-packageinfo-as-source-of-truth)
    - [7.1.4 PackageInfo.FromHeader Method](#714-packageinfofromheader-method)
    - [7.1.5 PackageHeader Structure](#715-packageheader-structure)
    - [7.1.6 PackageHeader.ToHeader Method](#716-packageheadertoheader-method)
  - [7.2 SignatureInfo Structure](#72-signatureinfo-structure)
  - [7.3 SecurityStatus Structure](#73-securitystatus-structure)
  - [7.4 Package Information Methods](#74-package-information-methods)
    - [7.4.1 Package GetPackageInfo Method](#741-packagegetpackageinfo-method)
    - [7.4.2 Package RefreshPackageInfo Method](#742-packagerefreshpackageinfo-method)
- [8. PathMetadata System](#8-pathmetadata-system)
  - [8.1 PathMetadata Structures](#81-pathmetadata-structures)
    - [8.1.1 PathMetadataType Type](#811-pathmetadatatype-type)
    - [8.1.2 PathMetadataEntry Structure](#812-pathmetadataentry-structure)
    - [8.1.3 PathInheritance Structure](#813-pathinheritance-structure)
    - [8.1.4 PathMetadata Structure](#814-pathmetadata-structure)
    - [8.1.5 PathFileSystem Structure](#815-pathfilesystem-structure)
    - [8.1.6 ACLEntry Structure](#816-aclentry-structure)
    - [8.1.7 PathMetadataEntry Tag Management](#817-pathmetadataentry-tag-management)
    - [8.1.8 PathMetadataEntry Methods](#818-pathmetadataentry-methods)
    - [8.1.9 PathMetadataEntry Validation Methods](#819-pathmetadataentry-validation-methods)
    - [8.1.10 `PathInfo` Structure](#8110-pathinfo-structure)
    - [8.1.11 `FilePathAssociation` Structure](#8111-filepathassociation-structure)
    - [8.1.12 `DestPathOverride` Structure](#8112-destpathoverride-structure)
    - [8.1.13 `DestPathInput` Interface](#8113-destpathinput-interface)
  - [8.2 `PathMetadata` Management Methods](#82-pathmetadata-management-methods)
    - [8.2.1 Core `PathMetadata` CRUD Operations](#821-core-pathmetadata-crud-operations)
    - [Path Information Query Methods](#path-information-query-methods)
    - [Path Association Methods](#path-association-methods)
    - [Special Metadata File Management](#special-metadata-file-management)
    - [Special Metadata File Creation Helpers](#special-metadata-file-creation-helpers)
    - [8.2.2 Package GetPathInfo Method](#822-packagegetpathinfo-method)
    - [8.2.3 Package ListPaths Method](#823-packagelistpaths-method)
    - [8.2.4 Package ListDirectories Method](#824-packagelistdirectories-method)
    - [8.2.5 Package GetDirectoryCount Method](#825-packagegetdirectorycount-method)
    - [8.2.6 Package GetPathHierarchy Method](#826-packagegetpathhierarchy-method)
    - [8.2.7 Package AssociateFileWithPath Method](#827-packageassociatefilewithpath-method)
    - [8.2.8 Package UpdateFilePathAssociations Method](#828-packageupdatefilepathassociations-method)
  - [8.3 Special Metadata File Management](#83-special-metadata-file-management)
    - [8.3.1 Special File Type Requirements](#831-special-file-type-requirements)
    - [8.3.2 Special File Types](#832-special-file-types)
    - [8.3.3 PackageHeader Flags](#833-packageheader-flags)
    - [8.3.4 FileEntry Requirements](#834-fileentry-requirements)
    - [8.3.5 Implementation Details](#835-implementation-details)
  - [8.4 Path Association System](#84-path-association-system)
    - [8.4.1 FileEntry Path Properties](#841-fileentry-path-properties)
    - [8.4.2 PathMetadataEntry Filesystem Properties](#842-pathmetadataentry-filesystem-properties)
    - [8.4.3 Association Management](#843-association-management)
  - [8.5 File-Path Association](#85-file-path-association)
    - [8.5.1 File-Path Association Query Methods](#851-file-path-association-query-methods)
    - [8.5.2 Path Hierarchy Analysis Methods](#852-path-hierarchy-analysis-methods)
    - [8.5.3 Symbolic Link Management Methods](#853-symbolic-link-management-methods)
    - [8.5.4 Symlink Validation Methods](#854-symlink-validation-methods)
    - [8.5.5 PathTree Structure](#855-pathtree-structure)
    - [8.5.6 PathNode Structure](#856-pathnode-structure)
    - [8.5.7 PathStats Structure](#857-pathstats-structure)
    - [8.5.8 SymlinkEntry Structure](#858-symlinkentry-structure)
    - [8.5.9 SymlinkMetadata Structure](#859-symlinkmetadata-structure)
    - [8.5.10 SymlinkFileSystem Structure](#8510-symlinkfilesystem-structure)
    - [8.5.11 Symlink Validation Methods (Details)](#8511-symlink-validation-methods-details)
    - [8.5.12 Symlink Creation from Duplicate Paths](#8512-symlink-creation-from-duplicate-paths)

---

## 0. Overview

This document defines the package metadata API for the NovusPack system, including comment management, AppID/VendorID operations, special metadata file types, and security validation for metadata fields.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [Metadata System](metadata.md) - Format/schema and special file definitions (Source of Truth)

## 1. Comment Management

This section describes comment management operations for packages.

### 1.1 Package-Level Comment Methods

This section describes package-level comment management methods.

#### 1.1.1 Package.SetComment Method

```go
// SetComment sets or updates the package comment
// Returns *PackageError on failure
func (p *Package) SetComment(comment string) error
```

#### 1.1.2 Package.GetComment Method

```go
// GetComment retrieves the current package comment
func (p *Package) GetComment() string
```

#### 1.1.3 Package.ClearComment Method

```go
// ClearComment removes the package comment
// Returns *PackageError on failure
func (p *Package) ClearComment() error
```

#### 1.1.4 Package.HasComment Method

```go
// HasComment checks if the package has a comment
func (p *Package) HasComment() bool
```

### 1.2 PackageComment Structure

```go
// PackageComment represents the optional package comment section
type PackageComment struct {
    CommentLength uint32   // Length of comment including null terminator
    Comment       string   // UTF-8 encoded comment string (null-terminated)
    Reserved      [3]uint8 // Reserved for future use (must be 0)
}
```

**Size**: Variable (4 + comment_length + 3 bytes)

**Purpose**: Provides structured storage for package comment data with validation and serialization support.

### 1.3 PackageComment Methods

This section describes PackageComment methods for comment operations.

#### 1.3.1 PackageComment.Size Method

```go
// Size returns the size of the package comment
func (pc *PackageComment) Size() int
```

#### 1.3.2 PackageComment.WriteTo Method

```go
// WriteTo writes the comment to a writer
func (pc *PackageComment) WriteTo(w io.Writer) (int64, error)
```

#### 1.3.3 PackageComment.ReadFrom Method

```go
// ReadFrom reads the comment from a reader
func (pc *PackageComment) ReadFrom(r io.Reader) (int64, error)
```

#### 1.3.4 PackageComment.Validate Method

```go
// Validate validates the package comment
// Returns *PackageError on failure
func (pc *PackageComment) Validate() error
```

**Purpose**: Provides low-level access to package comment data and serialization.

**Size Returns**: `int` indicating the size of the comment in bytes

##### 1.3.4.1 WriteTo Parameters

- `w`: Writer to write comment data to

**WriteTo Returns**: Number of bytes written and error

##### 1.3.4.2 ReadFrom Parameters

- `r`: Reader to read comment data from

**ReadFrom Returns**: Number of bytes read and error

##### 1.3.4.3 Validate Returns

**Validate Returns**: Error if comment is invalid

##### 1.3.4.4 Error Conditions

- `ErrTypeValidation`: Comment format is invalid, comment exceeds size limits
- `ErrTypeIO`: I/O error during read/write operations

##### 1.3.4.5 Example Usage

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

#### 1.3.5 NewPackageComment Function

Creates a new PackageComment with proper initialization.

```go
// NewPackageComment creates and returns a new PackageComment with zero values
func NewPackageComment() *PackageComment
```

Returns a new PackageComment instance with all fields initialized to their zero values:

- `CommentLength` set to 0
- `Comment` set to empty string
- `Reserved` bytes all set to 0

This is equivalent to an empty comment state and is the primary way to create a new PackageComment instance.

**Note**: For unmarshaling PackageComment instances from binary data, see [ReadFrom](#13-packagecomment-methods) in the PackageComment Methods section.

### 1.4 Comment Security Validation

This section describes comment security validation functions.

#### 1.4.1 ValidateComment Function

```go
// ValidateComment validates comment content for security issues
// Returns *PackageError on failure
func ValidateComment(comment string) error
```

#### 1.4.2 SanitizeComment Function

```go
// SanitizeComment sanitizes comment content to prevent injection attacks
// Returns *PackageError on failure
func SanitizeComment(comment string) (string, error)
```

#### 1.4.3 ValidateCommentEncoding Function

```go
// ValidateCommentEncoding validates UTF-8 encoding of comment
// Returns *PackageError on failure
func ValidateCommentEncoding(comment []byte) error
```

#### 1.4.4 CheckCommentLength Function

```go
// CheckCommentLength validates comment length against limits
// Returns *PackageError on failure
func CheckCommentLength(comment string) error
```

#### 1.4.5 DetectInjectionPatterns Function

```go
// DetectInjectionPatterns scans comment for malicious patterns
func DetectInjectionPatterns(comment string) ([]string, error)
```

### 1.5 Signature Comment Security

This section describes signature comment security features.

#### 1.5.1 ValidateSignatureComment Function

```go
// ValidateSignatureComment validates signature comment for security issues
// Returns *PackageError on failure
func ValidateSignatureComment(comment string) error
```

#### 1.5.2 SanitizeSignatureComment Function

```go
// SanitizeSignatureComment sanitizes signature comment content
// Returns *PackageError on failure
func SanitizeSignatureComment(comment string) (string, error)
```

#### 1.5.3 CheckSignatureCommentLength Function

```go
// CheckSignatureCommentLength validates signature comment length
// Returns *PackageError on failure
func CheckSignatureCommentLength(comment string) error
```

#### 1.5.4 AuditSignatureComment Function

```go
// AuditSignatureComment logs signature comment for security auditing
// Returns *PackageError on failure
func AuditSignatureComment(comment string, signatureIndex int) error
```

## 2. AppID Management

This section describes AppID management operations.

### 2.1 Package.SetAppID Method

```go
// SetAppID sets or updates the package AppID
// Returns *PackageError on failure
func (p *Package) SetAppID(appID uint64) error
```

#### 2.1.1 Package.GetAppID Method

```go
// GetAppID retrieves the current package AppID
func (p *Package) GetAppID() uint64
```

#### 2.1.2 Package.ClearAppID Method

```go
// ClearAppID removes the package AppID (set to 0)
// Returns *PackageError on failure
func (p *Package) ClearAppID() error
```

#### 2.1.3 Package.HasAppID Method

```go
// HasAppID checks if the package has an AppID (non-zero)
func (p *Package) HasAppID() bool
```

#### 2.1.4 Package.GetAppIDInfo Method

```go
// GetAppIDInfo gets detailed AppID information if available
func (p *Package) GetAppIDInfo() AppIDInfo
```

## 3. VendorID Management

This section describes VendorID management operations.

### 3.1 Package.SetVendorID Method

```go
// SetVendorID sets or updates the package VendorID
// Returns *PackageError on failure
func (p *Package) SetVendorID(vendorID uint32) error
```

#### 3.1.1 Package.GetVendorID Method

```go
// GetVendorID retrieves the current package VendorID
func (p *Package) GetVendorID() uint32
```

#### 3.1.2 Package.ClearVendorID Method

```go
// ClearVendorID removes the package VendorID (set to 0)
// Returns *PackageError on failure
func (p *Package) ClearVendorID() error
```

#### 3.1.3 Package.HasVendorID Method

```go
// HasVendorID checks if the package has a VendorID (non-zero)
func (p *Package) HasVendorID() bool
```

#### 3.1.4 Package.GetVendorIDInfo Method

```go
// GetVendorIDInfo gets detailed VendorID information if available
func (p *Package) GetVendorIDInfo() VendorIDInfo
```

## 4. Combined Management

This section describes combined package identity management operations.

### 4.1 Package.SetPackageIdentity Method

```go
// SetPackageIdentity sets both VendorID and AppID
// Returns *PackageError on failure
func (p *Package) SetPackageIdentity(vendorID uint32, appID uint64) error
```

#### 4.1.1 Package.GetPackageIdentity Method

```go
// GetPackageIdentity gets both VendorID and AppID
func (p *Package) GetPackageIdentity() (uint32, uint64)
```

#### 4.1.2 Package.ClearPackageIdentity Method

```go
// ClearPackageIdentity clears both VendorID and AppID
// Returns *PackageError on failure
func (p *Package) ClearPackageIdentity() error
```

See [Package Information Methods](#74-package-information-methods) for `GetPackageInfo`.

## 5. Special Metadata File Types

NovusPack supports special metadata file types (see [File Types System - Special Files](file_type_system.md#339-special-file-types-65000-65535)) that provide structured metadata and package management capabilities.

### 5.1 Package Metadata File (Type 65000)

This section describes the package metadata file type and operations.

#### 5.1.1 Package.AddMetadataFile Method

```go
// AddMetadataFile adds a YAML metadata file to the package
// Returns *PackageError on failure
func (p *Package) AddMetadataFile(metadataYAML []byte) error
```

#### 5.1.2 Package.GetMetadataFile Method

```go
// GetMetadataFile retrieves metadata from the special metadata file
// Returns *PackageError on failure
func (p *Package) GetMetadataFile() ([]byte, error)
```

#### 5.1.3 Package.UpdateMetadataFile Method

```go
// UpdateMetadataFile updates the package metadata file
// Returns *PackageError on failure
func (p *Package) UpdateMetadataFile(metadataYAML []byte) error
```

#### 5.1.4 Package.RemoveMetadataFile Method

```go
// RemoveMetadataFile removes the package metadata file
// Returns *PackageError on failure
func (p *Package) RemoveMetadataFile() error
```

#### 5.1.5 Package.HasMetadataFile Method

```go
// HasMetadataFile checks if package has a metadata file
func (p *Package) HasMetadataFile() bool
```

**Purpose**: Contains structured YAML metadata about the package including:

- Package description and version information
- Author and license details
- Build and compilation metadata
- Custom package-specific data

### 5.2 Package Manifest File (Type 65001)

This section describes the package manifest file type and operations.

#### 5.2.1 Package.AddManifestFile Method

```go
// AddManifestFile adds a package manifest file
// Returns *PackageError on failure
func (p *Package) AddManifestFile(manifest ManifestData) error
```

#### 5.2.2 Package.GetManifestFile Method

```go
// GetManifestFile retrieves the package manifest
// Returns *PackageError on failure
func (p *Package) GetManifestFile() (ManifestData, error)
```

#### 5.2.3 Package.UpdateManifestFile Method

```go
// UpdateManifestFile updates the package manifest
// Returns *PackageError on failure
func (p *Package) UpdateManifestFile(updates ManifestData) error
```

#### 5.2.4 Package.RemoveManifestFile Method

```go
// RemoveManifestFile removes the package manifest
// Returns *PackageError on failure
func (p *Package) RemoveManifestFile() error
```

#### 5.2.5 Package.HasManifestFile Method

```go
// HasManifestFile checks if package has a manifest file
func (p *Package) HasManifestFile() bool
```

**Purpose**: Defines the package structure and dependencies including:

- File organization and structure
- Dependency requirements
- Installation instructions
- Package relationships

### 5.3 Package Index File (Type 65002)

This section describes the package index file type and operations.

#### 5.3.1 Package.AddIndexFile Method

```go
// AddIndexFile adds a package index file
// Returns *PackageError on failure
func (p *Package) AddIndexFile(index IndexData) error
```

#### 5.3.2 Package.GetIndexFile Method

```go
// GetIndexFile retrieves the package index
// Returns *PackageError on failure
func (p *Package) GetIndexFile() (IndexData, error)
```

#### 5.3.3 Package.UpdateIndexFile Method

```go
// UpdateIndexFile updates the package index
// Returns *PackageError on failure
func (p *Package) UpdateIndexFile(updates IndexData) error
```

#### 5.3.4 Package.RemoveIndexFile Method

```go
// RemoveIndexFile removes the package index
// Returns *PackageError on failure
func (p *Package) RemoveIndexFile() error
```

#### 5.3.5 Package.HasIndexFile Method

```go
// HasIndexFile checks if package has an index file
func (p *Package) HasIndexFile() bool
```

**Purpose**: Provides file navigation and indexing including:

- File location mappings
- Content-based indexing
- Search and navigation data
- File relationship mappings

### 5.4 Package Signature File (Type 65003)

This section describes the package signature file type and operations.

#### 5.4.1 Package.AddSignatureFile Method

```go
// AddSignatureFile adds a digital signature file
// Returns *PackageError on failure
func (p *Package) AddSignatureFile(signature SignatureData) error
```

#### 5.4.2 Package.GetSignatureFile Method

```go
// GetSignatureFile retrieves the signature file
// Returns *PackageError on failure
func (p *Package) GetSignatureFile() (SignatureData, error)
```

#### 5.4.3 Package.UpdateSignatureFile Method

```go
// UpdateSignatureFile updates the signature file
// Returns *PackageError on failure
func (p *Package) UpdateSignatureFile(updates SignatureData) error
```

#### 5.4.4 Package.RemoveSignatureFile Method

```go
// RemoveSignatureFile removes the signature file
// Returns *PackageError on failure
func (p *Package) RemoveSignatureFile() error
```

#### 5.4.5 Package.HasSignatureFile Method

```go
// HasSignatureFile checks if package has a signature file
func (p *Package) HasSignatureFile() bool
```

**Purpose**: Contains digital signature information including:

- Signature metadata and timestamps
- Public key information
- Signature validation data
- Trust chain information

### 5.5 Special File Management

This section describes special file management operations.

#### 5.5.1 Package.GetSpecialFiles Method

```go
// GetSpecialFiles returns all special files in the package
func (p *Package) GetSpecialFiles() ([]SpecialFileInfo, error)
```

#### 5.5.2 Package.GetSpecialFileByType Method

```go
// GetSpecialFileByType retrieves special file by type
func (p *Package) GetSpecialFileByType(fileType FileType) (SpecialFileInfo, error)
```

#### 5.5.3 Package.RemoveSpecialFile Method

```go
// RemoveSpecialFile removes a special file by type
// Returns *PackageError on failure
func (p *Package) RemoveSpecialFile(fileType FileType) error
```

#### 5.5.4 Package.ValidateSpecialFiles Method

```go
// ValidateSpecialFiles validates all special files
// Returns *PackageError on failure
func (p *Package) ValidateSpecialFiles() error
```

#### 5.5.5 Special File Data Structures

This section describes special file data structures used in the API.

##### 5.5.5.1 SpecialFileInfo Structure

```go
// SpecialFileInfo contains information about special metadata files in the package.
type SpecialFileInfo struct {
    Type        FileType // File type (see [File Types System - Special Files](file_type_system.md#339-special-file-types-65000-65535))
    Name        string   // Special file name (e.g., "__NVPK_META_240__.nvpkmeta")
    Size        int64    // File size in bytes
    Offset      int64    // Offset in package
    Data        []byte   // File content
    Valid       bool     // Whether file is valid
    Error       string   // Error message if invalid
}
```

##### 5.5.5.2 ManifestData Structure

```go
// ManifestData contains manifest file data structure.
type ManifestData struct {
    Version     string            // Manifest version
    Package     PackageInfo       // Package information
    Dependencies []Dependency     // Package dependencies
    Structure   []FileStructure   // File organization
    Install     InstallInfo       // Installation instructions
}
```

##### 5.5.5.3 IndexData Structure

```go
// IndexData contains index file data structure.
type IndexData struct {
    Version     string            // Index version
    Files       []FileIndex       // File index entries
    Navigation  NavigationData    // Navigation structure
    Search      SearchIndex       // Search index data
}
```

##### 5.5.5.4 SignatureData Structure

```go
// SignatureData contains signature file data structure.
type SignatureData struct {
    Version     string            // Signature format version
    Signatures  []SignatureInfo   // Signature information
    TrustChain  []TrustInfo       // Trust chain data
    Validation  ValidationData    // Validation metadata
}
```

## 6. Metadata-Only Packages

Metadata-only packages are NovusPack packages that contain no regular content files (FileCount = 0).
These packages may contain special metadata files (see [File Types System - Special Files](file_type_system.md#339-special-file-types-65000-65535)) or may be completely empty.
They serve specific purposes in package management and distribution systems, including placeholders and namespace reservations.

### 6.1 Metadata-Only Package Definition

A metadata-only package is defined as:

- **FileCount = 0**: No regular content files
- **IsMetadataOnly flag = 1**: Package header Bit 7 (metadata-only) is set
- **TotalSize = 0**: No uncompressed content data

Metadata-only packages may or may not contain special metadata files.
Packages with 0 files but no special metadata files are valid and represent empty or placeholder packages.

### 6.2 Valid Use Cases

This section describes valid use cases for metadata-only packages.

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

#### 6.2.5 Empty and Placeholder Packages

- **Empty packages**: Packages with no files or metadata, used as placeholders
- **Reserved packages**: Packages claiming a namespace or identifier
- **Future expansion**: Packages created for future content addition

### 6.3 Security Considerations

This section describes security considerations for metadata-only packages.

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

This section describes the metadata-only package API.

#### 6.4.1 Package.IsMetadataOnlyPackage Method

```go
// IsMetadataOnlyPackage checks if package contains only metadata files
func (p *Package) IsMetadataOnlyPackage() bool
```

#### 6.4.2 Package.AddMetadataOnlyFile Method

```go
// AddMetadataOnlyFile adds a special metadata file to a metadata-only package
// Returns *PackageError on failure
func (p *Package) AddMetadataOnlyFile(fileType FileType, data []byte) error
```

#### 6.4.3 Package.GetMetadataOnlyFiles Method

```go
// GetMetadataOnlyFiles returns all metadata files in the package
// Returns *PackageError on failure
func (p *Package) GetMetadataOnlyFiles() ([]SpecialFileInfo, error)
```

#### 6.4.4 Package.ValidateMetadataOnlyIntegrity Method

```go
// ValidateMetadataOnlyIntegrity validates metadata-only package integrity
// Returns *PackageError on failure
func (p *Package) ValidateMetadataOnlyIntegrity() error
```

See [6.4.5 Package.ValidateMetadataOnlyPackage Method](#645-packagevalidatemetadataonlypackage-method) for `ValidateMetadataOnlyPackage`.

#### 6.4.5 Package.ValidateMetadataOnlyPackage Method

```go
// ValidateMetadataOnlyPackage performs comprehensive validation:
//     Ensure FileCount == 0
//     Ensure IsMetadataOnly flag (Bit 7) is set in header
//     Validate all special metadata files (if present)
//     Check for malicious metadata patterns (if metadata files present)
//     Verify signature scope includes all metadata (if signatures present)
//     Ensure metadata consistency
//     Validate package structure
func (p *Package) ValidateMetadataOnlyPackage() error
```

#### 6.4.6 Write Operation Requirements

Write operations automatically determine and set the metadata-only flag based on package state:

- **Automatic Detection**: Write operations automatically detect metadata-only status when FileCount = 0 (no regular content files)
- **Header Flag**: Write operations MUST set the metadata-only flag (Bit 7) when FileCount = 0
- **Flag Synchronization**: PackageInfo.IsMetadataOnly MUST be synchronized to header flag before writing
- **Validation**: Write operations MUST validate that FileCount = 0 matches IsMetadataOnly flag state
- **Empty Packages**: Empty packages with no special metadata files are valid and automatically marked as metadata-only

Note: There is no separate function to create metadata-only packages.
Metadata-only status is determined automatically during write operations based on FileCount.
To create a metadata-only package, create a package with no regular files using standard package creation methods.

#### 6.4.7 Security Considerations for Metadata-Only Packages

- **Optional signatures**: Signatures are recommended but not mandatory for metadata-only packages
- **Enhanced validation**: Stricter validation requirements for packages with special metadata files
- **Trust verification**: Higher trust requirements for metadata-only packages with special metadata
- **Audit logging**: Enhanced logging for metadata-only package operations with special metadata

## 7. Package Information Structures

This section describes package information structures used in the API.

### 7.1 PackageInfo Structure

The PackageInfo structure provides comprehensive package information and metadata:

```go
// PackageInfo contains comprehensive package information and metadata.
type PackageInfo struct {
    // Basic Package Information
    FormatVersion         uint32    // Package file format version
    FileCount             int       // Number of regular content files in the package (excludes special metadata files, types 65000-65535)
    FilesUncompressedSize int64     // Total uncompressed size of all regular files (excludes special metadata files)
    FilesCompressedSize   int64     // Total compressed size of all regular files (excludes special metadata files)

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
    IsImmutable    bool      // Whether package is immutable (signed)

    // Timestamps
    Created        time.Time // Package creation timestamp
    Modified       time.Time // Package modification timestamp

    // Version Tracking
    PackageDataVersion uint32 // Tracks changes to package data content (file additions, removals, data modifications)
    MetadataVersion    uint32 // Tracks changes to package metadata (comment and identity changes)

    // Package Features
    HasMetadataFiles    bool      // Whether package has metadata files
    HasPerFileTags      bool      // Whether package has per-file tags (path metadata with properties)
    HasExtendedAttrs    bool      // Whether package has extended attributes (filesystem metadata)
    HasEncryptedData    bool      // Whether package contains encrypted files
    HasCompressedData   bool      // Whether package contains compressed files
    IsMetadataOnly      bool      // Whether package contains no regular files (FileCount = 0)

    // Package Compression
    PackageCompression  uint8     // Package compression type (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    IsPackageCompressed bool      // Whether the entire package is compressed
    PackageOriginalSize int64     // Original package size before compression (0 if not compressed)
    PackageCompressedSize int64   // Compressed package size (0 if not compressed)
    PackageCompressionRatio float64 // Compression ratio (0.0-1.0, 0.0 if not compressed)
}
```

#### 7.1.1 PackageInfo Scope and Exclusions

`PackageInfo` provides lightweight package-level information and does NOT include:

- Individual `FileEntry` metadata (use `ListFiles()` or `GetMetadata()` for file-level details)
- Special metadata file contents (use `GetMetadata()` for special file data)
- Path metadata entries (use `GetMetadata()` for path metadata)
- Full signature data (only signature summary information is included; see `Signatures` field for `SignatureInfo` details)

For comprehensive metadata including file entries and special metadata files, use `GetMetadata()` which returns `PackageMetadata` containing `PackageInfo` plus all file and metadata details.

#### 7.1.2 NewPackageInfo Function

Creates a new PackageInfo with proper initialization.

```go
// NewPackageInfo creates a new PackageInfo with default values
func NewPackageInfo() *PackageInfo
```

Returns a new PackageInfo instance with all fields initialized to their zero values or specification defaults:

- All numeric fields set to 0
- All boolean fields set to false
- String fields set to empty strings
- Slice fields initialized to empty slices
- Time fields set to zero time (will be set on package creation)

There is no package-wide security level.

This ensures that all PackageInfo initialization is centralized and aligned with specification defaults.

**Note**: PackageInfo instances are typically created by the package system during package creation or opening. This factory function provides a consistent way to initialize PackageInfo instances when needed.

#### 7.1.3. PackageInfo As Source of Truth

PackageInfo serves different roles depending on the operation phase:

1. **During Package Open/Read Operations**:

   - Package header is read from disk (header is the on-disk source of truth)
   - Header flags and metadata are used to populate PackageInfo fields
   - PackageInfo is synchronized from header state

2. **During In-Memory Operations**:

   - **PackageInfo is the source of truth** for all package metadata and flags
   - All package operations read from and update PackageInfo
   - Header fields may become stale during in-memory operations

3. **During Package Write/Save Operations**:
   - The header is synchronized from PackageInfo before serialization.
   - PackageInfo values are written to header fields.
   - PackageInfo versions are written to the header version fields.

**Flag Synchronization**: Functions like `UpdateSpecialMetadataFlags()` currently update both header flags and PackageInfo fields simultaneously during in-memory operations. This ensures consistency, but the architectural intent is that PackageInfo should be the source of truth during in-memory operations, with the header being updated from PackageInfo during write operations.

#### 7.1.4 PackageInfo.FromHeader Method

Synchronizes a PackageInfo instance from an on-disk PackageHeader.

This helper provides the standard, centralized mapping from header fields and flags to PackageInfo fields.

```go
// FromHeader synchronizes PackageInfo fields from the provided PackageHeader.
//
// This method must only copy data that is represented in the header.
// It must not compute derived values that require scanning file entries or reading file data.
//
// Returns *PackageError on failure.
func (pi *PackageInfo) FromHeader(header *PackageHeader) error
```

FromHeader must synchronize at least the following fields:

- FormatVersion from header.FormatVersion.
- VendorID from header.VendorID.
- AppID from header.AppID.
- PackageDataVersion from header.PackageDataVersion.
- MetadataVersion from header.MetadataVersion.
- Created from header.CreatedTime.
- Modified from header.ModifiedTime.
- HasComment from header.CommentSize (or equivalent comment presence indicator).
- HasSignatures and IsImmutable from header.SignatureOffset (signed packages are immutable).
- PackageCompression and IsPackageCompressed from the compression type encoded in header.Flags (bits 15-8).
- Feature flags derived from header.Flags (bits 0-7), including HasMetadataFiles, HasPerFileTags, HasExtendedAttrs, HasEncryptedData, HasCompressedData, and IsMetadataOnly.

#### 7.1.5 PackageHeader Structure

```go
// PackageHeader represents the fixed-size header of a NovusPack (.nvpk) file
// Size: 112 bytes (fixed)
type PackageHeader struct {
    Magic              uint32 // Package identifier (0x4E56504B "NVPK")
    FormatVersion      uint32 // Format version (current: 1)
    Flags              uint32 // Package-level features and options
    PackageDataVersion uint32 // Tracks changes to package data content
    MetadataVersion    uint32 // Tracks changes to package metadata
    PackageCRC         uint32 // CRC32 of package content (0 if skipped)
    CreatedTime        uint64 // Package creation timestamp (Unix nanoseconds)
    ModifiedTime       uint64 // Package modification timestamp (Unix nanoseconds)
    LocaleID           uint32 // Locale identifier for path encoding
    Reserved           uint32 // Reserved for future use (must be 0)
    AppID              uint64 // Application/game identifier (0 if not associated)
    VendorID           uint32 // Storefront/platform identifier (0 if not associated)
    CreatorID          uint32 // Creator identifier (reserved for future use)
    IndexStart         uint64 // Offset to file index from start of file
    IndexSize          uint64 // Size of file index in bytes
    ArchiveChainID     uint64 // Archive chain identifier
    ArchivePartInfo    uint32 // Combined part number and total parts
    CommentSize        uint32 // Size of package comment in bytes (0 if no comment)
    CommentStart       uint64 // Offset to package comment from start of file
    SignatureOffset    uint64 // Offset to signatures block from start of file
}
```

**Purpose**: Provides comprehensive metadata and navigation information for the entire package.

**Cross-Reference**: See [Package File Format - Package Header](package_file_format.md#2-package-header) for complete header specification.

#### 7.1.6 PackageHeader.ToHeader Method

Synchronizes an on-disk PackageHeader instance from an in-memory PackageInfo instance.

This helper provides the standard, centralized mapping from PackageInfo fields to the header fields and flags written to disk.

```go
// ToHeader synchronizes PackageHeader fields from the provided PackageInfo.
//
// This method must only write fields that are represented in the header.
// It must not mutate fields that are computed by the writer pipeline (for example IndexStart, IndexSize, and CRC).
//
// Returns *PackageError on failure.
func (h *PackageHeader) ToHeader(pi *PackageInfo) error
```

ToHeader must synchronize at least the following fields:

- FormatVersion from pi.FormatVersion.
- VendorID from pi.VendorID.
- AppID from pi.AppID.
- PackageDataVersion from pi.PackageDataVersion.
- MetadataVersion from pi.MetadataVersion.
- CreatedTime from pi.Created.
- ModifiedTime from pi.Modified.
- Compression type bits (15-8) from pi.PackageCompression.
- Feature flags (bits 0-7) from the corresponding PackageInfo booleans.
- CommentSize and CommentStart must reflect the serialized comment section computed by the writer.
- SignatureOffset must reflect the serialized signatures section computed by the signer.

### 7.2 SignatureInfo Structure

```go
// SignatureInfo contains signature information for a package.
type SignatureInfo struct {
    Index         int       // Signature index in the package
    Type          uint32    // Signature type (ML-DSA, SLH-DSA, PGP, X.509)
    Size          uint32    // Size of signature data in bytes
    Offset        uint64    // Offset to signature data from start of file
    Flags         uint32    // Signature-specific flags
    Timestamp     uint32    // Unix timestamp when signature was created
    Comment       string    // Signature comment (if any)
    Algorithm     string    // Algorithm name/description
    SecurityLevel int       // Algorithm security level (v2, signature algorithm specific)
    Valid         bool      // Whether signature is valid
    Trusted       bool      // Whether signature is trusted
    Error         string    // Error message if validation failed
}
```

### 7.3 SecurityStatus Structure

```go
// SecurityStatus contains the security status of a package.
type SecurityStatus struct {
    SignatureCount      int                           // Number of signatures
    ValidSignatures     int                           // Number of valid signatures
    TrustedSignatures   int                           // Number of trusted signatures
    SignatureResults    []SignatureValidationResult   // Individual results
    HasChecksums        bool                          // Checksums present
    ChecksumsValid      bool                          // Checksums valid
    ValidationErrors    []string                      // Validation errors
}
```

### 7.4 Package Information Methods

This section describes package information methods.

#### 7.4.1 Package.GetPackageInfo Method

```go
// GetPackageInfo returns comprehensive package information
func (p *Package) GetPackageInfo() (*PackageInfo, error)
```

#### 7.4.2 Package.RefreshPackageInfo Method

```go
// RefreshPackageInfo refreshes package information from the file on-disk
// Returns *PackageError on failure
func (p *Package) RefreshPackageInfo(ctx context.Context) error
```

See [GetSecurityStatus](api_security.md#11-multiple-signature-validation-incremental) for the security status method.

## 8. PathMetadata System

The path metadata system provides structured path definitions (files, directories, and symlinks) with metadata, inheritance rules, and filesystem properties for the NovusPack system using special metadata files.

**Tag Inheritance Model**: Tag inheritance only works with `PathMetadataEntry` instances, not `FileEntry` instances directly.
This is because `FileEntry` can have multiple paths (via `FileEntry.Paths`), while `PathMetadataEntry` represents a single path.
Inheritance is resolved per-path by walking up the `PathMetadataEntry.ParentPath` chain.
Each `PathMetadataEntry` can have its own inheritance chain, allowing different paths for the same file content to inherit different tags.

**Cross-Reference**: For file-path association methods and FileEntry path operations, see [FileEntry API - Path Management](api_file_mgmt_file_entry.md#5-path-management).

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

### 8.1 PathMetadata Structures

**Note**: PathEntry is defined in the generics package and is shared by both FileEntry and PathMetadataEntry.
See [Generic Types and Patterns - PathEntry](api_generics.md#13-pathentry-type) for complete specification.

#### 8.1.1 PathMetadataType Type

```go
// PathMetadataType represents the type of path entry
type PathMetadataType uint8

const (
    PathMetadataTypeFile            PathMetadataType = 0 // Regular file
    PathMetadataTypeDirectory       PathMetadataType = 1 // Regular directory
    PathMetadataTypeFileSymlink     PathMetadataType = 2 // Symlink to a file
    PathMetadataTypeDirectorySymlink PathMetadataType = 3 // Symlink to a directory
)
```

#### 8.1.2 PathMetadataEntry Structure

```go
// PathMetadataEntry represents a path (file, directory, or symlink) with metadata, inheritance rules, and filesystem properties
type PathMetadataEntry struct {
    // Path is the path entry (minimal PathEntry from generics package)
    // All paths are stored with a leading "/" to indicate the package root
    // - Directory paths must end with "/" (e.g., "/assets/")
    // - File paths must NOT end with "/" (e.g., "/assets/file.txt")
    // - Root path is represented as "/" (the package root)
    // - All paths use forward slashes ("/") as separators
    Path        generics.PathEntry         `yaml:"path"`

    // Type indicates whether this is a file, directory, file symlink, or directory symlink
    Type        PathMetadataType           `yaml:"type"`

    // Tags are path-specific tags (typed tags)
    Tags  []*generics.Tag[any]       `yaml:"tags"`

    // Inheritance controls tag inheritance behavior (optional, only for directories)
    Inheritance *PathInheritance           `yaml:"inheritance,omitempty"` // Inheritance settings (nil for files)

    // Metadata contains path metadata (optional, only for directories)
    Metadata    *PathMetadata              `yaml:"metadata,omitempty"`    // Path metadata (nil for files)

    // Destination extraction override paths.
    //
    // These fields define optional destination extraction directory overrides for this path.
    // They are interpreted by ExtractPath destination resolution.
    //
    // Destinations may be absolute or relative.
    // Relative destinations may include "." and ".." segments.
    //
    // Relative destinations are resolved relative to the default extraction directory for the path.
    // This is the directory the path would be extracted to under session base with no overrides.
    //
    // Cross-Reference: See [File Extraction API - ExtractPath Destination Resolution](api_file_mgmt_extraction.md#151-extractpath-destination-resolution).
    DestPath    string                     `yaml:"dest_path,omitempty"`
    DestPathWin string                     `yaml:"dest_path_win,omitempty"`

    // FileSystem contains filesystem-specific properties
    FileSystem  PathFileSystem             `yaml:"filesystem"`

    // Path hierarchy (runtime only, not stored in file)
    ParentPath  *PathMetadataEntry         `yaml:"-"` // Pointer to parent path (nil for root)

    // FileEntry associations (runtime only, not stored in file)
    // Since FileEntry and PathMetadataEntry are in the same package, direct references are used.
    AssociatedFileEntries []*FileEntry `yaml:"-"` // FileEntry instances associated with this path
}
```

#### 8.1.3 PathInheritance Structure

```go
// PathInheritance controls tag inheritance behavior (for directories only)
type PathInheritance struct {
    Enabled  bool `yaml:"enabled"`  // Whether this path provides inheritance
    Priority int  `yaml:"priority"` // Inheritance priority (higher = more specific)
}
```

#### 8.1.4 PathMetadata Structure

```go
// PathMetadata contains path metadata (for directories only)
type PathMetadata struct {
    Created     string `yaml:"created"`     // Path creation time (ISO8601)
    Modified    string `yaml:"modified"`    // Last modification time (ISO8601)
    Description string `yaml:"description"` // Human-readable description
}
```

#### 8.1.5 PathFileSystem Structure

```go
// PathFileSystem contains filesystem-specific properties
type PathFileSystem struct {
    // Execute permissions (always captured)
    IsExecutable bool           `yaml:"is_executable"`     // Whether file has any execute permission bits set (tracked by default)

    // Unix/Linux properties (optional, captured when PreservePermissions is enabled)
    Mode    *uint32            `yaml:"mode,omitempty"`    // File/directory permissions and type (Unix-style)
    UID     *uint32            `yaml:"uid,omitempty"`     // User ID
    GID     *uint32            `yaml:"gid,omitempty"`     // Group ID
    ACL     []ACLEntry         `yaml:"acl,omitempty"`     // Access Control List

    // Timestamps (Unix nanoseconds since epoch)
    ModTime    uint64         `yaml:"mod_time,omitempty"`    // Modification time
    CreateTime uint64         `yaml:"create_time,omitempty"` // Creation time
    AccessTime uint64         `yaml:"access_time,omitempty"` // Access time

    // Symbolic link support
    LinkTarget string         `yaml:"link_target,omitempty"` // Target path for symbolic links (empty if not a symlink)

    // Windows properties (optional)
    WindowsAttrs *uint32      `yaml:"windows_attrs,omitempty"` // Windows attributes

    // Extended attributes (optional)
    ExtendedAttrs map[string]string `yaml:"extended_attrs,omitempty"` // Extended attributes

    // Filesystem flags (optional)
    Flags   *uint16            `yaml:"flags,omitempty"`   // Filesystem-specific flags
}
```

#### 8.1.6 ACLEntry Structure

```go
// ACLEntry represents an Access Control List entry
type ACLEntry struct {
    Type    string `yaml:"type"`    // "user", "group", "other", "mask"
    ID      *uint32 `yaml:"id,omitempty"` // User/Group ID (nil for "other")
    Perms   string  `yaml:"perms"`  // Permissions (e.g., "rwx", "r--")
}
```

#### 8.1.7 PathMetadataEntry Tag Management

All tag operations use typed tags for type safety.
See [FileEntry API - Tag Management](api_file_mgmt_file_entry.md#3-tag-management) for complete tag system documentation.

**Tag Management Architecture**: Tags are managed directly on `PathMetadataEntry` and `FileEntry` instances.
There are no Package-level tag management functions.

When deciding where to set tags, consider the following guidance:

- Set tags on `FileEntry` when the same tags need to apply to all paths for a file with multiple paths.
- Set tags on `PathMetadataEntry` when tags should only apply to a specific path.
- `PathMetadataEntry` tags can participate in inheritance; `FileEntry` tags do not.

**Note**: These are standalone functions rather than methods due to Go's limitation of not supporting generic methods on non-generic types.
See [Generic Types and Patterns](api_generics.md) for details.

##### 8.1.7.1 GetPathMetaTags Function

```go
// GetPathMetaTags returns all tags as typed tags for a PathMetadataEntry
// Returns *PackageError on failure
func GetPathMetaTags(pme *PathMetadataEntry) ([]*Tag[any], error)
```

##### 8.1.7.2 GetPathMetaTagsByType Function

```go
// GetPathMetaTagsByType returns all tags of a specific type for a PathMetadataEntry
// Returns a slice of Tag pointers with the specified type parameter T
// Only tags matching the type T and corresponding TagValueType are returned
// Returns *PackageError on failure
func GetPathMetaTagsByType[T any](pme *PathMetadataEntry) ([]*Tag[T], error)
```

##### 8.1.7.3 AddPathMetaTags Function

```go
// AddPathMetaTags adds multiple new tags with type safety to a PathMetadataEntry
// Returns *PackageError if any tag with the same key already exists
func AddPathMetaTags(pme *PathMetadataEntry, tags []*Tag[any]) error
```

##### 8.1.7.4 SetPathMetaTags Function

```go
// SetPathMetaTags updates existing tags from a slice of typed tags for a PathMetadataEntry
// Returns *PackageError if any tag key does not already exist
// Only modifies tags that already exist; does not create new tags
func SetPathMetaTags(pme *PathMetadataEntry, tags []*Tag[any]) error
```

##### 8.1.7.5 GetPathMetaTag Function

```go
// GetPathMetaTag retrieves a type-safe tag by key from a PathMetadataEntry
// Returns the tag pointer and an error. If the tag is not found, returns (nil, nil).
// If an underlying error occurs, returns (nil, error).
// Returns *PackageError on failure
// If the tag type is unknown, use GetPathMetaTag[any](pme, "key") to retrieve the tag and inspect its Type field
func GetPathMetaTag[T any](pme *PathMetadataEntry, key string) (*Tag[T], error)
```

##### 8.1.7.6 AddPathMetaTag Function

```go
// AddPathMetaTag adds a new tag with type safety to a PathMetadataEntry
// Returns *PackageError if a tag with the same key already exists
func AddPathMetaTag[T any](pme *PathMetadataEntry, key string, value T, tagType TagValueType) error
```

##### 8.1.7.7 SetPathMetaTag Function

```go
// SetPathMetaTag updates an existing tag with type safety for a PathMetadataEntry
// Returns *PackageError if the tag key does not already exist
// Only modifies existing tags; does not create new tags
func SetPathMetaTag[T any](pme *PathMetadataEntry, key string, value T, tagType TagValueType) error
```

##### 8.1.7.8 RemovePathMetaTag Function

```go
// RemovePathMetaTag removes a tag by key from a PathMetadataEntry
// Returns *PackageError on failure
func RemovePathMetaTag(pme *PathMetadataEntry, key string) error
```

##### 8.1.7.9 HasPathMetaTag Function

```go
// HasPathMetaTag checks if a tag with the specified key exists on a PathMetadataEntry
func HasPathMetaTag(pme *PathMetadataEntry, key string) bool
```

#### 8.1.8 PathMetadataEntry Methods

This section describes PathMetadataEntry methods.

##### 8.1.8.1 PathMetadataEntry.SetPath Method

```go
// Path management methods for PathMetadataEntry
func (pme *PathMetadataEntry) SetPath(path string)
```

##### 8.1.8.2 PathMetadataEntry.GetPath Method

```go
// GetPath returns the path as stored (Unix-style with forward slashes).
// For platform-specific display, use GetPathForPlatform() or convert manually.
func (pme *PathMetadataEntry) GetPath() string
```

##### 8.1.8.3 PathMetadataEntry.GetPathForPlatform Method

```go
// GetPathForPlatform returns the path converted for the specified platform.
// On Windows: converts forward slashes to backslashes.
// On Unix/Linux: returns the path as stored (with forward slashes).
func (pme *PathMetadataEntry) GetPathForPlatform(isWindows bool) string
```

##### 8.1.8.4 PathMetadataEntry.GetPathEntry Method

```go
// GetPathEntry returns the PathEntry representation of this path metadata entry.
func (pme *PathMetadataEntry) GetPathEntry() generics.PathEntry
```

##### 8.1.8.5 PathMetadataEntry.GetType Method

```go
// Type and symlink methods for PathMetadataEntry
func (pme *PathMetadataEntry) GetType() PathMetadataType
```

##### 8.1.8.6 PathMetadataEntry.IsDirectory Method

```go
// IsDirectory returns true if this path metadata entry represents a directory.
func (pme *PathMetadataEntry) IsDirectory() bool
```

##### 8.1.8.7 PathMetadataEntry.IsFile Method

```go
// IsFile returns true if this path metadata entry represents a file.
func (pme *PathMetadataEntry) IsFile() bool
```

##### 8.1.8.8 PathMetadataEntry.IsSymlink Method

```go
// IsSymlink returns true if this path metadata entry represents a symlink.
func (pme *PathMetadataEntry) IsSymlink() bool
```

##### 8.1.8.9 PathMetadataEntry.GetLinkTarget Method

```go
// GetLinkTarget returns the target path of the symlink.
func (pme *PathMetadataEntry) GetLinkTarget() string
```

##### 8.1.8.10 PathMetadataEntry.ResolveSymlink Method

```go
// ResolveSymlink resolves the symlink to its final target path.
func (pme *PathMetadataEntry) ResolveSymlink() string
```

##### 8.1.8.11 PathMetadataEntry.SetParentPath Method

```go
// Parent path management methods for PathMetadataEntry
func (pme *PathMetadataEntry) SetParentPath(parent *PathMetadataEntry)
```

##### 8.1.8.12 PathMetadataEntry.GetParentPath Method

```go
// GetParentPath returns the parent path metadata entry.
func (pme *PathMetadataEntry) GetParentPath() *PathMetadataEntry
```

##### 8.1.8.13 PathMetadataEntry.GetParentPathString Method

```go
// GetParentPathString returns the parent path as a string.
func (pme *PathMetadataEntry) GetParentPathString() string
```

##### 8.1.8.14 PathMetadataEntry.GetDepth Method

```go
// GetDepth returns the depth of this path in the directory hierarchy.
func (pme *PathMetadataEntry) GetDepth() int
```

##### 8.1.8.15 PathMetadataEntry.IsRoot Method

```go
// IsRoot returns true if this path metadata entry represents the root path.
func (pme *PathMetadataEntry) IsRoot() bool
```

##### 8.1.8.16 PathMetadataEntry.GetAncestors Method

```go
// GetAncestors returns all ancestor path metadata entries up to the root.
func (pme *PathMetadataEntry) GetAncestors() []*PathMetadataEntry
```

##### 8.1.8.17 PathMetadataEntry.GetInheritedTags Method

```go
// Tag inheritance methods for PathMetadataEntry
// These methods resolve inheritance by walking up the ParentPath chain
func (pme *PathMetadataEntry) GetInheritedTags() ([]*Tag[any], error)
```

##### 8.1.8.18 PathMetadataEntry.GetEffectiveTags Method

```go
// GetEffectiveTags returns all tags for this PathMetadataEntry, including:
// 1. Tags directly on this PathMetadataEntry
// 2. Tags inherited from parent PathMetadataEntry instances (path hierarchy)
// 3. Tags from associated FileEntry instances (treated as if applied to this PathMetadataEntry)
func (pme *PathMetadataEntry) GetEffectiveTags() ([]*Tag[any], error)
```

##### 8.1.8.19 PathMetadataEntry.AssociateWithFileEntry Method

```go
// FileEntry association methods for PathMetadataEntry
// AssociateWithFileEntry associates this PathMetadataEntry with a FileEntry
// The association is established if the PathMetadataEntry.Path.Path matches one of the FileEntry.Paths
// Returns *PackageError on failure
func (pme *PathMetadataEntry) AssociateWithFileEntry(fe *FileEntry) error
```

##### 8.1.8.20 PathMetadataEntry.GetAssociatedFileEntries Method

```go
// GetAssociatedFileEntries returns all FileEntry instances associated with this PathMetadataEntry
// Returns empty slice if no FileEntry instances are associated
func (pme *PathMetadataEntry) GetAssociatedFileEntries() []*FileEntry
```

#### 8.1.9 PathMetadataEntry Validation Methods

This section defines validation behavior for PathMetadataEntry instances.

##### 8.1.9.1 PathMetadataEntry.Validate Method

```go
// Validate validates the PathMetadataEntry state and returns an error on failure.
func (pme *PathMetadataEntry) Validate() error
```

#### 8.1.10 PathInfo Structure

```go
// PathInfo provides runtime path metadata information
type PathInfo struct {
    Entry      PathMetadataEntry // Path metadata entry data
    FileCount  int           // Number of files in this path
    SubDirs    []string      // Immediate subdirectories
    ParentPath string        // Parent path
    Depth      int           // Path depth (0 = root)
}
```

#### 8.1.11 FilePathAssociation Structure

```go
// FilePathAssociation links files to their path metadata
type FilePathAssociation struct {
    FilePath     string        // File path
    Path         string        // Parent path
    PathMetadata *PathInfo     // Path metadata information (nil if no path metadata)
    InheritedTags []Tag        // Tags inherited from path hierarchy
    EffectiveTags []Tag        // All tags including inheritance
}
```

#### 8.1.12 DestPathOverride Structure

```go
// DestPathOverride specifies destination extraction directory overrides.
//
// A nil field means "no override specified" for that field.
type DestPathOverride struct {
    DestPath    *string
    DestPathWin *string
}
```

#### 8.1.13 DestPathInput Interface

```go
// DestPathInput is the allowed input type set for SetDestPath.
//
// DestPathInput supports:
// - string: a single destination string
// - map[string]string: a map with keys "DestPath" and/or "DestPathWin"
//
// Note: The map form uses string keys for ergonomics in callers.
// Keys other than "DestPath" and "DestPathWin" MUST be rejected with ErrTypeValidation.
type DestPathInput interface {
    ~string | ~map[string]string
}
```

### 8.2 PathMetadata Management Methods

The PathMetadata management API provides methods for creating, reading, updating, and deleting path metadata entries.

#### 8.2.1 Core PathMetadata CRUD Operations

This section describes core CRUD operations for PathMetadata.

##### 8.2.1.1 Package.GetPathMetadata Method

```go
// GetPathMetadata retrieves all path metadata entries from the package
// Returns *PackageError on failure
func (p *Package) GetPathMetadata(ctx context.Context) ([]*PathMetadataEntry, error)
```

Returns all path metadata entries currently stored in the package.

##### 8.2.1.2 Package.SetPathMetadata Method

```go
// SetPathMetadata replaces all path metadata entries in the package
// Returns *PackageError on failure
func (p *Package) SetPathMetadata(ctx context.Context, entries []*PathMetadataEntry) error
```

Replaces all existing path metadata entries with the provided entries.

##### 8.2.1.3 Package.AddPathMetadata Method

```go
// AddPathMetadata adds a new path metadata entry to the package
// Returns *PackageError on failure
func (p *Package) AddPathMetadata(ctx context.Context, path string, pathType PathMetadataType, properties map[string]string, inheritance *PathInheritance, metadata *PathMetadata) error
```

Adds a new path metadata entry with the specified path, type, properties, inheritance rules, and metadata.

##### 8.2.1.4 Package.RemovePathMetadata Method

```go
// RemovePathMetadata removes a path metadata entry by path
// Returns *PackageError on failure
func (p *Package) RemovePathMetadata(ctx context.Context, path string) error
```

Removes the path metadata entry for the specified path.

##### 8.2.1.5 Package.UpdatePathMetadata Method

```go
// UpdatePathMetadata updates an existing path metadata entry
// Returns *PackageError on failure
func (p *Package) UpdatePathMetadata(ctx context.Context, path string, properties map[string]string, inheritance *PathInheritance, metadata *PathMetadata) error
```

Updates the path metadata entry for the specified path with new properties, inheritance rules, and metadata.

##### 8.2.1.6 Package.SetDestPath Method

```go
// SetDestPath sets destination extraction directory overrides for a stored path.
//
// This is a pure in-memory operation.
//
// storedPath MUST be treated as a stored package path.
// If storedPath does not begin with "/", the implementation MUST prefix "/" before matching or creating entries.
//
// If no PathMetadataEntry exists for storedPath, SetDestPath MUST create one.
// The new entry type MUST be inferred from storedPath:
// - "/" is PathMetadataTypeDirectory
// - paths ending with "/" are PathMetadataTypeDirectory
// - all other paths are PathMetadataTypeFile
//
// Returns *PackageError on failure.
func (p *Package) SetDestPath(storedPath string, override DestPathOverride) error
```

Sets destination extraction directory overrides for a stored path.
This is a pure in-memory operation that creates a PathMetadataEntry if one doesn't exist.

##### 8.2.1.7 SetDestPath Function

```go
// SetDestPathTyped is a generic helper for SetDestPath.
//
// This helper exists to allow compile-time type checking of the dest input.
// It converts dest to DestPathOverride, then delegates to Package.SetDestPath.
//
// If dest is a string, it MUST be parsed to determine which destination field to set.
// If the string is a Windows-only absolute path (drive letter like "C:\\" or "C:/", or UNC path like "\\\\server\\share"),
// it MUST be stored as DestPathWin.
// Otherwise, it MUST be stored as DestPath.
func SetDestPath[T DestPathInput](p *Package, storedPath string, dest T) error
```

Generic helper for SetDestPath that provides compile-time type checking.

##### 8.2.1.8 Directory Metadata Convenience Methods

These methods are wrappers around [AddPathMetadata](#8213-packageaddpathmetadata-method), [RemovePathMetadata](#8214-packageremovepathmetadata-method), and [UpdatePathMetadata](#8215-packageupdatepathmetadata-method).
They specifically handle directory paths (paths ending with `/`).

##### 8.2.1.9 Package.AddDirectoryMetadata Method

```go
// AddDirectoryMetadata adds directory path metadata (metadata-only, does not add files)
// Returns *PackageError on failure
func (p *Package) AddDirectoryMetadata(ctx context.Context, path string, properties map[string]string, inheritance *PathInheritance, metadata *PathMetadata) error
```

##### 8.2.1.10 Package.RemoveDirectoryMetadata Method

```go
// RemoveDirectoryMetadata removes directory path metadata (metadata-only, does not remove files)
// Returns *PackageError on failure
func (p *Package) RemoveDirectoryMetadata(ctx context.Context, path string) error
```

##### 8.2.1.11 Package.UpdateDirectoryMetadata Method

```go
// UpdateDirectoryMetadata updates directory path metadata (metadata-only, does not modify files)
// Returns *PackageError on failure
func (p *Package) UpdateDirectoryMetadata(ctx context.Context, path string, properties map[string]string, inheritance *PathInheritance, metadata *PathMetadata) error
```

**Note**: These methods are metadata-only operations.
To add or remove files within a directory, use [AddDirectory](api_file_mgmt_addition.md#25-packageadddirectory-method) or [RemoveDirectory](api_file_mgmt_removal.md#42-packageremovedirectory-method) from the File Management API.

##### 8.2.1.12 Package.ValidatePathMetadata Method

```go
// ValidatePathMetadata validates all path metadata entries
// Returns *PackageError on failure
func (p *Package) ValidatePathMetadata() error
```

##### 8.2.1.13 Package.GetPathConflicts Method

```go
// GetPathConflicts returns a list of paths with conflicting metadata
// Returns *PackageError on failure
func (p *Package) GetPathConflicts() ([]string, error)
```

Validates path metadata consistency and identifies conflicts.

#### Path Information Query Methods

- [GetPathInfo](#822-packagegetpathinfo-method) - Returns complete path information for a specific path
- [ListPaths](#823-packagelistpaths-method) - Returns all paths in the package
- [ListDirectories](#824-packagelistdirectories-method) - Returns all directory paths
- [GetDirectoryCount](#825-packagegetdirectorycount-method) - Returns the total number of directories
- [GetPathHierarchy](#826-packagegetpathhierarchy-method) - Returns the path hierarchy as a map

#### Path Association Methods

- [AssociateFileWithPath](#827-packageassociatefilewithpath-method) - Links a file to its path metadata
- [UpdateFilePathAssociations](#828-packageupdatefilepathassociations-method) - Rebuilds all file-path associations

See [Path Association System](#84-path-association-system) for additional association methods including `DisassociateFileFromPath` and `GetFilePathAssociations`.

#### Special Metadata File Management

- [SavePathMetadataFile](#835-implementation-details) - Creates and saves the path metadata file
- [LoadPathMetadataFile](#835-implementation-details) - Loads and parses the path metadata file
- [UpdateSpecialMetadataFlags](#835-implementation-details) - Updates package header flags

See [Special Metadata File Management](#83-special-metadata-file-management) for details.

#### Special Metadata File Creation Helpers

##### 8.2.1.14 Package.CreateSpecialMetadataFile Method

```go
// CreateSpecialMetadataFile creates a special metadata FileEntry
// Returns *PackageError on failure
func (p *Package) CreateSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string, content []byte) (*FileEntry, error)
```

##### 8.2.1.15 Package.UpdateSpecialMetadataFile Method

```go
// UpdateSpecialMetadataFile updates an existing special metadata file
// Returns *PackageError on failure
func (p *Package) UpdateSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string, content []byte) error
```

##### 8.2.1.16 Package.RemoveSpecialMetadataFile Method

```go
// RemoveSpecialMetadataFile removes a special metadata file
// Returns *PackageError on failure
func (p *Package) RemoveSpecialMetadataFile(ctx context.Context, fileType uint16, fileName string) error
```

Helper methods for creating, updating, and removing special metadata files.

#### 8.2.2 Package.GetPathInfo Method

Returns complete path information including metadata, inheritance, and associations for a specific path.

##### 8.2.2.1 GetPathInfo Parameters

- `path`: Path string to query

##### 8.2.2.2 GetPathInfo Returns

- `*PathInfo`: Complete path information structure
- `error`: Any error that occurred during the query

#### 8.2.3 Package.ListPaths Method

Returns all paths in the package with their complete metadata.

##### 8.2.3.1 ListPaths Returns

- `[]PathInfo`: Slice of all path information entries
- `error`: Any error that occurred during the query

#### 8.2.4 Package.ListDirectories Method

Returns all directory paths in the package with their complete metadata.
Only returns paths that are marked as directories (paths ending with `/` or with directory metadata).

##### 8.2.4.1 ListDirectories Returns

- `[]PathInfo`: Slice of directory path information entries
- `error`: Any error that occurred during the query

##### 8.2.4.2 ListDirectories Use Cases

- List all directories in the package
- Query directory permissions and metadata
- Build directory tree visualizations
- Verify directory metadata consistency

#### 8.2.5 Package.GetDirectoryCount Method

Returns the total number of directories in the package.

##### 8.2.5.1 GetDirectoryCount Returns

- `int`: Total count of directory paths
- `error`: Any error that occurred during the query

##### 8.2.5.2 GetDirectoryCount Use Cases

- Package statistics
- Progress tracking during extraction
- Directory structure validation
- Quick checks without loading full path data

#### 8.2.6 Package.GetPathHierarchy Method

Returns the path hierarchy as a map of parent paths to their child paths.

##### 8.2.6.1 GetPathHierarchy Returns

- `map[string][]string`: Map of parent path to slice of child paths
- `error`: Any error that occurred during the query

##### 8.2.6.2 GetPathHierarchy Use Cases

- Build directory tree structures
- Resolve path inheritance hierarchies
- Validate path relationships
- Efficient traversal of path hierarchies

#### 8.2.7 Package.AssociateFileWithPath Method

**AssociateFileWithPath** links a file to its path metadata, establishing the association between a FileEntry and its corresponding PathMetadataEntry.

The function performs the following operations:

1. Validates the context for cancellation/timeout.
2. Retrieves the FileEntry by path using `findFileEntryByPath()`.
3. Retrieves the path metadata entry by path using `findPathMetadataByPath()`.
4. Establishes the association by calling `FileEntry.AssociateWithPathMetadata()` with the retrieved PathMetadataEntry.
5. Sets parent path association for hierarchy traversal:
   - Derives the parent path using `filepath.Dir` on the path (after normalizing trailing slashes)
   - Validates that the parent path is different from the path and not the current directory (`.`)
   - If valid, retrieves the parent path metadata entry and sets `PathMetadataEntry.ParentPath` to establish the hierarchy
   - **This operation never fails** - missing parents are valid for root paths or when path metadata is optional
   - `ParentPath` is a runtime-only convenience field (not persisted) that is `nil` for root paths
6. Does not modify inheritance or tags - those are handled separately on PathMetadataEntry when needed.

#### 8.2.8 Package.UpdateFilePathAssociations Method

**UpdateFilePathAssociations** rebuilds all file-path associations for all files in the package, ensuring that every file is properly linked to its path metadata.

The function performs the following operations:

1. Retrieves all files in the package using `ListFiles`.
2. Retrieves all path metadata entries using `GetPathMetadata`.
3. Builds a path metadata map for efficient lookup:
   - Creates a map keyed by path strings
   - Maps each path to its corresponding PathMetadataEntry pointer
4. Associates each file with its path metadata:
   - Iterates through all files and their associated paths (`FileEntry.Paths`)
   - For each path in `FileEntry.Paths`, extracts the path string and looks up the corresponding `PathMetadataEntry` in the map
   - If a matching `PathMetadataEntry` exists, calls `FileEntry.AssociateWithPathMetadata(pme)` to establish the association
   - This adds the `PathMetadataEntry` to `FileEntry.PathMetadataEntries[pathString]` and adds the `FileEntry` to `PathMetadataEntry.AssociatedFileEntries`
   - Processes all paths for each file to handle files with multiple path entries
5. Establishes parent path associations:
   - For each `PathMetadataEntry`, derives the parent path using `filepath.Dir` on `PathMetadataEntry.Path.Path` (after normalizing trailing slashes)
   - Looks up the parent `PathMetadataEntry` in the map
   - If found, sets `PathMetadataEntry.ParentPath` to the parent `PathMetadataEntry` pointer
   - **This operation never fails** - missing parents are valid for root paths or when path metadata is optional
   - `ParentPath` is a runtime-only convenience field (not persisted) that remains `nil` when parents are not found
   - This builds the path hierarchy needed for inheritance resolution

This function ensures that the entire package has consistent file-path metadata associations, enabling proper tag inheritance and filesystem property management across all files.

### 8.3 Special Metadata File Management

Special metadata files must be saved with specific flags and file types to ensure proper recognition and processing.

#### 8.3.1 Special File Type Requirements

- Must use special file types (see [File Types System - Special Files](file_type_system.md#339-special-file-types-65000-65535))
- Must have reserved file names (e.g., `__NVPK_PATH_65001__.nvpkpath`)
- Must support either uncompressed data or LZ4-compressed data with automatic decompression on read
- Must have proper package header flags set

#### 8.3.2 Special File Types

- **Type 65000**: Package metadata (`__NVPK_PKG_65000__.yaml`)
- **Type 65001**: Path metadata (`__NVPK_PATH_65001__.nvpkpath`)
- **Type 65002**: Symbolic link metadata (`__NVPK_SYMLINK_65002__.nvpksym`)
- **Type 65003-65535**: Reserved for future use

#### 8.3.3 PackageHeader Flags

- **Bit 6**: Has special metadata files (set to 1 when special files exist)
- **Bit 5**: Has per-file tags (set to 1 if path metadata provides inheritance)

#### 8.3.4 FileEntry Requirements

- `Type` field set to appropriate special file type (65001 for path metadata)
- `CompressionType` set to 0 (no compression) or 2 (LZ4 compression)
- `EncryptionType` set to 0x00 (no encryption) - special files should not be encrypted
- `Tags` should include `file_type=special_metadata` and `metadata_type=path`

**Automatic Decompression**: When loading path metadata files (via `LoadPathMetadataFile` or when accessing the FileEntry data), the system automatically decompresses the file data if `CompressionType` is set to LZ4 (2).
This works similarly to how package-level compression is automatically handled - the decompression occurs transparently when the file data is accessed.
See [FileEntry API - Data Management](api_file_mgmt_file_entry.md#4-data-management) for details on automatic decompression behavior.

#### 8.3.5 Implementation Details

This section describes implementation details for path metadata operations.

##### 8.3.5.1 Package.SavePathMetadataFile Method

```go
// SavePathMetadataFile creates and saves the path metadata file
func (p *Package) SavePathMetadataFile(ctx context.Context) error
```

##### 8.3.5.2 Package.UpdateSpecialMetadataFlags Method

```go
// UpdateSpecialMetadataFlags updates package header flags based on special files
func (p *Package) UpdateSpecialMetadataFlags() error
```

**SavePathMetadataFile** creates and saves the path metadata file as a special metadata file in the package.

The function performs the following operations:

1. Retrieves the current path metadata entries using `GetPathMetadata`.
2. Marshals the path entries to YAML format, wrapping them in a map with the key `"paths"`.
3. Creates a special metadata FileEntry using `CreateSpecialMetadataFile` with:
   - File type `65001` (path metadata)
   - File name `"__NVPK_PATH_65001__.nvpkpath"`
   - The marshaled YAML data as content
   - `CompressionType` set to `0` (none) or `2` (LZ4) depending on the requested behavior
4. Sets the required tags on the FileEntry:
   - `file_type` = `"special_metadata"` (string tag)
   - `metadata_type` = `"path"` (string tag)
   - `format` = `"yaml"` (string tag)
   - `version` = `1` (integer tag)
5. Updates the package header flags by calling `UpdateSpecialMetadataFlags`.

**LoadPathMetadataFile** loads and parses the path metadata file from the package.

The function performs the following operations:

1. Locates the path metadata special file (type 65001) in the package.
2. Loads the file data using `FileEntry.LoadData()`, which automatically decompresses the data if `CompressionType` is set to LZ4 (2).
3. Parses the YAML content to extract path metadata entries.
4. Updates the package's internal path metadata state with the loaded entries.

The automatic decompression behavior is similar to how package-level compression is handled - when `FileEntry.LoadData()` is called, it detects the `CompressionType` field and automatically decompresses the data before returning it.
This ensures that path metadata files compressed with LZ4 are transparently decompressed when loaded, without requiring explicit decompression calls.
See [Package Compression API](api_package_compression.md) for details on compression and decompression behavior.

**UpdateSpecialMetadataFlags** updates the package header flags and PackageInfo to reflect the current state of special metadata files, per-file tags, and extended attributes.

##### 8.3.5.3 Data Flow Architecture

- **During OpenPackage/Read**: Header flags are read from disk and used to populate PackageInfo (header is source of truth)
- **During In-Memory Operations**: PackageInfo is the source of truth for package state
- **During SavePackage/Write**: PackageInfo is used to update header flags before writing to disk (PackageInfo is source of truth)

**Current Implementation Note**: This function currently updates both header flags and PackageInfo simultaneously during in-memory operations. When write operations are fully implemented, the write path will use PackageInfo as the source of truth to update header flags before serialization.

The function:

1. Checks for the presence of special metadata files.
2. Checks for per-file tags (path metadata properties).
3. Checks for extended attributes (filesystem metadata in path metadata entries).
4. Updates the package header flags accordingly:
   - Sets Bit 6 (FlagHasSpecialMetadata) if special metadata files exist
   - Sets Bit 5 (FlagHasPerFileTags) if per-file tags are present
   - Sets Bit 3 (FlagHasExtendedAttrs) if extended attributes are present
5. Updates the corresponding PackageInfo fields to match header flags.

This ensures the package header and PackageInfo accurately reflect the package structure and capabilities during in-memory operations.

### 8.4 Path Association System

The path association system links FileEntry objects to their corresponding PathMetadataEntry metadata, enabling tag inheritance and filesystem property management.

#### 8.4.1 FileEntry Path Properties

**Note**: Since inheritance is now handled only on `PathMetadataEntry` (not `FileEntry`), `FileEntry` no longer maintains `ParentPath` or `InheritedTags` properties.
Inheritance is resolved per-path by accessing the associated `PathMetadataEntry` and walking up its `ParentPath` chain.

#### 8.4.2 PathMetadataEntry Filesystem Properties

- `Mode` - Unix/Linux path permissions (octal)
- `UID`/`GID` - User and Group IDs
- `ACL` - Access Control List entries
- `WindowsAttrs` - Windows path attributes
- `ExtendedAttrs` - Extended attributes map
- `Flags` - Filesystem-specific flags

#### 8.4.3 Association Management

File-path associations are managed at the struct level, linking `FileEntry` instances to `PathMetadataEntry` instances directly.
Associations are established by matching paths between `FileEntry.Paths` and `PathMetadataEntry.Path`.

**FileEntry Association Methods**: See [FileEntry API - Path Management](api_file_mgmt_file_entry.md#5-path-management) for `FileEntry.AssociateWithPathMetadata()` and `FileEntry.GetPathMetadataForPath()`.

**PathMetadataEntry Association Methods**: See [Path Metadata Structures - PathMetadataEntry Methods](#818-pathmetadataentry-methods) for `PathMetadataEntry.AssociateWithFileEntry()` and `PathMetadataEntry.GetAssociatedFileEntries()`.

**Package-Level Association Management**: See [PathMetadata Management Methods](#82-pathmetadata-management-methods) for `Package.UpdateFilePathAssociations()`.

**Note**: Tag management is performed directly on FileEntry or PathMetadataEntry instances using AddTag(), SetTag(), GetTag(), RemoveTag(), etc.
To get inherited or effective tags for a path, use PathMetadataEntry.GetInheritedTags() or PathMetadataEntry.GetEffectiveTags() after retrieving the PathMetadataEntry for the specific path.
GetEffectiveTags() includes tags from associated FileEntry instances, treating them as if they were directly applied to the PathMetadataEntry.

### 8.5 File-Path Association

This section describes file-path association operations.

#### 8.5.1 File-Path Association Query Methods

This section describes file-path association query methods.

##### 8.5.1.1 Package.GetFilePathAssociation Method

```go
// File-path association query methods (Package-level)
// These methods work with path strings to find and return associated structs
func (p *Package) GetFilePathAssociation(filePath string) (*FilePathAssociation, error)
```

##### 8.5.1.2 Package.GetFilesInPath Method

```go
// GetFilesInPath returns all file entries within the specified path.
func (p *Package) GetFilesInPath(path string) ([]*FileEntry, error)
```

##### 8.5.1.3 Package.GetPathFiles Method

```go
// GetPathFiles returns all file entries associated with the specified path.
func (p *Package) GetPathFiles(path string) ([]*FileEntry, error)
```

#### 8.5.2 Path Hierarchy Analysis Methods

This section describes path hierarchy analysis methods.

##### 8.5.2.1 Package.GetPathTree Method

```go
// Path hierarchy analysis
func (p *Package) GetPathTree() (*PathTree, error)
```

##### 8.5.2.2 Package.GetPathStats Method

```go
// GetPathStats returns statistics for all paths in the package.
func (p *Package) GetPathStats() (map[string]PathStats, error)
```

#### 8.5.3 Symbolic Link Management Methods

This section describes symbolic link management methods.

##### 8.5.3.1 Package.AddSymlink Method

```go
// AddSymlink adds a symbolic link to the package
// Parameters:
//   - symlink: SymlinkEntry to add
// Returns:
//   - Error if validation fails or symlink cannot be added
// Validation:
//   - Calls ValidateSymlinkPaths() to ensure paths are valid and within package root
//   - Verifies target exists as FileEntry or PathMetadataEntry directory
//   - Returns ErrTypeValidation, ErrTypeSecurity, or ErrTypeNotFound on validation failure
// Returns *PackageError on failure
func (p *Package) AddSymlink(symlink SymlinkEntry) error
```

##### 8.5.3.2 Package.RemoveSymlink Method

```go
// Returns *PackageError on failure
func (p *Package) RemoveSymlink(sourcePath string) error
```

##### 8.5.3.3 Package.GetSymlink Method

```go
// Returns *PackageError on failure
func (p *Package) GetSymlink(sourcePath string) (*SymlinkEntry, error)
```

##### 8.5.3.4 Package.ListSymlinks Method

```go
// Returns *PackageError on failure
func (p *Package) ListSymlinks() ([]SymlinkEntry, error)
```

##### 8.5.3.5 Package.UpdateSymlink Method

```go
// Returns *PackageError on failure
func (p *Package) UpdateSymlink(sourcePath string, symlink SymlinkEntry) error
```

##### 8.5.3.6 Package.SaveSymlinkMetadataFile Method

```go
// Returns *PackageError on failure
func (p *Package) SaveSymlinkMetadataFile(ctx context.Context, symlink SymlinkEntry) error
```

##### 8.5.3.7 Package.LoadSymlinkMetadataFile Method

```go
// Returns *PackageError on failure
func (p *Package) LoadSymlinkMetadataFile(ctx context.Context, fileEntry *FileEntry) (*SymlinkEntry, error)
```

#### 8.5.4 Symlink Validation Methods

This section describes symlink validation methods.

##### 8.5.4.1 Package.ValidateSymlinkPaths Method

```go
// ValidateSymlinkPaths validates symlink source and target paths
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - sourcePath: The symlink source path (where the symlink will be created)
//   - targetPath: The symlink target path (where the symlink points to)
// Returns:
//   - Error if validation fails
// Validation performed:
//   - Both paths are package-relative (start with "/")
//   - Both paths are within package root (no ".." escapes)
//   - Target path exists as FileEntry or PathMetadataEntry directory
//   - Returns ErrTypeValidation for invalid paths
//   - Returns ErrTypeSecurity for paths escaping package root
//   - Returns ErrTypeNotFound if target does not exist
func (p *Package) ValidateSymlinkPaths(ctx context.Context, sourcePath, targetPath string) error
```

##### 8.5.4.2 Package.TargetExists Method

```go
// TargetExists checks if a path exists as FileEntry or directory PathMetadataEntry
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - path: The path to check (package-relative with leading "/")
// Returns:
//   - true if path exists as FileEntry or directory PathMetadataEntry
//   - false otherwise
func (p *Package) TargetExists(ctx context.Context, path string) bool
```

##### 8.5.4.3 Package.ValidatePathWithinPackageRoot Method

```go
// ValidatePathWithinPackageRoot validates that a path is within package root
// Parameters:
//   - path: The path to validate (package-relative)
// Returns:
//   - Normalized path if valid
//   - Error if path escapes package root or is invalid
//   - Returns ErrTypeValidation for invalid format
//   - Returns ErrTypeSecurity for paths escaping package root
func (p *Package) ValidatePathWithinPackageRoot(path string) (string, error)
```

#### 8.5.5 PathTree Structure

```go
// PathTree represents the complete path hierarchy
type PathTree struct {
    Root       *PathNode `json:"root"`
    TotalDirs  int       `json:"total_dirs"`
    TotalFiles int       `json:"total_files"`
}
```

#### 8.5.6 PathNode Structure

```go
// PathNode represents a node in the path tree
type PathNode struct {
    Path         string       `json:"path"`
    PathMetadata *PathInfo    `json:"path_metadata,omitempty"`
    Files        []*FileEntry `json:"files"`
    Children     []*PathNode  `json:"children"`
}
```

#### 8.5.7 PathStats Structure

```go
// PathStats provides statistics for a path
type PathStats struct {
    FileCount      int    `json:"file_count"`
    TotalSize      int64  `json:"total_size"`
    CompressedSize int64  `json:"compressed_size"`
    LastModified   string `json:"last_modified"`
}
```

#### 8.5.8 SymlinkEntry Structure

```go
// SymlinkEntry represents a symbolic link with metadata
type SymlinkEntry struct {
    SourcePath    string    `yaml:"source_path"`    // Original symlink path
    TargetPath    string    `yaml:"target_path"`    // Target path (resolved)
    Properties    []Tag     `yaml:"properties"`     // Symlink-specific tags
    Metadata      SymlinkMetadata `yaml:"metadata"` // Symlink metadata
    FileSystem    SymlinkFileSystem `yaml:"filesystem"` // Filesystem properties
}
```

#### 8.5.9 SymlinkMetadata Structure

```go
// SymlinkMetadata contains symlink creation and modification information
type SymlinkMetadata struct {
    Created   time.Time `yaml:"created"`   // When symlink was created
    Modified  time.Time `yaml:"modified"`  // When symlink was last modified
    Description string  `yaml:"description,omitempty"` // Optional description
}
```

#### 8.5.10 SymlinkFileSystem Structure

```go
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

#### 8.5.11 Symlink Validation Methods (Details)

The following methods provide validation for symlink operations:

##### 8.5.11.1 ValidateSymlinkPaths Method (Validation Details)

See [8.5.4.1 Package.ValidateSymlinkPaths Method](#8541-packagevalidatesymlinkpaths-method) for the method signature.

Validates that symlink source and target paths are valid and within package boundaries.

##### 8.5.11.2 Validation Rules

- Both paths must be package-relative (start with "/")
- Both paths must be within package root (no ".." escapes)
- Target path must exist as FileEntry or PathMetadataEntry directory
- Returns `ErrTypeValidation` for invalid paths
- Returns `ErrTypeSecurity` for paths escaping package root
- Returns `ErrTypeNotFound` if target does not exist

##### 8.5.11.3 TargetExists Method (Validation Details)

See [8.5.4.2 Package.TargetExists Method](#8542-packagetargetexists-method) for the method signature.

Checks if a path exists as FileEntry or directory PathMetadataEntry.

##### 8.5.11.4 Use Cases

- Verify symlink targets before creating symlinks
- Validate paths during conversion operations
- Check path existence for validation workflows

##### 8.5.11.5 ValidatePathWithinPackageRoot Method (Validation Details)

See [8.5.4.3 Package.ValidatePathWithinPackageRoot Method](#8543-packagevalidatepathwithinpackageroot-method) for the method signature.

Validates that a path is within package root and does not escape boundaries.

##### 8.5.11.6 ValidatePathWithinPackageRoot Validation Rules

- Path must be package-relative format
- No ".." components that escape package root
- Returns normalized path if valid
- Returns `ErrTypeValidation` for invalid format
- Returns `ErrTypeSecurity` for paths escaping package root

#### 8.5.12 Symlink Creation from Duplicate Paths

Symlinks can be created automatically when converting duplicate path entries on a FileEntry to symlinks.

##### 8.5.12.1 Conversion Process

1. Select primary path (explicit, custom selector, or lexicographic)
2. Create SymlinkEntry for each non-primary path pointing to primary
3. Create PathMetadataEntry with `Type: PathMetadataTypeFileSymlink`
4. Update FileEntry to have single primary path

##### 8.5.12.2 Package Root Boundary Enforcement

- All symlink paths (source and target) MUST be within package root
- No external filesystem references allowed
- Validation performed via `ValidateSymlinkPaths()` before creation

##### 8.5.12.3 Target Existence Validation

- Symlink targets MUST exist before symlink creation
- Verified via `TargetExists()` check
- Returns `ErrTypeNotFound` if target does not exist

See [File Update API - ConvertPathsToSymlinks](api_file_mgmt_updates.md#1711-packageconvertpathstosymlinks-method) for complete conversion API documentation.

##### 8.5.12.4 Example: Converting Duplicate Paths to Symlinks

```go
// Find FileEntry with multiple paths
entry, _ := pkg.GetFileByPath(ctx, "/app/bin/main")
// Assume entry has paths: ["/app/bin/main", "/usr/local/bin/main", "/opt/main"]

// Convert to symlinks with explicit primary path
options := &SymlinkConvertOptions{
    PrimaryPath: Option.Some("/app/bin/main"),
    PreservePathMetadata: Option.Some(true),
}
updatedEntry, symlinks, err := pkg.ConvertPathsToSymlinks(ctx, entry, options)
// Result:
// - updatedEntry has single path: "/app/bin/main"
// - symlinks contains 2 SymlinkEntry objects pointing to "/app/bin/main"
```
