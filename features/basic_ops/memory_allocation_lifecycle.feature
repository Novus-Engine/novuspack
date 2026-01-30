@domain:basic_ops @m2 @REQ-API_BASIC-215 @spec(api_basic_operations.md#3372-memory-allocation-lifecycle)
Feature: Memory allocation lifecycle

  @REQ-API_BASIC-215 @happy
  Scenario: Memory allocation lifecycle defines allocation and deallocation patterns
    Given in-memory structures used by a package
    When memory is allocated during operations
    Then allocation lifecycle defines when allocations are created
    And deallocation or release occurs during cleanup and close
    And memory lifecycle avoids unbounded growth across operations
    And memory lifecycle supports large package handling
    And lifecycle aligns with documented memory management strategies

