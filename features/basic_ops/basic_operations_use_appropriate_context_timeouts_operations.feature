@domain:basic_ops @m2 @REQ-API_BASIC-081 @spec(api_basic_operations.md#913-use-appropriate-context-timeouts)
Feature: Basic Operations: Use Appropriate Context Timeouts

  @REQ-API_BASIC-081 @happy
  Scenario: Context timeouts are used for long-running operations
    Given a long-running package operation
    When context timeout is configured
    Then timeout prevents indefinite blocking
    And operation is cancelled if timeout is exceeded
    And timeout is appropriate for operation duration

  @REQ-API_BASIC-081 @happy
  Scenario: Context timeouts are set based on expected operation duration
    Given a package operation with known duration
    When context timeout is configured
    Then timeout matches expected operation duration
    And timeout provides sufficient time for operation
    And timeout prevents unnecessary delays

  @REQ-API_BASIC-081 @happy
  Scenario: Context timeouts handle timeout errors gracefully
    Given a package operation with context timeout
    And operation exceeds timeout duration
    When timeout error occurs
    Then timeout error is handled gracefully
    And error follows structured error format
    And error indicates context timeout

  @REQ-API_BASIC-081 @happy
  Scenario: Different operations use different timeout values
    Given multiple package operations with different durations
    When context timeouts are configured
    Then each operation has appropriate timeout
    And fast operations use shorter timeouts
    And slow operations use longer timeouts

  @REQ-API_BASIC-081 @error
  Scenario: Operations without timeouts can hang indefinitely
    Given a long-running package operation
    And context without timeout
    When operation takes excessive time
    Then operation can block indefinitely
    And resources may not be released
    And cancellation is not possible
