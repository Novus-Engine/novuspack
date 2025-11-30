@domain:generics @m6 @spec(api_generics.md#33-error-handling)
Feature: Generic error wrappers behave consistently

  @REQ-GEN-002 @error
  Scenario: Error wrapper patterns preserve error code consistently
    Given an operation that returns an error with code
    When the error is wrapped generically
    Then error code is preserved in wrapped error
    And error code remains consistent
    And error code is accessible

  @REQ-GEN-002 @error
  Scenario: Error wrapper patterns preserve error message consistently
    Given an operation that returns an error with message
    When the error is wrapped generically
    Then error message is preserved in wrapped error
    And error message remains consistent
    And error message is accessible

  @REQ-GEN-002 @error
  Scenario: Error wrapper patterns preserve error context consistently
    Given an operation that returns an error with context
    When the error is wrapped generically
    Then error context is preserved in wrapped error
    And error context remains consistent
    And error context is accessible

  @REQ-GEN-002 @error
  Scenario: Error wrapper patterns maintain structured error format
    Given structured errors
    When errors are wrapped generically
    Then structured error format is maintained
    And error structure remains consistent
    And error inspection is supported

  @REQ-GEN-002 @happy
  Scenario: Generic error wrappers work with Result type
    Given a Result type with error
    When Result error is wrapped generically
    Then error wrapping is consistent
    And Result error handling is maintained
    And type-safe error handling is provided
