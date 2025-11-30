@domain:compression @m2 @REQ-COMPR-070 @REQ-COMPR-071 @REQ-COMPR-072 @REQ-COMPR-073 @REQ-COMPR-074 @REQ-COMPR-075 @spec(api_package_compression.md#132-intelligent-defaults-and-memory-management)
Feature: Memory strategy defaults and configuration

  @REQ-COMPR-071 @REQ-COMPR-072 @happy
  Scenario: Conservative strategy uses 25% of available RAM
    Given a system with 8GB available RAM
    When conservative memory strategy is selected
    Then maximum memory usage is set to 2GB (25%)
    And compression respects memory limit

  @REQ-COMPR-071 @REQ-COMPR-073 @happy
  Scenario: Balanced strategy uses 50% of available RAM as default
    Given a system with 16GB available RAM
    When balanced memory strategy is selected
    Then maximum memory usage is set to 8GB (50%)
    And this is the default strategy
    And compression respects memory limit

  @REQ-COMPR-071 @REQ-COMPR-074 @happy
  Scenario: Aggressive strategy uses 75% of available RAM
    Given a system with 32GB available RAM
    When aggressive memory strategy is selected
    Then maximum memory usage is set to 24GB (75%)
    And compression uses more memory for better performance

  @REQ-COMPR-071 @REQ-COMPR-075 @happy
  Scenario: Custom strategy allows custom memory configuration
    Given a system with available RAM
    When custom memory strategy is configured with MaxMemoryUsage
    Then specified memory limit is used
    And compression respects custom limit

  @REQ-COMPR-075 @happy
  Scenario: Custom strategy supports arbitrary memory limits
    Given compression configuration
    When custom strategy is set with MaxMemoryUsage of 512MB
    Then memory usage is limited to 512MB
    And compression adapts to available memory
