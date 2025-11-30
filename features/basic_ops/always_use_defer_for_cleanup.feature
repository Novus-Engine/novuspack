@domain:basic_ops @m2 @REQ-API_BASIC-079 @spec(api_basic_operations.md#911-always-use-defer-for-cleanup)
Feature: Always Use Defer for Cleanup

  @REQ-API_BASIC-079 @happy
  Scenario: Defer statements ensure cleanup even when errors occur
    Given a package operation that requires cleanup
    When defer Close is used for cleanup
    And an error occurs during operation
    Then cleanup still executes via defer
    And resources are properly released
    And cleanup happens regardless of errors

  @REQ-API_BASIC-079 @happy
  Scenario: Defer prevents resource leaks
    Given package operations with resources
    When defer statements are used consistently
    Then all resources are cleaned up
    And resource leaks are prevented
    And cleanup is guaranteed to execute

  @REQ-API_BASIC-079 @happy
  Scenario: Defer provides consistent cleanup behavior
    Given various package operations
    When defer cleanup pattern is applied
    Then cleanup behavior is consistent
    And cleanup occurs at function exit
    And cleanup order is predictable

  @REQ-API_BASIC-079 @error
  Scenario: Missing defer leads to resource leaks on errors
    Given a package operation without defer
    And an error occurs during operation
    Then cleanup may not execute
    And resources may leak
    And defer should always be used
