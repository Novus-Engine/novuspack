@domain:basic_ops @m2 @REQ-API_BASIC-034 @spec(api_basic_operations.md#432-behavior)
Feature: CreateWithOptions behavior

  @REQ-API_BASIC-034 @happy
  Scenario: CreateWithOptions uses Create internally and applies options
    Given a package needs to be created
    When CreateWithOptions is called with options
    Then Create is called internally
    And options are applied to package configuration
    And package is configured according to options

  @REQ-API_BASIC-034 @happy
  Scenario: CreateWithOptions configures package with provided options
    Given CreateWithOptions options
    When CreateWithOptions is called
    Then package structure is configured
    And package metadata is initialized
    And options are applied correctly
    And package is ready for file operations

  @REQ-API_BASIC-034 @happy
  Scenario: CreateWithOptions prepares package in memory
    Given CreateWithOptions call
    When package is created
    Then package exists in memory
    And package is not yet written to disk
    And package can be configured further

  @REQ-API_BASIC-034 @error
  Scenario: CreateWithOptions validates options before applying
    Given invalid CreateWithOptions options
    When CreateWithOptions is called
    Then validation error is returned
    And error indicates invalid options
    And package is not created
