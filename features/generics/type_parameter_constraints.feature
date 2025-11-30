@domain:generics @m2 @REQ-GEN-006 @spec(api_generics.md#32-type-parameter-constraints)
Feature: Type Parameter Constraints

  @REQ-GEN-006 @happy
  Scenario: Type parameter constraints use most restrictive constraint
    Given a generic type or function
    When type parameters are defined
    Then most restrictive constraint that works is used
    And comparable constraint is preferred over any
    And interface constraints are used for behavior requirements

  @REQ-GEN-006 @happy
  Scenario: Comparable constraint is used when possible
    Given a generic type requiring comparison
    When comparable constraint is applied
    Then type parameter must be comparable
    And comparison operations are available
    And type safety is enforced at compile time

  @REQ-GEN-006 @happy
  Scenario: Interface constraints enforce behavior requirements
    Given a generic type requiring specific behavior
    When interface constraint is applied
    Then type parameter must implement interface
    And required methods are available
    And type safety is enforced at compile time

  @REQ-GEN-006 @error
  Scenario: Type parameter constraints validate at compile time
    Given a generic type with type parameter constraints
    When type parameter violates constraint
    Then compile-time error is generated
    And type safety is enforced
    And invalid type instantiations are prevented

  @REQ-GEN-006 @happy
  Scenario: Type parameter constraints support multiple constraints
    Given a generic type or function
    When multiple constraints are specified
    Then type parameter must satisfy all constraints
    And constraint intersection is enforced
    And type safety is maintained
