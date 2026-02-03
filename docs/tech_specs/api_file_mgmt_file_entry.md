# NovusPack Technical Specifications - FileEntry API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. FileEntry Structure](#1-fileentry-structure)
  - [1.1 `FileEntry` Structure Definition](#11-fileentry-structure-definition)
  - [1.2 Related Type Definitions](#12-related-type-definitions)
  - [1.3 Runtime-Only Fields](#13-runtime-only-fields)
- [2. FileEntry Creation](#2-fileentry-creation)
  - [2.1 NewFileEntry Function](#21-newfileentry-function)
    - [2.1.1 `NewFileEntry` Function](#211-newfileentry-function-signature)
    - [2.1.2 Example Usage](#212-example-usage)
- [3. Tag Management](#3-tag-management)
  - [3.1 Tag Placement Guidance](#31-tag-placement-guidance)
    - [3.1.1 Example: File and Path Tag Placement](#311-example-file-and-path-tag-placement)
    - [3.1.2 Tag Management Function Signatures](#312-tag-management-function-signatures)
  - [3.2 Tag Type](#32-tag-type)
  - [3.3 Tag Operations Usage](#33-tag-operations-usage)
    - [3.3.1 Getting All Tags](#331-getting-all-tags)
    - [3.3.2 Getting Tags by Type](#332-getting-tags-by-type)
    - [3.3.3 Adding Multiple Tags](#333-adding-multiple-tags)
    - [3.3.4 Updating Multiple Tags](#334-updating-multiple-tags)
    - [3.3.5 Getting Individual Tags](#335-getting-individual-tags)
    - [3.3.6 Adding Individual Tags](#336-adding-individual-tags)
    - [3.3.7 Updating Individual Tags](#337-updating-individual-tags)
- [4. Data Management](#4-data-management)
  - [4.1 Basic Data Operations](#41-basic-data-operations)
    - [4.1.1 FileEntry LoadData Method](#411-fileentry-loaddata-method)
    - [4.1.2 `FileEntry.UnloadData` Method](#412-fileentryunloaddata-method)
    - [4.1.3 `FileEntry.GetData` Method](#413-fileentrygetdata-method)
    - [4.1.4 `FileEntry.SetData` Method](#414-fileentrysetdata-method)
  - [4.2 Temporary File Operations](#42-temporary-file-operations)
    - [4.2.1 `FileEntry.CreateTempFile` Method](#421-fileentrycreatetempfile-method)
    - [4.2.2 `FileEntry.StreamToTempFile` Method](#422-fileentrystreamtotempfile-method)
    - [4.2.3 `FileEntry.WriteToTempFile` Method](#423-fileentrywritetotempfile-method)
    - [4.2.4 `FileEntry.ReadFromTempFile` Method](#424-fileentryreadfromtempfile-method)
    - [4.2.5 `FileEntry.CleanupTempFile` Method](#425-fileentrycleanuptempfile-method)
  - [4.3 ProcessingState Management](#43-processingstate-management)
    - [4.3.1 `FileEntry.GetProcessingState` Method](#431-fileentrygetprocessingstate-method)
    - [4.3.2 `FileEntry.SetProcessingState` Method](#432-fileentrysetprocessingstate-method)
  - [4.4 Source Tracking (CurrentSource/OriginalSource)](#44-source-tracking-currentsourceoriginalsource)
    - [4.4.1 `FileEntry.SetCurrentSource` Method](#441-fileentrysetcurrentsource-method)
    - [4.4.2 `FileEntry.GetCurrentSource` Method](#442-fileentrygetcurrentsource-method)
    - [4.4.3 `FileEntry.SetOriginalSource` Method](#443-fileentrysetoriginalsource-method)
    - [4.4.4 `FileEntry.GetOriginalSource` Method](#444-fileentrygetoriginalsource-method)
    - [4.4.5 `FileEntry.IsCurrentSourceTempFile` Method](#445-fileentryiscurrentsourcetempfile-method)
    - [4.4.6 `FileEntry.HasOriginalSource` Method](#446-fileentryhasoriginalsource-method)
    - [4.4.7 FileEntry SetOriginalSourceFromPackage Method](#447-fileentrysetoriginalsourcefrompackage-method)
    - [4.4.8 FileEntry CopyCurrentToOriginal Method](#448-fileentrycopycurrenttooriginal-method)
  - [4.5 Multi-Stage Transformation Pipeline](#45-multi-stage-transformation-pipeline)
    - [4.5.1 FileEntry InitializeTransformPipeline Method](#451-fileentryinitializetransformpipeline-method)
    - [4.5.2 FileEntry GetTransformPipeline Method](#452-fileentrygettransformpipeline-method)
    - [4.5.3 FileEntry ExecuteTransformStage Method](#453-fileentryexecutetransformstage-method)
    - [4.5.4 FileEntry ResumeTransformation Method](#454-fileentryresumetransformation-method)
    - [4.5.5 FileEntry CleanupTransformPipeline Method](#455-fileentrycleanuptransformpipeline-method)
    - [4.5.6 FileEntry ValidateSources Method](#456-fileentryvalidatesources-method)
- [5. Path Management](#5-path-management)
  - [5.1 `FileEntry.HasSymlinks` Method](#51-fileentryhassymlinks-method)
  - [5.2 FileEntry GetSymlinkPaths Method](#52-fileentrygetsymlinkpaths-method)
  - [5.3 FileEntry GetPrimaryPath Method](#53-fileentrygetprimarypath-method)
  - [5.4 `FileEntry.ResolveAllSymlinks` Method](#54-fileentryresolveallsymlinks-method)
  - [5.5 FileEntry AssociateWithPathMetadata Method](#55-fileentryassociatewithpathmetadata-method)
  - [5.6 FileEntry GetPathMetadataForPath Method](#56-fileentrygetpathmetadataforpath-method)
- [6. Marshaling](#6-marshaling)
  - [6.1 Marshal Methods](#61-marshal-methods)
    - [6.1.1 FileEntry MarshalMeta Method](#611-fileentrymarshalmeta-method)
    - [6.1.2 FileEntry MarshalData Method](#612-fileentrymarshaldata-method)
    - [6.1.3 FileEntry Marshal Method](#613-fileentrymarshal-method)
  - [6.2 WriteTo Methods](#62-writeto-methods)
    - [6.2.1 FileEntry WriteMetaTo Method](#621-fileentrywritemetato-method)
    - [6.2.2 FileEntry WriteDataTo Method](#622-fileentrywritedatato-method)
    - [6.2.3 FileEntry WriteTo Method](#623-fileentrywriteto-method)
  - [6.3 Unmarshal Function](#63-unmarshalfileentry-function)
  - [6.4 FileEntry Marshaling Purpose](#64-fileentry-marshaling-purpose)
  - [6.5 FileEntry Marshaling Usage Notes](#65-fileentry-marshaling-usage-notes)
- [7. FileEntry Properties](#7-fileentry-properties)
  - [7.1 IsCompressed](#71-fileentryiscompressed-method)
  - [7.2 HasEncryptionKey](#72-fileentryhasencryptionkey-method)
  - [7.3 GetEncryptionType](#73-fileentrygetencryptiontype-method)
  - [7.4 IsEncrypted](#74-fileentryisencrypted-method)
  - [7.5 FileEntry Properties Purpose](#75-fileentry-properties-purpose)
  - [7.6 IsCompressed Returns](#76-iscompressed-returns)
  - [7.7 HasEncryptionKey Returns](#77-hasencryptionkey-returns)
  - [7.8 GetEncryptionType Returns](#78-getencryptiontype-returns)
  - [7.9 IsEncrypted Returns](#79-isencrypted-returns)
- [8. FileEntry Compression](#8-fileentry-compression)
  - [8.1 FileEntry Compress Method](#81-fileentrycompress-method)
  - [8.2 FileEntry Decompress Method](#82-fileentrydecompress-method)
  - [8.3 FileEntry GetCompressionInfo Method](#83-fileentrygetcompressioninfo-method)
  - [8.4 FileEntry Compression Purpose](#84-fileentry-compression-purpose)
  - [8.5 Compress Parameters](#85-compress-parameters)
  - [8.6 Decompress Parameters](#86-decompress-parameters)
  - [8.7 GetCompressionInfo Parameters](#87-getcompressioninfo-parameters)
  - [8.8 FileEntry Compression Error Conditions](#88-fileentry-compression-error-conditions)
  - [8.9 FileEntry Compression Usage Notes](#89-fileentry-compression-usage-notes)
- [9. FileEntry Encryption](#9-fileentry-encryption)
  - [9.1 SetEncryptionKey Method](#91-fileentrysetencryptionkey-method)
  - [9.2 FileEntry Encrypt Method](#92-fileentryencrypt-method)
  - [9.3 FileEntry Decrypt Method](#93-fileentrydecrypt-method)
  - [9.4 FileEntry UnsetEncryptionKey Method](#94-fileentryunsetencryptionkey-method)
  - [9.5 FileEntry Encryption Purpose](#95-fileentry-encryption-purpose)
  - [9.6 SetEncryptionKey Parameters](#96-setencryptionkey-parameters)
  - [9.7 Encrypt/Decrypt Parameters](#97-encryptdecrypt-parameters)
  - [9.8 FileEntry Encryption Error Conditions](#98-fileentry-encryption-error-conditions)
  - [9.9 FileEntry Encryption Usage Notes](#99-fileentry-encryption-usage-notes)
- [10. FileEntry Data Management](#10-fileentry-data-management)
  - [10.1 LoadData Method](#101-fileentryloaddata-method)
  - [10.2 FileEntry ProcessData Method](#102-fileentryprocessdata-method)
  - [10.3 FileEntry Data Management Purpose](#103-fileentry-data-management-purpose)
  - [10.4 LoadData Behavior](#104-loaddata-behavior)
  - [10.5 ProcessData Behavior](#105-processdata-behavior)
  - [10.6 FileEntry Data Management Error Conditions](#106-fileentry-data-management-error-conditions)
  - [10.7 FileEntry Data Management Usage Notes](#107-fileentry-data-management-usage-notes)
  - [11. HashEntry Type](#11-hashentry-struct)
- [12. HashType Type](#12-hashtype-type)
- [13. HashPurpose Type](#13-hashpurpose-type)
- [14. TagValueType Type](#14-tagvaluetype-type)
- [15. `ProcessingState` Type](#15-processingstate-type)
- [16. `FileSource` Structure](#16-filesource-structure)
- [17. `OptionalData` Structure](#17-optionaldata-structure)
- [18. OptionalDataType Type](#18-optionaldatatype-type)
- [19. Tag Generic Type](#19-tag-generic-type)
  - [19.1 `Tag` Type Definition](#191-tag-struct)
  - [19.2 `NewTag` Function](#192-newtag-function)
  - [19.3 Tag GetValue Method](#193-tagtgetvalue-method)
  - [19.4 `Tag.SetValue` Method](#194-tagtsetvalue-method)

---

## 0. Overview

This document specifies the FileEntry data model and FileEntry-scoped methods.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [Core Package Interface](api_core.md)
- [Package Metadata API](api_metadata.md)
- [Generic Types and Patterns](api_generics.md)

## 1. FileEntry Structure

This section defines the FileEntry structure and its components.

### 1.1 FileEntry Structure Definition

FileEntry represents a single file content with its metadata. A FileEntry can have multiple paths (aliases) stored in the `Paths` array, allowing the same file content to be accessible via different paths within the package.

**Multiple Paths Behavior**: When a FileEntry has multiple paths and one path is removed:

- Only the specified path is removed from the `Paths` array
- The `PathCount` field is decremented
- If the removed path is the last path (PathCount becomes 0), the entire FileEntry is removed from the package
- This ensures file content is preserved as long as at least one path references it

```go
// FileEntry represents a FileEntry in the package with complete metadata.
type FileEntry struct {
    // Static fields (64 bytes total)
    FileID             uint64    // Unique file identifier (8 bytes)
    OriginalSize       uint64    // Original file size before processing (8 bytes)
    StoredSize         uint64    // Final file size after compression/encryption (8 bytes)
    RawChecksum        uint32    // CRC32 of raw file content (4 bytes)
    StoredChecksum     uint32    // CRC32 of processed file content (4 bytes)
    FileVersion        uint32    // File data version (4 bytes)
    MetadataVersion    uint32    // File metadata version (4 bytes)
    PathCount          uint16    // Number of paths (2 bytes)
    Type               FileType  // File type identifier (2 bytes)
    CompressionType    uint8     // Compression algorithm identifier (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    CompressionLevel   uint8     // Compression level (1 byte)
    EncryptionType     uint8     // Encryption algorithm identifier (1 byte)
    HashCount          uint8     // Number of hash entries (1 byte)
    HashDataOffset     uint32    // Offset to hash data (4 bytes)
    HashDataLen        uint16    // Length of hash data (2 bytes)
    OptionalDataLen    uint16    // Length of optional data (2 bytes)
    OptionalDataOffset uint32    // Offset to optional data (4 bytes)
    Reserved           uint32    // Reserved for future use (4 bytes)

    // Variable-length data (populated on demand)
    Paths              []generics.PathEntry    // File paths (minimal PathEntry from generics package)
    Hashes             []HashEntry    // Content hashes
    OptionalData       OptionalData   // Structured optional data

    // Convenience properties (computed from OptionalData)
    Tags               []*Tag[any]    // Direct access to OptionalData.Tags (typed tags)

    // Data management (runtime only, not stored in file)
    EntryOffset        uint64               // Absolute offset to the FileEntry metadata start in the package file
    Data               []byte               // File content in memory (only for small files being processed)
    IsDataLoaded       bool                 // Whether data is currently loaded in memory
    ProcessingState    ProcessingState      // Current processing state of the file (data-state model)

    // Source tracking (unified approach replacing SourceFile/SourceOffset/SourceSize/TempFilePath/IsTempFile)
    CurrentSource      *FileSource          // Current data source (may be original, temp, or final stage)
    OriginalSource     *FileSource          // Original source before transformations (nil for new files)
    TransformPipeline  *TransformPipeline   // Transformation stages (nil for simple operations)

    // PathMetadataEntry associations (runtime only, not stored in file)
    // Maps each path in FileEntry.Paths to its corresponding PathMetadataEntry.
    // Inheritance is now handled only on PathMetadataEntry, not FileEntry.
    // To access inherited tags, use the associated PathMetadataEntry's GetInheritedTags() method.
    PathMetadataEntries map[string]*PathMetadataEntry // Path -> PathMetadataEntry mapping
}
```

### 1.2 Related Type Definitions

This document defines `FileEntry` as the primary data model.
Supporting types are defined under their own H2 sections for easier navigation.

- Hashes: [HashEntry Type](#11-hashentry-struct), [HashType Type](#12-hashtype-type), [HashPurpose Type](#13-hashpurpose-type)
- Tags: [TagValueType Type](#14-tagvaluetype-type), [Tag Generic Type](#19-tag-generic-type)
- Processing and pipelines:
  - [ProcessingState Type](#15-processingstate-type)
  - [FileSource Structure](#16-filesource-structure)
  - [File Transformation Pipelines: TransformPipeline Structure](api_file_mgmt_transform_pipelines.md#22-transformpipeline-structure)
  - [File Transformation Pipelines: TransformStage Structure](api_file_mgmt_transform_pipelines.md#23-transformstage-structure)
  - [File Transformation Pipelines: TransformType Type](api_file_mgmt_transform_pipelines.md#24-transformtype-type)
- Optional data: [OptionalData Structure](#17-optionaldata-structure), [OptionalDataType Type](#18-optionaldatatype-type)

### 1.3 Runtime-Only Fields

The following `FileEntry` fields are runtime-only and MUST NOT be serialized into the package file format.

- `EntryOffset`: Absolute offset, in bytes, to the start of the FileEntry metadata in the package file.
- The file data payload start offset is derived as `EntryOffset + TotalSize()`.
- `EntryOffset` MUST be populated during package open operations after the file index is loaded.
- If the offset is unknown (for example, staged in-memory entries that are not yet written), it MUST be set to `0`.
- `CurrentSource`, `OriginalSource`, `TransformPipeline`: Source tracking and multi-stage transformation state (see [Source Tracking (CurrentSource/OriginalSource)](#44-source-tracking-currentsourceoriginalsource)).

Migration from previous source fields:

| Old Field         | New Equivalent                | Notes                         |
| ----------------- | ----------------------------- | ----------------------------- |
| `fe.SourceFile`   | `fe.CurrentSource.File`       | Nil check CurrentSource first |
| `fe.SourceOffset` | `fe.CurrentSource.Offset`     | Nil check CurrentSource first |
| `fe.SourceSize`   | `fe.CurrentSource.Size`       | Nil check CurrentSource first |
| `fe.TempFilePath` | `fe.CurrentSource.FilePath`   | Only if IsTempFile is true    |
| `fe.IsTempFile`   | `fe.CurrentSource.IsTempFile` | Nil check CurrentSource first |

**Note**: The previous SetSourceFile, GetSourceFile, SetTempPath, and GetTempPath methods have been removed.
Use SetCurrentSource and GetCurrentSource helper methods instead (see [Source Tracking (CurrentSource/OriginalSource)](#44-source-tracking-currentsourceoriginalsource)).

## 2. FileEntry Creation

This section describes how FileEntry instances are created.

### 2.1 NewFileEntry Function

Creates a new FileEntry with proper tag synchronization.

#### 2.1.1 NewFileEntry Function Signature

```go
// NewFileEntry creates a new FileEntry with proper tag synchronization.
func NewFileEntry() *FileEntry
```

Returns a new FileEntry instance with all fields initialized and tag system ready for use.

This is the primary way to create a new FileEntry instance. The function ensures that:

- All fields are properly initialized to their zero values
- The tag system is ready for use with proper synchronization
- The FileEntry is in a valid state for further operations

#### 2.1.2 Example Usage

```go
// Create a new FileEntry.
fe := NewFileEntry()

// Now you can work with the FileEntry.
// Add tags, set paths, load data, etc.
```

**Note**: For unmarshaling FileEntry instances from binary data, see [UnmarshalFileEntry](#63-unmarshalfileentry-function).

## 3. Tag Management

All tag operations use typed tags for type safety.
The tag system only supports typed tags; untyped tags are not supported.

**Tag Management Architecture**: Tags are managed directly on FileEntry and PathMetadataEntry instances.

To manage tags for a file, retrieve the FileEntry (e.g., via `GetFileByPath()`) and call tag methods on that instance.
To manage tags for a path, work with PathMetadataEntry instances directly.

### 3.1 Tag Placement Guidance

Use the following guidance to decide where to set tags:

- **Set tags on `FileEntry`** when the same tags need to apply to **all paths** for a file that has multiple paths.

  - FileEntry tags are shared across all paths in the `FileEntry.Paths` slice.
  - Use this for file-level metadata that should be consistent regardless of which path is used to access the file.
  - FileEntry tags are automatically included in the effective tags for each associated PathMetadataEntry (as if they were directly applied to the PathMetadataEntry).

- **Set tags on `PathMetadataEntry`** when tags should only apply to a **specific path**.
  - PathMetadataEntry tags are path-specific and do not affect other paths for the same file content.
  - Use this when different paths for the same file content need different metadata (e.g., different permissions, different tags).
  - PathMetadataEntry tags can participate in inheritance (see [Path Metadata System - Tag Inheritance](api_metadata.md#8-pathmetadata-system)).
  - When calling `GetEffectiveTags()` on a PathMetadataEntry, it includes: direct PathMetadataEntry tags, inherited path tags, and associated FileEntry tags.

#### 3.1.1 Example: File and Path Tag Placement

```go
// File with multiple paths (all relative, like zip files).
fe := &FileEntry{
    Paths: []generics.PathEntry{
        {Path: "documents/reports/analysis.pdf"},
        {Path: "archive/2024/analysis.pdf"},
    },
}

// Set file-level tag (applies to both paths)
AddFileEntryTag(fe, "version", "1.0.0", TagValueTypeString)

// Set path-specific tag (only applies to documents/reports/analysis.pdf).
// Note: In practice, you would retrieve PathMetadataEntry from the package.
pathInfo1, _ := package.GetPathInfo(ctx, "documents/reports/analysis.pdf")
if pathInfo1 != nil {
    AddPathMetaTag(&pathInfo1.Entry, "department", "finance", TagValueTypeString)
}

// Different path-specific tag (only applies to archive/2024/analysis.pdf)
pathInfo2, _ := package.GetPathInfo(ctx, "archive/2024/analysis.pdf")
if pathInfo2 != nil {
    AddPathMetaTag(&pathInfo2.Entry, "department", "archives", TagValueTypeString)
}
```

#### 3.1.2 Tag Management Function Signatures

This section defines the function signatures for tag management operations.

##### 3.1.2.1 GetFileEntryTags Function

```go
// GetFileEntryTags returns all tags as typed tags for a FileEntry.
// Returns a slice of Tag pointers, where each tag maintains its type information.
// Returns *PackageError on failure (corruption, I/O).
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTags(fe *FileEntry) ([]*Tag[any], error)
```

##### 3.1.2.2 GetFileEntryTagsByType Function

```go
// GetFileEntryTagsByType returns all tags of a specific type for a FileEntry.
// Returns a slice of Tag pointers with the specified type parameter T.
// Only tags matching the type T and corresponding TagValueType are returned.
// Returns *PackageError on failure (corruption, I/O).
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTagsByType[T any](fe *FileEntry) ([]*Tag[T], error)
```

##### 3.1.2.3 GetFileEntryTag Function

```go
// GetFileEntryTag retrieves a type-safe tag by key from a FileEntry.
// Returns the tag pointer and an error. If the tag is not found, returns (nil, nil).
// If an underlying error occurs (corruption, I/O), returns (nil, error).
// Returns *PackageError on failure.
// If the tag type is unknown, use GetFileEntryTag[any](fe, "key") to retrieve the tag and inspect its Type field.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTag[T any](fe *FileEntry, key string) (*Tag[T], error)
```

##### 3.1.2.4 AddFileEntryTags Function

```go
// AddFileEntryTags adds multiple new tags with type safety to a FileEntry.
// Returns *PackageError if any tag with the same key already exists.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func AddFileEntryTags(fe *FileEntry, tags []*Tag[any]) error
```

##### 3.1.2.5 SetFileEntryTags Function

```go
// SetFileEntryTags updates existing tags from a slice of typed tags for a FileEntry.
// Returns *PackageError if any tag key does not already exist.
// Only modifies tags that already exist; does not create new tags.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SetFileEntryTags(fe *FileEntry, tags []*Tag[any]) error
```

##### 3.1.2.6 AddFileEntryTag Function

```go
// AddFileEntryTag adds a new tag with type safety to a FileEntry.
// Returns *PackageError if a tag with the same key already exists
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func AddFileEntryTag[T any](fe *FileEntry, key string, value T, tagType TagValueType) error
```

##### 3.1.2.7 SetFileEntryTag Function

```go
// SetFileEntryTag updates an existing tag with type safety for a FileEntry.
// Returns *PackageError if the tag key does not already exist
// Only modifies existing tags; does not create new tags
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SetFileEntryTag[T any](fe *FileEntry, key string, value T, tagType TagValueType) error
```

##### 3.1.2.8 RemoveFileEntryTag Function

```go
// RemoveFileEntryTag removes a tag by key from a FileEntry.
// Returns *PackageError on failure.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func RemoveFileEntryTag(fe *FileEntry, key string) error
```

##### 3.1.2.9 HasFileEntryTag Function

```go
// HasFileEntryTag checks if a tag with the specified key exists on a FileEntry
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func HasFileEntryTag(fe *FileEntry, key string) bool
```

##### 3.1.2.10 HasFileEntryTags Function

```go
// HasFileEntryTags checks if the FileEntry has any tags.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func HasFileEntryTags(fe *FileEntry) bool
```

##### 3.1.2.11 SyncFileEntryTags Function

```go
// SyncFileEntryTags synchronizes tags with the underlying storage for a FileEntry.
// Returns *PackageError on failure.
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SyncFileEntryTags(fe *FileEntry) error
```

### 3.2 Tag Type

The `Tag[T]` type provides type-safe tag operations.
All tags are stored and accessed as typed tags.

The canonical Tag type definition is in [Tag Generic Type](#19-tag-generic-type).

**Cross-Reference**: For the underlying generic type patterns, see [Option Type](api_generics.md#11-option-type).

### 3.3 Tag Operations Usage

This section describes how to use tag operations in practice.

#### 3.3.1 Getting All Tags

```go
// Get all tags as typed tags.
tags, err := GetFileEntryTags(fe)
if err != nil {
    // Handle underlying error (corruption, I/O).
    fmt.Printf("Error getting tags: %v\n", err)
    return
}
for _, tag := range tags {
    // Each tag is a *Tag[any], access value based on tag.Type.
    switch tag.Type {
    case TagValueTypeString:
        if strTag, ok := tag.(*Tag[string]); ok {
            fmt.Printf("Tag %s: %s\n", tag.Key, strTag.GetValue())
        }
    case TagValueTypeInteger:
        if intTag, ok := tag.(*Tag[int64]); ok {
            fmt.Printf("Tag %s: %d\n", tag.Key, intTag.GetValue())
        }
    }
}
```

#### 3.3.2 Getting Tags by Type

```go
// Get all string tags.
stringTags, err := GetFileEntryTagsByType[string](fe)
if err != nil {
    // Handle underlying error (corruption, I/O).
    fmt.Printf("Error getting string tags: %v\n", err)
    return
}
for _, tag := range stringTags {
    fmt.Printf("String tag %s: %s\n", tag.Key, tag.GetValue())
}

// Get all integer tags.
intTags, err := GetFileEntryTagsByType[int64](fe)
if err != nil {
    fmt.Printf("Error getting integer tags: %v\n", err)
    return
}
for _, tag := range intTags {
    fmt.Printf("Integer tag %s: %d\n", tag.Key, tag.GetValue())
}

// Get all boolean tags.
boolTags, err := GetFileEntryTagsByType[bool](fe)
if err != nil {
    fmt.Printf("Error getting boolean tags: %v\n", err)
    return
}
for _, tag := range boolTags {
    fmt.Printf("Boolean tag %s: %v\n", tag.Key, tag.GetValue())
}
```

#### 3.3.3 Adding Multiple Tags

```go
// Add multiple new tags (fails if any key already exists).
tags := []*Tag[any]{
    NewTag("author", "John Doe", TagValueTypeString),
    NewTag("version", int64(1), TagValueTypeInteger),
    NewTag("published", true, TagValueTypeBoolean),
}
err := AddFileEntryTags(fe, tags)
```

#### 3.3.4 Updating Multiple Tags

```go
// Update multiple existing tags (fails if any key does not exist)
tags := []*Tag[any]{
    NewTag("author", "Jane Doe", TagValueTypeString),
    NewTag("version", int64(2), TagValueTypeInteger),
}
err := SetFileEntryTags(fe, tags)
```

#### 3.3.5 Getting Individual Tags

```go
// Get a tag with type safety when you know the type.
authorTag, err := GetFileEntryTag[string](fe, "author")
if err != nil {
    // Handle underlying error (corruption, I/O).
    fmt.Printf("Error getting tag: %v\n", err)
} else if authorTag != nil {
    // Tag found.
    fmt.Printf("Author: %s\n", authorTag.GetValue())
}
// If both err and authorTag are nil, tag was not found (normal case)

versionTag, err := GetFileEntryTag[int64](fe, "version")
if err != nil {
    fmt.Printf("Error getting tag: %v\n", err)
} else if versionTag != nil {
    fmt.Printf("Version: %d\n", versionTag.GetValue())
}

// Get a tag when you don't know the type.
unknownTag, err := GetFileEntryTag[any](fe, "someKey")
if err != nil {
    fmt.Printf("Error getting tag: %v\n", err)
} else if unknownTag != nil {
    // Inspect the tag's Type field to determine how to handle it
    switch unknownTag.Type {
    case TagValueTypeString:
        if strValue, ok := unknownTag.Value.(string); ok {
            fmt.Printf("String tag: %s\n", strValue)
        }
    case TagValueTypeInteger:
        if intValue, ok := unknownTag.Value.(int64); ok {
            fmt.Printf("Integer tag: %d\n", intValue)
        }
    // ... handle other types
    }
}
```

#### 3.3.6 Adding Individual Tags

```go
// Add a new tag (fails if tag with same key already exists).
err := AddFileEntryTag(fe, "author", "Jane Doe", TagValueTypeString)
err := AddFileEntryTag(fe, "version", int64(2), TagValueTypeInteger)
err := AddFileEntryTag(fe, "published", true, TagValueTypeBoolean)
```

#### 3.3.7 Updating Individual Tags

```go
// Update an existing tag (fails if tag does not exist)
err := SetFileEntryTag(fe, "author", "Jane Doe", TagValueTypeString)
err := SetFileEntryTag(fe, "version", int64(2), TagValueTypeInteger)
err := SetFileEntryTag(fe, "published", false, TagValueTypeBoolean)
```

**Note**: The old untyped `Tag` struct and type-specific methods (`SetStringTag`, `GetStringTag`, etc.) are no longer supported.
All tag operations must use the typed tag system.

## 4. Data Management

This section describes data management operations for FileEntry.

### 4.1 Basic Data Operations

This section describes basic data operations for FileEntry.

#### 4.1.1 FileEntry LoadData Method

See [LoadData Method](#101-fileentryloaddata-method) for the complete method definition and documentation.

#### 4.1.2 FileEntry.UnloadData Method

```go
// UnloadData unloads file data from memory.
func (fe *FileEntry) UnloadData()
```

#### 4.1.3 FileEntry.GetData Method

```go
// GetData returns the in-memory file data.
// Returns *PackageError on failure.
func (fe *FileEntry) GetData() ([]byte, error)
```

#### 4.1.4 FileEntry.SetData Method

```go
// SetData sets the in-memory file data
func (fe *FileEntry) SetData(data []byte)
```

### 4.2 Temporary File Operations

This section describes operations for working with temporary files.

#### 4.2.1 FileEntry.CreateTempFile Method

```go
// CreateTempFile creates a temporary file for staging.
// Returns *PackageError on failure.
func (fe *FileEntry) CreateTempFile(ctx context.Context) error
```

#### 4.2.2 FileEntry.StreamToTempFile Method

```go
// StreamToTempFile streams data to a temporary file
// Returns *PackageError on failure.
func (fe *FileEntry) StreamToTempFile(ctx context.Context) error
```

#### 4.2.3 FileEntry.WriteToTempFile Method

```go
// WriteToTempFile writes data to a temporary file.
// Returns *PackageError on failure.
func (fe *FileEntry) WriteToTempFile(ctx context.Context, data []byte) error
```

#### 4.2.4 FileEntry.ReadFromTempFile Method

```go
// ReadFromTempFile reads data from a temporary file
func (fe *FileEntry) ReadFromTempFile(ctx context.Context, offset, size int64) ([]byte, error)
```

#### 4.2.5 FileEntry.CleanupTempFile Method

```go
// CleanupTempFile removes temporary files.
func (fe *FileEntry) CleanupTempFile(ctx context.Context) error
```

### 4.3 ProcessingState Management

This section describes how processing state is managed for FileEntry.

#### 4.3.1 FileEntry.GetProcessingState Method

```go
// GetProcessingState returns the current processing state
func (fe *FileEntry) GetProcessingState() ProcessingState
```

#### 4.3.2 FileEntry.SetProcessingState Method

```go
// SetProcessingState sets the current processing state.
func (fe *FileEntry) SetProcessingState(state ProcessingState)
```

### 4.4 Source Tracking (CurrentSource/OriginalSource)

CurrentSource is a `*FileSource` that tracks the current location and state of file data.
It replaces the previous separate fields (`SourceFile`, `SourceOffset`, `SourceSize`, `TempFilePath`, `IsTempFile`) with a unified structure.

CurrentSource may point to:

- The opened package file when the file's data is stored in the package (`IsPackage=true`).
- A temporary file created during multi-stage transformations (`IsTempFile=true`).
- An external source file when the FileEntry represents data staged from disk (`IsExternal=true`).

The FileSource structure provides all necessary information:

- File handle (`File`) and path (`FilePath`) for accessing the data.
- Byte range (`Offset`, `Size`) specifying where data starts and how much to read.
- Type flags (`IsPackage`, `IsTempFile`, `IsExternal`) indicating the source category.

OriginalSource is a `*FileSource` that preserves the original data source before any transformations are applied.
This is essential for multi-stage extraction and update flows where the original package file location must be preserved for resume, audit, and cleanup behavior.

TransformPipeline is a `*TransformPipeline` that tracks multiple sequential transformation stages for large files or multi-step operations.
Each stage reads from an InputSource and writes to an OutputSource, and `CurrentSource` is advanced to the latest completed stage output.

#### 4.4.1 FileEntry.SetCurrentSource Method

```go
// SetCurrentSource sets the current data source for the FileEntry
// Returns *PackageError if source is invalid.
func (fe *FileEntry) SetCurrentSource(source *FileSource) error
```

#### 4.4.2 FileEntry.GetCurrentSource Method

```go
// GetCurrentSource returns the current data source
// Returns nil if no current source is set.
func (fe *FileEntry) GetCurrentSource() *FileSource
```

#### 4.4.3 FileEntry.SetOriginalSource Method

```go
// SetOriginalSource sets the original data source before transformations.
func (fe *FileEntry) SetOriginalSource(source *FileSource) error
```

#### 4.4.4 FileEntry.GetOriginalSource Method

```go
// GetOriginalSource returns the original data source.
// Returns nil if no original source is tracked (e.g., new files)
func (fe *FileEntry) GetOriginalSource() *FileSource
```

#### 4.4.5 FileEntry.IsCurrentSourceTempFile Method

```go
// IsCurrentSourceTempFile returns true if current source is a temporary file.
func (fe *FileEntry) IsCurrentSourceTempFile() bool
```

#### 4.4.6 FileEntry.HasOriginalSource Method

```go
// HasOriginalSource returns true if original source is tracked
func (fe *FileEntry) HasOriginalSource() bool
```

#### 4.4.7 FileEntry.SetOriginalSourceFromPackage Method

```go
// SetOriginalSourceFromPackage creates original source pointing to package file.
func (fe *FileEntry) SetOriginalSourceFromPackage(packageFile *os.File, packagePath string)
```

#### 4.4.8 FileEntry.CopyCurrentToOriginal Method

```go
// CopyCurrentToOriginal saves current source as original (for transformations)
func (fe *FileEntry) CopyCurrentToOriginal()
```

### 4.5 Multi-Stage Transformation Pipeline

This section describes the multi-stage transformation pipeline system.

#### 4.5.1 FileEntry.InitializeTransformPipeline Method

```go
// InitializeTransformPipeline creates a new transformation pipeline.
func (fe *FileEntry) InitializeTransformPipeline(stages []TransformStage) error
```

#### 4.5.2 FileEntry.GetTransformPipeline Method

```go
// GetTransformPipeline returns the current transformation pipeline.
// Returns nil if no pipeline is active.
func (fe *FileEntry) GetTransformPipeline() *TransformPipeline
```

#### 4.5.3 FileEntry.ExecuteTransformStage Method

```go
// ExecuteTransformStage executes a specific stage in the pipeline
// Returns *PackageError on failure.
func (fe *FileEntry) ExecuteTransformStage(ctx context.Context, stageIndex int) error
```

#### 4.5.4 FileEntry.ResumeTransformation Method

```go
// ResumeTransformation resumes pipeline from last completed stage.
// Returns *PackageError on failure.
func (fe *FileEntry) ResumeTransformation(ctx context.Context) error
```

#### 4.5.5 FileEntry.CleanupTransformPipeline Method

```go
// CleanupTransformPipeline cleans up all temporary files in pipeline.
// Returns *PackageError on failure.
func (fe *FileEntry) CleanupTransformPipeline() error
```

#### 4.5.6 FileEntry.ValidateSources Method

```go
// ValidateSources validates CurrentSource, OriginalSource, and pipeline consistency
// Returns *PackageError if validation fails.
func (fe *FileEntry) ValidateSources() error
```

**Note**: The following methods have been removed and replaced by CurrentSource/OriginalSource management:

- `SetSourceFile(file *os.File, offset, size int64)` - Use `SetCurrentSource` instead
- `GetSourceFile() (*os.File, int64, int64)` - Use `GetCurrentSource` instead
- `SetTempPath(path string)` - Use `SetCurrentSource` with a FileSource that has IsTempFile=true
- `GetTempPath() string` - Use `GetCurrentSource().FilePath` after checking IsTempFile

## 5. Path Management

Cross-Reference: For path metadata structures, PathMetadataEntry tag population, and package-level path management, see [Package Metadata API - Path Metadata System](api_metadata.md#8-pathmetadata-system).

**Note**: PathEntry is defined in the generics package. See [Generic Types and Patterns - PathEntry](api_generics.md#13-pathentry-type) for complete specification. Path metadata (permissions, timestamps, ownership, tags) is stored separately in PathMetadataEntry structures. See [Package Metadata API - Path Metadata](api_metadata.md#8-pathmetadata-system) for details.

**Cross-Reference**: For complete package path semantics, validation rules, and normalization requirements, see [Package Path Semantics](api_core.md#2-package-path-semantics).

**Path Storage and Display**: All paths have different formats for internal storage versus user display.
Internal storage uses a leading `/` to ensure full path references within the package (e.g., `/path/to/file.txt`).
The leading `/` indicates the package root, not the OS filesystem root.
When displaying paths to end users, the leading `/` MUST be stripped (e.g., displayed as `path/to/file.txt` on Unix-like systems or `path\to\file.txt` on Windows).
File listings, path displays, and extraction operations show paths without leading `/`.
See [Package Path Semantics](api_core.md#2-package-path-semantics) for complete details.

### 5.1 FileEntry.HasSymlinks Method

```go
// HasSymlinks returns true if the FileEntry has any symlink paths.
func (fe *FileEntry) HasSymlinks() bool
```

### 5.2 FileEntry.GetSymlinkPaths Method

```go
// GetSymlinkPaths returns all symlink paths associated with this FileEntry.
func (fe *FileEntry) GetSymlinkPaths() []generics.PathEntry
```

### 5.3 FileEntry.GetPrimaryPath Method

```go
// GetPrimaryPath returns the primary path in display format (no leading slash).
//
// The returned string MUST use forward slashes as separators.
// For platform-specific filesystem display, convert manually or use path conversion utilities.
func (fe *FileEntry) GetPrimaryPath() string
```

### 5.4 FileEntry.ResolveAllSymlinks Method

```go
// ResolveAllSymlinks resolves all symlink paths to their target paths.
func (fe *FileEntry) ResolveAllSymlinks() []string
```

### 5.5 FileEntry.AssociateWithPathMetadata Method

```go
// PathMetadataEntry association methods for FileEntry.
// AssociateWithPathMetadata associates this FileEntry with a PathMetadataEntry.
// The association is based on matching paths between FileEntry.Paths and PathMetadataEntry.Path.
// Returns *PackageError on failure.
func (fe *FileEntry) AssociateWithPathMetadata(pme *PathMetadataEntry) error
```

### 5.6 FileEntry.GetPathMetadataForPath Method

```go
// GetPathMetadataForPath returns the PathMetadataEntry associated with a specific path in this FileEntry.
// Returns nil if no PathMetadataEntry is associated with the given path.
func (fe *FileEntry) GetPathMetadataForPath(path string) *PathMetadataEntry
```

### 5.7 FileEntry.GetPaths Method

```go
// GetPaths returns all paths associated with this FileEntry.
func (fe *FileEntry) GetPaths() []generics.PathEntry
```

### 5.8 FileEntry.GetFileID Method

```go
// GetFileID returns the unique file identifier.
func (fe *FileEntry) GetFileID() uint64
```

### 5.9 FileEntry.GetParentPath Method

```go
// GetParentPath returns the parent directory path for the primary path.
func (fe *FileEntry) GetParentPath() string
```

### 5.10 FileEntry.GetDirectoryDepth Method

```go
// GetDirectoryDepth returns the depth of the primary path in the hierarchy.
func (fe *FileEntry) GetDirectoryDepth() int
```

### 5.11 FileEntry.IsRootRelative Method

```go
// IsRootRelative returns true if the primary path is root-relative (no parent path).
func (fe *FileEntry) IsRootRelative() bool
```

## 6. Marshaling

This section describes marshaling operations for FileEntry.

### 6.1 Marshal Methods

This section describes marshal methods for FileEntry.

#### 6.1.1 FileEntry.MarshalMeta Method

```go
// MarshalMeta marshals the FileEntry metadata (header + variable data) to bytes.
// Marshals the complete FileEntry metadata structure including:
//   - Fixed 64-byte header
//   - Path entries
//   - Hash data
//   - Optional data (including tags)
// Returns *PackageError on failure.
func (fe *FileEntry) MarshalMeta() ([]byte, error)
```

#### 6.1.2 FileEntry.MarshalData Method

```go
// MarshalData marshals the FileEntry data (file content) to bytes.
// Marshals the file data content (already processed with compression/encryption)
// Returns *PackageError on failure.
func (fe *FileEntry) MarshalData() ([]byte, error)
```

#### 6.1.3 FileEntry.Marshal Method

```go
// Marshal marshals both FileEntry metadata and data.
// Returns metadata and data as separate byte slices for flexible writing.
// Returns *PackageError on failure.
func (fe *FileEntry) Marshal() (metadata, data []byte, err error)
```

### 6.2 WriteTo Methods

This section describes WriteTo methods for FileEntry.

#### 6.2.1 FileEntry.WriteMetaTo Method

```go
// WriteMetaTo writes the FileEntry metadata to a writer.
// Implements efficient streaming for large metadata.
// Returns *PackageError on failure.
func (fe *FileEntry) WriteMetaTo(w io.Writer) (int64, error)
```

#### 6.2.2 FileEntry.WriteDataTo Method

```go
// WriteDataTo writes the FileEntry data to a writer.
// Implements efficient streaming for large files.
// Returns *PackageError on failure.
func (fe *FileEntry) WriteDataTo(w io.Writer) (int64, error)
```

#### 6.2.3 FileEntry.WriteTo Method

```go
// WriteTo writes both metadata and data to a writer.
// Implements io.WriterTo interface.
// Returns *PackageError on failure.
func (fe *FileEntry) WriteTo(w io.Writer) (int64, error)
```

### 6.3 UnmarshalFileEntry Function

```go
// UnmarshalFileEntry unmarshals a FileEntry from binary data.
// Unmarshals the FileEntry with proper tag synchronization.
func UnmarshalFileEntry(data []byte) (*FileEntry, error)
```

### 6.4 FileEntry Marshaling Purpose

Provides marshaling methods for FileEntry metadata and data, supporting both byte-slice and streaming writer interfaces.

### 6.5 FileEntry Marshaling Usage Notes

- The byte-slice methods (`MarshalMeta`, `MarshalData`, `Marshal`) are convenient for small files or when you need the data in memory.
- The writer-based methods (`WriteMetaTo`, `WriteDataTo`, `WriteTo`) provide memory-efficient streaming alternatives suitable for large files.
- Writer methods follow the same pattern as `PackageComment.WriteTo` for consistency.
- Choose byte-slice methods for simplicity, or writer methods for memory efficiency with large files.

### 6.6 FileEntry Binary Format Methods

This section describes FileEntry methods tied directly to the binary format contract.

#### 6.6.1 FileEntry.ReadFrom Method

```go
// ReadFrom reads FileEntry metadata from a reader.
// Implements io.ReaderFrom.
// Returns *PackageError on failure.
func (fe *FileEntry) ReadFrom(r io.Reader) (int64, error)
```

#### 6.6.2 FileEntry.Validate Method

```go
// Validate validates the FileEntry state.
// Returns *PackageError on failure.
func (fe *FileEntry) Validate() error
```

#### 6.6.3 FileEntry.FixedSize Method

```go
// FixedSize returns the size of the fixed FileEntry metadata section in bytes.
func (fe *FileEntry) FixedSize() int
```

#### 6.6.4 FileEntry.VariableSize Method

```go
// VariableSize returns the size of the variable-length FileEntry metadata section in bytes.
func (fe *FileEntry) VariableSize() int
```

#### 6.6.5 FileEntry.TotalSize Method

```go
// TotalSize returns the total size of the FileEntry metadata (fixed + variable) in bytes.
func (fe *FileEntry) TotalSize() int
```

## 7. FileEntry Properties

This section describes properties and accessors for FileEntry.

### 7.1 FileEntry.IsCompressed Method

```go
// IsCompressed checks if the file is compressed.
func (entry *FileEntry) IsCompressed() bool
```

### 7.2 FileEntry.HasEncryptionKey Method

```go
// HasEncryptionKey checks if the file has an encryption key set
func (entry *FileEntry) HasEncryptionKey() bool
```

### 7.3 FileEntry.GetEncryptionType Method

```go
// GetEncryptionType returns the encryption type used for this file.
func (entry *FileEntry) GetEncryptionType() EncryptionType
```

### 7.4 FileEntry.IsEncrypted Method

```go
// IsEncrypted checks if the file is encrypted
func (entry *FileEntry) IsEncrypted() bool
```

### 7.5 FileEntry Properties Purpose

Provides access to FileEntry properties.

### 7.6 IsCompressed Returns

`bool` indicating if the file is compressed

### 7.7 HasEncryptionKey Returns

`bool` indicating if the file has an encryption key

### 7.8 GetEncryptionType Returns

`EncryptionType` value indicating the encryption algorithm used (e.g., EncryptionNone, EncryptionAES256GCM, EncryptionMLKEM768)

### 7.9 IsEncrypted Returns

`bool` indicating if the file is encrypted (equivalent to `GetEncryptionType() != EncryptionNone`)

## 8. FileEntry Compression

This section describes compression operations for FileEntry.

### 8.1 FileEntry.Compress Method

```go
// Compress compresses the FileEntry content using the specified compression type.
// Returns *PackageError on failure.
func (fe *FileEntry) Compress(ctx context.Context, compressionType uint8) error
```

### 8.2 FileEntry.Decompress Method

```go
// Decompress decompresses the FileEntry content
// Returns *PackageError on failure.
func (fe *FileEntry) Decompress(ctx context.Context) error
```

### 8.3 FileEntry.GetCompressionInfo Method

```go
// GetCompressionInfo gets compression information for the FileEntry.
func (fe *FileEntry) GetCompressionInfo() (*FileCompressionInfo, error)
```

**Note**: The `FileCompressionInfo` type is defined in [File Compression API](api_file_mgmt_compression.md#41-filecompressioninfo-struct-definition).

### 8.4 FileEntry Compression Purpose

Manages compression for individual file entries.

For complete documentation including Package-level convenience methods, see [File Compression API](api_file_mgmt_compression.md).

### 8.5 Compress Parameters

- `ctx`: Context for cancellation and timeout handling
- `compressionType`: Compression algorithm to use

### 8.6 Decompress Parameters

- `ctx`: Context for cancellation and timeout handling

### 8.7 GetCompressionInfo Parameters

None (operates on the FileEntry instance).

### 8.8 FileEntry Compression Error Conditions

- `ErrTypeValidation`: FileEntry is invalid
- `ErrTypeValidation`: File is already compressed (for Compress)
- `ErrTypeValidation`: File is not compressed (for Decompress)
- `ErrTypeCompression`: Failed to compress or decompress file content
- `ErrTypeContext`: Context was cancelled
- `ErrTypeContext`: Context timeout exceeded

### 8.9 FileEntry Compression Usage Notes

Compression and decompression become durable only after Write, SafeWrite, or FastWrite completes successfully.

For Package-level convenience methods that operate on files by path, see [File Compression API](api_file_mgmt_compression.md).

## 9. FileEntry Encryption

This section describes encryption operations for FileEntry.

### 9.1 FileEntry.SetEncryptionKey Method

```go
// SetEncryptionKey sets the encryption key for the file
func (entry *FileEntry) SetEncryptionKey(key *EncryptionKey) error
```

### 9.2 FileEntry.Encrypt Method

```go
// Encrypt encrypts data using the file's encryption key.
func (entry *FileEntry) Encrypt(data []byte) ([]byte, error)
```

### 9.3 FileEntry.Decrypt Method

```go
// Decrypt decrypts data using the file's encryption key.
func (entry *FileEntry) Decrypt(data []byte) ([]byte, error)
```

### 9.4 FileEntry.UnsetEncryptionKey Method

```go
// UnsetEncryptionKey removes the encryption key from the file
func (entry *FileEntry) UnsetEncryptionKey()
```

### 9.5 FileEntry Encryption Purpose

Manages encryption for individual file entries.

### 9.6 SetEncryptionKey Parameters

- `key`: Encryption key to set

### 9.7 Encrypt/Decrypt Parameters

- `data`: Data to encrypt or decrypt

### 9.8 FileEntry Encryption Error Conditions

- `ErrTypeEncryption`: Invalid encryption key
- `ErrTypeEncryption`: Encryption operation failed
- `ErrTypeEncryption`: Decryption operation failed

### 9.9 FileEntry Encryption Usage Notes

SetEncryptionKey and UnsetEncryptionKey manage the encryption key for the FileEntry.

## 10. FileEntry Data Management

This section describes data management operations for FileEntry.

### 10.1 FileEntry.LoadData Method

```go
// LoadData loads the file data into memory.
func (fe *FileEntry) LoadData(ctx context.Context) error
```

### 10.2 FileEntry.ProcessData Method

```go
// ProcessData processes the file data (compression, encryption, etc.)
func (fe *FileEntry) ProcessData(ctx context.Context) error
```

### 10.3 FileEntry Data Management Purpose

Manages file data loading and processing.

### 10.4 LoadData Behavior

- Loads file content from package into memory
- Prepares data for access and processing
- May trigger decompression or decryption

### 10.5 ProcessData Behavior

- Applies compression and encryption to file data
- Updates FileEntry metadata
- Prepares data for storage

### 10.6 FileEntry Data Management Error Conditions

- `ErrTypeIO`: I/O error during data operations
- `ErrTypeEncryption`: Failed to decrypt data
- `ErrTypeCompression`: Failed to decompress data

### 10.7 FileEntry Data Management Usage Notes

LoadData and UnloadData manage the file's content in memory, while ProcessData applies compression/encryption.

## 11. HashEntry Struct

```go
// HashEntry represents a hash with type and purpose.
type HashEntry struct {
    Type    HashType    // Hash algorithm type
    Purpose HashPurpose // Hash purpose
    Data    []byte      // Hash data
}
```

## 12. HashType Type

```go
// HashType represents hash algorithm types.
type HashType uint8
const (
    HashTypeSHA256   HashType = 0x00  // SHA-256 (32 bytes) - Standard cryptographic hash
    HashTypeSHA512   HashType = 0x01  // SHA-512 (64 bytes) - Stronger cryptographic hash
    HashTypeBLAKE3   HashType = 0x02  // BLAKE3 (32 bytes) - Fast cryptographic hash
    HashTypeXXH3     HashType = 0x03  // XXH3 (8 bytes) - Ultra-fast non-cryptographic hash
    HashTypeBLAKE2b  HashType = 0x04  // BLAKE2b (64 bytes) - Cryptographic hash with configurable output
    HashTypeBLAKE2s  HashType = 0x05  // BLAKE2s (32 bytes) - Cryptographic hash optimized for 32-bit systems
    HashTypeSHA3_256 HashType = 0x06  // SHA-3-256 (32 bytes) - SHA-3 family hash
    HashTypeSHA3_512 HashType = 0x07  // SHA-3-512 (64 bytes) - SHA-3 family hash
    HashTypeCRC32    HashType = 0x08  // CRC32 (4 bytes) - Fast checksum for error detection
    HashTypeCRC64    HashType = 0x09  // CRC64 (8 bytes) - Stronger checksum for error detection
)
```

## 13. HashPurpose Type

```go
// HashPurpose represents hash purposes
type HashPurpose uint8
const (
    HashPurposeContentVerification HashPurpose = 0x00  // Content verification - Verify file content integrity
    HashPurposeDeduplication       HashPurpose = 0x01  // Deduplication - Identify duplicate content
    HashPurposeIntegrityCheck      HashPurpose = 0x02  // Integrity check - General integrity verification
    HashPurposeFastLookup          HashPurpose = 0x03  // Fast lookup - Quick content identification
    HashPurposeErrorDetection      HashPurpose = 0x04  // Error detection - Detect data corruption
)
```

## 14. TagValueType Type

```go
// TagValueType represents the type of a tag value.
type TagValueType uint8
const (
    // Basic Types
    TagValueTypeString      TagValueType = 0x00  // String value
    TagValueTypeInteger     TagValueType = 0x01  // 64-bit signed integer
    TagValueTypeFloat       TagValueType = 0x02  // 64-bit floating point number
    TagValueTypeBoolean     TagValueType = 0x03  // Boolean value

    // Structured Data
    TagValueTypeJSON        TagValueType = 0x04  // JSON-encoded object or array
    TagValueTypeYAML        TagValueType = 0x05  // YAML-encoded data
    TagValueTypeStringList  TagValueType = 0x06  // Comma-separated list of strings

    // Identifiers
    TagValueTypeUUID        TagValueType = 0x07  // UUID string
    TagValueTypeHash        TagValueType = 0x08  // Hash/checksum string
    TagValueTypeVersion     TagValueType = 0x09  // Semantic version string

    // Time
    TagValueTypeTimestamp   TagValueType = 0x0A  // ISO8601 timestamp

    // Network/Communication
    TagValueTypeURL         TagValueType = 0x0B  // URL string
    TagValueTypeEmail       TagValueType = 0x0C  // Email address

    // File System
    TagValueTypePath        TagValueType = 0x0D  // File system path
    TagValueTypeMimeType    TagValueType = 0x0E  // MIME type string

    // Localization
    TagValueTypeLanguage    TagValueType = 0x0F  // Language code (ISO 639-1)

    // NovusPack Special Files
    TagValueTypeNovusPackMetadata TagValueType = 0x10  // NovusPack special metadata file reference
)
```

## 15. ProcessingState Type

```go
// ProcessingState defines the current state of file data transformations.
// This tracks what processing has been applied to the data to inform operations
// what additional processing is needed. This uses a data-state model rather than
// a workflow model.
type ProcessingState uint8

const (
    ProcessingStateRaw                     ProcessingState = iota // Raw (unprocessed) data
    ProcessingStateCompressed              ProcessingState         // Compressed but not encrypted
    ProcessingStateEncrypted               ProcessingState         // Encrypted but not compressed
    ProcessingStateCompressedAndEncrypted  ProcessingState         // Both compressed and encrypted
)
```

## 16. FileSource Structure

```go
// FileSource represents a source location for file data (original or intermediate).
type FileSource struct {
    File         *os.File  // File handle (may be nil if closed)
    FilePath     string    // Path to file (for reopening if needed)
    Offset       int64     // Offset in file where data starts
    Size         int64     // Size of data at this location
    IsPackage    bool      // True if this is the package file itself
    IsTempFile   bool      // True if this is a temporary file
    IsExternal   bool      // True if this is an external source file
}
```

## 17. OptionalData Structure

```go
// OptionalData represents structured optional data for a FileEntry.
type OptionalData struct {
    Tags                []*Tag[any]         // Per-file tags data (DataType 0x00)
    PathEncoding        *uint8              // Path encoding type (DataType 0x01)
    PathFlags           *uint8              // Path handling flags (DataType 0x02)
    CompressionDictID   *uint32             // Dictionary ID for solid compression (DataType 0x03)
    SolidGroupID        *uint32             // Solid compression group ID (DataType 0x04)
    FileSystemFlags     *uint16             // File system specific flags (DataType 0x05)
    WindowsAttributes   *uint32             // Windows file attributes (DataType 0x06)
    ExtendedAttributes  map[string]string   // Unix extended attributes (DataType 0x07)
    ACLData             []byte              // Access Control List data (DataType 0x08)
    CustomData          map[uint8][]byte    // Custom data for reserved types (0x09-0xFF)
}
```

## 18. OptionalDataType Type

```go
// OptionalDataType represents the type of optional data
type OptionalDataType uint8
const (
    OptionalDataTypeTags                OptionalDataType = 0x00  // TagsData
    OptionalDataTypePathEncoding        OptionalDataType = 0x01  // PathEncoding
    OptionalDataTypePathFlags           OptionalDataType = 0x02  // PathFlags
    OptionalDataTypeCompressionDictID   OptionalDataType = 0x03  // CompressionDictionaryID
    OptionalDataTypeSolidGroupID        OptionalDataType = 0x04  // SolidGroupID
    OptionalDataTypeFileSystemFlags     OptionalDataType = 0x05  // FileSystemFlags
    OptionalDataTypeWindowsAttributes   OptionalDataType = 0x06  // WindowsAttributes
    OptionalDataTypeExtendedAttributes  OptionalDataType = 0x07  // ExtendedAttributes
    OptionalDataTypeACLData             OptionalDataType = 0x08  // ACLData
)
```

## 19. Tag Generic Type

This section describes the generic Tag type used throughout the FileEntry API.

### 19.1 Tag Struct

```go
// Tag represents a type-safe tag with a specific value type.
type Tag[T any] struct {
    Key   string
    Value T
    Type  TagValueType
}
```

### 19.2 NewTag Function

```go
// NewTag creates a new type-safe tag with the specified key, value, and type.
func NewTag[T any](key string, value T, tagType TagValueType) *Tag[T]
```

### 19.3 Tag[T].GetValue Method

```go
// GetValue returns the type-safe value of the tag.
func (t *Tag[T]) GetValue() T
```

### 19.4 Tag[T].SetValue Method

```go
// SetValue sets the type-safe value of the tag.
func (t *Tag[T]) SetValue(value T)
```
