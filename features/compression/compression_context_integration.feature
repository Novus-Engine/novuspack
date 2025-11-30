@domain:compression @m2 @REQ-COMPR-013 @spec(api_package_compression.md#02-context-integration)
Feature: Compression Context Integration

  @REQ-COMPR-013 @happy
  Scenario: All compression methods accept context.Context parameter
    Given a compression operation
    When compression methods are examined
    Then all methods accept context.Context as first parameter
    And context parameter is required for all compression operations
    And context follows standard Go context patterns

  @REQ-COMPR-013 @happy
  Scenario: Compression methods respect context cancellation
    Given a compression operation
    And a context with cancellation support
    When compression method is called
    Then method checks for context cancellation
    And operation can be cancelled via context
    And cancellation is respected during operation

  @REQ-COMPR-013 @happy
  Scenario: Compression methods respect context timeout
    Given a compression operation
    And a context with timeout
    When compression method is called
    Then method checks for context timeout
    And operation is terminated when timeout expires
    And timeout is respected during operation

  @REQ-COMPR-013 @error
  Scenario: Compression methods return error on context cancellation
    Given a compression operation
    And a cancelled context
    When compression method is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-COMPR-013 @error
  Scenario: Compression methods return error on context timeout
    Given a compression operation
    And a context that times out
    When compression method is called
    Then context timeout error is returned
    And error type is context timeout
    And error follows structured error format

  @REQ-COMPR-013 @happy
  Scenario: Context integration enables graceful operation termination
    Given a compression operation
    And a context with cancellation capability
    When context is cancelled during operation
    Then operation terminates gracefully
    And resources are cleaned up
    And no partial state is left
