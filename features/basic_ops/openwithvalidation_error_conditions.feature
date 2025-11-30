@domain:basic_ops @m2 @REQ-API_BASIC-045 @spec(api_basic_operations.md#522-openwithvalidation-error-conditions)
Feature: OpenWithValidation Error Conditions

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation returns all errors from Open method
    Given a NovusPack package instance
    And a valid context
    And conditions that cause Open to fail
    When OpenWithValidation is called
    Then same errors as Open method are returned
    And error types match Open method errors
    And error details are preserved

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation returns validation error when package validation fails
    Given a NovusPack package instance
    And a valid context
    And a package file with validation failure
    When OpenWithValidation is called
    Then validation error is returned
    And error indicates package validation failed
    And error follows structured error format

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation returns error for invalid signatures
    Given a NovusPack package instance
    And a valid context
    And a package file with invalid signatures
    When OpenWithValidation is called
    Then validation error is returned
    And error indicates invalid signatures
    And error specifies signature validation failure

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation returns corruption error when checksums don't match
    Given a NovusPack package instance
    And a valid context
    And a package file with mismatched checksums
    When OpenWithValidation is called
    Then corruption error is returned
    And error indicates checksum mismatch
    And error specifies data integrity issue

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation returns corruption error for data integrity issues
    Given a NovusPack package instance
    And a valid context
    And a package file with data integrity problems
    When OpenWithValidation is called
    Then corruption error is returned
    And error indicates data integrity issue
    And error follows structured error format

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation handles file not found errors
    Given a NovusPack package instance
    And a valid context
    And a non-existent package file path
    When OpenWithValidation is called
    Then file not found error is returned
    And error indicates package file does not exist

  @REQ-API_BASIC-045 @error
  Scenario: OpenWithValidation handles invalid package format errors
    Given a NovusPack package instance
    And a valid context
    And a file with invalid package format
    When OpenWithValidation is called
    Then invalid format error is returned
    And error indicates format validation failure
