@domain:security @m2 @REQ-SEC-107 @spec(api_security.md#54-ml-kem-key-information)
Feature: ML-KEM Key Information

  @REQ-SEC-107 @happy
  Scenario: ML-KEM key information provides GetPublicKey method
    Given an open NovusPack package
    And a valid context
    And ML-KEM key instance
    When GetPublicKey is called
    Then public key data is returned
    And public key is []byte type
    And public key can be used for encryption
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-107 @happy
  Scenario: ML-KEM key information provides GetLevel method
    Given an open NovusPack package
    And a valid context
    And ML-KEM key instance
    When GetLevel is called
    Then security level (1-5) is returned
    And level indicates key security strength
    And level matches key generation level

  @REQ-SEC-107 @happy
  Scenario: ML-KEM key information provides Clear method
    Given an open NovusPack package
    And a valid context
    And ML-KEM key instance with sensitive data
    When Clear is called
    Then sensitive key data is cleared from memory
    And private key data is securely wiped
    And key instance is safe for disposal

  @REQ-SEC-107 @happy
  Scenario: ML-KEM key information provides key metadata access
    Given an open NovusPack package
    And a valid context
    And ML-KEM key instance
    When key information is retrieved
    Then public key is accessible
    And security level is accessible
    And key metadata supports key management

  @REQ-SEC-107 @happy
  Scenario: ML-KEM key information provides secure cleanup
    Given an open NovusPack package
    And a valid context
    And ML-KEM key instance after use
    When Clear is called for cleanup
    Then sensitive data is securely removed
    And memory is properly cleaned
    And key instance is safe for garbage collection

  @REQ-SEC-107 @error
  Scenario: ML-KEM key information handles invalid key access
    Given an open NovusPack package
    And a valid context
    And nil or invalid ML-KEM key
    When key information is accessed
    Then appropriate error is returned
    And error indicates invalid key
    And error follows structured error format
