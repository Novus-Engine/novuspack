@domain:basic_ops @m2 @REQ-API_BASIC-042 @spec(api_basic_operations.md#513-openpackage-error-conditions)
Feature: OpenPackage error conditions

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns I/O error for non-existent file
    Given a Package instance
    And a path to non-existent file
    When OpenPackage is called
    Then I/O error is returned
    And error indicates file does not exist
    And package is not opened

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns validation error for invalid path
    Given a Package instance
    And an invalid or malformed path
    When OpenPackage is called
    Then validation error is returned
    And error indicates invalid path
    And package is not opened

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns validation error for invalid package format
    Given a Package instance
    And a file that is not a valid NovusPack package
    When OpenPackage is called
    Then validation error is returned
    And error indicates invalid package format
    And package is not opened

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns I/O error on file system errors
    Given a Package instance
    And a path with file system issues
    When OpenPackage is called
    Then I/O error is returned
    And error indicates file system problem
    And error provides details about failure

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns security error for insufficient permissions
    Given a Package instance
    And a path with insufficient read permissions
    When OpenPackage is called
    Then security error is returned
    And error indicates insufficient permissions
    And package is not opened

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns context error on cancellation
    Given a Package instance
    And a cancelled context
    When OpenPackage is called with cancelled context
    Then context error is returned
    And error type is context cancellation

  @REQ-API_BASIC-042 @error
  Scenario: OpenPackage returns context error on timeout
    Given a Package instance
    And a context with expired timeout
    When OpenPackage is called
    Then context error is returned
    And error type is context timeout
