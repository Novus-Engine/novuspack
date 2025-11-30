@domain:generics @m2 @REQ-GEN-003 @spec(api_generics.md#33-error-handling)
Feature: Generic error wrapping patterns

  @happy
  Scenario: Generic error wrapping preserves error chain
    Given an error and generic context
    When error is wrapped with generic context
    Then error chain is preserved
    And type-safe context is accessible
    And error unwrapping works correctly

  @error
  Scenario: Generic error operations handle type mismatches
    Given generic error operations
    When type mismatch occurs
    Then appropriate error is returned
    And type safety is maintained

  @REQ-GEN-005 @error
  Scenario: Generic methods with context respect cancellation
    Given a generic method with context parameter
    And a cancelled context
    When generic method is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-GEN-006 @error
  Scenario: Generic type parameters are validated
    Given a generic method with type constraints
    When method is called with type violating constraints
    Then compilation or runtime error occurs
    And type safety is enforced

  @REQ-GEN-007 @error
  Scenario: Generic validator functions validate input
    Given a generic validator function
    When validator is called with invalid input
    Then structured validation error is returned
    And error indicates validation failure

  @REQ-GEN-008 @error
  Scenario: Context errors propagate through generic error handling
    Given a generic method that returns context error
    When error is handled generically
    Then context error is preserved
    And error type remains ErrTypeContext
