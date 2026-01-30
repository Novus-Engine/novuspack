@domain:core @m2 @REQ-CORE-161 @spec(api_core.md#1021-packageerror-type-definition)
Feature: PackageError type definition defines structured error structure

  @REQ-CORE-161 @happy
  Scenario: PackageError type defines structured error structure
    Given an error condition in package operations
    When a structured error is created
    Then PackageError type is used for the structure
    And the type definition matches the specification
    And the structure supports type-safe context
