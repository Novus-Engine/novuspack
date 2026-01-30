# Multi-Stage Transformation Pipeline Requirements

## Pipeline Structure and Management

- REQ-PIPELINE-001: TransformPipeline tracks ordered transformation stages with stage list, current stage index, and completion status. [api_file_management.md#2112-pipeline-structure](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-002: FileSource unifies source location tracking with file handle, path, offset, size, and type flags (IsPackage, IsTempFile, IsExternal). [api_file_management.md#119-fileentry-data-source-fields-runtime-only](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-003: CurrentSource tracks current data location (original, intermediate, or final stage) replacing separate SourceFile/SourceOffset/SourceSize/TempFilePath/IsTempFile fields. [api_file_management.md#1191-currentsource-unified-data-source-tracking](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-004: OriginalSource preserves original data source before transformations for package-based operations, resume, and audit trails. [api_file_management.md#1192-originalsource-transformation-origin-tracking](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-005: TransformPipeline is nil for simple operations not requiring multi-stage processing. [api_file_management.md#1193-transformpipeline-multi-stage-transformation-support](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-006: InitializeTransformPipeline creates new transformation pipeline with ordered stages. [api_file_management.md#146-multi-stage-transformation-pipeline](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-007: GetTransformPipeline returns current transformation pipeline (nil if no pipeline active). [api_file_management.md#146-multi-stage-transformation-pipeline](../tech_specs/api_file_mgmt_file_entry.md)

## Stage Types and Execution

- REQ-PIPELINE-008: TransformStage represents individual transformation with type, input source, output source, completion status, and error tracking. [api_file_management.md#2112-pipeline-structure](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-009: TransformType identifies transformation types: None (0x00), Compress (0x01), Decompress (0x02), Encrypt (0x03), Decrypt (0x04), Verify (0x05), Custom (0xFF). [api_file_management.md#11-fileentry-structure](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-010: ExecuteTransformStage executes specific stage in pipeline by index with context support. [api_file_management.md#146-multi-stage-transformation-pipeline](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-011: Each stage reads from InputSource and writes to OutputSource (temporary file for intermediate stages). [api_file_management.md#2112-pipeline-structure](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-012: CurrentSource is updated to point to OutputSource after each completed stage. [api_file_management.md#2112-pipeline-structure](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-013: ProcessingState transitions at each stage (Raw => Compressed => CompressedAndEncrypted for addition, reverse for extraction). [api_file_mgmt_transform_pipelines.md#12-processingstate-transitions](../tech_specs/api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
- REQ-PIPELINE-014: System automatically determines required transformation sequence based on CompressionType and EncryptionType fields. [api_file_management.md#5142-multi-stage-pipeline-extraction-large-files](../tech_specs/api_file_mgmt_extraction.md)

## Resource Management

- REQ-PIPELINE-015: Disk space pre-flight checks estimate required space based on file sizes and transformation types before starting pipeline. [api_file_management.md#2114-disk-space-management](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-016: System checks available disk space in temp directory and returns ErrTypeIO with descriptive message if insufficient. [api_file_management.md#2114-disk-space-management](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-017: Space estimation follows guidelines: decrypt (same size), decompress (2-10x), compress (0.1-0.9x), pipeline requires source + all intermediate stages simultaneously. [api_file_management.md#2114-disk-space-management](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-018: MaxTransformStages limits pipeline depth (default: 10) to prevent resource exhaustion and memory leaks. [api_file_management.md#2119-configuration-options](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-019: System returns ErrTypeValidation if pipeline exceeds MaxTransformStages limit. [api_file_management.md#2119-configuration-options](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-020: Temporary files created in system temp directory with unique names for isolation. [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)

## Security and Cleanup

- REQ-PIPELINE-021: All intermediate temporary files tracked in TransformPipeline for comprehensive cleanup. [api_file_management.md#2116-intermediate-stage-cleanup](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-022: On success, all intermediate temporary files removed after final stage completes, only final output retained, CurrentSource restored to OriginalSource. [api_file_management.md#2116-intermediate-stage-cleanup](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-023: On failure, all created temporary files removed regardless of which stage failed, cleanup handles missing files gracefully. [api_file_management.md#2116-intermediate-stage-cleanup](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-024: CleanupTransformPipeline method provides manual cleanup of all temporary files in pipeline. [api_file_management.md#146-multi-stage-transformation-pipeline](../tech_specs/api_file_mgmt_file_entry.md)
- REQ-PIPELINE-025: Temporary files for encrypted content use context-aware security: encrypted on disk when possible, with exception for decrypt operations where user intends to decrypt. [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-026: Secure cleanup with overwrites for sensitive temp files, using OS temp directory permissions as baseline security. [api_file_management.md#2117-temporary-file-security](../tech_specs/api_file_mgmt_transform_pipelines.md)

## Configuration and Validation

- REQ-PIPELINE-027: MaxTransformStages configurable via AddFileOptions with default 10, covers typical 3-stage operations with generous headroom. [api_file_management.md#2825-multi-stage-transformation-pipeline-options](../tech_specs/api_file_mgmt_addition.md)
- REQ-PIPELINE-028: ValidateProcessingState option enables ProcessingState validation (default: false) for debugging and strict validation scenarios. [api_file_management.md#2825-multi-stage-transformation-pipeline-options](../tech_specs/api_file_mgmt_addition.md)
- REQ-PIPELINE-029: When ValidateProcessingState enabled, system validates ProcessingState matches actual data state in CurrentSource and returns ErrTypeValidation if mismatch detected. [api_file_management.md#2825-multi-stage-transformation-pipeline-options](../tech_specs/api_file_mgmt_addition.md)
- REQ-PIPELINE-030: ValidateSources method validates CurrentSource, OriginalSource, and pipeline consistency. [api_file_mgmt_transform_pipelines.md#112-pipeline-validation](../tech_specs/api_file_mgmt_transform_pipelines.md#112-pipeline-validation)
- REQ-PIPELINE-031: Validation checks: CurrentSource set for active processing, pipeline has at least one stage if present, CurrentSource matches final stage output when completed, OriginalSource type valid. [api_file_mgmt_transform_pipelines.md#112-pipeline-validation](../tech_specs/api_file_mgmt_transform_pipelines.md#112-pipeline-validation)

## Error Handling and Resume

- REQ-PIPELINE-032: ResumeTransformation resumes pipeline from last completed stage after interruption. [api_file_management.md#2118-error-recovery-and-resume](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-033: Resume checks last completed stage, verifies output exists, retries if missing, continues from next stage if present, updates ProcessingState at each transition. [api_file_management.md#2118-error-recovery-and-resume](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-034: Stage failures store error in TransformStage.Error, partial temporary files cleaned up automatically, TransformPipeline.Completed remains false until all stages succeed. [api_file_management.md#2118-error-recovery-and-resume](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-035: System supports retry of failed stage or entire pipeline after error. [api_file_management.md#2118-error-recovery-and-resume](../tech_specs/api_file_mgmt_transform_pipelines.md)

## Parallel Execution

- REQ-PIPELINE-036: System supports parallel stage execution where safe (checksum verification alongside decompression, multiple independent extractions, non-conflicting operations). [api_file_management.md#2115-parallel-stage-execution](../tech_specs/api_file_mgmt_transform_pipelines.md)
- REQ-PIPELINE-037: All stages run sequentially by default unless explicitly marked safe for parallelism, prioritizing correctness over performance. [api_file_management.md#2115-parallel-stage-execution](../tech_specs/api_file_mgmt_transform_pipelines.md)
