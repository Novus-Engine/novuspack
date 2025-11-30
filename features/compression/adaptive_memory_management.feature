@domain:compression @m2 @REQ-COMPR-080 @spec(api_package_compression.md#1323-adaptive-memory-management)
Feature: Adaptive Memory Management

  @REQ-COMPR-080 @happy
  Scenario: Adaptive memory management monitors available memory during compression
    Given a compression operation in progress
    And adaptive memory management is enabled
    When compression runs
    Then available memory is continuously monitored
    And memory usage is tracked during operation
    And memory pressure is detected

  @REQ-COMPR-080 @happy
  Scenario: Adaptive memory management dynamically adjusts chunk size
    Given a compression operation in progress
    And adaptive chunking is enabled
    And memory pressure is detected
    When adaptive memory management responds
    Then chunk size is dynamically reduced
    And memory usage stays within limits
    And operation continues successfully

  @REQ-COMPR-080 @happy
  Scenario: Adaptive memory management increases chunk size when memory available
    Given a compression operation in progress
    And adaptive chunking is enabled
    And additional memory becomes available
    When adaptive memory management responds
    Then chunk size is increased for better performance
    And memory is utilized efficiently
    And compression performance improves

  @REQ-COMPR-080 @happy
  Scenario: Adaptive memory management prevents out of memory errors
    Given a compression operation with memory constraints
    And adaptive memory management is enabled
    When memory limits are approached
    Then chunk size is reduced automatically
    And disk buffering is enabled if needed
    And out of memory errors are prevented

  @REQ-COMPR-080 @happy
  Scenario: Adaptive memory management adjusts based on system load
    Given a compression operation
    And adaptive memory management is enabled
    And system load changes
    When memory management adapts
    Then chunk size adjusts to system conditions
    And memory usage adapts to available resources
    And operation remains stable
