@domain:file_mgmt @m2 @REQ-FILEMGMT-102 @spec(api_file_mgmt_file_entry.md#83-encryptdecrypt-parameters)
Feature: Encrypt/Decrypt Parameter Specification

  @REQ-FILEMGMT-102 @happy
  Scenario: FileEntry Encrypt parameters include data to encrypt
    Given an open NovusPack package
    And a valid context
    And a FileEntry with encryption key
    And data to encrypt
    When Encrypt is called on FileEntry
    Then data parameter contains plaintext bytes to encrypt
    And encryption uses file's encryption key
    And encrypted data is returned

  @REQ-FILEMGMT-102 @happy
  Scenario: FileEntry Decrypt parameters include data to decrypt
    Given an open NovusPack package
    And a valid context
    And an encrypted FileEntry
    And encrypted data
    When Decrypt is called on FileEntry
    Then data parameter contains ciphertext bytes to decrypt
    And decryption uses file's encryption key
    And decrypted data is returned

  @REQ-FILEMGMT-102 @error
  Scenario: Encrypt handles missing encryption key
    Given an open NovusPack package
    And a valid context
    And a FileEntry without encryption key
    And data to encrypt
    When Encrypt is called on FileEntry without key
    Then a structured error is returned
    And error indicates missing encryption key
    And error follows structured error format

  @REQ-FILEMGMT-102 @error
  Scenario: Decrypt handles missing encryption key
    Given an open NovusPack package
    And a valid context
    And a FileEntry without encryption key
    And encrypted data
    When Decrypt is called on FileEntry without key
    Then a structured error is returned
    And error indicates missing encryption key
    And error follows structured error format

  @REQ-FILEMGMT-102 @error
  Scenario: Encrypt handles encryption failures
    Given an open NovusPack package
    And a valid context
    And a FileEntry with encryption key
    And data that fails to encrypt
    When Encrypt is called and encryption fails
    Then a structured error is returned
    And error indicates encryption failure
    And error follows structured error format

  @REQ-FILEMGMT-102 @error
  Scenario: Decrypt handles decryption failures
    Given an open NovusPack package
    And a valid context
    And an encrypted FileEntry
    And invalid encrypted data
    When Decrypt is called and decryption fails
    Then a structured error is returned
    And error indicates decryption failure
    And error follows structured error format
