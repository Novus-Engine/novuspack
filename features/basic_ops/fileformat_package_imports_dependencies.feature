@domain:basic_ops @m2 @REQ-API_BASIC-166 @spec(api_basic_operations.md#2151-file-format-package-imports)
Feature: fileformat package imports

  @REQ-API_BASIC-166 @happy
  Scenario: fileformat package imports define its dependency boundaries
    Given the fileformat package
    When dependencies are evaluated
    Then fileformat imports define which packages it depends on
    And dependencies are limited to what is required for binary format structures
    And dependency boundaries avoid unnecessary coupling to higher-level packages
    And dependency boundaries support maintainable package evolution
    And imports remain consistent with the documented dependency graph

