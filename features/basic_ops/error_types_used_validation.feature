@domain:basic_ops @m2 @REQ-API_BASIC-066 @spec(api_basic_operations.md#821-error-types-used)
Feature: Error Types Used

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeValidation is used for validation errors
    Given package operations with invalid parameters
    When validation errors occur
    Then ErrTypeValidation errors are returned
    And errors indicate input validation failure
    And errors enable validation error handling

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeIO is used for I/O errors
    Given package file operations
    When file system errors occur
    Then ErrTypeIO errors are returned
    And errors indicate I/O operation failure
    And errors enable I/O error handling

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeSecurity is used for security errors
    Given package security operations
    When security errors occur
    Then ErrTypeSecurity errors are returned
    And errors indicate security-related failure
    And errors enable security error handling

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeUnsupported is used for unsupported features
    Given package operations with unsupported features
    When unsupported operations are attempted
    Then ErrTypeUnsupported errors are returned
    And errors indicate feature not supported
    And errors enable fallback handling

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeContext is used for context errors
    Given package operations with context
    When context cancellation or timeout occurs
    Then ErrTypeContext errors are returned
    And errors indicate context-related failure
    And errors enable cancellation handling

  @REQ-API_BASIC-066 @happy
  Scenario: ErrTypeCorruption is used for data corruption errors
    Given package file operations
    When data corruption is detected
    Then ErrTypeCorruption errors are returned
    And errors indicate data integrity issue
    And errors enable corruption handling
