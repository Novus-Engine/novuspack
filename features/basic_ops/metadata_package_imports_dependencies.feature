@domain:basic_ops @m2 @REQ-API_BASIC-171 @spec(api_basic_operations.md#218-metadata-package-imports)
Feature: metadata package imports

  @REQ-API_BASIC-171 @happy
  Scenario: metadata package imports define dependency boundaries
    Given the metadata package
    When dependencies are evaluated
    Then metadata package imports define which packages it depends on
    And dependencies avoid coupling to unrelated capability packages
    And dependency boundaries support maintainable evolution of metadata functionality
    And imports are consistent with the documented dependency graph
    And import boundaries reduce risk of circular dependencies

