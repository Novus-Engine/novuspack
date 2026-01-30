@domain:file_mgmt @m2 @REQ-FILEMGMT-025 @spec(api_file_mgmt_queries.md#43-listencryptedfiles) @spec(api_file_mgmt_file_entry.md#9-fileentry-encryption)
Feature: Get encrypted files list

  @happy
  Scenario: GetEncryptedFiles returns all encrypted files
    Given an open package with multiple encrypted files
    When GetEncryptedFiles is called
    Then all encrypted file paths are returned
    And list contains all encrypted files
    And unencrypted files are excluded

  @happy
  Scenario: GetEncryptedFiles returns empty list when no encrypted files
    Given an open package with no encrypted files
    When GetEncryptedFiles is called
    Then empty list is returned
