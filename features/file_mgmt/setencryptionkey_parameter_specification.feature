@domain:file_mgmt @m2 @REQ-FILEMGMT-101 @spec(api_file_management.md#1122-setencryptionkey-parameters)
Feature: SetEncryptionKey Parameter Specification

  @REQ-FILEMGMT-101 @happy
  Scenario: SetEncryptionKey parameters include encryption key
    Given an open NovusPack package
    And a valid context
    And a FileEntry
    And an encryption key
    When SetEncryptionKey is called on FileEntry
    Then key parameter is accepted
    And encryption key is set for the file entry
    And file entry encryption is configured

  @REQ-FILEMGMT-101 @happy
  Scenario: SetEncryptionKey configures file encryption
    Given an open NovusPack package
    And a valid context
    And a FileEntry
    And a valid encryption key
    When SetEncryptionKey is called
    Then file entry encryption key is set
    And file can be encrypted using the key
    And encryption configuration is complete

  @REQ-FILEMGMT-101 @error
  Scenario: SetEncryptionKey handles invalid encryption key
    Given an open NovusPack package
    And a valid context
    And a FileEntry
    And an invalid encryption key
    When SetEncryptionKey is called with invalid key
    Then a structured error is returned
    And error indicates invalid encryption key
    And error follows structured error format
