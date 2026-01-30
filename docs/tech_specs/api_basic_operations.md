# NovusPack Technical Specifications - Basic Operations API

## Table of Contents

- [Table of Contents](#table-of-contents)
- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Context Integration](#1-context-integration)
- [2. Go API v1 Package Organization](#2-go-api-v1-package-organization)
  - [2.1 Module Path](#21-module-path)
    - [2.1.1 API Package Structure](#211-api-package-structure)
    - [2.1.2 Root Package Purpose](#212-root-package-purpose)
    - [2.1.3 Root Package Example Import](#213-root-package-example-import)
    - [2.1.4 File Format Package: `fileformat`](#214-file-format-package-fileformat)
    - [2.1.5 File Format Package Key Types](#215-file-format-package-key-types)
    - [2.1.6 File Format Package Example Import](#216-file-format-package-example-import)
    - [2.1.7 Metadata Package Purpose](#217-metadata-package-purpose)
    - [2.1.8 Metadata Package Imports](#218-metadata-package-imports)
    - [2.1.9 Generics Package: `generics`](#219-generics-package-generics)
    - [2.1.10 Generics Package Key Types](#2110-generics-package-key-types)
    - [2.1.11 Generics Package Example Import](#2111-generics-package-example-import)
    - [2.1.12 Error Handling Package Purpose](#2112-error-handling-package-purpose)
    - [2.1.13 Error Handling Package Imports](#2113-error-handling-package-imports)
    - [2.1.14 Signatures Package: `signatures`](#2114-signatures-package-signatures)
    - [2.1.15 Signatures Package Key Types](#2115-signatures-package-key-types)
    - [2.1.16 Signatures Package Example Import](#2116-signatures-package-example-import)
    - [2.1.17 Internal Package Purpose](#2117-internal-package-purpose)
  - [2.2 Import Patterns](#22-import-patterns)
    - [Recommended Import Pattern](#recommended-import-pattern)
    - [Direct Subpackage Imports](#direct-subpackage-imports)
    - [Avoiding Circular Dependencies](#avoiding-circular-dependencies)
    - [2.2.1 Dependency Graph](#221-dependency-graph)
  - [2.3 Package Aliases](#23-package-aliases)
- [3. Package Structure and Loading](#3-package-structure-and-loading)
  - [3.1 Package Implementation Structure](#31-package-implementation-structure)
  - [3.2 Package Loading Process](#32-package-loading-process)
    - [3.2.1 `Package.loadSpecialMetadataFiles` Method](#321-packageloadspecialmetadatafiles-method)
    - [3.2.2 `Package.loadPathMetadata` Method](#322-packageloadpathmetadata-method)
    - [3.2.3 `Package.updateFilePathAssociations` Method](#323-packageupdatefilepathassociations-method)
  - [3.3 Package Implementation Details](#33-package-implementation-details)
    - [3.3.1 Package Structure Implementation](#331-package-structure-implementation)
    - [3.3.2 Data Loading Strategies](#332-data-loading-strategies)
    - [3.3.3 State Management](#333-state-management)
    - [3.3.4 Memory Management](#334-memory-management)
    - [3.3.5 Data Relationships](#335-data-relationships)
    - [3.3.6 Thread Safety and Concurrency](#336-thread-safety-and-concurrency)
    - [3.3.7 Resource Lifecycle](#337-resource-lifecycle)
- [4. Package Format Constants](#4-package-format-constants)
  - [4.1 Format Constants](#41-format-constants)
- [5. Package Lifecycle Operations](#5-package-lifecycle-operations)
  - [5.1 Package Lifecycle - Always Use Defer for Cleanup](#51-package-lifecycle---always-use-defer-for-cleanup)
  - [5.2 Package Lifecycle - Check Package State Before Operations](#52-package-lifecycle---check-package-state-before-operations)
  - [5.3 Package Lifecycle - Use appropriate context timeouts](#53-package-lifecycle---use-appropriate-context-timeouts)
- [6. `NewPackage` Constructor Function](#6-newpackage-function)
  - [6.1 NewPackage Behavior](#61-newpackage-behavior)
  - [6.2 NewPackage Example Usage](#62-newpackage-example-usage)
- [7. `NewPackageWithOptions` Constructor Function](#7-newpackagewithoptions-function)
  - [7.1 NewPackageWithOptions Parameters](#71-newpackagewithoptions-parameters)
  - [7.2 NewPackageWithOptions Behavior](#72-newpackagewithoptions-behavior)
  - [7.3 NewPackageWithOptions Error Conditions](#73-newpackagewithoptions-error-conditions)
  - [7.4 NewPackageWithOptions Example Usage](#74-newpackagewithoptions-example-usage)
  - [7.5 NewPackageWithOptions Without Path](#75-newpackagewithoptions-without-path)
  - [7.6 `CreateOptions` Structure](#76-createoptions-structure)
- [8. `Package.SetTargetPath` Method](#8-packagesettargetpath-method)
  - [8.1 Package.SetTargetPath Parameters](#81-packagesettargetpath-parameters)
  - [8.2 Package.SetTargetPath Behavior](#82-packagesettargetpath-behavior)
  - [8.3 Package.SetTargetPath Method Error Conditions](#83-packagesettargetpath-method-error-conditions)
  - [8.4 Package.SetTargetPath Example Usage](#84-packagesettargetpath-example-usage)
  - [8.5 Package.SetTargetPath vs NewPackageWithOptions](#85-packagesettargetpath-vs-newpackagewithoptions)
- [9. Package Configuration](#9-package-configuration)
  - [9.1 `PackageConfig` Structure](#91-packageconfig-structure)
    - [9.1.1 PackageConfig Fields](#911-packageconfig-fields)
  - [9.2 PackageConfig Backward Compatibility](#92-packageconfig-backward-compatibility)
  - [9.3 `PathHandling` Type](#93-pathhandling-type)
- [10. `OpenPackage`](#10-openpackage-function)
  - [10.1 OpenPackage Parameters](#101-openpackage-parameters)
  - [10.2 OpenPackage Behavior](#102-openpackage-behavior)
  - [10.3 OpenPackage Method Error Conditions](#103-openpackage-method-error-conditions)
    - [10.3.1 OpenPackage Example Usage](#1031-openpackage-example-usage)
- [11. Opening Packages as Read-Only](#11-opening-packages-as-read-only)
  - [11.1 Read-Only Enforcement Mechanism](#111-read-only-enforcement-mechanism)
    - [11.1.1 Mutating Methods That Must Be Rejected](#1111-mutating-methods-that-must-be-rejected)
  - [11.2 `OpenPackageReadOnly` Function](#112-openpackagereadonly-function)
    - [11.2.1 OpenPackageReadOnly Behavior](#1121-openpackagereadonly-behavior)
    - [11.2.2 OpenPackageReadOnly Method Error Conditions](#1122-openpackagereadonly-method-error-conditions)
  - [11.3 `readOnlyPackage` Structure](#113-readonlypackage-struct)
  - [11.4 `readOnlyPackage.readOnlyError` Helper](#114-readonlypackagereadonlyerror-method)
  - [11.5 `ReadOnlyErrorContext` Structure](#115-readonlyerrorcontext-structure)
  - [11.6 `readOnlyPackage` Implementation Methods](#116-readonlypackage-implementation-methods)
- [12. `OpenBrokenPackage` Function](#12-openbrokenpackage-function)
- [13. `Package.Close` Method](#13-packageclose-method)
  - [13.1 `Package.Close` Behavior](#131-packageclose-behavior)
  - [13.2 `Package.Close` Method Error Conditions](#132-packageclose-method-error-conditions)
  - [13.3 `Package.Close` Example Usage](#133-packageclose-example-usage)
- [14. `Package.CloseWithCleanup` Method](#14-packageclosewithcleanup-method)
  - [14.1 `Package.CloseWithCleanup` Behavior](#141-packageclosewithcleanup-behavior)
- [15. `Package.Validate` Method](#15-packagevalidate-method)
  - [15.1 `Package.Validate` Behavior](#151-packagevalidate-behavior)
  - [15.2 `Package.Validate` Method Error Conditions](#152-packagevalidate-method-error-conditions)
  - [15.3 `Package.Validate` Example Usage](#153-packagevalidate-example-usage)
- [16. `Package.Defragment` Method](#16-packagedefragment-method)
  - [16.1 `Package.Defragment` Behavior](#161-packagedefragment-behavior)
  - [16.2 `Package.Defragment` Error Conditions](#162-packagedefragment-error-conditions)
  - [16.3 `Package.Defragment` Example Usage](#163-packagedefragment-example-usage)
- [17. `Package.GetInfo` Method](#17-packagegetinfo-method)
  - [17.1 `Package.GetInfo` Error Conditions](#171-packagegetinfo-error-conditions)
  - [17.2 `Package.GetInfo` Example Usage](#172-packagegetinfo-example-usage)
- [18. Header Inspection](#18-header-inspection)
  - [18.1 Header Inspection Use Cases](#181-header-inspection-use-cases)
  - [18.2 ReadHeader vs ReadHeaderFromPath](#182-readheader-vs-readheaderfrompath)
  - [18.3 `ReadHeader` Function](#183-readheader-function)
  - [18.4 `ReadHeaderFromPath` Function](#184-readheaderfrompath-function)
    - [18.4.1 `ReadHeaderFromPath` Parameters](#1841-readheaderfrompath-parameters)
    - [18.4.2 `ReadHeaderFromPath` Error Conditions](#1842-readheaderfrompath-error-conditions)
    - [18.4.3 `ReadHeaderFromPath` Example Usage](#1843-readheaderfrompath-example-usage)
  - [18.5 `Package.ReadHeader` Method](#185-packagereadheader-method)
    - [18.5.1 `Package.ReadHeader` Parameters](#1851-packagereadheader-parameters)
    - [18.5.2 `Package.ReadHeader` Error Conditions](#1852-packagereadheader-error-conditions)
  - [18.6 `Package.IsOpen` Method](#186-packageisopen-method)
  - [18.7 `Package.IsReadOnly` Method](#187-packageisreadonly-method)
  - [18.8 `Package.GetPath` Method](#188-packagegetpath-method)
- [19. Package Session Base Management](#19-package-session-base-management)
  - [19.1 Session Base for File Addition](#191-session-base-for-file-addition)
  - [19.2 Session Base for File Extraction](#192-session-base-for-file-extraction)
  - [19.3 Session Base Lifecycle](#193-session-base-lifecycle)
  - [19.4 `Package.SetSessionBase` Method](#194-packagesetsessionbase-method)
    - [19.4.1 `Package.SetSessionBase` Parameters](#1941-packagesetsessionbase-parameters)
    - [19.4.2 Package.SetSessionBase Returns](#1942-packagesetsessionbase-returns)
    - [19.4.3 Package.SetSessionBase Example Usage](#1943-packagesetsessionbase-example-usage)
  - [19.5 Package.GetSessionBase Method](#195-packagegetsessionbase-method)
    - [19.5.1 Package.GetSessionBase Returns](#1951-packagegetsessionbase-returns)
    - [19.5.2 Package.GetSessionBase Example](#1952-packagegetsessionbase-example)
  - [19.6 `Package.ClearSessionBase` Method](#196-packageclearsessionbase-method)
    - [19.6.1 Package.ClearSessionBase Example](#1961-packageclearsessionbase-example)
  - [19.7 Package.HasSessionBase Method](#197-packagehassessionbase-method)
    - [19.7.1 Package.HasSessionBase Returns](#1971-packagehassessionbase-returns)
    - [19.7.2 Package.HasSessionBase Example](#1972-packagehassessionbase-example)
- [20. Structured Error System](#20-structured-error-system)
  - [20.1 Error Types Used](#201-error-types-used)
  - [20.2 `PackageErrorContext` Structure](#202-packageerrorcontext-structure)
  - [20.3 `SecurityErrorContext` Structure](#203-securityerrorcontext-structure)
  - [20.4 `IOErrorContext` Structure](#204-ioerrorcontext-structure)
  - [20.5 Creating Errors with Context](#205-creating-errors-with-context)
  - [20.6 Error Inspection](#206-error-inspection)
- [21. Error Handling Best Practices](#21-error-handling-best-practices)
  - [21.1 Always Check for Errors](#211-always-check-for-errors)
  - [21.2 Use Structured Errors for Better Debugging](#212-use-structured-errors-for-better-debugging)
  - [21.3 Use Context for Cancellation](#213-use-context-for-cancellation)
  - [21.4 Handle Different Error Types Appropriately](#214-handle-different-error-types-appropriately)
  - [21.5 Clean Up Resources](#215-clean-up-resources)
  - [21.6 Wrap Errors with Context](#216-wrap-errors-with-context)
- [22. Resource Management](#22-resource-management)
  - [22.1 Use Context for Resource Management](#221-use-context-for-resource-management)
  - [22.2 Handle Cleanup Errors Gracefully](#222-handle-cleanup-errors-gracefully)

## 0. Overview

This document defines the basic package operations for the NovusPack system, covering the fundamental lifecycle operations of creating, opening, and closing packages.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Package Metadata API](api_metadata.md) - Comment management, AppID/VendorID, and metadata operations
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation

## 1. Context Integration

See [Core API - Context Integration](api_core.md#02-context-integration) for the complete context integration specification.

The NovusPack Basic Operations API follows the same context integration patterns as the core API.

These methods assume the following imports:

```go
import (
    "context"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)
// FileEntry, PathMetadataEntry, PackageInfo, PackageHeader, PackageIndex from other API modules
```

## 2. Go API V1 Package Organization

The NovusPack Go API is organized into multiple subpackages within `api/go` to provide clear separation of concerns and avoid circular dependencies.

### 2.1 Module Path

The module path for all packages is:

```text
github.com/novus-engine/novuspack/api/go
```

#### 2.1.1 API Package Structure

The API is organized into the following subpackages:

##### 2.1.1.1. Root Package: Novuspack

**Import Path**: `github.com/novus-engine/novuspack/api/go`

#### 2.1.2 Root Package Purpose

The root package provides the main public API and re-exports types from domain-specific subpackages for convenience.

##### 2.1.2.1 Root Package Key Features

- Re-exports all public types from subpackages
- Provides unified import path for most use cases
- Implements the main `Package` interface and lifecycle operations

#### 2.1.3 Root Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go"
```

##### 2.1.3.1 Re-exported Types

- From `fileformat`: `PackageHeader`, `FileIndex`, `IndexEntry`
- From `generics`: `PathEntry`, `Option`, `Result`, `Tag`, `Strategy`, `Validator`
- From `metadata`: `PackageComment`, `PackageInfo`, `FileEntry`, `HashEntry`, `OptionalDataEntry`, `ProcessingState`
- From `signatures`: `Signature`, `SignatureInfo`
- From `pkgerrors`: `ErrorType`, `PackageError`

#### 2.1.4. File Format Package: Fileformat

**Import Path**: `github.com/novus-engine/novuspack/api/go/fileformat`

##### 2.1.4.1 File Format Package Purpose

Provides structures and operations for the NovusPack binary file format, including headers, file entries, and data management.

#### 2.1.5 File Format Package Key Types

- `PackageHeader` - Package header structure
- `FileIndex` - File index structure
- `IndexEntry` - File index entry structure

##### 2.1.5.1 File Format Package Imports

- `github.com/novus-engine/novuspack/api/go/generics` - For `PathEntry` and generic types
- `github.com/novus-engine/novuspack/api/go/metadata` - For `PathMetadataEntry` associations
- `github.com/novus-engine/novuspack/api/go/pkgerrors` - For error handling

#### 2.1.6 File Format Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go/fileformat"
```

##### 2.1.6.1. Metadata Package: Metadata

**Import Path**: `github.com/novus-engine/novuspack/api/go/metadata`

#### 2.1.7 Metadata Package Purpose

Provides metadata structures including package information, path metadata, directory entries, and security levels.

##### 2.1.7.1 Metadata Package Key Types

- `PackageInfo` - Package metadata information
- `PathMetadataEntry` - Path metadata with inheritance and filesystem properties
- `FileEntry` - FileEntry with paths, hashes, and optional data
- `HashEntry` - Hash entry structure
- `OptionalDataEntry` - Optional data entry structure
- `ProcessingState` - File processing state enumeration
- `PackageComment` - Package comment structure

#### 2.1.8 Metadata Package Imports

- `github.com/novus-engine/novuspack/api/go/generics` - For `PathEntry`, `Tag`, and interface types
- `github.com/novus-engine/novuspack/api/go/pkgerrors` - For error handling

**Note**: `FileEntry`, `HashEntry`, `OptionalDataEntry`, and `ProcessingState` are defined in the `metadata` package, not the `fileformat` package.

##### 2.1.8.1 Metadata Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go/metadata"
```

#### 2.1.9. Generics Package: Generics

**Import Path**: `github.com/novus-engine/novuspack/api/go/generics`

##### 2.1.9.1 Generics Package Purpose

Provides generic types, patterns, and shared interfaces to avoid circular dependencies between domain packages.

#### 2.1.10 Generics Package Key Types

- `PathEntry` - Minimal path structure used by both `FileEntry` and `PathMetadataEntry`
- `Option[T]` - Type-safe optional values
- `Result[T, E]` - Type-safe result type
- `Tag[T]` - Generic tag structure
- `Strategy[T, U]` - Strategy pattern interface
- `Validator[T]` - Validation interface

##### 2.1.10.1 Generics Package Imports

- No internal NovusPack imports (standalone package)

#### 2.1.11 Generics Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go/generics"
```

##### 2.1.11.1. Error Handling Package: Pkgerrors

**Import Path**: `github.com/novus-engine/novuspack/api/go/pkgerrors`

#### 2.1.12 Error Handling Package Purpose

Provides structured error handling with typed errors and validation context.

##### 2.1.12.1 Error Handling Package Key Types

- `ErrorType` - Error type enumeration
- `PackageError` - Structured error type
- `ValidationErrorContext` - Validation error context structure

#### 2.1.13 Error Handling Package Imports

- Standard library only (no internal NovusPack imports)

##### 2.1.13.1 Error Handling Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go/pkgerrors"
```

#### 2.1.14. Signatures Package: Signatures

**Import Path**: `github.com/novus-engine/novuspack/api/go/signatures`

##### 2.1.14.1 Signatures Package Purpose

Provides digital signature structures and operations for package integrity verification.
Signature management and signature validation are deferred to v2.
V1 only enforces signed package immutability based on signature presence.

#### 2.1.15 Signatures Package Key Types

- `Signature` - Digital signature structure
- `SignatureInfo` - Signature information structure

##### 2.1.15.1 Signatures Package Imports

- `github.com/novus-engine/novuspack/api/go/pkgerrors` - For error handling

#### 2.1.16 Signatures Package Example Import

```go
import "github.com/novus-engine/novuspack/api/go/signatures"
```

##### 2.1.16.1. Internal Package: Internal

**Import Path**: `github.com/novus-engine/novuspack/api/go/internal`

#### 2.1.17 Internal Package Purpose

Provides internal helper functions used by the main package. This package is not part of the public API and should not be imported by external code.

##### 2.1.17.1 Internal Package Key Features

- Internal helper functions for file operations
- Not exported for external use

**Note**: This package is not part of the public API. Do not import it in external code.

### 2.2 Import Patterns

#### Recommended Import Pattern

For most use cases, import the root package:

```go
import "github.com/novus-engine/novuspack/api/go"
```

This provides access to all re-exported types and the main `Package` interface.

#### Direct Subpackage Imports

For advanced use cases or when you need access to package-specific functionality not re-exported, import subpackages directly:

```go
import (
    "github.com/novus-engine/novuspack/api/go/fileformat"
    "github.com/novus-engine/novuspack/api/go/metadata"
    "github.com/novus-engine/novuspack/api/go/pkgerrors"
)
```

#### Avoiding Circular Dependencies

The package structure is designed to avoid circular dependencies:

- `generics` package has no internal NovusPack imports and provides shared interfaces
- `pkgerrors` package has no internal NovusPack imports
- `metadata` package contains both `FileEntry` and `PathMetadataEntry`, which can directly reference each other since they're in the same package
- `signatures` only imports `pkgerrors`

#### 2.2.1 Dependency Graph

The following diagram shows the dependency relationships between packages:

```text
novuspack (root)
├── fileformat
│   ├── generics (shared types and interfaces)
│   └── pkgerrors
├── metadata
│   ├── generics (shared types and interfaces)
│   │   └── Contains: FileEntry, HashEntry, OptionalDataEntry, ProcessingState
│   └── pkgerrors
├── signatures
│   └── pkgerrors
├── generics (standalone, no internal imports)
└── pkgerrors (standalone, no internal imports)
```

### 2.3 Package Aliases

When importing multiple packages, use aliases to avoid naming conflicts:

```go
import (
    novuspack "github.com/novus-engine/novuspack/api/go"
    fileformat "github.com/novus-engine/novuspack/api/go/fileformat"
    metadata "github.com/novus-engine/novuspack/api/go/metadata"
)
```

## 3. Package Structure and Loading

This section describes the internal structure of packages and how they are loaded from disk.

### 3.1 Package Implementation Structure

The NovusPack API uses an interface-based design where the `Package` interface is implemented by the concrete `filePackage` struct.

The `Package` interface provides the public API for package operations, while `filePackage` contains the actual implementation.

The canonical `Package` interface definition is specified in [Core Package Interface API - Package Interface](api_core.md#11-package-interface).

The `Package` interface provides a unified API that combines PackageReader, PackageWriter, lifecycle operations, file management, metadata operations, compression operations, and session base management.

For session base management (used for both file addition and extraction operations), see [Package Session Base Management](#19-package-session-base-management).

### 3.2 Package Loading Process

When a package is opened, the following initialization sequence occurs:

1. **Load package header** - Validates magic number and version
2. **Load package info** - Retrieves metadata (comment, VendorID, AppID)
3. **Load file entries** - Reads file index and entry structures
4. **Load special metadata files** - Processes special file types (65000-65535)
5. **Load path metadata** - Parses path structure from YAML
6. **Update file-path associations** - Links files to their path metadata

The package opening functions are:

- [OpenPackage](#10-openpackage-function) - Opens an existing package from disk and validates the package structure during open
- [OpenPackageReadOnly](#112-openpackagereadonly-function) - Opens a package in read-only mode, enforcing read-only behavior for both in-memory modifications and writes to disk
- [OpenBrokenPackage](#12-openbrokenpackage-function) - Opens a package that may be invalid or partially corrupted, intended for repair workflows

The package loading process uses internal methods to complete initialization:

- [loadSpecialMetadataFiles](#321-packageloadspecialmetadatafiles-method) - Loads all special metadata files (file types 65000-65535)
- [loadPathMetadata](#322-packageloadpathmetadata-method) - Loads path metadata from special files and parses YAML structure
- [updateFilePathAssociations](#323-packageupdatefilepathassociations-method) - Links file entries to their corresponding path metadata entries

#### 3.2.1 Package.loadSpecialMetadataFiles Method

```go
// loadSpecialMetadataFiles loads all special metadata files
// Returns *PackageError on failure
func (p *Package) loadSpecialMetadataFiles(ctx context.Context) error
```

This method loads all special metadata files (file types 65000-65535) from the package.

#### 3.2.2 Package.loadPathMetadata Method

```go
// loadPathMetadata loads path metadata from special files
// Returns *PackageError on failure
func (p *Package) loadPathMetadata(ctx context.Context) error
```

This method loads path metadata from special files and parses the YAML structure.

#### 3.2.3 Package.updateFilePathAssociations Method

```go
// updateFilePathAssociations links files to their path metadata
// Returns *PackageError on failure
func (p *Package) updateFilePathAssociations(ctx context.Context) error
```

This method links file entries to their corresponding path metadata entries.

### 3.3 Package Implementation Details

This section documents the internal implementation details of the NovusPack package system, including data loading strategies, state management, memory management, and resource lifecycle.

**Note**: The canonical `filePackage` struct definition is specified in [Core Package Interface API - filePackage Implementation](api_core.md#111-filepackage-struct).

#### 3.3.1 Package Structure Implementation

The NovusPack API uses an interface-based design where the `Package` interface is implemented by the concrete `filePackage` struct.

See [Core Package Interface API - filePackage Implementation](api_core.md#111-filepackage-struct) for the complete struct definition and field descriptions.

#### 3.3.2 Data Loading Strategies

The NovusPack implementation uses a hybrid loading strategy that balances memory efficiency with performance.

##### 3.3.2.1 Eager Loading Immediate

The following data is loaded immediately when a package is opened:

- Package header: Validated and loaded from file start.
- File index: Loaded from header.IndexStart offset.
- Package info: Populated from header and index data.

##### 3.3.2.2 On-Demand Loading (Lazy)

The following data is loaded only when needed:

- File entries: FileEntries slice is initially empty.
  - Individual entries are loaded when file operations require them.
  - Entries are located using the index (FileID => Offset mapping).
  - This minimizes memory usage for packages with many files.

- Path metadata: PathMetadataEntries slice is initially empty.
  - Loaded when path metadata operations are performed.
  - Loaded from special metadata files (type 65001).
  - Parsed from YAML format.

##### 3.3.2.3 Loading Triggers

File entries are loaded when:

- File read operations are performed.
- File metadata queries require entry data.
- File-path associations are updated.
- Package validation requires entry inspection.

Path metadata is loaded when:

- Path metadata queries are performed.
- Tag inheritance operations require path hierarchy.
- File-path associations are updated.
- Path metadata management operations are called.

##### 3.3.2.4 Memory Management Benefits

- Large packages can be opened without loading all file entries into memory.
- Memory usage scales with actual usage, not package size.
- File handles remain open for efficient on-demand reading.
- Resources are released when Close() is called.

#### 3.3.3 State Management

The implementation does not track a separate "Created" state as a dedicated field.
Instead, state is derived from a small set of fields and invariants.

The only explicit lifecycle flag is `isOpen`.
Additional state is inferred from whether the package has an associated `FilePath`.
Additional state is inferred from whether a file handle exists.
Additional state is inferred from whether metadata caches are still present in memory.

##### 3.3.3.1 Lifecycle States

Lifecycle state definitions are expressed as predicates over fields.

- **New (unconfigured)**: Package returned by NewPackage(), not yet associated with a file.
  - Predicate: `FilePath == ""`.
  - Predicate: `isOpen == false`.
  - Notes: `Info` is initialized with default values but is not considered loaded from a package file.

- **Open (configured, in-memory)**: Package has been configured via NewPackageWithOptions().
  - Predicate: `FilePath != ""`.
  - Predicate: `isOpen == true`.
  - Predicate: `fileHandle == nil`.
  - Notes: No on-disk package is opened for reading.

- **Open (file-backed)**: Package has been opened via OpenPackage().
  - Predicate: `FilePath != ""`.
  - Predicate: `isOpen == true`.
  - Predicate: `fileHandle != nil`.
  - Notes: Header/index and required metadata caches are loaded as specified by OpenPackage eager loading requirements.

- **Closed (cached)**: Package has been closed via Close(), but cached metadata may remain.
  - Predicate: `isOpen == false`.
  - Predicate: `fileHandle == nil`.
  - Notes: Pure in-memory read operations MAY be allowed if required metadata caches remain in memory.
  - The `isOpen` flag being false is sufficient to indicate the package is closed.

- **Closed (cleaned)**: Package has been closed via CloseWithCleanup(), and caches are cleared.
  - Predicate: `isOpen == false`.
  - Predicate: `fileHandle == nil`.
  - Predicate: metadata caches required by pure in-memory reads have been cleared.
  - Notes: Pure in-memory read operations MUST fail in this state.
  - The `isOpen` flag being false is sufficient to indicate the package is closed.

##### 3.3.3.2 State Transitions

Valid state transitions:

- NewPackage() => New state.
- New state => NewPackageWithOptions() => Open (configured, in-memory) state.
- Open (configured, in-memory) state => Close() => Closed (cached) state.
- Open (configured, in-memory) state => CloseWithCleanup() => Closed (cleaned) state.
- OpenPackage() => Open (file-backed) state.
- Open (file-backed) state => Close() => Closed (cached) state.
- Open (file-backed) state => CloseWithCleanup() => Closed (cleaned) state.
- Any state => Close() => Closed (cached) state (idempotent).

##### 3.3.3.3 State-Dependent Operations

- GetInfo(): Available in Open states, and in Closed (cached) state if metadata remains in memory.
- GetMetadata(): Available in Open states, and in Closed (cached) state if metadata remains in memory.
- ListFiles(): Available in Open states, and in Closed (cached) state if metadata remains in memory.
- Validate(): Available only in Open (file-backed) state.
- ReadFile(): Available only in Open (file-backed) state.
- AddFile(), AddFileFromMemory(), RemoveFile(): Available in Open states.
- Close(): Available in any state (idempotent).
- CloseWithCleanup(): Available in any state.

##### 3.3.3.4 State Validation

The implementation validates state before operations:

- Closed (cleaned) state: Pure in-memory read operations MUST return validation errors.
- Closed (cached) state: I/O operations MUST return validation errors.
- Closed (cached) state: Pure in-memory read operations MAY succeed if metadata remains cached after Close().
- Open (file-backed) state: Required for I/O read operations and validation operations.
- Open states: Required for in-memory write operations (for example AddFile(), RemoveFile()).

#### 3.3.4 Memory Management

The implementation uses several strategies to manage memory efficiently.

##### 3.3.4.1 On-Demand Loading

File entries are not loaded into memory until needed.
This allows opening very large packages without excessive memory usage.

##### 3.3.4.2 Resource Cleanup

- File handles are closed when Close() is called.
- Memory allocated for cached metadata MAY remain after Close() to support in-memory read operations.
- CloseWithCleanup() clears in-memory caches (slices and maps) as part of cleanup.

##### 3.3.4.3 Large Package Handling

For packages that don't fit in memory:

- File entries are read from disk as needed.
- File data is streamed rather than loaded entirely.
- Index provides efficient random access to entries.
- File handles remain open for efficient I/O.

##### 3.3.4.4 Memory Allocation Patterns

- Slices are pre-allocated with appropriate capacity when size is known.
- Maps are initialized with make() when first needed.
- FileEntry data is allocated only when LoadData() is called.
- Temporary buffers are reused where possible.

#### 3.3.5 Data Relationships

The package implementation maintains several relationships between data structures.

##### 3.3.5.1 FileEntries ↔ PathMetadataEntries

- Association: FileEntry.PathMetadataEntries map links paths to PathMetadataEntry instances (direct references since both types are in the metadata package).
- Bidirectional: PathMetadataEntry.AssociatedFileEntries links back to FileEntry instances.
- Matching: Associations are established by matching path strings.
- Purpose: Enables per-path tag inheritance and filesystem properties.

##### 3.3.5.2 SpecialFiles Mapping

- Organization: Special files are organized by file type ID (65000-65535).
- Access: Direct lookup via file type ID.
- Types: Metadata files, signature files, manifest files, etc.
- Loading: Loaded during package opening.

##### 3.3.5.3 Header => Index => Entries

- Navigation: Header contains index location (IndexStart, IndexSize).
- Index: Contains entry count and FileID => Offset mappings.
- Entries: Located using index offsets, read on demand.
- Flow: Header => Index => IndexEntry => FileEntry.

##### 3.3.5.4 PackageInfo Aggregation

- Source: Aggregated from header, index, and entry data.
- Updates: Updated as package state changes.
- Caching: Cached to avoid repeated calculations.
- Access: Available via GetInfo() method.

#### 3.3.6 Thread Safety and Concurrency

The current implementation has specific concurrency characteristics.

##### 3.3.6.1 Current Limitations

- Package instances should not be shared across goroutines.
- No internal locking is performed.
- Concurrent operations on the same package instance are not safe.

##### 3.3.6.2 Safe Operations

- Multiple packages can be opened concurrently in different goroutines.
- Read operations on different packages are safe.
- Write operations on different packages are safe.

##### 3.3.6.3 Resource Locking

- No mutexes or locks are used in the current implementation.
- File I/O operations use standard Go file handles (not thread-safe).
- Package state fields are not protected by locks.

##### 3.3.6.4 Future Considerations

- Thread-safe operations would require adding mutexes.
- Read-write locks could enable concurrent reads.
- Per-operation locking could enable limited concurrency.

#### 3.3.7 Resource Lifecycle

The package implementation manages several types of resources throughout its lifecycle.

##### 3.3.7.1 File Handle Lifecycle

- Opened: During OpenPackage() operation, file handle is opened for reading.
- Active: File handle remains open while package is in Open state.
- Closed: During Close() operation, file handle is closed.
- Cleanup: File handle is set to nil after closing.

##### 3.3.7.2 Memory Allocation Lifecycle

- Initialization: Slices and maps are allocated during package creation or opening.
- Growth: Slices grow as entries are loaded on demand.
- Cleanup: Memory is released when package is closed.
- Garbage Collection: Go runtime handles final cleanup.

##### 3.3.7.3 State Flag Lifecycle

- isOpen: Set to true during OpenPackage(), false during Close().
- closed: Set to true during Close(), prevents further operations.
- Validation: Flags are checked before operations to ensure valid state.

##### 3.3.7.4 Error Recovery

- File handle errors: File handle is closed, package transitions to Closed state.
- Memory errors: Allocated resources are released, error is returned.
- State errors: Operations return validation errors, package state is preserved.

## 4. Package Format Constants

```go
const (
    // NVPKMagic is the magic number for .nvpk files
    NVPKMagic = 0x4E56504B // "NVPK" in hex

    // NVPKVersion is the current version of the .nvpk format
    NVPKVersion = 1

    // HeaderSize is the fixed size of the package header
    // See: Package File Format - Package Header for authoritative definition
    HeaderSize int64 = 112
)
```

These constants define fundamental values for the NovusPack format.

### 4.1 Format Constants

- `NVPKMagic`: Package identifier (0x4E56504B "NVPK")
- `NVPKVersion`: Current format version (1)
- `HeaderSize`: Fixed header size in bytes (see [Package File Format - Package Header](package_file_format.md#2-package-header))

**Usage**: Validate package header magic number and version before processing.

## 5. Package Lifecycle Operations

The NovusPack system follows a simple lifecycle pattern:

1. **Create** - Create a new package (in-memory)
2. **Open** - Open an existing package
3. **Operations** - Perform various operations (add files, metadata, etc.)
4. **Write** - Write the package to disk
5. **Close** - Close the package and release resources

### 5.1 Package Lifecycle - Always Use Defer for Cleanup

Always use defer statements to ensure resources are properly cleaned up, even when errors occur.
This prevents resource leaks and ensures consistent cleanup behavior.

### 5.2 Package Lifecycle - Check Package State Before Operations

Always verify that a package is in the correct state before performing operations.
Check if the package is open, not read-only, and in a valid state for the intended operation.

### 5.3. Package Lifecycle - Use Appropriate Context Timeouts

Use appropriate context timeouts for long-running operations to prevent indefinite blocking.
Set timeouts based on the expected operation duration and handle timeout errors gracefully.

## 6. NewPackage Function

```go
// NewPackage creates a new empty package
func NewPackage() (Package, error)
```

This function creates a new, empty NovusPack package in memory with default values.
The package exists only in memory until written to disk using one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`).

Returns a new `Package` instance with:

- Default header values (magic number, version, timestamps)
- Empty file index
- Empty package comment
- Closed state set to false

### 6.1 NewPackage Behavior

- Creates package structure in memory only (no file I/O operations performed)
- Initializes package with standard NovusPack header
- Sets creation timestamp to current time
- Prepares package for file operations
- Package must be written to disk using one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`) before it can be persisted

### 6.2 NewPackage Example Usage

```go
package, err := NewPackage()
if err != nil {
    return err
}
defer package.Close()
```

## 7. NewPackageWithOptions Function

```go
// NewPackageWithOptions creates a new package with specified configuration options
// Returns *PackageError on failure
func NewPackageWithOptions(ctx context.Context, options CreateOptions) (Package, error)
```

This function creates a new NovusPack package in memory with the specified configuration options.
**This function does not write to disk** - it only creates and configures the package structure in memory.
The package file is only written to disk when one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`) is called.

**Path Validation**: If `Path` is provided in options, this function validates that the provided path is valid and the target directory is writable, even though it doesn't write to disk. This ensures early detection of path-related issues.****

### 7.1 NewPackageWithOptions Parameters

- `ctx`: Context for cancellation and timeout handling
- `options`: Package configuration options

### 7.2 NewPackageWithOptions Behavior

- Creates package structure in memory (same as `NewPackage`)
- Initializes package with standard NovusPack header
- Sets creation timestamp to current time
- If `Path` is provided:
  - Validates that the provided path is valid and well-formed
  - Validates that the target directory exists and is writable (fails if directory does not exist)
  - Stores the target path for later writing operations
- Applies provided options:
  - Sets package comment if `Comment` is provided
  - Sets VendorID if provided
  - Sets AppID if provided
  - Stores file permissions for later use during Write operations
  - Enables package-level compression if `CompressPackage` is true
  - Sets session base path if `SessionBase` is provided (uses `SetSessionBase` internally)
- The package remains in an unsigned state until written (compressed if `CompressPackage` is true)
- No file I/O operations are performed on the target file - package remains in memory

**Note**: While the target file is not created during `NewPackageWithOptions`, if a path is provided, the path and directory are validated to ensure they exist and are writable. This enables early error detection before file operations begin. The parent directory must already exist - `NewPackageWithOptions` will not create missing parent directories.

### 7.3 NewPackageWithOptions Error Conditions

- **Validation Errors** (see [`ErrTypeValidation`](#201-error-types-used)):
  - Invalid or malformed file path (if Path is provided)
  - Target directory does not exist (parent directories are not created)
  - Target directory is not writable
  - Insufficient permissions to create file in target directory
- **Security Errors**: Insufficient permissions to access target directory (see [`ErrTypeSecurity`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded
- **Additional Validation Errors**: Invalid option values (e.g., invalid VendorID format)

### 7.4 NewPackageWithOptions Example Usage

```go
// Create package with options (still in memory)
options := CreateOptions{
    Path: Option.Some("/path/to/game-package.nvpk"),
    Comment: Option.Some("My Game Package"),
    VendorID: Option.Some(uint32(0x00000001)), // Steam
    AppID: Option.Some(uint64(0x00000000000002DA)), // CS:GO
    Permissions: Option.Some(os.FileMode(0644)),
    CompressPackage: true,
    SessionBase: Option.Some("/base/path"),
}

package, err := NewPackageWithOptions(ctx, options)
if err != nil {
    return err
}
defer package.Close()

// Add files, metadata, etc...
// ...

// Write to disk
err = package.Write(ctx)
if err != nil {
    return err
}
```

### 7.5 NewPackageWithOptions Without Path

If `Path` is not provided in options, the package is created without a target path.
The path can be set later using `SetTargetPath`:

```go
// Example: Create package without path
options := CreateOptions{
    Comment: Option.Some("My Package"),
    CompressPackage: true,
}

package, err := NewPackageWithOptions(ctx, options)
if err != nil {
    return err
}

// Set path later
err = package.SetTargetPath(ctx, "/path/to/package.nvpk")
if err != nil {
    return err
}
```

### 7.6 CreateOptions Structure

```go
// CreateOptions represents options for creating a package.
// CreateOptions allows configuring package creation with metadata,
// comments, and identifiers.
type CreateOptions struct {
    Path          Option[string]    // Target file system path where package will be written
    Comment       Option[string]    // Initial package comment
    VendorID      Option[uint32]    // Vendor identifier
    AppID         Option[uint64]    // Application identifier
    Permissions   Option[os.FileMode] // File permissions (default: 0644)
    CompressPackage bool            // Enable package-level compression (default: false)
    SessionBase   Option[string]    // Session base path for file operations
}
```

## 8. Package.SetTargetPath Method

```go
// SetTargetPath changes the package's target write path
// Returns *PackageError on failure
func (p *Package) SetTargetPath(ctx context.Context, path string) error
```

This function changes the target path for an existing package that will be used when calling `Write`, `SafeWrite`, or `FastWrite`.
This is useful when you want to write an existing package (either newly created or opened from disk) to a different location.

**Early Validation**: This function validates the path and target directory immediately (requiring minimal I/O), enabling early error detection before write operations begin.
Consistent with `Create`, this ensures path-related issues are caught early.

**Path Validation**: This function validates that the provided path is valid and the target directory is writable, even though it doesn't write to disk.
This validation requires minimal filesystem I/O to check directory existence and permissions.

**Signature Clearing**: If the package is signed and the new path differs from the current path, this function MUST clear all signature information from the in-memory package.
This is required because signed packages are immutable and writing to a new location creates a new, unsigned package.

**Important**: Signature clearing only occurs when the new path differs from the current path.
If `SetTargetPath` is called with the same path as the current path, signatures are NOT cleared.

See [Package Writing API - Writing Signed Package Content to New Path](api_writing.md#43-writing-signed-package-content-to-new-path) for complete signature clearing behavior.

### 8.1 Package.SetTargetPath Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: New file system path where the package will be written

### 8.2 Package.SetTargetPath Behavior

- Validates that the provided path is valid and well-formed
- Validates that the target directory exists and is writable (requires minimal filesystem I/O)
- If the new path differs from the current path and the package is signed, clears signature information (see [Package Writing API - Writing Signed Package Content to New Path](api_writing.md#43-writing-signed-package-content-to-new-path))
- Updates the package's internal target path
- Does not create or modify files (validation only)

### 8.3 Package.SetTargetPath Method Error Conditions

- **Validation Errors** (see [`ErrTypeValidation`](#201-error-types-used)):
  - Invalid or malformed file path
  - Target directory does not exist
  - Target directory is not writable
  - Insufficient permissions to create file in target directory
- **Security Errors**: Insufficient permissions to access target directory (see [`ErrTypeSecurity`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded

### 8.4 Package.SetTargetPath Example Usage

```go
// Open existing package
pkg, err := OpenPackage(ctx, "/path/to/existing.nvpk")
if err != nil {
    return err
}

// Make some modifications
err = pkg.AddFile(ctx, "/path/to/newfile.txt", nil)
if err != nil {
    return err
}

// Change target path to write to a new location (validates directory early)
err = pkg.SetTargetPath(ctx, "/path/to/modified-package.nvpk")
if err != nil {
    return err
}

// Write to the new location
err = pkg.SafeWrite(ctx, false)
if err != nil {
    return err
}
```

### 8.5. Package.SetTargetPath Vs NewPackageWithOptions

- `NewPackageWithOptions`: Used for initial package creation with configuration options, including optional path
- `SetTargetPath`: Used to change the write path on an existing package (created or opened)
- Both validate the target path and directory
- Both clear signatures when the path changes on a signed package
- `NewPackageWithOptions` creates and configures a new package; `SetTargetPath` only changes the path on an existing package

## 9. Package Configuration

This section describes package configuration options and settings.

### 9.1 PackageConfig Structure

```go
// Package configuration
type PackageConfig struct {
    // DefaultPathHandling specifies default behavior for multiple paths
    // Default: PathHandlingHardLinks (backward compatible)
    DefaultPathHandling PathHandling

    // AutoConvertToSymlinks enables automatic conversion of duplicate paths to symlinks
    // Default: false (backward compatible)
    AutoConvertToSymlinks bool
}
```

Provides package-level configuration for path handling behavior during file addition operations.

#### 9.1.1 PackageConfig Fields

- `DefaultPathHandling`: Specifies default behavior for multiple paths pointing to the same content
  - `PathHandlingHardLinks` (1): Store multiple paths as hard links (backward compatible default)
  - `PathHandlingSymlinks` (2): Convert additional paths to symlinks
  - `PathHandlingPreserve` (3): Preserve original filesystem behavior
  - Default: `PathHandlingHardLinks` (maintains backward compatibility)

- `AutoConvertToSymlinks`: Enables automatic conversion of duplicate paths to symlinks during deduplication
  - When `true`, deduplication creates symlinks instead of adding paths to existing FileEntry
  - Default: `false` (backward compatible)

### 9.2 PackageConfig Backward Compatibility

Default values ensure backward compatibility:

- `DefaultPathHandling = PathHandlingHardLinks`: Maintains existing behavior
- `AutoConvertToSymlinks = false`: No automatic conversion unless explicitly enabled
- Existing packages continue to work without changes

### 9.3 PathHandling Type

```go
// PathHandling specifies how to handle multiple paths pointing to the same content
type PathHandling uint8

const (
    PathHandlingDefault     PathHandling = 0 // Use package default
    PathHandlingHardLinks   PathHandling = 1 // Store multiple paths as hard links (current behavior)
    PathHandlingSymlinks    PathHandling = 2 // Convert additional paths to symlinks
    PathHandlingPreserve    PathHandling = 3 // Preserve original filesystem behavior (detect and respect symlinks/hardlinks)
)
```

## 10. OpenPackage Function

```go
// OpenPackage opens an existing package from the specified path.
// It validates the on-disk package structure during open.
// Returns *PackageError on failure
func OpenPackage(ctx context.Context, path string) (Package, error)
```

This function opens an existing NovusPack package file for reading.
This function validates the package header, index, and required invariants before returning a Package instance.

### 10.1 OpenPackage Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: File system path to the existing package file

### 10.2 OpenPackage Behavior

- Opens the package file for reading.
- Reads and validates the package header.
- Determines whether the package is compressed (header flags bits 15-8 != 0).
- If the package is compressed:
  - Reads the metadata index located immediately after the header.
  - Uses the metadata index to locate the compressed FileEntry metadata blocks.
  - Decompresses and reads all FileEntry metadata.
  - Locates, decompresses, and reads the compressed FileEntry index.
  - Reads the package comment, if present.
  - Detects signature presence for immutability enforcement, if present.
- If the package is not compressed:
  - Locates and reads all FileEntry metadata.
  - Locates and reads the FileEntry index.
  - Reads the package comment, if present.
  - Detects signature presence for immutability enforcement, if present.
- Validates the file index structure and offsets.
- Prepares the package for read operations.

### 10.3 OpenPackage Method Error Conditions

- **Validation Errors**: Invalid format, invalid header, invalid index structure, or violated invariants (see [`ErrTypeValidation`](#201-error-types-used)).
- **Unsupported Errors**: Package version not supported (see [`ErrTypeUnsupported`](#201-error-types-used)).
- **Compression Errors**: Failed to decompress compressed FileEntry metadata or the compressed FileEntry index during open.
- **Security Errors**: Insufficient permissions to open file (see [`ErrTypeSecurity`](#201-error-types-used)).
- **I/O Errors**: File system errors during opening.
- **Context Errors**: Context cancellation or timeout exceeded.

#### 10.3.1 OpenPackage Example Usage

```go
package, err := OpenPackage(ctx, "/path/to/existing-package.nvpk")
if err != nil {
    return err
}
defer package.Close()
```

## 11. Opening Packages As Read-Only

This section outlines the handling of read only packages.

### 11.1 Read-Only Enforcement Mechanism

The read-only mode must be enforced without duplicating OpenPackage parsing or validation logic.

In the API, OpenPackageReadOnly should call OpenPackage and then return a wrapper type that implements the Package interface.

The wrapper must delegate read operations to the underlying package and must reject all mutating operations.

The wrapper must return structured errors for rejected operations.

The wrapper should prevent callers from type-asserting the returned Package to the writable implementation type.

This can be achieved by returning a distinct wrapper type as the dynamic type behind the Package interface.

#### 11.1.1 Mutating Methods That Must Be Rejected

The wrapper must reject all methods that mutate package state in memory or write to disk.

This includes all PackageWriter methods, all state-changing metadata setters, and lifecycle methods that change the target path or package configuration for writing.

At minimum, the wrapper must reject Create, SetTargetPath, Defragment, AddFile, AddFileFromMemory, AddFilePattern, AddDirectory, RemoveFile, RemoveFilePattern, Write, SafeWrite, FastWrite, SetComment, ClearComment, SetAppID, ClearAppID, SetVendorID, ClearVendorID, SetPackageIdentity, and ClearPackageIdentity.

### 11.2 OpenPackageReadOnly Function

```go
// OpenPackageReadOnly opens a package in a read-only mode.
// It validates the on-disk package structure during open.
// The returned package must reject any attempt to mutate state or write to disk.
// Returns *PackageError on failure.
func OpenPackageReadOnly(ctx context.Context, path string) (Package, error)
```

This function opens an existing NovusPack package file for reading.
This function must enforce immutability for the returned package.

#### 11.2.1 OpenPackageReadOnly Behavior

- Opens the package file for reading.
- Reads and validates the package header.
- Determines whether the package is compressed (header flags bits 15-8 != 0).
- If the package is compressed:
  - Reads the metadata index located immediately after the header.
  - Uses the metadata index to locate the compressed FileEntry metadata blocks.
  - Decompresses and reads all FileEntry metadata.
  - Locates, decompresses, and reads the compressed FileEntry index.
  - Reads the package comment, if present.
  - Detects signature presence for immutability enforcement, if present.
- If the package is not compressed:
  - Locates and reads all FileEntry metadata.
  - Locates and reads the FileEntry index.
  - Reads the package comment, if present.
  - Detects signature presence for immutability enforcement, if present.
- Validates the file index structure and offsets.
- Returns a Package wrapper that enforces read-only behavior.
- Rejects in-memory mutation operations with a structured error.
- Rejects write operations to disk with a structured error.

#### 11.2.2 OpenPackageReadOnly Method Error Conditions

- All errors from `OpenPackage`.
- **Security Errors**: A write or mutation operation is attempted on a read-only package (see [`ErrTypeSecurity`](#201-error-types-used)).

### 11.3 readOnlyPackage Struct

```go
// readOnlyPackage is a wrapper that enforces read-only behavior for a Package.
//
// This wrapper must be the dynamic type stored behind the Package interface returned by OpenPackageReadOnly.
//
// This prevents callers from type-asserting to the writable implementation type.
type readOnlyPackage struct {
    inner Package
}

var _ Package = (*readOnlyPackage)(nil)
```

The `readOnlyPackage` type wraps the inner Package and implements the Package interface:

- **Read operations** (ReadFile, ListFiles, GetMetadata, GetInfo, Validate, Close, IsOpen, GetComment, HasComment, GetAppID, HasAppID, GetVendorID, HasVendorID, GetPackageIdentity) delegate directly to the inner Package.
- **Mutating operations** (Create, Defragment, AddFile, AddFileFromMemory, RemoveFile, Write, SafeWrite, FastWrite, SetComment, ClearComment, SetAppID, ClearAppID, SetVendorID, ClearVendorID, SetPackageIdentity, ClearPackageIdentity) return a read-only error via the [`readOnlyError` helper method](#114-readonlypackagereadonlyerror-method).

### 11.4 readOnlyPackage.readOnlyError Method

```go
// readOnlyError creates a structured security error for read-only enforcement.
// This helper method is used by all mutating operations to return consistent errors.
func (p *readOnlyPackage) readOnlyError(operation string) error {
    return pkgerrors.NewPackageError(pkgerrors.ErrTypeSecurity, "package is read-only", nil, ReadOnlyErrorContext{
        Operation: operation,
    })
}
```

### 11.5 ReadOnlyErrorContext Structure

```go
// ReadOnlyErrorContext provides typed context for read-only enforcement errors.
type ReadOnlyErrorContext struct {
    Operation string
}
```

### 11.6. ReadOnlyPackage Implementation Methods

The following implementation shows how `OpenPackageReadOnly` reuses `OpenPackage` and wraps the returned Package to enforce read-only behavior.
See [OpenPackageReadOnly](#112-openpackagereadonly-function) for the function signature.

The implementation pattern:

```go
// Implementation body for OpenPackageReadOnly
pkg, err := OpenPackage(ctx, path)
if err != nil {
    return nil, err
}

return &readOnlyPackage{inner: pkg}, nil
```

## 12. OpenBrokenPackage Function

```go
// OpenBrokenPackage opens a package that may be invalid or partially corrupted.
// This function is intended for repair workflows.
// TODO: Define repair APIs and the minimum guarantees of the returned package.
// Returns *PackageError on failure.
func OpenBrokenPackage(ctx context.Context, path string) (Package, error)
```

This function opens a package that may be invalid or partially corrupted.
This function is intended to support repair workflows and forensic inspection.
This function is not required to enforce the same validation guarantees as OpenPackage.
This function must return structured errors when I/O fails or when the file is unreadable.
This function should expose enough internal state to enable repair operations.

TODO: Specify repair operations and the allowed state transitions.

## 13. Package.Close Method

```go
// Close closes the package and releases resources
// Returns *PackageError on failure
func (p *Package) Close() error
```

This function closes the package file and releases all associated resources.

### 13.1 Package.Close Behavior

- Closes the package file handle
- Releases memory buffers and caches
- Clears package state and metadata
- Performs any necessary cleanup operations
- Does not modify the package file (use Write methods to save changes)

### 13.2 Package.Close Method Error Conditions

- **I/O Errors**: File system errors during closing (see [`ErrTypeIO`](#201-error-types-used))
- **Validation Errors**: Package is not currently open (see [`ErrTypeValidation`](#201-error-types-used))

### 13.3 Package.Close Example Usage

```go
// Example:
err := package.Close()
if err != nil {
    return err
}
```

## 14. Package.CloseWithCleanup Method

```go
// CloseWithCleanup closes the package and performs cleanup operations
// Returns *PackageError on failure
func (p *Package) CloseWithCleanup(ctx context.Context) error
```

This function closes the package and performs additional cleanup operations.

### 14.1 Package.CloseWithCleanup Behavior

- Closes the package file
- Performs cleanup operations (defragmentation, optimization)
- Releases all resources, including in-memory metadata caches
- May take longer than standard Close due to cleanup operations

## 15. Package.Validate Method

```go
// Validate validates package format, structure, and integrity
// Returns *PackageError on failure
func (p *Package) Validate(ctx context.Context) error
```

This function performs comprehensive validation of the package format, structure, and integrity.

### 15.1 Package.Validate Behavior

- Validates package header format and version
- Checks FileEntry structure and consistency
- Verifies data section integrity and checksums
- Detects signature presence for immutability enforcement, but does not validate signature contents in v1
- Ensures package follows NovusPack specifications
- Returns detailed error information for any issues found

### 15.2 Package.Validate Method Error Conditions

- **Validation Errors**: Package not open, invalid format, validation failed (see [`ErrTypeValidation`](#201-error-types-used))
- **Corruption Errors**: Checksum mismatches (see [`ErrTypeCorruption`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded (see [`ErrTypeContext`](#201-error-types-used))

### 15.3 Package.Validate Example Usage

```go
// Example
err := package.Validate(ctx)
if err != nil {
    return err
}
```

## 16. Package.Defragment Method

```go
// Defragment optimizes package structure and removes unused space
// Returns *PackageError on failure
func (p *Package) Defragment(ctx context.Context) error
```

This function optimizes package structure by removing unused space and reorganizing data for better performance.

### 16.1 Package.Defragment Behavior

- Removes unused space from deleted files
- Reorganizes file entries for optimal access
- Compacts data sections to reduce file size
- Updates internal indexes and references
- Preserves all package metadata and signatures
- May take significant time for large packages

### 16.2 Package.Defragment Error Conditions

- **Validation Errors**: Package not open, read-only mode (see [`ErrTypeValidation`](#201-error-types-used))
- **I/O Errors**: File system errors during defragmentation (see [`ErrTypeIO`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded (see [`ErrTypeContext`](#201-error-types-used))

### 16.3 Package.Defragment Example Usage

```go
// Example:
err := package.Defragment(ctx)
if err != nil {
    return err
}
```

## 17. Package.GetInfo Method

```go
// GetInfo gets basic package information
func (p *Package) GetInfo() (*PackageInfo, error)
```

This function retrieves comprehensive information about the current package.

Returns a `PackageInfo` structure containing:

- Basic package information (file count, sizes)
- Package identity (VendorID, AppID)
- Package comment and metadata
- Digital signature information
- Security and compression status
- Timestamps and feature flags

### 17.1 Package.GetInfo Error Conditions

- **Validation Errors**: Package not currently open (see [`ErrTypeValidation`](#201-error-types-used))

### 17.2 Package.GetInfo Example Usage

```go
// Get comprehensive package information
info, err := package.GetInfo()
if err != nil {
    return err
}
fmt.Printf("Package has %d files\n", info.FileCount)
fmt.Printf("Package version: %d\n", info.Version)
```

## 18. Header Inspection

These are low-level functions for header-only inspection without opening the full package.

### 18.1 Header Inspection Use Cases

- Validate .nvpk file format without loading package data
- Inspect package metadata before deciding to open
- Debugging corrupted or partially readable packages
- Stream processing where only header information is needed
- Quick validation of package files without full I/O overhead

### 18.2. ReadHeader Vs ReadHeaderFromPath

- `ReadHeader`: Use when you have an existing `io.Reader` or need fine-grained control over file operations
- `ReadHeaderFromPath`: Use when you want a simple, one-line header read from a file path with automatic file management

### 18.3 ReadHeader Function

```go
// ReadHeader reads the package header from a reader
func ReadHeader(ctx context.Context, reader io.Reader) (*Header, error)
```

### 18.4 ReadHeaderFromPath Function

```go
// ReadHeaderFromPath reads the package header from a file path
func ReadHeaderFromPath(ctx context.Context, path string) (*PackageHeader, error)
```

Reads the package header from a file path.
This is a convenience function that opens the file, reads the header, and closes the file automatically.

For more control over the file handle or to read from other sources, use `ReadHeader` with an `io.Reader`.

#### 18.4.1 ReadHeaderFromPath Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: File system path to the package file

Returns `*PackageHeader` and error.

#### 18.4.2 ReadHeaderFromPath Error Conditions

- **Validation Errors**: Invalid package header format, invalid file path (see [`ErrTypeValidation`](#201-error-types-used))
- **Unsupported Errors**: Unsupported package version (see [`ErrTypeUnsupported`](#201-error-types-used))
- **I/O Errors**: File not found, permission denied, file system errors (see [`ErrTypeIO`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded (see [`ErrTypeContext`](#201-error-types-used))

#### 18.4.3 ReadHeaderFromPath Example Usage

```go
// Example:
header, err := ReadHeaderFromPath(ctx, "/path/to/package.nvpk")
if err != nil {
    return err
}
fmt.Printf("Format Version: %d\n", header.FormatVersion)
fmt.Printf("Magic: 0x%08X\n", header.Magic)
```

### 18.5 Package.ReadHeader Method

Reads the package header from an `io.Reader`.
This function is useful when you already have an open file handle or need to read from a stream.

#### 18.5.1 Package.ReadHeader Parameters

- `ctx`: Context for cancellation and timeout handling
- `reader`: Input stream to read header from

Returns `*PackageHeader` and error.

#### 18.5.2 Package.ReadHeader Error Conditions

- **Validation Errors**: Invalid package header format (see [`ErrTypeValidation`](#201-error-types-used))
- **Unsupported Errors**: Unsupported package version (see [`ErrTypeUnsupported`](#201-error-types-used))
- **Context Errors**: Context cancellation or timeout exceeded (see [`ErrTypeContext`](#201-error-types-used))

##### 18.5.2.1 Package.ReadHeader Example Usage

```go
// Example:
file, err := os.Open("/path/to/package.nvpk")
if err != nil {
    return err
}
defer file.Close()

header, err := ReadHeader(ctx, file)
if err != nil {
    return err
}
fmt.Printf("Format Version: %d\n", header.FormatVersion)
```

### 18.6 Package.IsOpen Method

```go
// IsOpen checks if the package is currently open
func (p *Package) IsOpen() bool
```

### 18.7 Package.IsReadOnly Method

```go
// IsReadOnly checks if the package is in read-only mode
func (p *Package) IsReadOnly() bool
```

### 18.8 Package.GetPath Method

```go
// GetPath returns the current package file path
func (p *Package) GetPath() string
```

These functions provide information about the current package state.

## 19. Package Session Base Management

The `Package.sessionBase` property controls what the package considers to be the "base" path for converting between absolute filesystem paths and stored package paths.

The `sessionBase` is a runtime-only property that is used for bidirectional path conversion:

- **File Addition (Construction):** Converts absolute filesystem paths to stored package paths
- **File Extraction:** Converts stored package paths to absolute filesystem paths

The session base is a runtime-only property that persists during package operations and is used to automatically derive paths in both directions.

### 19.1 Session Base for File Addition

When an absolute filesystem path is provided to file addition operations ([`AddFile`](api_file_mgmt_addition.md#21-packageaddfile-method), [`AddFilePattern`](api_file_mgmt_addition.md#24-packageaddfilepattern-method), or [`AddDirectory`](api_file_mgmt_addition.md#25-packageadddirectory-method)) without an explicit `BasePath` in `AddFileOptions`, the session base is automatically established from that first absolute path.

Once established, the session base persists for all subsequent file operations within the same package construction session, ensuring consistent path derivation across multiple file additions.

For complete details on how session base affects path derivation during file addition, see [File Addition API - Session Base](api_file_mgmt_addition.md#264-session-base---package-level-automatic-basepath) and [AddFileOptions: Path Determination](api_file_mgmt_addition.md#26-addfileoptions-path-determination).

### 19.2 Session Base for File Extraction

When extracting files using [`ExtractPath`](api_file_mgmt_extraction.md#1-extractpath-package-method), the session base is used as the default extraction root directory.

The session base can be set explicitly using [`SetSessionBase`](#194-packagesetsessionbase-method) before extraction operations, or via the `SessionBase` option in [`ExtractPathOptions`](api_file_mgmt_extraction.md#2-extractpathoptions-struct).

If the session base is not set and is required for default-relative extraction, extraction operations will fail with `ErrTypeValidation`.

For complete details on how session base affects extraction destination resolution, see [File Extraction API - Destination Resolution](api_file_mgmt_extraction.md#151-extractpath-destination-resolution).

### 19.3 Session Base Lifecycle

The session base can be set explicitly using [`SetSessionBase`](#194-packagesetsessionbase-method) before any file operations, or via the `SessionBase` option in [`NewPackageWithOptions`](#7-newpackagewithoptions-function).

The session base is runtime-only and is not persisted to disk.

### 19.4 Package.SetSessionBase Method

```go
// SetSessionBase explicitly sets the package-level session base path
// This method allows setting the session base before any file operations
// Returns *PackageError on failure (e.g., invalid path format)
func (p *Package) SetSessionBase(basePath string) error
```

Sets the package-level session base path explicitly.

This is useful when you want to control the base path before adding files, rather than letting it be established automatically from the first absolute path.

#### 19.4.1 Package.SetSessionBase Parameters

- `basePath`: The filesystem base directory to use for path derivation (must be an absolute path)

#### 19.4.2 Package.SetSessionBase Returns

- `error`: `*PackageError` with `ErrTypeValidation` if the path format is invalid

#### 19.4.3 Package.SetSessionBase Example Usage

```go
// Example: Setting session base for file addition
pkg := NewPackage()
err := pkg.SetSessionBase("/home/user/project")
if err != nil {
    return err
}

// Now all absolute paths will be relative to /home/user/project
entry, err := pkg.AddFile(ctx, "/home/user/project/src/main.go", nil)
// Stored as: /src/main.go

// Example: Setting session base for file extraction
pkg2, err := OpenPackage(ctx, "/path/to/package.nvp")
if err != nil {
    return err
}
err = pkg2.SetSessionBase("/tmp/extract")
if err != nil {
    return err
}

// Extract files will be written to /tmp/extract/...
err = pkg2.ExtractPath(ctx, "/src/main.go", false, nil)
// Extracted to: /tmp/extract/src/main.go
```

### 19.5 Package.GetSessionBase Method

```go
// GetSessionBase returns the current session base path
// Returns empty string if no session base has been established
func (p *Package) GetSessionBase() string
```

#### 19.5.1 Package.GetSessionBase Returns

Returns the current session base path for inspection or logging:

- `string`: Current session base path, or empty string if no session base is set

#### 19.5.2 Package.GetSessionBase Example

```go
base := pkg.GetSessionBase()
if base != "" {
    fmt.Printf("Session base: %s\n", base)
}
```

### 19.6 Package.ClearSessionBase Method

```go
// ClearSessionBase clears the package-level session base
// Subsequent absolute paths will establish a new session base
func (p *Package) ClearSessionBase()
```

Clears the current session base, allowing a new base to be established from the next absolute path.

This is useful when you want to change the base path strategy mid-construction.

#### 19.6.1 Package.ClearSessionBase Example

```go
// Example: Add files with one base
pkg.SetSessionBase("/home/user/project1")
pkg.AddFile(ctx, "/home/user/project1/file1.txt", nil)

// Example: Switch to different base
pkg.ClearSessionBase()
pkg.SetSessionBase("/home/user/project2")
pkg.AddFile(ctx, "/home/user/project2/file2.txt", nil)
```

### 19.7 Package.HasSessionBase Method

```go
// HasSessionBase returns true if a session base is currently set
func (p *Package) HasSessionBase() bool
```

Checks whether a session base is currently active.

#### 19.7.1 Package.HasSessionBase Returns

- `bool`: true if a session base is set, false otherwise

#### 19.7.2 Package.HasSessionBase Example

```go
if !pkg.HasSessionBase() {
    // Set explicit base before adding files
    pkg.SetSessionBase("/home/user/project")
}
```

## 20. Structured Error System

The NovusPack API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities.

**Usage**: Create structured errors with rich context for different error scenarios.

### 20.1 Error Types Used

The NovusPack Basic Operations API uses the following error types from the structured error system:

- `ErrTypeValidation`: Input validation errors, invalid parameters, format errors
- `ErrTypeIO`: I/O errors, file system operations, network issues
- `ErrTypeSecurity`: Security-related errors, access denied, authentication failures
- `ErrTypeUnsupported`: Unsupported features, versions, or operations
- `ErrTypeContext`: Context cancellation, timeout, and lifecycle errors
- `ErrTypeCorruption`: Data corruption, checksum failures, integrity violations

### 20.2 PackageErrorContext Structure

```go
// Define error context types
type PackageErrorContext struct {
    Path      string
    Operation string
}
```

### 20.3 SecurityErrorContext Structure

```go
// SecurityErrorContext provides typed context for security-related errors.
// This context structure is used with structured errors to provide additional
// diagnostic information for security operations.
type SecurityErrorContext struct {
    Path string
    User string
}
```

### 20.4 IOErrorContext Structure

```go
// IOErrorContext provides typed context for I/O-related errors.
// This context structure is used with structured errors to provide additional
// diagnostic information for file system and I/O operations.
type IOErrorContext struct {
    File   string
    Offset int64
}
```

### 20.5 Creating Errors with Context

```go
// Example: Create a validation error with typed context
err := NewPackageError(ErrTypeValidation, "package file not found", nil, PackageErrorContext{
    Path:      "/path/to/package.nvpk",
    Operation: "Open",
})

// Example: Wrap an existing error with typed context
err := WrapErrorWithContext(io.ErrUnexpectedEOF, ErrTypeIO, "unexpected end of file", IOErrorContext{
    File:   "package.nvpk",
    Offset: 1024,
})

// Example: Create a security error with typed context
err := NewPackageError(ErrTypeSecurity, "permission denied", nil, SecurityErrorContext{
    Path: "/path/to/package.nvpk",
    User: "anonymous",
})

// Example: Create an I/O error with typed context
err := NewPackageError(ErrTypeIO, "failed to read package file", nil, PackageErrorContext{
    Path:      "/path/to/package.nvpk",
    Operation: "Open",
})
```

### 20.6 Error Inspection

**Usage**: Check error types and handle them with appropriate logging and context extraction.

## 21. Error Handling Best Practices

Best practices for error handling.

### 21.1 Always Check for Errors

Always check for errors after calling package operations and handle them appropriately.
Never ignore error return values as they indicate critical failures that must be addressed.

### 21.2 Use Structured Errors for Better Debugging

Use the structured error system to provide rich context for debugging.
Wrap errors with additional context information to help identify the source of problems and provide better error messages to users.

### 21.3 Use Context for Cancellation

Use context timeouts and cancellation to prevent operations from hanging indefinitely.
Set appropriate timeouts for long-running operations and handle context cancellation gracefully.

### 21.4 Handle Different Error Types Appropriately

Handle different error types with appropriate responses.
Provide user-friendly messages for validation errors, log security errors, and implement retry logic for I/O errors.
Use the structured error system to determine the appropriate handling strategy.

### 21.5 Clean Up Resources

Always clean up resources properly using defer statements.
Ensure packages are closed even when errors occur, and handle cleanup errors gracefully to prevent resource leaks.

### 21.6 Wrap Errors with Context

Wrap errors with additional context information to provide better debugging information.
Include relevant details such as file paths, operation names, and parameter values in error messages.

## 22. Resource Management

This section describes best practices for managing resources in package operations.

### 22.1 Use Context for Resource Management

Use context for resource management and cancellation.
Pass context to long-running operations and handle context cancellation to ensure proper resource cleanup and operation termination.

### 22.2 Handle Cleanup Errors Gracefully

Handle cleanup errors gracefully by logging warnings rather than failing.
Use defer functions to ensure cleanup occurs even when errors happen, and log cleanup failures as warnings rather than errors.
