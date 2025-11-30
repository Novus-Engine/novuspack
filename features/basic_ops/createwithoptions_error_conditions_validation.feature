@domain:basic_ops @m2 @REQ-API_BASIC-035 @spec(api_basic_operations.md#433-createwithoptions-error-conditions)
Feature: CreateWithOptions error conditions

  @REQ-API_BASIC-035 @error
  Scenario: CreateWithOptions inherits Create validation errors
    Given invalid CreateWithOptions parameters
    When CreateWithOptions is called
    Then Create validation errors are returned
    And invalid path errors are returned
    And directory not exist errors are returned
    And directory not writable errors are returned

  @REQ-API_BASIC-035 @error
  Scenario: CreateWithOptions returns error for invalid option values
    Given CreateWithOptions with invalid option values
    When CreateWithOptions is called
    Then validation error is returned
    And error indicates invalid option values
    And error specifies which options are invalid

  @REQ-API_BASIC-035 @error
  Scenario: CreateWithOptions returns error for invalid VendorID format
    Given CreateWithOptions with invalid VendorID format
    When CreateWithOptions is called
    Then validation error is returned
    And error indicates invalid VendorID format

  @REQ-API_BASIC-035 @error
  Scenario: CreateWithOptions returns security error for insufficient permissions
    Given a path with insufficient permissions
    When CreateWithOptions is called
    Then security error is returned
    And error indicates insufficient permissions

  @REQ-API_BASIC-035 @error
  Scenario: CreateWithOptions respects context cancellation
    Given a context that is cancelled
    When CreateWithOptions is called with cancelled context
    Then context error is returned
    And error type is context cancellation
