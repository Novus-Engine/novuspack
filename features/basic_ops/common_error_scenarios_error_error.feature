@domain:basic_ops @m2 @REQ-API_BASIC-070 @spec(api_basic_operations.md#833-common-error-scenarios)
Feature: Common Error Scenarios

  @REQ-API_BASIC-070 @happy
  Scenario: Common error scenarios demonstrate typical error cases
    Given various package operations
    When common error scenarios are encountered
    Then validation errors are demonstrated
    And I/O errors are demonstrated
    And security errors are demonstrated
    And context errors are demonstrated

  @REQ-API_BASIC-070 @error
  Scenario: File not found error scenario
    Given a package operation requiring file
    And file does not exist
    When operation is attempted
    Then file not found error is returned
    And error provides file path context

  @REQ-API_BASIC-070 @error
  Scenario: Invalid path error scenario
    Given a package operation
    And invalid or malformed path is provided
    When operation is attempted
    Then invalid path error is returned
    And error indicates path format issue

  @REQ-API_BASIC-070 @error
  Scenario: Package not open error scenario
    Given a package operation
    And package is not open
    When operation is attempted
    Then package not open error is returned
    And error indicates package must be open

  @REQ-API_BASIC-070 @error
  Scenario: Context cancellation error scenario
    Given a long-running package operation
    And context is cancelled
    When operation continues
    Then context cancellation error is returned
    And error indicates context was cancelled
