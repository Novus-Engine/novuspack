@domain:file_mgmt @m2 @REQ-FILEMGMT-024 @spec(api_file_management.md#1114-getencryptiontype-returns)
Feature: Get file encryption type

  @happy
  Scenario: GetFileEncryptionType returns encryption type for file
    Given an open package with file that has encryption
    When GetFileEncryptionType is called with file path
    Then encryption type is returned
    And type is valid EncryptionType value

  @error
  Scenario: GetFileEncryptionType fails for non-existent file
    Given an open package
    When GetFileEncryptionType is called with non-existent path
    Then structured validation error is returned
