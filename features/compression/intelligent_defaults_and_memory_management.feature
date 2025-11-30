@domain:compression @m2 @REQ-COMPR-070 @spec(api_package_compression.md#132-intelligent-defaults-and-memory-management)
Feature: Intelligent defaults and memory management

  @REQ-COMPR-070 @happy
  Scenario: System auto-detects optimal memory settings
    Given compression operations
    When system queries available system RAM
    Then optimal memory settings are automatically detected
    And memory limits are calculated based on system capabilities
    And automatic optimization is applied

  @REQ-COMPR-070 @happy
  Scenario: Memory strategy defaults are automatically selected
    Given compression operations on systems with different RAM
    When compression is performed
    Then Conservative strategy is selected for systems with <4GB RAM
    And Balanced strategy is selected for systems with 4-16GB RAM
    And Aggressive strategy is selected for systems with >16GB RAM

  @REQ-COMPR-070 @happy
  Scenario: Chunk size is automatically calculated
    Given compression operations without explicit chunk size
    When chunk size calculation is needed
    Then optimal chunk size is calculated as 25% of allocated memory limit
    And chunk size fits within memory constraints
    And multiple concurrent operations are supported

  @REQ-COMPR-070 @happy
  Scenario: Worker count is automatically detected
    Given compression operations requiring worker pool
    When worker count is needed
    Then number of available CPU cores is automatically detected
    And worker count is set accordingly
    And optimal parallel processing is enabled without overloading

  @REQ-COMPR-070 @happy
  Scenario: Adaptive memory management adjusts dynamically
    Given compression operations with adaptive management
    When compression is performed
    Then available memory is continuously monitored
    And chunk size is reduced if memory pressure detected
    And disk fallback occurs automatically when memory limits hit
