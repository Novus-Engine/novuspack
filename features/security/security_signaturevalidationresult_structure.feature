@domain:security @m2 @REQ-SEC-079 @spec(api_security.md#23-signaturevalidationresult-struct)
Feature: Security: SignatureValidationResult Structure

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains signature index
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then Index field contains signature index in package
    And index identifies signature position
    And index enables result lookup

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains signature type
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then Type field contains signature type identifier
    And type identifies signature algorithm
    And type enables type-specific processing

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains validity status
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then Valid field indicates whether signature is valid
    And validity status is boolean
    And validity enables validation assessment

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains trust status
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then Trusted field indicates whether signature is trusted
    And trust status is boolean
    And trust enables trust assessment

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains error message
    Given an open NovusPack package
    And package with invalid signature
    When SignatureValidationResult struct is examined
    Then Error field contains error message if validation failed
    And error message is string
    And error provides validation failure details

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains timestamp
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then Timestamp field contains when signature was created
    And timestamp is uint32
    And timestamp enables temporal validation

  @REQ-SEC-079 @happy
  Scenario: SignatureValidationResult struct contains public key
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult struct is examined
    Then PublicKey field contains public key used for validation if available
    And public key is byte array
    And public key enables key verification
