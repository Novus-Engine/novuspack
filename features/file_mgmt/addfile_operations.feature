@domain:file_mgmt @m2 @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @spec(api_file_management.md#2-add-file-operations)
Feature: AddFile Operations

  @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @happy
  Scenario: Addfile operations support unified file addition interface
    Given an open writable package
    And FileSource providing file data
    When AddFile is called with path, source, and options
    Then unified file addition interface is used
    And file is added to package
    And created FileEntry is returned

  @REQ-FILEMGMT-049 @REQ-FILEMGMT-052 @REQ-FILEMGMT-001 @happy
  Scenario: Addfile returns created fileentry with metadata
    Given an open writable package
    And FileSource providing file data
    When AddFile is called
    Then created FileEntry is returned
    And FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums

  @REQ-FILEMGMT-001 @happy
  Scenario: Adding a file updates index and metadata
    Given an open writable package
    And FileSource with file data
    When AddFile is called
    Then package index is updated with new file entry
    And package metadata is updated
    And file count is incremented

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile reads file content from FileSource
    Given an open writable package
    And FileSource providing file data
    When AddFile is called
    Then file content is read from FileSource
    And streaming is used for large files when supported
    And memory is managed efficiently

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile uses AddFileOptions for configuration
    Given an open writable package
    And FileSource providing file data
    And AddFileOptions with compression and encryption settings
    When AddFile is called with options
    Then compression settings are applied
    And encryption settings are applied
    And file processing follows options

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile supports various FileSource implementations
    Given an open writable package
    When AddFile is called with different FileSource types
    Then filesystem files are supported via FilePathSource
    And in-memory data is supported via MemorySource
    And custom sources are supported via FileSource interface

  @REQ-FILEMGMT-001 @happy
  Scenario: AddFile automatically closes FileSource when done
    Given an open writable package
    And FileSource providing file data
    When AddFile completes
    Then FileSource is automatically closed
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
    And FileSource with oversized content
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
