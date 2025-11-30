@domain:basic_ops @m2 @REQ-API_BASIC-044 @spec(api_basic_operations.md#521-openwithvalidation-behavior)
Feature: OpenWithValidation Operation Behavior

  @REQ-API_BASIC-044 @happy
  Scenario: OpenWithValidation opens package file
    Given a NovusPack package instance
    And a valid context
    And an existing package file
    When OpenWithValidation is called
    Then package file is opened
    And file handle is established
    And package is ready for validation

  @REQ-API_BASIC-044 @happy
  Scenario: OpenWithValidation performs full package validation
    Given a NovusPack package instance
    And a valid context
    And an existing valid package file
    When OpenWithValidation is called
    Then package structure is validated
    And package checksums are verified
    And package signatures are validated
    And package integrity is confirmed

  @REQ-API_BASIC-044 @happy
  Scenario: OpenWithValidation ensures package integrity before operations
    Given a NovusPack package instance
    And a valid context
    And a valid package file
    When OpenWithValidation is called successfully
    Then package integrity is verified
    And package is ready for subsequent operations
    And no integrity issues are present

  @REQ-API_BASIC-044 @happy
  Scenario: OpenWithValidation returns detailed error information on failure
    Given a NovusPack package instance
    And a valid context
    And a package file with validation issues
    When OpenWithValidation is called
    Then detailed error information is returned
    And error specifies which validation failed
    And error indicates structure, checksum, or signature issue
    And error follows structured error format

  @REQ-API_BASIC-044 @error
  Scenario: OpenWithValidation fails when package structure is invalid
    Given a NovusPack package instance
    And a valid context
    And a package file with invalid structure
    When OpenWithValidation is called
    Then validation error is returned
    And error indicates structure validation failure
    And package is not opened
