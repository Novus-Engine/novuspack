@domain:security @security @m4 @spec(api_security.md#5-ml-kem-key-management)
Feature: Enforce ML-KEM key rules

  @REQ-SEC-003 @security @happy
  Scenario: ML-KEM key structure adheres to ML-KEM requirements
    Given an open NovusPack package
    And a valid context
    And ML-KEM key pair with valid structure
    When ML-KEM key structure is examined
    Then PublicKey contains ML-KEM public key data
    And PrivateKey contains ML-KEM private key data
    And Level contains security level (1-5)
    And key structure follows ML-KEM requirements

  @REQ-SEC-003 @security @happy
  Scenario: ML-KEM key generation adheres to ML-KEM rules
    Given a valid context
    And valid security level (1-5)
    When GenerateMLKEMKey is called
    Then ML-KEM key pair is generated
    And key follows ML-KEM key structure
    And key security level matches requested level
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-003 @security @happy
  Scenario: ML-KEM encryption operations adhere to ML-KEM rules
    Given an open NovusPack package
    And a valid context
    And valid ML-KEM key
    And plaintext data
    When ML-KEM encryption is performed
    Then encryption uses ML-KEM key exchange
    And ciphertext follows ML-KEM format
    And encryption adheres to ML-KEM requirements
    When ML-KEM decryption is performed
    Then decryption uses ML-KEM key
    And plaintext matches original data

  @REQ-SEC-003 @security @error
  Scenario: ML-KEM key generation rejects invalid security levels
    Given a valid context
    And invalid security level (not 1, 3, or 5)
    When GenerateMLKEMKey is called
    Then ErrInvalidSecurityLevel error is returned
    And error indicates security level must be 1, 3, or 5
    And error follows structured error format

  @REQ-SEC-003 @security @error
  Scenario: ML-KEM operations reject invalid keys
    Given an open NovusPack package
    And a valid context
    And invalid or nil ML-KEM key
    When ML-KEM encryption is performed
    Then key validation error is returned
    And error indicates invalid key format
    And error follows structured error format

  @REQ-SEC-003 @security @happy
  Scenario: ML-KEM key information provides key metadata
    Given an open NovusPack package
    And a valid context
    And ML-KEM key with metadata
    When ML-KEM key information is retrieved
    Then key security level is provided
    And key generation timestamp is provided
    And key expiration information is provided
    And key metadata follows ML-KEM requirements
