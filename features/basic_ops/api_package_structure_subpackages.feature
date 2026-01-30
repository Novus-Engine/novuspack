@domain:basic_ops @m2 @REQ-API_BASIC-112 @spec(api_basic_operations.md#211-api-package-structure)
Feature: API package structure

  @REQ-API_BASIC-112 @happy
  Scenario: API package structure defines subpackage organization
    Given the Go API v1 package layout
    When public packages are grouped by concern
    Then subpackages are defined for major capability areas
    And each subpackage has a clear responsibility boundary
    And shared types avoid circular dependencies between subpackages
    And the organization supports common consumer workflows
    And the documented structure matches the intended public API layout

