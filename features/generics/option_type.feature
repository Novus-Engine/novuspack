@domain:generics @m2 @REQ-GEN-003 @spec(api_generics.md#11-option-type)
Feature: Option Type

  @REQ-GEN-003 @happy
  Scenario: Option type provides type-safe optional configuration values
    Given a generic Option type
    When Option is used for optional values
    Then type-safe optional handling is provided
    And Option wraps values of any type
    And Option indicates if value is set

  @REQ-GEN-003 @happy
  Scenario: Option Set method sets optional value
    Given an Option type instance
    When Set is called with a value
    Then value is stored in Option
    And Option indicates value is set
    And Option can be retrieved later

  @REQ-GEN-003 @happy
  Scenario: Option Get method retrieves value and presence flag
    Given an Option type with set value
    When Get is called
    Then value is returned
    And boolean indicates if value was set
    And type-safe retrieval is provided

  @REQ-GEN-003 @happy
  Scenario: Option GetOrDefault returns value or default
    Given an Option type
    And a default value
    When GetOrDefault is called with default
    Then value is returned if Option is set
    And default value is returned if Option is not set
    And type-safe default handling is provided

  @REQ-GEN-003 @happy
  Scenario: Option IsSet checks if value is set
    Given an Option type
    When IsSet is called
    Then boolean indicates if value is set
    And Option state can be queried
    And type-safe state checking is provided

  @REQ-GEN-003 @happy
  Scenario: Option Clear removes set value
    Given an Option type with set value
    When Clear is called
    Then value is removed
    And Option indicates value is not set
    And Option can be reused
