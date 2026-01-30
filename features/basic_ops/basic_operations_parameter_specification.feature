@skip @domain:basic_ops @m2 @REQ-API_BASIC-029 @spec(api_basic_operations.md#421-create-parameters)
Feature: Basic Operations Parameter Specification

# This feature captures parameter-level constraints for key basic operations.
# Detailed runnable scenarios live in the dedicated basic_ops feature files.

  @REQ-API_BASIC-029 @REQ-API_BASIC-020 @validation
  Scenario: Create parameters include context and target path and are validated early
    Given a new Package instance
    And an existing writable directory
    When Create is called with a context and a path in that directory
    Then the package stores the target path for later write operations
    And Create validates the target directory exists and is writable
    And Create does not create any on-disk package file

  @REQ-API_BASIC-031 @error
  Scenario: Create rejects a missing parent directory
    Given a new Package instance
    And a path whose parent directory does not exist
    When Create is called with that path
    Then a structured validation error is returned
    And the error indicates the parent directory must already exist
