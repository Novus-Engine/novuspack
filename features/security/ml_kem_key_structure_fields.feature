@domain:security @m2 @REQ-SEC-094 @spec(api_security.md#512-fields)
Feature: ML-KEM Key Structure Fields

  @REQ-SEC-094 @happy
  Scenario: ML-KEM key structure PublicKey field stores public key data
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When PublicKey field is examined
    Then PublicKey contains ML-KEM public key data
    And PublicKey is []byte type
    And PublicKey is used for encryption operations

  @REQ-SEC-094 @happy
  Scenario: ML-KEM key structure PrivateKey field stores private key data
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When PrivateKey field is examined
    Then PrivateKey contains ML-KEM private key data
    And PrivateKey is []byte type
    And PrivateKey is used for decryption operations

  @REQ-SEC-094 @happy
  Scenario: ML-KEM key structure Level field stores security level
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When Level field is examined
    Then Level contains security level (1-5)
    And Level is int type
    And Level indicates key security strength
    And higher level indicates more secure key

  @REQ-SEC-094 @happy
  Scenario: ML-KEM key structure fields provide complete key pair information
    Given an open NovusPack package
    And a valid context
    And ML-KEM key structure
    When all fields are examined
    Then PublicKey provides encryption capability
    And PrivateKey provides decryption capability
    And Level provides security strength information
    And fields together represent complete key pair
