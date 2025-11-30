@domain:basic_ops @m2 @REQ-API_BASIC-048 @spec(api_basic_operations.md#612-close-error-conditions)
Feature: Close Error Conditions

  @REQ-API_BASIC-048 @error
  Scenario: Close returns I/O error on file system errors
    Given an open NovusPack package
    And a file system error occurs during closing
    When Close is called
    Then I/O error is returned
    And error indicates file system issue
    And partial cleanup is attempted

  @REQ-API_BASIC-048 @error
  Scenario: Close returns validation error when package is not open
    Given a NovusPack package that is not open
    When Close is called
    Then validation error is returned
    And error indicates package is not currently open
    And error prevents invalid close operation

  @REQ-API_BASIC-048 @error
  Scenario: Close handles file handle closure errors
    Given an open NovusPack package
    And file handle cannot be closed properly
    When Close is called
    Then I/O error is returned
    And error indicates closure failure
    And error provides details about failure
