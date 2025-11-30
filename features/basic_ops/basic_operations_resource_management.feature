@domain:basic_ops @m2 @REQ-API_BASIC-085 @spec(api_basic_operations.md#93-resource-management)
Feature: Basic Operations Resource Management

  @REQ-API_BASIC-085 @happy
  Scenario: Context is used for resource management
    Given a long-running package operation
    When context is passed to operation
    Then context supports resource management
    And context enables cancellation
    And resources are managed through context lifecycle

  @REQ-API_BASIC-085 @happy
  Scenario: Context cancellation ensures resource cleanup
    Given a package operation with resources
    And a context for cancellation
    When context is cancelled during operation
    Then operation is terminated
    And resources are cleaned up
    And resource leaks are prevented

  @REQ-API_BASIC-085 @happy
  Scenario: Cleanup errors are handled gracefully
    Given a package operation requiring cleanup
    When cleanup operation fails
    Then cleanup error is logged as warning
    And cleanup does not fail entire operation
    And warnings are logged rather than errors

  @REQ-API_BASIC-085 @happy
  Scenario: Defer functions ensure cleanup occurs
    Given a package operation
    When defer function is used for cleanup
    Then cleanup occurs even when errors happen
    And cleanup is guaranteed to execute
    And resource cleanup is consistent

  @REQ-API_BASIC-085 @error
  Scenario: Resource leaks are prevented by proper cleanup
    Given multiple package operations
    When operations complete or fail
    Then all resources are released
    And no resource leaks occur
    And resources are tracked and managed
