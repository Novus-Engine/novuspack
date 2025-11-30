@domain:basic_ops @m2 @REQ-API_BASIC-087 @spec(api_basic_operations.md#932-handle-cleanup-errors-gracefully)
Feature: Handle Cleanup Errors Gracefully

  @REQ-API_BASIC-087 @happy
  Scenario: Cleanup errors are logged as warnings rather than failing
    Given a package cleanup operation
    When cleanup operation fails
    Then cleanup error is logged as warning
    And cleanup failure does not fail entire operation
    And warnings are logged rather than errors

  @REQ-API_BASIC-087 @happy
  Scenario: Defer functions ensure cleanup occurs even when errors happen
    Given a package operation with cleanup
    When defer cleanup is used
    And an error occurs during operation
    Then cleanup still executes via defer
    And cleanup errors are logged as warnings
    And operation failure does not prevent cleanup

  @REQ-API_BASIC-087 @happy
  Scenario: Cleanup errors do not prevent resource cleanup
    Given package operations requiring cleanup
    When cleanup error occurs
    Then cleanup error is handled gracefully
    And other cleanup operations continue
    And resource cleanup is attempted

  @REQ-API_BASIC-087 @error
  Scenario: Cleanup errors are handled to prevent resource leaks
    Given package cleanup operations
    When cleanup errors occur
    Then errors are logged but not propagated
    And cleanup attempts continue
    And resource leaks are minimized
