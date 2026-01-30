@domain:file_mgmt @m2 @REQ-FILEMGMT-034 @spec(api_file_mgmt_file_entry.md#9-fileentry-encryption)
Feature: Get file encryption information

  @happy
  Scenario: GetFileEncryptionInfo returns encryption information
    Given an open package with encrypted file
    When GetFileEncryptionInfo is called with file path
    Then FileEncryptionInfo is returned
    And encryption type is included
    And key size is included
    And encryption metadata is included

  @happy
  Scenario: GetFileEncryptionInfo returns info for unencrypted file
    Given an open package with unencrypted file
    When GetFileEncryptionInfo is called with file path
    Then FileEncryptionInfo is returned
    And encryption type indicates no encryption

  @error
  Scenario: GetFileEncryptionInfo fails for non-existent file
    Given an open package
    When GetFileEncryptionInfo is called with non-existent path
    Then structured validation error is returned

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: GetFileEncryptionInfo validates path parameter
    Given an open package
    When GetFileEncryptionInfo is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: GetFileEncryptionInfo respects context cancellation
    Given an open package with files
    And a cancelled context
    When GetFileEncryptionInfo is called
    Then structured context error is returned
    And error type is context cancellation
