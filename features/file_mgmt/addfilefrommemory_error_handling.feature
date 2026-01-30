@domain:file_mgmt @m2 @REQ-FILEMGMT-206 @spec(api_file_mgmt_addition.md#226-addfilefrommemory-error-conditions)
Feature: AddFileFromMemory Error Handling

  @REQ-FILEMGMT-206 @error
  Scenario: Invalid path format error
    Given invalid or malformed path
    When AddFileFromMemory is called
    Then validation error is returned
    And error type is ErrTypeValidation
    And error indicates invalid path format
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Empty path error
    Given empty string as path
    When AddFileFromMemory is called
    Then validation error is returned
    And error indicates path cannot be empty
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Nil data pointer error
    Given nil data pointer
    When AddFileFromMemory is called
    Then validation error is returned
    And error type is ErrTypeValidation
    And error indicates data cannot be nil
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Package state error - not open
    Given package is not open
    When AddFileFromMemory is called
    Then state error is returned
    And error type is ErrTypeState
    And error indicates package must be open
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Package state error - read-only
    Given package is opened read-only
    When AddFileFromMemory is called
    Then state error is returned
    And error indicates package is read-only
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Context cancellation error
    Given context that is cancelled
    When AddFileFromMemory is called
    Then context error is returned
    And error indicates operation cancelled
    And file is not added to package
    And package state remains unchanged

  @REQ-FILEMGMT-206 @error
  Scenario: Context timeout error
    Given context with short timeout
    When AddFileFromMemory is called and timeout exceeds
    Then timeout error is returned
    And error indicates deadline exceeded
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Path already exists error
    Given path already exists in package
    And overwrite option is not set
    When AddFileFromMemory is called with same path
    Then validation error is returned
    And error indicates path already exists
    And existing file is not modified

  @REQ-FILEMGMT-206 @error
  Scenario: Processing error with encryption
    Given encryption options are provided
    And encryption setup fails
    When AddFileFromMemory is called
    Then processing error is returned
    And error indicates encryption failure
    And file is not added to package

  @REQ-FILEMGMT-206 @error
  Scenario: Structured error with context
    Given AddFileFromMemory operation fails
    When error is returned
    Then error is structured PackageError
    And error contains operation context
    And error includes attempted path
    And error includes failure reason
    And error is suitable for debugging
