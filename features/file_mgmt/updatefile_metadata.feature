@domain:file_mgmt @m2 @REQ-FILEMGMT-144 @REQ-FILEMGMT-012 @spec(api_file_management.md#513-returns)
Feature: UpdateFile

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-012 @happy
  Scenario: Updatefile returns updated fileentry with metadata
    Given an open writable package
    And existing file entry
    When UpdateFile is called
    Then updated FileEntry is returned
    And FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-012 @happy
  Scenario: UpdateFile returns FileEntry with complete metadata
    Given an open writable package
    And existing file entry
    And FileSource with new content
    When UpdateFile is called
    Then returned FileEntry includes updated size
    And returned FileEntry includes updated checksums
    And returned FileEntry includes updated timestamps
    And returned FileEntry includes compression status
    And returned FileEntry includes encryption details

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-145 @happy
  Scenario: UpdateFile behavior includes content update and metadata preservation
    Given an open writable package
    And existing file entry
    When UpdateFile updates content
    Then file content is updated from FileSource
    And file entry metadata is updated (size, checksums, timestamps)
    And file path and basic metadata are preserved
    And package metadata and file count are updated

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-145 @happy
  Scenario: UpdateFile automatically closes FileSource when done
    Given an open writable package
    And existing file entry
    And FileSource providing content
    When UpdateFile completes
    Then FileSource is automatically closed
    And resources are released
    And cleanup is performed

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-146 @error
  Scenario: UpdateFile handles file not found errors
    Given an open writable package
    And file entry does not exist or is invalid
    When UpdateFile is called
    Then ErrFileNotFound error is returned
    And error indicates file entry issue
    And error follows structured error format

  @REQ-FILEMGMT-144 @REQ-FILEMGMT-146 @error
  Scenario: UpdateFile handles processing errors
    Given an open writable package
    And existing file entry
    And FileSource that fails during processing
    When UpdateFile encounters processing error
    Then structured processing error is returned
    And error indicates processing failure
    And error follows structured error format

  @REQ-FILEMGMT-144 @error
  Scenario: UpdateFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When UpdateFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
