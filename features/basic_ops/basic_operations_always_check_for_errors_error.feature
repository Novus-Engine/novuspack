@domain:basic_ops @m2 @REQ-API_BASIC-072 @spec(api_basic_operations.md#841-always-check-for-errors)
Feature: Basic Operations: Always Check for Errors

  @REQ-API_BASIC-072 @happy
  Scenario: Error return values are always checked
    Given a package operation that returns error
    When operation is called
    Then error return value is checked
    And error is handled appropriately
    And error is never ignored

  @REQ-API_BASIC-072 @happy
  Scenario: Checking errors prevents silent failures
    Given a package operation that may fail
    When error is checked after operation
    Then failures are detected immediately
    And failures are handled properly
    And silent failures are prevented

  @REQ-API_BASIC-072 @happy
  Scenario: Error checking is consistent across operations
    Given various package operations
    When operations are performed
    Then all error return values are checked
    And error checking pattern is consistent
    And no errors are ignored

  @REQ-API_BASIC-072 @error
  Scenario: Ignored errors lead to undefined behavior
    Given a package operation
    When error return value is ignored
    Then operation failure may go undetected
    And package state may be inconsistent
    And errors should always be checked
