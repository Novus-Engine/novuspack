@domain:file_mgmt @m2 @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @spec(api_file_mgmt_addition.md#2-add-file-operations) @spec(api_basic_operations.md#46-package-configuration) @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management) @spec(api_file_mgmt_transform_pipelines.md#19-configuration-options) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions) @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management) @spec(api_file_mgmt_transform_pipelines.md#19-configuration-options) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: AddFile Operations

  @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @happy
  Scenario: Addfile operations support unified file addition interface
    Given an open writable package
    And a filesystem file path
    When AddFile is called with path and options
    Then unified file addition interface is used
    And file is added to package
    And created FileEntry is returned

  @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @happy
  Scenario: Addfile returns created fileentry with metadata
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then created FileEntry is returned
    And FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums

  @REQ-FILEMGMT-001 @happy
  Scenario: Adding a file updates index and metadata
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then package index is updated with new file entry
    And package metadata is updated
    And file count is incremented

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile reads file content from filesystem path
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path
    And streaming is used for large files when needed
    And memory is managed efficiently

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile uses AddFileOptions for configuration
    Given an open writable package
    And a filesystem file path
    And AddFileOptions with compression and encryption settings
    When AddFile is called with options
    Then compression settings are applied
    And encryption settings are applied
    And file processing follows options

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile closes file handles and releases resources
    Given an open writable package
    And a filesystem file path
    When AddFile completes
    Then file handles are closed
    And resources are released
    And cleanup is performed

  @REQ-FILEMGMT-001 @error
  Scenario: AddFile returns error when package not open
    Given a package that is not open
    When AddFile is called
    Then ErrPackageNotOpen error is returned
    And error follows structured error format

  @REQ-FILEMGMT-001 @error
  Scenario: AddFile validates content size limits
    Given an open writable package
    And a file exceeding size limits
    When AddFile is called
    Then structured validation error is returned
    And error indicates size limit exceeded
    And error follows structured error format

  @REQ-FILEMGMT-001 @error
  Scenario: AddFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When AddFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format

  @REQ-FILEMGMT-342 @REQ-PIPELINE-013 @happy
  Scenario: AddFile uses multi-stage pipeline for large files with compression and encryption
    Given an open writable package
    And a large file exceeding memory threshold
    And AddFileOptions with compression and encryption enabled
    When AddFile is called
    Then multi-stage transformation pipeline is created
    And stage 1 compresses data to temporary file
    And stage 2 encrypts data to temporary file
    And ProcessingState transitions from Raw to Compressed to CompressedAndEncrypted
    And file is written to package from final stage
    And all temporary files are cleaned up

  @REQ-FILEMGMT-343 @REQ-PIPELINE-015 @error
  Scenario: AddFile checks disk space before multi-stage pipeline
    Given an open writable package
    And a large file requiring multi-stage pipeline
    And AddFileOptions with compression and encryption
    And insufficient disk space for temporary files
    When AddFile is called
    Then disk space error is returned
    And error indicates insufficient space
    And no temporary files are created

  @REQ-FILEMGMT-347 @REQ-PIPELINE-018 @happy
  Scenario: AddFile respects MaxTransformStages configuration
    Given an open writable package
    And AddFileOptions with MaxTransformStages set to 5
    And a file requiring 3 transformation stages
    When AddFile is called
    Then pipeline is created successfully
    And all 3 stages execute
    And file is added to package

  @REQ-FILEMGMT-360 @REQ-FILEMGMT-361 @happy
  Scenario: Add duplicate file with PathHandlingSymlinks
    Given an open writable package
    And a file entry already exists with content "test data"
    When AddFile is called with duplicate content and PathHandling set to PathHandlingSymlinks
    Then a symlink should be created pointing to the original file entry
    And the original file entry should have PathCount = 1

  @REQ-FILEMGMT-360 @happy
  Scenario: Add duplicate file with PathHandlingHardLinks
    Given an open writable package
    And a file entry already exists with content "test data"
    When AddFile is called with duplicate content and PathHandling set to PathHandlingHardLinks
    Then the path should be added to the existing file entry
    And the file entry should have PathCount = 2

  @REQ-FILEMGMT-360 @happy
  Scenario: Add duplicate file with PathHandlingPreserve
    Given an open writable package
    And a file entry already exists with content "test data"
    When AddFile is called with duplicate content and PathHandling set to PathHandlingPreserve
    Then the system should detect and preserve original filesystem behavior
