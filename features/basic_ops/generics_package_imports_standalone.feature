@domain:basic_ops @m2 @REQ-API_BASIC-176 @spec(api_basic_operations.md#21101-generics-package-imports)
Feature: generics package imports

  @REQ-API_BASIC-176 @happy
  Scenario: generics package remains standalone with minimal dependencies
    Given the generics package
    When imports are evaluated
    Then the package imports only what is required for generic utilities
    And the package avoids internal imports that would create cycles
    And generic abstractions remain reusable across capability packages
    And import boundaries remain consistent with the dependency graph
    And the package stays lightweight to avoid pulling in heavy dependencies

