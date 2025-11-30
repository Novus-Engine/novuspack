@domain:file_mgmt @m2 @REQ-FILEMGMT-121 @spec(api_file_management.md#1244-use-context-for-cancellation-with-structured-errors)
Feature: Use Context for Cancellation with Structured Errors

  @REQ-FILEMGMT-121 @happy
  Scenario: Context cancellation is used with structured errors
    Given an open NovusPack package
    And a context that can be cancelled
    When file management operations are performed
    Then context cancellation is supported
    And structured context errors are returned on cancellation
    And error context includes cancellation information

  @REQ-FILEMGMT-121 @happy
  Scenario: Context timeout handling uses structured errors
    Given an open NovusPack package
    And a context with timeout
    When file management operations are performed
    Then context timeout handling is supported
    And structured context timeout errors are returned
    And error context includes timeout information

  @REQ-FILEMGMT-121 @happy
  Scenario: Context cancellation and timeout provide appropriate error context
    Given an open NovusPack package
    And a context with cancellation or timeout
    When context errors occur
    Then structured errors include appropriate context
    And error context helps identify cancellation or timeout
    And error handling is improved with context
