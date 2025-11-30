@domain:file_mgmt @m2 @REQ-FILEMGMT-212 @REQ-FILEMGMT-215 @REQ-FILEMGMT-216 @spec(api_file_management.md#923-getfilebychecksum)
Feature: GetFileByChecksum

  @REQ-FILEMGMT-212 @happy
  Scenario: GetFileByChecksum finds file entry by CRC32 checksum
    Given an open NovusPack package
    And a valid context
    And files with CRC32 checksums exist
    When GetFileByChecksum is called with checksum
    Then file entry with matching checksum is returned
    And boolean true is returned if found
    And boolean false is returned if not found

  @REQ-FILEMGMT-215 @happy
  Scenario: GetFileByChecksum returns FileEntry and boolean when found
    Given an open NovusPack package
    And a valid context
    And a file exists with known CRC32 checksum
    When GetFileByChecksum is called with matching checksum
    Then FileEntry with matching checksum is returned
    And boolean true is returned indicating found
    And FileEntry contains complete file information

  @REQ-FILEMGMT-215 @happy
  Scenario: GetFileByChecksum returns nil and false when not found
    Given an open NovusPack package
    And a valid context
    And a non-existent CRC32 checksum
    When GetFileByChecksum is called with non-existent checksum
    Then nil FileEntry is returned
    And boolean false is returned indicating not found
    And no error is returned for not found case

  @REQ-FILEMGMT-216 @happy
  Scenario: GetFileByChecksum supports fast content identification
    Given an open NovusPack package
    And a valid context
    And files with CRC32 checksums exist
    When GetFileByChecksum is called with checksum
    Then file is identified quickly by checksum
    And fast lookup is performed
    And content identification is efficient

  @REQ-FILEMGMT-216 @happy
  Scenario: GetFileByChecksum supports quick duplicate detection
    Given an open NovusPack package
    And a valid context
    And files with matching checksums exist
    When GetFileByChecksum is called with checksum
    Then duplicate files can be detected quickly
    And checksum-based duplicate detection is enabled
    And lightweight duplicate matching is supported

  @REQ-FILEMGMT-216 @happy
  Scenario: GetFileByChecksum supports lightweight file matching
    Given an open NovusPack package
    And a valid context
    And files exist in the package
    When GetFileByChecksum is used for file matching
    Then lightweight file matching is performed
    And CRC32 checksum enables fast matching
    And file matching is efficient
