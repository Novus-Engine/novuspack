@domain:security @m2 @REQ-SEC-095 @spec(api_security.md#52-ml-kem-key-generation)
Feature: ML-KEM Key Generation with Security Levels

  @REQ-SEC-095 @happy
  Scenario: GenerateMLKEMKey generates new ML-KEM key pair
    Given an open NovusPack package
    And a valid context
    And security level (1-5)
    When GenerateMLKEMKey is called with level
    Then new ML-KEM key pair is generated
    And MLKEMKey instance is returned
    And key pair is generated at specified security level

  @REQ-SEC-095 @happy
  Scenario: GenerateMLKEMKey generates key at security level 1
    Given an open NovusPack package
    And a valid context
    And security level 1
    When GenerateMLKEMKey is called with level 1
    Then key pair is generated at level 1
    And key size matches level 1 requirements
    And key provides 128-bit security

  @REQ-SEC-095 @happy
  Scenario: GenerateMLKEMKey generates key at security level 3
    Given an open NovusPack package
    And a valid context
    And security level 3
    When GenerateMLKEMKey is called with level 3
    Then key pair is generated at level 3
    And key size matches level 3 requirements
    And key provides 192-bit security

  @REQ-SEC-095 @happy
  Scenario: GenerateMLKEMKey generates key at security level 5
    Given an open NovusPack package
    And a valid context
    And security level 5
    When GenerateMLKEMKey is called with level 5
    Then key pair is generated at level 5
    And key size matches level 5 requirements
    And key provides 256-bit security

  @REQ-SEC-095 @happy
  Scenario: GenerateMLKEMKey accepts context for cancellation and timeout
    Given an open NovusPack package
    And a valid context
    When GenerateMLKEMKey is called
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-099 @REQ-SEC-011 @error
  Scenario: GenerateMLKEMKey fails with invalid security level
    Given an open NovusPack package
    And invalid security level (not 1-5)
    And a valid context
    When GenerateMLKEMKey is called with invalid level
    Then ErrInvalidSecurityLevel error is returned
    And error indicates security level must be 1-5

  @REQ-SEC-099 @error
  Scenario: GenerateMLKEMKey fails if key generation fails
    Given an open NovusPack package
    And system conditions preventing key generation
    And a valid context
    When GenerateMLKEMKey is called
    Then ErrKeyGenerationFailed error is returned
    And error indicates key generation failure

  @REQ-SEC-099 @REQ-SEC-011 @error
  Scenario: GenerateMLKEMKey fails with cancelled context
    Given an open NovusPack package
    And a cancelled context
    When GenerateMLKEMKey is called
    Then ErrContextCancelled error is returned
    And error indicates context cancellation

  @REQ-SEC-099 @REQ-SEC-011 @error
  Scenario: GenerateMLKEMKey fails with context timeout
    Given an open NovusPack package
    And a context with timeout
    And timeout expires during key generation
    When GenerateMLKEMKey is called
    Then ErrContextTimeout error is returned
    And error indicates context timeout
