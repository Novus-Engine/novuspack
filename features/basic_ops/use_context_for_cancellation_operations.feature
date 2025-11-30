@domain:basic_ops @m2 @REQ-API_BASIC-074 @spec(api_basic_operations.md#843-use-context-for-cancellation)
Feature: Use Context for Cancellation

  @REQ-API_BASIC-074 @happy
  Scenario: Context timeouts prevent operations from hanging indefinitely
    Given a package operation
    And a context with timeout configured
    When operation exceeds timeout duration
    Then operation is cancelled
    And operation does not hang indefinitely
    And timeout prevents indefinite blocking

  @REQ-API_BASIC-074 @happy
  Scenario: Context cancellation is used for long-running operations
    Given a long-running package operation
    And a context with cancellation support
    When context cancellation is requested
    Then operation is cancelled gracefully
    And cancellation is handled properly
    And resources are cleaned up

  @REQ-API_BASIC-074 @happy
  Scenario: Context cancellation is handled gracefully
    Given a package operation
    And a cancelled context
    When operation is attempted with cancelled context
    Then cancellation error is returned
    And error follows structured error format
    And error handling is graceful

  @REQ-API_BASIC-074 @happy
  Scenario: Appropriate timeouts are set for long-running operations
    Given a long-running package operation
    When context timeout is configured
    Then timeout matches operation duration
    And timeout is appropriate for operation type
    And timeout prevents unnecessary blocking

  @REQ-API_BASIC-074 @happy
  Scenario: Context cancellation supports graceful shutdown
    Given a package operation in progress
    And a context with cancellation capability
    When graceful shutdown is requested
    Then context cancellation stops operation
    And operation terminates cleanly
    And resources are released

  @REQ-API_BASIC-074 @error
  Scenario: Operations without context cancellation cannot be stopped
    Given a long-running package operation
    And context without cancellation support
    When operation needs to be stopped
    Then operation cannot be cancelled
    And operation continues until completion
    And resources may be held unnecessarily
