@domain:basic_ops @m2 @REQ-API_BASIC-033 @spec(api_basic_operations.md#431-createwithoptions-parameters) @spec(api_basic_operations.md#633-createwithoptions-parameters)
Feature: CreateWithOptions parameters

  @REQ-API_BASIC-033 @happy
  Scenario: CreateWithOptions accepts option-based creation parameters
    Given a target package path
    And an options list describing package configuration
    When CreateWithOptions is called with a path and options
    Then the path parameter identifies the package to create
    And the options parameter is used for option-based creation configuration
    And options are applied as part of the creation process
    And the package is configured in memory based on the provided options
    And option-based creation supports complex configurations compared to defaults

