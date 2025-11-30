@domain:compression @m2 @REQ-COMPR-073 @spec(api_package_compression.md#13212-balanced-strategy-50-ram-default)
Feature: Balanced strategy uses 50% RAM as default

  @REQ-COMPR-073 @happy
  Scenario: Balanced strategy allocates 50% of RAM
    Given compression operations with MemoryStrategyBalanced
    When memory strategy is applied
    Then 50% of total RAM is allocated for compression
    And memory allocation is balanced
    And performance and system stability are balanced

  @REQ-COMPR-073 @happy
  Scenario: Balanced strategy is default for systems with 4-16GB RAM
    Given a system with 4-16GB RAM
    When compression operation is performed
    Then Balanced strategy is automatically selected as default
    And 50% of RAM is allocated
    And strategy provides optimal balance

  @REQ-COMPR-073 @happy
  Scenario: Balanced strategy provides good compression speed while leaving system responsive
    Given compression operations with Balanced strategy
    When compression is performed
    Then good compression speed is achieved
    And system remains responsive
    And balance between performance and stability is maintained

  @REQ-COMPR-073 @happy
  Scenario: Balanced strategy optimizes performance while maintaining system stability
    Given compression operations requiring balance
    When MemoryStrategyBalanced is used
    Then performance is optimized
    And system stability is maintained
    And balance meets typical use case requirements
