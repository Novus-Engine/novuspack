@domain:signatures @m2 @v2 @REQ-SIG-053 @spec(api_signatures.md#521-specific-signature-error-types)
Feature: Specific Signature Error Types

  @REQ-SIG-053 @happy
  Scenario: Specific signature error types are defined with typed contexts
    Given a NovusPack package
    When specific signature error types are examined
    Then ErrSignatureInvalid is defined with SignatureErrorContext
    And ErrSignatureNotFound is defined with typed context
    And ErrSignatureValidationFailed is defined with SignatureErrorContext
    And ErrUnsupportedSignatureType is defined with UnsupportedErrorContext
    And ErrInvalidSignatureData is defined with ValidationErrorContext
    And ErrSignatureTooLarge is defined with ValidationErrorContext
    And ErrKeyNotFound is defined with SecurityErrorContext
    And ErrInvalidKey is defined with SecurityErrorContext
    And ErrSignatureGenerationFailed is defined with SignatureErrorContext
    And ErrSignatureCorrupted is defined with typed context

  @REQ-SIG-053 @happy
  Scenario: SignatureErrorContext provides signature operation context
    Given a NovusPack package
    When signature error with SignatureErrorContext is created
    Then context contains SignatureIndex field
    And context contains Algorithm field
    And context contains Operation field
    And context enables detailed error debugging

  @REQ-SIG-053 @happy
  Scenario: UnsupportedErrorContext provides unsupported operation context
    Given a NovusPack package
    When unsupported signature type error is created
    Then UnsupportedErrorContext contains SignatureType field
    And context contains SupportedTypes field with list of supported types
    And context contains Operation field
    And context enables informative error messages

  @REQ-SIG-053 @happy
  Scenario: SecurityErrorContext provides security operation context
    Given a NovusPack package
    When security-related signature error is created
    Then SecurityErrorContext contains KeyID field
    And context contains KeyType field
    And context contains Operation field
    And context enables security error analysis

  @REQ-SIG-053 @error
  Scenario: Specific signature error types handle error conditions correctly
    Given a NovusPack package
    When signature operation encounters specific error
    Then appropriate error type constant is returned
    And typed context structure is included
    And error provides detailed debugging information
    And error follows structured error format
