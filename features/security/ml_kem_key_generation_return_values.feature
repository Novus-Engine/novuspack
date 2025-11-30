@domain:security @m2 @REQ-SEC-098 @spec(api_security.md#523-returns)
Feature: ML-KEM Key Generation Return Values

  @REQ-SEC-098 @happy
  Scenario: ML-KEM key generation returns new MLKEMKey instance
    Given an open NovusPack package
    And a valid context
    And valid security level
    When GenerateMLKEMKey is called
    Then new MLKEMKey instance is returned
    And key instance contains PublicKey field
    And key instance contains PrivateKey field
    And key instance contains Level field
    And key instance follows ML-KEM key structure

  @REQ-SEC-098 @happy
  Scenario: ML-KEM key generation returns key with correct security level
    Given an open NovusPack package
    And a valid context
    And security level 3
    When GenerateMLKEMKey is called with level 3
    Then returned key has Level field set to 3
    And key security level matches requested level
    And key provides appropriate security strength

  @REQ-SEC-098 @happy
  Scenario: ML-KEM key generation returns complete key pair
    Given an open NovusPack package
    And a valid context
    And valid security level
    When GenerateMLKEMKey is called
    Then returned key contains valid public key data
    And returned key contains valid private key data
    And public and private keys form valid key pair
    And key pair can be used for encryption and decryption

  @REQ-SEC-098 @error
  Scenario: ML-KEM key generation returns error on failure
    Given an open NovusPack package
    And a valid context
    And invalid security level or generation failure
    When GenerateMLKEMKey is called
    Then error is returned instead of key
    And error indicates specific failure reason
    And error follows structured error format
