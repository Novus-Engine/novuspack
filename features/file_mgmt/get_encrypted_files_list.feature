@domain:file_mgmt @m2 @REQ-FILEMGMT-025 @spec(api_file_management.md#8-file-encryption-operations)
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
