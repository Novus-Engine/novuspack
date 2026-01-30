# NovusPack Technical Specifications - File Management Error Handling

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Management Error Handling](#1-file-management-error-handling)
  - [1.1 Error Type Categories](#11-error-type-categories)
  - [1.2 Structured Error Examples](#12-structured-error-examples)
    - [1.2.1 Creating File Management Errors](#121-creating-file-management-errors)

---

## 0. Overview

This document specifies file-management-specific error handling guidance and examples.
The core structured error system is defined in [Structured Error System](api_core.md).

### 0.1 Cross-References

- [Structured Error System](api_core.md)
- [Basic Operations API - Error Handling Best Practices](api_basic_operations.md#21-error-handling-best-practices)
- [File Addition API](api_file_mgmt_addition.md)
- [File Extraction API](api_file_mgmt_extraction.md)
- [File Removal API](api_file_mgmt_removal.md)
- [File Update API](api_file_mgmt_updates.md)

## 1. File Management Error Handling

File management operations MUST return structured errors.
Operation-specific error condition lists are defined in each operation specification.

### 1.1 Error Type Categories

The file management API uses structured errors with the following error type categories:

| ErrorType          | Common Scenarios                                                                                                                           |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ |
| ErrTypeValidation  | File not found, file already exists, invalid path, invalid pattern, content too large, no files found, package not open, package read-only |
| ErrTypeIO          | I/O errors during file operations                                                                                                          |
| ErrTypeEncryption  | Unsupported encryption type, encryption/decryption failures, key generation failures, invalid key                                          |
| ErrTypeCompression | Decompression operation failures                                                                                                           |
| ErrTypeSecurity    | Invalid security level                                                                                                                     |
| ErrTypeContext     | Context cancelled, context timeout                                                                                                         |

### 1.2 Structured Error Examples

This section provides examples of structured errors used in file management operations.

#### 1.2.1 Creating File Management Errors

See [FileErrorContext Structure](api_core.md#104-packageerror-structure) for the complete structure definition.

Example usage:

```go
// Example: File not found with typed context
// Note: FileErrorContext is defined in api_core.md
type ExampleFileErrorContext struct {
    Path      string
    Operation string
    Size      int64
}

err := NewPackageError(ErrTypeValidation, "file not found", nil, ExampleFileErrorContext{
    Path:      "/path/to/file",
    Operation: "ExtractPath",
    Size:      0,
})

// Encryption failure with typed context
err := NewPackageError(ErrTypeEncryption, "encryption failed", nil, EncryptionErrorContext{
    Algorithm: "AES-256-GCM",
    KeySize:   256,
    FilePath:  "/path/to/file",
})

// Example:  Pattern matching error with typed context
type PatternErrorContext struct {
    Pattern   string
    Directory string
    Operation string
}

err := NewPackageError(ErrTypeValidation, "no files found", nil, PatternErrorContext{
    Pattern:   "*.txt",
    Directory: "/src",
    Operation: "AddFilePattern",
})
```
