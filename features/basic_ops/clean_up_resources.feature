@domain:basic_ops @m2 @REQ-API_BASIC-076 @spec(api_basic_operations.md#845-clean-up-resources)
Feature: Clean Up Resources

  @REQ-API_BASIC-076 @happy
  Scenario: Resources are cleaned up properly using defer statements
    Given package operations with resources
    When defer statements are used for cleanup
    Then resources are cleaned up properly
    And cleanup occurs even when errors happen
    And resource cleanup is consistent

  @REQ-API_BASIC-076 @happy
  Scenario: Packages are closed even when errors occur
    Given a package operation that may fail
    When error occurs during operation
    Then package is still closed via defer
    And cleanup happens regardless of errors
    And resource leaks are prevented

  @REQ-API_BASIC-076 @happy
  Scenario: Cleanup errors are handled gracefully
    Given a cleanup operation
    When cleanup operation fails
    Then cleanup error is logged as warning
    And cleanup failure does not fail entire operation
    And warnings are logged rather than errors

  @REQ-API_BASIC-076 @error
  Scenario: Missing cleanup leads to resource leaks
    Given package operations without proper cleanup
    When operations complete or fail
    Then resources may not be released
    And resource leaks may occur
    And defer cleanup should always be used
