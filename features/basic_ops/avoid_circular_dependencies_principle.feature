@domain:basic_ops @m2 @REQ-API_BASIC-118 @spec(api_basic_operations.md#avoiding-circular-dependencies)
Feature: Avoiding circular dependencies

  @REQ-API_BASIC-118 @happy
  Scenario: Package structure design avoids circular dependencies
    Given the Go API v1 package structure
    When new packages and shared types are designed
    Then circular dependencies are avoided by design
    And shared types are extracted into appropriate shared packages
    And dependency direction remains consistent across the codebase
    And packages remain buildable without import cycles
    And API organization supports long-term maintainability

