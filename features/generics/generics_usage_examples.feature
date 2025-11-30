@domain:generics @m2 @REQ-GEN-009 @spec(api_generics.md#111-option-type-usage-examples)
Feature: Generics Usage Examples

  @REQ-GEN-009 @happy
  Scenario: Option type usage demonstrates creating option with value
    Given a string value
    When Option Set is called with value
    Then Option stores the value
    And Option IsSet returns true
    And Option Get returns value and true

  @REQ-GEN-009 @happy
  Scenario: Option type usage demonstrates GetOrDefault with set value
    Given an Option with set value
    And a default value
    When GetOrDefault is called with default
    Then set value is returned
    And default value is not used
    And Option value is retrieved

  @REQ-GEN-009 @happy
  Scenario: Option type usage demonstrates GetOrDefault with unset value
    Given an Option without set value
    And a default value
    When GetOrDefault is called with default
    Then default value is returned
    And Option value is not used
    And default fallback works

  @REQ-GEN-009 @happy
  Scenario: Option type usage demonstrates checking if value is set
    Given an Option type
    When IsSet is called
    Then boolean indicates value presence
    And Option state can be queried
    And conditional logic can be applied

  @REQ-GEN-009 @happy
  Scenario: Option type usage demonstrates optional value patterns
    Given generic Option type
    When Option is used in different contexts
    Then optional value patterns are demonstrated
    And type-safe optional handling is shown
    And best practices are illustrated
