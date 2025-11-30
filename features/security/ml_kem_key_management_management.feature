@domain:security @m2 @REQ-CRYPTO-003 @spec(security.md#341-ml-kem-key-management)
Feature: ML-KEM Key Management

  @REQ-CRYPTO-003 @happy
  Scenario: ML-KEM key management provides quantum-safe key generation
    Given an open NovusPack package
    And a valid context
    And security level requirements
    When ML-KEM key management generates keys
    Then ML-KEM keys are generated at specified security levels (1-5)
    And key generation uses appropriate security levels
    And generated keys are quantum-safe

  @REQ-CRYPTO-003 @happy
  Scenario: ML-KEM key management provides encryption and decryption operations
    Given an open NovusPack package
    And a valid context
    And ML-KEM key pair
    When ML-KEM key management performs operations
    Then encryption operations encrypt data using ML-KEM keys
    And decryption operations decrypt data using ML-KEM keys
    And operations provide quantum-safe encryption

  @REQ-CRYPTO-003 @happy
  Scenario: ML-KEM key management provides key access
    Given an open NovusPack package
    And a valid context
    And ML-KEM key pair
    When ML-KEM key management provides access
    Then public keys can be retrieved
    And security level information can be retrieved
    And key metadata supports key management

  @REQ-CRYPTO-003 @happy
  Scenario: ML-KEM key management provides secure key structure
    Given an open NovusPack package
    And a valid context
    And ML-KEM key pair
    When ML-KEM key management stores keys
    Then public and private key data are securely stored
    And key structure ensures secure storage
    And key format supports security requirements

  @REQ-CRYPTO-003 @happy
  Scenario: ML-KEM key management provides quantum-safe key handling
    Given an open NovusPack package
    And a valid context
    And ML-KEM key management system
    When ML-KEM key management is used
    Then all key operations are quantum-safe
    And key generation is quantum-safe
    And key storage is secure
    And key access is controlled
