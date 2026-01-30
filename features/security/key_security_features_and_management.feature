@domain:security @m2 @v2 @REQ-SEC-015 @REQ-SEC-048 @spec(security.md#211-key-security-features) @spec(api_security.md#5-ml-kem-key-structure-and-operations)
Feature: Key Security Features and Management

  @REQ-SEC-015 @happy
  Scenario: Key security features define signature security capabilities
    Given an open NovusPack package
    And signature keys
    When key security features are examined
    Then secure key generation is provided
    And secure key storage is provided
    And secure key distribution is provided
    And signature security capabilities are defined

  @REQ-SEC-015 @happy
  Scenario: Key security features support multiple key types
    Given an open NovusPack package
    And signature keys
    When key security features are examined
    Then ML-DSA keys are supported
    And SLH-DSA keys are supported
    And PGP keys are supported
    And X.509 keys are supported

  @REQ-SEC-015 @happy
  Scenario: Key security features provide key validation
    Given an open NovusPack package
    And signature keys
    When key security features are examined
    Then key validation is provided
    And invalid keys are detected
    And key validation ensures security

  @REQ-SEC-048 @happy
  Scenario: Key management provides secure key handling
    Given an open NovusPack package
    And key management system
    When key management is examined
    Then secure key generation is provided
    And secure key storage is provided
    And secure key access is provided
    And secure key handling is ensured

  @REQ-SEC-048 @happy
  Scenario: Key management provides private key protection
    Given an open NovusPack package
    And private keys
    When key management is examined
    Then private keys are protected
    And access control for private keys is implemented
    And private key security is ensured

  @REQ-SEC-048 @happy
  Scenario: Key management provides public key distribution
    Given an open NovusPack package
    And public keys
    When key management is examined
    Then public key distribution is provided
    And public keys are accessible for validation
    And public key distribution supports verification

  @REQ-SEC-048 @happy
  Scenario: Key management supports key lifecycle management
    Given an open NovusPack package
    And key management system
    When key management is examined
    Then key creation is supported
    And key rotation is supported
    And key revocation is supported
    And key expiration is supported

  @REQ-SEC-009 @error
  Scenario: Key management validation fails with invalid key parameters
    Given an open NovusPack package
    And invalid key parameters
    When key management operation is called
    Then structured validation error is returned
    And error indicates invalid key parameters
