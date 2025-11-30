@domain:basic_ops @m2 @REQ-API_BASIC-086 @spec(api_basic_operations.md#931-use-context-for-resource-management)
Feature: Use Context for Resource Management

  @REQ-API_BASIC-086 @happy
  Scenario: Context is used for resource management
    Given a package operation that uses resources
    And a valid context
    When operation is performed with context
    Then context supports resource management
    And resources are tracked via context
    And context enables resource lifecycle control

  @REQ-API_BASIC-086 @happy
  Scenario: Context is used for resource cancellation
    Given a package operation that uses resources
    And a context with cancellation support
    When context cancellation is requested
    Then resources are released
    And resource cleanup is triggered
    And cancellation ensures proper resource management

  @REQ-API_BASIC-086 @happy
  Scenario: Context is passed to long-running operations
    Given a long-running package operation
    And a context for resource management
    When operation is started with context
    Then context is passed to operation
    And context supports resource tracking
    And context enables operation cancellation

  @REQ-API_BASIC-086 @happy
  Scenario: Context cancellation ensures proper resource cleanup
    Given a package operation with allocated resources
    And a context with cancellation capability
    When context is cancelled
    Then resources are cleaned up
    And resource cleanup is performed
    And operation termination ensures resource release

  @REQ-API_BASIC-086 @happy
  Scenario: Context enables operation termination for resource management
    Given a package operation in progress
    And a context for resource management
    When operation needs to be terminated
    Then context cancellation terminates operation
    And resources are released on termination
    And termination is clean and controlled

  @REQ-API_BASIC-086 @error
  Scenario: Operations without context cannot manage resources properly
    Given a package operation that uses resources
    And no context for resource management
    When operation is performed
    Then resource management is limited
    And resources may not be properly tracked
    And cancellation may not be possible
