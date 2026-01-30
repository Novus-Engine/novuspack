# NovusPack Go API Definitions Index

- [0. Overview](#0-overview)
- [1. Package Interface Types](#1-package-interface-types)
  - [1.1 Package Lifecycle Methods](#11-package-lifecycle-methods)
  - [1.2 Package File Management Methods](#12-package-file-management-methods)
  - [1.3 Package Information and Queries Methods](#13-package-information-and-queries-methods)
  - [1.4 Package Metadata Methods](#14-package-metadata-methods)
  - [1.5 Package Compression Methods](#15-package-compression-methods)
  - [1.6 Package Path and Configuration Methods](#16-package-path-and-configuration-methods)
  - [1.7 Package Signature Management Methods](#17-package-signature-management-methods)
  - [1.8 Package Other Methods](#18-package-other-methods)
  - [1.9 Package Helper Functions](#19-package-helper-functions)
- [2. PackageReader Interface Types](#2-packagereader-interface-types)
  - [2.1 PackageReader Methods](#21-packagereader-methods)
    - [2.1.1 PackageReader Read Operations](#211-packagereader-read-operations)
    - [2.1.2 PackageReader Query Operations](#212-packagereader-query-operations)
    - [2.1.3 PackageReader Other Methods](#213-packagereader-other-methods)
  - [2.2 PackageReader Helper Functions](#22-packagereader-helper-functions)
- [3. PackageWriter Interface Types](#3-packagewriter-interface-types)
  - [3.1 PackageWriter Methods](#31-packagewriter-methods)
    - [3.1.1 PackageWriter Write Operations](#311-packagewriter-write-operations)
    - [3.1.2 PackageWriter Other Methods](#312-packagewriter-other-methods)
  - [3.2 PackageWriter Helper Functions](#32-packagewriter-helper-functions)
- [4. FileEntry Types](#4-fileentry-types)
  - [4.1 FileEntry Methods](#41-fileentry-methods)
    - [4.1.1 FileEntry Data Management Methods](#411-fileentry-data-management-methods)
    - [4.1.2 FileEntry Transformation Methods](#412-fileentry-transformation-methods)
  - [4.2 FileEntry Helper Functions](#42-fileentry-helper-functions)
- [5. Metadata Types](#5-metadata-types)
  - [5.1 Metadata Methods](#51-metadata-methods)
  - [5.2 Metadata Helper Functions](#52-metadata-helper-functions)
- [6. Compression Types](#6-compression-types)
  - [6.1 Compression Methods](#61-compression-methods)
  - [6.2 Compression Helper Functions](#62-compression-helper-functions)
- [7. Encryption and Security Types](#7-encryption-and-security-types)
  - [7.1 Encryption and Security Methods](#71-encryption-and-security-methods)
  - [7.2 Encryption and Security Helper Functions](#72-encryption-and-security-helper-functions)
- [8. Signature Types](#8-signature-types)
  - [8.1 Signature Methods](#81-signature-methods)
  - [8.2 Signature Helper Functions](#82-signature-helper-functions)
- [9. Streaming and Buffer Types](#9-streaming-and-buffer-types)
  - [9.1 Streaming and Buffer Methods](#91-streaming-and-buffer-methods)
  - [9.2 Streaming and Buffer Helper Functions](#92-streaming-and-buffer-helper-functions)
- [10. Deduplication Types](#10-deduplication-types)
  - [10.1 Deduplication Methods](#101-deduplication-methods)
  - [10.2 Deduplication Helper Functions](#102-deduplication-helper-functions)
- [11. FileType System Types](#11-filetype-system-types)
  - [11.1 FileType System Methods](#111-filetype-system-methods)
  - [11.2 FileType System Helper Functions](#112-filetype-system-helper-functions)
- [12. Generic Types](#12-generic-types)
  - [12.1 Generic Methods](#121-generic-methods)
  - [12.2 Generic Helper Functions](#122-generic-helper-functions)
- [13. Error Types](#13-error-types)
  - [13.1 Error Methods](#131-error-methods)
  - [13.2 Error Helper Functions](#132-error-helper-functions)
- [14. Other Types](#14-other-types)
  - [14.1 Other Type Methods](#141-other-type-methods)
- [15. General Helper Functions](#15-general-helper-functions)
  - [15.1 General Validation Functions](#151-general-validation-functions)
  - [15.2 General Utility Functions](#152-general-utility-functions)

## 0. Overview

This document provides a comprehensive index of all NovusPack API functions, types, and structures, with direct links to their detailed documentation.
Use this index to quickly locate specific API elements across the documentation.

## 1. Package Interface Types

- **`filePackage`** - [filePackage Struct](api_core.md#111-filepackage-struct)
  - filePackage is the concrete implementation of the Package interface.
    - filePackage is documented in the linked spec.
- **`Package`** - [Package](api_core.md#11-package-interface)
  - Package defines the main interface for NovusPack package operations.
  - Package combines PackageReader and PackageWriter interfaces, providing complete package lifecycle management including opening, closing, and defragmentation operations.
    - Package is documented in the linked spec.

### 1.1 Package Lifecycle Methods

- **`Package.Close`** - [Package.Close](api_basic_operations.md#13-packageclose-method)
  - Close closes the package and releases resources Returns *PackageError on failure.
- **`Package.Validate`** - [Package.Validate](api_basic_operations.md#15-packagevalidate-method)
  - Validate validates package format, structure, and integrity.
- **`Package.ValidateIntegrity`** - [Package.ValidateIntegrity](api_security.md#111-packagevalidateintegrity-method)
  - ValidateIntegrity validates package integrity (checksums and structural consistency).
- **`Package.Defragment`** - [Package.Defragment](api_basic_operations.md#16-packagedefragment-method)
  - Defragment optimizes the package layout and removes unused space.

### 1.2 Package File Management Methods

- **`Package.AddDirectory`** - [Package.AddDirectory](api_file_mgmt_addition.md#25-packageadddirectory-method)
  - AddDirectory recursively adds files from a filesystem directory into the package.
- **`Package.AddFile`** - [Package.AddFile](api_file_mgmt_addition.md#21-packageaddfile-method)
  - AddFile adds a file to the package.
- **`Package.AddFileFromMemory`** - [Package.AddFileFromMemory](api_file_mgmt_addition.md#22-packageaddfilefrommemory-method)
  - AddFileFromMemory adds a file to the package from in-memory data.
- **`Package.AddFileWithEncryption`** - [Package.AddFileWithEncryption](api_file_mgmt_addition.md#23-packageaddfilewithencryption-method)
  - AddFileWithEncryption adds a file to the package and configures encryption for the entry.
- **`Package.AddFilePattern`** - [Package.AddFilePattern](api_file_mgmt_addition.md#24-packageaddfilepattern-method)
  - AddFilePattern adds files matching a filesystem pattern into the package.
- **`Package.ExtractPath`** - [Package.ExtractPath](api_file_mgmt_extraction.md#12-packageextractpath-method)
  - ExtractPath extracts a file or directory subtree from the package to disk.
- **`Package.RemoveDirectory`** - [Package.RemoveDirectory](api_file_mgmt_removal.md#42-packageremovedirectory-method)
  - RemoveDirectory removes all files within a directory path from the package.
  - High-level counterpart to AddDirectory.
  - Returns *PackageError on failure.
- **`Package.RemoveFile`** - [Package.RemoveFile](api_file_mgmt_removal.md#22-packageremovefile-method)
  - RemoveFile removes a file from the package.
  - High-level counterpart to AddFile.
  - Returns *PackageError on failure.
- **`Package.RemoveFilePattern`** - [Package.RemoveFilePattern](api_file_mgmt_removal.md#32-packageremovefilepattern-method)
  - RemoveFilePattern removes files matching a pattern from the package.
  - High-level counterpart to AddFilePattern.
  - Returns *PackageError on failure.
- **`Package.UpdateFilePattern`** - [Package.UpdateFilePattern](api_file_mgmt_updates.md#12-packageupdatefilepattern-method)
  - UpdateFilePattern updates files matching a pattern using a source directory and options.
- **`Package.AddFilePath`** - [Package.AddFilePath](api_file_mgmt_updates.md#14-packageaddfilepath-method)
  - AddFilePath adds an additional stored path to an existing FileEntry.
- **`Package.RemoveFilePath`** - [Package.RemoveFilePath](api_file_mgmt_updates.md#15-packageremovefilepath-method)
  - RemoveFilePath removes a stored path from an existing FileEntry.
- **`Package.AddFileHash`** - [Package.AddFileHash](api_file_mgmt_updates.md#16-packageaddfilehash-method)
  - AddFileHash adds a hash entry to a FileEntry for integrity or deduplication.
- **`Package.SetSessionBase`** - [Package.SetSessionBase](api_basic_operations.md#194-packagesetsessionbase-method)
  - SetSessionBase explicitly sets the package-level session base path This method allows setting the session base before any file operations Returns *PackageError on failure (e.g., invalid path format).
- **`Package.ClearSessionBase`** - [Package.ClearSessionBase](api_basic_operations.md#196-packageclearsessionbase-method)
  - ClearSessionBase clears the current session base path.

### 1.3 Package Information and Queries Methods

- **`Package.FileExists`** - [Package.FileExists](api_file_mgmt_queries.md#111-packagefileexists-method)
  - FileExists checks if a file with the given path exists in the package.
- **`Package.FindEntriesByPathPatterns`** - [Package.FindEntriesByPathPatterns](api_file_mgmt_queries.md#331-packagefindentriesbypathpatterns-method)
  - FindEntriesByPathPatterns gets files matching patterns from the package.
- **`Package.FindEntriesByType`** - [Package.FindEntriesByType](api_file_mgmt_queries.md#321-packagefindentriesbytype-method)
  - FindEntriesByType finds all FileEntry objects of a specific type.
- **`Package.GetFileByChecksum`** - [Package.GetFileByChecksum](api_file_mgmt_queries.md#251-packagegetfilebychecksum-method)
  - GetFileByChecksum gets a FileEntry by CRC32 checksum Returns *PackageError if file not found.
- **`Package.GetFileByFileID`** - [Package.GetFileByFileID](api_file_mgmt_queries.md#231-packagegetfilebyfileid-method)
  - GetFileByFileID gets a FileEntry by its unique FileID Returns *PackageError if file not found.
- **`Package.GetFileByHash`** - [Package.GetFileByHash](api_file_mgmt_queries.md#241-packagegetfilebyhash-method)
  - GetFileByHash gets a FileEntry by content hash Returns *PackageError if file not found.
- **`Package.GetFileByOffset`** - [Package.GetFileByOffset](api_file_mgmt_queries.md#221-packagegetfilebyoffset-method)
  - GetFileByOffset gets a FileEntry by offset Returns *PackageError if file not found.
- **`Package.GetFileByPath`** - [Package.GetFileByPath](api_file_mgmt_queries.md#211-packagegetfilebypath-method)
  - GetFileByPath gets a FileEntry by path Returns *PackageError if file not found.
- **`Package.GetFileCount`** - [Package.GetFileCount](api_file_mgmt_queries.md#411-packagegetfilecount-method)
  - GetFileCount returns the total number of regular content files in the package Excludes special metadata files (types 65000-65535).
- **`Package.GetPath`** - [Package.GetPath](api_basic_operations.md#188-packagegetpath-method)
  - GetPath returns the current package file path.
- **`Package.GetSessionBase`** - [Package.GetSessionBase](api_basic_operations.md#195-packagegetsessionbase-method)
  - GetSessionBase returns the current session base path Returns empty string if no session base has been established.
- **`Package.IsOpen`** - [Package.IsOpen](api_basic_operations.md#186-packageisopen-method)
  - IsOpen checks if the package is currently open.
- **`Package.IsReadOnly`** - [Package.IsReadOnly](api_basic_operations.md#187-packageisreadonly-method)
  - IsReadOnly checks if the package is in read-only mode.
- **`Package.ListFiles`** - [Package.ListFiles](api_file_mgmt_queries.md#112-packagelistfiles-method)
  - ListFiles returns all file entries in the package.
- **`Package.GetSecurityStatus`** - [Package.GetSecurityStatus](api_security.md#112-packagegetsecuritystatus-method)
  - GetSecurityStatus returns the current security status of the package.
- **`Package.ListEncryptedFiles`** - [Package.ListEncryptedFiles](api_file_mgmt_queries.md#431-packagelistencryptedfiles-method)
  - ListEncryptedFiles returns encrypted file entries in the package.
- **`Package.FindExistingEntryByCRC32`** - [Package.FindExistingEntryByCRC32](api_deduplication.md#311-packagefindexistingentrybycrc32-method)
  - FindExistingEntryByCRC32 finds an existing entry by size and CRC32 (deduplication helper).
- **`Package.FindExistingEntryMultiLayer`** - [Package.FindExistingEntryMultiLayer](api_deduplication.md#312-packagefindexistingentrymultilayer-method)
  - FindExistingEntryMultiLayer finds an existing entry using multi-layer deduplication checks.
- **`Package.AddPathToExistingEntry`** - [Package.AddPathToExistingEntry](api_deduplication.md#313-packageaddpathtoexistingentry-method)
  - AddPathToExistingEntry adds an additional path to an existing entry as part of deduplication.
- **`Package.SetTargetPath`** - [Package.SetTargetPath](api_basic_operations.md#8-packagesettargetpath-method)
  - SetTargetPath changes the package's target write path Returns *PackageError on failure.
- **`Package.updateFilePathAssociations`** - [Package.updateFilePathAssociations](api_basic_operations.md#323-packageupdatefilepathassociations-method)
  - updateFilePathAssociations links files to their path metadata Returns *PackageError on failure.

### 1.4 Package Metadata Methods

- **`Package.AddDirectoryMetadata`** - [Package.AddDirectoryMetadata](api_metadata.md#8219-packageadddirectorymetadata-method)
  - AddDirectoryMetadata adds directory path metadata (metadata-only, does not add files) Returns *PackageError on failure.
- **`Package.AddIndexFile`** - [Package.AddIndexFile](api_metadata.md#531-packageaddindexfile-method)
  - AddIndexFile adds a package index file Returns *PackageError on failure.
- **`Package.AddManifestFile`** - [Package.AddManifestFile](api_metadata.md#521-packageaddmanifestfile-method)
  - AddManifestFile adds a package manifest file Returns *PackageError on failure.
- **`Package.AddMetadataFile`** - [Package.AddMetadataFile](api_metadata.md#511-packageaddmetadatafile-method)
  - AddMetadataFile adds a YAML metadata file to the package Returns *PackageError on failure.
- **`Package.AddMetadataOnlyFile`** - [Package.AddMetadataOnlyFile](api_metadata.md#642-packageaddmetadataonlyfile-method)
  - AddMetadataOnlyFile adds a special metadata file to a metadata-only package Returns *PackageError on failure.
- **`Package.AddPathMetadata`** - [Package.AddPathMetadata](api_metadata.md#8213-packageaddpathmetadata-method)
  - AddPathMetadata adds a new path metadata entry to the package Returns *PackageError on failure.
- **`Package.AddSignatureFile`** - [Package.AddSignatureFile](api_metadata.md#541-packageaddsignaturefile-method)
  - AddSignatureFile adds a digital signature file Returns *PackageError on failure.
- **`Package.AddSymlink`** - [Package.AddSymlink](api_metadata.md#8531-packageaddsymlink-method)
  - AddSymlink adds a symbolic link to the package.
  - Parameter: symlink: SymlinkEntry to add.
  - Return: Error if validation fails or symlink cannot be added.
  - Validation: Calls ValidateSymlinkPaths() to ensure paths are valid and within package root.
  - Validation: Verifies target exists as FileEntry or PathMetadataEntry directory.
  - Validation: Returns ErrTypeValidation, ErrTypeSecurity, or ErrTypeNotFound on validation failure Returns *PackageError on failure.
- **`Package.ClearAppID`** - [Package.ClearAppID](api_metadata.md#212-packageclearappid-method)
  - ClearAppID removes the package AppID (set to 0) Returns *PackageError on failure.
- **`Package.ClearComment`** - [Package.ClearComment](api_metadata.md#113-packageclearcomment-method)
  - ClearComment removes the package comment Returns *PackageError on failure.
- **`Package.ClearPackageIdentity`** - [Package.ClearPackageIdentity](api_metadata.md#412-packageclearpackageidentity-method)
  - ClearPackageIdentity clears both VendorID and AppID Returns *PackageError on failure.
- **`Package.ClearVendorID`** - [Package.ClearVendorID](api_metadata.md#312-packageclearvendorid-method)
  - ClearVendorID removes the package VendorID (set to 0) Returns *PackageError on failure.
- **`Package.ConvertPathsToSymlinks`** - [Package.ConvertPathsToSymlinks](api_file_mgmt_updates.md#1711-packageconvertpathstosymlinks-method)
  - ConvertPathsToSymlinks converts duplicate paths on a FileEntry to symlinks.
  - Parameter: ctx: Context for cancellation and timeout.
  - Parameter: entry: FileEntry with multiple paths (PathCount > 1).
  - Parameter: options: Path-to-symlink conversion options (primary path selection, metadata preservation).
  - Return: Updated FileEntry with single path.
  - Return: Slice of created SymlinkEntry objects.
- **`Package.ConvertSymlinksToHardLinks`** - [Package.ConvertSymlinksToHardLinks](api_file_mgmt_updates.md#1721-packageconvertsymlinkstohardlinks-method)
  - ConvertSymlinksToHardLinks converts symlinks back to hard links (reverse operation).
  - Parameter: ctx: Context for cancellation and timeout.
  - Parameter: symlinkEntry: SymlinkEntry to convert back to hard link.
  - Return: Updated FileEntry with additional path added.
  - Return: Error if conversion fails.
  - Behavior: Removes SymlinkEntry from package.
- **`Package.ConvertAllPathsToSymlinks`** - [Package.ConvertAllPathsToSymlinks](api_file_mgmt_updates.md#1712-packageconvertallpathstosymlinks-method)
  - ConvertAllPathsToSymlinks converts all eligible multi-path entries to use symlinks.
- **`Package.ConvertAllSymlinksToHardLinks`** - [Package.ConvertAllSymlinksToHardLinks](api_file_mgmt_updates.md#1722-packageconvertallsymlinkstohardlinks-method)
  - ConvertAllSymlinksToHardLinks converts all symlinks back to hard links.
- **`Package.GetMultiPathEntries`** - [Package.GetMultiPathEntries](api_file_mgmt_updates.md#1731-packagegetmultipathentries-method)
  - GetMultiPathEntries returns FileEntry objects with multiple stored paths.
- **`Package.GetMultiPathCount`** - [Package.GetMultiPathCount](api_file_mgmt_updates.md#1732-packagegetmultipathcount-method)
  - GetMultiPathCount returns the number of entries with multiple stored paths.
- **`Package.CreateSpecialMetadataFile`** - [Package.CreateSpecialMetadataFile](api_metadata.md#82114-packagecreatespecialmetadatafile-method)
  - CreateSpecialMetadataFile creates a special metadata FileEntry Returns *PackageError on failure.
- **`Package.FindEntriesByTag`** - [Package.FindEntriesByTag](api_file_mgmt_queries.md#311-packagefindentriesbytag-method)
  - FindEntriesByTag finds all FileEntry objects with a specific tag.
- **`Package.GetAppID`** - [Package.GetAppID](api_metadata.md#211-packagegetappid-method)
  - GetAppID retrieves the current package AppID.
- **`Package.GetAppIDInfo`** - [Package.GetAppIDInfo](api_metadata.md#214-packagegetappidinfo-method)
  - GetAppIDInfo gets detailed AppID information if available.
- **`Package.GetComment`** - [Package.GetComment](api_metadata.md#112-packagegetcomment-method)
  - GetComment retrieves the current package comment.
- **`Package.GetFilePathAssociation`** - [Package.GetFilePathAssociation](api_metadata.md#8511-packagegetfilepathassociation-method)
  - Package.GetFilePathAssociation File-path association query methods (Package-level) These methods work with path strings to find and return associated structs.
- **`Package.GetFilesInPath`** - [Package.GetFilesInPath](api_metadata.md#8512-packagegetfilesinpath-method)
  - GetFilesInPath returns all file entries within the specified path.
- **`Package.GetIndexFile`** - [Package.GetIndexFile](api_metadata.md#532-packagegetindexfile-method)
  - GetIndexFile retrieves the package index Returns *PackageError on failure.
- **`Package.GetManifestFile`** - [Package.GetManifestFile](api_metadata.md#522-packagegetmanifestfile-method)
  - GetManifestFile retrieves the package manifest Returns *PackageError on failure.
- **`Package.GetMetadataFile`** - [Package.GetMetadataFile](api_metadata.md#512-packagegetmetadatafile-method)
  - GetMetadataFile retrieves metadata from the special metadata file Returns *PackageError on failure.
- **`Package.GetMetadataIndexOffset`** - [Package.GetMetadataIndexOffset](api_package_compression.md#7292-packagegetmetadataindexoffset-method)
  - GetMetadataIndexOffset returns the offset to metadata index Returns fixed offset 112 bytes (PackageHeaderSize) when compression enabled Returns *PackageError if package is not compressed (no metadata index).
- **`Package.GetMetadataOnlyFiles`** - [Package.GetMetadataOnlyFiles](api_metadata.md#643-packagegetmetadataonlyfiles-method)
  - GetMetadataOnlyFiles returns all metadata files in the package Returns *PackageError on failure.
- **`Package.GetPackageIdentity`** - [Package.GetPackageIdentity](api_metadata.md#411-packagegetpackageidentity-method)
  - GetPackageIdentity gets both VendorID and AppID.
- **`Package.GetPackageInfo`** - [Package.GetPackageInfo](api_metadata.md#741-packagegetpackageinfo-method)
  - GetPackageInfo returns comprehensive package information.
- **`Package.GetPathConflicts`** - [Package.GetPathConflicts](api_metadata.md#82113-packagegetpathconflicts-method)
  - GetPathConflicts returns a list of paths with conflicting metadata Returns *PackageError on failure.
- **`Package.GetPathFiles`** - [Package.GetPathFiles](api_metadata.md#8513-packagegetpathfiles-method)
  - GetPathFiles returns all file entries associated with the specified path.
- **`Package.GetPathMetadata`** - [Package.GetPathMetadata](api_metadata.md#8211-packagegetpathmetadata-method)
  - GetPathMetadata retrieves all path metadata entries from the package Returns *PackageError on failure.
- **`Package.GetPathStats`** - [Package.GetPathStats](api_metadata.md#8522-packagegetpathstats-method)
  - GetPathStats returns statistics for all paths in the package.
- **`Package.GetPathTree`** - [Package.GetPathTree](api_metadata.md#8521-packagegetpathtree-method)
  - Package.GetPathTree Path hierarchy analysis.
- **`Package.GetSignatureFile`** - [Package.GetSignatureFile](api_metadata.md#542-packagegetsignaturefile-method)
  - GetSignatureFile retrieves the signature file Returns *PackageError on failure.
- **`Package.GetSpecialFileByType`** - [Package.GetSpecialFileByType](api_metadata.md#552-packagegetspecialfilebytype-method)
  - GetSpecialFileByType retrieves special file by type.
- **`Package.GetSpecialFiles`** - [Package.GetSpecialFiles](api_metadata.md#551-packagegetspecialfiles-method)
  - GetSpecialFiles returns all special files in the package.
- **`Package.GetSymlink`** - [Package.GetSymlink](api_metadata.md#8533-packagegetsymlink-method)
  - Package.GetSymlink Returns *PackageError on failure.
- **`Package.GetVendorID`** - [Package.GetVendorID](api_metadata.md#311-packagegetvendorid-method)
  - GetVendorID retrieves the current package VendorID.
- **`Package.GetVendorIDInfo`** - [Package.GetVendorIDInfo](api_metadata.md#314-packagegetvendoridinfo-method)
  - GetVendorIDInfo gets detailed VendorID information if available.
- **`Package.HasAppID`** - [Package.HasAppID](api_metadata.md#213-packagehasappid-method)
  - HasAppID checks if the package has an AppID (non-zero).
- **`Package.HasComment`** - [Package.HasComment](api_metadata.md#114-packagehascomment-method)
  - HasComment checks if the package has a comment.
- **`Package.HasIndexFile`** - [Package.HasIndexFile](api_metadata.md#535-packagehasindexfile-method)
  - HasIndexFile checks if package has an index file.
- **`Package.HasManifestFile`** - [Package.HasManifestFile](api_metadata.md#525-packagehasmanifestfile-method)
  - HasManifestFile checks if package has a manifest file.
- **`Package.HasMetadataFile`** - [Package.HasMetadataFile](api_metadata.md#515-packagehasmetadatafile-method)
  - HasMetadataFile checks if package has a metadata file.
- **`Package.HasMetadataIndex`** - [Package.HasMetadataIndex](api_package_compression.md#7291-packagehasmetadataindex-method)
  - HasMetadataIndex checks if package has metadata index (compression enabled) Returns true if header flags bits 15-8 != 0.
- **`Package.HasSignatureFile`** - [Package.HasSignatureFile](api_metadata.md#545-packagehassignaturefile-method)
  - HasSignatureFile checks if package has a signature file.
- **`Package.HasVendorID`** - [Package.HasVendorID](api_metadata.md#313-packagehasvendorid-method)
  - HasVendorID checks if the package has a VendorID (non-zero).
- **`Package.IsMetadataOnlyPackage`** - [Package.IsMetadataOnlyPackage](api_metadata.md#641-packageismetadataonlypackage-method)
  - IsMetadataOnlyPackage checks if package contains only metadata files.
- **`Package.ListSymlinks`** - [Package.ListSymlinks](api_metadata.md#8534-packagelistsymlinks-method)
  - Package.ListSymlinks Returns *PackageError on failure.
- **`Package.loadPathMetadata`** - [Package.loadPathMetadata](api_basic_operations.md#322-packageloadpathmetadata-method)
  - loadPathMetadata loads path metadata from special files Returns *PackageError on failure.
- **`Package.loadSpecialMetadataFiles`** - [Package.loadSpecialMetadataFiles](api_basic_operations.md#321-packageloadspecialmetadatafiles-method)
  - loadSpecialMetadataFiles loads all special metadata files Returns *PackageError on failure.
- **`Package.LoadSymlinkMetadataFile`** - [Package.LoadSymlinkMetadataFile](api_metadata.md#8537-packageloadsymlinkmetadatafile-method)
  - Package.LoadSymlinkMetadataFile Returns *PackageError on failure.
- **`Package.RefreshPackageInfo`** - [Package.RefreshPackageInfo](api_metadata.md#742-packagerefreshpackageinfo-method)
  - RefreshPackageInfo refreshes package information from the file on-disk Returns *PackageError on failure.
- **`Package.RemoveDirectoryMetadata`** - [Package.RemoveDirectoryMetadata](api_metadata.md#82110-packageremovedirectorymetadata-method)
  - RemoveDirectoryMetadata removes directory path metadata (metadata-only, does not remove files) Returns *PackageError on failure.
- **`Package.RemoveIndexFile`** - [Package.RemoveIndexFile](api_metadata.md#534-packageremoveindexfile-method)
  - RemoveIndexFile removes the package index Returns *PackageError on failure.
- **`Package.RemoveManifestFile`** - [Package.RemoveManifestFile](api_metadata.md#524-packageremovemanifestfile-method)
  - RemoveManifestFile removes the package manifest Returns *PackageError on failure.
- **`Package.RemoveMetadataFile`** - [Package.RemoveMetadataFile](api_metadata.md#514-packageremovemetadatafile-method)
  - RemoveMetadataFile removes the package metadata file Returns *PackageError on failure.
- **`Package.RemovePathMetadata`** - [Package.RemovePathMetadata](api_metadata.md#8214-packageremovepathmetadata-method)
  - RemovePathMetadata removes a path metadata entry by path Returns *PackageError on failure.
- **`Package.RemoveSignatureFile`** - [Package.RemoveSignatureFile](api_metadata.md#544-packageremovesignaturefile-method)
  - RemoveSignatureFile removes the signature file Returns *PackageError on failure.
- **`Package.RemoveSpecialFile`** - [Package.RemoveSpecialFile](api_metadata.md#553-packageremovespecialfile-method)
  - RemoveSpecialFile removes a special file by type Returns *PackageError on failure.
- **`Package.RemoveSpecialMetadataFile`** - [Package.RemoveSpecialMetadataFile](api_metadata.md#82116-packageremovespecialmetadatafile-method)
  - RemoveSpecialMetadataFile removes a special metadata file Returns *PackageError on failure.
- **`Package.RemoveSymlink`** - [Package.RemoveSymlink](api_metadata.md#8532-packageremovesymlink-method)
  - Package.RemoveSymlink Returns *PackageError on failure.
- **`Package.SavePathMetadataFile`** - [Package.SavePathMetadataFile](api_metadata.md#8351-packagesavepathmetadatafile-method)
  - SavePathMetadataFile creates and saves the path metadata file.
- **`Package.SaveSymlinkMetadataFile`** - [Package.SaveSymlinkMetadataFile](api_metadata.md#8536-packagesavesymlinkmetadatafile-method)
  - Package.SaveSymlinkMetadataFile Returns *PackageError on failure.
- **`Package.SetAppID`** - [Package.SetAppID](api_metadata.md#21-packagesetappid-method)
  - SetAppID sets or updates the package AppID Returns *PackageError on failure.
- **`Package.SetComment`** - [Package.SetComment](api_metadata.md#111-packagesetcomment-method)
  - SetComment sets or updates the package comment Returns *PackageError on failure.
- **`Package.SetDestPath`** - [Package.SetDestPath](api_metadata.md#8216-packagesetdestpath-method)
  - SetDestPath sets destination extraction directory overrides for a stored path.
  - This is a pure in-memory operation.
  - storedPath MUST be treated as a stored package path.
  - If storedPath does not begin with "/", the implementation MUST prefix "/" before matching or creating entries.
  - If no PathMetadataEntry exists for storedPath, SetDestPath MUST create one.
  - The new entry type MUST be inferred from storedPath.
- **`Package.SetPackageIdentity`** - [Package.SetPackageIdentity](api_metadata.md#41-packagesetpackageidentity-method)
  - SetPackageIdentity sets both VendorID and AppID Returns *PackageError on failure.
- **`Package.SetPathMetadata`** - [Package.SetPathMetadata](api_metadata.md#8212-packagesetpathmetadata-method)
  - SetPathMetadata replaces all path metadata entries in the package Returns *PackageError on failure.
- **`Package.SetVendorID`** - [Package.SetVendorID](api_metadata.md#31-packagesetvendorid-method)
  - SetVendorID sets or updates the package VendorID Returns *PackageError on failure.
- **`Package.TargetExists`** - [Package.TargetExists](api_metadata.md#8542-packagetargetexists-method)
  - TargetExists checks if a path exists as FileEntry or directory PathMetadataEntry.
  - Parameter: ctx: Context for cancellation and timeout.
  - Parameter: path: The path to check (package-relative with leading "/").
  - Return: true if path exists as FileEntry or directory PathMetadataEntry.
  - Return: false otherwise.
- **`Package.UpdateDirectoryMetadata`** - [Package.UpdateDirectoryMetadata](api_metadata.md#82111-packageupdatedirectorymetadata-method)
  - UpdateDirectoryMetadata updates directory path metadata (metadata-only, does not modify files) Returns *PackageError on failure.
- **`Package.UpdateFile`** - [Package.UpdateFile](api_file_mgmt_updates.md#11-packageupdatefile-method)
  - UpdateFile updates file content and metadata in the package The new file data is read from the sourceFilePath on the filesystem.
  - The storedPath identifies which file in the package to update.
- **`Package.UpdateFileMetadata`** - [Package.UpdateFileMetadata](api_file_mgmt_updates.md#13-packageupdatefilemetadata-method)
  - UpdateFileMetadata updates file metadata without changing content.
- **`Package.UpdateIndexFile`** - [Package.UpdateIndexFile](api_metadata.md#533-packageupdateindexfile-method)
  - UpdateIndexFile updates the package index Returns *PackageError on failure.
- **`Package.UpdateManifestFile`** - [Package.UpdateManifestFile](api_metadata.md#523-packageupdatemanifestfile-method)
  - UpdateManifestFile updates the package manifest Returns *PackageError on failure.
- **`Package.UpdateMetadataFile`** - [Package.UpdateMetadataFile](api_metadata.md#513-packageupdatemetadatafile-method)
  - UpdateMetadataFile updates the package metadata file Returns *PackageError on failure.
- **`Package.UpdatePathMetadata`** - [Package.UpdatePathMetadata](api_metadata.md#8215-packageupdatepathmetadata-method)
  - UpdatePathMetadata updates an existing path metadata entry Returns *PackageError on failure.
- **`Package.UpdateSignatureFile`** - [Package.UpdateSignatureFile](api_metadata.md#543-packageupdatesignaturefile-method)
  - UpdateSignatureFile updates the signature file Returns *PackageError on failure.
- **`Package.UpdateSpecialMetadataFile`** - [Package.UpdateSpecialMetadataFile](api_metadata.md#82115-packageupdatespecialmetadatafile-method)
  - UpdateSpecialMetadataFile updates an existing special metadata file Returns *PackageError on failure.
- **`Package.UpdateSpecialMetadataFlags`** - [Package.UpdateSpecialMetadataFlags](api_metadata.md#8352-packageupdatespecialmetadataflags-method)
  - UpdateSpecialMetadataFlags updates package header flags based on special files.
- **`Package.UpdateSymlink`** - [Package.UpdateSymlink](api_metadata.md#8535-packageupdatesymlink-method)
  - Package.UpdateSymlink Returns *PackageError on failure.
- **`Package.ValidateMetadataOnlyIntegrity`** - [Package.ValidateMetadataOnlyIntegrity](api_metadata.md#644-packagevalidatemetadataonlyintegrity-method)
  - ValidateMetadataOnlyIntegrity validates metadata-only package integrity Returns *PackageError on failure.
- **`Package.ValidateMetadataOnlyPackage`** - [Package.ValidateMetadataOnlyPackage](api_metadata.md#645-packagevalidatemetadataonlypackage-method)
  - ValidateMetadataOnlyPackage performs comprehensive validation.
  - Validation: Ensure FileCount == 0 Ensure IsMetadataOnly flag (Bit 7) is set in header Validate all special metadata files (if present) Check for malicious metadata patterns (if metadata files present) Verify signature scope includes all metadata (if signatures present) Ensure metadata consistency Validate package structure.
- **`Package.ValidatePathMetadata`** - [Package.ValidatePathMetadata](api_metadata.md#82112-packagevalidatepathmetadata-method)
  - ValidatePathMetadata validates all path metadata entries Returns *PackageError on failure.
- **`Package.ValidatePathWithinPackageRoot`** - [Package.ValidatePathWithinPackageRoot](api_metadata.md#8543-packagevalidatepathwithinpackageroot-method)
  - ValidatePathWithinPackageRoot validates that a path is within package root.
  - Parameter: path: The path to validate (package-relative).
  - Return: Normalized path if valid.
  - Return: Error if path escapes package root or is invalid.
  - Return: Returns ErrTypeValidation for invalid format.
  - Return: Returns ErrTypeSecurity for paths escaping package root.
- **`Package.ValidateSpecialFiles`** - [Package.ValidateSpecialFiles](api_metadata.md#554-packagevalidatespecialfiles-method)
  - ValidateSpecialFiles validates all special files Returns *PackageError on failure.
- **`Package.ValidateSymlinkPaths`** - [Package.ValidateSymlinkPaths](api_metadata.md#8541-packagevalidatesymlinkpaths-method)
  - ValidateSymlinkPaths validates symlink source and target paths.
  - Parameter: ctx: Context for cancellation and timeout.
  - Parameter: sourcePath: The symlink source path (where the symlink will be created).
  - Parameter: targetPath: The symlink target path (where the symlink points to).
  - Return: Error if validation fails Validation performed:.
  - Return: Both paths are package-relative (start with "/").

### 1.5 Package Compression Methods

- **`Package.CanCompressPackage`** - [Package.CanCompressPackage](api_package_compression.md#7282-packagecancompresspackage-method)
  - CanCompressPackage checks if package can be compressed (not signed).
- **`Package.CompressFile`** - [Package.CompressFile](api_file_mgmt_compression.md#16-packagecompressfile-method)
  - CompressFile compresses an existing file in the package by path This is a convenience wrapper that looks up the FileEntry and calls Compress Returns *PackageError on failure.
- **`Package.CompressGeneric`** - [7.2.10.1 Package.CompressGeneric Method](api_package_compression.md#72101-packagecompressgeneric-method)
  - Package.CompressGeneric Generic compression methods for type-safe operations CompressionStrategy[T] embeds Strategy[T, T] from the generics package See [Core Generic Types](api_generics.md#1-core-generic-types) for base strategy pattern.
- **`Package.CompressPackage`** - [Package.CompressPackage](api_package_compression.md#41-packagecompresspackage-method)
  - CompressPackage compresses package content in memory Compresses file entries and data separately using LZ4 for metadata and specified type for data Compresses file index with LZ4 Creates metadata index for fast access (NOT header, metadata index, comment, or signatures) Signed packages cannot be compressed Returns *PackageError on failure.
- **`Package.CompressPackageConcurrent`** - [Package.CompressPackageConcurrent](api_package_compression.md#831-packagecompresspackageconcurrent-method)
  - CompressPackageConcurrent compresses package content using worker pool Returns *PackageError on failure.
- **`Package.compressPackageContent`** - [Package.compressPackageContent](api_package_compression.md#731-packagecompresspackagecontent-method)
  - Package.compressPackageContent Internal compression methods (used by CompressPackage and Write) Returns *PackageError on failure.
- **`Package.CompressPackageFile`** - [Package.CompressPackageFile](api_package_compression.md#61-packagecompresspackagefile-method)
  - CompressPackageFile compresses package content and writes to specified path Compresses file entries + data + index (NOT header, comment, or signatures) Signed packages cannot be compressed Returns *PackageError on failure.
- **`Package.CompressPackageStream`** - [Package.CompressPackageStream](api_package_compression.md#51-packagecompresspackagestream-method)
  - CompressPackageStream compresses large package content using streaming Uses temporary files and chunked processing to handle files of any size Configuration determines the level of optimization and memory management Returns *PackageError on failure.
- **`Package.DecompressFile`** - [Package.DecompressFile](api_file_mgmt_compression.md#26-packagedecompressfile-method)
  - DecompressFile decompresses an existing file in the package by path This is a convenience wrapper that looks up the FileEntry and calls Decompress Returns *PackageError on failure.
- **`Package.DecompressGeneric`** - [7.2.10.2 Package.DecompressGeneric Method](api_package_compression.md#72102-packagedecompressgeneric-method)
  - DecompressGeneric decompresses data using a generic compression strategy.
- **`Package.DecompressPackage`** - [Package.DecompressPackage](api_package_compression.md#42-packagedecompresspackage-method)
  - DecompressPackage decompresses the package in memory Decompresses all compressed content Returns *PackageError on failure.
- **`Package.DecompressPackageConcurrent`** - [Package.DecompressPackageConcurrent](api_package_compression.md#832-packagedecompresspackageconcurrent-method)
  - DecompressPackageConcurrent decompresses package content using worker pool Returns *PackageError on failure.
- **`Package.decompressPackageContent`** - [Package.decompressPackageContent](api_package_compression.md#732-packagedecompresspackagecontent-method)
  - Package.decompressPackageContent Returns *PackageError on failure.
- **`Package.DecompressPackageFile`** - [Package.DecompressPackageFile](api_package_compression.md#62-packagedecompresspackagefile-method)
  - DecompressPackageFile decompresses the package and writes to specified path Decompresses all compressed content and writes uncompressed package Returns *PackageError on failure.
- **`Package.DecompressPackageStream`** - [Package.DecompressPackageStream](api_package_compression.md#52-packagedecompresspackagestream-method)
  - DecompressPackageStream decompresses large package content using streaming Uses streaming to manage memory efficiently for large packages Returns *PackageError on failure.
- **`Package.GetFileCompressionInfo`** - [Package.GetFileCompressionInfo](api_file_mgmt_compression.md#35-packagegetfilecompressioninfo-method)
  - GetFileCompressionInfo gets compression information for a file by path This is a convenience wrapper that looks up the FileEntry and calls GetCompressionInfo.
- **`Package.GetPackageCompressedSize`** - [Package.GetPackageCompressedSize](api_package_compression.md#727-packagegetpackagecompressedsize-method)
  - GetPackageCompressedSize returns the compressed size Returns *PackageError if package is not compressed.
- **`Package.GetPackageCompressionInfo`** - [Package.GetPackageCompressionInfo](api_package_compression.md#722-packagegetpackagecompressioninfo-method)
  - GetPackageCompressionInfo returns package compression information.
- **`Package.GetPackageCompressionRatio`** - [Package.GetPackageCompressionRatio](api_package_compression.md#725-packagegetpackagecompressionratio-method)
  - GetPackageCompressionRatio returns the compression ratio Returns *PackageError if package is not compressed.
- **`Package.GetPackageCompressionType`** - [Package.GetPackageCompressionType](api_package_compression.md#724-packagegetpackagecompressiontype-method)
  - GetPackageCompressionType returns the package compression type Returns compression type from header flags bits 15-8 Returns *PackageError if package is not compressed.
- **`Package.GetPackageOriginalSize`** - [Package.GetPackageOriginalSize](api_package_compression.md#726-packagegetpackageoriginalsize-method)
  - GetPackageOriginalSize returns the original size before compression Returns *PackageError if package is not compressed.
- **`Package.IsPackageCompressed`** - [Package.IsPackageCompressed](api_package_compression.md#723-packageispackagecompressed-method)
  - IsPackageCompressed checks if the package is compressed Checks header flags bits 15-8 for compression type.
- **`Package.ListCompressedFiles`** - [Package.ListCompressedFiles](api_file_mgmt_queries.md#421-packagelistcompressedfiles-method)
  - ListCompressedFiles returns all compressed file entries.
- **`Package.SetPackageCompressionType`** - [Package.SetPackageCompressionType](api_package_compression.md#7281-packagesetpackagecompressiontype-method)
  - SetPackageCompressionType sets the package compression type (without compressing) Returns *PackageError on failure.
- **`Package.ValidateCompressionData`** - [7.2.10.3 Package.ValidateCompressionData Method](api_package_compression.md#72103-packagevalidatecompressiondata-method)
  - Package.ValidateCompressionData Returns *PackageError on failure.

### 1.6 Package Path and Configuration Methods

This subsection groups package path and runtime configuration methods (for example, target path and session base management).

### 1.7 Package Signature Management Methods

- **`Package.AddSignature`** - [Package.AddSignature](api_signatures.md#111-packageaddsignature-method)
  - AddSignature adds a new digital signature (appends incrementally) LOW-LEVEL: Use when you have pre-computed signature data Automatically sets the "Has signatures" bit (Flags Bit 0) and SignatureOffset if this is the first signature Ensures signature integrity by validating header state before signing Returns *PackageError on failure.
- **`Package.ClearAllSignatures`** - [Package.ClearAllSignatures](api_signatures.md#116-packageclearallsignatures-method)
  - ClearAllSignatures removes all signatures Returns *PackageError on failure.
- **`Package.GetAllSignatures`** - [Package.GetAllSignatures](api_signatures.md#115-packagegetallsignatures-method)
  - GetAllSignatures gets all signatures.
- **`Package.GetSignature`** - [Package.GetSignature](api_signatures.md#114-packagegetsignature-method)
  - GetSignature gets signature by index Returns *PackageError on failure.
- **`Package.GetSignatureCount`** - [Package.GetSignatureCount](api_signatures.md#113-packagegetsignaturecount-method)
  - GetSignatureCount gets total number of signatures.
- **`Package.GetSignatureStatus`** - [Package.GetSignatureStatus](api_signatures.md#2713-packagegetsignaturestatus-method)
  - GetSignatureStatus returns the current signature status of the package.
- **`Package.RemoveSignature`** - [Package.RemoveSignature](api_signatures.md#112-packageremovesignature-method)
  - RemoveSignature removes signature by index and all later signatures Returns *PackageError on failure.
- **`Package.SignPackage`** - [Package.SignPackage](api_signatures.md#2811-packagesignpackage-method)
  - HIGH-LEVEL: Use when you have a private key and want to generate + add signature Internally calls AddSignature after generating signature data Returns *PackageError on failure.
- **`Package.SignPackageWithKeyFile`** - [Package.SignPackageWithKeyFile](api_signatures.md#2812-packagesignpackagewithkeyfile-method)
  - Package.SignPackageWithKeyFile Returns *PackageError on failure.
- **`Package.SignPackageWithPGP`** - [Package.SignPackageWithPGP](api_signatures.md#2821-packagesignpackagewithpgp-method)
  - Package.SignPackageWithPGP PGP-specific signing (internally calls AddSignature) Returns *PackageError on failure.
- **`Package.SignPackageWithPGPKeyring`** - [Package.SignPackageWithPGPKeyring](api_signatures.md#2822-packagesignpackagewithpgpkeyring-method)
  - Package.SignPackageWithPGPKeyring Returns *PackageError on failure.
- **`Package.SignPackageWithX509`** - [Package.SignPackageWithX509](api_signatures.md#2831-packagesignpackagewithx509-method)
  - Package.SignPackageWithX509 X.509-specific signing (internally calls AddSignature) Returns *PackageError on failure.
- **`Package.SignPackageWithX509Chain`** - [Package.SignPackageWithX509Chain](api_signatures.md#2832-packagesignpackagewithx509chain-method)
  - Package.SignPackageWithX509Chain Returns *PackageError on failure.
- **`Package.UpdateSignature`** - [Package.UpdateSignature](api_signatures.md#284-packageupdatesignature-method)
  - UpdateSignature replaces the most recent signature with new signature data Internally calls AddSignature after removing the previous signature Returns *PackageError on failure.
- **`Package.ValidateAllSignatures`** - [Package.ValidateAllSignatures](api_security.md#1131-packagevalidateallsignatures-method)
  - ValidateAllSignatures validates all signatures in order.
- **`Package.ValidatePGPSignature`** - [Package.ValidatePGPSignature](api_signatures.md#2721-packagevalidatepgpsignature-method)
  - Package.ValidatePGPSignature PGP-specific validation Returns *PackageError on failure.
- **`Package.ValidatePGPSignatureWithKey`** - [Package.ValidatePGPSignatureWithKey](api_signatures.md#2722-packagevalidatepgpsignaturewithkey-method)
  - Package.ValidatePGPSignatureWithKey Returns *PackageError on failure.
- **`Package.ValidateSignature`** - [Package.ValidateSignature](api_signatures.md#2711-packagevalidatesignature-method)
  - Package.ValidateSignature General signature validation Returns *PackageError on failure.
- **`Package.ValidateSignatureChain`** - [Package.ValidateSignatureChain](api_security.md#1134-packagevalidatesignaturechain-method)
  - ValidateSignatureChain validates signature chain integrity.
- **`Package.ValidateSignatureIndex`** - [Package.ValidateSignatureIndex](api_security.md#1133-packagevalidatesignatureindex-method)
  - ValidateSignatureIndex validates signature by index Returns *PackageError on failure.
- **`Package.ValidateSignatureType`** - [Package.ValidateSignatureType](api_security.md#1132-packagevalidatesignaturetype-method)
  - ValidateSignatureType validates signatures of specific type.
- **`Package.ValidateSignatureWithKey`** - [Package.ValidateSignatureWithKey](api_signatures.md#2712-packagevalidatesignaturewithkey-method)
  - Package.ValidateSignatureWithKey Returns *PackageError on failure.
- **`Package.ValidateX509Signature`** - [Package.ValidateX509Signature](api_signatures.md#2731-packagevalidatex509signature-method)
  - Package.ValidateX509Signature X.509-specific validation Returns *PackageError on failure.
- **`Package.ValidateX509SignatureWithChain`** - [Package.ValidateX509SignatureWithChain](api_signatures.md#2732-packagevalidatex509signaturewithchain-method)
  - Package.ValidateX509SignatureWithChain Returns *PackageError on failure.

### 1.8 Package Other Methods

- **`Package.SafeWrite`** - [Package.SafeWrite](api_writing.md#11-packagesafewrite-method)
  - SafeWrite writes the package atomically to the configured target path.
- **`Package.FastWrite`** - [Package.FastWrite](api_writing.md#21-packagefastwrite-method)
  - FastWrite performs in-place updates to an existing package file.
- **`Package.Write`** - [Package.Write](api_writing.md#533-packagewrite-method)
  - Write selects the appropriate write strategy (SafeWrite or FastWrite) based on package state.
- **`Package.EncryptFile`** - [Package.EncryptFile](api_security.md#461-packageencryptfile-method)
  - EncryptFile encrypts a file in the package.
- **`Package.DecryptFile`** - [Package.DecryptFile](api_security.md#462-packagedecryptfile-method)
  - DecryptFile decrypts a file in the package.
- **`Package.ValidateFileEncryption`** - [Package.ValidateFileEncryption](api_security.md#463-packagevalidatefileencryption-method)
  - ValidateFileEncryption validates encryption state and metadata for a file.
- **`Package.GetFileEncryptionInfo`** - [Package.GetFileEncryptionInfo](api_security.md#464-packagegetfileencryptioninfo-method)
  - GetFileEncryptionInfo returns encryption information for a file.

### 1.9 Package Helper Functions

- **`NewPackage`** - [6. NewPackage Function](api_basic_operations.md#6-newpackage-function)
  - NewPackage creates a new empty package.
- **`NewPackageComment`** - [1.3.5 NewPackageComment Function](api_metadata.md#135-newpackagecomment-function)
  - NewPackageComment creates and returns a new PackageComment with zero values.
- **`NewPackageError`** - [10.5.1 NewPackageError Function](api_core.md#1051-newpackageerror-function)
  - NewPackageError creates a structured error with type-safe context All errors must include typed context for type safety.
- **`NewPackageHeader`** - [2.8.2 NewPackageHeader Function](package_file_format.md#282-newpackageheader-function)
  - NewPackageHeader creates and returns a new PackageHeader with default values.
- **`NewPackageInfo`** - [7.1.2 NewPackageInfo Function](api_metadata.md#712-newpackageinfo-function)
  - NewPackageInfo creates a new PackageInfo with default values.
- **`NewPackageWithOptions`** - [7. NewPackageWithOptions Function](api_basic_operations.md#7-newpackagewithoptions-function)
  - NewPackageWithOptions creates a new package with specified configuration options Returns *PackageError on failure.
- **`NormalizePackagePath`** - [Normalizepackagepath](api_core.md#21-path-normalization-rules)
  - NormalizePackagePath normalizes a package-internal path for consistent comparison and storage.
  - Applies separator normalization, dot-segment canonicalization, leading slash, NFC, and path-length checks.
  - Returns the normalized path with leading "/" or an error if the path is invalid or would escape the package root.
- **`OpenBrokenPackage`** - [Openbrokenpackage](api_basic_operations.md#12-openbrokenpackage-function)
  - OpenBrokenPackage opens a package that may be invalid or partially corrupted.
  - This function is intended for repair workflows.
  - Returns *PackageError on failure.
- **`OpenPackage`** - [Openpackage](api_basic_operations.md#10-openpackage-function)
  - OpenPackage opens an existing package from the specified path.
  - It validates the on-disk package structure during open.
  - Returns *PackageError on failure.
- **`OpenPackageReadOnly`** - [Openpackagereadonly](api_basic_operations.md#112-openpackagereadonly-function)
  - OpenPackageReadOnly opens a package in a read-only mode.
  - It validates the on-disk package structure during open.
  - The returned package must reject any attempt to mutate state or write to disk.
  - Returns *PackageError on failure.
- **`ReadHeader`** - [Readheader](api_basic_operations.md#183-readheader-function)
  - ReadHeader reads the package header from a reader.
- **`ReadHeaderFromPath`** - [Readheaderfrompath](api_basic_operations.md#184-readheaderfrompath-function)
  - ReadHeaderFromPath reads the package header from a file path.
- **`ToDisplayPath`** - [Todisplaypath](api_core.md#122-todisplaypath-function)
  - ToDisplayPath converts a stored package path to display format by stripping the leading slash.
- **`ValidatePackagePath`** - [Validatepackagepath](api_core.md#123-validatepackagepath-function)
  - ValidatePackagePath validates a package path according to package path semantics.
  - Validates path format, rejects empty or whitespace-only paths, normalizes separators, canonicalizes dot segments, and ensures the path does not escape the package root.
- **`ValidatePathLength`** - [Validatepathlength](api_core.md#124-validatepathlength-function)
  - ValidatePathLength validates platform path length portability constraints.
  - Return: warnings: Non-fatal portability warnings.
  - Error: ErrTypeValidation when the path exceeds the hard limit.

## 2. PackageReader Interface Types

- **`PackageReader`** - [1.2 PackageReader Interface](api_core.md#12-packagereader-interface)
  - PackageReader defines the read-only interface for opened packages.

### 2.1 PackageReader Methods

This subsection groups `PackageReader` methods by operational category.

#### 2.1.1 PackageReader Read Operations

This subsection groups file read and streaming-oriented `PackageReader` operations.

#### 2.1.2 PackageReader Query Operations

This subsection groups query and inspection `PackageReader` operations.

#### 2.1.3 PackageReader Other Methods

This subsection groups remaining `PackageReader` operations not classified above.

### 2.2 PackageReader Helper Functions

This subsection groups helper functions related to `PackageReader` usage.

## 3. PackageWriter Interface Types

- **`PackageWriter`** - [1.3 PackageWriter Interface](api_core.md#13-packagewriter-interface)
  - PackageWriter defines the write interface for persisting a package to disk.

### 3.1 PackageWriter Methods

This subsection groups `PackageWriter` methods by operational category.

#### 3.1.1 PackageWriter Write Operations

This subsection groups write and persistence-oriented `PackageWriter` operations.

#### 3.1.2 PackageWriter Other Methods

This subsection groups remaining `PackageWriter` operations not classified above.

### 3.2 PackageWriter Helper Functions

This subsection groups helper functions related to `PackageWriter` usage.

## 4. FileEntry Types

- **`FileEntry`** - [1.1 FileEntry Structure Definition](api_file_mgmt_file_entry.md#11-fileentry-structure-definition)
  - FileEntry represents a FileEntry in the package with complete metadata.
- **`HashEntry`** - [11. HashEntry Struct](api_file_mgmt_file_entry.md#11-hashentry-struct)
  - HashEntry represents a hash with type and purpose.
- **`HashType`** - [12. HashType Type](api_file_mgmt_file_entry.md#12-hashtype-type)
  - HashType represents hash algorithm types.
- **`HashPurpose`** - [13. HashPurpose Type](api_file_mgmt_file_entry.md#13-hashpurpose-type)
  - HashPurpose represents the purpose for a hash (for example, deduplication vs integrity).
- **`ProcessingState`** - [15. ProcessingState Type](api_file_mgmt_file_entry.md#15-processingstate-type)
  - ProcessingState defines the current state of file data transformations.
- **`FileSource`** - [16. FileSource Structure](api_file_mgmt_file_entry.md#16-filesource-structure)
  - FileSource represents a source location for file data (original or intermediate).
- **`OptionalData`** - [17. OptionalData Structure](api_file_mgmt_file_entry.md#17-optionaldata-structure)
  - OptionalData represents structured optional data for a FileEntry.
- **`OptionalDataType`** - [18. OptionalDataType Type](api_file_mgmt_file_entry.md#18-optionaldatatype-type)
  - OptionalDataType identifies optional data payload types.

### 4.1 FileEntry Methods

This subsection groups `FileEntry` methods by operational category.

#### 4.1.1 FileEntry Data Management Methods

- **`FileEntry.CleanupTempFile`** - [FileEntry.CleanupTempFile](api_file_mgmt_file_entry.md#425-fileentrycleanuptempfile-method)
  - CleanupTempFile removes temporary files.
- **`FileEntry.CopyCurrentToOriginal`** - [FileEntry.CopyCurrentToOriginal](api_file_mgmt_file_entry.md#448-fileentrycopycurrenttooriginal-method)
  - CopyCurrentToOriginal saves current source as original (for transformations).
- **`FileEntry.GetCurrentSource`** - [FileEntry.GetCurrentSource](api_file_mgmt_file_entry.md#442-fileentrygetcurrentsource-method)
  - GetCurrentSource returns the current data source Returns nil if no current source is set.
- **`FileEntry.GetData`** - [FileEntry.GetData](api_file_mgmt_file_entry.md#413-fileentrygetdata-method)
  - GetData returns the in-memory file data.
  - Returns *PackageError on failure.
- **`FileEntry.GetEncryptionType`** - [FileEntry.GetEncryptionType](api_file_mgmt_file_entry.md#73-fileentrygetencryptiontype-method)
  - GetEncryptionType returns the encryption type used for this file.
- **`FileEntry.GetOriginalSource`** - [FileEntry.GetOriginalSource](api_file_mgmt_file_entry.md#444-fileentrygetoriginalsource-method)
  - GetOriginalSource returns the original data source.
  - Returns nil if no original source is tracked (e.g., new files).
- **`FileEntry.GetPrimaryPath`** - [FileEntry.GetPrimaryPath](api_file_mgmt_file_entry.md#53-fileentrygetprimarypath-method)
  - GetPrimaryPath returns the primary path in display format (no leading slash).
  - The returned string MUST use forward slashes as separators.
  - For platform-specific filesystem display, convert manually or use path conversion utilities.
- **`FileEntry.GetProcessingState`** - [FileEntry.GetProcessingState](api_file_mgmt_file_entry.md#431-fileentrygetprocessingstate-method)
  - GetProcessingState returns the current processing state.
- **`FileEntry.GetSymlinkPaths`** - [FileEntry.GetSymlinkPaths](api_file_mgmt_file_entry.md#52-fileentrygetsymlinkpaths-method)
  - GetSymlinkPaths returns all symlink paths associated with this FileEntry.
- **`FileEntry.GetTransformPipeline`** - [FileEntry.GetTransformPipeline](api_file_mgmt_file_entry.md#452-fileentrygettransformpipeline-method)
  - GetTransformPipeline returns the current transformation pipeline.
  - Returns nil if no pipeline is active.
- **`FileEntry.AssociateWithPathMetadata`** - [FileEntry.AssociateWithPathMetadata](api_file_mgmt_file_entry.md#55-fileentryassociatewithpathmetadata-method)
  - AssociateWithPathMetadata associates the FileEntry with a PathMetadataEntry.
- **`FileEntry.GetPathMetadataForPath`** - [FileEntry.GetPathMetadataForPath](api_file_mgmt_file_entry.md#56-fileentrygetpathmetadataforpath-method)
  - GetPathMetadataForPath returns the PathMetadataEntry for a given stored path.
- **`FileEntry.IsCompressed`** - [FileEntry.IsCompressed](api_file_mgmt_file_entry.md#71-fileentryiscompressed-method)
  - IsCompressed returns true if the file is compressed.
- **`FileEntry.GetCompressionInfo`** - [FileEntry.GetCompressionInfo](api_file_mgmt_file_entry.md#83-fileentrygetcompressioninfo-method)
  - GetCompressionInfo returns compression details for the entry.
- **`FileEntry.HasEncryptionKey`** - [FileEntry.HasEncryptionKey](api_file_mgmt_file_entry.md#72-fileentryhasencryptionkey-method)
  - HasEncryptionKey checks if the file has an encryption key set.
- **`FileEntry.HasOriginalSource`** - [FileEntry.HasOriginalSource](api_file_mgmt_file_entry.md#446-fileentryhasoriginalsource-method)
  - HasOriginalSource returns true if original source is tracked.
- **`FileEntry.HasSymlinks`** - [FileEntry.HasSymlinks](api_file_mgmt_file_entry.md#51-fileentryhassymlinks-method)
  - HasSymlinks returns true if the FileEntry has any symlink paths.
- **`FileEntry.IsCurrentSourceTempFile`** - [FileEntry.IsCurrentSourceTempFile](api_file_mgmt_file_entry.md#445-fileentryiscurrentsourcetempfile-method)
  - IsCurrentSourceTempFile returns true if current source is a temporary file.
- **`FileEntry.IsEncrypted`** - [FileEntry.IsEncrypted](api_file_mgmt_file_entry.md#74-fileentryisencrypted-method)
  - IsEncrypted checks if the file is encrypted.
- **`FileEntry.LoadData`** - [FileEntry.LoadData](api_file_mgmt_file_entry.md#101-fileentryloaddata-method)
  - LoadData loads the file data into memory.
- **`FileEntry.Marshal`** - [FileEntry.Marshal](api_file_mgmt_file_entry.md#613-fileentrymarshal-method)
  - Marshal marshals both FileEntry metadata and data.
  - Returns metadata and data as separate byte slices for flexible writing.
  - Returns *PackageError on failure.
- **`FileEntry.MarshalMeta`** - [FileEntry.MarshalMeta](api_file_mgmt_file_entry.md#611-fileentrymarshalmeta-method)
  - MarshalMeta marshals FileEntry metadata to bytes.
- **`FileEntry.MarshalData`** - [FileEntry.MarshalData](api_file_mgmt_file_entry.md#612-fileentrymarshaldata-method)
  - MarshalData marshals FileEntry data to bytes.
- **`FileEntry.ReadFromTempFile`** - [FileEntry.ReadFromTempFile](api_file_mgmt_file_entry.md#424-fileentryreadfromtempfile-method)
  - ReadFromTempFile reads data from a temporary file.
- **`FileEntry.CreateTempFile`** - [FileEntry.CreateTempFile](api_file_mgmt_file_entry.md#421-fileentrycreatetempfile-method)
  - CreateTempFile creates a temporary file for staging file data.
- **`FileEntry.ResolveAllSymlinks`** - [FileEntry.ResolveAllSymlinks](api_file_mgmt_file_entry.md#54-fileentryresolveallsymlinks-method)
  - ResolveAllSymlinks resolves all symlink paths to their target paths.
- **`FileEntry.StreamToTempFile`** - [FileEntry.StreamToTempFile](api_file_mgmt_file_entry.md#422-fileentrystreamtotempfile-method)
  - StreamToTempFile streams data to a temporary file Returns *PackageError on failure.
- **`FileEntry.UnloadData`** - [FileEntry.UnloadData](api_file_mgmt_file_entry.md#412-fileentryunloaddata-method)
  - UnloadData unloads file data from memory.
- **`FileEntry.UnsetEncryptionKey`** - [FileEntry.UnsetEncryptionKey](api_file_mgmt_file_entry.md#94-fileentryunsetencryptionkey-method)
  - UnsetEncryptionKey removes the encryption key from the file.
- **`FileEntry.ValidateSources`** - [FileEntry.ValidateSources](api_file_mgmt_file_entry.md#456-fileentryvalidatesources-method)
  - ValidateSources validates CurrentSource, OriginalSource, and pipeline consistency Returns *PackageError if validation fails.
- **`FileEntry.WriteDataTo`** - [FileEntry.WriteDataTo](api_file_mgmt_file_entry.md#622-fileentrywritedatato-method)
  - WriteDataTo writes the FileEntry data to a writer.
  - Implements efficient streaming for large files.
  - Returns *PackageError on failure.
- **`FileEntry.WriteMetaTo`** - [FileEntry.WriteMetaTo](api_file_mgmt_file_entry.md#621-fileentrywritemetato-method)
  - WriteMetaTo writes the FileEntry metadata to a writer.
  - Implements efficient streaming for large metadata.
  - Returns *PackageError on failure.
- **`FileEntry.WriteTo`** - [FileEntry.WriteTo](api_file_mgmt_file_entry.md#623-fileentrywriteto-method)
  - WriteTo writes both metadata and data to a writer.
  - Implements io.WriterTo interface.
  - Returns *PackageError on failure.
- **`FileEntry.WriteToTempFile`** - [FileEntry.WriteToTempFile](api_file_mgmt_file_entry.md#423-fileentrywritetotempfile-method)
  - WriteToTempFile writes data to a temporary file.
  - Returns *PackageError on failure.

#### 4.1.2 FileEntry Transformation Methods

- **`FileEntry.CleanupTransformPipeline`** - [FileEntry.CleanupTransformPipeline](api_file_mgmt_file_entry.md#455-fileentrycleanuptransformpipeline-method)
  - CleanupTransformPipeline cleans up all temporary files in pipeline.
  - Returns *PackageError on failure.
- **`FileEntry.Compress`** - [FileEntry.Compress](api_file_mgmt_file_entry.md#81-fileentrycompress-method)
  - Compress applies compression to the FileEntry data.
- **`FileEntry.Decrypt`** - [FileEntry.Decrypt](api_file_mgmt_file_entry.md#93-fileentrydecrypt-method)
  - Decrypt decrypts data using the file's encryption key.
- **`FileEntry.Decompress`** - [FileEntry.Decompress](api_file_mgmt_file_entry.md#82-fileentrydecompress-method)
  - Decompress reverses compression on the FileEntry data.
- **`FileEntry.Encrypt`** - [FileEntry.Encrypt](api_file_mgmt_file_entry.md#92-fileentryencrypt-method)
  - Encrypt encrypts data using the file's encryption key.
- **`FileEntry.InitializeTransformPipeline`** - [FileEntry.InitializeTransformPipeline](api_file_mgmt_file_entry.md#451-fileentryinitializetransformpipeline-method)
  - InitializeTransformPipeline creates a new transformation pipeline.
- **`FileEntry.ExecuteTransformStage`** - [FileEntry.ExecuteTransformStage](api_file_mgmt_file_entry.md#453-fileentryexecutetransformstage-method)
  - ExecuteTransformStage executes a single transformation stage.
- **`FileEntry.ProcessData`** - [FileEntry.ProcessData](api_file_mgmt_file_entry.md#102-fileentryprocessdata-method)
  - ProcessData processes FileEntry data through the configured pipeline.
- **`FileEntry.ResumeTransformation`** - [FileEntry.ResumeTransformation](api_file_mgmt_file_entry.md#454-fileentryresumetransformation-method)
  - ResumeTransformation resumes pipeline from last completed stage.
  - Returns *PackageError on failure.
- **`FileEntry.SetCurrentSource`** - [FileEntry.SetCurrentSource](api_file_mgmt_file_entry.md#441-fileentrysetcurrentsource-method)
  - SetCurrentSource sets the current data source for the FileEntry Returns *PackageError if source is invalid.
- **`FileEntry.SetData`** - [FileEntry.SetData](api_file_mgmt_file_entry.md#414-fileentrysetdata-method)
  - SetData sets the in-memory file data.
- **`FileEntry.SetEncryptionKey`** - [FileEntry.SetEncryptionKey](api_file_mgmt_file_entry.md#91-fileentrysetencryptionkey-method)
  - SetEncryptionKey sets the encryption key for the file.
- **`FileEntry.SetOriginalSource`** - [FileEntry.SetOriginalSource](api_file_mgmt_file_entry.md#443-fileentrysetoriginalsource-method)
  - SetOriginalSource sets the original data source before transformations.
- **`FileEntry.SetOriginalSourceFromPackage`** - [FileEntry.SetOriginalSourceFromPackage](api_file_mgmt_file_entry.md#447-fileentrysetoriginalsourcefrompackage-method)
  - SetOriginalSourceFromPackage creates original source pointing to package file.
- **`FileEntry.SetProcessingState`** - [FileEntry.SetProcessingState](api_file_mgmt_file_entry.md#432-fileentrysetprocessingstate-method)
  - SetProcessingState sets the current processing state.

### 4.2 FileEntry Helper Functions

- **`NewFileEntry`** - [2.1.1 NewFileEntry Function Signature](api_file_mgmt_file_entry.md#211-newfileentry-function-signature)
  - NewFileEntry creates a new FileEntry with proper tag synchronization.
- **`UnmarshalFileEntry`** - [Unmarshalfileentry](api_file_mgmt_file_entry.md#63-unmarshalfileentry-function)
  - UnmarshalFileEntry unmarshals a FileEntry from binary data.
  - Unmarshals the FileEntry with proper tag synchronization.

## 5. Metadata Types

- **`ACLEntry`** - [8.1.6 ACLEntry Structure](api_metadata.md#816-aclentry-structure)
  - ACLEntry represents an Access Control List entry.
- **`DestPathInput`** - [8.1.12 DestPathInput Interface](api_metadata.md#8112-destpathinput-interface)
  - DestPathInput is the allowed input type set for SetDestPath.
  - DestPathInput supports.
  - string: a single destination string.
  - map[string]string: a map with keys "DestPath" and/or "DestPathWin" Note: The map form uses string keys for ergonomics in callers.
    Keys other than "DestPath" and "DestPathWin" MUST be rejected with ErrTypeValidation.
- **`DestPathOverride`** - [8.1.11 DestPathOverride Structure](api_metadata.md#8111-destpathoverride-structure)
  - DestPathOverride specifies destination extraction directory overrides.
  - A nil field means "no override specified" for that field.
- **`FileMetadataUpdate`** - [1.3.3 FileMetadataUpdate Structure](api_file_mgmt_updates.md#133-filemetadataupdate-structure)
  - FileMetadataUpdate contains metadata updates for a FileEntry.
- **`FilePathAssociation`** - [8.1.10 FilePathAssociation Structure](api_metadata.md#8110-filepathassociation-structure)
  - FilePathAssociation links files to their path metadata.
- **`IndexData`** - [5.5.5.3 IndexData Structure](api_metadata.md#5553-indexdata-structure)
  - IndexData contains index file data structure.
- **`ManifestData`** - [5.5.5.2 ManifestData Structure](api_metadata.md#5552-manifestdata-structure)
  - ManifestData contains manifest file data structure.
- **`PackageComment`** - [1.2 PackageComment Structure](api_metadata.md#12-packagecomment-structure)
  - PackageComment represents the optional package comment section.
- **`PackageHeader`** - [7.1.5 PackageHeader Structure](api_metadata.md#715-packageheader-structure)
  - PackageHeader represents the fixed-size header of a NovusPack (.nvpk) file Size: 112 bytes (fixed).
- **`PackageInfo`** - [7.1 PackageInfo Structure](api_metadata.md#71-packageinfo-structure)
  - PackageInfo contains comprehensive package information and metadata.
- **`PathFileSystem`** - [8.1.5 PathFileSystem Structure](api_metadata.md#815-pathfilesystem-structure)
  - PathFileSystem contains filesystem-specific properties.
- **`PathInfo`** - [8.1.9 PathInfo Structure](api_metadata.md#819-pathinfo-structure)
  - PathInfo provides runtime path metadata information.
- **`PathInheritance`** - [8.1.3 PathInheritance Structure](api_metadata.md#813-pathinheritance-structure)
  - PathInheritance controls tag inheritance behavior (for directories only).
- **`PathMetadata`** - [8.1.4 PathMetadata Structure](api_metadata.md#814-pathmetadata-structure)
  - PathMetadata contains path metadata (for directories only).
- **`PathMetadataEntry`** - [8.1.2 PathMetadataEntry Structure](api_metadata.md#812-pathmetadataentry-structure)
  - PathMetadataEntry represents a path (file, directory, or symlink) with metadata, inheritance rules, and filesystem properties.
- **`PathMetadataPatch`** - [2.8.2 PathMetadataPatch Struct](api_file_mgmt_addition.md#282-pathmetadatapatch-struct)
  - PathMetadataPatch specifies persisted PathMetadataEntry fields to create or update at add time.
  - This patch is applied to the PathMetadataEntry for the derived stored path.
  - It does not change the stored path itself.
  - Cross-Reference.
  - Cross-Reference: [Package Metadata API.
  - Cross-Reference: PathMetadataEntry Structure](api_metadata.md#812-pathmetadataentry-structure).
- **`PathMetadataType`** - [8.1.1 PathMetadataType Type](api_metadata.md#811-pathmetadatatype-type)
  - PathMetadataType represents the type of path entry.
- **`PathNode`** - [8.5.6 PathNode Structure](api_metadata.md#856-pathnode-structure)
  - PathNode represents a node in the path tree.
- **`PathStats`** - [8.5.7 PathStats Structure](api_metadata.md#857-pathstats-structure)
  - PathStats provides statistics for a path.
- **`PathTree`** - [8.5.5 PathTree Structure](api_metadata.md#855-pathtree-structure)
  - PathTree represents the complete path hierarchy.
- **`SecurityStatus`** - [Securitystatus](api_metadata.md#73-securitystatus-structure)
  - SecurityStatus contains the security status of a package.
- **`SignatureData`** - [Signaturedata](api_metadata.md#5554-signaturedata-structure)
  - SignatureData contains signature file data structure.
- **`SignatureInfo`** - [Signatureinfo](api_metadata.md#72-signatureinfo-structure)
  - SignatureInfo contains signature information for a package.
- **`SpecialFileInfo`** - [5.5.5.1 SpecialFileInfo Structure](api_metadata.md#5551-specialfileinfo-structure)
  - SpecialFileInfo contains information about special metadata files in the package.
- **`SymlinkEntry`** - [8.5.8 SymlinkEntry Structure](api_metadata.md#858-symlinkentry-structure)
  - SymlinkEntry represents a symbolic link with metadata.
- **`SymlinkFileSystem`** - [8.5.10 SymlinkFileSystem Structure](api_metadata.md#8510-symlinkfilesystem-structure)
  - SymlinkFileSystem contains filesystem-specific properties for symlinks.
- **`SymlinkMetadata`** - [8.5.9 SymlinkMetadata Structure](api_metadata.md#859-symlinkmetadata-structure)
  - SymlinkMetadata contains symlink creation and modification information.
- **`Tag`** - [19.1 Tag Struct](api_file_mgmt_file_entry.md#191-tag-struct)
  - Tag represents a type-safe tag with a specific value type.
- **`TagValueType`** - [14. TagValueType Type](api_file_mgmt_file_entry.md#14-tagvaluetype-type)
  - TagValueType represents the type of a tag value.
- **`TransformStage`** - [2.3 TransformStage Structure](api_file_mgmt_transform_pipelines.md#23-transformstage-structure)
  - TransformStage represents a single transformation stage.

### 5.1 Metadata Methods

- **`PackageComment.ReadFrom`** - [PackageComment.ReadFrom](api_metadata.md#133-packagecommentreadfrom-method)
  - ReadFrom reads the comment from a reader.
- **`PackageComment.Size`** - [PackageComment.Size](api_metadata.md#131-packagecommentsize-method)
  - Size returns the size of the package comment.
- **`PackageComment.Validate`** - [PackageComment.Validate](api_metadata.md#134-packagecommentvalidate-method)
  - Validate validates the package comment Returns *PackageError on failure.
- **`PackageComment.WriteTo`** - [PackageComment.WriteTo](api_metadata.md#132-packagecommentwriteto-method)
  - WriteTo writes the comment to a writer.
- **`PackageHeader.ToHeader`** - [PackageHeader.ToHeader](api_metadata.md#716-packageheadertoheader-method)
  - ToHeader synchronizes PackageHeader fields from the provided PackageInfo.
  - This method must only write fields that are represented in the header.
  - It must not mutate fields that are computed by the writer pipeline (for example IndexStart, IndexSize, and CRC).
  - Returns *PackageError on failure.
- **`PackageInfo.FromHeader`** - [PackageInfo.FromHeader](api_metadata.md#714-packageinfofromheader-method)
  - FromHeader synchronizes PackageInfo fields from the provided PackageHeader.
  - This method must only copy data that is represented in the header.
  - It must not compute derived values that require scanning file entries or reading file data.
  - Returns *PackageError on failure.
- **`PathMetadataEntry.AssociateWithFileEntry`** - [8.1.8.19 PathMetadataEntry.AssociateWithFileEntry Method](api_metadata.md#81819-pathmetadataentryassociatewithfileentry-method)
  - PathMetadataEntry.AssociateWithFileEntry FileEntry association methods for PathMetadataEntry AssociateWithFileEntry associates this PathMetadataEntry with a FileEntry The association is established if the PathMetadataEntry.Path.Path matches one of the FileEntry.Paths Returns *PackageError on failure.
- **`PathMetadataEntry.GetAncestors`** - [8.1.8.16 PathMetadataEntry.GetAncestors Method](api_metadata.md#81816-pathmetadataentrygetancestors-method)
  - GetAncestors returns all ancestor path metadata entries up to the root.
- **`PathMetadataEntry.GetAssociatedFileEntries`** - [8.1.8.20 PathMetadataEntry.GetAssociatedFileEntries Method](api_metadata.md#81820-pathmetadataentrygetassociatedfileentries-method)
  - GetAssociatedFileEntries returns all FileEntry instances associated with this PathMetadataEntry Returns empty slice if no FileEntry instances are associated.
- **`PathMetadataEntry.GetDepth`** - [8.1.8.14 PathMetadataEntry.GetDepth Method](api_metadata.md#81814-pathmetadataentrygetdepth-method)
  - GetDepth returns the depth of this path in the directory hierarchy.
- **`PathMetadataEntry.GetEffectiveTags`** - [8.1.8.18 PathMetadataEntry.GetEffectiveTags Method](api_metadata.md#81818-pathmetadataentrygeteffectivetags-method)
  - GetEffectiveTags returns all tags for this PathMetadataEntry, including: 1.
  - Tags directly on this PathMetadataEntry 2.
  - Tags inherited from parent PathMetadataEntry instances (path hierarchy) 3.
  - Tags from associated FileEntry instances (treated as if applied to this PathMetadataEntry).
- **`PathMetadataEntry.GetInheritedTags`** - [8.1.8.17 PathMetadataEntry.GetInheritedTags Method](api_metadata.md#81817-pathmetadataentrygetinheritedtags-method)
  - PathMetadataEntry.GetInheritedTags Tag inheritance methods for PathMetadataEntry These methods resolve inheritance by walking up the ParentPath chain.
- **`PathMetadataEntry.GetLinkTarget`** - [8.1.8.9 PathMetadataEntry.GetLinkTarget Method](api_metadata.md#8189-pathmetadataentrygetlinktarget-method)
  - GetLinkTarget returns the target path of the symlink.
- **`PathMetadataEntry.GetParentPath`** - [8.1.8.12 PathMetadataEntry.GetParentPath Method](api_metadata.md#81812-pathmetadataentrygetparentpath-method)
  - GetParentPath returns the parent path metadata entry.
- **`PathMetadataEntry.GetParentPathString`** - [8.1.8.13 PathMetadataEntry.GetParentPathString Method](api_metadata.md#81813-pathmetadataentrygetparentpathstring-method)
  - GetParentPathString returns the parent path as a string.
- **`PathMetadataEntry.GetPath`** - [8.1.8.2 PathMetadataEntry.GetPath Method](api_metadata.md#8182-pathmetadataentrygetpath-method)
  - GetPath returns the path as stored (Unix-style with forward slashes).
  - For platform-specific display, use GetPathForPlatform() or convert manually.
- **`PathMetadataEntry.GetPathEntry`** - [8.1.8.4 PathMetadataEntry.GetPathEntry Method](api_metadata.md#8184-pathmetadataentrygetpathentry-method)
  - GetPathEntry returns the PathEntry representation of this path metadata entry.
- **`PathMetadataEntry.GetPathForPlatform`** - [8.1.8.3 PathMetadataEntry.GetPathForPlatform Method](api_metadata.md#8183-pathmetadataentrygetpathforplatform-method)
  - GetPathForPlatform returns the path converted for the specified platform.
  - On Windows: converts forward slashes to backslashes.
  - On Unix/Linux: returns the path as stored (with forward slashes).
- **`PathMetadataEntry.GetType`** - [8.1.8.5 PathMetadataEntry.GetType Method](api_metadata.md#8185-pathmetadataentrygettype-method)
  - PathMetadataEntry.GetType Type and symlink methods for PathMetadataEntry.
- **`PathMetadataEntry.IsDirectory`** - [8.1.8.6 PathMetadataEntry.IsDirectory Method](api_metadata.md#8186-pathmetadataentryisdirectory-method)
  - IsDirectory returns true if this path metadata entry represents a directory.
- **`PathMetadataEntry.IsFile`** - [8.1.8.7 PathMetadataEntry.IsFile Method](api_metadata.md#8187-pathmetadataentryisfile-method)
  - IsFile returns true if this path metadata entry represents a file.
- **`PathMetadataEntry.IsRoot`** - [8.1.8.15 PathMetadataEntry.IsRoot Method](api_metadata.md#81815-pathmetadataentryisroot-method)
  - IsRoot returns true if this path metadata entry represents the root path.
- **`PathMetadataEntry.IsSymlink`** - [8.1.8.8 PathMetadataEntry.IsSymlink Method](api_metadata.md#8188-pathmetadataentryissymlink-method)
  - IsSymlink returns true if this path metadata entry represents a symlink.
- **`PathMetadataEntry.ResolveSymlink`** - [8.1.8.10 PathMetadataEntry.ResolveSymlink Method](api_metadata.md#81810-pathmetadataentryresolvesymlink-method)
  - ResolveSymlink resolves the symlink to its final target path.
- **`PathMetadataEntry.SetParentPath`** - [8.1.8.11 PathMetadataEntry.SetParentPath Method](api_metadata.md#81811-pathmetadataentrysetparentpath-method)
  - PathMetadataEntry.SetParentPath Parent path management methods for PathMetadataEntry.
- **`PathMetadataEntry.SetPath`** - [8.1.8.1 PathMetadataEntry.SetPath Method](api_metadata.md#8181-pathmetadataentrysetpath-method)
  - PathMetadataEntry.SetPath Path management methods for PathMetadataEntry.
- **`Tag.GetValue`** - [Tag.GetValue](api_file_mgmt_file_entry.md#193-tagtgetvalue-method)
  - GetValue returns the type-safe value of the tag.
- **`Tag.SetValue`** - [Tag.SetValue](api_file_mgmt_file_entry.md#194-tagtsetvalue-method)
  - SetValue sets the type-safe value of the tag.

### 5.2 Metadata Helper Functions

- **`AddFileEntryTag`** - [Addfileentrytag](api_file_mgmt_file_entry.md#3126-addfileentrytag-function)
  - AddFileEntryTag adds a new tag with type safety to a FileEntry.
  - Returns *PackageError if a tag with the same key already exists Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`AddFileEntryTags`** - [Addfileentrytags](api_file_mgmt_file_entry.md#3124-addfileentrytags-function)
  - AddFileEntryTags adds multiple new tags with type safety to a FileEntry.
  - Returns *PackageError if any tag with the same key already exists.
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`AddPathMetaTag`** - [Addpathmetatag](api_metadata.md#8176-addpathmetatag-function)
  - AddPathMetaTag adds a new tag with type safety to a PathMetadataEntry Returns *PackageError if a tag with the same key already exists.
- **`AddPathMetaTags`** - [Addpathmetatags](api_metadata.md#8173-addpathmetatags-function)
  - AddPathMetaTags adds multiple new tags with type safety to a PathMetadataEntry Returns *PackageError if any tag with the same key already exists.
- **`AuditSignatureComment`** - [1.5.4 AuditSignatureComment Function](api_metadata.md#154-auditsignaturecomment-function)
  - AuditSignatureComment logs signature comment for security auditing Returns *PackageError on failure.
- **`CheckCommentLength`** - [Checkcommentlength](api_metadata.md#144-checkcommentlength-function)
  - CheckCommentLength validates comment length against limits Returns *PackageError on failure.
- **`CheckSignatureCommentLength`** - [1.5.3 CheckSignatureCommentLength Function](api_metadata.md#153-checksignaturecommentlength-function)
  - CheckSignatureCommentLength validates signature comment length Returns *PackageError on failure.
- **`DetectInjectionPatterns`** - [Detectinjectionpatterns](api_metadata.md#145-detectinjectionpatterns-function)
  - DetectInjectionPatterns scans comment for malicious patterns.
- **`GetFileEntryTag`** - [Getfileentrytag](api_file_mgmt_file_entry.md#3123-getfileentrytag-function)
  - GetFileEntryTag retrieves a type-safe tag by key from a FileEntry.
  - Returns the tag pointer and an error.
  - If the tag is not found, returns (nil, nil).
  - If an underlying error occurs (corruption, I/O), returns (nil, error).
  - Returns *PackageError on failure.
  - If the tag type is unknown, use GetFileEntryTag[any](fe, "key") to retrieve the tag and inspect its Type field.
- **`GetFileEntryTags`** - [Getfileentrytags](api_file_mgmt_file_entry.md#3121-getfileentrytags-function)
  - GetFileEntryTags returns all tags as typed tags for a FileEntry.
  - Returns a slice of Tag pointers, where each tag maintains its type information.
  - Returns *PackageError on failure (corruption, I/O).
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`GetFileEntryTagsByType`** - [Getfileentrytagsbytype](api_file_mgmt_file_entry.md#3122-getfileentrytagsbytype-function)
  - GetFileEntryTagsByType returns all tags of a specific type for a FileEntry.
  - Returns a slice of Tag pointers with the specified type parameter T.
  - Only tags matching the type T and corresponding TagValueType are returned.
  - Returns *PackageError on failure (corruption, I/O).
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`GetPathMetaTag`** - [Getpathmetatag](api_metadata.md#8175-getpathmetatag-function)
  - GetPathMetaTag retrieves a type-safe tag by key from a PathMetadataEntry Returns the tag pointer and an error.
  - If the tag is not found, returns (nil, nil).
  - If an underlying error occurs, returns (nil, error).
  - Returns *PackageError on failure If the tag type is unknown, use GetPathMetaTag[any](pme, "key") to retrieve the tag and inspect its Type field.
- **`GetPathMetaTags`** - [Getpathmetatags](api_metadata.md#8171-getpathmetatags-function)
  - GetPathMetaTags returns all tags as typed tags for a PathMetadataEntry Returns *PackageError on failure.
- **`GetPathMetaTagsByType`** - [Getpathmetatagsbytype](api_metadata.md#8172-getpathmetatagsbytype-function)
  - GetPathMetaTagsByType returns all tags of a specific type for a PathMetadataEntry Returns a slice of Tag pointers with the specified type parameter T Only tags matching the type T and corresponding TagValueType are returned Returns *PackageError on failure.
- **`HasFileEntryTag`** - [Hasfileentrytag](api_file_mgmt_file_entry.md#3129-hasfileentrytag-function)
  - HasFileEntryTag checks if a tag with the specified key exists on a FileEntry Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`HasFileEntryTags`** - [Hasfileentrytags](api_file_mgmt_file_entry.md#31210-hasfileentrytags-function)
  - HasFileEntryTags checks if the FileEntry has any tags.
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`HasPathMetaTag`** - [Haspathmetatag](api_metadata.md#8179-haspathmetatag-function)
  - HasPathMetaTag checks if a tag with the specified key exists on a PathMetadataEntry.
- **`NewTag`** - [Newtag](api_file_mgmt_file_entry.md#192-newtag-function)
  - NewTag creates a new type-safe tag with the specified key, value, and type.
- **`RemoveFileEntryTag`** - [Removefileentrytag](api_file_mgmt_file_entry.md#3128-removefileentrytag-function)
  - RemoveFileEntryTag removes a tag by key from a FileEntry.
  - Returns *PackageError on failure.
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`RemovePathMetaTag`** - [Removepathmetatag](api_metadata.md#8178-removepathmetatag-function)
  - RemovePathMetaTag removes a tag by key from a PathMetadataEntry Returns *PackageError on failure.
- **`SanitizeComment`** - [Sanitizecomment](api_metadata.md#142-sanitizecomment-function)
  - SanitizeComment sanitizes comment content to prevent injection attacks Returns *PackageError on failure.
- **`SanitizeSignatureComment`** - [1.5.2 SanitizeSignatureComment Function](api_metadata.md#152-sanitizesignaturecomment-function)
  - SanitizeSignatureComment sanitizes signature comment content Returns *PackageError on failure.
- **`SetDestPath`** - [Setdestpath](api_metadata.md#8217-setdestpath-function)
  - SetDestPathTyped is a generic helper for SetDestPath.
  - This helper exists to allow compile-time type checking of the dest input.
  - It converts dest to DestPathOverride, then delegates to Package.SetDestPath.
  - If dest is a string, it MUST be parsed to determine which destination field to set.
  - If the string is a Windows-only absolute path (drive letter like "C:\\" or "C:/", or UNC path like "\\\\server\\share"), it MUST be stored as DestPathWin.
  - Otherwise, it MUST be stored as DestPath.
- **`SetFileEntryTag`** - [Setfileentrytag](api_file_mgmt_file_entry.md#3127-setfileentrytag-function)
  - SetFileEntryTag updates an existing tag with type safety for a FileEntry.
  - Returns *PackageError if the tag key does not already exist Only modifies existing tags; does not create new tags Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`SetFileEntryTags`** - [Setfileentrytags](api_file_mgmt_file_entry.md#3125-setfileentrytags-function)
  - SetFileEntryTags updates existing tags from a slice of typed tags for a FileEntry.
  - Returns *PackageError if any tag key does not already exist.
  - Only modifies tags that already exist; does not create new tags.
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`SetPathMetaTag`** - [Setpathmetatag](api_metadata.md#8177-setpathmetatag-function)
  - SetPathMetaTag updates an existing tag with type safety for a PathMetadataEntry Returns *PackageError if the tag key does not already exist Only modifies existing tags; does not create new tags.
- **`SetPathMetaTags`** - [Setpathmetatags](api_metadata.md#8174-setpathmetatags-function)
  - SetPathMetaTags updates existing tags from a slice of typed tags for a PathMetadataEntry Returns *PackageError if any tag key does not already exist Only modifies tags that already exist; does not create new tags.
- **`SyncFileEntryTags`** - [Syncfileentrytags](api_file_mgmt_file_entry.md#31211-syncfileentrytags-function)
  - SyncFileEntryTags synchronizes tags with the underlying storage for a FileEntry.
  - Returns *PackageError on failure.
  - Note: This is a standalone function rather than a method due to Go's limitation of not supporting generic methods on non-generic types.
  - See api_generics.md for details.
- **`ValidateComment`** - [Validatecomment](api_metadata.md#141-validatecomment-function)
  - ValidateComment validates comment content for security issues Returns *PackageError on failure.
- **`ValidateCommentEncoding`** - [Validatecommentencoding](api_metadata.md#143-validatecommentencoding-function)
  - ValidateCommentEncoding validates UTF-8 encoding of comment Returns *PackageError on failure.
- **`ValidateSignatureComment`** - [Validatesignaturecomment](api_metadata.md#151-validatesignaturecomment-function)
  - ValidateSignatureComment validates signature comment for security issues Returns *PackageError on failure.

## 6. Compression Types

- **`AdvancedCompressionStrategy`** - [2.1.3 AdvancedCompressionStrategy Interface](api_package_compression.md#213-advancedcompressionstrategy-interface)
  - AdvancedCompressionStrategy for compression with additional validation and metrics.
- **`ByteCompressionStrategy`** - [Bytecompressionstrategy](api_package_compression.md#212-bytecompressionstrategy-interface)
  - ByteCompressionStrategy is the concrete implementation for []byte data.
- **`Compression`** - [3.5 Generic Compression Interface](api_package_compression.md#35-generic-compression-interface)
  - Compression provides type-safe compression for any data type.
- **`CompressionConfig`** - [9.1.1 CompressionConfig Structure](api_package_compression.md#911-compressionconfig-structure)
  - CompressionConfig extends Config for compression-specific settings.
- **`CompressionConfigBuilder`** - [9.1.2.1 CompressionConfigBuilder Struct](api_package_compression.md#9121-compressionconfigbuilder-struct)
  - CompressionConfigBuilder provides fluent configuration building for compression.
- **`CompressionErrorContext`** - [4.3.1.2 CompressionErrorContext Structure](api_package_compression.md#14312-compressionerrorcontext-structure)
  - CompressionErrorContext Define error context types.
- **`CompressionFileOperations`** - [Compressionfileoperations](api_package_compression.md#34-compressionfileoperations-interface)
  - CompressionFileOperations provides file-based compression operations.
- **`CompressionInfo`** - [Compressioninfo](api_package_compression.md#31-compressioninfo-interface)
  - CompressionInfo provides read-only access to compression information.
- **`CompressionJob`** - [2.2.4 CompressionJob Structure](api_package_compression.md#224-compressionjob-structure)
  - CompressionJob represents a unit of work for compression (extends Job).
- **`CompressionOperations`** - [Compressionoperations](api_package_compression.md#32-compressionoperations-interface)
  - CompressionOperations provides basic compression/decompression operations.
- **`CompressionResource`** - [8.4.2 CompressionResource Structure](api_package_compression.md#842-compressionresource-structure)
  - CompressionResource represents a compression-specific resource.
- **`CompressionResourcePool`** - [8.4.1.1 CompressionResourcePool Struct](api_package_compression.md#8411-compressionresourcepool-struct)
  - CompressionResourcePool manages compression-specific resources.
- **`CompressionStrategy`** - [2.1.1 CompressionStrategy Interface Definition](api_package_compression.md#211-compressionstrategy-interface-definition)
  - CompressionStrategy extends Strategy[T, T] for compression operations Both input and output are the same type T The Strategy.Type() method returns "compression" as the category.
- **`CompressionStreaming`** - [Compressionstreaming](api_package_compression.md#33-compressionstreaming-interface)
  - CompressionStreaming provides streaming compression for large packages.
- **`CompressionValidationRule`** - [9.2.2 CompressionValidationRule Structure](api_package_compression.md#922-compressionvalidationrule-structure)
  - CompressionValidationRule represents a compression-specific validation rule.
- **`CompressionValidator`** - [9.2.1.1 CompressionValidator Struct](api_package_compression.md#9211-compressionvalidator-struct)
  - CompressionValidator provides compression-specific validation.
- **`CompressionWorkerPool`** - [8.2.1 CompressionWorkerPool Structure](api_package_compression.md#821-compressionworkerpool-structure)
  - CompressionWorkerPool extends WorkerPool for compression operations.
- **`FileCompressionInfo`** - [4.1 FileCompressionInfo Struct Definition](api_file_mgmt_compression.md#41-filecompressioninfo-struct-definition)
  - FileCompressionInfo contains file compression details.
- **`LZ4Strategy`** - [Lz4strategy](api_package_compression.md#222-lz4strategy-structure)
  - LZ4Strategy LZ4 compression strategy with generic support.
- **`LZMAStrategy`** - [Lzmastrategy](api_package_compression.md#223-lzmastrategy-structure)
  - LZMAStrategy LZMA compression strategy with generic support.
- **`MemoryStrategy`** - [MemoryStrategy](api_package_compression.md#215-memorystrategy-type)
  - MemoryStrategy defines the memory management approach for compression operations.
- **`MemoryErrorContext`** - [Memoryerrorcontext](api_package_compression.md#14314-memoryerrorcontext-structure)
  - MemoryErrorContext provides error context for memory-related compression errors.
- **`PackageCompressionInfo`** - [1.3 PackageCompressionInfo Struct](api_package_compression.md#13-packagecompressioninfo-struct)
  - PackageCompressionInfo contains package compression details.
- **`StreamConfig`** - [Streamconfig](api_package_compression.md#214-streamconfig-structure)
  - StreamConfig handles streaming compression for files of any size.
- **`UnsupportedCompressionErrorContext`** - [4.3.1.3 UnsupportedCompressionErrorContext Structure](api_package_compression.md#14313-unsupportedcompressionerrorcontext-structure)
  - UnsupportedCompressionErrorContext provides error context for unsupported compression type errors.
- **`ZstandardStrategy`** - [Zstandardstrategy](api_package_compression.md#221-zstandardstrategy-structure)
  - ZstandardStrategy Zstandard compression strategy with generic support.

### 6.1 Compression Methods

- **`CompressionConfigBuilder.Build`** - [CompressionConfigBuilder.Build](api_package_compression.md#9127-compressionconfigbuilderbuild-method)
  - Build constructs and returns the final compression configuration.
- **`CompressionConfigBuilder.WithCompressionLevel`** - [CompressionConfigBuilder.WithCompressionLevel](api_package_compression.md#9124-compressionconfigbuilderwithcompressionlevel-method)
  - WithCompressionLevel sets the compression level for the compression configuration builder.
- **`CompressionConfigBuilder.WithCompressionType`** - [CompressionConfigBuilder.WithCompressionType](api_package_compression.md#9123-compressionconfigbuilderwithcompressiontype-method)
  - WithCompressionType sets the compression type for the configuration.
- **`CompressionConfigBuilder.WithMemoryStrategy`** - [CompressionConfigBuilder.WithMemoryStrategy](api_package_compression.md#9126-compressionconfigbuilderwithmemorystrategy-method)
  - WithMemoryStrategy sets the memory strategy for the configuration.
- **`CompressionConfigBuilder.WithSolidCompression`** - [CompressionConfigBuilder.WithSolidCompression](api_package_compression.md#9125-compressionconfigbuilderwithsolidcompression-method)
  - WithSolidCompression enables or disables solid compression for the configuration.
- **`CompressionResourcePool.AcquireCompressionResource`** - [CompressionResourcePool.AcquireCompressionResource](api_package_compression.md#8412-compressionresourcepoolacquirecompressionresource-method)
  - CompressionResourcePool.AcquireCompressionResource Compression-specific resource management methods.
- **`CompressionResourcePool.GetCompressionResourceStats`** - [CompressionResourcePool.GetCompressionResourceStats](api_package_compression.md#8414-compressionresourcepoolgetcompressionresourcestats-method)
  - GetCompressionResourceStats returns statistics about compression resource usage.
- **`CompressionResourcePool.ReleaseCompressionResource`** - [CompressionResourcePool.ReleaseCompressionResource](api_package_compression.md#8413-compressionresourcepoolreleasecompressionresource-method)
  - CompressionResourcePool.ReleaseCompressionResource Returns *PackageError on failure.
- **`CompressionValidator.AddCompressionRule`** - [CompressionValidator.AddCompressionRule](api_package_compression.md#9212-compressionvalidatoraddcompressionrule-method)
  - AddCompressionRule adds a compression validation rule to the validator.
- **`CompressionValidator.ValidateCompressionData`** - [CompressionValidator.ValidateCompressionData](api_package_compression.md#9213-compressionvalidatorvalidatecompressiondata-method)
  - CompressionValidator.ValidateCompressionData Returns *PackageError on failure.
- **`CompressionValidator.ValidateDecompressionData`** - [CompressionValidator.ValidateDecompressionData](api_package_compression.md#9214-compressionvalidatorvalidatedecompressiondata-method)
  - CompressionValidator.ValidateDecompressionData Returns *PackageError on failure.
- **`CompressionWorkerPool.CompressConcurrently`** - [CompressionWorkerPool.CompressConcurrently](api_package_compression.md#822-compressionworkerpooltcompressconcurrently-method)
  - CompressionWorkerPool.CompressConcurrently Compression-specific methods.
- **`CompressionWorkerPool.DecompressConcurrently`** - [CompressionWorkerPool.DecompressConcurrently](api_package_compression.md#823-compressionworkerpooltdecompressconcurrently-method)
  - DecompressConcurrently decompresses multiple data items concurrently using a worker pool.
- **`CompressionWorkerPool.GetCompressionStats`** - [CompressionWorkerPool.GetCompressionStats](api_package_compression.md#824-compressionworkerpooltgetcompressionstats-method)
  - GetCompressionStats returns statistics about compression operations performed by the worker pool.

### 6.2 Compression Helper Functions

- **`NewCompressionConfigBuilder`** - [9.1.2.2 NewCompressionConfigBuilder Function](api_package_compression.md#9122-newcompressionconfigbuilder-function)
  - NewCompressionConfigBuilder creates a new compression configuration builder.

## 7. Encryption and Security Types

- **`AES256GCMFileHandler`** - [4.5.2.1 AES256GCMFileHandler Structure](api_security.md#4521-aes256gcmfilehandler-structure)
  - AES256GCMFileHandler Built-in file encryption handlers.
- **`ByteEncryptionStrategy`** - [Byteencryptionstrategy](api_security.md#412-byteencryptionstrategy-interface)
  - ByteEncryptionStrategy is the concrete implementation for []byte data.
- **`ChaCha20Poly1305FileHandler`** - [Chacha20poly1305filehandler](api_security.md#4522-chacha20poly1305filehandler-structure)
  - ChaCha20Poly1305FileHandler provides file encryption using ChaCha20-Poly1305 algorithm.
- **`EncryptionAlgorithm`** - [3.1.1 EncryptionAlgorithm Type](api_security.md#311-encryptionalgorithm-type)
  - EncryptionAlgorithm enumeration.
- **`EncryptionConfig`** - [4.3.1 EncryptionConfig Structure](api_security.md#431-encryptionconfig-structure)
  - EncryptionConfig provides type-safe encryption configuration.
- **`EncryptionConfigBuilder`** - [4.3.2.1 EncryptionConfigBuilder Struct](api_security.md#4321-encryptionconfigbuilder-struct)
  - EncryptionConfigBuilder provides fluent configuration building for encryption.
- **`EncryptionErrorContext`** - [4.7.1 EncryptionErrorContext Struct](api_security.md#471-encryptionerrorcontext-struct)
  - EncryptionErrorContext provides type-safe error context for encryption operations.
- **`EncryptionKey`** - [4.1.3.1 EncryptionKey Struct Type Definition](api_security.md#4131-encryptionkey-struct-type-definition)
  - EncryptionKey provides type-safe key management Uses Option[T] internally for type-safe key storage.
- **`EncryptionStrategy`** - [4.1.1 EncryptionStrategy Interface Definition](api_security.md#411-encryptionstrategy-interface-definition)
  - EncryptionStrategy extends Strategy[T, T] for encryption operations Both input and output are the same type T The Strategy.Type() method returns "encryption" as the category.
- **`EncryptionType`** - [3.1.2 EncryptionType Alias](api_security.md#312-encryptiontype-alias)
  - EncryptionType is a v1 alias of EncryptionAlgorithm.
  - This exists for compatibility with existing APIs.
- **`EncryptionValidationRule`** - [4.4.2 EncryptionValidationRule Alias](api_security.md#442-encryptionvalidationrule-alias)
  - EncryptionValidationRule is an alias for the generic ValidationRule.
- **`EncryptionValidator`** - [4.4.1 EncryptionValidator Struct](api_security.md#441-encryptionvalidator-struct)
  - EncryptionValidator provides type-safe encryption validation.
- **`FileEncryptionHandler`** - [4.5.1 FileEncryptionHandler Interface](api_security.md#451-fileencryptionhandler-interface)
  - FileEncryptionHandler provides file-specific encryption operations.
- **`MLKEMFileHandler`** - [Mlkemfilehandler](api_security.md#4523-mlkemfilehandler-structure)
  - MLKEMFileHandler provides file encryption using ML-KEM (post-quantum) algorithm.
- **`MLKEMKey`** - [5.1 MLKEMKey Struct](api_security.md#51-mlkemkey-struct)
  - MLKEMKey ML-KEM Key Structure.
- **`SecurityErrorContext`** - [Securityerrorcontext](api_basic_operations.md#203-securityerrorcontext-structure)
  - SecurityErrorContext provides typed context for security-related errors.
  - This context structure is used with structured errors to provide additional diagnostic information for security operations.

### 7.1 Encryption and Security Methods

- **`EncryptionConfigBuilder.Build`** - [EncryptionConfigBuilder.Build](api_security.md#4327-encryptionconfigbuildertbuild-method)
  - Build constructs and returns the final encryption configuration.
- **`EncryptionConfigBuilder.WithAuthenticationTag`** - [EncryptionConfigBuilder.WithAuthenticationTag](api_security.md#4326-encryptionconfigbuildertwithauthenticationtag-method)
  - WithAuthenticationTag enables or disables authentication tag generation for the configuration.
- **`EncryptionConfigBuilder.WithEncryptionType`** - [EncryptionConfigBuilder.WithEncryptionType](api_security.md#4323-encryptionconfigbuildertwithencryptiontype-method)
  - WithEncryptionType sets the encryption type for the configuration.
- **`EncryptionConfigBuilder.WithKeySize`** - [EncryptionConfigBuilder.WithKeySize](api_security.md#4324-encryptionconfigbuildertwithkeysize-method)
  - WithKeySize sets the encryption key size for the encryption configuration builder.
- **`EncryptionConfigBuilder.WithRandomIV`** - [EncryptionConfigBuilder.WithRandomIV](api_security.md#4325-encryptionconfigbuildertwithrandomiv-method)
  - WithRandomIV enables or disables random IV generation for the configuration.
- **`EncryptionKey.GetKey`** - [EncryptionKey.GetKey](api_security.md#4133-encryptionkeytgetkey-method)
  - GetKey returns the encryption key material.
- **`EncryptionKey.IsExpired`** - [EncryptionKey.IsExpired](api_security.md#4136-encryptionkeytisexpired-method)
  - IsExpired returns true if the encryption key has expired.
- **`EncryptionKey.IsValid`** - [EncryptionKey.IsValid](api_security.md#4135-encryptionkeytisvalid-method)
  - IsValid returns true if the encryption key is valid.
- **`EncryptionKey.SetKey`** - [EncryptionKey.SetKey](api_security.md#4134-encryptionkeytsetkey-method)
  - SetKey sets the encryption key material.
- **`EncryptionValidator.AddEncryptionRule`** - [EncryptionValidator.AddEncryptionRule](api_security.md#443-encryptionvalidatortaddencryptionrule-method)
  - AddEncryptionRule adds an encryption validation rule to the validator.
- **`EncryptionValidator.ValidateDecryptionData`** - [EncryptionValidator.ValidateDecryptionData](api_security.md#445-encryptionvalidatortvalidatedecryptiondata-method)
  - EncryptionValidator.ValidateDecryptionData Returns *PackageError on failure.
- **`EncryptionValidator.ValidateEncryptionData`** - [EncryptionValidator.ValidateEncryptionData](api_security.md#444-encryptionvalidatortvalidateencryptiondata-method)
  - EncryptionValidator.ValidateEncryptionData Returns *PackageError on failure.
- **`EncryptionValidator.ValidateEncryptionKey`** - [EncryptionValidator.ValidateEncryptionKey](api_security.md#446-encryptionvalidatortvalidateencryptionkey-method)
  - EncryptionValidator.ValidateEncryptionKey Returns *PackageError on failure.
- **`MLKEMKey.Clear`** - [MLKEMKey.Clear](api_security.md#533-mlkemkeyclear-method)
  - Clear clears sensitive key data from memory.
- **`MLKEMKey.Decrypt`** - [MLKEMKey.Decrypt](api_security.md#522-mlkemkeydecrypt-method)
  - Decrypt decrypts ciphertext using ML-KEM key.
- **`MLKEMKey.Encrypt`** - [MLKEMKey.Encrypt](api_security.md#521-mlkemkeyencrypt-method)
  - Encrypt encrypts plaintext using ML-KEM key.
- **`MLKEMKey.GetLevel`** - [MLKEMKey.GetLevel](api_security.md#532-mlkemkeygetlevel-method)
  - GetLevel returns the security level of the key.
- **`MLKEMKey.GetPublicKey`** - [MLKEMKey.GetPublicKey](api_security.md#531-mlkemkeygetpublickey-method)
  - GetPublicKey returns the public key data.

### 7.2 Encryption and Security Helper Functions

- **`GetEncryptionTypeName`** - [3.2.2 GetEncryptionTypeName Function](api_security.md#322-getencryptiontypename-function)
  - GetEncryptionTypeName returns the human-readable name of the encryption type.
- **`IsValidEncryptionType`** - [3.2.1 IsValidEncryptionType Function](api_security.md#321-isvalidencryptiontype-function)
  - IsValidEncryptionType checks if the encryption type is valid.
- **`NewEncryptionConfigBuilder`** - [4.3.2.2 NewEncryptionConfigBuilder Function](api_security.md#4322-newencryptionconfigbuilder-function)
  - NewEncryptionConfigBuilder creates a new encryption configuration builder.
- **`NewEncryptionKey`** - [4.1.3.2 NewEncryptionKey Function](api_security.md#4132-newencryptionkey-function)
  - NewEncryptionKey creates a new encryption key with the specified type, ID, and key material.

## 8. Signature Types

- **`ByteSignatureStrategy`** - [Bytesignaturestrategy](api_signatures.md#412-bytesignaturestrategy-interface)
  - ByteSignatureStrategy is the concrete implementation for []byte data.
- **`Signature`** - [4.1.4.1 Signature Struct](api_signatures.md#4141-signature-struct)
  - Signature provides type-safe signature data Uses Option[T] internally for type-safe signature data storage.
- **`SignatureConfig`** - [4.3.1 SignatureConfig Structure](api_signatures.md#431-signatureconfig-structure)
  - SignatureConfig provides type-safe signature configuration.
- **`SignatureConfigBuilder`** - [4.3.2.1 SignatureConfigBuilder Struct](api_signatures.md#4321-signatureconfigbuilder-struct)
  - SignatureConfigBuilder provides fluent configuration building for signatures.
- **`SignatureErrorContext`** - [5.3.1 SignatureErrorContext Structure](api_signatures.md#531-signatureerrorcontext-structure)
  - SignatureErrorContext Signature-specific error context types.
- **`SignatureStrategy`** - [4.1.1 SignatureStrategy Interface Definition](api_signatures.md#411-signaturestrategy-interface-definition)
  - SignatureStrategy extends Strategy[T, Signature[T]] for signature operations The Process method from Strategy represents the Sign operation The Strategy.Type() method returns "signature" as the category.
- **`SignatureValidationResult`** - [Signaturevalidationresult](api_security.md#23-signaturevalidationresult-struct)
  - SignatureValidationResult provides information about individual signature validation results.
- **`SignatureValidationRule`** - [Signaturevalidationrule](api_signatures.md#442-signaturevalidationrule-alias)
  - SignatureValidationRule is an alias for the generic ValidationRule.
- **`SignatureValidator`** - [4.4.1 SignatureValidator Struct](api_signatures.md#441-signaturevalidator-struct)
  - SignatureValidator provides type-safe signature validation.
- **`SigningKey`** - [Signingkey](api_signatures.md#4131-signingkey-struct)
  - SigningKey provides type-safe key management for signatures Stores private key material only.
  - public keys are handled separately for verification Uses Option[T] internally for type-safe key storage All private key material must be handled within runtime/secret.Do for security.
- **`UnsupportedErrorContext`** - [Unsupportederrorcontext](api_signatures.md#532-unsupportederrorcontext-structure)
  - UnsupportedErrorContext provides error context for unsupported signature type errors.
- **`ValidationErrorContext`** - [Validationerrorcontext](api_signatures.md#534-validationerrorcontext-structure)
  - ValidationErrorContext provides error context for signature validation errors.

### 8.1 Signature Methods

- **`Signature.GetData`** - [Signature.GetData](api_signatures.md#4143-signaturetgetdata-method)
  - GetData returns the signature data.
- **`Signature.GetSignatureType`** - [Signature.GetSignatureType](api_signatures.md#4146-signaturetgetsignaturetype-method)
  - GetSignatureType returns the type of the signature.
- **`Signature.IsValid`** - [Signature.IsValid](api_signatures.md#4145-signaturetisvalid-method)
  - IsValid returns true if the signature is valid.
- **`Signature.SetData`** - [Signature.SetData](api_signatures.md#4144-signaturetsetdata-method)
  - SetData sets the signature data.
- **`SignatureConfigBuilder.Build`** - [SignatureConfigBuilder.Build](api_signatures.md#4327-signatureconfigbuildertbuild-method)
  - Build constructs and returns the final signature configuration.
- **`SignatureConfigBuilder.WithKeySize`** - [SignatureConfigBuilder.WithKeySize](api_signatures.md#4324-signatureconfigbuildertwithkeysize-method)
  - WithKeySize sets the key size for the configuration.
- **`SignatureConfigBuilder.WithMetadata`** - [SignatureConfigBuilder.WithMetadata](api_signatures.md#4326-signatureconfigbuildertwithmetadata-method)
  - WithMetadata enables or disables metadata inclusion for the configuration.
- **`SignatureConfigBuilder.WithSignatureType`** - [SignatureConfigBuilder.WithSignatureType](api_signatures.md#4323-signatureconfigbuildertwithsignaturetype-method)
  - WithSignatureType sets the signature type for the configuration.
- **`SignatureConfigBuilder.WithTimestamp`** - [SignatureConfigBuilder.WithTimestamp](api_signatures.md#4325-signatureconfigbuildertwithtimestamp-method)
  - WithTimestamp enables or disables timestamp inclusion for the configuration.
- **`SignatureValidator.AddSignatureRule`** - [SignatureValidator.AddSignatureRule](api_signatures.md#443-signaturevalidatortaddsignaturerule-method)
  - AddSignatureRule adds a signature validation rule to the validator.
- **`SignatureValidator.ValidateSignatureData`** - [SignatureValidator.ValidateSignatureData](api_signatures.md#444-signaturevalidatortvalidatesignaturedata-method)
  - SignatureValidator.ValidateSignatureData Returns *PackageError on failure.
- **`SignatureValidator.ValidateSignatureFormat`** - [SignatureValidator.ValidateSignatureFormat](api_signatures.md#446-signaturevalidatortvalidatesignatureformat-method)
  - SignatureValidator.ValidateSignatureFormat Returns *PackageError on failure.
- **`SignatureValidator.ValidateSignatureKey`** - [SignatureValidator.ValidateSignatureKey](api_signatures.md#445-signaturevalidatortvalidatesignaturekey-method)
  - SignatureValidator.ValidateSignatureKey Returns *PackageError on failure.
- **`SigningKey.GetKey`** - [SigningKey.GetKey](api_signatures.md#4133-signingkeytgetkey-method)
  - GetKey returns the signing key material.
- **`SigningKey.IsExpired`** - [SigningKey.IsExpired](api_signatures.md#4136-signingkeytisexpired-method)
  - IsExpired returns true if the signing key has expired.
- **`SigningKey.IsValid`** - [SigningKey.IsValid](api_signatures.md#4135-signingkeytisvalid-method)
  - IsValid returns true if the signing key is valid.
- **`SigningKey.SetKey`** - [SigningKey.SetKey](api_signatures.md#4134-signingkeytsetkey-method)
  - SetKey sets the signing key material.

### 8.2 Signature Helper Functions

- **`NewSignature`** - [4.1.4.2 NewSignature Function](api_signatures.md#4142-newsignature-function)
  - Note: Current implementation is simplified for v1 (signatures deferred to v2) Future v2 implementation: func NewSignature[T any](sigType SignatureType, data []byte) *Signature[T].
- **`NewSignatureConfigBuilder`** - [4.3.2.2 NewSignatureConfigBuilder Function](api_signatures.md#4322-newsignatureconfigbuilder-function)
  - NewSignatureConfigBuilder creates a new signature configuration builder.
- **`NewSigningKey`** - [Newsigningkey](api_signatures.md#4132-newsigningkey-function)
  - NewSigningKey creates a new signing key with the specified type, ID, and key material.

## 9. Streaming and Buffer Types

- **`BufferConfig`** - [2.2.2 BufferConfig Struct](api_streaming.md#222-bufferconfig-struct)
  - BufferConfig configures buffer pool behavior and limits.
- **`BufferPool`** - [2.2.1.1 BufferPool Struct Type Definition](api_streaming.md#2211-bufferpool-struct-type-definition)
  - BufferPool manages buffers of any type.
- **`ChunkMode`** - [Chunkmode](api_streaming.md#3215-chunkmode-type)
  - ChunkMode defines how chunks are processed concurrently.
- **`FileStream`** - [Filestream](api_streaming.md#121-filestream-struct)
  - FileStream provides streaming access to file data with buffering and progress tracking.
- **`StreamingConcurrencyConfig`** - [3.2.1.4 StreamingConcurrencyConfig Structure](api_streaming.md#3214-streamingconcurrencyconfig-structure)
  - StreamingConcurrencyConfig defines streaming-specific concurrency settings.
- **`StreamingConfig`** - [4.2.1.1 StreamingConfig Structure](api_streaming.md#4211-streamingconfig-structure)
  - StreamingConfig extends Config for streaming-specific settings.
- **`StreamingConfigBuilder`** - [4.2.1.2 StreamingConfigBuilder Struct](api_streaming.md#4212-streamingconfigbuilder-struct)
  - StreamingConfigBuilder provides fluent configuration building for streaming.
- **`StreamingConfigDefaults`** - [4.3.1 StreamingConfigDefaults Structure](api_streaming.md#431-streamingconfigdefaults-structure)
  - StreamingConfigDefaults represents default streaming configuration values.
- **`StreamingJob`** - [3.2.1.3 StreamingJob Structure](api_streaming.md#3213-streamingjob-structure)
  - StreamingJob represents a unit of streaming work.
- **`StreamingWorker`** - [3.2.1.2 StreamingWorker Structure](api_streaming.md#3212-streamingworker-structure)
  - StreamingWorker represents a single streaming worker.
- **`StreamingWorkerPool`** - [3.2.1.1 StreamingWorkerPool Structure](api_streaming.md#3211-streamingworkerpool-structure)
  - StreamingWorkerPool manages concurrent streaming workers.

### 9.1 Streaming and Buffer Methods

- **`BufferPool.Get`** - [BufferPool.Get](api_streaming.md#2213-bufferpooltget-method)
  - Get retrieves a buffer of the specified size from the pool.
- **`BufferPool.GetStats`** - [BufferPool.GetStats](api_streaming.md#2314-bufferpooltgetstats-method)
  - GetStats returns statistics about buffer pool usage.
- **`BufferPool.Put`** - [BufferPool.Put](api_streaming.md#2214-bufferpooltput-method)
  - Put returns a buffer to the pool for reuse.
- **`BufferPool.SetMaxTotalSize`** - [BufferPool.SetMaxTotalSize](api_streaming.md#2322-bufferpooltsetmaxtotalsize-method)
  - SetMaxTotalSize sets the maximum total size for all buffers in the pool.
- **`BufferPool.TotalSize`** - [BufferPool.TotalSize](api_streaming.md#2321-bufferpoolttotalsize-method)
  - BufferPool.TotalSize Additional BufferPool methods.
- **`FileStream.Close`** - [FileStream.Close](api_streaming.md#1323-filestreamclose-method)
  - FileStream.Close Returns *PackageError on failure.
- **`FileStream.EstimatedTimeRemaining`** - [FileStream.EstimatedTimeRemaining](api_streaming.md#1336-filestreamestimatedtimeremaining-method)
  - EstimatedTimeRemaining returns an estimate of the time remaining to complete the stream read.
- **`FileStream.GetStats`** - [FileStream.GetStats](api_streaming.md#1331-filestreamgetstats-method)
  - GetStats returns statistics about the stream's read operations.
- **`FileStream.IsClosed`** - [FileStream.IsClosed](api_streaming.md#1334-filestreamisclosed-method)
  - IsClosed returns true if the stream has been closed.
- **`FileStream.Position`** - [FileStream.Position](api_streaming.md#1333-filestreamposition-method)
  - Position returns the current read position in the stream.
- **`FileStream.Progress`** - [FileStream.Progress](api_streaming.md#1335-filestreamprogress-method)
  - Progress returns progress information about the stream read operation.
- **`FileStream.Read`** - [FileStream.Read](api_streaming.md#1341-filestreamread-method)
  - Read reads data from the stream into the provided buffer.
- **`FileStream.ReadAt`** - [FileStream.ReadAt](api_streaming.md#1342-filestreamreadat-method)
  - ReadAt reads data from the stream at the specified offset.
- **`FileStream.ReadChunk`** - [FileStream.ReadChunk](api_streaming.md#1321-filestreamreadchunk-method)
  - FileStream.ReadChunk Returns *PackageError on failure.
- **`FileStream.Seek`** - [FileStream.Seek](api_streaming.md#1322-filestreamseek-method)
  - FileStream.Seek Returns *PackageError on failure.
- **`FileStream.Size`** - [FileStream.Size](api_streaming.md#1332-filestreamsize-method)
  - Size returns the total size of the stream in bytes.
- **`StreamingConfigBuilder.Build`** - [StreamingConfigBuilder.Build](api_streaming.md#4218-streamingconfigbuilderbuild-method)
  - Build constructs and returns the final streaming configuration.
- **`StreamingConfigBuilder.WithChunkProcessingMode`** - [StreamingConfigBuilder.WithChunkProcessingMode](api_streaming.md#4215-streamingconfigbuilderwithchunkprocessingmode-method)
  - WithChunkProcessingMode sets the chunk processing mode for the configuration.
- **`StreamingConfigBuilder.WithMaxStreamsPerWorker`** - [StreamingConfigBuilder.WithMaxStreamsPerWorker](api_streaming.md#4216-streamingconfigbuilderwithmaxstreamsperworker-method)
  - WithMaxStreamsPerWorker sets the maximum number of streams per worker for the configuration.
- **`StreamingConfigBuilder.WithStreamBufferSize`** - [StreamingConfigBuilder.WithStreamBufferSize](api_streaming.md#4214-streamingconfigbuilderwithstreambuffersize-method)
  - WithStreamBufferSize sets the stream buffer size for the configuration.
- **`StreamingConfigBuilder.WithStreamTimeout`** - [StreamingConfigBuilder.WithStreamTimeout](api_streaming.md#4217-streamingconfigbuilderwithstreamtimeout-method)
  - WithStreamTimeout sets the stream timeout for the configuration.
- **`StreamingWorkerPool.GetStreamingStats`** - [StreamingWorkerPool.GetStreamingStats](api_streaming.md#333-streamingworkerpoolgetstreamingstats-method)
  - GetStreamingStats returns current streaming worker pool statistics.
- **`StreamingWorkerPool.Start`** - [StreamingWorkerPool.Start](api_streaming.md#3311-streamingworkerpoolstart-method)
  - Start initializes and starts the streaming worker pool Returns *PackageError on failure.
- **`StreamingWorkerPool.Stop`** - [StreamingWorkerPool.Stop](api_streaming.md#3312-streamingworkerpoolstop-method)
  - Stop gracefully shuts down the streaming worker pool Returns *PackageError on failure.
- **`StreamingWorkerPool.SubmitStreamingJob`** - [StreamingWorkerPool.SubmitStreamingJob](api_streaming.md#3321-streamingworkerpoolsubmitstreamingjob-method)
  - SubmitStreamingJob submits a streaming job to the worker pool Returns *PackageError on failure.

### 9.2 Streaming and Buffer Helper Functions

- **`CreateStreamingConfig`** - [4.3.2 CreateStreamingConfig Function](api_streaming.md#432-createstreamingconfig-function)
  - CreateStreamingConfig creates a streaming configuration with intelligent defaults.
- **`DefaultBufferConfig`** - [2.6.1 DefaultBufferConfig Function](api_streaming.md#261-defaultbufferconfig-function)
  - DefaultBufferConfig returns a buffer configuration with default values.
- **`GetStreamingConfigDefaults`** - [4.3.4 GetStreamingConfigDefaults Function](api_streaming.md#434-getstreamingconfigdefaults-function)
  - GetStreamingConfigDefaults returns default streaming configuration values.
- **`NewBufferPool`** - [2.3.2 NewBufferPool Function](api_generics.md#232-newbufferpool-function)
  - NewBufferPool creates a new buffer pool for the specified type.
- **`NewFileStream`** - [Newfilestream](api_streaming.md#131-newfilestream-function)
  - NewFileStream creates a new file stream with the specified configuration.
- **`NewStreamingConfigBuilder`** - [4.2.1.3 NewStreamingConfigBuilder Function](api_streaming.md#4213-newstreamingconfigbuilder-function)
  - NewStreamingConfigBuilder creates a new streaming configuration builder.
- **`ProcessStreamsConcurrently`** - [Processstreamsconcurrently](api_streaming.md#3322-processstreamsconcurrently-function)
  - ProcessStreamsConcurrently processes multiple streams concurrently.
- **`ValidateStreamingConfig`** - [Validatestreamingconfig](api_streaming.md#433-validatestreamingconfig-function)
  - ValidateStreamingConfig validates streaming configuration settings Returns *PackageError on failure.

## 10. Deduplication Types

This section groups deduplication-related types, methods, and helpers. For details, see [Deduplication API](api_deduplication.md).

### 10.1 Deduplication Methods

This subsection groups deduplication methods (for example, duplicate detection and conversion workflows).

### 10.2 Deduplication Helper Functions

This subsection groups helper functions used by deduplication operations.

## 11. FileType System Types

- **`FileType`** - [3.1 FileType Type](file_type_system.md#31-filetype-type)
  - FileType represents a file type identifier Note: This is the authoritative definition.
  - All other references should link to this document.

### 11.1 FileType System Methods

This subsection groups methods related to file type classification and file type => compression selection logic.

### 11.2 FileType System Helper Functions

- **`DetermineFileType`** - [Determinefiletype](file_type_system.md#411-determinefiletype-function-detection-process)
  - DetermineFileType uses a sophisticated multi-stage detection process to identify file types.
- **`IsAudioFile`** - [Isaudiofile](file_type_system.md#216-isaudiofile-function)
  - IsAudioFile returns true if file type is within audio file range (7000-7999).
- **`IsBinaryFile`** - [Isbinaryfile](file_type_system.md#211-isbinaryfile-function)
  - IsBinaryFile returns true if file type is within binary file range (0-999).
- **`IsConfigFile`** - [Isconfigfile](file_type_system.md#214-isconfigfile-function)
  - IsConfigFile returns true if file type is within config file range (4000-4999).
- **`IsImageFile`** - [Isimagefile](file_type_system.md#215-isimagefile-function)
  - IsImageFile returns true if file type is within image file range (5000-6999).
- **`IsScriptFile`** - [Isscriptfile](file_type_system.md#213-isscriptfile-function)
  - IsScriptFile returns true if file type is within script file range (2000-3999).
- **`IsSpecialFile`** - [Isspecialfile](file_type_system.md#219-isspecialfile-function)
  - IsSpecialFile returns true if file type is within special file range (65000-65535).
- **`IsSystemFile`** - [Issystemfile](file_type_system.md#218-issystemfile-function)
  - IsSystemFile returns true if file type is within system file range (10000-10999).
- **`IsTextFile`** - [Istextfile](file_type_system.md#212-istextfile-function)
  - IsTextFile returns true if file type is within text file range (1000-1999).
- **`IsVideoFile`** - [Isvideofile](file_type_system.md#217-isvideofile-function)
  - IsVideoFile returns true if file type is within video file range (8000-9999).
- **`SelectCompressionType`** - [2.2.1 SelectCompressionType Function](file_type_system.md#221-selectcompressiontype-function)
  - SelectCompressionType selects the appropriate compression algorithm based on file type.
  - Skip compression for already compressed formats: Returns CompressionNone for JPEG, PNG, GIF, MP3, MP4, OGG, FLAC files.
  - Special file handling: Uses IsSpecialFile() to check for special file types.
  - FileTypeSignature: Never compress signature files (returns CompressionNone).
  - Special file handling: FileTypeMetadata, FileTypeManifest, FileTypeIndex: Always compress YAML special files (returns CompressionZstd).
  - Other special files: Default compression (returns CompressionZstd).
    Text-based files: Returns CompressionZstd for text, script, and config files (good compression for text).
    Binary media files: Returns CompressionLZ4 for image, audio, and video files (fast compression for binary data).
    Default: Returns CompressionZstd as default compression method.

## 12. Generic Types

- **`ConcurrencyConfig`** - [1.8.4 ConcurrencyConfig Structure](api_generics.md#184-concurrencyconfig-structure)
  - ConcurrencyConfig defines worker pool and thread safety settings.
- **`Config`** - [1.10.1 Config Structure](api_generics.md#1101-config-structure)
  - Config provides type-safe configuration for any data type.
- **`ConfigBuilder`** - [1.10.2.1 ConfigBuilder Struct](api_generics.md#11021-configbuilder-struct)
  - ConfigBuilder provides fluent configuration building.
- **`Job`** - [1.8.3 Job Structure](api_generics.md#183-job-structure)
  - Job represents a unit of work for concurrent processing.
- **`Option`** - [1.1.1 Option Struct](api_generics.md#111-option-struct)
  - Option provides type-safe optional configuration values.
- **`PathEntry`** - [1.3.1.1 PathEntry Struct](api_generics.md#1311-pathentry-struct)
  - PathEntry represents a minimal file or directory path.
- **`Result`** - [1.2.1 Result Struct](api_generics.md#121-result-struct)
  - Result represents a value that may be an error.
- **`Strategy`** - [1.6 Strategy Interface](api_generics.md#16-strategy-interface)
  - Strategy defines a generic strategy pattern for processing different data types.
- **`ThreadSafetyMode`** - [1.8.5 ThreadSafetyMode Type](api_generics.md#185-threadsafetymode-type)
  - ThreadSafetyMode defines the level of thread safety guarantees.
- **`ValidationRule`** - [1.7.2.1 ValidationRule Struct](api_generics.md#1721-validationrule-struct)
  - ValidationRule represents a single validation rule.
- **`Validator`** - [1.7.1 Validator Interface Definition](api_generics.md#171-validator-interface-definition)
  - Validator defines a generic validation interface for type-safe validation.
- **`Worker`** - [1.8.2 Worker Structure](api_generics.md#182-worker-structure)
  - Worker represents a single worker in a WorkerPool.
- **`WorkerPool`** - [1.8.1 WorkerPool Structure](api_generics.md#181-workerpool-structure)
  - WorkerPool manages concurrent workers for any data type.

### 12.1 Generic Methods

- **`ConfigBuilder.Build`** - [ConfigBuilder.Build](api_generics.md#11027-configbuildertbuild-method)
  - Build constructs and returns the final configuration.
- **`ConfigBuilder.WithChunkSize`** - [ConfigBuilder.WithChunkSize](api_generics.md#11023-configbuildertwithchunksize-method)
  - WithChunkSize sets the chunk size for the configuration.
- **`ConfigBuilder.WithCompressionLevel`** - [ConfigBuilder.WithCompressionLevel](api_generics.md#11025-configbuildertwithcompressionlevel-method)
  - WithCompressionLevel sets the compression level for the configuration.
- **`ConfigBuilder.WithMemoryUsage`** - [ConfigBuilder.WithMemoryUsage](api_generics.md#11024-configbuildertwithmemoryusage-method)
  - WithMemoryUsage sets the memory usage limit for the configuration.
- **`ConfigBuilder.WithStrategy`** - [ConfigBuilder.WithStrategy](api_generics.md#11026-configbuildertwithstrategy-method)
  - WithStrategy sets the processing strategy for the configuration.
- **`Option.Clear`** - [Option.Clear](api_generics.md#116-optiontclear-method)
  - Clear clears the option value.
- **`Option.Get`** - [Option.Get](api_generics.md#113-optiontget-method)
  - Get returns the value and a boolean indicating if the value is set.
- **`Option.GetOrDefault`** - [Option.GetOrDefault](api_generics.md#114-optiontgetordefault-method)
  - GetOrDefault returns the value if set, otherwise returns the default value.
- **`Option.IsSet`** - [Option.IsSet](api_generics.md#115-optiontisset-method)
  - IsSet returns true if the option has a value set.
- **`Option.Set`** - [Option.Set](api_generics.md#112-optiontset-method)
  - Set sets the option value.
- **`PathEntry.GetPath`** - [PathEntry.GetPath](api_generics.md#1316-pathentrygetpath-method)
  - GetPath returns the path string as stored (Unix-style with forward slashes).
- **`PathEntry.GetPathForPlatform`** - [PathEntry.GetPathForPlatform](api_generics.md#1317-pathentrygetpathforplatform-method)
  - GetPathForPlatform returns the path string converted for the specified platform On Windows, converts forward slashes to backslashes On Unix/Linux, returns the path as stored (with forward slashes).
- **`PathEntry.ReadFrom`** - [PathEntry.ReadFrom](api_generics.md#1314-pathentryreadfrom-method)
  - ReadFrom reads a PathEntry from the provided io.Reader Implements io.ReaderFrom interface Returns number of bytes read and any error encountered.
- **`PathEntry.Size`** - [PathEntry.Size](api_generics.md#1313-pathentrysize-method)
  - Size returns the total size of the PathEntry in bytes Formula: 2 (PathLength) + PathLength (Path).
- **`PathEntry.Validate`** - [PathEntry.Validate](api_generics.md#1312-pathentryvalidate-method)
  - Validate performs validation checks on the PathEntry Returns error if PathLength doesn't match Path length, or if Path is empty/invalid.
- **`PathEntry.WriteTo`** - [PathEntry.WriteTo](api_generics.md#1315-pathentrywriteto-method)
  - WriteTo writes a PathEntry to the provided io.Writer Implements io.WriterTo interface Returns number of bytes written and any error encountered.
- **`Result.IsErr`** - [Result.IsErr](api_generics.md#126-resulttiserr-method)
  - IsErr returns true if the Result contains an error.
- **`Result.IsOk`** - [Result.IsOk](api_generics.md#125-resulttisok-method)
  - IsOk returns true if the Result contains a value (no error).
- **`Result.Unwrap`** - [Result.Unwrap](api_generics.md#124-resulttunwrap-method)
  - Unwrap returns the value and error from the Result.
- **`ValidationRule.Validate`** - [ValidationRule.Validate](api_generics.md#1722-validationruletvalidate-method)
  - ValidationRule.Validate Returns *PackageError on failure.
- **`WorkerPool.GetWorkerStats`** - [WorkerPool.GetWorkerStats](api_generics.md#194-workerpooltgetworkerstats-method)
  - GetWorkerStats returns current worker pool statistics.
- **`WorkerPool.Start`** - [WorkerPool.Start](api_generics.md#191-workerpooltstart-method)
  - Start initializes and starts the worker pool Returns *PackageError on failure.
- **`WorkerPool.Stop`** - [WorkerPool.Stop](api_generics.md#192-workerpooltstop-method)
  - Stop gracefully shuts down the worker pool Returns *PackageError on failure.
- **`WorkerPool.SubmitJob`** - [WorkerPool.SubmitJob](api_generics.md#193-workerpooltsubmitjob-method)
  - SubmitJob submits a job to the worker pool Returns *PackageError on failure.

### 12.2 Generic Helper Functions

- **`ComposeValidators`** - [ComposeValidators](api_generics.md#223-composevalidators-function)
  - ComposeValidators creates a validator that runs multiple validators.
- **`Err`** - [Err](api_generics.md#123-err-function)
  - Err creates a failed Result with the given error.
- **`NewConfigBuilder`** - [NewConfigBuilder](api_generics.md#233-newconfigbuilder-function-streaming-configuration)
  - NewConfigBuilder creates a new configuration builder.
- **`Ok`** - [Ok](api_generics.md#122-ok-function)
  - Ok creates a successful Result with the given value.
- **`ProcessConcurrently`** - [ProcessConcurrently](api_generics.md#195-processconcurrently-function)
  - ProcessConcurrently processes multiple items concurrently using a worker pool.
- **`ValidateAll`** - [ValidateAll](api_generics.md#222-validateall-function)
  - ValidateAll validates multiple values using a validator.
- **`ValidateWith`** - [ValidateWith](api_generics.md#221-validatewith-function)
  - ValidateWith validates a single value using a validator.

## 13. Error Types

- **`ErrorType`** - [10.2 ErrorType Types and Categories](api_core.md#102-errortype-types-and-categories)
  - ErrorType categorizes errors for programmatic handling.
- **`IOErrorContext`** - [20.4 IOErrorContext Structure](api_basic_operations.md#204-ioerrorcontext-structure)
  - IOErrorContext provides typed context for I/O-related errors.
- **`PackageError`** - [10.4 PackageError Structure](api_core.md#104-packageerror-structure)
  - PackageError represents a structured error in package operations.
- **`PackageErrorContext`** - [20.2 PackageErrorContext Structure](api_basic_operations.md#202-packageerrorcontext-structure)
  - PackageErrorContext provides typed context for common package operation errors.
- **`PatternErrorContext`** - [1.2.1 Creating File Management Errors](api_file_mgmt_errors.md#121-creating-file-management-errors)
  - PatternErrorContext provides typed context for pattern matching errors.
- **`ReadOnlyErrorContext`** - [11.5 ReadOnlyErrorContext Structure](api_basic_operations.md#115-readonlyerrorcontext-structure)
  - ReadOnlyErrorContext provides typed context for read-only enforcement errors.

### 13.1 Error Methods

- **`PackageError.Is`** - [10.4.3 PackageError.Is Method](api_core.md#1043-packageerroris-method)
  - Is implements error matching for error comparison.
- **`PackageError.Error`** - [10.4.1 PackageError.Error Method](api_core.md#1041-packageerrorerror-method)
  - Error returns the formatted error string.
- **`PackageError.Unwrap`** - [10.4.2 PackageError.Unwrap Method](api_core.md#1042-packageerrorunwrap-method)
  - Unwrap returns the underlying cause error.

### 13.2 Error Helper Functions

- **`AddErrorContext`** - [Adderrorcontext](api_core.md#1055-adderrorcontext-function)
  - AddErrorContext adds type-safe context to errors.
- **`AsPackageError`** - [Aspackageerror](api_core.md#1053-aspackageerror-function)
  - AsPackageError checks if an error is a PackageError and returns it if found.
- **`GetErrorContext`** - [Geterrorcontext](api_core.md#1054-geterrorcontext-function)
  - GetErrorContext retrieves type-safe context from errors.
- **`MapError`** - [Maperror](api_core.md#1056-maperror-function)
  - MapError transforms an error with a generic mapper function.
- **`WrapErrorWithContext`** - [Wraperrorwithcontext](api_core.md#1052-wraperrorwithcontext-function)
  - WrapErrorWithContext wraps an error with type-safe context.

## 14. Other Types

- **`AddFileOptions`** - [2.8 AddFileOptions Struct](api_file_mgmt_addition.md#28-addfileoptions-struct)
  - AddFileOptions configures file addition behavior (path determination, metadata preservation, and processing options).
- **`CreateOptions`** - [7.6 CreateOptions Structure](api_basic_operations.md#76-createoptions-structure)
  - CreateOptions configures package creation behavior (including initial metadata and settings).
- **`DestPathSpec`** - [DestPathSpec](api_file_mgmt_extraction.md#1511-destpathspec-struct)
  - DestPathSpec configures a destination path override for extraction.
- **`ExtractPathOptions`** - [2. ExtractPathOptions Struct](api_file_mgmt_extraction.md#2-extractpathoptions-struct)
  - ExtractPathOptions configures filesystem extraction behavior.
- **`FileIndex`** - [6.1.2 FileIndex Struct](package_file_format.md#612-fileindex-struct)
  - FileIndex represents the file index section of a package.
- **`FileInfo`** - [1.2.4 FileInfo Structure](api_core.md#124-fileinfo-structure)
  - FileInfo provides lightweight file information for listing operations.
- **`IndexEntry`** - [6.1.1 IndexEntry Struct](package_file_format.md#611-indexentry-struct)
  - IndexEntry represents a single file index entry.
- **`PackageConfig`** - [9.1 PackageConfig Structure](api_basic_operations.md#91-packageconfig-structure)
  - PackageConfig provides package-level configuration for path handling behavior.
- **`PathHandling`** - [9.3 PathHandling Type](api_basic_operations.md#93-pathhandling-type)
  - PathHandling specifies how to handle multiple paths pointing to the same content.
- **`RecoveryFileHeader`** - [RecoveryFileHeader](api_writing.md#2721-recoveryfileheader-structure)
  - RecoveryFileHeader contains header information for recovery files used by writing operations.
- **`RemoveDirectoryOptions`** - [4.4 RemoveDirectoryOptions Struct](api_file_mgmt_removal.md#44-removedirectoryoptions-struct)
  - RemoveDirectoryOptions configures directory removal behavior.
- **`SymlinkConvertOptions`** - [1.7 SymlinkConvertOptions Struct](api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
  - SymlinkConvertOptions configures path-to-symlink conversion behavior.
- **`TransformPipeline`** - [2.2 TransformPipeline Structure](api_file_mgmt_transform_pipelines.md#22-transformpipeline-structure)
  - TransformPipeline tracks a multi-stage transformation pipeline for large or multi-step operations.
- **`TransformType`** - [2.4 TransformType Type](api_file_mgmt_transform_pipelines.md#24-transformtype-type)
  - TransformType identifies the type of a transformation stage (compress, encrypt, etc.).
- **`readOnlyPackage`** - [11.3 readOnlyPackage Struct](api_basic_operations.md#113-readonlypackage-struct)
  - readOnlyPackage is a wrapper that enforces read-only behavior for a Package.

### 14.1 Other Type Methods

- **`readOnlyPackage.readOnlyError`** - [readOnlyPackage.readOnlyError](api_basic_operations.md#114-readonlypackagereadonlyerror-method)
  - readOnlyError creates a structured security error for read-only enforcement.

## 15. General Helper Functions

This section groups general-purpose helper functions referenced across the specs.

### 15.1 General Validation Functions

This subsection groups validation helpers (for example, format and input validation utilities).

### 15.2 General Utility Functions

- **`NewFileIndex`** - [NewFileIndex](package_file_format.md#613-newfileindex-function)
  - NewFileIndex creates and returns a new FileIndex with zero values.
