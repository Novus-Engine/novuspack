# NovusPack Technical Specifications - File Management Best Practices

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Path Management](#1-file-path-management)
  - [1.1 Use consistent path formats](#11-use-consistent-path-formats)
  - [1.2 Path Normalization](#12-path-normalization)
  - [1.3 Path Extraction](#13-path-extraction)
  - [1.4 Validate paths before use](#14-validate-paths-before-use)
- [2. Encryption Management](#2-encryption-management)
  - [2.1 Generate keys externally and use with NovusPack](#21-generate-keys-externally-and-use-with-novuspack)
  - [2.2 Secure key management](#22-secure-key-management)
- [3. Performance Considerations](#3-performance-considerations)
  - [3.1 Use patterns for bulk operations](#31-use-patterns-for-bulk-operations)
  - [3.2 Handle large files with streaming](#32-handle-large-files-with-streaming)
  - [3.3 Use appropriate context timeouts](#33-use-appropriate-context-timeouts)
  - [3.4 Symlink Conversion Performance](#34-symlink-conversion-performance)

---

## 0. Overview

This document provides usage guidance and best practices for file management.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Management API Index](api_file_mgmt_index.md)
- [File Addition API](api_file_mgmt_addition.md)
- [File Update API](api_file_mgmt_updates.md)
- [Security Validation API](api_security.md)

## 1. File Path Management

Best practices for managing file paths.

### 1.1. Use Consistent Path Formats

**Path Storage and Display**: Paths have different formats for internal storage versus user display.

```go
// STORAGE FORMAT (internal - with leading slash)
// Good: Stored paths have leading slash
storedPath := "/documents/subfolder/file.txt"

// DISPLAY FORMAT (user-facing - without leading slash)
// Good: Display paths strip leading slash
displayPath := "documents/subfolder/file.txt"  // Unix/Linux display
displayPath := "documents\\subfolder\\file.txt" // Windows display

// Bad: Using stored format for display
displayPath := "/documents/file.txt"  // Wrong: don't show leading slash to users
```

### 1.2 Path Normalization

Input without leading slash `documents/file.txt` is normalized to stored format `/documents/file.txt`.

Filesystem input path mapping for `AddFile` and `AddFilePattern` is specified in:

- [File Addition API: Filesystem Input Path And Stored Path Derivation](api_file_mgmt_addition.md#213-filesystem-input-path-and-stored-path-derivation)
- [File Addition API: AddFileOptions: Path Determination](api_file_mgmt_addition.md#26-addfileoptions-path-determination)

### 1.3 Path Extraction

The stored path `/path/to/file.txt` is extracted to `path/to/file.txt` on Unix and `path\\to\\file.txt` on Windows.

The leading slash is never shown to end users.

### 1.4. Validate Paths Before Use

```go
if !isValidPath(path) {
    return fmt.Errorf("invalid file path: %s", path)
}
```

## 2. Encryption Management

Best practices for managing encrypted files & their keys.

### 2.1. Generate Keys Externally and Use with NovusPack

NovusPack does not provide key generation functions.
Generate encryption keys using appropriate external libraries:

```go
// Example: Generate AES-256-GCM key using crypto/rand
import "crypto/rand"

aesKey := make([]byte, 32) // 256 bits
if _, err := rand.Read(aesKey); err != nil {
    return err
}

// Wrap in NovusPack EncryptionKey structure
key := NewEncryptionKey(EncryptionAES256GCM, "key-id", aesKey)
defer key.Clear()

// Add encrypted file
inputPath := "/path/to/file"
options := &AddFileOptions{
    StoredPath: Option.Some("/data/file.bin"),
}
entry, err := package.AddFileWithEncryption(ctx, inputPath, key, options)

// Example: Generate ML-KEM key using Cloudflare CIRCL
import "github.com/cloudflare/circl/kem/mlkem"

pub, priv, err := mlkem.Scheme768().GenerateKeyPair()
if err != nil {
    return err
}

mlkemKey := NewEncryptionKey(EncryptionMLKEM768, "key-id", &MLKEMKey{
    PublicKey:  pub.Bytes(),
    PrivateKey: priv.Bytes(),
    Level:      3,
})
defer mlkemKey.Clear()

// Add encrypted file with ML-KEM
inputPath = "/path/to/secret.bin"
options = &AddFileOptions{
    StoredPath: Option.Some("/data/secret.bin"),
}
entry, err = package.AddFileWithEncryption(ctx, inputPath, mlkemKey, options)

// Or use AddFile directly with options
options = &AddFileOptions{
    EncryptionKey: Option.Some(mlkemKey),
    Compress:      Option.Some(true),
}
entry, err = package.AddFile(ctx, inputPath, options)
```

### 2.2. Secure Key Management

```go
// Always clear keys after use
defer key.Clear()

// Handle key material with runtime/secret.Do when possible
import "runtime/secret"

inputPath := "/path/to/file"
options := &AddFileOptions{
    StoredPath: Option.Some("/data/file.bin"),
}

secret.Do(func() {
    // Work with key material here
    entry, err := package.AddFileWithEncryption(ctx, inputPath, key, options)
})
```

All key management operations that handle private key material should use Go's `runtime/secret` package (experimental in Go 1.26+) to protect sensitive cryptographic data in memory.
Encryption operations that access keys should be wrapped within `runtime/secret.Do` to ensure that sensitive key material is promptly erased from memory registers and stack frames after use.
This reduces the risk of keys being exposed through memory dumps or side-channel attacks.
The `Clear` method should also use `runtime/secret.Do` when zeroing out sensitive key data to ensure that the cleanup operation itself does not leave traces in memory.

See [ML-KEM Key Structure and Operations](api_security.md) for detailed key handling requirements.

## 3. Performance Considerations

Performance best practices.

### 3.1. Use Patterns for Bulk Operations

Use AddFilePattern for bulk operations instead of individual AddFile calls for better performance.

### 3.2. Handle Large Files with Streaming

For very large files, consider using the streaming API (see [Streaming and Buffer Management](api_streaming.md)).

### 3.3. Use Appropriate Context Timeouts

Use context.WithTimeout for operations that may take a long time to prevent indefinite blocking.

### 3.4 Symlink Conversion Performance

When converting multiple paths to symlinks, use batch operations for better performance:

```go
// Batch conversion with progress reporting
options := &SymlinkConvertOptions{
    PreservePathMetadata: Option.Some(true),
}
converted, symlinks, err := pkg.ConvertAllPathsToSymlinks(ctx, options, func(current, total int) {
    fmt.Printf("Converting %d/%d\n", current, total)
})
```

Batch conversion processes all multi-path entries in a single pass, reducing overhead compared to individual conversions.
