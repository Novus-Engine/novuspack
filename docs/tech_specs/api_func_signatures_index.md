# NovusPack Technical Specifications - API Type and Function Signatures Index

- [0. Overview](#0-overview)
- [1. Basic Operations](#1-basic-operations)
  - [1.1 Package Lifecycle](#11-package-lifecycle)
  - [1.2 Package Operations](#12-package-operations)
  - [1.3 Package State](#13-package-state)
- [2. File Management](#2-file-management)
  - [2.1 Basic File Operations](#21-basic-file-operations)
  - [2.2 Encryption-Aware File Operations](#22-encryption-aware-file-operations)
  - [2.3 File Pattern Operations](#23-file-pattern-operations)
  - [2.4 Deduplication Operations](#24-deduplication-operations)
  - [2.5 File Information and Queries](#25-file-information-and-queries)
  - [2.6 FileEntry Methods](#26-fileentry-methods)
- [3. Package Writing](#3-package-writing)
  - [3.1 Write Operations](#31-write-operations)
  - [3.2 Write Strategy Selection](#32-write-strategy-selection)
- [4. Package Compression](#4-package-compression)
  - [4.1 In-Memory Compression](#41-in-memory-compression)
  - [4.2 File-Based Compression](#42-file-based-compression)
  - [4.3 Compression Status](#43-compression-status)
- [5. Digital Signatures](#5-digital-signatures)
  - [5.1 Signature Management](#51-signature-management)
  - [5.2 Signature Validation](#52-signature-validation)
- [6. Package Metadata](#6-package-metadata)
  - [6.1 Comment Management](#61-comment-management)
  - [6.2 AppID/VendorID Management](#62-appidvendorid-management)
  - [6.3 Package Information Structures](#63-package-information-structures)
  - [6.4 Directory Metadata System](#64-directory-metadata-system)
  - [6.5 Signature Comment Security](#65-signature-comment-security)
- [7. Streaming and Buffer Management](#7-streaming-and-buffer-management)
  - [7.1 File Streaming](#71-file-streaming)
  - [7.2 Buffer Management](#72-buffer-management)
  - [7.3 Additional Streaming Topics](#73-additional-streaming-topics)
- [8. Data Types and Structures](#8-data-types-and-structures)
  - [8.1 Package Information](#81-package-information)
  - [8.2 File Management](#82-file-management)
  - [8.3 Compression](#83-compression)
  - [8.4 Streaming](#84-streaming)
  - [8.5 Package Creation](#85-package-creation)
- [9. Error Types](#9-error-types)
  - [9.1 Basic Operations Errors](#91-basic-operations-errors)
  - [9.2 File Management Errors](#92-file-management-errors)
  - [9.3 Compression Errors](#93-compression-errors)
  - [9.4 Encryption Errors](#94-encryption-errors)
  - [9.5 Signature Errors](#95-signature-errors)
  - [9.6 Security Validation Errors](#96-security-validation-errors)
  - [9.7 Structured Error Utilities](#97-structured-error-utilities)
- [10. Generics](#10-generics)
  - [10.1 Core Generic Types](#101-core-generic-types)
  - [10.2 Generic Patterns](#102-generic-patterns)

## 0. Overview

This document provides a comprehensive index of all NovusPack API functions, types, and structures, with direct links to their detailed documentation.
Use this index to quickly locate specific API elements across the documentation.

## 1. Basic Operations

### 1.1 Package Lifecycle

- **`NewPackage`** - [Basic Operations API - Package Creation](api_basic_operations.md#21-package-constructor)
  - Creates a new empty package with default values
- **`OpenPackage`** - [Basic Operations API - Package Opening](api_basic_operations.md#3-package-opening)
  - Convenience helper to open a package by path
- **`Create`** - [Basic Operations API - Package Creation](api_basic_operations.md#32-create-method)
  - Creates a new package at the specified path
- **`CreateWithOptions`** - [Basic Operations API - Package Creation](api_basic_operations.md#33-create-with-options)
  - Creates a new package with specified options
- **`Open`** - [Basic Operations API - Package Opening](api_basic_operations.md#3-package-opening)
  - Opens an existing package from the specified path
- **`OpenWithValidation`** - [Basic Operations API - Package Opening](api_basic_operations.md#3-package-opening)
  - Opens a package and performs full validation
- **`Close`** - [Basic Operations API - Package Closing](api_basic_operations.md#4-package-closing)
  - Closes the package and releases resources
- **`CloseWithCleanup`** - [Basic Operations API - Package Closing](api_basic_operations.md#4-package-closing)
  - Closes the package and performs cleanup operations

### 1.2 Package Operations

- **`Validate`** - [Basic Operations API - Package Validation](api_basic_operations.md#61-package-validation)
  - Validates package format, structure, and integrity
- **`Defragment`** - [Basic Operations API - Package Defragmentation](api_basic_operations.md#52-package-defragmentation)
  - Optimizes package structure and removes unused space
- **`GetInfo`** - [Basic Operations API - Package Information](api_basic_operations.md#53-package-information)
  - Gets comprehensive package information
- **`ReadHeader`** - [Basic Operations API - Header Inspection](api_basic_operations.md#54-header-inspection)
  - Reads the package header from a reader
- **`NewBuilder`** - [Basic Operations API - Package Creation](api_basic_operations.md#2-package-creation)
  - Creates a new PackageBuilder

### 1.3 Package State

- **`IsOpen`** - [Basic Operations API - Check Package State](api_basic_operations.md#55-check-package-state)
  - Checks if the package is currently open
- **`IsReadOnly`** - [Basic Operations API - Check Package State](api_basic_operations.md#55-check-package-state)
  - Checks if the package is in read-only mode
- **`GetPath`** - [Basic Operations API - Check Package State](api_basic_operations.md#55-check-package-state)
  - Returns the current package file path

## 2. File Management

### 2.1 Basic File Operations

- **`AddFile`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - Adds a file to the package from any data source using a unified interface
- **`FileSource`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - Interface for providing file data from various sources (with context support)
- **`FilePathSource`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - FileSource implementation for filesystem files
- **`MemorySource`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - FileSource implementation for in-memory data
- **`NewFilePathSource`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - Creates a FileSource from a filesystem path
- **`NewMemorySource`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - Creates a FileSource from byte data
- **`AddFileOptions`** - [File Management API - Add File](api_file_management.md#11-add-file)
  - Unified configuration options for all file addition operations
- **`RemoveFile`** - [File Management API - Remove File](api_file_management.md#12-remove-file)
  - Removes a file from the package
- **`ExtractFile`** - [File Management API - Extract File](api_file_management.md#13-extract-file)
  - Extracts file content from the package
- **`CompressFile`** - [File Management API - File Compression Operations](api_file_management.md#14-file-compression-operations)
  - Compresses an existing file in the package
- **`DecompressFile`** - [File Management API - File Compression Operations](api_file_management.md#14-file-compression-operations)
  - Decompresses an existing file in the package
- **`GetFileCompressionInfo`** - [File Management API - File Compression Operations](api_file_management.md#14-file-compression-operations)
  - Gets compression information for a file
- **`FileCompressionInfo`** - [File Management API - File Compression Operations](api_file_management.md#14-file-compression-operations)
  - Contains file compression details
- **`EncryptFile`** - [File Management API - File Encryption Operations](api_file_management.md#15-file-encryption-operations)
  - Encrypts an existing file in the package
- **`DecryptFile`** - [File Management API - File Encryption Operations](api_file_management.md#15-file-encryption-operations)
  - Decrypts an existing file in the package
- **`GetFileEncryptionInfo`** - [File Management API - File Encryption Operations](api_file_management.md#15-file-encryption-operations)
  - Gets encryption information for a file
- **`FileEncryptionInfo`** - [File Management API - File Encryption Operations](api_file_management.md#15-file-encryption-operations)
  - Contains file encryption details

### 2.2 Encryption-Aware File Operations

- **`AddFileWithEncryption`** - [File Management API - Add File with Encryption](api_file_management.md#21-add-file-with-encryption)
  - Adds a file with specific encryption type
- **`GetFileEncryptionType`** - [File Management API - Get File Encryption Type](api_file_management.md#22-get-file-encryption-type)
  - Returns the encryption type for a specific file
- **`GetEncryptedFiles`** - [File Management API - Get Encrypted Files](api_file_management.md#23-get-encrypted-files)
  - Returns a list of all encrypted files in the package
- **`IsValidEncryptionType`** - [Security Validation API - Encryption Type System](api_security.md#32-encryption-type-validation)
  - Validates encryption type values
- **`GetEncryptionTypeName`** - [Security Validation API - Encryption Type System](api_security.md#31-encryption-type-definition)
  - Returns human-readable encryption type name
- **`GenerateMLKEMKey`** - [Security Validation API - ML-KEM Key Management](api_security.md#5-ml-kem-key-management)
  - Generates an ML-KEM key pair

### 2.3 File Pattern Operations

- **`AddFilePattern`** - [File Management API - Add File Pattern](api_file_management.md#31-add-file-pattern)
  - Adds files matching a pattern with configurable options
- **`GetPatterns`** - [File Management API - Add File Pattern](api_file_management.md#31-add-file-pattern)
  - Gets files matching patterns from the package

### 2.4 Deduplication Operations

- **`FindExistingEntryByCRC32`** - [File Management API - File Deduplication](api_file_management.md#61-file-deduplication)
  - Finds existing entry by CRC32 checksum
- **`FindExistingEntryMultiLayer`** - [File Management API - File Deduplication](api_file_management.md#61-file-deduplication)
  - Performs multi-layer deduplication
- **`AddPathToExistingEntry`** - [File Management API - File Deduplication](api_file_management.md#61-file-deduplication)
  - Adds a path to an existing entry

### 2.5 File Information and Queries

- **`FileExists`** - [File Management API - File Existence and Properties](api_file_management.md#71-file-existence-and-properties)
  - Checks if a file exists in the package
  - Returns detailed file entry information for a specific file
- **`ListFiles`** - [File Management API - File Existence and Properties](api_file_management.md#71-file-existence-and-properties)
  - Returns a list of all files in the package
- **`FindEntriesByPathPatterns`** - [File Management API - File Existence and Properties](api_file_management.md#71-file-existence-and-properties)
  - Returns file entries matching path patterns
- **`GetFileByPath`** - [File Management API - File Existence and Properties](api_file_management.md#71-file-existence-and-properties)
  - Gets a file entry by path
- **`GetFileByOffset`** - [File Management API - File Existence and Properties](api_file_management.md#71-file-existence-and-properties)
  - Gets a file entry by offset
- **`GetFileByFileID`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Gets a file entry by its unique FileID
- **`GetFileByHash`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Gets a file entry by content hash
- **`GetFileByChecksum`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Gets a file entry by CRC32 checksum
- **`FindEntriesByTag`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Finds all file entries with a specific tag
- **`FindEntriesByType`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Finds all file entries of a specific type
- **`GetFileCount`** - [File Management API - File Lookup by Metadata](api_file_management.md#92-file-lookup-by-metadata)
  - Gets the total number of files in the package

### 2.6 FileEntry Methods

- **`IsCompressed`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Checks if the file is compressed
- **`HasEncryptionKey`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Checks if the file has an encryption key set
- **`ToBinaryFormat`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Converts the file entry to binary format
- **`ParseFileEntry`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Parses raw data into a FileEntry
- **`NewFileEntry`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Constructs a new FileEntry
- **`LoadFileEntry`** - [File Management API - File Entry Properties](api_file_management.md#81-file-entry-properties)
  - Loads a FileEntry from raw data
- **`SetEncryptionKey`** - [File Management API - File Entry Encryption](api_file_management.md#82-file-entry-encryption)
  - Sets the encryption key for the file
- **`Encrypt`** - [File Management API - File Entry Encryption](api_file_management.md#82-file-entry-encryption)
  - Encrypts data using the file's encryption key
- **`Decrypt`** - [File Management API - File Entry Encryption](api_file_management.md#82-file-entry-encryption)
  - Decrypts data using the file's encryption key
- **`UnsetEncryptionKey`** - [File Management API - File Entry Encryption](api_file_management.md#82-file-entry-encryption)
  - Removes the encryption key from the file
- **`LoadData`** - [File Management API - File Entry Data Management](api_file_management.md#83-file-entry-data-management)
  - Loads the file data into memory
- **`ProcessData`** - [File Management API - File Entry Data Management](api_file_management.md#83-file-entry-data-management)
  - Processes the file data (compression, encryption, etc.)

## 3. Package Writing

### 3.1 Write Operations

- **`Write`** - [Package Writing API - Write Method](api_writing.md#1-write-operations)
  - General write method with compression handling
- **`SafeWrite`** - [Package Writing API - SafeWrite Method](api_writing.md#1-write-operations)
  - Atomic write with temp file strategy
- **`FastWrite`** - [Package Writing API - FastWrite Method](api_writing.md#1-write-operations)
  - In-place updates for existing packages

### 3.2 Write Strategy Selection

- **`SelectWriteStrategy`** - [Package Writing API - Write Strategy Selection](api_writing.md#2-write-strategy-selection)
  - Selects appropriate write strategy based on package state

## 4. Package Compression

### 4.1 In-Memory Compression

- **`CompressPackage`** - [Package Compression API - In-Memory Compression](api_package_compression.md#4-in-memory-compression-methods)
  - Compresses package content in memory
- **`DecompressPackage`** - [Package Compression API - In-Memory Compression](api_package_compression.md#4-in-memory-compression-methods)
  - Decompresses the package in memory

### 4.2 File-Based Compression

- **`CompressPackageFile`** - [Package Compression API - File-Based Compression](api_package_compression.md#6-file-based-compression-methods)
  - Compresses package content and writes to specified path
- **`DecompressPackageFile`** - [Package Compression API - File-Based Compression](api_package_compression.md#6-file-based-compression-methods)
  - Decompresses the package and writes to specified path

### 4.3 Compression Status

- **`GetPackageCompressionInfo`** - [Package Compression API - Compression Status](api_package_compression.md#72-compression-status-methods)
  - Returns package compression information
- **`IsPackageCompressed`** - [Package Compression API - Compression Status](api_package_compression.md#72-compression-status-methods)
  - Checks if the package is compressed
- **`GetPackageCompressionType`** - [Package Compression API - Compression Status](api_package_compression.md#72-compression-status-methods)
  - Returns the package compression type
- **`SetPackageCompressionType`** - [Package Compression API - Compression Status](api_package_compression.md#72-compression-status-methods)
  - Sets the package compression type
- **`CanCompressPackage`** - [Package Compression API - Compression Status](api_package_compression.md#72-compression-status-methods)
  - Checks if package can be compressed

## 5. Digital Signatures

### 5.1 Signature Management

- **`AddSignature`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Adds a new digital signature to the package
- **`RemoveSignature`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Removes signature by index and all later signatures
- **`GetSignatureCount`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Returns the total number of signatures in the package
- **`GetSignature`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Retrieves signature information by index
- **`GetAllSignatures`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Returns all signatures in the package
- **`ClearAllSignatures`** - [Digital Signature API - Signature Management](api_signatures.md#1-signature-management)
  - Removes all signatures (resets package to unsigned state)

### 5.2 Signature Validation

- **`ValidateAllSignatures`** - [Digital Signature API - Signature Validation](api_signatures.md#27-signature-validation)
  - Validates all signatures in order
- **`ValidateSignatureType`** - [Digital Signature API - Signature Validation](api_signatures.md#27-signature-validation)
  - Validates signatures of specific type
- **`ValidateSignatureIndex`** - [Digital Signature API - Signature Validation](api_signatures.md#27-signature-validation)
  - Validates signature by index
- **`ValidateSignatureWithKey`** - [Digital Signature API - Signature Validation](api_signatures.md#27-signature-validation)
  - Validates signature with specific public key
- **`ValidateSignatureChain`** - [Digital Signature API - Signature Validation](api_signatures.md#27-signature-validation)
  - Validates signature chain integrity

## 6. Package Metadata

### 6.1 Comment Management

- **`SetComment`** - [Package Metadata API - Comment Management](api_metadata.md#1-comment-management)
  - Sets the package comment
- **`GetComment`** - [Package Metadata API - Comment Management](api_metadata.md#1-comment-management)
  - Gets the current package comment
- **`ClearComment`** - [Package Metadata API - Comment Management](api_metadata.md#1-comment-management)
  - Clears the package comment
- **`HasComment`** - [Package Metadata API - Comment Management](api_metadata.md#1-comment-management)
  - Checks if the package has a comment
- **`Size`** - [Package Metadata API - PackageComment Methods](api_metadata.md#11-packagecomment-methods)
  - Returns the size of the package comment
- **`WriteTo`** - [Package Metadata API - PackageComment Methods](api_metadata.md#11-packagecomment-methods)
  - Writes the comment to a writer
- **`ReadFrom`** - [Package Metadata API - PackageComment Methods](api_metadata.md#11-packagecomment-methods)
  - Reads the comment from a reader
- **`Validate`** - [Package Metadata API - PackageComment Methods](api_metadata.md#11-packagecomment-methods)
  - Validates the package comment
- **`ValidateComment`** - [Package Metadata API - Comment Security Validation](api_metadata.md#12-comment-security-validation)
  - Validates comment content for security issues
- **`SanitizeComment`** - [Package Metadata API - Comment Security Validation](api_metadata.md#12-comment-security-validation)
  - Sanitizes comment content to prevent injection
- **`ValidateCommentEncoding`** - [Package Metadata API - Comment Security Validation](api_metadata.md#12-comment-security-validation)
  - Validates UTF-8 encoding of comment
- **`CheckCommentLength`** - [Package Metadata API - Comment Security Validation](api_metadata.md#12-comment-security-validation)
  - Validates comment length limits
- **`DetectInjectionPatterns`** - [Package Metadata API - Comment Security Validation](api_metadata.md#12-comment-security-validation)
  - Detects malicious patterns in comment content

### 6.2 AppID/VendorID Management

- **`SetAppID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Sets the application identifier
- **`GetAppID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Gets the current application identifier
- **`ClearAppID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Clears the application identifier
- **`HasAppID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Checks if application identifier is set
- **`SetVendorID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Sets the vendor/platform identifier
- **`GetVendorID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Gets the current vendor identifier
- **`ClearVendorID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Clears the vendor identifier
- **`HasVendorID`** - [Package Metadata API - AppID/VendorID Management](api_metadata.md#2-appidvendorid-management)
  - Checks if vendor identifier is set

### 6.3 Package Information Structures

- **`GetPackageInfo`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Gets comprehensive package information
- **`GetSecurityStatus`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Gets current security status
- **`RefreshPackageInfo`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Refreshes package information cache

### 6.4 Directory Metadata System

- **Directory structures and APIs** - [Package Metadata API - Directory Metadata System](api_metadata.md#8-directory-metadata-system)
  - Structures, file formats, and management methods for directory metadata

### 6.5 Signature Comment Security

- **`ValidateSignatureComment`** - [Package Metadata API - Signature Comment Security](api_metadata.md#13-signature-comment-security)
  - Validates signature comment for security issues
- **`SanitizeSignatureComment`** - [Package Metadata API - Signature Comment Security](api_metadata.md#13-signature-comment-security)
  - Sanitizes signature comment content
- **`CheckSignatureCommentLength`** - [Package Metadata API - Signature Comment Security](api_metadata.md#13-signature-comment-security)
  - Validates signature comment length
- **`AuditSignatureComment`** - [Package Metadata API - Signature Comment Security](api_metadata.md#13-signature-comment-security)
  - Audits signature comment for security logging

## 7. Streaming and Buffer Management

### 7.1 File Streaming

- **`NewFileStream`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - Creates a new file stream for large files
- **`ReadChunk`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - Reads a chunk of data from the stream
- **`Seek`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - Seeks to a specific position in the stream
- **`Close`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - Closes the file stream
- **`GetStats`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - Gets streaming statistics
- **`Size`** - [Streaming API - Stream Information](api_streaming.md#141-stream-information)
  - Returns the total size of the stream
- **`Position`** - [Streaming API - Stream Information](api_streaming.md#141-stream-information)
  - Returns the current position in the stream
- **`IsClosed`** - [Streaming API - Stream Information](api_streaming.md#141-stream-information)
  - Checks if the stream is closed
- **`Progress`** - [Streaming API - Progress Monitoring](api_streaming.md#142-progress-monitoring)
  - Returns detailed progress information
- **`EstimatedTimeRemaining`** - [Streaming API - Progress Monitoring](api_streaming.md#142-progress-monitoring)
  - Estimates time remaining for completion
- **`Read`** - [Streaming API - Standard Go Interfaces](api_streaming.md#143-standard-go-interfaces)
  - Implements io.Reader interface
- **`ReadAt`** - [Streaming API - Standard Go Interfaces](api_streaming.md#143-standard-go-interfaces)
  - Implements io.ReaderAt interface

### 7.2 Buffer Management

- **`NewBufferPool`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Creates a new buffer pool
- **`Get`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Gets a buffer from the pool
- **`Put`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Returns a buffer to the pool
- **`GetStats`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Gets buffer pool statistics
- **`TotalSize`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Returns the total size of all buffers in the pool
- **`SetMaxTotalSize`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Sets the maximum total size for the buffer pool
- **`Close`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Closes the buffer pool

### 7.3 Additional Streaming Topics

- See [Streaming API](api_streaming.md) for progress monitoring, standard interfaces, and advanced configuration
- **`DefaultBufferConfig`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Returns a default BufferConfig
- **`NewStreamingConfigBuilder`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Creates a builder for streaming configuration
- **`CreateStreamingConfig`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Creates a StreamingConfig
- **`ValidateStreamingConfig`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Validates a StreamingConfig
- **`GetStreamingConfigDefaults`** - [Streaming API - Additional BufferPool Methods](api_streaming.md#231-additional-bufferpool-methods)
  - Returns default streaming configuration values

## 8. Data Types and Structures

### 8.1 Package Information

- **`PackageInfo`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Comprehensive package information
- **`SignatureInfo`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Detailed signature information
- **`SecurityStatus`** - [Package Metadata API - Package Information Structures](api_metadata.md#7-package-information-structures)
  - Current security status

### 8.2 File Management

- **`EncryptionType`** - [File Management API - Encryption Type System](api_file_management.md#3-encryption-type-system)
  - Encryption algorithm enumeration
- **`MLKEMKey`** - [File Management API - ML-KEM Key Management](api_file_management.md#4-ml-kem-key-management)
  - ML-KEM key pair structure

### 8.3 Compression

- **`PackageCompressionInfo`** - [Package Compression API - Compression Information Structure](api_package_compression.md#13-compression-information-structure)
  - Package compression information
- **`PackageCompression`** - [Package Compression API - Compression Types](api_package_compression.md#12-compression-types)
  - Compression type constants

### 8.4 Streaming

- **`FileStream`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - File streaming interface
- **`StreamConfig`** - [Streaming API - File Streaming](api_streaming.md#1-file-streaming)
  - File stream configuration
- **`BufferPool`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Buffer pool interface
- **`BufferConfig`** - [Streaming API - Buffer Management](api_streaming.md#2-buffer-management)
  - Buffer pool configuration
  - Progress and statistics structures

### 8.5 Package Creation

- **`CreateOptions`** - [Basic Operations API - Package Creation](api_basic_operations.md#2-package-creation)
  - Package creation options
- **`NPKMagic`** - [Basic Operations API - Package Constants](api_basic_operations.md#11-package-format-constants)
  - Magic number for .npk files
- **`NPKVersion`** - [Basic Operations API - Package Constants](api_basic_operations.md#11-package-format-constants)
  - Current version of the .npk format
- **`HeaderSize`** - [Package File Format - Package Header](package_file_format.md#2-package-header)
  - Fixed size of the package header (authoritative definition)

## 9. Error Types

### 9.1 Basic Operations Errors

- **`ErrFileNotFound`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Package file not found
- **`ErrFileExists`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - File already exists
- **`ErrInvalidPath`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Invalid file path
- **`ErrPermissionDenied`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Permission denied
- **`ErrIOError`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - I/O error
- **`ErrPackageNotOpen`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Package is not open
- **`ErrValidationFailed`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Package validation failed
- **`ErrContextCancelled`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Context cancelled
- **`ErrContextTimeout`** - [Basic Operations API - Error Handling](api_basic_operations.md#6-error-handling)
  - Context timeout

### 9.2 File Management Errors

- **`ErrContentTooLarge`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - File content too large
- **`ErrNoFilesFound`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - No files found matching pattern
- **`ErrUnsupportedEncryption`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Unsupported encryption type
- **`ErrEncryptionFailed`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Encryption failed
- **`ErrDecryptionFailed`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Decryption failed
- **`ErrDecompressionFailed`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Decompression failed
- **`ErrInvalidSecurityLevel`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Invalid security level
- **`ErrKeyGenerationFailed`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Key generation failed
- **`ErrInvalidKey`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Invalid key
- **`ErrPackageReadOnly`** - [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Package is read-only

### 9.3 Compression Errors

- See [Package Compression API - Error Handling](api_package_compression.md#12-error-handling)
  - Compression failed
  - Decompression failed
  - Unsupported compression type
  - Invalid compression parameters

### 9.4 Encryption Errors

- See [File Management API - Error Handling](api_file_management.md#7-error-handling)
  - Encryption failed
  - Decryption failed
  - Invalid encryption key
  - Unsupported encryption type

### 9.5 Signature Errors

- See [Digital Signature API - Error Handling](api_signatures.md#5-error-handling)
  - Signature creation failed
  - Signature validation failed
  - Unsupported signature type
  - Invalid signature data

### 9.6 Security Validation Errors

- See [Security Validation API - Package Validation](api_security.md#1-package-validation)
  - Validation failed
  - Integrity check failed
  - Trust chain invalid

### 9.7 Structured Error Utilities

- **`NewPackageError`** - [Core Package Interface - Structured Error System](api_core.md#11-structured-error-system)
  - Creates a new structured error
- **`WrapError`** - [Core Package Interface - Structured Error System](api_core.md#11-structured-error-system)
  - Wraps an existing error with structured information
- **`IsPackageError`** - [Core Package Interface - Structured Error System](api_core.md#11-structured-error-system)
  - Checks whether an error is a PackageError
- **`GetErrorType`** - [Core Package Interface - Structured Error System](api_core.md#11-structured-error-system)
  - Retrieves the structured error type

## 10. Generics

### 10.1 Core Generic Types

- **Option** - [Generic Types and Patterns - Option](api_generics.md#11-option-type)
  - Optional value container
- **Result** - [Generic Types and Patterns - Result](api_generics.md#12-result-type)
  - Success/error encapsulation
- **Collection** - [Generic Types and Patterns - Collection](api_generics.md#13-collection-interface)
  - Generic collections and iterators
- **Strategy** - [Generic Types and Patterns - Strategy](api_generics.md#15-strategy-interface)
  - Pluggable behavior pattern
- **Validator** - [Generic Types and Patterns - Validator](api_generics.md#16-validator-interface)
  - Validation abstraction

### 10.2 Generic Patterns

- **Collection Operations** - [Generic Function Patterns - Collection Operations](api_generics.md#21-collection-operations)
- **Validation Functions** - [Generic Function Patterns - Validation Functions](api_generics.md#22-validation-functions)
- **Factory Functions** - [Generic Function Patterns - Factory Functions](api_generics.md#23-factory-functions)
