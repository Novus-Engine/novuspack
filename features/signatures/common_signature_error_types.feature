@domain:signatures @m2 @REQ-SIG-052 @spec(api_signatures.md#52-common-signature-error-types)
Feature: Common Signature Error Types

  @REQ-SIG-052 @happy
  Scenario: Specific signature error types are defined
    Given a NovusPack package
    When signature error types are examined
    Then ErrSignatureInvalid is defined for invalid signature format
    And ErrSignatureNotFound is defined for signature not found
    And ErrSignatureValidationFailed is defined for validation failures
    And ErrUnsupportedSignatureType is defined for unsupported algorithms
    And ErrInvalidSignatureData is defined for invalid data format
    And ErrSignatureTooLarge is defined for oversized signatures
    And ErrKeyNotFound is defined for missing signing keys
    And ErrInvalidKey is defined for invalid key format
    And ErrSignatureGenerationFailed is defined for generation failures
    And ErrSignatureCorrupted is defined for corrupted signature data

  @REQ-SIG-052 @happy
  Scenario: Error types provide granular error handling
    Given a NovusPack package
    When signature operation encounters specific error
    Then appropriate error type constant is returned
    And error type enables granular error handling
    And error type enables error categorization
    And error type supports error handling patterns

  @REQ-SIG-052 @error
  Scenario: Invalid signature returns ErrSignatureInvalid
    Given a NovusPack package
    And an invalid signature format
    When signature validation is attempted
    Then ErrSignatureInvalid error is returned
    And error indicates invalid signature format or data
    And error follows structured error format

  @REQ-SIG-052 @error
  Scenario: Missing signature returns ErrSignatureNotFound
    Given a NovusPack package
    And a signature index that does not exist
    When signature is retrieved by index
    Then ErrSignatureNotFound error is returned
    And error indicates signature not found at index
    And error follows structured error format
