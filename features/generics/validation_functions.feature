@domain:generics @m2 @REQ-GEN-007 @spec(api_generics.md#22-validation-functions)
Feature: Validation Functions

  @REQ-GEN-007 @happy
  Scenario: ValidateWith validates single value using validator
    Given a value to validate
    And a Validator implementation
    And a valid context
    When ValidateWith is called with context, value, and validator
    Then value is validated before processing
    And error is returned if validation fails
    And type-safe validation is provided

  @REQ-GEN-007 @happy
  Scenario: ValidateAll validates multiple values using validator
    Given multiple values to validate
    And a Validator implementation
    And a valid context
    When ValidateAll is called with context, values, and validator
    Then all values are validated
    And slice of errors is returned
    And each value validation error is captured
    And type-safe batch validation is provided

  @REQ-GEN-007 @happy
  Scenario: ComposeValidators creates composite validator
    Given multiple Validator implementations
    When ComposeValidators is called with validators
    Then composite validator is created
    And composite validator runs all validators
    And validation errors from all validators are collected
    And type-safe validator composition is provided

  @REQ-GEN-007 @error
  Scenario: Validation functions prevent processing invalid values
    Given an invalid value
    And a Validator implementation that rejects the value
    And a valid context
    When ValidateWith is called
    Then validation fails before processing
    And error indicates validation failure
    And invalid value is not processed

  @REQ-GEN-007 @happy
  Scenario: Validation functions support generic type parameters
    Given generic validation functions
    When validation functions are used with different types
    Then type safety is enforced at compile time
    And validation functions work with any type
    And generic patterns are reusable
