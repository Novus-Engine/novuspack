@domain:basic_ops @m2 @REQ-API_BASIC-065 @spec(api_basic_operations.md#82-common-error-types)
Feature: Basic Operations: Common Error Types

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeValidation categorizes validation errors
    Given a package operation with invalid parameters
    When validation error occurs
    Then ErrTypeValidation error is returned
    And error indicates input validation failure
    And error enables appropriate validation handling

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeIO categorizes I/O errors
    Given a package file operation
    When file system error occurs
    Then ErrTypeIO error is returned
    And error indicates I/O operation failure
    And error enables retry logic handling

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeSecurity categorizes security errors
    Given a package security operation
    When security error occurs
    Then ErrTypeSecurity error is returned
    And error indicates security-related failure
    And error enables security handling

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeUnsupported categorizes unsupported feature errors
    Given a package operation with unsupported feature
    When unsupported operation is attempted
    Then ErrTypeUnsupported error is returned
    And error indicates feature not supported
    And error enables fallback handling

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeContext categorizes context errors
    Given a package operation with context
    When context cancellation or timeout occurs
    Then ErrTypeContext error is returned
    And error indicates context-related failure
    And error enables cancellation handling

  @REQ-API_BASIC-065 @happy
  Scenario: ErrTypeCorruption categorizes data corruption errors
    Given a package file operation
    When data corruption is detected
    Then ErrTypeCorruption error is returned
    And error indicates data integrity issue
    And error enables corruption handling
