@domain:basic_ops @m2 @REQ-API_BASIC-203 @spec(api_basic_operations.md#3342-resource-cleanup)
Feature: Resource cleanup

  @REQ-API_BASIC-203 @happy
  Scenario: Resource cleanup releases resources reliably
    Given a package that acquires resources during operations
    When the package is closed or an operation fails
    Then resource cleanup releases resources reliably
    And cleanup occurs for both success and failure paths
    And cleanup avoids leaking file handles or memory
    And cleanup is consistent with resource lifecycle rules
    And cleanup supports long-running processes without accumulation

