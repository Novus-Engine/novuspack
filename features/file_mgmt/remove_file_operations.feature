@domain:file_mgmt @m2 @spec(api_file_mgmt_removal.md#4-file-removal-operations)
Feature: Remove File Operations

  @REQ-FILEMGMT-136 @REQ-FILEMGMT-139 @happy
  Scenario: Remove file operations support file removal from packages
    Given an open writable NovusPack package
    And file "documents/data.txt" exists
    When RemoveFile is called with path
    Then file is removed from package
    And file entry is removed from package index
    And directory state or tombstones are updated
    And package metadata and file count are updated

  @REQ-FILEMGMT-139 @happy
  Scenario: RemoveFile removes file by path
    Given an open writable package
    And file exists at path "data.txt"
    When RemoveFile is called with path
    Then file entry is removed from package index
    And file data is marked as deleted
    And package metadata is updated

  @REQ-FILEMGMT-139 @happy
  Scenario: RemoveFile removes file by virtual path
    Given an open writable package
    And file exists at path "documents/data.txt"
    When RemoveFile is called with path
    Then file is located by path
    And file entry is removed
    And directory state is updated

  @REQ-FILEMGMT-139 @happy
  Scenario: File removal preserves package integrity and signatures
    Given a signed writable package
    And file exists in package
    When RemoveFile is called
    Then package integrity is preserved
    And signatures remain valid
    And package structure is maintained

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile returns error when package not open
    Given a package that is not open
    When RemoveFile is called
    Then ErrPackageNotOpen error is returned
    And error follows structured error format

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile returns error when file not found
    Given an open writable package
    And file does not exist
    When RemoveFile is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile returns error when package is read-only
    Given a read-only open package
    When RemoveFile is called
    Then ErrPackageReadOnly error is returned
    And error indicates read-only mode
    And error follows structured error format

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
