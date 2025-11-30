@domain:security @m2 @REQ-SEC-081 @spec(api_security.md#3-encryption-type-system)
Feature: Encryption Type System

  @REQ-SEC-081 @happy
  Scenario: Encryption type system defines EncryptionNone value
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionNone is examined
    Then EncryptionNone indicates no encryption
    And EncryptionNone is default encryption type
    And EncryptionNone is used for unencrypted files

  @REQ-SEC-081 @happy
  Scenario: Encryption type system defines EncryptionAES256GCM value
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionAES256GCM is examined
    Then EncryptionAES256GCM indicates AES-256-GCM symmetric encryption
    And EncryptionAES256GCM provides traditional encryption
    And EncryptionAES256GCM is used for compatibility

  @REQ-SEC-081 @happy
  Scenario: Encryption type system defines EncryptionMLKEM value
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionMLKEM is examined
    Then EncryptionMLKEM indicates ML-KEM post-quantum encryption
    And EncryptionMLKEM provides quantum-safe encryption
    And EncryptionMLKEM is default for new packages

  @REQ-SEC-081 @happy
  Scenario: Encryption type system provides IsValidEncryptionType validation
    Given an open NovusPack package
    And a valid context
    And encryption type value
    When IsValidEncryptionType is called
    Then valid encryption types return true
    And invalid encryption types return false
    And validation ensures type safety

  @REQ-SEC-081 @happy
  Scenario: Encryption type system provides GetEncryptionTypeName utility
    Given an open NovusPack package
    And a valid context
    And encryption type value
    When GetEncryptionTypeName is called
    Then human-readable name is returned
    And EncryptionNone returns "None"
    And EncryptionAES256GCM returns "AES-256-GCM"
    And EncryptionMLKEM returns "ML-KEM"

  @REQ-SEC-081 @happy
  Scenario: Encryption type system defines available encryption algorithms
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When encryption type system is examined
    Then system defines available encryption algorithms for file encryption
    And system supports traditional and quantum-safe encryption
    And system enables per-file encryption selection
