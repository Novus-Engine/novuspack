@domain:signatures @m2 @REQ-SIG-051 @spec(api_signatures.md#51-structured-error-system)
Feature: Signature Structured Error System

  @REQ-SIG-051 @happy
  Scenario: Structured error system uses core error types
    Given a NovusPack package
    When signature operations use structured errors
    Then ErrTypeSignature is used for digital signature validation errors
    And ErrTypeValidation is used for input validation errors
    And ErrTypeSecurity is used for security-related errors
    And ErrTypeUnsupported is used for unsupported features
    And ErrTypeCorruption is used for data corruption errors

  @REQ-SIG-051 @happy
  Scenario: Signature errors follow structured error format
    Given a NovusPack package
    And an open NovusPack package
    When signature operation encounters an error
    Then structured error is returned
    And error contains error type categorization
    And error contains error message
    And error contains context information
    And error follows structured error system format

  @REQ-SIG-051 @happy
  Scenario: Structured errors provide categorization for debugging
    Given a NovusPack package
    When signature validation fails
    Then error type indicates signature validation error category
    And error message describes the specific failure
    And error context provides debugging information
    And error format enables programmatic error handling

  @REQ-SIG-051 @error
  Scenario: Structured errors handle error categorization correctly
    Given a NovusPack package
    When different error conditions occur
    Then validation errors use ErrTypeValidation
    And signature errors use ErrTypeSignature
    And security errors use ErrTypeSecurity
    And unsupported operation errors use ErrTypeUnsupported
    And corruption errors use ErrTypeCorruption
