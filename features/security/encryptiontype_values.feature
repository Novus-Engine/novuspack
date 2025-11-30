@domain:security @m2 @REQ-SEC-083 @spec(api_security.md#312-values)
Feature: EncryptionType Values

  @REQ-SEC-083 @happy
  Scenario: Encryption type values define EncryptionNone
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionNone value is examined
    Then EncryptionNone indicates no encryption (default)
    And value is used for unencrypted files
    And value is iota starting value

  @REQ-SEC-083 @happy
  Scenario: Encryption type values define EncryptionAES256GCM
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionAES256GCM value is examined
    Then EncryptionAES256GCM indicates AES-256-GCM symmetric encryption
    And value provides traditional encryption
    And value is used for compatibility

  @REQ-SEC-083 @happy
  Scenario: Encryption type values define EncryptionMLKEM
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When EncryptionMLKEM value is examined
    Then EncryptionMLKEM indicates ML-KEM post-quantum encryption
    And value provides quantum-safe encryption
    And value is default for new packages

  @REQ-SEC-083 @happy
  Scenario: Encryption type values provide complete enumeration
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When encryption type values are examined
    Then all available encryption algorithms are defined
    And values enable encryption type selection
    And values support per-file encryption configuration
