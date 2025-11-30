@domain:core @m1 @REQ-CORE-011 @spec(api_core.md#111-error-types-and-categories)
Feature: Error type categories and usage

  @happy
  Scenario: ErrTypeValidation is used for input validation errors
    Given invalid input parameters
    When validation fails
    Then ErrTypeValidation error is returned
    And error indicates validation failure

  @happy
  Scenario: ErrTypeIO is used for file I/O errors
    Given file system operations
    When I/O error occurs
    Then ErrTypeIO error is returned
    And error indicates I/O failure

  @happy
  Scenario: ErrTypeSecurity is used for security-related errors
    Given security-sensitive operations
    When security check fails
    Then ErrTypeSecurity error is returned
    And error indicates security violation

  @happy
  Scenario: ErrTypeCorruption is used for data corruption
    Given package data
    When corruption is detected
    Then ErrTypeCorruption error is returned
    And error indicates data integrity issue

  @happy
  Scenario: ErrTypeUnsupported is used for unsupported features
    Given an operation
    When feature is unsupported
    Then ErrTypeUnsupported error is returned
    And error indicates unsupported feature

  @happy
  Scenario: ErrTypeContext is used for context errors
    Given a long-running operation
    When context is cancelled
    Then ErrTypeContext error is returned
    And error indicates context cancellation

  @happy
  Scenario: ErrTypeEncryption is used for encryption errors
    Given encryption operations
    When encryption fails
    Then ErrTypeEncryption error is returned
    And error indicates encryption failure

  @happy
  Scenario: ErrTypeCompression is used for compression errors
    Given compression operations
    When compression fails
    Then ErrTypeCompression error is returned
    And error indicates compression failure

  @happy
  Scenario: ErrTypeSignature is used for signature errors
    Given signature operations
    When signature validation fails
    Then ErrTypeSignature error is returned
    And error indicates signature failure

  @REQ-CORE-015 @REQ-CORE-017 @error
  Scenario: Context timeout errors use ErrTypeContext
    Given a long-running operation
    And a context with timeout
    When operation exceeds timeout
    Then ErrTypeContext error is returned
    And error indicates timeout occurred
    And error is structured with context information

  @REQ-CORE-015 @REQ-CORE-016 @error
  Scenario: Context cancellation errors use ErrTypeContext
    Given an operation in progress
    And a cancelled context
    When context cancellation is detected
    Then ErrTypeContext error is returned
    And error indicates cancellation occurred
    And error is structured with context information

  @REQ-CORE-018 @error
  Scenario: Input validation errors use ErrTypeValidation
    Given an operation requiring input
    When invalid input parameters are provided
    Then ErrTypeValidation error is returned
    And error indicates validation failure
    And error message is clear and descriptive
