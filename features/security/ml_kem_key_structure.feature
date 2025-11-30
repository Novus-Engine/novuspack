@domain:security @m2 @REQ-SEC-092 @spec(api_security.md#51-ml-kem-key-structure)
Feature: ML-KEM Key Structure

  @REQ-SEC-092 @happy
  Scenario: ML-KEM key structure contains PublicKey field
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When ML-KEM key structure is examined
    Then PublicKey contains ML-KEM public key data
    And PublicKey is []byte type
    And PublicKey is used for encryption operations

  @REQ-SEC-092 @happy
  Scenario: ML-KEM key structure contains PrivateKey field
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When ML-KEM key structure is examined
    Then PrivateKey contains ML-KEM private key data
    And PrivateKey is []byte type
    And PrivateKey is used for decryption operations

  @REQ-SEC-092 @happy
  Scenario: ML-KEM key structure contains Level field
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When ML-KEM key structure is examined
    Then Level contains security level (1-5)
    And Level is int type
    And Level indicates key security strength

  @REQ-SEC-092 @happy
  Scenario: ML-KEM key structure represents key pair
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When ML-KEM key structure is used
    Then structure represents complete ML-KEM key pair
    And key pair consists of public and private keys
    And key pair includes associated security level

  @REQ-SEC-092 @happy
  Scenario: ML-KEM key structure follows ML-KEM requirements
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When ML-KEM key structure is validated
    Then key structure adheres to ML-KEM specifications
    And key format matches ML-KEM standard
    And key data is properly structured

  @REQ-SEC-092 @error
  Scenario: ML-KEM key structure handles invalid key data
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure with invalid data
    When ML-KEM key structure is validated
    Then validation error is returned
    And error indicates invalid key format
    And error follows structured error format
