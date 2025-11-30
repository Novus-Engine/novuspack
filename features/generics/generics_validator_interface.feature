@domain:generics @m2 @REQ-GEN-015 @spec(api_generics.md#16-validator-interface)
Feature: Generics Validator Interface

  @REQ-GEN-015 @happy
  Scenario: Validator interface provides type-safe validation pattern
    Given a generic Validator interface
    When Validator is implemented for a type
    Then type-safe validation pattern is provided
    And Validate method is available

  @REQ-GEN-015 @happy
  Scenario: Validator Validate validates values
    Given a Validator implementation
    And a value to validate
    And a valid context
    When Validate is called with context and value
    Then error is nil if validation succeeds
    And error is returned if validation fails
    And type-safe validation is provided

  @REQ-GEN-015 @error
  Scenario: Validator returns structured errors on validation failure
    Given a Validator implementation
    And an invalid value
    And a valid context
    When Validate is called and validation fails
    Then structured error is returned
    And error indicates validation failure
    And error provides validation details

  @REQ-GEN-015 @happy
  Scenario: ValidationRule provides single validation rule
    Given a ValidationRule with name, predicate, and message
    When ValidationRule Validate is called
    Then predicate function determines validation result
    And message is used in error if validation fails
    And type-safe rule validation is provided

  @REQ-GEN-015 @happy
  Scenario: Validator pattern supports multiple validation rules
    Given multiple ValidationRule instances
    When rules are composed together
    Then all rules can be applied to values
    And validation logic is reusable
    And type safety is maintained
