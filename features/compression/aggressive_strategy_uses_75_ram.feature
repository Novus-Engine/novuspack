@domain:compression @m2 @REQ-COMPR-074 @spec(api_package_compression.md#13213-aggressive-strategy-75-ram)
Feature: Aggressive strategy uses 75% RAM

  @REQ-COMPR-074 @happy
  Scenario: Aggressive strategy allocates 75% of RAM
    Given compression operations with MemoryStrategyAggressive
    When memory strategy is applied
    Then 75% of total RAM is allocated for compression
    And memory allocation is aggressive
    And maximum performance is prioritized

  @REQ-COMPR-074 @happy
  Scenario: Aggressive strategy is default for systems with more than 16GB RAM
    Given a system with more than 16GB RAM
    When compression operation is performed
    Then Aggressive strategy is automatically selected
    And 75% of RAM is allocated
    And strategy maximizes compression performance

  @REQ-COMPR-074 @happy
  Scenario: Aggressive strategy provides maximum performance for dedicated compression systems
    Given compression operations on dedicated compression system
    When MemoryStrategyAggressive is used
    Then maximum performance is achieved
    And compression speed is maximized
    And system is dedicated to compression tasks

  @REQ-COMPR-074 @happy
  Scenario: Aggressive strategy uses maximum available resources
    Given compression operations requiring maximum resources
    When Aggressive strategy is used
    Then maximum available resources are used
    And compression performance is maximized
    And resource utilization is aggressive
