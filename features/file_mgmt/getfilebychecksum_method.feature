@domain:file_mgmt @m2 @REQ-FILEMGMT-288 @spec(api_file_mgmt_queries.md#251-package-getfilebychecksum-method)
Feature: Package GetFileByChecksum method gets a file entry by CRC32 checksum

  @REQ-FILEMGMT-288 @happy
  Scenario: GetFileByChecksum gets file entry by CRC32 checksum
    Given an open NovusPack package with file entries and checksums
    When GetFileByChecksum is called with a CRC32 checksum
    Then a file entry by CRC32 checksum is returned when found
    And the behavior matches the GetFileByChecksum method specification
    And error is returned when no file matches checksum
    And checksum is validated before lookup
