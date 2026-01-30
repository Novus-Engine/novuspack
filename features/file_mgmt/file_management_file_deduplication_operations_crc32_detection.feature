@domain:file_mgmt @m2 @REQ-FILEMGMT-020 @spec(api_deduplication.md#31-file-deduplication)
Feature: File Management: File Deduplication Operations (CRC32 Detection)

  @happy
  Scenario: FindExistingEntryByCRC32 finds duplicate by checksum
    Given an open NovusPack package with existing files
    When FindExistingEntryByCRC32 is called with checksum
    Then existing entry with matching checksum is returned if found
    And duplicate detection works correctly

  @happy
  Scenario: FindExistingEntryMultiLayer performs comprehensive deduplication
    Given an open NovusPack package with existing files
    When FindExistingEntryMultiLayer is called
    Then deduplication uses multiple hash types
    And deduplication uses content verification
    And duplicate detection is accurate

  @happy
  Scenario: AddPathToExistingEntry links additional path to existing content
    Given an open writable NovusPack package with existing file
    When AddPathToExistingEntry is called with new path
    Then path is added to existing entry
    And content is shared
    And storage is optimized

  @happy
  Scenario: Deduplication prevents duplicate content storage
    Given an open writable NovusPack package
    When duplicate file is added
    Then deduplication detects duplicate
    And content is not duplicated
    And storage is optimized

  @error
  Scenario: Deduplication operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When deduplication operation is called
    Then structured context error is returned
