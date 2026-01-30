@domain:basic_ops @m2 @REQ-API_BASIC-133 @spec(api_basic_operations.md#632-createwithoptions-function) @spec(api_basic_operations.md#7-newpackagewithoptions-function) @spec(api_basic_operations.md#newpackagewithoptions-parameters) @spec(api_basic_operations.md#newpackagewithoptions-behavior) @spec(api_basic_operations.md#newpackagewithoptions-error-conditions)
Feature: CreateWithOptions configures a package for later writing

  @REQ-API_BASIC-133 @happy
  Scenario: CreateWithOptions and NewPackageWithOptions apply specified options for later writing
    Given a target package path
    And a set of creation options
    When CreateWithOptions is called to create a package with options
    Then the package is configured in memory according to the specified options
    And the configuration persists in the in-memory package state for later writing
    And option application follows the documented parameter and behavior rules
    And invalid option combinations produce a structured error as specified

