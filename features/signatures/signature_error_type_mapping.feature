@domain:signatures @m2 @REQ-SIG-054 @spec(api_signatures.md#522-error-type-mapping)
Feature: Signature Error Type Mapping

  @REQ-SIG-054 @happy
  Scenario: Error messages map to error types correctly
    Given a NovusPack package
    When error type mapping is examined
    Then ErrSignatureInvalid maps to ErrTypeSignature
    And ErrSignatureNotFound maps to ErrTypeValidation
    And ErrSignatureValidationFailed maps to ErrTypeSignature
    And ErrUnsupportedSignatureType maps to ErrTypeUnsupported
    And ErrInvalidSignatureData maps to ErrTypeValidation
    And ErrSignatureTooLarge maps to ErrTypeValidation
    And ErrKeyNotFound maps to ErrTypeSecurity
    And ErrInvalidKey maps to ErrTypeSecurity
    And ErrSignatureGenerationFailed maps to ErrTypeSignature
    And ErrSignatureCorrupted maps to ErrTypeCorruption

  @REQ-SIG-054 @happy
  Scenario: Error type mapping enables error categorization
    Given a NovusPack package
    When signature error occurs
    Then error message maps to appropriate error type
    And error type enables error categorization
    And error type supports structured error handling
    And error type enables programmatic error handling

  @REQ-SIG-054 @happy
  Scenario: Error type mapping supports legacy error handling
    Given a NovusPack package
    When legacy error is encountered
    Then error maps to structured error type
    And error maintains backward compatibility
    And error enables migration to structured errors
