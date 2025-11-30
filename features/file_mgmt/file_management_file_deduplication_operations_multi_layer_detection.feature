@domain:file_mgmt @m2 @REQ-FILEMGMT-195 @REQ-FILEMGMT-196 @REQ-FILEMGMT-197 @REQ-FILEMGMT-198 @REQ-FILEMGMT-199 @REQ-FILEMGMT-200 @REQ-FILEMGMT-201 @spec(api_file_management.md#9-deduplication-operations)
Feature: File Management: File Deduplication Operations (Multi-Layer Detection)

  @REQ-FILEMGMT-196 @REQ-FILEMGMT-200 @happy
  Scenario: FindExistingEntryByCRC32 detects duplicate by checksum
    Given an open package
    And file "original.txt" exists with CRC32 checksum 0xABCD1234
    When FindExistingEntryByCRC32 is called with checksum 0xABCD1234
    Then the file entry for "original.txt" is returned
    And duplicate can be identified

  @REQ-FILEMGMT-197 @happy
  Scenario: FindExistingEntryByCRC32 accepts CRC32 checksum parameter
    Given an open package
    When FindExistingEntryByCRC32 is called with checksum 0x12345678
    Then the operation completes
    And result indicates whether duplicate exists

  @REQ-FILEMGMT-198 @happy
  Scenario: FindExistingEntryMultiLayer performs comprehensive duplicate detection
    Given an open package
    And file "data.bin" exists with size 1024, CRC32 0x1111, and content hash
    When FindExistingEntryMultiLayer is called with matching size, checksum, and content
    Then the file entry for "data.bin" is returned
    And verification confirms exact duplicate

  @REQ-FILEMGMT-198 @happy
  Scenario: FindExistingEntryMultiLayer verifies with multiple layers
    Given an open package
    And candidate file has matching CRC32 but different content hash
    When FindExistingEntryMultiLayer is called
    Then no duplicate is returned
    And false positives are avoided through multi-layer verification

  @REQ-FILEMGMT-199 @happy
  Scenario: AddPathToExistingEntry adds path to duplicate entry
    Given an open package
    And existing file entry for "original.txt"
    When AddPathToExistingEntry is called with entry and new path "copy.txt"
    Then the file entry now has both "original.txt" and "copy.txt" paths
    And content is shared between paths

  @REQ-FILEMGMT-200 @happy
  Scenario: Deduplication searches for existing files with matching checksums
    Given an open package
    And file with CRC32 0xAAAA exists
    When duplicate detection is performed for file with CRC32 0xAAAA
    Then matching file is found
    And duplicate storage is avoided

  @REQ-FILEMGMT-200 @happy
  Scenario: Deduplication performs multi-layer verification
    Given an open package
    And candidate file for deduplication
    When multi-layer verification is performed
    Then CRC32 check occurs first
    Then content hash verification occurs
    And verification confirms exact match

  @REQ-FILEMGMT-200 @happy
  Scenario: Deduplication adds paths to existing entries when duplicates found
    Given an open package
    And existing file entry for "data.bin"
    When duplicate file is detected with same content
    Then new path is added to existing entry
    And storage space is reduced

  @REQ-FILEMGMT-201 @happy
  Scenario: Deduplication supports simple CRC32-based lookup
    Given an open package
    When FindExistingEntryByCRC32 is used for duplicate detection
    Then fast CRC32-based lookup is performed
    And duplicates can be quickly identified

  @REQ-FILEMGMT-201 @happy
  Scenario: Deduplication supports multi-layer verification for accuracy
    Given an open package
    When FindExistingEntryMultiLayer is used for duplicate detection
    Then CRC32 check is performed
    Then content hash verification is performed
    And accurate duplicate detection is achieved

  @REQ-FILEMGMT-197 @error
  Scenario: FindExistingEntryByCRC32 handles non-existent checksums gracefully
    Given an open package
    When FindExistingEntryByCRC32 is called with non-matching checksum
    Then nil is returned
    And no error occurs

  @REQ-FILEMGMT-198 @error
  Scenario: FindExistingEntryMultiLayer handles mismatched parameters
    Given an open package
    When FindExistingEntryMultiLayer is called with mismatched size
    Then no duplicate is returned
    And verification correctly rejects non-duplicates
