@domain:file_mgmt @m2 @REQ-FILEMGMT-195 @spec(api_file_management.md#9-deduplication-operations)
Feature: File Deduplication Operations

  @REQ-FILEMGMT-195 @happy
  Scenario: Deduplication operations provide duplicate detection and management
    Given an open NovusPack package
    When deduplication operations are performed
    Then duplicate detection is provided
    And duplicate management is provided
    And storage optimization is achieved

  @REQ-FILEMGMT-195 @happy
  Scenario: FindExistingEntryByCRC32 searches for duplicates by checksum
    Given an open package
    And file "original.txt" exists with CRC32 checksum 0xABCD1234
    When FindExistingEntryByCRC32 is called with checksum 0xABCD1234
    Then the file entry for "original.txt" is returned
    And duplicate can be identified by checksum

  @REQ-FILEMGMT-195 @happy
  Scenario: FindExistingEntryMultiLayer performs comprehensive duplicate detection
    Given an open package
    And file "data.bin" exists with size 1024, CRC32 0x1111, and content hash
    When FindExistingEntryMultiLayer is called with matching size, checksum, and content
    Then the file entry for "data.bin" is returned
    And verification confirms exact duplicate
    And multi-layer verification prevents false positives

  @REQ-FILEMGMT-195 @happy
  Scenario: AddPathToExistingEntry links additional path to existing content
    Given an open package
    And existing file entry for "original.txt"
    When AddPathToExistingEntry is called with entry and new path "copy.txt"
    Then the file entry now has both "original.txt" and "copy.txt" paths
    And content is shared between paths
    And storage space is reduced

  @REQ-FILEMGMT-195 @happy
  Scenario: Deduplication reduces storage by sharing content
    Given an open package
    And duplicate file is detected
    When AddPathToExistingEntry is used to link duplicate
    Then new path is added to existing entry
    And content is not duplicated
    And storage is optimized

  @REQ-FILEMGMT-195 @error
  Scenario: FindExistingEntryByCRC32 handles non-existent checksums gracefully
    Given an open package
    When FindExistingEntryByCRC32 is called with non-matching checksum
    Then nil is returned
    And no error occurs
    And operation completes gracefully

  @REQ-FILEMGMT-195 @error
  Scenario: Deduplication operations respect context cancellation
    Given an open package
    And a cancelled context
    When deduplication operation is called
    Then structured context error is returned
    And error follows structured error format
