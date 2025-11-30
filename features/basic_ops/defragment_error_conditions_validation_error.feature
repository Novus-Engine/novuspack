@domain:basic_ops @m2 @REQ-API_BASIC-056 @spec(api_basic_operations.md#722-defragment-error-conditions)
Feature: Defragment error conditions

  @REQ-API_BASIC-056 @error
  Scenario: Defragment returns validation error when package is not open
    Given a NovusPack package that is not open
    When Defragment is called
    Then validation error is returned
    And error indicates package is not open
    And defragmentation is rejected

  @REQ-API_BASIC-056 @error
  Scenario: Defragment returns validation error in read-only mode
    Given an open NovusPack package in read-only mode
    When Defragment is called
    Then validation error is returned
    And error indicates package is read-only
    And defragmentation requires writable package

  @REQ-API_BASIC-056 @error
  Scenario: Defragment returns I/O error on file system errors
    Given an open NovusPack package
    And file system error occurs during defragmentation
    When Defragment is called
    Then I/O error is returned
    And error indicates file system issue
    And error provides details about failure

  @REQ-API_BASIC-056 @error
  Scenario: Defragment returns context error on cancellation
    Given an open NovusPack package
    And a cancelled context
    When Defragment is called with cancelled context
    Then context error is returned
    And error type is context cancellation

  @REQ-API_BASIC-056 @error
  Scenario: Defragment returns context error on timeout
    Given an open NovusPack package
    And a context with expired timeout
    When Defragment is called
    Then context error is returned
    And error type is context timeout
