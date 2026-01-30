# NovusPack Technical Specifications - File Transformation Pipelines

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Multi-Stage Transformation Pipelines](#1-multi-stage-transformation-pipelines)
  - [1.1 Pipeline Structure](#11-pipeline-structure)
  - [1.2 ProcessingState Transitions](#12-processingstate-transitions)
  - [1.3 Pipeline Execution Model](#13-pipeline-execution-model)
  - [1.4 Disk Space Management](#14-disk-space-management)
  - [1.5 Parallel Stage Execution](#15-parallel-stage-execution)
  - [1.6 Intermediate Stage Cleanup](#16-intermediate-stage-cleanup)
  - [1.7 Temporary File Security](#17-temporary-file-security)
  - [1.8 Error Recovery and Resume](#18-error-recovery-and-resume)
  - [1.9 Configuration Options](#19-configuration-options)
  - [1.10 Example: Multi-Stage Extraction](#110-example-multi-stage-extraction)
  - [1.11 Example: Multi-Stage Addition](#111-example-multi-stage-addition)
  - [1.12 Pipeline Validation](#112-pipeline-validation)
- [2. Pipeline Type Definitions](#2-pipeline-type-definitions)
  - [2.1 FileSource Structure](#21-filesource-structure)
  - [2.2 TransformPipeline Structure](#22-transformpipeline-structure)
  - [2.3 TransformStage Structure](#23-transformstage-structure)
  - [2.4 TransformType Type](#24-transformtype-type)

---

## 0. Overview

This document specifies the multi-stage transformation pipeline system used for large file processing.
It is extracted from the File Management API specification.

### 0.1 Cross-References

- [File Addition API](api_file_mgmt_addition.md)
- [File Extraction API](api_file_mgmt_extraction.md)
- [FileEntry API](api_file_mgmt_file_entry.md)
- [Core Package Interface](api_core.md)

## 1. Multi-Stage Transformation Pipelines

This section documents the multi-stage transformation pipeline system that enables memory-efficient processing of large files through sequential transformation stages.

Multi-stage transformation pipelines enable processing of files too large to fit in memory by breaking complex transformations (compress, encrypt, decrypt, decompress) into sequential stages with intermediate temporary files.

Key Benefits:

- Memory-efficient processing of files exceeding available RAM
- Resume capability from last completed stage after interruption
- Proper tracking and cleanup of all intermediate temporary files
- Clear ProcessingState transitions at each stage

Common Use Cases:

- Large file extraction: Package (compressed+encrypted) => Decrypt => Decompress => Final output
- Large file addition: External file (raw) => Compress => Encrypt => Package
- File updates: Package => Decrypt => Decompress => Recompress => Re-encrypt => Package

### 1.1 Pipeline Structure

A transformation pipeline consists of:

1. **OriginalSource**: Initial data location before any transformations
2. **TransformPipeline**: Ordered sequence of transformation stages
3. **TransformStage**: Individual transformation step (compress, decompress, encrypt, decrypt, verify)
4. **CurrentSource**: Current data location after latest completed stage

Each TransformStage tracks:

- **StageType**: Type of transformation (TransformType enum)
- **InputSource**: Where to read data for this stage
- **OutputSource**: Where stage output is written (nil if not yet executed)
- **Completed**: Whether stage finished successfully
- **Error**: Any error that occurred during stage execution

### 1.2 ProcessingState Transitions

ProcessingState uses a **data-state model** (not workflow model) to track what transformations have been applied:

- **ProcessingStateRaw**: Raw (unprocessed) data
- **ProcessingStateCompressed**: Compressed but not encrypted
- **ProcessingStateEncrypted**: Encrypted but not compressed (rare)
- **ProcessingStateCompressedAndEncrypted**: Both compressed and encrypted

ProcessingState transitions as transformations are applied:

Extraction Example (decrypt => decompress):

1. Start: CompressedAndEncrypted (data in package)
2. After decrypt: Compressed (intermediate temp file)
3. After decompress: Raw (final output)

Addition Example (compress => encrypt):

1. Start: Raw (external file)
2. After compress: Compressed (intermediate temp file)
3. After encrypt: CompressedAndEncrypted (ready for package)

### 1.3 Pipeline Execution Model

The pipeline executes stages in order.
Each stage reads from its InputSource and writes to an OutputSource.
The OutputSource becomes the InputSource for the next stage.

The implementation MUST treat the pipeline as an internal mechanism.
Callers MUST NOT be required to manage intermediate stage files.

When a stage completes successfully, the implementation MUST update:

- CurrentSource to the stage OutputSource.
- ProcessingState to match the data-state model.
- TransformStage.Completed to true for the stage.

If a stage fails, the implementation MUST:

- Record the failure in TransformStage.Error.
- Mark the stage as not completed.
- Clean up any partial stage output.

### 1.4 Disk Space Management

Before starting multi-stage transformations, the system performs pre-flight disk space checks:

- Estimates required space based on file sizes and transformation types
- Checks available disk space in temp directory
- Returns `ErrTypeIO` with descriptive message if insufficient space
- Prevents partial transformations that would fail due to disk space

Space Estimation Guidelines:

- Decrypt: Approximately same size as encrypted data
- Decompress: 2-10x size depending on compression ratio
- Compress: 0.1-0.9x size depending on compression ratio and content
- Pipeline requires space for source + all intermediate stages simultaneously

### 1.5 Parallel Stage Execution

The system supports parallel execution of stages where safe:

Safe for Parallelization:

- Checksum verification alongside decompression (read-only validation)
- Multiple independent file extractions
- Any stages that don't modify shared state

Sequential by Default:

- All stages run sequentially unless explicitly marked safe for parallelism
- Prioritizes correctness over performance
- Parallel execution only when guaranteed non-conflicting

### 1.6 Intermediate Stage Cleanup

The system manages temporary files throughout the pipeline lifecycle:

Temporary files are created in the system temporary directory by default.

During Execution:

- Each completed stage creates a new temporary file
- Previous stage output becomes input for next stage
- All temporary files tracked in TransformPipeline

On Success:

- All intermediate temporary files removed after final stage completes
- Only final output retained
- CurrentSource restored to OriginalSource

On Failure:

- All created temporary files removed regardless of which stage failed
- TransformPipeline tracks all temporary files for cleanup
- Cleanup handles missing files gracefully (no error if already removed)

Manual Cleanup:

Use `CleanupTransformPipeline()` to manually clean up temporary files:

```go
err := fe.CleanupTransformPipeline()
```

### 1.7 Temporary File Security

Temporary files are handled with context-aware security:

General Rule: Temp files for encrypted content should be encrypted on disk when possible.

Decrypt Operations Exception: If user is decrypting content to write to disk, unencrypted temp files in same directory are acceptable (user intent is to decrypt).

Security Requirements:

- Atomic operations where possible
- Secure cleanup with overwrites for sensitive temp files
- Use OS temp directory permissions as baseline security
- Priority 6 implementation: Full temp file encryption for encrypted content

### 1.8 Error Recovery and Resume

The pipeline system supports resume from last completed stage:

Resume Capability:

```go
err := fe.ResumeTransformation(ctx)
```

Resume Behavior:

1. Checks last completed stage (CurrentStage index)
2. Verifies last stage output still exists
3. If output missing, retries that stage
4. If output present, continues from next stage
5. Updates ProcessingState at each stage transition

Error Handling:

- Stage failures store error in TransformStage.Error
- Partial temporary files cleaned up automatically
- TransformPipeline.Completed remains false until all stages succeed
- Can retry failed stage or entire pipeline

### 1.9 Configuration Options

**MaxTransformStages** (default: 10):

- Limits pipeline depth to prevent resource exhaustion
- Typical usage: 2-3 stages (decrypt => decompress)
- Default allows generous headroom for complex scenarios
- Returns `ErrTypeValidation` if pipeline exceeds limit

ValidateProcessingState (default: false):

- When enabled, validates ProcessingState matches actual data state
- Useful for debugging and strict validation scenarios
- Disabled by default for performance
- Returns `ErrTypeValidation` if mismatch detected

Set via AddFileOptions:

```go
options := &AddFileOptions{
    MaxTransformStages: Option.Some(15),
    ValidateProcessingState: Option.Some(true),
}
```

### 1.10 Example: Multi-Stage Extraction

```go
// File in package: compressed + encrypted, 50GB raw
ctx := context.Background()

// Read large file (uses pipeline automatically)
data, err := pkg.ReadFile(ctx, "/large_file.bin")
if err != nil {
    return err
}

// Pipeline executed behind the scenes:
// 1. OriginalSource: Package file at offset X (CompressedAndEncrypted)
// 2. Stage 1 Decrypt: => temp1.dat (Compressed)
// 3. Stage 2 Decompress: => temp2.dat (Raw)
// 4. Cleanup: Remove temp1.dat and temp2.dat
// 5. Return: Raw data from temp2.dat
```

### 1.11 Example: Multi-Stage Addition

```go
// Add large file with compression and encryption, 50GB raw
options := &AddFileOptions{
    Compress: Option.Some(true),
    CompressionType: Option.Some(uint8(1)), // Zstd
    EncryptionKey: Option.Some(encKey),
}

fe, err := pkg.AddFile(ctx, "/path/to/large_file.bin", options)
if err != nil {
    return err
}

// Pipeline executed behind the scenes:
// 1. OriginalSource: External file (Raw)
// 2. Stage 1 Compress: => temp1.dat (Compressed)
// 3. Stage 2 Encrypt: => temp2.dat (CompressedAndEncrypted)
// 4. Write to package from temp2.dat
// 5. Cleanup: Remove temp1.dat and temp2.dat
```

### 1.12 Pipeline Validation

Validate pipeline state and sources:

```go
err := fe.ValidateSources()
if err != nil {
    // Handle validation error
    // Returns *PackageError with ErrTypeValidation
}
```

Validation checks:

- CurrentSource is set when ProcessingState indicates active processing
- TransformPipeline has at least one stage if present
- CurrentSource matches final pipeline stage output when completed
- OriginalSource type is valid (IsPackage or IsExternal)

## 2. Pipeline Type Definitions

This section defines the canonical types used by the multi-stage pipeline model.

### 2.1 FileSource Structure

`FileSource` is a general data-source primitive used by `FileEntry` for both pipeline and non-pipeline cases.
Its canonical definition is in [FileEntry API: FileSource Structure](api_file_mgmt_file_entry.md#16-filesource-structure).

### 2.2 TransformPipeline Structure

```go
// TransformPipeline tracks a multi-stage transformation pipeline
type TransformPipeline struct {
    Stages       []TransformStage  // All transformation stages (in order)
    CurrentStage int               // Index of current stage (-1 if not started)
    Completed    bool              // True if all stages completed successfully
}
```

### 2.3 TransformStage Structure

```go
// TransformStage represents a single transformation stage
type TransformStage struct {
    StageType    TransformType  // Type of transformation
    InputSource  *FileSource    // Input for this stage
    OutputSource *FileSource    // Output from this stage (nil if not yet executed)
    Completed    bool           // True if this stage completed successfully
    Error        error          // Error if this stage failed
}
```

### 2.4 TransformType Type

```go
// TransformType identifies the type of transformation
type TransformType uint8

const (
    TransformTypeNone        TransformType = 0x00  // No transformation
    TransformTypeCompress    TransformType = 0x01  // Compression
    TransformTypeDecompress  TransformType = 0x02  // Decompression
    TransformTypeEncrypt     TransformType = 0x03  // Encryption
    TransformTypeDecrypt     TransformType = 0x04  // Decryption
    TransformTypeVerify      TransformType = 0x05  // Checksum verification
    TransformTypeCustom      TransformType = 0xFF  // Custom transformation
)
```
