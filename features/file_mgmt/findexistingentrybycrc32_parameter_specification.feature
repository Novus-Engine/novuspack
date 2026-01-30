@domain:file_mgmt @m2 @REQ-FILEMGMT-197 @spec(api_deduplication.md#312-findexistingentrybycrc32-parameters)
Feature: FindExistingEntryByCRC32 Parameter Specification

  @REQ-FILEMGMT-197 @happy
  Scenario: FindExistingEntryByCRC32 parameters include rawChecksum
    Given an open NovusPack package
    And files with CRC32 checksums exist
    And a CRC32 checksum to search for
    When FindExistingEntryByCRC32 is called
    Then rawChecksum parameter is accepted
    And rawChecksum specifies CRC32 checksum value
    And existing file entry with matching checksum is found

  @REQ-FILEMGMT-197 @happy
  Scenario: FindExistingEntryByCRC32 supports simple CRC32-based lookup
    Given an open NovusPack package
    And files with CRC32 checksums exist
    When FindExistingEntryByCRC32 is called with matching checksum
    Then file entry with matching CRC32 is returned
    And CRC32-based lookup is efficient
    And duplicate detection is enabled

  @REQ-FILEMGMT-197 @happy
  Scenario: FindExistingEntryByCRC32 returns nil when no match found
    Given an open NovusPack package
    And a non-matching CRC32 checksum
    When FindExistingEntryByCRC32 is called with non-matching checksum
    Then nil is returned
    And no file entry is found
    And no error is returned for not found case
