@domain:generics @m2 @REQ-GEN-021 @spec(api_generics.md#3-best-practices)
Feature: Generics Best Practices

  @REQ-GEN-021 @happy
  Scenario: Best practices recommend using Result for operations that may fail
    Given operations that may fail
    When error handling is needed
    Then Result type is used
    And type-safe error handling is provided
    And error handling is consistent

  @REQ-GEN-021 @happy
  Scenario: Best practices recommend using Option for optional values
    Given optional configuration values
    When optional values are needed
    Then Option type is used
    And type-safe optional handling is provided
    And optional value patterns are consistent

  @REQ-GEN-021 @happy
  Scenario: Best practices recommend using most restrictive constraint
    Given generic types and functions
    When type parameter constraints are defined
    Then most restrictive constraint that works is used
    And comparable constraint is preferred over any
    And interface constraints are used for behavior requirements

  @REQ-GEN-021 @happy
  Scenario: Best practices recommend documenting type parameter constraints
    Given generic types and functions
    When documentation is written
    Then type parameter constraints are documented
    And usage examples are provided
    And generic vs non-generic versions are explained

  @REQ-GEN-021 @happy
  Scenario: Best practices recommend testing with multiple type instantiations
    Given generic types and functions
    When testing is performed
    Then multiple type instantiations are tested
    And type-specific test cases are used
    And compile-time type safety is verified
