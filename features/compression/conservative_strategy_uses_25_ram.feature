@domain:compression @m2 @REQ-COMPR-072 @spec(api_package_compression.md#13211-conservative-strategy-25-ram)
Feature: Conservative strategy uses 25% RAM

  @REQ-COMPR-072 @happy
  Scenario: Conservative strategy allocates 25% of RAM
    Given compression operations with MemoryStrategyConservative
    When memory strategy is applied
    Then 25% of total RAM is allocated for compression
    And memory allocation is conservative
    And system stability is prioritized

  @REQ-COMPR-072 @happy
  Scenario: Conservative strategy is default for systems with less than 4GB RAM
    Given a system with less than 4GB RAM
    When compression operation is performed
    Then Conservative strategy is automatically selected
    And 25% of RAM is allocated
    And strategy ensures system stability

  @REQ-COMPR-072 @happy
  Scenario: Conservative strategy is used when other processes need memory
    Given compression operations on system with limited RAM
    When other processes need memory
    Then Conservative strategy should be used
    And memory allocation leaves resources for other processes
    And system remains responsive

  @REQ-COMPR-072 @happy
  Scenario: Conservative strategy ensures system stability
    Given compression operations with Conservative strategy
    When compression is performed
    Then system stability is ensured
    And memory usage is limited
    And system remains responsive during compression
