@domain:file_mgmt @m2 @REQ-FILEMGMT-022 @spec(api_file_management.md#112-file-entry-encryption)
Feature: FileEntry encryption key management

  @happy
  Scenario: UnsetEncryptionKey removes encryption key
    Given a FileEntry with encryption key
    When UnsetEncryptionKey is called
    Then encryption key is removed
    And HasEncryptionKey returns false
    And encryption cannot be performed

  @happy
  Scenario: UnsetEncryptionKey on entry without key has no effect
    Given a FileEntry without encryption key
    When UnsetEncryptionKey is called
    Then no error occurs
    And file entry state is unchanged

  @error
  Scenario: Encryption operations fail after UnsetEncryptionKey
    Given a FileEntry with encryption key removed
    When Encrypt is called
    Then structured validation error is returned
    And error indicates no encryption key
