@domain:file_mgmt @m2 @REQ-FILEMGMT-111 @spec(api_file_management.md#121-structured-error-system)
Feature: File Management: Structured Error System

  @REQ-FILEMGMT-111 @happy
  Scenario: Structured error system provides error categorization
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then errors are categorized by type
    And error types provide classification
    And structured error format is used

  @REQ-FILEMGMT-111 @happy
  Scenario: Structured error system provides error context
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then errors include context information
    And error details are comprehensive
    And debugging capabilities are enhanced

  @REQ-FILEMGMT-111 @happy
  Scenario: Structured error system integrates with core error system
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then structured error system is used
    And errors follow core error patterns
    And error system is consistent with core API

  @REQ-FILEMGMT-111 @error
  Scenario: Structured error system handles various error types
    Given an open NovusPack package
    And a valid context
    When different error conditions occur
    Then validation errors are returned as structured errors
    And I/O errors are returned as structured errors
    And context errors are returned as structured errors
    And all errors follow structured error format

  @REQ-FILEMGMT-111 @happy
  Scenario: Structured error system supports error inspection
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then errors can be inspected programmatically
    And error types can be checked
    And error context can be retrieved
    And error handling is improved
