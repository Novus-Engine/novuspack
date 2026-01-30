@domain:file_mgmt @m2 @REQ-FILEMGMT-012 @spec(api_file_mgmt_updates.md#113-updatefile-returns)
Feature: UpdateFile

  @REQ-FILEMGMT-012 @happy
  Scenario: Updatefile returns updated fileentry with metadata
    Given an open writable package
    And an existing file in the package
    And a filesystem file path
    When UpdateFile is called
    Then updated FileEntry is returned
    And FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums

  @REQ-FILEMGMT-012 @happy
  Scenario: UpdateFile returns FileEntry with complete metadata
    Given an open writable package
    And an existing file in the package
    And a filesystem file path
    When UpdateFile is called
    Then returned FileEntry includes updated size
    And returned FileEntry includes updated checksums
    And returned FileEntry includes updated timestamps
    And returned FileEntry includes compression status
    And returned FileEntry includes encryption details

  @REQ-FILEMGMT-012 @happy
  Scenario: UpdateFile behavior includes content update and metadata preservation
    Given an open writable package
    And an existing file in the package
    And a filesystem file path
    When UpdateFile updates content
    Then file content is updated from filesystem path
    And file entry metadata is updated (size, checksums, timestamps)
    And file path and basic metadata are preserved
    And package metadata and file count are updated

  @REQ-FILEMGMT-012 @happy
  Scenario: UpdateFile releases file handles when done
    Given an open writable package
    And an existing file in the package
    And a filesystem file path
    When UpdateFile completes
    Then file handles are closed
    And resources are released
    And cleanup is performed

  @REQ-FILEMGMT-012 @error
  Scenario: UpdateFile handles file not found errors
    Given an open writable package
    And no matching stored package path exists
    When UpdateFile is called
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-FILEMGMT-012 @error
  Scenario: UpdateFile handles processing errors
    Given an open writable package
    And an existing file in the package
    And a filesystem file path
    When UpdateFile encounters processing error
    Then structured processing error is returned
    And error indicates processing failure
    And error follows structured error format

  @REQ-FILEMGMT-012 @error
  Scenario: UpdateFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When UpdateFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
