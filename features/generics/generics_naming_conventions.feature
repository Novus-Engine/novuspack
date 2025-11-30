@domain:generics @m2 @REQ-GEN-022 @spec(api_generics.md#31-naming-conventions)
Feature: Generics Naming Conventions

  @REQ-GEN-022 @happy
  Scenario: Naming conventions use descriptive type parameter names
    Given generic types and functions
    When type parameters are named
    Then descriptive names like T, U, V are used for simple cases
    And meaningful names like Key, Value, Element are used for complex cases
    And naming is consistent and clear

  @REQ-GEN-022 @happy
  Scenario: Naming conventions use simple descriptive names without Generic prefix
    Given generic types and functions
    When types are named
    Then simple descriptive names like Option, Config, WorkerPool are used
    And Generic prefix is not used in names
    And generic syntax indicates genericity

  @REQ-GEN-022 @happy
  Scenario: Naming conventions let generic syntax indicate genericity
    Given generic types and functions
    When types are defined
    Then generic syntax [T any] indicates genericity
    And type names themselves are simple and descriptive
    And naming is clear and intuitive

  @REQ-GEN-022 @happy
  Scenario: Naming conventions maintain consistency across generic types
    Given multiple generic types
    When naming conventions are applied
    Then naming is consistent across types
    And naming patterns are predictable
    And generic naming is standardized
