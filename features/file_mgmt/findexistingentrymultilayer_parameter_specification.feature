@domain:file_mgmt @m2 @REQ-FILEMGMT-198 @spec(api_file_management.md#913-findexistingentrymultilayer-parameters)
Feature: FindExistingEntryMultiLayer Parameter Specification

  @REQ-FILEMGMT-198 @happy
  Scenario: FindExistingEntryMultiLayer parameters include originalSize, rawChecksum, and content
    Given an open NovusPack package
    And files exist in the package
    And file content to verify
    When FindExistingEntryMultiLayer is called
    Then originalSize parameter specifies original file size
    And rawChecksum parameter specifies CRC32 checksum
    And content parameter contains file content for verification
    And multi-layer verification is performed

  @REQ-FILEMGMT-198 @happy
  Scenario: FindExistingEntryMultiLayer performs multi-layer verification
    Given an open NovusPack package
    And files exist in the package
    And file content with matching size and CRC32
    When FindExistingEntryMultiLayer is called
    Then CRC32 checksum verification is performed
    And content hash verification is performed
    And accurate duplicate detection is ensured

  @REQ-FILEMGMT-198 @happy
  Scenario: FindExistingEntryMultiLayer returns file entry and processed content
    Given an open NovusPack package
    And files exist in the package
    And matching file content
    When FindExistingEntryMultiLayer is called
    Then existing FileEntry is returned if match found
    And processed content is returned
    And both return values are provided
