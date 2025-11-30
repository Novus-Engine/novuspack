@domain:basic_ops @m2 @REQ-API_BASIC-049 @spec(api_basic_operations.md#613-close-example-usage)
Feature: Close Example Usage

  @REQ-API_BASIC-049 @happy
  Scenario: Close is used with error checking
    Given a NovusPack package operation
    When Close is called after operations
    Then error return value is checked
    And error is handled appropriately
    And package cleanup is verified

  @REQ-API_BASIC-049 @happy
  Scenario: Close is used with defer for guaranteed cleanup
    Given a NovusPack package operation
    When defer Close is used at function start
    And function operations execute
    Then Close is guaranteed to execute
    And cleanup happens even on errors
    And resource leaks are prevented

  @REQ-API_BASIC-049 @happy
  Scenario: Close example demonstrates proper usage pattern
    Given a package is opened or created
    When package operations are performed
    And Close is called with error checking
    Then package lifecycle is complete
    And resources are properly released
    And error handling is demonstrated
