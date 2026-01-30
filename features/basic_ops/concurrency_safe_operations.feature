@domain:basic_ops @m2 @REQ-API_BASIC-211 @spec(api_basic_operations.md#3362-safe-operations)
Feature: Safe operations for concurrent access

  @REQ-API_BASIC-211 @happy
  Scenario: Safe operations define what can be used concurrently
    Given concurrent goroutines operating on a package
    When concurrency guidance is applied
    Then safe operations are identified for concurrent access
    And safe operations avoid mutation or use synchronization
    And safe operations behave predictably under concurrent reads
    And safe operations align with documented thread safety modes
    And consumers can rely on the safe operation set for parallel workloads

