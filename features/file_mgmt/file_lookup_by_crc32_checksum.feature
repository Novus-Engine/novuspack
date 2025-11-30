@domain:file_mgmt @m2 @REQ-FILEMGMT-212 @REQ-FILEMGMT-213 @REQ-FILEMGMT-214 @REQ-FILEMGMT-215 @REQ-FILEMGMT-216 @spec(api_file_management.md#923-getfilebychecksum)
Feature: File lookup by CRC32 checksum

  @REQ-FILEMGMT-213 @REQ-FILEMGMT-214 @happy
  Scenario: GetFileByChecksum finds file by CRC32 checksum
    Given an open package with file "data.bin"
    And the file has CRC32 checksum 0x12345678
    When GetFileByChecksum is called with checksum 0x12345678
    Then the file entry for "data.bin" is returned
    And the found flag is true

  @REQ-FILEMGMT-213 @REQ-FILEMGMT-214 @happy
  Scenario: GetFileByChecksum returns false when checksum not found
    Given an open package
    When GetFileByChecksum is called with non-existent checksum 0xFFFFFFFF
    Then nil file entry is returned
    And the found flag is false

  @REQ-FILEMGMT-216 @happy
  Scenario: GetFileByChecksum supports fast duplicate detection
    Given an open package with multiple files
    And some files have matching CRC32 checksums
    When GetFileByChecksum is called for a checksum with duplicates
    Then one of the files with matching checksum is returned
    And checksum-based deduplication can be performed

  @REQ-FILEMGMT-216 @happy
  Scenario: GetFileByChecksum works with deduplication operations
    Given an open package
    And FindExistingEntryByCRC32 has found a duplicate
    When GetFileByChecksum is called with the same CRC32
    Then the same file entry is returned
    And deduplication operations are consistent

  @REQ-FILEMGMT-214 @error
  Scenario: GetFileByChecksum respects context cancellation
    Given an open package
    And a cancelled context
    When GetFileByChecksum is called
    Then a structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-215 @happy
  Scenario: GetFileByChecksum returns FileEntry with checksum verification
    Given an open package with file
    When GetFileByChecksum is called
    Then the returned FileEntry contains the matching checksum
    And checksum can be verified against file content
