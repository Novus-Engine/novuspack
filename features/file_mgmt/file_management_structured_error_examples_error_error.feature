@domain:file_mgmt @m2 @REQ-FILEMGMT-116 @spec(api_file_management.md#123-structured-error-examples)
Feature: File Management: Structured Error Examples

  @REQ-FILEMGMT-116 @happy
  Scenario: Structured error examples demonstrate file not found error with context
    Given an open NovusPack package
    And a valid context
    And a non-existent file path
    When ExtractFile operation fails with file not found
    Then structured error is created with ErrTypeValidation
    And error includes path context
    And error includes operation context
    And error provides debugging information

  @REQ-FILEMGMT-116 @happy
  Scenario: Structured error examples demonstrate encryption failure with context
    Given an open NovusPack package
    And a valid context
    And an encryption operation
    When encryption fails
    Then structured error is created with ErrTypeEncryption
    And error includes algorithm context
    And error includes key size context
    And error includes file path context
    And error provides debugging information

  @REQ-FILEMGMT-116 @happy
  Scenario: Structured error examples demonstrate pattern matching error with context
    Given an open NovusPack package
    And a valid context
    And a file pattern
    When AddFilePattern operation fails with no files found
    Then structured error is created with ErrTypeValidation
    And error includes pattern context
    And error includes directory context
    And error includes operation context
    And error provides debugging information

  @REQ-FILEMGMT-116 @happy
  Scenario: Structured error examples provide rich context for debugging
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then structured errors include relevant context
    And error context helps identify problem source
    And error context provides operation details
    And error context enhances debugging capabilities

  @REQ-FILEMGMT-116 @error
  Scenario: Structured error examples handle various error scenarios
    Given an open NovusPack package
    And a valid context
    When different error conditions occur
    Then validation errors are created with appropriate context
    And I/O errors are created with appropriate context
    And encryption errors are created with appropriate context
    And all errors follow structured error format
