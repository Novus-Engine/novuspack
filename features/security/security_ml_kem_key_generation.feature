@domain:security @m2 @REQ-SEC-006 @REQ-SEC-007 @REQ-SEC-008 @REQ-SEC-009 @REQ-SEC-010 @REQ-SEC-011 @spec(api_security.md#5-ml-kem-key-management)
Feature: ML-KEM key generation

  @happy
  Scenario: GenerateMLKEMKey generates ML-KEM key pair
    When GenerateMLKEMKey is called
    Then MLKEMKey structure is returned
    And public key is included
    And private key is included
    And key pair is valid

  @happy
  Scenario: Generated ML-KEM keys can be used for encryption
    Given a generated MLKEMKey pair
    When keys are used for encryption
    Then encryption succeeds with public key
    And decryption succeeds with private key
    And key pair is functional

  @error
  Scenario: GenerateMLKEMKey fails if key generation is unavailable
    Given system where ML-KEM is not available
    When GenerateMLKEMKey is called
    Then structured error is returned
    And error indicates key generation failure

  @REQ-SEC-007 @REQ-SEC-008 @error
  Scenario: Security operations validate encryption type parameter
    Given security operation context
    When IsValidEncryptionType or GetEncryptionTypeName is called with invalid type
    Then structured validation error is returned
    And error indicates invalid encryption type

  @REQ-SEC-007 @REQ-SEC-010 @error
  Scenario: GenerateMLKEMKey validates ML-KEM level parameter
    Given security operation context
    When GenerateMLKEMKey is called with unsupported level
    Then structured validation error is returned
    And error indicates unsupported ML-KEM level

  @REQ-SEC-007 @REQ-SEC-009 @error
  Scenario: Security operations validate key parameters
    Given security operation context
    When encryption operation is called with nil key
    Then structured validation error is returned
    And error indicates invalid key

  @REQ-SEC-007 @REQ-SEC-011 @error
  Scenario: Security operations respect context cancellation
    Given security operation context
    And a cancelled context
    When security operation is called
    Then structured context error is returned
    And error type is context cancellation
    And operation stops
