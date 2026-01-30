@domain:basic_ops @m2 @REQ-API_BASIC-101 @spec(api_basic_operations.md#443-settargetpath-error-conditions)
Feature: SetTargetPath Error Handling

  @REQ-API_BASIC-101 @error
  Scenario: Invalid path format error
    Given a package needs target path change
    When SetTargetPath is called with invalid path format
    Then validation error is returned
    And error type is ErrTypeValidation
    And error message indicates invalid path
    And package state remains unchanged

  @REQ-API_BASIC-101 @error
  Scenario: Malformed file path error
    Given a package needs target path change
    When SetTargetPath is called with malformed path
    Then validation error is returned
    And error indicates path format issues
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Non-existent directory error
    Given a package needs target path change
    When SetTargetPath is called with path to non-existent directory
    Then validation error is returned
    And error indicates directory does not exist
    And error provides directory path information
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Directory not writable error
    Given a package needs target path change
    When SetTargetPath is called with non-writable directory
    Then validation error is returned
    And error type is ErrTypeSecurity or ErrTypeValidation
    And error indicates insufficient write permissions
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Insufficient permissions to create file
    Given a package needs target path change
    When SetTargetPath is called with directory lacking file creation permissions
    Then security error is returned
    And error indicates insufficient permissions
    And error provides permission details
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Context cancellation error
    Given a package needs target path change
    And context is set to cancel during validation
    When SetTargetPath is called
    Then context cancellation is detected
    And error indicates context cancelled
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Context timeout error
    Given a package needs target path change
    And context has very short timeout
    When SetTargetPath is called and timeout exceeds
    Then timeout error is returned
    And error indicates context deadline exceeded
    And target path is not updated

  @REQ-API_BASIC-101 @error
  Scenario: Error includes structured context
    Given a package needs target path change
    When SetTargetPath fails with validation error
    Then error is structured PackageError
    And error contains context with attempted path
    And error contains context with failure reason
    And error is suitable for debugging
