@domain:file_mgmt @m2 @REQ-FILEMGMT-069 @spec(api_file_mgmt_addition.md#286-encryption-options)
Feature: File Encryption Options

  @REQ-FILEMGMT-069 @happy
  Scenario: Encryption options control encryption behavior
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with Encrypt option is set
    Then encryption behavior is controlled
    And Encrypt option enables or disables encryption
    And EncryptionType option specifies algorithm
    And EncryptionKey option can override EncryptionType

  @REQ-FILEMGMT-069 @happy
  Scenario: Encryption options support encryption type selection
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with EncryptionType is set
    Then encryption type is configured
    And EncryptionType option specifies encryption algorithm
    And default encryption type is EncryptionNone

  @REQ-FILEMGMT-069 @happy
  Scenario: Encryption options support specific encryption key
    Given an open NovusPack package
    And a file to be added
    And a valid context
    And an encryption key is available
    When AddFileOptions with EncryptionKey is set
    Then specific encryption key is used
    And EncryptionKey option overrides EncryptionType
    And key-based encryption is applied

  @REQ-FILEMGMT-069 @happy
  Scenario: Encryption options have default values
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with nil or default values is used
    Then encryption defaults to false
    And encryption type defaults to EncryptionNone
    And encryption key defaults to nil

  @REQ-FILEMGMT-069 @error
  Scenario: Encryption options validate encryption settings
    Given an open NovusPack package
    And a file to be added
    And an invalid encryption type
    And a valid context
    When AddFileOptions with invalid encryption type is used
    Then a structured error is returned
    And error indicates unsupported encryption type
    And error follows structured error format

  @REQ-FILEMGMT-069 @error
  Scenario: Encryption options validate encryption key
    Given an open NovusPack package
    And a file to be added
    And an invalid encryption key
    And a valid context
    When AddFileOptions with invalid encryption key is used
    Then a structured error is returned
    And error indicates invalid encryption key
    And error follows structured error format
