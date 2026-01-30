@domain:basic_ops @m2 @REQ-API_BASIC-051 @spec(api_basic_operations.md#713-validate-example-usage)
Feature: Validate Example Usage

  @REQ-API_BASIC-051 @happy
  Scenario: Validate example demonstrates basic package validation
    Given an open NovusPack package
    And a valid context
    When Validate is called with context
    Then package validation is performed
    And no error is returned for valid package

  @REQ-API_BASIC-051 @happy
  Scenario: Validate example includes error handling pattern
    Given an open NovusPack package
    And a valid context
    When Validate is called with error checking
    Then error result is checked
    And error handling follows standard Go pattern
    And function returns early on error

  @REQ-API_BASIC-051 @happy
  Scenario: Validate example follows standard Go error handling
    Given a code example demonstrating Validate usage
    When example code is examined
    Then Validate call includes error check
    And error is handled with if err != nil pattern
    And function returns error on failure

  @REQ-API_BASIC-051 @happy
  Scenario: Validate example uses context parameter correctly
    Given a code example demonstrating Validate usage
    When example code is examined
    Then context is passed as parameter
    And context is used from calling function
    And context supports standard Go patterns

  @REQ-API_BASIC-051 @error
  Scenario: Validate example handles validation errors
    Given an open NovusPack package
    And package has validation issues
    And a valid context
    When Validate is called
    Then validation error is returned
    And error indicates validation failure
    And error follows structured error format

  @REQ-API_BASIC-051 @error
  Scenario: Validate example handles corruption errors
    Given an open NovusPack package
    And package has corruption issues
    And a valid context
    When Validate is called
    Then corruption error is returned
    And error indicates corruption detected
    And error follows structured error format
