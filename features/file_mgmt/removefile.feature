@domain:file_mgmt @m2 @REQ-FILEMGMT-188 @REQ-FILEMGMT-077 @spec(api_file_management.md#3121-removefile)
Feature: RemoveFile

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-077 @happy
  Scenario: Removefile method removes files from package
    Given an open writable package
    And file entry exists
    When RemoveFile is called with file entry reference
    Then file is removed from package
    And file entry is removed from package index
    And removal operation completes successfully

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-077 @happy
  Scenario: RemoveFile removes file by FileEntry identifier
    Given an open writable package
    And file entry with FileID 123 exists
    When RemoveFile is called with the file entry
    Then file with FileID 123 is removed
    And file entry is removed from index
    And package metadata reflects removal

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-080 @happy
  Scenario: RemoveFile updates index and directory state
    Given an open writable package
    And file exists in directory structure
    When RemoveFile is called
    Then package index is updated
    And directory state reflects removal
    And file count is decremented

  @REQ-FILEMGMT-188 @happy
  Scenario: RemoveFile marks file data as deleted
    Given an open writable package
    And file with data exists
    When RemoveFile is called
    Then file data is marked as deleted
    And space is reclaimed during defragmentation
    And file entry is removed immediately

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-078 @happy
  Scenario: RemoveFile accepts context and file identifier parameters
    Given an open writable package
    And valid context
    And file entry reference
    When RemoveFile is called with context and entry
    Then operation uses context for cancellation
    And operation uses entry for file identification
    And operation completes successfully

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-081 @error
  Scenario: RemoveFile handles missing file errors
    Given an open writable package
    And invalid or non-existent file entry
    When RemoveFile is called
    Then ErrFileNotFound error is returned
    And error indicates missing file
    And error follows structured error format

  @REQ-FILEMGMT-188 @REQ-FILEMGMT-081 @error
  Scenario: RemoveFile handles package state errors
    Given a package that is not open
    When RemoveFile is called
    Then ErrPackageNotOpen error is returned
    And error indicates package state issue
    And error follows structured error format

  @REQ-FILEMGMT-188 @error
  Scenario: RemoveFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
