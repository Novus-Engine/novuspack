@domain:file_mgmt @m2 @REQ-FILEMGMT-136 @spec(api_file_mgmt_removal.md#41-removefile)
Feature: RemoveFile

  @REQ-FILEMGMT-136 @happy
  Scenario: RemoveFile method removes files from package
    Given an open writable package
    And file exists in package
    When RemoveFile is called
    Then file is removed from package
    And file entry is removed from package index
    And removal operation completes successfully

  @REQ-FILEMGMT-136 @happy
  Scenario: RemoveFile removes file by path
    Given an open writable package
    And file exists at path "data.bin"
    When RemoveFile is called with path
    Then file at path "data.bin" is removed
    And file entry is removed from index
    And package metadata reflects removal

  @REQ-FILEMGMT-139 @happy
  Scenario: RemoveFile updates index and directory state
    Given an open writable package
    And file exists in directory structure
    When RemoveFile is called
    Then package index is updated
    And directory state reflects removal
    And file count is decremented

  @REQ-FILEMGMT-136 @happy
  Scenario: RemoveFile marks file data as deleted
    Given an open writable package
    And file with data exists
    When RemoveFile is called
    Then file data is marked as deleted
    And space is reclaimed during defragmentation
    And file entry is removed immediately

  @REQ-FILEMGMT-137 @happy
  Scenario: RemoveFile accepts context and path parameters
    Given an open writable package
    And valid context
    And file path "data.bin"
    When RemoveFile is called with path
    Then operation uses context for cancellation
    And operation uses path for file identification
    And operation completes successfully

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile handles missing file errors
    Given an open writable package
    And invalid or non-existent file entry
    When RemoveFile is called
    Then ErrFileNotFound error is returned
    And error indicates missing file
    And error follows structured error format

  @REQ-FILEMGMT-140 @error
  Scenario: RemoveFile handles package state errors
    Given a package that is not open
    When RemoveFile is called
    Then ErrPackageNotOpen error is returned
    And error indicates package state issue
    And error follows structured error format

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: RemoveFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
