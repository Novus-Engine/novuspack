@domain:basic_ops @m2 @REQ-API_BASIC-126 @spec(api_basic_operations.md#336-thread-safety-and-concurrency)
Feature: Thread safety and concurrency

  @REQ-API_BASIC-126 @happy
  Scenario: Thread safety and concurrency patterns define safe concurrent access
    Given concurrent goroutines operating on a package instance
    When concurrency rules are applied
    Then thread safety limitations are explicitly defined
    And operations safe for concurrent access are identified
    And locking and synchronization strategies are described
    And concurrent access avoids data races and inconsistent state
    And concurrency behavior supports predictable read and write operations

