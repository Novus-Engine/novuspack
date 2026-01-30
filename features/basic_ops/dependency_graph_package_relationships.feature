@domain:basic_ops @m2 @REQ-API_BASIC-115 @spec(api_basic_operations.md#221-dependency-graph)
Feature: Dependency graph

  @REQ-API_BASIC-115 @happy
  Scenario: Dependency graph defines package dependency relationships
    Given the Go API v1 package set
    When package dependencies are established
    Then a dependency graph defines which packages depend on others
    And dependencies flow from higher-level packages to lower-level utilities
    And shared abstractions reduce coupling between capability packages
    And dependency relationships are consistent with import boundaries
    And dependency design supports maintainable evolution of the API

