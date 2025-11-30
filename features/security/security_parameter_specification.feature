@domain:security @m2 @REQ-SEC-097 @REQ-SEC-103 @spec(api_security.md#522-parameters)
Feature: Security Parameter Specification

  @REQ-SEC-097 @happy
  Scenario: ML-KEM key generation accepts context parameter
    Given an open NovusPack package
    And a valid context
    When GenerateMLKEMKey is called
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-097 @happy
  Scenario: ML-KEM key generation accepts security level parameter
    Given an open NovusPack package
    And a valid context
    And security level (1-5)
    When GenerateMLKEMKey is called with level
    Then security level parameter is accepted
    And level determines key size and security
    And level parameter defines key creation interface

  @REQ-SEC-097 @happy
  Scenario: ML-KEM key generation parameters define key creation interface
    Given an open NovusPack package
    And key generation requirements
    When key generation parameters are examined
    Then context parameter enables cancellation and timeout
    And security level parameter determines key characteristics
    And parameters define complete key creation interface

  @REQ-SEC-103 @happy
  Scenario: ML-KEM encryption accepts context parameter
    Given an open NovusPack package
    And ML-KEM key
    And a valid context
    When Encrypt or Decrypt is called
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-103 @happy
  Scenario: ML-KEM encryption accepts plaintext parameter
    Given an open NovusPack package
    And ML-KEM key
    And plaintext data
    And a valid context
    When Encrypt is called with plaintext
    Then plaintext parameter is accepted
    And plaintext is encrypted
    And plaintext parameter defines encryption input

  @REQ-SEC-103 @happy
  Scenario: ML-KEM encryption accepts ciphertext parameter
    Given an open NovusPack package
    And ML-KEM key
    And ciphertext data
    And a valid context
    When Decrypt is called with ciphertext
    Then ciphertext parameter is accepted
    And ciphertext is decrypted
    And ciphertext parameter defines decryption input

  @REQ-SEC-103 @happy
  Scenario: ML-KEM encryption parameters define encryption interface
    Given an open NovusPack package
    And ML-KEM key
    And encryption requirements
    When encryption parameters are examined
    Then context parameter enables cancellation and timeout
    And plaintext/ciphertext parameters define data input
    And parameters define complete encryption interface

  @REQ-SEC-097 @REQ-SEC-103 @error
  Scenario: ML-KEM operations fail with invalid parameters
    Given an open NovusPack package
    And invalid parameters
    When ML-KEM operation is called
    Then structured validation error is returned
    And error indicates invalid parameters
